# [Indexer](https://pkg.go.dev/k8s.io/client-go/tools/cache#Indexer)

## Interface

The `Indexer` interface:

```go
type Indexer interface {
	Store
	// Index returns the stored objects whose set of indexed values
	// intersects the set of indexed values of the given object, for
	// the named index
	Index(indexName string, obj interface{}) ([]interface{}, error)
	// IndexKeys returns the storage keys of the stored objects whose
	// set of indexed values for the named index includes the given
	// indexed value
	IndexKeys(indexName, indexedValue string) ([]string, error)
	// ListIndexFuncValues returns all the indexed values of the given index
	ListIndexFuncValues(indexName string) []string
	// ByIndex returns the stored objects whose set of indexed values
	// for the named index includes the given indexed value
	ByIndex(indexName, indexedValue string) ([]interface{}, error)
	// GetIndexer return the indexers
	GetIndexers() Indexers

	// AddIndexers adds more indexers to this store.  If you call this after you already have data
	// in the store, the results are undefined.
	AddIndexers(newIndexers Indexers) error
}
```

Indexer extends `Store` with multiple indices.

The `Store` interface:

```go
type Store interface {

	// Add adds the given object to the accumulator associated with the given object's key
	Add(obj interface{}) error

	// Update updates the given object in the accumulator associated with the given object's key
	Update(obj interface{}) error

	// Delete deletes the given object from the accumulator associated with the given object's key
	Delete(obj interface{}) error

	// List returns a list of all the currently non-empty accumulators
	List() []interface{}

	// ListKeys returns a list of all the keys currently associated with non-empty accumulators
	ListKeys() []string

	// Get returns the accumulator associated with the given object's key
	Get(obj interface{}) (item interface{}, exists bool, err error)

	// GetByKey returns the accumulator associated with the given key
	GetByKey(key string) (item interface{}, exists bool, err error)

	// Replace will delete the contents of the store, using instead the
	// given list. Store takes ownership of the list, you should not reference
	// it after calling this function.
	Replace([]interface{}, string) error

	// Resync is meaningless in the terms appearing here but has
	// meaning in some implementations that have non-trivial
	// additional behavior (e.g., DeltaFIFO).
	Resync() error
}
```

## Usage

1. Create a indexer with `KeyFunc` and `Indexers`.

    ```go
	indexer := cache.NewIndexer(
		cache.MetaNamespaceKeyFunc, // Use <namespace>/<name> as a key if <namespace> exists, otherwise <name>
		cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, // default index function that indexes based on an object's namespace
	)
    ```

    Arguments:
    1. `KeyFunc`: KeyFunc knows how to make **a key from an object**.

        ```go
        type KeyFunc func(obj interface{}) (string, error)
        ```

        Example KeyFunc: [cache.MetaNamespaceKeyFunc](https://pkg.go.dev/k8s.io/client-go/tools/cache#MetaNamespaceKeyFunc) (e.g. object: Deployment with name `test` in `default` namespace -> key: `default/test`)

    1. `Indexers`: Indexers maps a name to an IndexFunc.

        We can have multiple indexes in a store (indexer). e.g. index by namespace, index by label, etc.

        ```go
        type Indexers map[string]IndexFunc
        ```

    Dependencies:
    1. `IndexFunc`: IndexFunc knows how to compute **the set of indexed values for an object**. This function determines which value to use for indexing.

        ```go
        type IndexFunc func(obj interface{}) ([]string, error)
        ```

        Example: [cache.MetaNamespaceIndexFunc](https://pkg.go.dev/k8s.io/client-go/tools/cache#MetaNamespaceIndexFunc) (e.g. object: Deployment with name `test` in `default` namespace -> set of indexed values: `[]string{"default"}`)

1. Set an object to the indexer.

    Indexer store the object with indexes based on the configured IndexFunc for each indexer.

    ```go
    err := indexer.Add(&appsv1.Deployment{...})
    ```

1. Get objects from the indexer.

    ```go
    // List indexer
	objs := indexer.List()
	fmt.Printf("indexer.List got %d objects\n", len(objs))
    ```

    ```go
	// Get object by key
	obj, exists, err := indexer.GetByKey("default/test")
    ```
