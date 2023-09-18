# [cluster](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.13.0/pkg/cluster/cluster.go)

[Cluster](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.13.0/pkg/cluster/cluster.go) provides various methods to interact with a cluster. Cluster is initialized and stored in [Manager](../manager/) with [cluster.New](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.13.0/pkg/cluster/cluster.go#L146).

Most of the fields in a cluster (scheme, cache, client, apiReader, recorderProvider, etc.) are used to injected to related components (Controller, EventHandlers, Sources, Predicates)

## Types

### 1. [Cluster](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.13.0/pkg/cluster/cluster.go#L39) interface

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

### 2. [cluster](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.13.0/pkg/cluster/internal.go#L34) struct

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

## [New](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.13.0/pkg/cluster/cluster.go#L146)

1. [SetOptionDefaults](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.13.0//pkg/cluster/cluster.go#L208)
	For more details, check below

1. Create a `mapper`
    ```go
    mapper, err := options.MapperProvider(config)
    ```
1. Create a `cache` with `NewCache` ([cache.New](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.13.0/pkg/cache/cache.go#L148))
    ```go
    cache, err := options.NewCache(config, cache.Options{Scheme: options.Scheme, Mapper: mapper, Resync: options.SyncPeriod, Namespace: options.Namespace})
    ```

    For more details, read [cache](../cache/README.md)

1. Create `apiReader`
    ```go
    apiReader, err := client.New(config, clientOptions)
    ```
1. Create a `writeObj` with `NewClient` ([DefaultNewClient](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.13.0/pkg/cluster/cluster.go#L259) -> [NewDelegatingClient](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.13.0/pkg/client/split.go#L44))
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


1. Create a `recorderProvider`
    ```go
    recorderProvider, err := options.newRecorderProvider(config, options.Scheme, options.Logger.WithName("events"), options.makeBroadcaster)
    ```
1. Create cluster
    ```go
    &cluster{
		config:           config,
		scheme:           options.Scheme,
		cache:            cache,
		fieldIndexes:     cache,
		client:           writeObj,
		apiReader:        apiReader,
		recorderProvider: recorderProvider,
		mapper:           mapper,
		logger:           options.Logger,
	}
    ```

## [SetOptionDefaults](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.13.0/pkg/cluster/cluster.go#L208)


|name|value|where to use|
|---|---|---|
|Scheme|scheme.Scheme||
|MapperProvider|`func(c *rest.Config) (meta.RESTMapper, error) {return apiutil.NewDynamicRESTMapper(c, nil)}`||
|NewClient|DefaultNewClient|
|NewCache|cache.New|
|newRecorderProvider|intrec.NewProvider|
|makeBroadcaster|`func() (record.EventBroadcaster, bool) {return record.NewBroadcaster(), true}`|
|Logger|logf.RuntimeLog.WithName("cluster")|

1. `options.Scheme = scheme.Scheme`(Use the Kubernetes client-go scheme if none is specified)
1. MapperProvider
    ```go
    options.MapperProvider = func(c *rest.Config) (meta.RESTMapper, error) {
		return apiutil.NewDynamicRESTMapper(c, nil)
	}
    ```
1. `options.NewClient = DefaultNewClient`
    ```go
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

    `GetClient()` returns `cluster.client`, a delegatingClient by default. For more details about `delegatingClient` you can check [client](../client/README.md)

1. `options.NewCache = cache.New`
1. `options.newRecorderProvider = intrec.NewProvider`
1. `record.NewBroadcaster()`
1. `options.Logger = logf.RuntimeLog.WithName("cluster")`

## [SetFields](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.13.0/pkg/cluster/internal.go#L67)



```go
func (c *cluster) SetFields(i interface{}) error {
	if _, err := inject.ConfigInto(c.config, i); err != nil {
		return err
	}
	if _, err := inject.ClientInto(c.client, i); err != nil {
		return err
	}
	if _, err := inject.APIReaderInto(c.apiReader, i); err != nil {
		return err
	}
	if _, err := inject.SchemeInto(c.scheme, i); err != nil {
		return err
	}
	if _, err := inject.CacheInto(c.cache, i); err != nil {
		return err
	}
	if _, err := inject.MapperInto(c.mapper, i); err != nil {
		return err
	}
	return nil
}
```

1. `cluster.SetFields` is called in [manager.SetFields](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.13.0/pkg/manager/internal.go#L196-L211)
1. `cluster.SetFields` injects `Config`, `Client`, `APIReader`, `Scheme`, `Cache` and `Mapper` into the specified `i`.
1. [manager.SetFields](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.13.0/pkg/manager/internal.go#L196-L211)'s usage:
    1. used for reconciler passed via builder in [controller](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.13.0/pkg/controller/controller.go#L138)
        ```go
        // Inject dependencies into Reconciler
        if err := mgr.SetFields(options.Reconciler); err != nil {
            return nil, err
        }
        ```

    1. used for runnables added to the Manager with [add function](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.13.0/pkg/manager/internal.go#L187)
        ```go
        // Add sets dependencies on i, and adds it to the list of Runnables to start.
        func (cm *controllerManager) Add(r Runnable) error {
            cm.Lock()
            defer cm.Unlock()
            return cm.add(r)
        }

        func (cm *controllerManager) add(r Runnable) error {
            // Set dependencies on the object
            if err := cm.SetFields(r); err != nil {
                return err
            }
            return cm.runnables.Add(r)
        }
        ```
        1. Controller is passed in [controller.New](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.13.0/pkg/controller/controller.go#L95)
