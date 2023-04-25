# [webhook](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/webhook) (WIP)

> Package webhook provides methods to build and bootstrap a webhook server.

```go
type Webhook struct {
	// Handler actually processes an admission request returning whether it was allowed or denied,
	// and potentially patches to apply to the handler.
	Handler Handler

	// RecoverPanic indicates whether the panic caused by webhook should be recovered.
	RecoverPanic bool

	// WithContextFunc will allow you to take the http.Request.Context() and
	// add any additional information such as passing the request path or
	// headers thus allowing you to read them from within the handler
	WithContextFunc func(context.Context, *http.Request) context.Context
	// contains filtered or unexported fields
}
```

```go
type Defaulter interface {
	runtime.Object
	Default()
}
```

```go
type Validator interface {
	runtime.Object
	ValidateCreate() error
	ValidateUpdate(old runtime.Object) error
	ValidateDelete() error
}
```

Example:

```go
// Create a manager
// Note: GetConfigOrDie will os.Exit(1) w/o any message if no kube-config can be found
mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{})
if err != nil {
	panic(err)
}

// Create a webhook server.
hookServer := &Server{
	Port: 8443,
}
if err := mgr.Add(hookServer); err != nil {
	panic(err)
}

// Register the webhooks in the server.
hookServer.Register("/mutating", mutatingHook)
hookServer.Register("/validating", validatingHook)

// Start the server by starting a previously-set-up manager
err = mgr.Start(ctrl.SetupSignalHandler())
if err != nil {
	// handle error
	panic(err)
}
```
