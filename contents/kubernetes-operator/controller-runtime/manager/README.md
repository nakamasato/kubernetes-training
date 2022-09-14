# [Manager](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/manager)

## types

### 1. Manager Interface

```go
type Manager interface {
	cluster.Cluster
	Add(Runnable) error
	Elected() <-chan struct{}
	AddMetricsExtraHandler(path string, handler http.Handler) error
	AddHealthzCheck(name string, check healthz.Checker) error
	AddReadyzCheck(name string, check healthz.Checker) error
	Start(ctx context.Context) error
	GetWebhookServer() *webhook.Server
	GetLogger() logr.Logger
	GetControllerOptions() v1alpha1.ControllerConfigurationSpec
}
```

### 2. [controllerManager](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.13.0/pkg/manager/internal.go#L66-L173)

```go
type controllerManager struct {
	sync.Mutex
	started bool

	stopProcedureEngaged *int64
	errChan              chan error
	runnables            *runnables

	// cluster holds a variety of methods to interact with a cluster. Required.
	cluster cluster.Cluster

    ...
}
```

### 3. [Runnable](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.13.0/pkg/manager/manager.go#L293-L298) interface

```go
type Runnable interface {
	Start(context.Context) error
}
```

### 4. [runnables](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.13.0/pkg/manager/runnable_group.go#L37-L45)

```go
type runnables struct {
	Webhooks       *runnableGroup
	Caches         *runnableGroup
	LeaderElection *runnableGroup
	Others         *runnableGroup
}

type runnableGroup struct {
	ctx    context.Context
	cancel context.CancelFunc

	start        sync.Mutex
	startOnce    sync.Once
	started      bool
	startQueue   []*readyRunnable
	startReadyCh chan *readyRunnable

	stop     sync.RWMutex
	stopOnce sync.Once
	stopped  bool

	errChan chan error
	ch chan *readyRunnable
	wg *sync.WaitGroup
}
```

## How `Manager` is initialized by [New](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.13.0/pkg/manager/manager.go#L336)

### 1. Set default values for Options fields wiht [setOptionsDefaults](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.13.0/pkg/manager/manager.go#L571)

```go
manager, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{})
```

```go
options = setOptionsDefaults(options)
```

|name|value|where is the option used|
|---|---|---|
|newResourceLock|leaderelection.NewResourceLock|`setLeaderElectionConfig`|
|newRecorderProvider|intrec.NewProvider|New to create recorderProvider|
|EventBroadcaster|`func() (record.EventBroadcaster, bool) {return record.NewBroadcaster(), true}`|as an option for `newRecorderProvider` in New|
|newMetricsListener|metrics.NewListener|New to create metricsListener|
|LeaseDuration|*defaultLeaseDuration|`setLeaderElectionConfig`|
|RenewDeadline|*defaultRenewDeadline|`setLeaderElectionConfig`|
|RetryPeriod|*defaultRetryPeriod|`setLeaderElectionConfig`|
|ReadinessEndpointName|defaultReadinessEndpoint||
|LivenessEndpointName|defaultLivenessEndpoint||
|newHealthProbeListener|defaultHealthProbeListener|
|GracefulShutdownTimeout|*defaultGracefulShutdownPeriod||
|Logger|log.Log||
|BaseContext|defaultBaseContext||

