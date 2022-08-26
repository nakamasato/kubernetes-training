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


Default:
1. Scheme: Use the Kubernetes client-go scheme if none is specified
1. MapperProvider: apiutil.NewDynamicRESTMapper(c)
1. NewClient: DefaultNewClient -> [NewDelegatingClient](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.12.3/pkg/client/split.go#L44) ->
    ```go
    // A delegating client forms a Client by composing separate reader, writer and
    // statusclient interfaces.  This way, you can have an Client that reads from a
    // cache and writes to the API server.
    client.NewDelegatingClient(client.NewDelegatingClientInput{
		CacheReader:     cache,
		Client:          c,
		UncachedObjects: uncachedObjects,
	})
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

1. NewCache: cache.New
1. newRecorderProvider: intrec.NewProvider
