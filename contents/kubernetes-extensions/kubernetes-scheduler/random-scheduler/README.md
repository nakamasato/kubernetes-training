# Kubernetes scheduler

## Components:
- `podQueue` channel
- `quit` channel
- `Scheduler`
    - ScheduleOne:
        1. Get a Pod from `podQueue`.
        1. findNode
            1. Get Nodes from lister.
            1. Filter out unschedulable nodes.
            1. Give a score for each node.
            1. Get the node with the highest score.
        1. bindNode
        1. emitEvent

## Steps:

1. Create `main`, `NewScheduler` and `Scheduler` struct with `Run` method.
    ```go
    package main

    import (
        "log"
        "k8s.io/api/core/v1"
    )

    func main()  {
        log.Println("Start a scheduler")

        podQueue := make(chan *v1.Pod, 300)
        defer close(podQueue)

        quit := make(chan struct{})
        defer close(quit)

        scheduler := NewScheduler(podQueue, quit)
        scheduler.Run()
    }

    type Scheduler struct {
    }

    func (s *Scheduler) Run() {

    }

    func NewScheduler(podQueue chan *v1.Pod, quit chan struct{}) Scheduler {
        return Scheduler{}
    }
    ```

    You can try running the empty scheduler:

    ```
    go run main.go
    ```
1. Update `Scheduler` struct (Prepare nodeInformer, podInformer, randomPredicate, and randomPriority):
    1. Update definitions.
        ```go
        type predicateFunc func(node *v1.Node, pod *v1.Pod) bool
        type priorityFunc func(node *v1.Node, pod *v1.Pod) int

        type Scheduler struct {
            clientset  *kubernetes.Clientset
            podQueue   chan *v1.Pod
            nodeLister listersv1.NodeLister
            predicates []predicateFunc
            priorities []priorityFunc
        }
        ```
    1. Initialize `clientset` in `NewScheduler`.
    1. Define `randomPredicate` and `randomPriority`.
        ```go
        func randomPredicate(node *v1.Node, pod *v1.Pod) bool {
            r := rand.Intn(2)
            return r == 0
        }

        func randomPriority(node *v1.Node, pod *v1.Pod) int {
            return rand.Intn(100)
        }
        ```
    1. Define `initInformers`.

        1. Create shared informer factory
        1. Create node informer with event handler for `AddFunc` (just print).
        1. Create pod informer with event handler for `AddFunc`.
            1. Check `NodeName == ""`: Unscheduled Pods
            1. Check `SchedulerName == scheduclerName`; this scheduler is specified.
        1. Start the factory.
        1. Return the nodeInformer lister.

        ```go
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
                    AddFunc: func(obj interface{}){
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
        ```

    1. Update `NewScheduler`:
        ```go
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
                clientset: clientset,
                podQueue: podQueue,
                nodeLister: initInformers(clientset, podQueue, quit),
                predicates: []predicateFunc{
                    randomPredicate,
                },
                priorities: []priorityFunc{
                    randomPriority,
                },
            }
        }
        ```

    1. Try running the scheduler at this point. (Nothing happens as `Run` is still empty.)
        ```
        go run main.go
        Start a scheduler
        Run is called
        ```
