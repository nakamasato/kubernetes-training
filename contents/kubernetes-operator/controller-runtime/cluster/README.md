# [cluster](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.12.3/pkg/cluster/cluster.go)

[Cluster](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.12.3/pkg/cluster/cluster.go) provides various methods to interact with a cluster. Cluster is initialized and stored in [Manager](../manager/) with [cluster.New](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.12.3/pkg/cluster/cluster.go#L146).

Most of the fields in a cluster (scheme, cache, client, apiReader, recorderProvider, etc.) are used to injected to related components (Controller, EventHandlers, Sources, Predicates)

## types

### 1. [Cluster](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.12.3/pkg/cluster/cluster.go#L39) interface

```go
type Cluster interface {
	SetFields(interface{}) error
	GetConfig() *rest.Config
	GetScheme() *runtime.Scheme
	GetClient() client.Client
	GetFieldIndexer() client.FieldIndexer
	GetCache() cache.Cache
	GetEventRecorderFor(name string) record.EventRecorder
	GetRESTMapper() meta.RESTMapper
	GetAPIReader() client.Reader
	Start(ctx context.Context) error
}
```

### 2. [cluster](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.12.3/pkg/cluster/internal.go#L34) struct

```go
type cluster struct {
	config *rest.Config
	scheme *runtime.Scheme // scheme is injected into Controllers, EventHandlers, Sources and Predicates.
	cache cache.Cache // injected is injected into Sources
	client client.Client // client is injected into Controllers (and EventHandlers, Sources and Predicates).
	apiReader client.Reader // apiReader is the reader that will make requests to the api server and not the cache.
	fieldIndexes client.FieldIndexer
	recorderProvider *intrec.Provider // recorderProvider is used to generate event recorders that will be injected into Controllers (and EventHandlers, Sources and Predicates).
	mapper meta.RESTMapper // mapper is used to map resources to kind, and map kind and version.
	logger logr.Logger
}
```


## Set fields

1. scheme: Use the Kubernetes client-go scheme if none is specified
1. mapper: Created with `MapperProvider`
    The following function is used if `MapperProvider` is not specified:
    ```go
    apiutil.NewDynamicRESTMapper(c)
    ```
1. cache: Created with `NewCache` (cache.New)
    ```go
    cache, err := options.NewCache(config, cache.Options{Scheme: options.Scheme, Mapper: mapper, Resync: options.SyncPeriod, Namespace: options.Namespace})
    ```
1. apiReader: Created with `client.New`
    ```go
    apiReader, err := client.New(config, clientOptions)
    ```
1. writeObj: Created with `NewClient` ([DefaultNewClient](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.12.3/pkg/cluster/cluster.go#L259) -> [NewDelegatingClient](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.12.3/pkg/client/split.go#L44))
    ```go
    writeObj, err := options.NewClient(cache, config, clientOptions, options.ClientDisableCacheFor...)
    ```

    <details>

    ```go
	if options.NewClient == nil {
		options.NewClient = DefaultNewClient
	}
    ```

    ```go
    // DefaultNewClient creates the default caching client.
    func DefaultNewClient(cache cache.Cache, config *rest.Config, options client.Options, uncachedObjects ...client.Object) (client.Client, error) {
        c, err := client.New(config, options)
        if err != nil {
            return nil, err
        }

        return client.NewDelegatingClient(client.NewDelegatingClientInput{
            CacheReader:     cache,
            Client:          c,
            UncachedObjects: uncachedObjects,
        })
    }
    ```

    ```go
    &delegatingClient{
		scheme: in.Client.Scheme(),
		mapper: in.Client.RESTMapper(),
		Reader: &delegatingReader{
			CacheReader:       in.CacheReader,
			ClientReader:      in.Client,
			scheme:            in.Client.Scheme(),
			uncachedGVKs:      uncachedGVKs,
			cacheUnstructured: in.CacheUnstructured,
		},
		Writer:       in.Client,
		StatusClient: in.Client,
	}
    ```

    </details>
1. recorderProvider: Created with `newRecorderProvider` (intrec.NewProvider)
    ```go
    recorderProvider, err := options.newRecorderProvider(config, options.Scheme, options.Logger.WithName("events"), options.makeBroadcaster)
    ```


```go
return &cluster{
    config:           config,
    scheme:           options.Scheme,
    cache:            cache,
    fieldIndexes:     cache,
    client:           writeObj, // client is a delegatingClient.
    apiReader:        apiReader,
    recorderProvider: recorderProvider,
    mapper:           mapper,
    logger:           options.Logger,
}, nil
```

`GetClient()` returns `cluster.client`, a delegatingCLient by default. For more details about `delegatingClient` you can check [client](../client)
