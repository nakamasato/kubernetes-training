# Kubernetes Tools

## https://pkg.go.dev/k8s.io/code-generator

Golang code-generators used to implement Kubernetes-style API types.

## https://pkg.go.dev/k8s.io/client-go

Go clients for talking to a kubernetes cluster.
- https://pkg.go.dev/k8s.io/client-go/tools
    - [cache](https://pkg.go.dev/k8s.io/client-go@v0.23.4/tools/cache): Package cache is a client-side caching mechanism.
## https://pkg.go.dev/k8s.io/apimachinery

Scheme, typing, encoding, decoding, and conversion packages for Kubernetes and Kubernetes-like API objects.
- [runtime](https://pkg.go.dev/k8s.io/apimachinery/pkg/runtime): Package runtime defines conversions between generic types and structs to map query strings to struct objects.
    - [Scheme](https://pkg.go.dev/k8s.io/apimachinery/pkg/runtime#Scheme): Scheme defines methods for serializing and deserializing API objects, a type registry for converting group, version, and kind information to and from Go schemas, and mappings between Go schemas of different versions. A scheme is the foundation for a versioned API and versioned configuration over time.
## https://pkg.go.dev/sigs.k8s.io/controller-runtime

The Kubernetes controller-runtime Project is a set of go libraries for building Controllers. It is leveraged by Kubebuilder and Operator SDK.
- Client
- Cache
- Manager
- Controller
- Webhook
- Reconciler
- Source
- EventHandler
- Predicate
