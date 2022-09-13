# inject

## Inject interface table

|interface name|required func|func to inject|Implemented by|
|---|---|---|---|
|Cache|InjectCache|CacheInto|Kind (source)|
|APIReader|InjectAPIReader|APIReaderInto||
|Config|InjectConfig|ConfigInto||
|Client|InjectClient|ClientInto||
|Scheme|InjectScheme|SchemeInto|DeferredFileLoader, Webhook|
|Stoppable|InjectStopChannel|StopChannelInto|Channel (source)|
|Mapper|InjectMapper|MapperInto|EnqueueRequestForOwner|
|Injector|InjectFunc|InjectorInto|Webhook, enqueueRequestsFromMapFunc, Controller, and, or (predicate), multiMutating, multiValidating, etc.|
|Logger|InjectLogger|LoggerInto|Webhook|

`Injector` interface

```go
// Func injects dependencies into i.
type Func func(i interface{}) error

// Injector is used by the ControllerManager to inject Func into Controllers.
type Injector interface {
	InjectFunc(f Func) error
}

// InjectorInto will set f and return the result on i if it implements Injector.  Returns
// false if i does not implement Injector.
func InjectorInto(f Func, i interface{}) (bool, error) {
	if ii, ok := i.(Injector); ok {
		return true, ii.InjectFunc(f)
	}
	return false, nil
}
```

`SetFields` set dependencies to the object that implments inject interface.

## Usage

Controller implement Injector with [InjectFunc](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.13.0/pkg/internal/controller/controller.go#L352)

```go
// InjectFunc implement SetFields.Injector.
func (c *Controller) InjectFunc(f inject.Func) error {
	c.SetFields = f
	return nil
}
```

With this implementation, any function can be injected to the `controller.SetFields` with `InjectorInto(func, controller)`.

This function is used in the [manager.SetFields](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.13.0/pkg/manager/internal.go#L187-L211)

```go
if _, err := inject.InjectorInto(cm.SetFields, i); err != nil {
```
This means set `clusterManager`'s `SetFields` function to `i`, specifically, `controller`. By the controller's `InjectFunc` implementation, **controller has exactly the same `SetFields` function as `clusterManager`**
