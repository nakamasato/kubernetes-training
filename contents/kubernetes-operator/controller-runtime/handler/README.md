# [handler](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/handler)

Package handler defines EventHandlers that enqueue reconcile.Requests in response to Create, Update, Deletion Events observed from Watching Kubernetes APIs. Users should provide a source.Source and handler.EventHandler to Controller.Watch in order to generate and enqueue reconcile.Request work items.qq

handler.EventHandler is an argument to Controller.Watch that enqueues reconcile.Requests in response to events.

1. *Unless you are implementing your own EventHandler, you can ignore the functions on the EventHandler interface.
1. Most users shouldn't need to implement their own EventHandler.*

## [EventHandler interface](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.13.0/pkg/handler/eventhandler.go)

```go
// * Use EnqueueRequestsFromMapFunc to transform an event for an object to a reconcile of an object
// of a different type - do this for events for types the Controller may be interested in, but doesn't create.
// (e.g. If Foo responds to cluster size events, map Node events to Foo objects.)
//
// Unless you are implementing your own EventHandler, you can ignore the functions on the EventHandler interface.
// Most users shouldn't need to implement their own EventHandler.
type EventHandler interface {
	// Create is called in response to an create event - e.g. Pod Creation.
	Create(event.CreateEvent, workqueue.RateLimitingInterface)

	// Update is called in response to an update event -  e.g. Pod Updated.
	Update(event.UpdateEvent, workqueue.RateLimitingInterface)

	// Delete is called in response to a delete event - e.g. Pod Deleted.
	Delete(event.DeleteEvent, workqueue.RateLimitingInterface)

	// Generic is called in response to an event of an unknown type or a synthetic event triggered as a cron or
	// external trigger request - e.g. reconcile Autoscaling, or a Webhook.
	Generic(event.GenericEvent, workqueue.RateLimitingInterface)
}
```

## [EnqueueRequestForObject](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.13.0/pkg/handler/enqueue.go#L36)

This is used by default in [builder.doWatch](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.13.0/pkg/builder/controller.go#L227). If you create an operator with kubebuilder, you're using this eventhandler.

1. `Create`, `Delete`, `Generic`:
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
