# builder


These implementation notes and diagrams target controller-runtime v0.25.1, pinned in the repository's go.mod. Use the signatures and call paths below for this version.

## Overview

![](overview.drawio.svg)

The main role of Builder is:

1. Create a Controller from the given Reconciler
1. Configure target resources for the controller
1. Register the controller to the Manager

About how the registered controllers are triggered, you can study in [Manager](../manager/). The controller is registered in the Manager's LeaderElection group by default. Controllers opting out of leader election are in Others. Both are started by `Manager.Start()` after cache readiness.

![](../manager/diagram.drawio.svg)


## Types

### [Builder](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.25.1/pkg/builder/controller.go)

```go
type Builder = TypedBuilder[reconcile.Request]

type TypedBuilder[request comparable] struct {
    forInput         ForInput
    ownsInput        []OwnsInput
    rawSources       []source.TypedSource[request]
    watchesInput     []WatchesInput[request]
    mgr              manager.Manager
    globalPredicates []predicate.Predicate
    ctrl             controller.TypedController[request]
    ctrlOptions      controller.TypedOptions[request]
    name             string
    newController    func(name string, mgr manager.Manager, options controller.TypedOptions[request]) (controller.TypedController[request], error)
}
```

## `ControllerManagedBy`: Initialize Builder with a Manager

Initialize a Builder with the specified manager.

```go
func ControllerManagedBy(m manager.Manager) *Builder {
    return TypedControllerManagedBy[reconcile.Request](m)
}
```

## [For](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.25.1/pkg/builder/controller.go), [Owns](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.25.1/pkg/builder/controller.go), and [Watches](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.25.1/pkg/builder/controller.go): Define what object to watch

![](for-owns-watches.drawio.svg)

1. `For(object client.Object, opts ...ForOption) *Builder`: only one resource can be configured. Same as
    ```go
    Watches(&corev1.Pod{}, &handler.EnqueueRequestForObject{})
    ```
3. `Owns(object client.Object, opts ...OwnsOption) *Builder`: Owns defines types of Objects being *generated* by the ControllerManagedBy, and configures the ControllerManagedBy to respond to create / delete / update events by **reconciling the owner object**. Same as the following code:
    ```go
    Watches(object, handler.EnqueueRequestForOwner(mgr.GetScheme(), mgr.GetRESTMapper(), ownerType, handler.OnlyControllerOwner()))
    ```
    [EnqueueRequestForOwner](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.25.1/pkg/handler/enqueue_owner.go): Extract owner object from ownerReferences and enqueue it to the queue.
3. `Watches(object client.Object, eventhandler handler.TypedEventHandler[client.Object, request], opts ...WatchesOption) *TypedBuilder[request]`: Watches exposes the lower-level ControllerManagedBy Watches functions through the builder. Consider using Owns or For instead of Watches directly.


Example:

```go
err = builder.
  ControllerManagedBy(mgr).  // Create the ControllerManagedBy
  For(&alpha1v1.Foo{}). // Foo is the Application API
  Owns(&corev1.Pod{}).       // Foo owns Pods created by it
  Complete(&FooReconciler{})
```

![](for-owns-example.drawio.svg)

## `Complete`: Receive reconciler and build a controller

Receive `reconcile.Reconciler` and call `Build`:

```go
func (blder *TypedBuilder[request]) Complete(r reconcile.TypedReconciler[request]) error {
    _, err := blder.Build(r)
    return err
}
```

## `Build`: Create a controller and return the controller.

```go
func (blder *TypedBuilder[request]) Build(r reconcile.TypedReconciler[request]) (controller.TypedController[request], error) {
    if r == nil {
        return nil, fmt.Errorf("must provide a non-nil Reconciler")
    }
    if blder.mgr == nil {
        return nil, fmt.Errorf("must provide a non-nil Manager")
    }
    if blder.forInput.err != nil {
        return nil, blder.forInput.err
    }

    if err := blder.doController(r); err != nil {
        return nil, err
    }

    if err := blder.doWatch(); err != nil {
        return nil, err
    }

    return blder.ctrl, nil
}
```

1. [bldr.doController](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.25.1/pkg/builder/controller.go) to construct the controller and register it with the Manager
    1. Create a new controller.
        ```go
        blder.ctrl, err = controller.NewTyped(controllerName, blder.mgr, ctrlOptions)
        ```
    1. the controller is added to the appropriate Manager runnable group by `Manager.Add(Runnable)` in `newController`. ([controller](../controller/README.md#how-controller-is-used))
1. [bldr.doWatch](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.25.1/pkg/builder/controller.go) to register watches for the target resources configured by `For`, `Owns`, and `Watches`.
    1. The actual implementation of `Watch` function is in the [controller](../controller)

## Convert `client.Object` to `Source`

1. [Controller.Watch](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.25.1/pkg/controller/controller.go) needs `Source` as the first argument.
    ```go
    Watch(src source.TypedSource[request]) error
    ```
1. `client.Object` is set in `ForInput`, `OwnsInput`, and `WatchesInput` for `For`, `Owns`, and `Watches` respectively.
1. Before calling `Controller.Watch`, the `client.Object` needs to be [project](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.25.1/pkg/builder/controller.go)ed into [Source](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.25.1/pkg/source/source.go) based on `objectProjection`:
    1. `projectAsNormal`: Use the object as it is. (**In most cases**)
    1. `projectAsMetadata`: Extract only metadata.

    ```go
    typeForSrc, err := blder.project(blder.forInput.object, blder.forInput.objectProjection)
    if err != nil {
        return err
    }
    hdlr := &handler.EnqueueRequestForObject{}
    src := source.Kind(blder.mgr.GetCache(), typeForSrc, hdlr, allPredicates...)
    if err := blder.ctrl.Watch(src); err != nil {
        return err
    }
    ```

    `source.Kind` constructs an internal Kind, which implements the [Source](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.25.1/pkg/source/source.go) interface.

    ```go
    type Kind[object client.Object, request comparable] struct {
        Type object

        Cache cache.Cache

        Handler handler.TypedEventHandler[object, request]

        Predicates []predicate.TypedPredicate[object]

        startedErr  chan error
        startCancel func()
    }
    ```

`For` uses the object's own key; `Owns` resolves ownerReferences and enqueues the owner's key; `Watches` uses your supplied mapping handler. Merely matching labels does not establish ownership. Predicates run before the handler, so filtering a Delete or an initial Add can suppress a needed reconciliation.

`WatchesRawSource` accepts an already configured source, such as a Channel, and does not automatically apply the builder's global event filter. Metadata projection watches PartialObjectMetadata: use that same representation for reads if you intend to avoid starting an additional typed informer.

`Complete` returns only the construction error; `Build` also returns the Controller. Both must be checked. The example's Foo type and FooReconciler stand for your own API and implementation.
