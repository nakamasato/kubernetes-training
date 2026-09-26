# [ListerWatcher](https://github.com/kubernetes/client-go/blob/v0.37.1/tools/cache/listwatch.go)

The examples target client-go v0.37.1. `ListWatch` supports both legacy and context-aware interfaces. Prefer `ListWithContext` / `WatchWithContext`; the methods without context are deprecated. The construction code below retains the legacy callbacks because the library still provides them for compatibility.

## type

### Interface

```go
// ListerWatcher is any object that knows how to perform an initial list and start a watch on a resource.
type ListerWatcherWithContext interface {
    ListerWithContext
    WatcherWithContext
}
```

```go
type ListerWithContext interface {
    ListWithContext(ctx context.Context, options metav1.ListOptions) (runtime.Object, error)
}

type WatcherWithContext interface {
    WatchWithContext(ctx context.Context, options metav1.ListOptions) (watch.Interface, error)
}
```

### ListWatch struct

```go
type ListWithContextFunc func(context.Context, metav1.ListOptions) (runtime.Object, error)
type WatchFuncWithContext func(context.Context, metav1.ListOptions) (watch.Interface, error)

type ListWatch struct {
    ListFunc ListFunc
    WatchFunc WatchFunc

    ListWithContextFunc  ListWithContextFunc
    WatchFuncWithContext WatchFuncWithContext

    DisableChunking bool
}
```

