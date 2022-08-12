# [DeltaFifo]()

## Overview

**DeltaFIFO** is a producer-consumer queue, where a [Reflector](../reflector) is intended to be the producer, and the consumer is whatever calls the Pop() method.

1. The actual data is stored in `items` in the for of `map[string]Deltas`.
1. The order is stored in `queue` as `[]string`.

```go
type DeltaFIFO struct {
	lock sync.RWMutex
	cond sync.Cond
	items map[string]Deltas
	queue []string
	populated bool
	initialPopulationCount int
	keyFunc KeyFunc
	knownObjects KeyListerGetter
	closed bool
	emitDeltaTypeReplaced bool
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

## Usage
