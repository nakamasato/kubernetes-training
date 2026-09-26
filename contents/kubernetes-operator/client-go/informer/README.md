# informer

## Overview

### Factory & Informers

![](informer-factory.drawio.svg)
### Single Informer

![](informer.drawio.svg)

***Informer*** monitors the changes of target resource. An informer is created for each of the target resources if you need to handle multiple resources (e.g. podInformer, deploymentInformer).

The snippets and diagrams below target client-go v0.37.1. Structs show implementation details rather than APIs to construct directly.

## types

### Interface [SharedInformerFactory](https://github.com/kubernetes/client-go/blob/v0.37.1/informers/factory.go)

```go
type SharedInformerFactory interface {
    internalinterfaces.SharedInformerFactory

    Start(stopCh <-chan struct{})

    StartWithContext(ctx context.Context)

    Shutdown()

    WaitForCacheSync(stopCh <-chan struct{}) map[reflect.Type]bool

    WaitForCacheSyncWithContext(ctx context.Context) cache.SyncResult

    ForResource(resource schema.GroupVersionResource) (GenericInformer, error)

    InformerFor(obj runtime.Object, newFunc internalinterfaces.NewInformerFunc) cache.SharedIndexInformer

    Admissionregistration() admissionregistration.Interface
    Internal() apiserverinternal.Interface
    Apps() apps.Interface
    Autoscaling() autoscaling.Interface
    Batch() batch.Interface
    Certificates() certificates.Interface
    Coordination() coordination.Interface
    Core() core.Interface
    Discovery() discovery.Interface
    Events() events.Interface
    Extensions() extensions.Interface
    Flowcontrol() flowcontrol.Interface
    Lifecycle() lifecycle.Interface
    Networking() networking.Interface
    Node() node.Interface
    Policy() policy.Interface
    Rbac() rbac.Interface
    Resource() resource.Interface
    Scheduling() scheduling.Interface
    Storage() storage.Interface
    Storagemigration() storagemigration.Interface
}
```

### Implementation [sharedInformerFactory](https://github.com/kubernetes/client-go/blob/v0.37.1/informers/factory.go)


```go
type sharedInformerFactory struct {
    client           kubernetes.Interface
    namespace        string
    tweakListOptions internalinterfaces.TweakListOptionsFunc
    lock             sync.Mutex
    defaultResync    time.Duration
    customResync     map[reflect.Type]time.Duration
    transform        cache.TransformFunc
    informerName     *cache.InformerName

    informers map[reflect.Type]cache.SharedIndexInformer
    startedInformers map[reflect.Type]bool
    wg sync.WaitGroup
    shuttingDown bool
}
```

Fields:

