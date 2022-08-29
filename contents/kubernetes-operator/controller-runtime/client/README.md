# [client](https://github.com/kubernetes-sigs/controller-runtime/tree/v0.12.3/pkg/client/client.go)

## [Client interface](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.12.3/pkg/client/interfaces.go#L101)

```go
// Client knows how to perform CRUD operations on Kubernetes objects.
type Client interface {
	Reader
	Writer
	StatusClient

	// Scheme returns the scheme this client is using.
	Scheme() *runtime.Scheme
	// RESTMapper returns the rest this client is using.
	RESTMapper() meta.RESTMapper
}
```

## [delegatingClient](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.12.3/pkg/client/split.go#L69)

```go
type delegatingClient struct {
	Reader
	Writer
	StatusClient

	scheme *runtime.Scheme
	mapper meta.RESTMapper
}
```

There's a function called [shouldBypassCache](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.12.3/pkg/client/split.go#L102) to check if the target object is cached or not. If cached, call [cacheReader](), otherwise call [clientReader]()


## How to use `client`

```go
func newClient(config *rest.Config, options Options) (*client, error) {
```

```go
c := &client{
    typedClient: typedClient{
        cache:      clientcache,
        paramCodec: runtime.NewParameterCodec(options.Scheme),
    },
    unstructuredClient: unstructuredClient{
        cache:      clientcache,
        paramCodec: noConversionParamCodec{},
    },
    metadataClient: metadataClient{
        client:     rawMetaClient,
        restMapper: options.Mapper,
    },
    scheme: options.Scheme,
    mapper: options.Mapper,
}
```

## How to use `delegatingClient`


Initialized in [NewDelegatingClient](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.12.3/pkg/client/split.go#L44)

Input:

```go
// NewDelegatingClientInput encapsulates the input parameters to create a new delegating client.
type NewDelegatingClientInput struct {
	CacheReader       Reader
	Client            Client
	UncachedObjects   []Object
	CacheUnstructured bool
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

[cacheReader](https://github.com/kubernetes-sigs/controller-runtime/blob/f46919744bee01060c9084a285e049afffd38c9d/pkg/cache/internal/cache_reader.go#L40):

```go
// CacheReader wraps a cache.Index to implement the client.CacheReader interface for a single type.
type CacheReader struct {
	// indexer is the underlying indexer wrapped by this cache.
	indexer cache.Indexer

	// groupVersionKind is the group-version-kind of the resource.
	groupVersionKind schema.GroupVersionKind

	// scopeName is the scope of the resource (namespaced or cluster-scoped).
	scopeName apimeta.RESTScopeName

	// disableDeepCopy indicates not to deep copy objects during get or list objects.
	// Be very careful with this, when enabled you must DeepCopy any object before mutating it,
	// otherwise you will mutate the object in the cache.
	disableDeepCopy bool
}
```
