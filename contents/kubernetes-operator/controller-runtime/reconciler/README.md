# [Reconciler](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/reconcile)

Controller logic is implemented in terms of Reconcilers ([pkg/reconcile](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/reconcile)). A Reconciler implements a function which takes a reconcile Request containing the name and namespace of the object to reconcile, reconciles the object, and returns a Response or an error indicating whether to requeue for a second round of processing.


## Interface

```go
type Reconciler interface {
	// Reconcile performs a full reconciliation for the object referred to by the Request.
	// The Controller will requeue the Request to be processed again if an error is non-nil or
	// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
	Reconcile(context.Context, Request) (Result, error)
}
```

## Implement

You can use either implementation of the `Reconciler` interface:
1. a reconciler struct with `Reconcile` function.
1. a `reconcile.Func`, which implements Reconciler interface:
	```go
	type Func func(context.Context, Request) (Result, error)
	```

([Controller](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.12.3/pkg/internal/controller/controller.go#L42) also implements Reconciler interface. The reconciler passed to `builder` is used inside the controller's `Reconcile` function.)
## How reconciler is used
Reconciler is passed to Controller [builder](../builder) when initializing controller (you can also check it in [Manager](../manager/)):

```go
ctrl.NewControllerManagedBy(mgr). // returns controller Builder
    For(&corev1.Pod{}). // defines the type of Object being reconciled
    Complete(podReconciler) // Complete builds the Application controller, and return error
```
