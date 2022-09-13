# [ListerWatcher](https://github.com/kubernetes/client-go/blob/v0.25.0/tools/cache/listwatch.go#L43)

## type

### Interface

```go
// ListerWatcher is any object that knows how to perform an initial list and start a watch on a resource.
type ListerWatcher interface {
	Lister
	Watcher
}
```

```go
// Lister is any object that knows how to perform an initial list.
type Lister interface {
	// List should return a list type object; the Items field will be extracted, and the
	// ResourceVersion field will be used to start the watch in the right place.
	List(options metav1.ListOptions) (runtime.Object, error)
}

// Watcher is any object that knows how to start a watch on a resource.
type Watcher interface {
	// Watch should begin a watch at the specified version.
	Watch(options metav1.ListOptions) (watch.Interface, error)
}
```

### ListWatch struct

```go
// ListFunc knows how to list resources
type ListFunc func(options metav1.ListOptions) (runtime.Object, error)

// WatchFunc knows how to watch resources
type WatchFunc func(options metav1.ListOptions) (watch.Interface, error)

// ListWatch knows how to list and watch a set of apiserver resources.  It satisfies the ListerWatcher interface.
// It is a convenience function for users of NewReflector, etc.
// ListFunc and WatchFunc must not be nil
type ListWatch struct {
	ListFunc  ListFunc
	WatchFunc WatchFunc
	// DisableChunking requests no chunking for this list watcher.
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

1. Created with [NewFilteredListWatchFromClient](https://github.com/kubernetes/client-go/blob/v0.25.0/tools/cache/listwatch.go#L80):

	```go
	func NewFilteredListWatchFromClient(c Getter, resource string, namespace string, optionsModifier func(options *metav1.ListOptions)) *ListWatch {
		listFunc := func(options metav1.ListOptions) (runtime.Object, error) {
			optionsModifier(&options)
			return c.Get().
				Namespace(namespace).
				Resource(resource).
				VersionedParams(&options, metav1.ParameterCodec).
				Do(context.TODO()).
				Get()
		}
		watchFunc := func(options metav1.ListOptions) (watch.Interface, error) {
			options.Watch = true
			optionsModifier(&options)
			return c.Get().
				Namespace(namespace).
				Resource(resource).
				VersionedParams(&options, metav1.ParameterCodec).
				Watch(context.TODO())
		}
		return &ListWatch{ListFunc: listFunc, WatchFunc: watchFunc}
	}
	```

1. Getter can be got from `clientset`: `clientset.CoreV1().RESTClient()`

	```go
	podListWatcher := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "pods", v1.NamespaceDefault, fields.Everything())
	```
1. listFunc and watchFunc is set in `NewFilteredListWatchFromClient`:
	1. `c.Get()` calls `NewRequest().Verb("GET")`
	1. [rest.NewRequest](https://github.com/kubernetes/client-go/blob/v0.25.0/rest/request.go#L126) creates and returns a request.
	1. [Request.Watch](https://github.com/kubernetes/client-go/blob/v0.25.0/rest/request.go#L607)
	    1. get retry func
	    1. in a for loop, run `retry.Before`, `client.Do(req)`, and `retry.After`

1. Then, listwatcher will be passed to informer.
	```go
	indexer, informer := cache.NewIndexerInformer(podListWatcher, &v1.Pod{}, 0, cache.ResourceEventHandlerFuncs{...
	```
	`ListWatcher` is used to initialize the store (FIFODelta) with `List` and keep the store up-to-date with `Watch`. For more details, you can check [informer](../informer/README.md)


## Example

Code: You can get the event through `w.ResultChan()`. The example is simplified version (no error handling).

```go
	w, err := podListWatcher.Watch(metav1.ListOptions{}) // returns watch.Interface
	if err != nil {
		klog.Fatal(err)
	}
loop:
	for {
		select {
		case event, ok := <-w.ResultChan():
			if !ok {
				break loop
			}

			meta, err := meta.Accessor(event.Object)
			if err != nil {
				continue
			}
			resourceVersion := meta.GetResourceVersion()
			klog.Infof("event: %s, resourceVersion: %s", event.Type, resourceVersion)
		}
	}
```

Run:

1. Start the ListerWatcher for Pods.
	```
	go run main.go
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

Referece: [example](https://github.com/kubernetes/client-go/blob/v0.25.0/examples/workqueue/main.go)
