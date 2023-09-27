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

## Validator

[Validator](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.2/pkg/webhook/admission#Validator) interface:

```go
type Validator interface {
	runtime.Object

	// ValidateCreate validates the object on creation.
	// The optional warnings will be added to the response as warning messages.
	// Return an error if the object is invalid.
	ValidateCreate() (warnings Warnings, err error)

	// ValidateUpdate validates the object on update. The oldObj is the object before the update.
	// The optional warnings will be added to the response as warning messages.
	// Return an error if the object is invalid.
	ValidateUpdate(old runtime.Object) (warnings Warnings, err error)

	// ValidateDelete validates the object on deletion.
	// The optional warnings will be added to the response as warning messages.
	// Return an error if the object is invalid.
	ValidateDelete() (warnings Warnings, err error)
}
```

The return value was updated in [controller-runtime@v0.15.0](https://github.com/kubernetes-sigs/controller-runtime/releases/tag/v0.15.0) ([⚠️ feat: new features about support warning with webhook #2014](https://github.com/kubernetes-sigs/controller-runtime/pull/2014)) from [[Feature Request]: Support "Warning" for Validation Webhook #1896](https://github.com/kubernetes-sigs/controller-runtime/issues/1896)

This is because Kubernets supports `warning` message in response for Admission webhook [ref](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#response) since 1.19:

> Admission webhooks can optionally return warning messages that are returned to the requesting client in HTTP Warning headers with a warning code of 299. Warnings can be sent with allowed or rejected admission responses.

## Changes

1. [v0.15.0](https://github.com/kubernetes-sigs/controller-runtime/releases/tag/v0.15.0)
	1. [Allow passing a custom webhook server controller-runtime#2293](https://github.com/kubernetes-sigs/controller-runtime/pull/2293) `webhook.Server` `struct` was changed to `interface`.
		```diff
		- hookServer := &Server{Port: 8443}
		+ hookServer := NewServer(Options{Port: 8443})
		```
	1. [⚠️ feat: new features about support warning with webhook #2014](https://github.com/kubernetes-sigs/controller-runtime/pull/2014) `Validator`, `CustomValidator` interface change: added warning to response of admission webhook.
		```diff
		type Validator interface {
			runtime.Object
		- 	ValidateCreate() error
		- 	ValidateUpdate(old runtime.Object) error
		- 	ValidateDelete() error

		+	// ValidateCreate validates the object on creation.
		+ 	// The optional warnings will be added to the response as warning messages.
		+ 	// Return an error if the object is invalid.
		+ 	ValidateCreate() (warnings Warnings, err error)
		+ 	// ValidateUpdate validates the object on update. The oldObj is the object before the update.
		+ 	// The optional warnings will be added to the response as warning messages.
		+ 	// Return an error if the object is invalid.
		+ 	ValidateUpdate(old runtime.Object) (warnings Warnings, err error)
		+ 	// ValidateDelete validates the object on deletion.
		+ 	// The optional warnings will be added to the response as warning messages.
		+ 	// Return an error if the object is invalid.
		+ 	ValidateDelete() (warnings Warnings, err error)
		}
		```