1. Find best node for a pod in `ScheduleOne` function.
    1. Define `ScheduleOne`.
        The role of `ScheduleOne`:
        1. Get a pod from `podQueue`.
        1. Get the fit node from `findNode`.

        ```go
        func (s *Scheduler) ScheduleOne() {
            p := <- s.podQueue
            log.Println("found a pod to schedule:", p.Namespace, "/", p.Name)

            node, err := s.findNode(p)
            if err != nil {
                log.Println("cannot find node that fits pod", err.Error())
                return
            }
            log.Printf("node %s is chosen for Pod [%s/%s]\n", node, p.Namespace, p.Name)
        }
    1. Define `findNode` (find the best node for a given pod. If no schedulable node, return error.)
        The role of `findNode`:
        1. Get nodes from the node lister.
        1. Return error if there's no node.
        1. Give a score with `prioritize` for each node.
        1. Return the node with highest score `findBestNode`.

        ```go
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
        ```

    1. Define `prioritize`: Give a score with `priorities` for each node.
        ```go
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
        ```
    1. `findBestNode`: Get the node with the highest score.
        ```go
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
        ```
    1. Update `Run` to call `ScheduleOne`.
        ```go
        func (s *Scheduler) Run(quit chan struct{}) {
	        log.Println("Run is called")
	        wait.Until(s.ScheduleOne, 0, quit)
        }
        ```

    1. Run and check the scheduler (Just choose a node):

        Run the scheduler:
        ```bash
        go run main.go
        2021/12/26 17:21:06 Start a scheduler
        2021/12/26 17:21:06 Run is called
        2021/12/26 17:21:06 New node is added. kind-control-plane
        ```

        Create a pod with `schedulerName: random-scheduler`:
        ```
        kubectl apply -f pod.yaml
        ```

        Scheduler's log:
        ```bash
        2021/12/26 17:21:06 found a pod to schedule: [default/nginx]
        2021/12/26 17:21:06 calculated priorities: map[kind-control-plane:47]
        2021/12/26 17:21:06 node kind-control-plane is chosen for Pod [default/nginx]
        ```

        ```
        kubectl get pod nginx
        NAME    READY   STATUS    RESTARTS   AGE
        nginx   0/1     Pending   0          3m55s
        ```
1. Bind Node to Pod.
    1. Add the following lines to `ScheduleOne`.

        ```go
            err = s.bindPod(p, node)
            if err != nil {
                log.Println("failed to bind pod", err.Error())
                return
            }
        ```
    1. Define `bindPod`.

        ```go
        import (
            "context"
            ...
            metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
            ...
        )
        ...
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
        ```
    1. Run the scheduler.
        1. Run the scheduler
            ```bash
            go run main.go
            2021/12/26 17:33:35 Start a scheduler
            2021/12/26 17:33:35 Run is called
            2021/12/26 17:33:35 New node is added. kind-control-plane
            ```
        1. Create a Pod.
            ```
            kubectl apply -f pod.yaml
            ```
        1. Check logs.
            ```bash
            2021/12/26 17:35:35 found a pod to schedule: [default/nginx]
            2021/12/26 17:35:35 calculated priorities: map[kind-control-plane:47]
            2021/12/26 17:35:35 node kind-control-plane is chosen for Pod [default/nginx]
            2021/12/26 17:35:35 pod [default/nginx] is successfully scheduled to node kind-control-plane
            ```
        1. Check pod.
            ```
            kubectl get pod nginx
            NAME    READY   STATUS    RESTARTS   AGE
            nginx   1/1     Running   0          26s
            ```
1. Emit event.
    1. Add `emitEvent` function.
    ```go
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
    ```

    1. Add the following lines to the `ScheduleOne`.

        ```go
            message := fmt.Sprintf("pod [%s/%s] is successfully scheduled to node %s", p.Namespace, p.Name, node)
            log.Println(message)

            err = s.emitEvent(p, message)
            if err != nil {
                log.Println("failed to emit scheduled event", err.Error())
                return
            }
        ```

    1. Check.
        1. Create a Pod.
            ```
            kubectl apply -f pod.yaml
            ```
        1. Run the scheduler.
            ```
            go run main.go
            2021/12/26 17:41:24 Start a scheduler
            2021/12/26 17:41:24 Run is called
            2021/12/26 17:41:24 New node is added. kind-control-plane
            2021/12/26 17:41:43 found a pod to schedule: [default/nginx]
            2021/12/26 17:41:43 calculated priorities: map[kind-control-plane:47]
            2021/12/26 17:41:43 node kind-control-plane is chosen for Pod [default/nginx]
            2021/12/26 17:41:43 pod [default/nginx] is successfully scheduled to node kind-control-plane
            ```
        1. Check event.

            ```
            kubectl get event | grep Scheduled
            68s         Normal   Scheduled   pod/nginx   pod [default/nginx] is successfully scheduled to node kind-control-plane
            ```

# Reference
- https://banzaicloud.com/blog/k8s-custom-scheduler/
- https://github.com/martonsereg/random-scheduler
