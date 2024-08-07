# Source

Component structure:
![](diagram.drawio.svg)

Dataflow:
![](dataflow.drawio.svg)

## [Source](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.17.0/pkg/source/source.go#L47-L51) interface

```go
type Source interface {
	// Start is internal and should be called only by the Controller to register an EventHandler with the Informer
	// to enqueue reconcile.Requests.
	Start(context.Context, handler.EventHandler, workqueue.RateLimitingInterface, ...predicate.Predicate) error
}

// SyncingSource is a source that needs syncing prior to being usable. The controller
// will call its WaitForSync prior to starting workers.
type SyncingSource interface {
	Source
	WaitForSync(ctx context.Context) error
}
```

## Implementations

1. Already removed in [Refactor source/handler/predicate packages to remove dep injection](https://github.com/kubernetes-sigs/controller-runtime/pull/2120) (from [v0.15.0](https://github.com/kubernetes-sigs/controller-runtime/releases/tag/v0.15.0))
    ~~`Kind` has `InjectCache` while `kindWithCache` doesn't.~~
1. [Kind](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.17.0/pkg/internal/source/kind.go#L20-L31) Kind is used to provide a source of **events originating inside the cluster** from Watches (e.g. Pod Create).
    ```go
    type Kind struct {
        Type client.Object
        cache cache.Cache
        started     chan error
        startCancel func()
    }
    ```
    1. Provide cache explicitly with the initialization method [`func Kind(cache cache.Cache, object client.Object) SyncingSource`](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.17.0/pkg/source/source.go#L61).
        ```go
        source.Kind(mgr.GetCache(), &corev1.Pod{})
        ```
1. [Channel](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.17.0/pkg/source/source.go#L70-L86): Channel is used to provide a source of **events originating outside the cluster** (e.g. GitHub Webhook callback).  **Channel requires the user to wire the external source** (eh.g. http handler) to write GenericEvents to the underlying channel.
1. [Informer](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.17.0/pkg/source/source.go#L185-L188): Informer is used to provide a source of **events originating inside the cluster** from Watches (e.g. Pod Create).
    ```go
    type Informer struct {
        // Informer is the controller-runtime Informer
        Informer cache.Informer
    }
    ```

What's the difference between `Informer` and `Kind`?

1. `Kind` gets informer from `cache.Cache`.
1. `Informer` needs to be initialized with `cache.Informer` directly.
1. `Kind` calls `WaitForCacheSync` after adding eventhandler by `AddEventHandler` while `Informer` doesn't.


## How `Source` is used

1. `Source` is initialized in `builder.doWatch` for each of `For`, `Owns`, and `Watches`:
    1. [For](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.17.0/pkg/builder/controller.go#L271-L275):
        ```go
        // Reconcile type
        typeForSrc, err := blder.project(blder.forInput.object, blder.forInput.objectProjection)
        if err != nil {
            return err
        }
        src := &source.Kind{Type: typeForSrc}
        ```
    1. [Owns](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.17.0/pkg/builder/controller.go#L289-L293):
        ```go
        typeForSrc, err := blder.project(own.object, own.objectProjection)
		if err != nil {
			return err
		}
		src := &source.Kind{Type: typeForSrc}
        ```
    1. [Watches](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.17.0/pkg/builder/controller.go#L316-L322):
        ```go
        // If the source of this watch is of type *source.Kind, project it.
		if srckind, ok := w.src.(*source.Kind); ok {
			typeForSrc, err := blder.project(srckind.Type, w.objectProjection)
			if err != nil {
				return err
			}
			srckind.Type = typeForSrc
		}
        ```
1. The initialized source is passed to `controller.Watch` in [builder.doWatch](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.17.0/pkg/builder/controller.go#L279) if the controller is initialized by [builder](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.17.0/pkg/builder/controller.go#L56)

    ```go
    if err := blder.ctrl.Watch(w.src, w.eventhandler, allPredicates...); err != nil {
        return err
    }
    ```
1. In [controller.Watch](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.17.0/pkg/internal/controller/controller.go#L136)
    1. `Cache` is injected from controller.
        ```go
        // Inject Cache into arguments
        if err := c.SetFields(src); err != nil {
            return err
        }
        ```
    1. `source.Start` is called with `EventHandler` and `Queue`
        ```go
        return src.Start(c.ctx, evthdler, c.Queue, prct...)
        ```
1. [Source.Start](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.17.0/pkg/internal/source/kind.go#L35)
    1. Get `informer` from the injected `cache`.
        ```go
        i, lastErr = ks.cache.GetInformer(ctx, ks.Type)
        ```
    1. Add the event handler with `AddEventHandler`
        ```go
        i.AddEventHandler(internal.EventHandler{Queue: queue, EventHandler: handler, Predicates: prct})
        ```
1. informer is started by `manager.Start()`.
    ```go
    manager.runnables.Cache.Start()
    ```
    1. cache is initialized in [cluster](../cluster/README.md#set-fields) when a [Manager](../manager/README.md#1-initialize-a-controllermanagerhttpsgithubcomkubernetes-sigscontroller-runtimeblobv0123pkgmanagerinternalgol66-with-newmanager) is created.

## Example Usage: Debugging your controller

if you want to check events of specific resource you can set by the following.

1. Set up scheme if you want to monitor CRD (Optional)
    ```go
    import (
        mysqlv1alpha1 "github.com/nakamasato/mysql-operator/api/v1alpha1" // Target CRD
	    utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	    clientgoscheme "k8s.io/client-go/kubernetes/scheme"
    )

    func init() {
    	utilruntime.Must(mysqlv1alpha1.AddToScheme(scheme))
    	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
    }
    ```

1. `cache.Get()`: Internally create informer if not exists.

    ```go
    pod := &v1.Pod{}
    cache.Get(ctx, client.ObjectKeyFromObject(pod), pod)

    mysqluser := &mysqlv1alpha1.MySQLUser{}
    cache.Get(ctx, client.ObjectKeyFromObject(mysqluser), mysqluser)
    ```
1. Start the cache.

    ```go
	go func() {
		if err := cache.Start(ctx); err != nil { // func (m *InformersMap) Start(ctx context.Context) error {
			log.Error(err, "failed to start cache")
		}
	}()
    ```
1. Create `Kind` for the target resource.
    ```go
    kind := source.Kind(cache, mysqluser)
    ```
1. Prepare `workqueue` and `eventHandler`.
    ```go
	queue := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "test")
	eventHandler := handler.Funcs{
		CreateFunc: func(e event.CreateEvent, q workqueue.RateLimitingInterface) {
			log.Info("CreateFunc is called", "object", e.Object.GetName())
			// queue.Add(WorkQueueItem{Event: "Create", Name: e.Object.GetName()})
		},
		UpdateFunc: func(e event.UpdateEvent, q workqueue.RateLimitingInterface) {
			log.Info("UpdateFunc is called", "objectNew", e.ObjectNew.GetName(), "objectOld", e.ObjectOld.GetName())
			// queue.Add(WorkQueueItem{Event: "Update", Name: e.ObjectNew.GetName()})
		},
		DeleteFunc: func(e event.DeleteEvent, q workqueue.RateLimitingInterface) {
			log.Info("DeleteFunc is called", "object", e.Object.GetName())
			// queue.Add(WorkQueueItem{Event: "Delete", Name: e.Object.GetName()})
		},
	}
    ```
1. Start `kind` with the prepared `eventHandler` and `queue`.
    ```go
	if err := kind.Start(ctx, eventHandler, queue); err != nil { // Get informer and set eventHandler
		log.Error(err, "")
	}
    ```

You can run:

1. Install your CRD
    ```
    kubectl apply -f https://raw.githubusercontent.com/nakamasato/mysql-operator/main/config/crd/bases/mysql.nakamasato.com_mysqlusers.yaml
    ```
1. Run the `kind`
    ```
    go run main.go
    ```

    <details>

    ```
    2022-09-15T06:58:43.895+0900    INFO    source-examples source start
    2022-09-15T06:58:44.070+0900    INFO    source-examples cache is created
    2022-09-15T06:58:44.071+0900    INFO    source-examples cache is started
    2022-09-15T06:58:44.096+0900    INFO    source-examples CreateFunc is called    {"object": "kube-apiserver-kind-control-plane"}
    2022-09-15T06:58:44.097+0900    INFO    source-examples CreateFunc is called    {"object": "kube-controller-manager-kind-control-plane"}
    2022-09-15T06:58:44.097+0900    INFO    source-examples CreateFunc is called    {"object": "kube-scheduler-kind-control-plane"}
    2022-09-15T06:58:44.097+0900    INFO    source-examples CreateFunc is called    {"object": "kube-proxy-zpj2w"}
    2022-09-15T06:58:44.097+0900    INFO    source-examples CreateFunc is called    {"object": "coredns-6d4b75cb6d-s2dhg"}
    2022-09-15T06:58:44.097+0900    INFO    source-examples CreateFunc is called    {"object": "coredns-6d4b75cb6d-25dbf"}
    2022-09-15T06:58:44.097+0900    INFO    source-examples CreateFunc is called    {"object": "etcd-kind-control-plane"}
    2022-09-15T06:58:44.097+0900    INFO    source-examples CreateFunc is called    {"object": "kindnet-8fjbg"}
    2022-09-15T06:58:44.097+0900    INFO    source-examples CreateFunc is called    {"object": "local-path-provisioner-9cd9bd544-xl67h"}
    2022-09-15T06:58:44.172+0900    INFO    source-examples kindWithCache is ready
    2022-09-15T06:58:44.172+0900    INFO    source-examples got item        {"item": {"Event":"Create","Name":"kube-apiserver-kind-control-plane"}}
    2022-09-15T06:58:44.172+0900    INFO    source-examples got item        {"item": {"Event":"Create","Name":"kube-controller-manager-kind-control-plane"}}
    2022-09-15T06:58:44.172+0900    INFO    source-examples got item        {"item": {"Event":"Create","Name":"kube-scheduler-kind-control-plane"}}
    2022-09-15T06:58:44.172+0900    INFO    source-examples got item        {"item": {"Event":"Create","Name":"kube-proxy-zpj2w"}}
    2022-09-15T06:58:44.172+0900    INFO    source-examples got item        {"item": {"Event":"Create","Name":"coredns-6d4b75cb6d-s2dhg"}}
    2022-09-15T06:58:44.172+0900    INFO    source-examples got item        {"item": {"Event":"Create","Name":"coredns-6d4b75cb6d-25dbf"}}
    2022-09-15T06:58:44.172+0900    INFO    source-examples got item        {"item": {"Event":"Create","Name":"etcd-kind-control-plane"}}
    2022-09-15T06:58:44.172+0900    INFO    source-examples got item        {"item": {"Event":"Create","Name":"kindnet-8fjbg"}}
    2022-09-15T06:58:44.172+0900    INFO    source-examples got item        {"item": {"Event":"Create","Name":"local-path-provisioner-9cd9bd544-xl67h"}}
    ```

    </details>

1. Create custom resource manually. You'll see the events related to the CRD.
    ```
    kubectl apply -f https://raw.githubusercontent.com/nakamasato/mysql-operator/main/config/samples/mysql_v1alpha1_mysqluser.yaml
    ```

    ```
    2024-07-27T11:44:33.871+0900    INFO    source-examples CreateFunc is called    {"object": "sample-user"}
    2024-07-27T11:44:33.872+0900    INFO    source-examples got item        {"item": {"Event":"Create","Name":"sample-user"}}
    ```

1. Create nginx Pod
    ```
    kubectl run nginx --image=nginx
    ```

    ```
    2024-07-27T11:45:13.787+0900    INFO    source-examples CreateFunc is called    {"object": "nginx"}
    2024-07-27T11:45:13.787+0900    INFO    source-examples got item        {"item": {"Event":"Create","Name":"nginx"}}
    2024-07-27T11:45:13.805+0900    INFO    source-examples UpdateFunc is called    {"objectNew": "nginx", "objectOld": "nginx"}
    2024-07-27T11:45:13.805+0900    INFO    source-examples got item        {"item": {"Event":"Update","Name":"nginx"}}
    2024-07-27T11:45:13.819+0900    INFO    source-examples UpdateFunc is called    {"objectNew": "nginx", "objectOld": "nginx"}
    2024-07-27T11:45:16.417+0900    INFO    source-examples UpdateFunc is called    {"objectNew": "nginx", "objectOld": "nginx"}
    ```