1. `client`: clientset to interact with API server
1. `namespace`: you can specify a namespace or all namespaces ([v1.NamespaceAll](https://github.com/kubernetes/client-go/blob/v0.37.1/informers/factory.go)) by default
1. `informers`: store created informers to start them when `factory.Start` is called.

The factory exposes API groups, each group exposes versions, and each version exposes resource-specific informer accessors. Accessor construction, shared-informer registration, and starting watches are separate steps.

How a new informer is created with a Factory:

1. Create a factory.

    ```go
    kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)
    ```

1. Obtain the Deployment informer accessor.

    ```go
    deploymentInformer := kubeInformerFactory.Apps().V1().Deployments()
    ```

    The v0.37.1 call chain is:

    | Call | Declared return type | Implementation |
    | --- | --- | --- |
    | [factory.Apps()](https://github.com/kubernetes/client-go/blob/v0.37.1/informers/factory.go#L379-L381) | `apps.Interface` | Calls `apps.New(f, f.namespace, f.tweakListOptions)` |
    | [apps.New(...)](https://github.com/kubernetes/client-go/blob/v0.37.1/informers/apps/interface.go#L45-L47) | `Interface` | Returns `&group{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}` |
    | [group.V1()](https://github.com/kubernetes/client-go/blob/v0.37.1/informers/apps/interface.go#L50-L52) | `v1.Interface` | Calls `v1.New(g.factory, g.namespace, g.tweakListOptions)` |
    | [v1.New(...)](https://github.com/kubernetes/client-go/blob/v0.37.1/informers/apps/v1/interface.go#L46-L48) | `Interface` | Returns `&version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}` |
    | [version.Deployments()](https://github.com/kubernetes/client-go/blob/v0.37.1/informers/apps/v1/interface.go#L61-L63) | `TypedDeploymentInformer` | Returns `&deploymentInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}` |

    In particular, Apps is a method on the concrete `sharedInformerFactory`; `kubeInformerFactory` is the sample variable, not a type or an upstream declaration. The method body is:

    ```go
    func (f *sharedInformerFactory) Apps() apps.Interface {
        return apps.New(f, f.namespace, f.tweakListOptions)
    }
    ```

    At this point the Deployment accessor exists, but it has not yet registered a shared informer or started API watches.

1. Obtain the shared informer and register handlers, usually while constructing your controller.

    The accessor's [Informer and TypedInformer methods](https://github.com/kubernetes/client-go/blob/v0.37.1/informers/apps/v1/deployment.go#L162-L168) are:

    ```go
    func (f *deploymentInformer) Informer() cache.SharedIndexInformer {
        return f.TypedInformer()
    }

    func (f *deploymentInformer) TypedInformer() DeploymentIndexInformer {
        return cache.NewTypedSharedIndexInformer[*apiappsv1.Deployment](f.factory.InformerFor(&apiappsv1.Deployment{}, f.defaultInformer))
    }
    ```

    [factory.InformerFor](https://github.com/kubernetes/client-go/blob/v0.37.1/informers/factory.go#L244-L268) locks the factory and looks up `reflect.TypeOf(obj)` in `f.informers`. If registered, it returns the existing informer. Otherwise, it selects the custom or default resync period, calls `newFunc(f.client, resyncPeriod)`, applies the factory transform if configured, and registers the result in `f.informers`.

    Here `newFunc` is [deploymentInformer.defaultInformer](https://github.com/kubernetes/client-go/blob/v0.37.1/informers/apps/v1/deployment.go#L158-L160):

    ```go
    func (f *deploymentInformer) defaultInformer(client kubernetes.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
        return NewTypedDeploymentInformerWithOptions(client, f.namespace, internalinterfaces.InformerOptions{ResyncPeriod: resyncPeriod, Indexers: cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, InformerName: f.factory.InformerName(), TweakListOptions: f.tweakListOptions})
    }
    ```

    [NewTypedDeploymentInformerWithOptions](https://github.com/kubernetes/client-go/blob/v0.37.1/informers/apps/v1/deployment.go#L118-L156) builds the Deployment ListWatch, wraps it with WatchList semantics, and calls `cache.NewSharedIndexInformerWithOptions`. The List/Watch callbacks use `client.AppsV1().Deployments(namespace)` and apply TweakListOptions. The resulting informer is wrapped with `cache.NewTypedSharedIndexInformer[*apiappsv1.Deployment]`. Construction sets up the objects; running them starts API requests.

    Register your handler and check the error:

    ```go
    _, err := deploymentInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
        AddFunc: handleAdd,
    })
    if err != nil {
        return err
    }
    ```

    `handleAdd` is your callback. Calling Lister also obtains the informer, because Lister needs its Indexer.

1. Start the registered informers.

    ```go
    kubeInformerFactory.StartWithContext(ctx)
    ```

    [StartWithContext](https://github.com/kubernetes/client-go/blob/v0.37.1/informers/factory.go#L167-L184) starts each registered informer that has not already started, using `informer.RunWithContext(ctx)`. Merely calling Apps().V1().Deployments() before Start is insufficient. If an informer is registered after Start, call StartWithContext again to start it.

    Wait for cache synchronization before reading through its Lister. Cancel the context before calling Shutdown to wait for informer goroutines. The channel-based Start used by the sample delegates to StartWithContext.

The next sections explain the shared informer's internals and RunWithContext lifecycle.

### Interface SharedInformer

- Interface:
    SharedInformer
    ```go
    // Selected methods; see the linked interface for all options.
    type SharedInformer interface {
        AddEventHandler(handler ResourceEventHandler) (ResourceEventHandlerRegistration, error)
        AddEventHandlerWithResyncPeriod(handler ResourceEventHandler, resyncPeriod time.Duration) (ResourceEventHandlerRegistration, error)
        AddEventHandlerWithOptions(handler ResourceEventHandler, options HandlerOptions) (ResourceEventHandlerRegistration, error)
        RemoveEventHandler(handle ResourceEventHandlerRegistration) error
        GetStore() Store
        GetController() Controller
        RunWithContext(ctx context.Context)
        HasSynced() bool
        LastSyncResourceVersion() string
        SetWatchErrorHandlerWithContext(handler WatchErrorHandlerWithContext) error
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

### Implementation sharedIndexInformer
```go
type sharedIndexInformer struct {
    indexer    Indexer
    controller Controller

    synced chan struct{}

    processor             *sharedProcessor
    cacheMutationDetector MutationDetector

    listerWatcher ListerWatcher

    objectType runtime.Object

    objectDescription string

    resyncCheckPeriod time.Duration
    defaultEventHandlerResyncPeriod time.Duration
    clock clock.Clock

    started, stopped bool
    startedLock      sync.Mutex

    blockDeltas sync.Mutex

    watchErrorHandler WatchErrorHandlerWithContext

    transform TransformFunc

    identifier InformerNameAndResource

    informerMetricsProvider InformerMetricsProvider

    keyFunc KeyFunc
}
```
`NewSharedIndexInformerWithOptions` initializes the Indexer, processor, mutation detector, and synchronization channels. `NewSharedIndexInformer` delegates to that constructor. The informer retains its ListerWatcher until `RunWithContext` creates the low-level controller and Reflector.

Components:
- [Indexer](../indexer)
- controller: explained below
- sharedProcessor: explained below
- [ListerWatcher](../listerwatcher/)


[sharedIndexInformer.Run](https://github.com/kubernetes/client-go/blob/v0.37.1/tools/cache/shared_informer.go):

1. Call `newQueueFIFO` to construct the informer queue. Its implementation depends on feature gates (including `InOrderInformers`); do not assume that every shared informer always uses DeltaFIFO. The queue accepts changes from the Reflector and supplies them to `Config.Process` or `Config.ProcessBatch`.
1. Create Controller with [New](https://github.com/kubernetes/client-go/blob/v0.37.1/tools/cache/controller.go)
1. Run [s.cacheMutationDetector.Run](https://github.com/kubernetes/client-go/blob/v0.37.1/tools/cache/mutation_detector.go)
1. Run `s.processor.run` <- start all listeners using a separate processor context, stopped after the low-level controller. listeners are added via `AddEventHandler`. (usually with `cache.ResourceEventHandlerFuncs{AddFunc: xx, UpdateFunc: xx, DeleteFunc: xx}`)
1. Run [s.controller.Run](https://github.com/kubernetes/client-go/blob/v0.37.1/tools/cache/controller.go) <- refer the controller section
        1. Create a new Reflector and call [r.Run](https://github.com/kubernetes/client-go/blob/v0.37.1/tools/cache/reflector.go) (ListAndWatch is called inside)

NewSharedInformer:

1. [NewSharedInformer](https://pkg.go.dev/k8s.io/client-go@v0.37.1/tools/cache#NewSharedInformer): call NewSharedIndexInformer with `Indexers{}`.
    ```go
    NewSharedIndexInformer(lw, exampleObject, defaultEventHandlerResyncPeriod, Indexers{})
    ```
1. [NewSharedIndexInformer](https://pkg.go.dev/k8s.io/client-go@v0.37.1/tools/cache#NewSharedIndexInformer)
    ```go
    func NewSharedIndexInformer(lw ListerWatcher, exampleObject runtime.Object, defaultEventHandlerResyncPeriod time.Duration, indexers Indexers) SharedIndexInformer {
        return NewSharedIndexInformerWithOptions(
            lw,
            exampleObject,
            SharedIndexInformerOptions{
                ResyncPeriod: defaultEventHandlerResyncPeriod,
                Indexers:     indexers,
            },
        )
    }
    ```

#### [sharedProcessor](https://github.com/kubernetes/client-go/blob/v0.37.1/tools/cache/shared_informer.go)

Role: hold a collection of listeners and distribute notification objects to them. `distribute` selects listeners; each listener buffers notifications and invokes the registered handler. Cache synchronization and completion of a handler's initial notifications are distinct; use the registration handle's `HasSynced` for the latter.

```go
type sharedProcessor struct {
    listenersStarted bool
    listenersLock    sync.RWMutex
    listenersRCond   *sync.Cond // Caller of Wait must hold a read lock on listenersLock.
    listeners map[*processorListener]bool
    clock     clock.Clock
    wg        wait.Group
}
```

1. `Listeners` are added for ResourceEventHandler via AddEventHandler
1. `distribute()` calls `listener.add` to propagate new events to each listener. `distribute()` is called by `informer.OnAdd`, `informer.OnUpdate`,  and `informer.OnDelete`
1. `run()` calls `listener.run` and `listener.pop` for all listeners.
`handler.OnAdd`, `handler.OnUpdate`, `handler.OnDelete` based on the notification type.

#### [Controller](https://github.com/kubernetes/client-go/blob/v0.37.1/tools/cache/controller.go)

Role: Run a reflector and enqueue item to Queue from ListerWatcher and process item from the queue with processfunc.

Interface:

```go
type Controller interface {
    RunWithContext(ctx context.Context)

    Run(stopCh <-chan struct{})

    HasSynced() bool

    HasSyncedChecker() DoneChecker

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

1. Most things are passed by `Config` (ListerWatcher, ObjectType, Queue (informer queue))

[Run](https://github.com/kubernetes/client-go/blob/v0.37.1/tools/cache/controller.go):

1. Create a Reflector with [NewReflectorWithOptions](https://github.com/kubernetes/client-go/blob/v0.37.1/tools/cache/reflector.go)
1. Run `reflector.RunWithContext` (details -> ref [reflector](../reflector))
    1. `ListAndWatchWithContext`
    1. `handleWatch`:
        1. event.Added -> store.Add
        1. event.Modified -> store.Update
        1. event.Deleted -> store.Delete (store = Queue)
1. Run [processLoop](https://github.com/kubernetes/client-go/blob/v0.37.1/tools/cache/controller.go) until the context is canceled.
    1. Pop item from the Queue and process it repeatedly. (Actual process is given by `Config.Process`, controller is just a container to execute `Process`)
        - `Config.Process`: [handleDeltas](https://github.com/kubernetes/client-go/blob/v0.37.1/tools/cache/shared_informer.go)
        `handleDeltas(logger, obj, isInInitialList)` calls [processDeltas(logger, s, s.indexer, deltas, isInInitialList, s.keyFunc)](https://github.com/kubernetes/client-go/blob/v0.37.1/tools/cache/controller.go)
            - `handler`: sharedIndexInformer
            - `clientState`: s.indexer
        - Keep indexer up-to-date by calling `indexer.Update()`, `indexer.Add()`, `indexer.Delete()`.
        - Distribute notification and add object to cacheMutationDetector by calling `sharedIndexInformer.OnUpdate()`, `sharedIndexInformer.OnAdd()`, `sharedIndexInformer.OnDelete()`

#### [MutationDetector](https://github.com/kubernetes/client-go/blob/v0.37.1/tools/cache/mutation_detector.go)

Role: Check if a cached object is mutated. Call failurefunc or panic if mutated.

1. By default, mutation detector is **not enabled**. (You can skip this components)
    ```go
    var mutationDetectionEnabled = false

    func init() {
        mutationDetectionEnabled, _ = strconv.ParseBool(os.Getenv("KUBE_CACHE_MUTATION_DETECTOR"))
    }
    ```
1. Run periodically calls [CompareObjects](https://github.com/kubernetes/client-go/blob/v0.37.1/tools/cache/mutation_detector.go).
1. CompareObjects compares `cached` and `copied` of `cacheObj` in `d.cachedObjs` and `d.retainedCachedObjs`.
    ```go
    type cacheObj struct {
        cached interface{}
        copied interface{}
    }
    ```
1. If any object is altered, call `failureFunc`. (if created with [NewCacheMutationDetector](https://github.com/kubernetes/client-go/blob/v0.37.1/tools/cache/mutation_detector.go), it doesn't have failureFunc, the program goes `panic`)
1. [AddObject](https://github.com/kubernetes/client-go/blob/v0.37.1/tools/cache/mutation_detector.go) adds an object to `d.addedObjs`.
1. You can enable the mutation detector to catch accidental mutation of shared cached objects. The updated sample calls `DeepCopy()` before changing labels, so it should not produce `panic: cache *v1.Pod modified`. Mutating the original object would trigger that failure.
    ```
    KUBE_CACHE_MUTATION_DETECTOR=true go run ./contents/kubernetes-operator/client-go/informer
    ```
## Example

1. Initialize clientset with `.kube/config`
1. Create an informer **factory** with the following line.
    ```go
    informerFactory := informers.NewSharedInformerFactory(kubeClient, time.Second*30)
    ```
    The second argument specifies ***ResyncPeriod***, which defines the interval of resync (*The resync operation consists of delivering to the handler an update notification for every object in the informer's local cache*). For more detail, please read [NewSharedInformer](https://pkg.go.dev/k8s.io/client-go@v0.37.1/tools/cache#NewSharedInformer)
1. Obtain a Pod informer accessor. Calling Informer() or Lister() registers its shared informer; starting the factory begins watching Pods.
    ```go
    podInformer := informerFactory.Core().V1().Pods()
    ```

    factory -> group -> version -> resource accessor (`TypedPodInformer`, which embeds `PodInformer`)

    ```go
    type PodInformer interface {
        Informer() cache.SharedIndexInformer
        Lister() v1.PodLister
    }
    ```

    1. `Informer()` returns `SharedIndexInformer`
        1. Delegate to `TypedInformer()`, which calls the factory's InformerFor and wraps the result with `cache.NewTypedSharedIndexInformer[*apicorev1.Pod]`.
        1. On first registration, `defaultInformer` constructs it through `NewTypedPodInformerWithOptions`, including the namespace index, informer name, and TweakListOptions.
        1. Reuse the registered shared informer on later calls. See the [Pod accessor implementation](https://github.com/kubernetes/client-go/blob/v0.37.1/informers/core/v1/pod.go).
    1. `Lister()` returns PodLister
        1. call `v1.NewPodLister(f.Informer().GetIndexer())`
        1. NewPodLister returns podLister with the given indexer.
            ```go
            type podLister struct {
                listers.ResourceIndexer[*corev1.Pod]
            }
            ```

1. Add event handlers (`AddFunc`, `UpdateFunc`, and `DeleteFunc`) to the pod informer.
    ```go
    _, err := podInformer.Informer().AddEventHandler(
        cache.ResourceEventHandlerFuncs{
            AddFunc:    handleAdd,
            UpdateFunc: handleUpdate,
            DeleteFunc: handleDelete,
        },
    )
    ```

    `handleAdd`, `handleUpdate`, and `handleDelete` define custom logic for each event. In this example, just print `"handleXXX is called"`

1. Create a signal-aware context and start the factory. Check the error returned by `AddEventHandler` first.
    ```go
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer stop()
    ch := ctx.Done()
    informerFactory.Start(ch)
    defer informerFactory.Shutdown()
    ```

1. Wait until the cache is synced.
    ```go
    cacheSynced := podInformer.Informer().HasSynced
    if ok := cache.WaitForCacheSync(ch, cacheSynced); !ok {
        log.Print("cache sync stopped")
        return
    }
    log.Println("cache is synced")
    ```

    [WaitForCacheSync](https://github.com/kubernetes/client-go/blob/v0.37.1/tools/cache/shared_informer.go)

    <details>

    ```go
    func WaitForCacheSync(stopCh <-chan struct{}, cacheSyncs ...InformerSynced) bool {
        err := wait.PollImmediateUntil(syncedPollPeriod,
            func() (bool, error) {
                for _, syncFunc := range cacheSyncs {
                    if !syncFunc() {
                        return false, nil
                    }
                }
                return true, nil
            },
            stopCh)
        if err != nil {
            return false
        }

        return true
    }
    ```

    The legacy stop-channel wrapper can be adapted to a context with [wait.ContextForChannel](https://github.com/kubernetes/apimachinery/blob/v0.37.1/pkg/util/wait/wait.go). New code can retain the context directly:

    ```go
    func ContextForChannel(parentCh <-chan struct{}) context.Context {
        return channelContext{stopCh: parentCh}
    }
    ```

    </details>

1. Wait for cancellation. The factory watches continuously without a separate polling loop.
    ```go
    <-ctx.Done()
    ```

## Run and check

Run from the repository root with Pod list/watch permissions. The logs below illustrate the event sequence; timestamps and Pod names depend on the cluster. The current sample also prints labels from a copied Pod. It does not modify the shared cache or API object. Delete keys use `DeletionHandlingMetaNamespaceKeyFunc`, which handles `DeletedFinalStateUnknown` tombstones.
1. Run
    ```
    go run ./contents/kubernetes-operator/client-go/informer
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
    ```
1. Create a `Pod` with name `nginx`.
    ```
    kubectl run nginx --image=nginx
    ```
1. Handlers are called by the events of the created `Pod`.
    ```
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
    2021/12/21 09:05:31 handleDelete is called for Pod (key: default/nginx)
    ```
1. Stop with Ctrl+C. Context cancellation stops watches, and `factory.Shutdown()` waits for informer goroutines.
1. Resync delivers Update notifications for cached objects every 30 seconds. It does not perform a fresh API List every 30 seconds.

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

## Tests

`go test ./contents/kubernetes-operator/client-go/informer` checks that handleAdd does not mutate a shared cached Pod and that deletion keys work for both objects and DeletedFinalStateUnknown tombstones. The live watch sequence above requires a cluster; these regression tests do not exercise API-server delivery.
