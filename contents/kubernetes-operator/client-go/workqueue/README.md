# Workqueue

A workqueue holds keys for a controller to process. It coalesces duplicate additions. If a key is added again while being processed, it can be processed again after `Done`.

Use the typed API:

```go
q := workqueue.NewTypedRateLimitingQueue(
    workqueue.DefaultTypedControllerRateLimiter[types.NamespacedName](),
)
defer q.ShutDown()
q.Add(types.NamespacedName{Namespace: "default", Name: "example"})
```

A worker can process one item as follows. Here, `sync` reads and reconciles the current state for the key:

```go
key, shutdown := q.Get()
if shutdown {
    return false
}
defer q.Done(key)
if err := sync(ctx, key); err != nil {
    q.AddRateLimited(key)
} else {
    q.Forget(key)
}
return true
```

- `Done` marks processing of an acquired item as complete, on both success and failure.
- `Forget` resets rate-limiter state, such as failure counts. It does not replace `Done`.
- `AddAfter` schedules an item after a duration; `AddRateLimited` uses the rate limiter to schedule it.
- `Get` blocks on an empty queue. Context cancellation alone does not unblock it; call `ShutDown` when stopping.

The [Source example](../../controller-runtime/source) uses a typed queue and calls `ShutDown` when its context is canceled. A controller-runtime Controller normally manages queues and retries internally, so its Reconciler returns a Result and an error.

References: [Workqueue API](https://pkg.go.dev/k8s.io/client-go@v0.37.1/util/workqueue), [official example](https://github.com/kubernetes/client-go/blob/v0.37.1/examples/workqueue/main.go).
