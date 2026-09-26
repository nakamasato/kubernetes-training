# [cluster](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.25.1/pkg/cluster/cluster.go)

Cluster owns the dependencies used to interact with one Kubernetes cluster. `manager.New` constructs it and embeds the `Cluster` interface. Reconcilers receive these dependencies explicitly through their fields or constructors.

These implementation notes target controller-runtime v0.25.1, pinned in the repository's go.mod. Use the signatures and call paths below for this version.

## Types

### 1. Cluster interface

`GetClient` provides cache-backed reads and direct writes. `GetAPIReader` bypasses the cache. `GetFieldIndexer` and `GetCache` share the same underlying cache.

```go
type Cluster interface {
    recorder.Provider

    GetHTTPClient() *http.Client

    GetConfig() *rest.Config

    GetCache() cache.Cache

    GetScheme() *runtime.Scheme

    GetClient() client.Client

    GetFieldIndexer() client.FieldIndexer

    GetRESTMapper() meta.RESTMapper

    GetAPIReader() client.Reader

    Start(ctx context.Context) error
}
```

### 2. cluster struct

The scheme maps Go types to GVKs; the RESTMapper maps those GVKs to API resources and their namespace scope. The shared HTTP client carries transport/authentication configuration. Event recorders publish Kubernetes Events rather than application logs.

```go
type cluster struct {
    config *rest.Config

    httpClient *http.Client
    scheme     *runtime.Scheme
    cache      cache.Cache
    client     client.Client

    apiReader client.Reader

    fieldIndexes client.FieldIndexer

    recorderProvider *intrec.Provider

    mapper meta.RESTMapper

    logger logr.Logger
}
```

## New

Follow [cluster.New](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.25.1/pkg/cluster/cluster.go) in this order:

1. Copy the REST config and call `setOptionsDefaults`. Errors creating an HTTP client stop construction.
2. Create the RESTMapper with `options.MapperProvider(config, options.HTTPClient)`. `apiutil.NewDynamicRESTMapper` requires a real HTTP client; passing nil is not a working default.
3. Fill `options.Cache` with Scheme, Mapper, and HTTPClient when omitted, then call `options.NewCache(config, cacheOpts)`.
4. Fill `options.Client` with the same dependencies. Ensure `clientOpts.Cache` exists and set its Reader to the cache unless the caller supplied one. `options.NewClient(config, clientOpts)` creates the read/write client.
5. Create the APIReader with `client.New` and no Cache option. Its Get/List go directly to the API server.
6. Create the recorder provider, then return the cluster with `fieldIndexes` pointing to the cache.

The relevant client wiring is:

```go
if clientOpts.Cache == nil {
    clientOpts.Cache = &client.CacheOptions{Unstructured: false}
}
if clientOpts.Cache.Reader == nil {
    clientOpts.Cache.Reader = cache
}
clientWriter, err := options.NewClient(config, clientOpts)
```

All errors must propagate to the Manager constructor. Neither constructing the cluster nor obtaining its client starts informers: the Manager starts the cache later. See [Cache](../cache/) for informer creation and synchronization, and [Client](../client/) for read routing.

## SetOptionDefaults

| Option | Default | Used for |
| --- | --- | --- |
| Scheme | client-go's `scheme.Scheme` | Built-in Kubernetes types; register custom types yourself |
| HTTPClient | `rest.HTTPClientFor(config)` | Shared API transport |
| MapperProvider | `apiutil.NewDynamicRESTMapper` | Discovery and REST mappings |
| NewClient | `client.New` | Cache-backed read/write client |
| NewCache | `cache.New` | Informer cache |
| newRecorderProvider | `intrec.NewProvider` | Event recorders |
| Logger | `logf.RuntimeLog.WithName("cluster")` | Cluster logging |

The broadcaster defaults support both the legacy record broadcaster and the events broadcaster. The defaulting implementation is:

```go
func setOptionsDefaults(options Options, config *rest.Config) (Options, error) {
    if options.HTTPClient == nil {
        var err error
        options.HTTPClient, err = rest.HTTPClientFor(config)
        if err != nil {
            return options, err
        }
    }

    if options.Scheme == nil {
        options.Scheme = scheme.Scheme
    }

    if options.MapperProvider == nil {
        options.MapperProvider = apiutil.NewDynamicRESTMapper
    }

    if options.NewClient == nil {
        options.NewClient = client.New
    }

    if options.NewCache == nil {
        options.NewCache = cache.New
    }

    if options.newRecorderProvider == nil {
        options.newRecorderProvider = intrec.NewProvider
    }

    evtCl, err := eventsv1client.NewForConfigAndClient(config, options.HTTPClient)
    if err != nil {
        return options, err
    }

    if options.EventBroadcaster == nil {
        options.makeBroadcaster = func() (record.EventBroadcaster, events.EventBroadcaster, bool) {
            return record.NewBroadcaster(), events.NewBroadcaster(&events.EventSinkImpl{Interface: evtCl}), true
        }
    } else {
        options.makeBroadcaster = func() (record.EventBroadcaster, events.EventBroadcaster, bool) {
            return options.EventBroadcaster, events.NewBroadcaster(&events.EventSinkImpl{Interface: evtCl}), false
        }
    }

    if options.Logger.GetSink() == nil {
        options.Logger = logf.RuntimeLog.WithName("cluster")
    }

    return options, nil
}
```

## SetFields migration: explicit dependency wiring

Older versions called `cluster.SetFields` from `manager.SetFields`, and Manager.Add injected dependencies into controllers, reconcilers, sources, handlers, and predicates. Those APIs have been removed. The relationship still exists, but dependencies are passed explicitly:

| Former injection | Current wiring |
| --- | --- |
| ConfigInto | `mgr.GetConfig()` passed to a constructor |
| ClientInto | `Client: mgr.GetClient()` in the reconciler |
| APIReaderInto | `APIReader: mgr.GetAPIReader()` where uncached reads are needed |
| SchemeInto | `Scheme: mgr.GetScheme()` |
| CacheInto for Kind | `source.Kind(mgr.GetCache(), object, handler)` |
| MapperInto for owner handler | `handler.EnqueueRequestForOwner(mgr.GetScheme(), mgr.GetRESTMapper(), owner)` |

```go
r := &MyReconciler{Client: mgr.GetClient(), Scheme: mgr.GetScheme()}
if err := ctrl.NewControllerManagedBy(mgr).For(&corev1.Pod{}).Complete(r); err != nil {
    return err
}
```

`MyReconciler` is your type implementing `Reconcile`. `Manager.Add` now registers lifecycle work; it does not populate arbitrary fields on that work. The [inject page](../inject/) retains the historical interfaces and shows the replacement patterns.
