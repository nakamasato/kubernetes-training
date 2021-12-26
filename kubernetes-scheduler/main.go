package main

import (
	"fmt"
	"log"
	"math/rand"
	"path/filepath"

	"k8s.io/api/core/v1"
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
	fmt.Println("Start a scheduler")

	podQueue := make(chan *v1.Pod, 300)
	defer close(podQueue)

	quit := make(chan struct{})
	defer close(quit)

	scheduler := NewScheduler(podQueue, quit)
	scheduler.Run()
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

func (s *Scheduler) Run() {
	fmt.Println("Run is called")
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
