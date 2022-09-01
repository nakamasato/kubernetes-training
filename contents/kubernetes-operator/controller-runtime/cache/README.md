# [cache](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/cache)

![](diagram.drawio.svg)

## [New](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.12.3/pkg/cache/cache.go#L148)

1. [Cache.New](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.12.3/pkg/cache/cache.go#L148) initializes and returns informerCache.
    ```go
    im := internal.NewInformersMap(config, opts.Scheme, opts.Mapper, *opts.Resync, opts.Namespace, selectorsByGVK, disableDeepCopyByGVK, transformByGVK)
    return &informerCache{InformersMap: im}, nil
    ```
    1. `disableDeepCopyByGVK` is determined from `opts`:
        ```go
        disableDeepCopyByGVK, err := convertToDisableDeepCopyByGVK(opts.UnsafeDisableDeepCopyByObject, opts.Scheme)
        ```
    1. The returned value is the following type:
        ```go
        type DisableDeepCopyByGVK map[schema.GroupVersionKind]bool
        ```
    1. Default value is empty: `internal.DisableDeepCopyByGVK{}`
1. [informerCache](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.12.3/pkg/cache/informer_cache.go#L49)
    ```go
    type informerCache struct {
        *internal.InformersMap
    }
    ```
1. [InformersMap](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.12.3/pkg/cache/internal/deleg_map.go#L34)
    ```go
    type InformersMap struct {
        // we abstract over the details of structured/unstructured/metadata with the specificInformerMaps
        // TODO(directxman12): genericize this over different projections now that we have 3 different maps

        structured   *specificInformersMap
        unstructured *specificInformersMap
        metadata     *specificInformersMap

        // Scheme maps runtime.Objects to GroupVersionKinds
        Scheme *runtime.Scheme
    }
    ```

    Initialized with [NewInformersMap](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.12.3/pkg/cache/internal/deleg_map.go#L48):
    ```go
    &InformersMap{
        structured:   newStructuredInformersMap(config, scheme, mapper, resync, namespace, selectors, disableDeepCopy, transformers),
        unstructured: newUnstructuredInformersMap(config, scheme, mapper, resync, namespace, selectors, disableDeepCopy, transformers),
        metadata:     newMetadataInformersMap(config, scheme, mapper, resync, namespace, selectors, disableDeepCopy, transformers),

        Scheme: scheme,
    }
    ```
1. [All newXXXInformersMap calls specificInformersMap](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.12.3/pkg/cache/internal/deleg_map.go#L110-L126):
    ```go
    // newStructuredInformersMap creates a new InformersMap for structured objects.
    func newStructuredInformersMap(config *rest.Config, scheme *runtime.Scheme, mapper meta.RESTMapper, resync time.Duration,
        namespace string, selectors SelectorsByGVK, disableDeepCopy DisableDeepCopyByGVK, transformers TransformFuncByObject) *specificInformersMap {
        return newSpecificInformersMap(config, scheme, mapper, resync, namespace, selectors, disableDeepCopy, transformers, createStructuredListWatch)
    }

    // newUnstructuredInformersMap creates a new InformersMap for unstructured objects.
    func newUnstructuredInformersMap(config *rest.Config, scheme *runtime.Scheme, mapper meta.RESTMapper, resync time.Duration,
        namespace string, selectors SelectorsByGVK, disableDeepCopy DisableDeepCopyByGVK, transformers TransformFuncByObject) *specificInformersMap {
        return newSpecificInformersMap(config, scheme, mapper, resync, namespace, selectors, disableDeepCopy, transformers, createUnstructuredListWatch)
    }

    // newMetadataInformersMap creates a new InformersMap for metadata-only objects.
    func newMetadataInformersMap(config *rest.Config, scheme *runtime.Scheme, mapper meta.RESTMapper, resync time.Duration,
        namespace string, selectors SelectorsByGVK, disableDeepCopy DisableDeepCopyByGVK, transformers TransformFuncByObject) *specificInformersMap {
        return newSpecificInformersMap(config, scheme, mapper, resync, namespace, selectors, disableDeepCopy, transformers, createMetadataListWatch)
    }
    ```
1. [specificInformersMap](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.12.3/pkg/cache/internal/informers_map.go#L89)
    Initialized with [newSpecificInformersMap](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.12.3/pkg/cache/internal/informers_map.go#L50)
    ```go
    ip := &specificInformersMap{
        config:            config,
        Scheme:            scheme,
        mapper:            mapper,
        informersByGVK:    make(map[schema.GroupVersionKind]*MapEntry),
        codecs:            serializer.NewCodecFactory(scheme),
        paramCodec:        runtime.NewParameterCodec(scheme),
        resync:            resync,
        startWait:         make(chan struct{}),
        createListWatcher: createListWatcher,
        namespace:         namespace,
        selectors:         selectors.forGVK,
        disableDeepCopy:   disableDeepCopy,
        transformers:      transformers,
    }
    ```

1. [MapEntry](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.12.3/pkg/cache/internal/informers_map.go#L79)

    ```go
    // MapEntry contains the cached data for an Informer.
    type MapEntry struct {
        // Informer is the cached informer
        Informer cache.SharedIndexInformer

        // CacheReader wraps Informer and implements the CacheReader interface for a single type
        Reader CacheReader
    }
    ```
1. [CacheReader](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.12.3/pkg/cache/internal/cache_reader.go#L40)
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
1. `MapEntry` and `CacheReader` is initialized in [addInformerToMap](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.12.3/pkg/cache/internal/informers_map.go#L247):

    ```go
	i := &MapEntry{
		Informer: ni,
		Reader: CacheReader{
			indexer:          ni.GetIndexer(),
			groupVersionKind: gvk,
			scopeName:        rm.Scope.Name(),
			disableDeepCopy:  ip.disableDeepCopy.IsDisabled(gvk),
		},
	}
    ```

    `disableDeepCopy` is determined by [isDisabled()](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.12.3/pkg/cache/internal/disabledeepcopy.go#L28) originally given from `opts.UnsafeDisableDeepCopyByObject`.

    ```go
    // IsDisabled returns whether a GroupVersionKind is disabled DeepCopy.
    func (disableByGVK DisableDeepCopyByGVK) IsDisabled(gvk schema.GroupVersionKind) bool {
        if d, ok := disableByGVK[gvk]; ok {
            return d
        } else if d, ok = disableByGVK[GroupVersionKindAll]; ok {
            return d
        }
        return false
    }
    ```
    As you can see, if it's empty (the default value), it returns `false`, which means **when cacheReader gets an obejct from the cache, it internally deepcopies the object**. You can modified the got object without deepcopying by yourself.
