package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"path/filepath"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	listersv1 "k8s.io/client-go/listers/core/v1"

	// "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const schedulerName = "random-scheduler"

func main() {
	log.Println("Start a scheduler")

	podQueue := make(chan *v1.Pod, 300)
	defer close(podQueue)

	quit := make(chan struct{})
	defer close(quit)

	scheduler := NewScheduler(podQueue, quit)
	scheduler.Run(quit)
}

type predicateFunc func(node *v1.Node, pod *v1.Pod) bool
type priorityFunc func(node *v1.Node, pod *v1.Pod) int

type Scheduler struct {
	clientset  *kubernetes.Clientset
	podQueue   chan *v1.Pod
	nodeLister listersv1.NodeLister
	predicates []predicateFunc
	priorities []priorityFunc
}

func (s *Scheduler) Run(quit chan struct{}) {
	log.Println("Run is called")
	wait.Until(s.ScheduleOne, 0, quit)
}

func (s *Scheduler) ScheduleOne() {
	p := <-s.podQueue
	log.Printf("found a pod to schedule: [%s/%s]\n", p.Namespace, p.Name)

	node, err := s.findNode(p)
	if err != nil {
		log.Println("cannot find node that fits pod", err.Error())
		return
	}
	log.Printf("node %s is chosen for Pod [%s/%s]\n", node, p.Namespace, p.Name)

	err = s.bindPod(p, node)
	if err != nil {
		log.Println("failed to bind pod", err.Error())
		return
	}
	message := fmt.Sprintf("pod [%s/%s] is successfully scheduled to node %s", p.Namespace, p.Name, node)
	log.Println(message)

	err = s.emitEvent(p, message)
	if err != nil {
		log.Println("failed to emit scheduled event", err.Error())
		return
	}
}

func (s *Scheduler) findNode(pod *v1.Pod) (string, error) {
	nodes, err := s.nodeLister.List(labels.Everything())
	if err != nil {
		return "", err
	}
	if len(nodes) == 0 {
		return "", errors.New("failed to find schedulable nodes")
	}
	priorities := s.prioritize(nodes, pod)
	return s.findBestNode(priorities), nil
}

func (s *Scheduler) prioritize(nodes []*v1.Node, pod *v1.Pod) map[string]int {
	priorities := make(map[string]int)
	for _, node := range nodes {
		for _, priority := range s.priorities {
			priorities[node.Name] += priority(node, pod)
		}
	}
	log.Println("calculated priorities:", priorities)
	return priorities
}

func (s *Scheduler) findBestNode(priorities map[string]int) string {
	var maxP int
	var bestNode string
	for node, p := range priorities {
		if p > maxP {
			maxP = p
			bestNode = node
		}
	}
	return bestNode
}

func (s *Scheduler) bindPod(pod *v1.Pod, node string) error {
	return s.clientset.CoreV1().Pods(pod.Namespace).Bind(
		context.Background(),
		&v1.Binding{
			ObjectMeta: metav1.ObjectMeta{Name: pod.Name, Namespace: pod.Namespace},
			Target:     v1.ObjectReference{APIVersion: "v1", Kind: "Node", Name: node},
		},
		metav1.CreateOptions{},
	)
}

func (s *Scheduler) emitEvent(p *v1.Pod, message string) error {
	timestamp := time.Now().UTC()
	_, err := s.clientset.CoreV1().Events(p.Namespace).Create(
		context.Background(),
		&v1.Event{
			Count:          1,
			Message:        message,
			Reason:         "Scheduled",
			LastTimestamp:  metav1.NewTime(timestamp),
			FirstTimestamp: metav1.NewTime(timestamp),
			Type:           "Normal",
			Source: v1.EventSource{
				Component: schedulerName,
			},
			InvolvedObject: v1.ObjectReference{
				Kind:      "Pod",
				Name:      p.Name,
				Namespace: p.Namespace,
				UID:       p.UID,
			},
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: p.Name + "-",
			},
		},
		metav1.CreateOptions{},
	)
	if err != nil {
		return err
	}
	return nil
}

func NewScheduler(podQueue chan *v1.Pod, quit chan struct{}) Scheduler {
	// In-Cluster configuration
	// config, err := rest.InClusterConfig()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	kubeConfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		log.Printf("Building config from flags, %s", err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	return Scheduler{
		clientset:  clientset,
		podQueue:   podQueue,
		nodeLister: initInformers(clientset, podQueue, quit),
		predicates: []predicateFunc{
			randomPredicate,
		},
		priorities: []priorityFunc{
			randomPriority,
		},
	}
}

func randomPredicate(node *v1.Node, pod *v1.Pod) bool {
	r := rand.Intn(2)
	return r == 0
}

func randomPriority(node *v1.Node, pod *v1.Pod) int {
	return rand.Intn(100)
}

func initInformers(clientset *kubernetes.Clientset, podQueue chan *v1.Pod, quit chan struct{}) listersv1.NodeLister {
	factory := informers.NewSharedInformerFactory(clientset, 0)
	nodeInformer := factory.Core().V1().Nodes()
	nodeInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				node, ok := obj.(*v1.Node)
				if !ok {
					log.Println("this is not a node")
					return
				}
				log.Printf("New node is added. %s\n", node.GetName())
			},
		},
	)

	podInformer := factory.Core().V1().Pods()
	podInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				pod, ok := obj.(*v1.Pod)
				if !ok {
					log.Println("This is not a pod")
					return
				}
				if pod.Spec.NodeName == "" && pod.Spec.SchedulerName == schedulerName {
					podQueue <- pod
				}
			},
		},
	)
	factory.Start(quit)
	return nodeInformer.Lister()
}
