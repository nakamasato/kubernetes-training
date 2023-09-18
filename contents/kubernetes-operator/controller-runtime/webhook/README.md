# [webhook](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/webhook) (WIP)

> Package webhook provides methods to build and bootstrap a webhook server.

```go
type Server interface {
	// NeedLeaderElection implements the LeaderElectionRunnable interface, which indicates
	// the webhook server doesn't need leader election.
	NeedLeaderElection() bool

	// Register marks the given webhook as being served at the given path.
	// It panics if two hooks are registered on the same path.
	Register(path string, hook http.Handler)

	// Start runs the server.
	// It will install the webhook related resources depend on the server configuration.
	Start(ctx context.Context) error

	// StartedChecker returns an healthz.Checker which is healthy after the
	// server has been started.
	StartedChecker() healthz.Checker

	// WebhookMux returns the servers WebhookMux
	WebhookMux() *http.ServeMux
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
hookServer := NewServer(Options{
	Port: 8443,
})
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

## Run

```
go run contents/kubernetes-operator/controller-runtime/webhook/main.go
panic: open /var/folders/c2/hjlk2kcn63s4kds9k2_ctdhc0000gp/T/k8s-webhook-server/serving-certs/tls.crt: no such file or directory

goroutine 1 [running]:
main.main()
        /Users/m.naka/repos/nakamasato/kubernetes-training/contents/kubernetes-operator/controller-runtime/webhook/main.go:56 +0x1c4
exit status 2
```

## Changes

1. [v0.15.0](https://github.com/kubernetes-sigs/controller-runtime/releases/tag/v0.15.0)
	1. [Allow passing a custom webhook server controller-runtime#2293](https://github.com/kubernetes-sigs/controller-runtime/pull/2293) `webhook.Server` `struct` was changed to `interface`.
		```diff
		- hookServer := &Server{Port: 8443}
		+ hookServer := NewServer(Options{Port: 8443})
		```
