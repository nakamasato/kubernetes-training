package main

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
	. "sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var (
	// Build webhooks used for the various server
	// configuration options
	//
	// These handlers could be also be implementations
	// of the AdmissionHandler interface for more complex
	// implementations.
	mutatingHook = &Admission{
		Handler: admission.HandlerFunc(func(ctx context.Context, req AdmissionRequest) AdmissionResponse {
			return Patched("some changes",
				JSONPatchOp{Operation: "add", Path: "/metadata/annotations/access", Value: "granted"},
				JSONPatchOp{Operation: "add", Path: "/metadata/annotations/reason", Value: "not so secret"},
			)
		}),
	}

	validatingHook = &Admission{
		Handler: admission.HandlerFunc(func(ctx context.Context, req AdmissionRequest) AdmissionResponse {
			return Denied("none shall pass!")
		}),
	}
)

func main() {
	// Create a manager
	// Note: GetConfigOrDie will os.Exit(1) w/o any message if no kube-config can be found
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{})
	if err != nil {
		panic(err)
	}

	// Create a webhook server.
	hookServer := NewServer(Options{Port: 8443})
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
}
