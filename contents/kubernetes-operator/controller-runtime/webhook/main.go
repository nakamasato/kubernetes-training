package main

import (
	"context"
	"encoding/json"
	"flag"
	"net/http"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var (
	mutatingHook   = &admission.Webhook{Handler: admission.HandlerFunc(mutate)}
	validatingHook = &admission.Webhook{
		Handler: admission.HandlerFunc(func(ctx context.Context, req admission.Request) admission.Response {
			return admission.Denied("none shall pass!")
		}),
	}
)

func mutate(ctx context.Context, req admission.Request) admission.Response {
	obj := &unstructured.Unstructured{}
	if err := json.Unmarshal(req.Object.Raw, obj); err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}
	annotations := obj.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}
	annotations["access"] = "granted"
	annotations["reason"] = "not so secret"
	obj.SetAnnotations(annotations)
	updated, err := json.Marshal(obj)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}
	// Compute a patch that also works when metadata.annotations is absent.
	return admission.PatchResponseFromRaw(req.Object.Raw, updated)
}

func main() {
	certDir := flag.String("cert-dir", "certs", "directory containing tls.crt and tls.key")
	flag.Parse()
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))
	hookServer := webhook.NewServer(webhook.Options{Port: 8443, CertDir: *certDir})
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{WebhookServer: hookServer})
	if err != nil {
		panic(err)
	}
	mgr.GetWebhookServer().Register("/mutating", mutatingHook)
	mgr.GetWebhookServer().Register("/validating", validatingHook)
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		panic(err)
	}
}