[watch.Interface](https://pkg.go.dev/k8s.io/apimachinery/pkg/watch#Interface):

```go
type Interface interface {
    // Stop stops watching. Will close the channel returned by ResultChan(). Releases
    // any resources used by the watch.
    Stop()

    // ResultChan returns a chan which will receive all the events. If an error occurs
    // or Stop() is called, the implementation will close this channel and
    // release any resources used by the watch.
    ResultChan() <-chan Event
}
```

## How ListWatch is used

1. Created with [NewFilteredListWatchFromClient](https://github.com/kubernetes/client-go/blob/v0.37.1/tools/cache/listwatch.go):

    ```go
    func NewFilteredListWatchFromClient(c Getter, resource string, namespace string, optionsModifier func(options *metav1.ListOptions)) *ListWatch {
        listFunc := func(options metav1.ListOptions) (runtime.Object, error) {
            optionsModifier(&options)
            return c.Get().
                Namespace(namespace).
                Resource(resource).
                VersionedParams(&options, metav1.ParameterCodec).
                Do(context.Background()).
                Get()
        }
        watchFunc := func(options metav1.ListOptions) (watch.Interface, error) {
            options.Watch = true
            optionsModifier(&options)
            return c.Get().
                Namespace(namespace).
                Resource(resource).
                VersionedParams(&options, metav1.ParameterCodec).
                Watch(context.Background())
        }
        listFuncWithContext := func(ctx context.Context, options metav1.ListOptions) (runtime.Object, error) {
            optionsModifier(&options)
            return c.Get().
                Namespace(namespace).
                Resource(resource).
                VersionedParams(&options, metav1.ParameterCodec).
                Do(ctx).
                Get()
        }
        watchFuncWithContext := func(ctx context.Context, options metav1.ListOptions) (watch.Interface, error) {
            options.Watch = true
            optionsModifier(&options)
            return c.Get().
                Namespace(namespace).
                Resource(resource).
                VersionedParams(&options, metav1.ParameterCodec).
                Watch(ctx)
        }
        return &ListWatch{
            ListFunc:             listFunc,
            WatchFunc:            watchFunc,
            ListWithContextFunc:  listFuncWithContext,
            WatchFuncWithContext: watchFuncWithContext,
        }
    }
    ```

1. Getter can be got from `clientset`: `clientset.CoreV1().RESTClient()`

    ```go
    podListWatcher := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "pods", v1.NamespaceDefault, fields.Everything())
    ```
1. Both legacy and context-aware callbacks are set in `NewFilteredListWatchFromClient`:
    1. `c.Get()` calls `NewRequest().Verb("GET")`
    1. [rest.NewRequest](https://github.com/kubernetes/client-go/blob/v0.37.1/rest/request.go) creates and returns a request.
    1. [Request.Watch](https://github.com/kubernetes/client-go/blob/v0.37.1/rest/request.go)
        1. get retry func
        1. in a for loop, run `retry.Before`, `client.Do(req)`, and `retry.After`

1. `ListWithContext` dispatches to `ListWithContextFunc` when available, otherwise to the legacy callback. `WatchWithContext` does the same for `WatchFuncWithContext`.
1. `Request.Watch(ctx)` creates the streaming HTTP request and returns a `watch.Interface`. Stop it with `Stop()` or by canceling its context. A closed result channel ends this example; a Reflector handles reconnects and expired resource versions.
1. Then, the ListWatch is passed to an informer.
    ```go
    indexer, informer := cache.NewIndexerInformer(
        podListWatcher, &v1.Pod{}, 0,
        cache.ResourceEventHandlerFuncs{AddFunc: handleAdd},
        cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
    )
    ```
    `ListerWatcher` is used to initialize the store (DeltaFIFO) with `List` and keep the store up-to-date with `Watch`. For more details, you can check [informer](../informer/README.md)


## Example

Read events through `w.ResultChan()`. First list the resource and use the collection's resourceVersion to start the watch without a gap. The context is canceled on Ctrl+C in [main.go](main.go).

```go
list, err := podListWatcher.ListWithContext(ctx, metav1.ListOptions{})
if err != nil {
    return err
}
listMeta, err := meta.ListAccessor(list)
if err != nil {
    return err
}
w, err := podListWatcher.WatchWithContext(ctx, metav1.ListOptions{
    ResourceVersion: listMeta.GetResourceVersion(),
})
if err != nil {
    return err
}
defer w.Stop()
for event := range w.ResultChan() {
    if event.Type == watch.Error {
        return apierrors.FromObject(event.Object)
    }
    objectMeta, err := meta.Accessor(event.Object)
    if err != nil {
        continue
    }
    klog.Infof("event: %s, resourceVersion: %s", event.Type, objectMeta.GetResourceVersion())
}
```

Run from the repository root with permission to list/watch Pods in the default namespace:

1. Start the ListerWatcher for Pods.
    ```
    go run ./contents/kubernetes-operator/client-go/listerwatcher
    I0913 07:54:29.394053   92277 main.go:57] resourceVersion: 2728
    ```
1. Create a Pod
    ```
    kubectl run nginx --image=nginx
    ```
    You'll see the event logs:
    ```
    I0913 07:54:29.394239   92277 main.go:64] items: 1
    I0913 07:54:29.401483   92277 main.go:84] event: ADDED, resourceVersion: 503
    I0913 07:55:20.475959   92277 main.go:84] event: MODIFIED, resourceVersion: 2789
    I0913 07:55:38.769688   92277 main.go:84] event: MODIFIED, resourceVersion: 2812
    ```
1. Patch the Pod
    ```
    kubectl patch pod nginx -p '{"metadata":{"annotations": {"key": "val"}}}' --type=merge
    ```
    You'll see `MODIFIED` in the logs.
1. Delete the Pod
    ```
    kubectl delete pod nginx
    ```

    ```
    I0913 08:02:23.486861   92277 main.go:84] event: MODIFIED, resourceVersion: 3294
    I0913 08:02:23.929399   92277 main.go:84] event: MODIFIED, resourceVersion: 3298
    I0913 08:02:24.239499   92277 main.go:84] event: MODIFIED, resourceVersion: 3300
    I0913 08:02:24.244502   92277 main.go:84] event: DELETED, resourceVersion: 3301
    ```

Reference: [example](https://github.com/kubernetes/client-go/blob/v0.37.1/examples/workqueue/main.go)

The timestamps and resource versions above illustrate an event sequence, not fixed output. The sample handles one List + Watch session. Use a Reflector or shared informer for continuous recovery after disconnections. Resource versions are opaque values, not counters to increment.
