# inject (historical APIs and current replacements)

The interfaces and code below describe controller-runtime v0.13.0. They are retained to explain the original diagrams and the old dependency propagation. They are not available in the pinned v0.25.1; use explicit constructors and fields as shown in the migration section.

## Inject interface table

|interface name|required func|func to inject|Implemented by|
|---|---|---|---|
|Cache|InjectCache|CacheInto|Kind (source)|
|APIReader|InjectAPIReader|APIReaderInto||
|Config|InjectConfig|ConfigInto||
|Client|InjectClient|ClientInto||
|Scheme|InjectScheme|SchemeInto|DeferredFileLoader, Webhook|
|Stoppable|InjectStopChannel|StopChannelInto|Channel (source)|
|Mapper|InjectMapper|MapperInto|EnqueueRequestForOwner|
|Injector|InjectFunc|InjectorInto|Webhook, enqueueRequestsFromMapFunc, Controller, and, or (predicate), multiMutating, multiValidating, etc.|
|Logger|InjectLogger|LoggerInto|Webhook|

`Injector` interface

```go
// Func injects dependencies into i.
type Func func(i interface{}) error

// Injector is used by the ControllerManager to inject Func into Controllers.
type Injector interface {
    InjectFunc(f Func) error
}

// InjectorInto will set f and return the result on i if it implements Injector.  Returns
// false if i does not implement Injector.
func InjectorInto(f Func, i interface{}) (bool, error) {
    if ii, ok := i.(Injector); ok {
        return true, ii.InjectFunc(f)
    }
    return false, nil
}
```

`SetFields` set dependencies to the object that implments inject interface.

## Usage in v0.13.0

Controller implement Injector with [InjectFunc](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.13.0/pkg/internal/controller/controller.go#L352)

```go
// InjectFunc implement SetFields.Injector.
func (c *Controller) InjectFunc(f inject.Func) error {
    c.SetFields = f
    return nil
}
```

With this implementation, any function can be injected to the `controller.SetFields` with `InjectorInto(func, controller)`.

This function is used in the [manager.SetFields](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.13.0/pkg/manager/internal.go#L187-L211)

```go
if _, err := inject.InjectorInto(cm.SetFields, i); err != nil {
```
This means set `clusterManager`'s `SetFields` function to `i`, specifically, `controller`. By the controller's `InjectFunc` implementation, **controller has exactly the same `SetFields` function as `clusterManager`**

## Current dependency wiring (v0.25.1)

The former flow was `Manager.Add → Manager.SetFields → Cluster.SetFields`, with Injector propagating SetFields into the Controller so Watch could populate sources, handlers, and predicates. Today Manager.Add registers lifecycle work; it does not inject those fields.

| Former interface | Current replacement |
| --- | --- |
| Cache / InjectCache | Pass `mgr.GetCache()` to `source.Kind` or `source.TypedKind` |
| Client / InjectClient | Initialize the reconciler's Client with `mgr.GetClient()` |
| APIReader / InjectAPIReader | Pass `mgr.GetAPIReader()` where uncached reads are required |
| Config / InjectConfig | Pass `mgr.GetConfig()` to constructors |
| Scheme / InjectScheme | Pass `mgr.GetScheme()` to constructors and decoders |
| Mapper / InjectMapper | Pass `mgr.GetRESTMapper()` to owner handlers |
| Stoppable / InjectStopChannel | Honor the context passed to Start |
| Logger / InjectLogger | Pass a logger explicitly or get it from context |
| Injector / InjectFunc | Constructor or struct-field wiring; no automatic recursive injection |

```go
r := &ReplicaSetReconciler{Client: mgr.GetClient()}
err := ctrl.NewControllerManagedBy(mgr).
    For(&appsv1.ReplicaSet{}).
    Owns(&corev1.Pod{}).
    Complete(r)
if err != nil {
    return err
}
```

See [Cluster](../cluster/#setfields-migration-explicit-dependency-wiring), [Builder](../builder/), and [Source](../source/) for the current construction paths. Custom event handlers and predicates must also receive their dependencies explicitly.