※ **New** in the table means [Manager.New](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.13.0/pkg/manager/manager.go#L336)


### 2. Initialize a [controllerManager](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.13.0/pkg/manager/internal.go#L66)

1. Initialize **Cluster**, which provides methods to interact with Kubernetes cluster
    ```go
    cluster, err := cluster.New(config, func(clusterOptions *cluster.Options) {
        clusterOptions.Scheme = options.Scheme
        clusterOptions.MapperProvider = options.MapperProvider
        clusterOptions.Logger = options.Logger
        clusterOptions.SyncPeriod = options.SyncPeriod
        clusterOptions.Namespace = options.Namespace
        clusterOptions.NewCache = options.NewCache
        clusterOptions.NewClient = options.NewClient
        clusterOptions.ClientDisableCacheFor = options.ClientDisableCacheFor
        clusterOptions.DryRunClient = options.DryRunClient
        clusterOptions.EventBroadcaster = options.EventBroadcaster //nolint:staticcheck
    })
    ```
    For more details, please check [cluster](../cluster/README.md).
    `Cluster` is also a runnable.
1. Initialize other necessary things like `recordProvider`, `runnables`, etc.

1. Initialize `controllerManager`

    ```go
    &controllerManager{
        ...
		cluster:                       cluster,
		runnables:                     runnables,
        ...
		recorderProvider:              recorderProvider,
	}
    ```
### 3. Bind a Controller to the Manager

Bind a Controller to the Manager using [NewControllerManagedBy](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.13.0/alias.go#L101)(alias for [builder.ControllerManagedBy](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.13.0/pkg/builder/controller.go#L66)).

```go
err = ctrl.
    NewControllerManagedBy(manager). // Create the Controller
    For(&appsv1.ReplicaSet{}).       // ReplicaSet is the Application API
    Owns(&corev1.Pod{}).             // ReplicaSet owns Pods created by it
    Complete(&ReplicaSetReconciler{Client: manager.GetClient()})
```

Internally, [builder.Build](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.13.0/pkg/builder/controller.go#L175) create a new controller and add it to `manager.runnables.Others` by `Manager.Add(Runnable)`.

You can also check [Builder](../builder) and [Internal process of adding a Controller to a Manager](#internal-process-of-adding-a-controller-to-a-manager)

### 4. [controllerManager.Start()](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.13.0/pkg/manager/internal.go#L399) calls `runnables.xxx.Start()` to start all runnables.
```go
err := cm.runnables.Webhooks.Start(cm.internalCtx);
...

err := cm.runnables.Caches.Start(cm.internalCtx);
...

err := cm.runnables.Others.Start(cm.internalCtx);
...
```
1. Controller will be in `runnables.Others` and you can check the actual `Start` logic in [controller](../controller).

## Internal process of adding a `Controller` to a `Manager`

1. `Manager.Add(Runnable)`: gets lock and calls `add(runnable)`.
    1. `cm.SetFields(r)`
        ```go
        if err := cm.cluster.SetFields(i); err != nil {
            return err
        }
        if _, err := inject.InjectorInto(cm.SetFields, i); err != nil {
            return err
        }
        if _, err := inject.StopChannelInto(cm.internalProceduresStop, i); err != nil {
            return err
        }
        if _, err := inject.LoggerInto(cm.logger, i); err != nil {
            return err
        }
        ```
        1. [cluster.SetFields](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.13.0/pkg/cluster/internal.go#L67) set dependencies on the object that implements the [inject](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.3/pkg/runtime/inject) interface. Specifically set the following cluster's field to the runnable (controller)
            1. `config` (`inject.ConfigInto(c.config, i)`)
            1. `client` (`inject.ClientInto(c.client, i)`)
            1. `apiReader` (`inject.APIReaderInto(c.apiReader, i)`)
            1. `scheme` (`inject.SchemeInto(c.scheme, i)`)
            1. `cache` (`inject.CacheInto(c.cache, i)`)
            1. `mapper` (`inject.MapperInto(c.mapper, i)`)
        1. `cm.SetFields` is set to `controller.SetFields` via `InjectorInto`. (details: [inject](../inject/)) <- `controller.SetFields` will be used for source, event handler and predicates in [Watch](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.13.0/pkg/internal/controller/controller.go#L129-L140).
        1. `StopChannelInto` and `Logger`.
    1. `cm.runnables.Add(r)`
        ```go
        type runnables struct {
            Webhooks       *runnableGroup
            Caches         *runnableGroup
            LeaderElection *runnableGroup
            Others         *runnableGroup
        }
        ```
        Add `r` based on the type.
        ```go
        func (r *runnables) Add(fn Runnable) error {
            switch runnable := fn.(type) {
            case hasCache:
                return r.Caches.Add(fn, func(ctx context.Context) bool {
                    return runnable.GetCache().WaitForCacheSync(ctx)
                })
            case *webhook.Server:
                return r.Webhooks.Add(fn, nil)
            case LeaderElectionRunnable:
                if !runnable.NeedLeaderElection() {
                    return r.Others.Add(fn, nil)
                }
                return r.LeaderElection.Add(fn, nil)
            default:
                return r.LeaderElection.Add(fn, nil)
            }
        }
        ```

## `Manager.GetClient()` and `GetScheme()`

1. The client, scheme and more are initialized and stored in the [cluster](../cluster) when a Manager is created.
1. The client, scheme and more are directly got from `cm.cluster.GetXXX()`
1. The client got by `GetClient()` is passed to `Reconciler` so you can manipulate objects in the Reconcile function.


## Example

1. Initialize with `NewManager`.

    ```go
    mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{})
    ```

    You can configure Options based on your requirements.
    example:

    ```go
    {
        Scheme:                 scheme,
        MetricsBindAddress:     metricsAddr,
        Port:                   9443,
        HealthProbeBindAddress: probeAddr,
        LeaderElection:         enableLeaderElection,
        LeaderElectionID:       "63ffe61d.example.com",
    }
    ```

1. Define a simple Reconciler

    ```go
	podReconciler := reconcile.Func(func(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
		fmt.Printf("podReconciler is called for %v\n", req)
		return reconcile.Result{}, nil
	})
    ```

    For more details about Reconciler, you can check [reconciler](../reconciler).

1. Set up Controller with `NewControllerManagedBy`

    ```go
    ctrl.NewControllerManagedBy(mgr). // returns controller Builder
        For(&corev1.Pod{}). // defines the type of Object being reconciled
        Complete(podReconciler) // Complete builds the Application controller, and return error
    ```

    1. `For`: define which resource to monitor.
    1. `Complete`: pass the reconciler to complete the controller.
    1. Internally, `NewControllerManagedBy` returns controller builder.
    1. Controller builder calls two functions in `Complete(reconcile.Reconciler)`
        1. [doController](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.13.0/pkg/builder/controller.go#L279): Set controller to the builder
            ```go
            blder.ctrl, err = newController(controllerName, blder.mgr, ctrlOptions)
            ```
        1. [doWatch](): call `blder.ctrl.Watch(src, hdler, allPredicates...)` for `For`, `Owns`, and `Watches`.

## Run

1. Run (initialize a Manager with podReconciler & deploymentReconciler)

    ```
    go run main.go
    2022-09-06T06:27:08.255+0900    INFO    controller-runtime.metrics      Metrics server is starting to listen    {"addr": ":8080"}
    2022-09-06T06:27:08.255+0900    INFO    Starting server {"path": "/metrics", "kind": "metrics", "addr": "[::]:8080"}
    2022-09-06T06:27:08.255+0900    INFO    Starting EventSource    {"controller": "pod", "controllerGroup": "", "controllerKind": "Pod", "source": "kind source: *v1.Pod"}
    2022-09-06T06:27:08.255+0900    INFO    Starting Controller     {"controller": "pod", "controllerGroup": "", "controllerKind": "Pod"}
    2022-09-06T06:27:08.255+0900    INFO    manager-examples        RunnableFunc is called
    2022-09-06T06:27:08.255+0900    INFO    Starting EventSource    {"controller": "deployment", "controllerGroup": "apps", "controllerKind": "Deployment", "source": "kind source: *v1.Deployment"}
    2022-09-06T06:27:08.255+0900    INFO    Starting Controller     {"controller": "deployment", "controllerGroup": "apps", "controllerKind": "Deployment"}
    2022-09-06T06:27:08.356+0900    INFO    Starting workers        {"controller": "pod", "controllerGroup": "", "controllerKind": "Pod", "worker count": 1}
    2022-09-06T06:27:08.357+0900    INFO    Starting workers        {"controller": "deployment", "controllerGroup": "apps", "controllerKind": "Deployment", "worker count": 1}
    2022-09-06T06:27:08.357+0900    INFO    manager-examples        podReconciler is called {"req": "kube-system/coredns-6d4b75cb6d-jtg59"}
    2022-09-06T06:27:08.357+0900    INFO    manager-examples        podReconciler is called {"req": "local-path-storage/local-path-provisioner-9cd9bd544-g89rs"}
    2022-09-06T06:27:08.357+0900    INFO    manager-examples        podReconciler is called {"req": "kube-system/kube-scheduler-kind-control-plane"}
    2022-09-06T06:27:08.357+0900    INFO    manager-examples        podReconciler is called {"req": "kube-system/kube-controller-manager-kind-control-plane"}
    2022-09-06T06:27:08.357+0900    INFO    manager-examples        podReconciler is called {"req": "kube-system/kube-proxy-7jsn6"}
    2022-09-06T06:27:08.357+0900    INFO    manager-examples        podReconciler is called {"req": "kube-system/coredns-6d4b75cb6d-k68r5"}
    2022-09-06T06:27:08.357+0900    INFO    manager-examples        podReconciler is called {"req": "kube-system/etcd-kind-control-plane"}
    2022-09-06T06:27:08.357+0900    INFO    manager-examples        podReconciler is called {"req": "kube-system/kube-apiserver-kind-control-plane"}
    2022-09-06T06:27:08.357+0900    INFO    manager-examples        podReconciler is called {"req": "kube-system/kindnet-6dj6q"}
    2022-09-06T06:27:08.358+0900    INFO    manager-examples        deploymentReconciler is called  {"req": "kube-system/coredns"}
    2022-09-06T06:27:08.358+0900    INFO    manager-examples        deploymentReconciler is called  {"req": "local-path-storage/local-path-provisioner"}
    ```

    The reconcile functions are called when cache is synced.

1. Create a Pod
    ```
    kubectl run nginx --image=nginx
    ```

    You'll see the following logs:

    ```
    2022-09-06T07:16:26.400+0900    INFO    manager-examples        podReconciler is called {"req": "default/nginx"}
    2022-09-06T07:16:26.519+0900    INFO    manager-examples        podReconciler is called {"req": "default/nginx"}
    2022-09-06T07:16:26.660+0900    INFO    manager-examples        podReconciler is called {"req": "default/nginx"}
    2022-09-06T07:16:32.547+0900    INFO    manager-examples        podReconciler is called {"req": "default/nginx"}
    ```
1. Delete the Pod
    ```
    kubectl delete pod nginx
    ```

    You'll see the logs again.
1. Create a Deployment
    ```
    kubectl create deploy nginx --image=nginx
    ```

    ```
    2022-09-06T07:17:04.963+0900    INFO    manager-examples        deploymentReconciler is called  {"req": "default/nginx"}
    2022-09-06T07:17:05.281+0900    INFO    manager-examples        deploymentReconciler is called  {"req": "default/nginx"}
    2022-09-06T07:17:05.320+0900    INFO    manager-examples        podReconciler is called {"req": "default/nginx-8f458dc5b-lnkqz"}
    2022-09-06T07:17:05.341+0900    INFO    manager-examples        podReconciler is called {"req": "default/nginx-8f458dc5b-lnkqz"}
    2022-09-06T07:17:05.342+0900    INFO    manager-examples        deploymentReconciler is called  {"req": "default/nginx"}
    2022-09-06T07:17:05.432+0900    INFO    manager-examples        podReconciler is called {"req": "default/nginx-8f458dc5b-lnkqz"}
    2022-09-06T07:17:05.461+0900    INFO    manager-examples        deploymentReconciler is called  {"req": "default/nginx"}
    2022-09-06T07:17:08.630+0900    INFO    manager-examples        podReconciler is called {"req": "default/nginx-8f458dc5b-lnkqz"}
    2022-09-06T07:17:08.674+0900    INFO    manager-examples        deploymentReconciler is called  {"req": "default/nginx"}
    ```

1. Delete the Deployment
    ```
    kubectl delete deploy nginx
    ```

    You'll see the logs again.
