# [handler](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.25.1/pkg/handler)

Package `handler` defines `EventHandlers` that enqueue `reconcile.Request`s in response to Create, Update, Deletion Events observed from Watching Kubernetes APIs. Construct a source with its handler and pass that source to `Controller.Watch` to generate and enqueue requests.

`handler.EventHandler` is the alias `TypedEventHandler[client.Object, reconcile.Request]`. Handler constructors and typed functions also support custom comparable request keys.

1. Unless you are implementing your own EventHandler, you can ignore the functions on the `EventHandler` interface.
1. Most users shouldn't need to implement their own EventHandler.

## [EventHandler interface](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.25.1/pkg/handler/eventhandler.go)

```go
type EventHandler = TypedEventHandler[client.Object, reconcile.Request]

type TypedEventHandler[object any, request comparable] interface {
    Create(context.Context, event.TypedCreateEvent[object], workqueue.TypedRateLimitingInterface[request])

    Update(context.Context, event.TypedUpdateEvent[object], workqueue.TypedRateLimitingInterface[request])

    Delete(context.Context, event.TypedDeleteEvent[object], workqueue.TypedRateLimitingInterface[request])

    Generic(context.Context, event.TypedGenericEvent[object], workqueue.TypedRateLimitingInterface[request])
}
```

## [EnqueueRequestForObject](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.25.1/pkg/handler/enqueue.go)

This is used by default in [builder.doWatch](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.25.1/pkg/builder/controller.go). If you create an operator with kubebuilder, you're using this eventhandler.
This function converts events received from the Source into `reconcile.Request`s object and enqueue them to the given queue.

1. `Create`, `Delete`, `Generic` map the event object to its namespace/name. Conceptually:
    ```go
    q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
        Name:      evt.Object.GetName(),
        Namespace: evt.Object.GetNamespace(),
    }})
    ```
1. `Update`: Enqueue ObjectNew (ObjectOld if ObjectNew doesn't exist)
    ```go
    q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
        Name:      evt.ObjectNew.GetName(),
        Namespace: evt.ObjectNew.GetNamespace(),
    }})
    ```

## How handlers are used

`source.Kind(cache, object, handler, predicates...)` attaches the handler to informer events. The internal event adapter checks predicates before invoking Create/Update/Delete/Generic with context and the controller's typed queue. Keep handlers fast: they schedule reconciliation, while API writes and business logic belong in Reconcile.

The actual Update implementation prefers ObjectNew and falls back to ObjectOld when it is absent:

```go
func (e *TypedEnqueueRequestForObject[T]) Update(ctx context.Context, evt event.TypedUpdateEvent[T], q workqueue.TypedRateLimitingInterface[reconcile.Request]) {
    switch {
    case !isNil(evt.ObjectNew):
        item := reconcile.Request{NamespacedName: types.NamespacedName{
            Name:      evt.ObjectNew.GetName(),
            Namespace: evt.ObjectNew.GetNamespace(),
        }}

        addToQueueUpdate(q, evt, item)
    case !isNil(evt.ObjectOld):
        item := reconcile.Request{NamespacedName: types.NamespacedName{
            Name:      evt.ObjectOld.GetName(),
            Namespace: evt.ObjectOld.GetNamespace(),
        }}

        addToQueueUpdate(q, evt, item)
    default:
        enqueueLog.Error(nil, "UpdateEvent received with no metadata", "event", evt)
    }
}
```

Create and Update use queue helpers that can lower the priority of initial-list or unchanged events on priority-aware queues. The mapping remains namespace/name. Multiple events can deduplicate into one request.

## Owners and custom mappings

- `EnqueueRequestForOwner(scheme, mapper, ownerType, handler.OnlyControllerOwner())` inspects ownerReferences and enqueues the owner's namespace/name. Builder.Owns configures this handler for the primary For type. Scheme resolves the owner GVK; RESTMapper determines whether it is namespaced.
- `EnqueueRequestsFromMapFunc` maps a watched object to one or more different objects. Use it for relationships that are not ownership, such as a ConfigMap referenced by several custom resources. On Update, mapping both old and new objects avoids losing requests when the relationship changes.
- `handler.TypedFuncs` lets a standalone consumer enqueue custom keys. The [Source sample](../source/) includes Kind, Namespace, and Name rather than allowing resources with the same name to collide.
