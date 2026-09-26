# Dependency wiring

The old `inject` package and `Manager.SetFields` flow are not part of the pinned controller-runtime API. Dependencies are now constructed once by the Manager's Cluster and passed explicitly to the components that need them. This page describes the current wiring; the related source, controller, and manager diagrams are updated to match it.

## Shared dependencies

`manager.New` constructs a Cluster. The Cluster creates the shared Cache, Client, Scheme, RESTMapper, and uncached APIReader. Retrieve the dependency from Manager and pass it to a constructor or struct field:

| Dependency | Current source | Typical recipient |
| --- | --- | --- |
| Cache | `mgr.GetCache()` | `source.Kind(cache, object, handler)` |
| Client | `mgr.GetClient()` | Reconciler's `Client` field |
| APIReader | `mgr.GetAPIReader()` | Code that needs an uncached read |
| Config | `mgr.GetConfig()` | A client or Kubernetes API constructor |
| Scheme | `mgr.GetScheme()` | A decoder or client options |
| RESTMapper | `mgr.GetRESTMapper()` | Owner handlers and object mapping |
| Logger | `ctrl.LoggerFrom(ctx)` or an explicit logger | Reconciler or runnable |

There is no automatic recursive field injection. Constructors and struct literals make dependencies visible and keep ownership clear.

## Source and controller construction

The Builder creates sources and registers them with a Controller. Cluster-backed sources receive a Cache when constructed. Custom handlers and predicates receive any needed dependencies directly:

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

The Builder and Controller manage source startup and synchronization. `Manager.Start(ctx)` supplies the lifecycle context; components should honor its cancellation instead of receiving an injected stop channel. See [Cluster](../cluster/), [Builder](../builder/), [Source](../source/), and [Manager](../manager/) for the construction and startup paths.
