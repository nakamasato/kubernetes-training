# [Reflector](https://github.com/kubernetes/client-go/blob/master/tools/cache/reflector.go)

## Overview

***Reflector*** watches a specified resource and causes all changes to be reflected in the given store.

```go
// Reflector watches a specified resource and causes all changes to be reflected in the given store.
type Reflector struct {
	name string
	expectedTypeName string
	expectedType reflect.Type
	expectedGVK *schema.GroupVersionKind
	store Store
	listerWatcher ListerWatcher
	backoffManager wait.BackoffManager
	initConnBackoffManager wait.BackoffManager
	resyncPeriod time.Duration
	ShouldResync func() bool
	clock clock.Clock
	paginatedResult bool
	lastSyncResourceVersion string
	isLastSyncResourceVersionUnavailable bool
	lastSyncResourceVersionMutex sync.RWMutex
	WatchListPageSize int64
	watchErrorHandler WatchErrorHandler
}
```

1. `store`: DeltaFIFO can be used for store.
1. `ListAndWatch` function is called in `Run`.
1. `ListAndWatch` function calls
    1. `list` func
    1. `w, err := listerwatcher.Watch()`
    1. `watchHandler(start, w, r.store, r.expectedType, r.expectedGVK, r.name, r.expectedTypeName, r.setLastSyncResourceVersion, r.clock, resyncerrc, stopCh)` -> `store.Add`, `store.Update`, `store.Delete`

## Usage
