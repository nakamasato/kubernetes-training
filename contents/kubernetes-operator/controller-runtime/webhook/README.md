# [webhook](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.25.1/pkg/webhook/server.go)

Webhook provides an HTTPS server and admission handlers. These notes target v0.25.1. Register paths on the Manager's configured WebhookServer so it participates in Manager startup and shutdown.

## Server interface

```go
type Server interface {
    NeedLeaderElection() bool

    Register(path string, hook http.Handler)

    Start(ctx context.Context) error

    StartedChecker() healthz.Checker

    WebhookMux() *http.ServeMux
}
```

The default server uses TLS, does not require leadership, and exposes a readiness checker. Registering an HTTP path does not create a Kubernetes MutatingWebhookConfiguration or ValidatingWebhookConfiguration; those are separate deployment resources.

## Defaulter and Validator interfaces

The current generic interfaces separate admission behavior from the object type:

```go
type Defaulter[T runtime.Object] interface {
    Default(ctx context.Context, obj T) error
}
```

```go
type Validator[T runtime.Object] interface {
    ValidateCreate(ctx context.Context, obj T) (warnings Warnings, err error)

    ValidateUpdate(ctx context.Context, oldObj, newObj T) (warnings Warnings, err error)

    ValidateDelete(ctx context.Context, obj T) (warnings Warnings, err error)
}
```

`CustomDefaulter` and `CustomValidator` alias the runtime.Object specializations. `WithDefaulter` / `WithValidator` construct admission.Webhook handlers. The lower-level `admission.Handler` lets you inspect Request and return a Response directly, as in this sample.

## Example

The [sample](main.go) creates `webhook.NewServer(webhook.Options{Port: 8443, CertDir: *certDir})`, passes it through `ctrl.Options{WebhookServer: hookServer}`, and registers `/mutating` and `/validating` through `mgr.GetWebhookServer()` before starting the Manager.

The mutating handler decodes to Unstructured so unknown fields survive, preserves existing annotations, initializes a missing annotations map, then computes a JSON Patch against the original request:

```go
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
    return admission.PatchResponseFromRaw(req.Object.Raw, updated)
}
```

The validating handler deliberately returns `admission.Denied("none shall pass!")`. This is an example response, not a useful production policy. Decode errors return an admission error; a successful mutation returns Allowed with patchType JSONPatch.

## Run

From the repository root, provide a kubeconfig and generate local test certificates in a temporary directory:

```sh
cert_dir=$(mktemp -d)
openssl req -x509 -newkey rsa:2048 -nodes \
  -keyout "$cert_dir/tls.key" -out "$cert_dir/tls.crt" \
  -days 1 -subj '/CN=localhost'
go run ./contents/kubernetes-operator/controller-runtime/webhook -cert-dir "$cert_dir"
```

In another terminal, submit a local AdmissionReview. `-k` is only for this temporary self-signed localhost certificate:

```sh
curl -sk https://localhost:8443/mutating \
  -H 'Content-Type: application/json' \
  -d '{"apiVersion":"admission.k8s.io/v1","kind":"AdmissionReview","request":{"uid":"demo","operation":"CREATE","kind":{"group":"","version":"v1","kind":"Pod"},"resource":{"group":"","version":"v1","resource":"pods"},"object":{"apiVersion":"v1","kind":"Pod","metadata":{"name":"demo"},"spec":{"containers":[{"name":"nginx","image":"nginx"}]}}}}'
```

Expect response.allowed=true and a base64-encoded JSON patch adding `access=granted` and `reason=not so secret` annotations. Send the same body to `/validating` to see allowed=false and the denial reason. HTTP success alone does not mean admission was allowed; inspect response.allowed.

Ctrl+C stops the server. Remove the temporary certificate directory afterward. For in-cluster use, provide a Service, a certificate valid for its DNS name, a trusted caBundle, and the appropriate webhook configuration/rules. The local HTTPS check does not test API-server routing to a deployed webhook.

## Validator

ValidateCreate receives the new object, ValidateUpdate receives old and new objects, and ValidateDelete receives the object being removed. They return warnings plus an error; errors deny the request, while warnings can accompany allowed or denied responses. Validation should not mutate objects.

The pinned controller-runtime version supports returning warnings with a validation result. Kubernetes can return these to the caller as HTTP Warning headers; see the [admission response documentation](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#response).

Run the mutation regression tests without a cluster:

```sh
go test ./contents/kubernetes-operator/controller-runtime/webhook
```

They exercise absent/existing annotations, field preservation, and invalid JSON. See [main_test.go](main_test.go) for the exact assertions.

## Configuration notes

1. `webhook.NewServer(webhook.Options{...})` returns the `webhook.Server` interface. Pass it through `ctrl.Options{WebhookServer: ...}` so Manager starts and stops it.
1. Generic `Validator[T]` methods receive a context and typed objects, and return warnings with an error. Errors deny the admission request; warnings can accompany either an allowed or denied response.
1. Register paths on `mgr.GetWebhookServer()`. Creating a path handler does not create the Kubernetes webhook configuration resources; deploy those separately.
