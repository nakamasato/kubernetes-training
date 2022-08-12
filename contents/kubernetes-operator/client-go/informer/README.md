# informer

## Overview

![](informer.drawio.svg)

***Informer*** monitors the changes of target resource. An informer is created for each of the target resources if you need to handle multiple resources (e.g. podInformer, deploymentInformer).


- Interface:
    SharedInformer
    ```go
    type SharedInformer interface {
        AddEventHandler(handler ResourceEventHandler)
        AddEventHandlerWithResyncPeriod(handler ResourceEventHandler, resyncPeriod time.Duration)
        GetStore() Store
        GetController() Controller
        Run(stopCh <-chan struct{})
        HasSynced() bool
        LastSyncResourceVersion() string
        SetWatchErrorHandler(handler WatchErrorHandler) error
        SetTransform(handler TransformFunc) error
    }
    ```
    SharedIndexInformer
    ```go
    type SharedIndexInformer interface {
        SharedInformer
        // AddIndexers add indexers to the informer before it starts.
        AddIndexers(indexers Indexers) error
        GetIndexer() Indexer
    }
    ```
- Implementation:
    ```go
    type sharedIndexInformer struct {
        indexer    Indexer
        controller Controller
        processor             *sharedProcessor
        cacheMutationDetector MutationDetector
        listerWatcher ListerWatcher
        objectType runtime.Object
        resyncCheckPeriod time.Duration
        defaultEventHandlerResyncPeriod time.Duration
        clock clock.Clock
        started, stopped bool
        startedLock      sync.Mutex
        blockDeltas sync.Mutex
        watchErrorHandler WatchErrorHandler
        transform TransformFunc
    }
    ```
    Components:
    - [Indexer](../indexer)
    - sharedProcessor
    - LiserWatcher


    [sharedIndexInformer.Run](https://github.com/kubernetes/client-go/blob/07171f82e7b66573fe0b161587ae23e47d598808/tools/cache/shared_informer.go#L397-L444):
    1. Create Delta Fifo by [NewDeltaFIFOWithOptions](https://github.com/kubernetes/client-go/blob/07171f82e7b66573fe0b161587ae23e47d598808/tools/cache/delta_fifo.go#L218)
    1. Create Controller with [New](https://github.com/kubernetes/client-go/blob/07171f82e7b66573fe0b161587ae23e47d598808/tools/cache/controller.go#L117)
    1. Run `s.cacheMutationDetector.Run`
    1. Run `s.processor.run` <- start all listeners. listners are added via `AddEventHandler`. (usually with `cache.ResourceEventHandlerFuncs{AddFunc: xx, UpdateFunc: xx, DeleteFunc: xx}`)
    1. Run `s.controller.Run(stopCh)` <- refer the controller section
- [NewSharedInformer](https://pkg.go.dev/k8s.io/client-go@v0.24.3/tools/cache#NewSharedInformer): call NewSharedIndexInformer with `Indexers{}`.
    ```go
    NewSharedIndexInformer(lw, exampleObject, defaultEventHandlerResyncPeriod, Indexers{})
    ```
- [NewSharedIndexInformer](https://pkg.go.dev/k8s.io/client-go@v0.24.3/tools/cache#NewSharedIndexInformer)
    ```go
    func NewSharedIndexInformer(lw ListerWatcher, exampleObject runtime.Object, defaultEventHandlerResyncPeriod time.Duration, indexers Indexers) SharedIndexInformer {
        realClock := &clock.RealClock{}
        sharedIndexInformer := &sharedIndexInformer{
            processor:                       &sharedProcessor{clock: realClock},
            indexer:                         NewIndexer(DeletionHandlingMetaNamespaceKeyFunc, indexers),
            listerWatcher:                   lw,
            objectType:                      exampleObject,
            resyncCheckPeriod:               defaultEventHandlerResyncPeriod,
            defaultEventHandlerResyncPeriod: defaultEventHandlerResyncPeriod,
            cacheMutationDetector:           NewCacheMutationDetector(fmt.Sprintf("%T", exampleObject)),
            clock:                           realClock,
        }
        return sharedIndexInformer
    }
    ```

- **sharedProcessor**

    ```go
    type sharedProcessor struct {
        listenersStarted bool
        listenersLock    sync.RWMutex
        listeners        []*processorListener
        syncingListeners []*processorListener
        clock            clock.Clock
        wg               wait.Group
    }
    ```

    1. Listeners are added for ResourceEventHandler via AddEventHandler
    1. Notification is added (addCh) to a listener by `handleDeletas(processDelta)` for each event from DeltaFifo.
    1. `run()` calls `handler.OnAdd`, `handler.OnUpdate`, `handler.OnDelete` based on the notification type.

- **Controller**

    Interface:

    ```go
    // Controller is a low-level controller that is parameterized by a
    // Config and used in sharedIndexInformer.
    type Controller interface {
        // Run does two things.  One is to construct and run a Reflector
        // to pump objects/notifications from the Config's ListerWatcher
        // to the Config's Queue and possibly invoke the occasional Resync
        // on that Queue.  The other is to repeatedly Pop from the Queue
        // and process with the Config's ProcessFunc.  Both of these
        // continue until `stopCh` is closed.
        Run(stopCh <-chan struct{})

        // HasSynced delegates to the Config's Queue
        HasSynced() bool

        // LastSyncResourceVersion delegates to the Reflector when there
        // is one, otherwise returns the empty string
        LastSyncResourceVersion() string
    }
    ```

    Implementation:
    ```go
    type controller struct {
        config         Config
        reflector      *Reflector
        reflectorMutex sync.RWMutex
        clock          clock.Clock
    }
    ```

    1. Most things are passed by `Config` (ListerWatcher, ObjectType, Queue (FifoDeltaQueue))

    [Run](https://github.com/kubernetes/client-go/blob/07171f82e7b66573fe0b161587ae23e47d598808/tools/cache/controller.go#L128):
    1. Create a Reflector with [NewReflector(lw ListerWatcher, expectedType interface{}, store Store, resyncPeriod time.Duration)](https://github.com/kubernetes/client-go/blob/07171f82e7b66573fe0b161587ae23e47d598808/tools/cache/reflector.go#L168)
    1. Run `reflector.Run` (details -> ref [reflector](../reflector))
        1. `ListAndWatch`
        1. `watchHandler` -> event.Added -> store.Add, event.Modified -> store.Update, event.Deleted -> store.Delete (store = Queue)
    1. Run [processLoop](https://github.com/kubernetes/client-go/blob/07171f82e7b66573fe0b161587ae23e47d598808/tools/cache/controller.go#L182) every second.
        1. Pop item from the Queue and process it repeatedly. (Actual process is given by `Config.Process`, controller is just a container to execute `Process`)
            - `Config.Process`: [HandleDeltas](https://github.com/kubernetes/client-go/blob/master/tools/cache/shared_informer.go#L566)
            `HandleDeltas` calls [processDeltas(s, s.indexer, s.transform, deltas)](https://github.com/kubernetes/client-go/blob/07171f82e7b66573fe0b161587ae23e47d598808/tools/cache/controller.go#L410)
                - `handler`: sharedIndexInformer
                - `clientState`: s.indexer
            - Keep indexer up-to-date by calling `indexer.Update()`, `indexer.Add()`, `indexer.Delete()`.
            - Distribute notification and add object to cacheMutationDetector by calling `sharedIndexInformer.OnUpdate()`, `sharedIndexInformer.OnAdd()`, `sharedIndexInformer.OnDelete()`


## Example

1. Initialize clientset with `.kube/config`
1. Create an informer **factory** with the following line.
    ```go
    informerFactory := informers.NewSharedInformerFactory(kubeClient, time.Second*30)
    ```
    The second argument specifies ***ResyncPeriod***, which defines the interval of resync (*The resync operation consists of delivering to the handler an update notification for every object in the informer's local cache*). For more detail, please read [NewSharedInformer](https://pkg.go.dev/k8s.io/client-go@v0.23.1/tools/cache#NewSharedInformer)
1. Create an informer for Pods, which watches Pod's changes.
    ```go
    podInformer := informerFactory.Core().V1().Pods()
    ```

    factory -> group -> version -> kind

    ```go
    type PodInformer interface {
        Informer() cache.SharedIndexInformer
        Lister() v1.PodLister
    }
    ```

    1. `Informer()` returns `SharedIndexInformer`
        1. call `f.factory.InformerFor(&corev1.Pod{}, f.defaultInformer)`
        1. create new informer with [NewFilteredPodInformer](https://github.com/kubernetes/client-go/blob/ee1a5aaf793a9ace9c433f5fb26a19058ed5f37c/informers/core/v1/pod.go#L58) if not exist
        1. return the informer
    1. `Lister()` returns PodLister
        1. call `v1.NewPodLister(f.Informer().GetIndexer())`
        1. NewPodLister returns podLister with the given indexer.
            ```go
            type podLister struct {
                indexer cache.Indexer
            }
            ```

1. Add event handlers (`AddFunc`, `UpdateFunc`, and `DeleteFunc`) to the pod informer.
    ```go
    podInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    handleAdd,
			UpdateFunc: handleUpdate,
			DeleteFunc: handleDelete,
		},
	)
    ```

    `handleAdd`, `handleUpdate`, and `handleDelete` define custom logic for each event. In this example, just print `"handleXXX is called"`
1. Create a stop channel and start the factory.
    ```go
    ch := make(chan struct{}) // stop channel
	informerFactory.Start(ch)
    ```
1. Wait until the cache is synced.
    ```go
    cacheSynced := podInformer.Informer().HasSynced
	if ok := cache.WaitForCacheSync(ch, cacheSynced); !ok {
		log.Printf("cache is not synced")
	}
	log.Println("cache is synced")
    ```
1. Run `run` function every 10 seconds
    ```go
    go wait.Until(run, time.Second*10, ch)
	<-ch
    ```

## Run and check
1. Run
    ```
    go run informer.go
    ```

1. All Pods are synced in the cache.

    ```
    2021/12/21 09:05:08 handleAdd is called for Pod (key: local-path-storage/local-path-provisioner-547f784dff-lhwfk)
    2021/12/21 09:05:08 handleAdd is called for Pod (key: kube-system/kube-scheduler-kind-control-plane)
    2021/12/21 09:05:08 handleAdd is called for Pod (key: kube-system/etcd-kind-control-plane)
    2021/12/21 09:05:08 handleAdd is called for Pod (key: kube-system/kube-apiserver-kind-control-plane)
    2021/12/21 09:05:08 handleAdd is called for Pod (key: kube-system/kindnet-nzc7p)
    2021/12/21 09:05:08 handleAdd is called for Pod (key: kube-system/coredns-558bd4d5db-b4wjg)
    2021/12/21 09:05:08 handleAdd is called for Pod (key: kube-system/kube-controller-manager-kind-control-plane)
    2021/12/21 09:05:08 handleAdd is called for Pod (key: kube-system/kube-proxy-vrcbc)
    2021/12/21 09:05:08 handleAdd is called for Pod (key: kube-system/coredns-558bd4d5db-8q78s)
    2021/12/21 09:05:08 handleAdd is called for Pod (key: default/foo-sample-688594b488-782kw)
    2021/12/21 09:05:08 cache is synced
    2021/12/21 09:05:08 run
    ```
1. Create a `Pod` with name `nginx`.
    ```
    kubectl run nginx --image=nginx
    ```
1. Handlers are called by the events of the created `Pod`.
    ```
    2021/12/21 09:05:18 run
    2021/12/21 09:05:20 handleAdd is called for Pod (key: default/nginx)
    2021/12/21 09:05:20 handleUpdate is called for Pod (key: default/nginx)
    2021/12/21 09:05:20 handleUpdate is called for Pod (key: default/nginx)
    ```
1. Delete the `Pod`
    ```
    kubectl delete po nginx
    ```
1. Handlers are called by the events of the Pod deletion.
    ```
    2021/12/21 09:05:29 handleUpdate is called for Pod (key: default/nginx)
    2021/12/21 09:05:30 handleUpdate is called for Pod (key: default/nginx)
    2021/12/21 09:05:31 handleUpdate is called for Pod (key: default/nginx)
    2021/12/21 09:05:31 handleUpdate is called for Pod (key: default/nginx)
    2021/12/21 09:05:31 handlDelete is called for Pod (key: default/nginx)
    ```
1. `run` function is called every 10 seconds.
    ```
    2021/12/21 09:26:08 run
    2021/12/21 09:26:18 run
    2021/12/21 09:26:28 run
    ```
1. The cached is resynced every 30 seconds.

    ```
    2021/12/21 09:27:08 handleUpdate is called for Pod (key: local-path-storage/local-path-provisioner-547f784dff-lhwfk)
    2021/12/21 09:27:08 handleUpdate is called for Pod (key: kube-system/kube-apiserver-kind-control-plane)
    2021/12/21 09:27:08 handleUpdate is called for Pod (key: kube-system/coredns-558bd4d5db-b4wjg)
    2021/12/21 09:27:08 handleUpdate is called for Pod (key: kube-system/kube-controller-manager-kind-control-plane)
    2021/12/21 09:27:08 handleUpdate is called for Pod (key: kube-system/coredns-558bd4d5db-8q78s)
    2021/12/21 09:27:08 handleUpdate is called for Pod (key: default/foo-sample-688594b488-782kw)
    2021/12/21 09:27:08 handleUpdate is called for Pod (key: kube-system/kube-scheduler-kind-control-plane)
    2021/12/21 09:27:08 handleUpdate is called for Pod (key: kube-system/etcd-kind-control-plane)
    2021/12/21 09:27:08 handleUpdate is called for Pod (key: kube-system/kindnet-nzc7p)
    2021/12/21 09:27:08 handleUpdate is called for Pod (key: kube-system/kube-proxy-vrcbc)
    ```

# reference
- https://adevjoe.com/post/client-go-informer/
- https://www.huweihuang.com/kubernetes-notes/code-analysis/kube-controller-manager/sharedIndexInformer.html
- https://yangxikun.com/kubernetes/2020/03/05/informer-lister.html