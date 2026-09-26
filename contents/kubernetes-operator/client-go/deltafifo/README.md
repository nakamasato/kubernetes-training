# [DeltaFIFO](https://pkg.go.dev/k8s.io/client-go/tools/cache#DeltaFIFO)

## Overview

**DeltaFIFO** is a producer-consumer queue, where a [Reflector](../reflector) is intended to be the producer, and the consumer is whatever calls the Pop() method.

1. The actual data is stored in `items` in the form of `map[string]Deltas`.
1. The order is stored in `queue` as `[]string`.

```go
type DeltaFIFO struct {
    logger klog.Logger

    name string

    lock sync.RWMutex
    cond sync.Cond

    items map[string]Deltas

    queue []string

    synced       chan struct{}
    syncedClosed bool

    populated bool
    initialPopulationCount int

    keyFunc KeyFunc

    knownObjects KeyListerGetter

    closed bool

    emitDeltaTypeReplaced bool

    transformer TransformFunc
}
```

```go
type Deltas []Delta
type DeltaType string
type Delta struct {
    Type   DeltaType
    Object interface{}
}
```

- Difference between DeltaFIFO and FIFO
    1. One is that the accumulator associated with a given object's key is not that object but rather a Deltas, which is a slice of Delta values for that object. Applying an object to a Deltas means to append a Delta except when the potentially appended Delta is a Deleted and the Deltas already ends with a Deleted. In that case the Deltas does not grow, although the terminal Deleted will be replaced by the new Deleted if the older Deleted's object is a DeletedFinalStateUnknown.

    1. The other difference is that DeltaFIFO has two additional ways that an object can be applied to an accumulator: Replaced and Sync. If EmitDeltaTypeReplaced is not set to true, Sync will be used in replace events for backwards compatibility. Sync is used for periodic resync events.
- DeltaFIFO solves this use case
    1. You want to process every object change (delta) at most once.
    1. When you process an object, you want to see everything that's happened to it since you last processed it.
    1. You want to process the deletion of some of the objects.
    1. You might want to periodically reprocess objects.

## Usage: how DeltaFIFO is used in informer

1. Create Indexer
1. Create DeltaFIFO
1. Call `fifo.Add(xx)` or `fifo.Update(xx)` or `fifo.Delete(xx)`
1. Call `fifo.Pop(process)` with `process` function `type PopProcessFunc func(interface{}, bool) error`, which converts the object into Deltas and process deltas with `processDeltas`.
    ```go
    func process(obj interface{}, isInInitialList bool) error { // type PopProcessFunc func(interface{}, bool) error
        if deltas, ok := obj.(cache.Deltas); ok {
            return processDeltas(deltas)
        }
        return errors.New("object given as Process argument is not Deltas")
    }
    ```
1. `processDeltas` updates/add/delete indexer.

For more details, you can check [informer](../informer/)

## Run

From the repository root, without a cluster:

```sh
go run ./contents/kubernetes-operator/client-go/deltafifo
```

The sample adds, updates, and deletes the same Pod, then processes its accumulated Deltas with one `Pop`. It prints `OnUpdate or OnAdd` twice and `OnDelete` once. Check errors from `Add`, `Update`, `Delete`, and `Pop`; close the queue when finished.

`isInInitialList` distinguishes initial-list items. The sample callback only logs; a real informer's callback updates the Indexer and notifies handlers. `KnownObjects` supplies existing state and does not update it automatically. An empty queue blocks `Pop` until work arrives or `Close` is called. Shared informers can use other queue implementations depending on version and feature gates; DeltaFIFO remains useful for studying changes grouped by key.
