# kubernetes scheduler

## Components:
- `podQueue` channel
- `quit` channel
- `Scheduler`

## Steps:

1. Create `main`, `NewScheduler` and `Scheduler` struct with `Run` method.
    ```go
    package main

    import (
        "fmt"
        "k8s.io/api/core/v1"
    )


    func main()  {
        fmt.Println("Start a scheduler")

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
