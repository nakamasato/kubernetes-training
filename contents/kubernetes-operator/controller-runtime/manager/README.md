# [Manager](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/manager)

## [controllerManager](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.12.3/pkg/manager/internal.go)

## Interface

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

## [controllerManager type](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.12.3/pkg/manager/internal.go#L66-L173)

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

[Runnable](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.12.3/pkg/manager/manager.go#L293-L298) interface:

```go
type Runnable interface {
	Start(context.Context) error
}
```

[runnables](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.12.3/pkg/manager/runnable_group.go#L37-L45)

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

## How Manager is used

1. Initialize a [controllerManager](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.12.3/pkg/manager/internal.go#L66) with NewManager
1. Bind a controller using `NewControllerManagedBy` with controller builder.
1. Internally, calls the functions:
    1. [bldr.doController](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.12.3/pkg/builder/controller.go#L191) to register the controler to the buidler
        1. Create a new controller and add it by `Manager.Add(Runnable)`
    1. [bldr.doWatch](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.12.3/pkg/builder/controller.go#L196) to start watching the target resources configured by `For`, `Owns`, and `Watches`.
        1. The actual implementation of `Watch` function is in the controller. You can also check [controller](../controller)
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
        1. [cluster.SetFields](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.12.3/pkg/cluster/internal.go#L67) set dependencies on the object that implements the [inject](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.3/pkg/runtime/inject) interface. Specifically set the following cluster's field to the runnable (controller)
            1. `config` (`inject.ConfigInto(c.config, i)`)
            1. `client` (`inject.ClientInto(c.client, i)`)
            1. `apiReader` (`inject.APIReaderInto(c.apiReader, i)`)
            1. `scheme` (`inject.SchemeInto(c.scheme, i)`)
            1. `cache` (`inject.CacheInto(c.cache, i)`)
            1. `mapper` (`inject.MapperInto(c.mapper, i)`)
        1. `cm.SetFields` is set to `controller.SetFields` via `InjectorInto`. (details: [inject](../inject/)) <- `controller.SetFields` will be used for source, event handler and predicates in [Watch](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.12.3/pkg/internal/controller/controller.go#L129-L140).
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

1. `controllerManager`'s `Start()` calls `runnables.xxx.Start()` to start all runnables.
    ```go
	if err := cm.runnables.Webhooks.Start(cm.internalCtx); err != nil {
		if !errors.Is(err, wait.ErrWaitTimeout) {
			return err
		}
	}

	// Start and wait for caches.
	if err := cm.runnables.Caches.Start(cm.internalCtx); err != nil {
		if !errors.Is(err, wait.ErrWaitTimeout) {
			return err
		}
	}

	// Start the non-leaderelection Runnables after the cache has synced.
	if err := cm.runnables.Others.Start(cm.internalCtx); err != nil {
		if !errors.Is(err, wait.ErrWaitTimeout) {
			return err
		}
	}
    ```
    1. Controller will be in `runnables.Others` and you can check the actual `Start` logic in [controller](../controller).

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
        1. [doController](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.12.3/pkg/builder/controller.go#L279): Set controller to the builder
            ```go
            blder.ctrl, err = newController(controllerName, blder.mgr, ctrlOptions)
            ```
        1. [doWatch](): call `blder.ctrl.Watch(src, hdler, allPredicates...)` for `For`, `Owns`, and `Watches`.

## Run

1. Run (initialize a Manager with podReconciler & deploymentReconciler)

    ```
    go run main.go
    deploymentReconciler is called for kube-system/coredns
    podReconciler is called for kube-system/storage-provisioner
    podReconciler is called for kube-system/vpnkit-controller
    podReconciler is called for kube-system/coredns-6d4b75cb6d-l82x2
    podReconciler is called for kube-system/etcd-docker-desktop
    podReconciler is called for kube-system/kube-controller-manager-docker-desktop
    podReconciler is called for kube-system/kube-scheduler-docker-desktop
    podReconciler is called for kube-system/coredns-6d4b75cb6d-t8tp4
    podReconciler is called for kube-system/kube-apiserver-docker-desktop
    podReconciler is called for kube-system/kube-proxy-q4rp5
    ```

    The reconcile functions are called when cache is synced.

1. Create a Pod
    ```
    kubectl run nginx --image=nginx
    ```

    You'll see the following logs:

    ```
    podReconciler is called for default/nginx
    podReconciler is called for default/nginx
    podReconciler is called for default/nginx
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
    deploymentReconciler is called for default/nginx
    podReconciler is called for default/nginx-8f458dc5b-f86nt
    deploymentReconciler is called for default/nginx
    podReconciler is called for default/nginx-8f458dc5b-f86nt
    deploymentReconciler is called for default/nginx
    podReconciler is called for default/nginx-8f458dc5b-f86nt
    deploymentReconciler is called for default/nginx
    podReconciler is called for default/nginx-8f458dc5b-f86nt
    deploymentReconciler is called for default/nginx
    ```

1. Delete the Deployment
    ```
    kubectl delete deploy nginx
    ```

    You'll see the logs again.
