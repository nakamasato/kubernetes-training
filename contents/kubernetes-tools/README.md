# Kubernetes Tools

## 1. https://pkg.go.dev/k8s.io/code-generator

Golang code-generators used to implement Kubernetes-style API types.

## 2. https://pkg.go.dev/k8s.io/client-go

Go clients for talking to a kubernetes cluster.
- https://pkg.go.dev/k8s.io/client-go/tools
    - [cache](https://pkg.go.dev/k8s.io/client-go@v0.23.4/tools/cache): Package cache is a client-side caching mechanism.
## 3. https://pkg.go.dev/k8s.io/apimachinery

Scheme, typing, encoding, decoding, and conversion packages for Kubernetes and Kubernetes-like API objects.
- [runtime](https://pkg.go.dev/k8s.io/apimachinery/pkg/runtime): Package runtime defines conversions between generic types and structs to map query strings to struct objects.
    - [Scheme](https://pkg.go.dev/k8s.io/apimachinery/pkg/runtime#Scheme): Scheme defines methods for serializing and deserializing API objects, a type registry for converting group, version, and kind information to and from Go schemas, and mappings between Go schemas of different versions. A scheme is the foundation for a versioned API and versioned configuration over time.
## 4. https://pkg.go.dev/sigs.k8s.io/controller-runtime

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

https://github.com/kubernetes-sigs/kubebuilder/blob/master/DESIGN.md#controller-runtime

## 5. https://pkg.go.dev/sigs.k8s.io/controller-tools/cmd/controller-gen

Generate Kubernetes controller stubs that sync configurable resource types

<details><summary>controller-gen -h</summary>

```
controller-gen --help
Generate Kubernetes API extension resources and code.

Usage:
  controller-gen [flags]

Examples:
        # Generate RBAC manifests and crds for all types under apis/,
        # outputting crds to /tmp/crds and everything else to stdout
        controller-gen rbac:roleName=<role name> crd paths=./apis/... output:crd:dir=/tmp/crds output:stdout

        # Generate deepcopy/runtime.Object implementations for a particular file
        controller-gen object paths=./apis/v1beta1/some_types.go

        # Generate OpenAPI v3 schemas for API packages and merge them into existing CRD manifests
        controller-gen schemapatch:manifests=./manifests output:dir=./manifests paths=./pkg/apis/...

        # Run all the generators for a given project
        controller-gen paths=./apis/...

        # Explain the markers for generating CRDs, and their arguments
        controller-gen crd -ww


Flags:
  -h, --detailed-help count   print out more detailed help
                              (up to -hhh for the most detailed output, or -hhhh for json output)
      --help                  print out usage and a summary of options
      --version               show version
  -w, --which-markers count   print out all markers available with the requested generators
                              (up to -www for the most detailed output, or -wwww for json output)


Options


generators

+webhook                                                                                                  package  generates (partial) {Mutating,Validating}WebhookConfiguration objects.
+schemapatch:manifests=<string>[,maxDescLen=<int>]                                                        package  patches existing CRDs with new schemata.
+rbac:roleName=<string>                                                                                   package  generates ClusterRole objects.
+object[:headerFile=<string>][,year=<string>]                                                             package  generates code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
+crd[:crdVersions=<[]string>][,maxDescLen=<int>][,preserveUnknownFields=<bool>][,trivialVersions=<bool>]  package  generates CustomResourceDefinition objects.


generic

+paths=<[]string>  package  represents paths and go-style path patterns to use as package roots.


output rules (optionally as output:<generator>:...)

+output:artifacts[:code=<string>],config=<string>  package  outputs artifacts to different locations, depending on whether they're package-associated or not.
+output:dir=<string>                               package  outputs each artifact to the given directory, regardless of if it's package-associated or not.
+output:none                                       package  skips outputting anything.
+output:stdout                                     package  outputs everything to standard-out, with no separation.
```

</details>

<details><summary>controller-gen crd -w</summary>

```
controller-gen crd -w

CRD

+groupName=<string>                                                                                                               package  specifies the API group name for this package.
+kubebuilder:printcolumn:JSONPath=<string>[,description=<string>][,format=<string>],name=<string>[,priority=<int>],type=<string>  type     adds a column to "kubectl get" output for this CRD.
+kubebuilder:resource[:categories=<[]string>][,path=<string>][,scope=<string>][,shortName=<[]string>][,singular=<string>]         type     configures naming and scope for a CRD.
+kubebuilder:skip                                                                                                                 package  don't consider this package as an API version.
+kubebuilder:skipversion                                                                                                          type     removes the particular version of the CRD from the CRDs spec.
+kubebuilder:storageversion                                                                                                       type     marks this version as the "storage version" for the CRD for conversion.
+kubebuilder:subresource:scale[:selectorpath=<string>],specpath=<string>,statuspath=<string>                                      type     enables the "/scale" subresource on a CRD.
+kubebuilder:subresource:status                                                                                                   type     enables the "/status" subresource on a CRD.
+versionName=<string>                                                                                                             package  overrides the API group version for this package (defaults to the package name).


CRD processing

+kubebuilder:pruning:PreserveUnknownFields      field  PreserveUnknownFields stops the apiserver from pruning fields which are not specified.
+kubebuilder:validation:XPreserveUnknownFields  type   PreserveUnknownFields stops the apiserver from pruning fields which are not specified.
+kubebuilder:validation:XPreserveUnknownFields  field  PreserveUnknownFields stops the apiserver from pruning fields which are not specified.


CRD validation

+kubebuilder:default=<any>                       field    sets the default value for this field.
+kubebuilder:validation:EmbeddedResource         field    EmbeddedResource marks a fields as an embedded resource with apiVersion, kind and metadata fields.
+kubebuilder:validation:Enum=<[]any>             type     specifies that this (scalar) field is restricted to the *exact* values specified here.
+kubebuilder:validation:Enum=<[]any>             field    specifies that this (scalar) field is restricted to the *exact* values specified here.
+kubebuilder:validation:ExclusiveMaximum=<bool>  field    indicates that the maximum is "up to" but not including that value.
+kubebuilder:validation:ExclusiveMaximum=<bool>  type     indicates that the maximum is "up to" but not including that value.
+kubebuilder:validation:ExclusiveMinimum=<bool>  type     indicates that the minimum is "up to" but not including that value.
+kubebuilder:validation:ExclusiveMinimum=<bool>  field    indicates that the minimum is "up to" but not including that value.
+kubebuilder:validation:Format=<string>          type     specifies additional "complex" formatting for this field.
+kubebuilder:validation:Format=<string>          field    specifies additional "complex" formatting for this field.
+kubebuilder:validation:MaxItems=<int>           field    specifies the maximum length for this list.
+kubebuilder:validation:MaxItems=<int>           type     specifies the maximum length for this list.
+kubebuilder:validation:MaxLength=<int>          field    specifies the maximum length for this string.
+kubebuilder:validation:MaxLength=<int>          type     specifies the maximum length for this string.
+kubebuilder:validation:Maximum=<int>            field    specifies the maximum numeric value that this field can have.
+kubebuilder:validation:Maximum=<int>            type     specifies the maximum numeric value that this field can have.
+kubebuilder:validation:MinItems=<int>           field    specifies the minimun length for this list.
+kubebuilder:validation:MinItems=<int>           type     specifies the minimun length for this list.
+kubebuilder:validation:MinLength=<int>          field    specifies the minimum length for this string.
+kubebuilder:validation:MinLength=<int>          type     specifies the minimum length for this string.
+kubebuilder:validation:Minimum=<int>            type     specifies the minimum numeric value that this field can have.
+kubebuilder:validation:Minimum=<int>            field    specifies the minimum numeric value that this field can have.
+kubebuilder:validation:MultipleOf=<int>         field    specifies that this field must have a numeric value that's a multiple of this one.
+kubebuilder:validation:MultipleOf=<int>         type     specifies that this field must have a numeric value that's a multiple of this one.
+kubebuilder:validation:Optional                 field    specifies that this field is optional, if fields are required by default.
+kubebuilder:validation:Optional                 package  specifies that all fields in this package are optional by default.
+kubebuilder:validation:Pattern=<string>         type     specifies that this string must match the given regular expression.
+kubebuilder:validation:Pattern=<string>         field    specifies that this string must match the given regular expression.
+kubebuilder:validation:Required                 field    specifies that this field is required, if fields are optional by default.
+kubebuilder:validation:Required                 package  specifies that all fields in this package are required by default.
+kubebuilder:validation:Type=<string>            field    overrides the type for this field (which defaults to the equivalent of the Go type).
+kubebuilder:validation:Type=<string>            type     overrides the type for this field (which defaults to the equivalent of the Go type).
+kubebuilder:validation:UniqueItems=<bool>       field    specifies that all items in this list must be unique.
+kubebuilder:validation:UniqueItems=<bool>       type     specifies that all items in this list must be unique.
+kubebuilder:validation:XEmbeddedResource        type     EmbeddedResource marks a fields as an embedded resource with apiVersion, kind and metadata fields.
+kubebuilder:validation:XEmbeddedResource        field    EmbeddedResource marks a fields as an embedded resource with apiVersion, kind and metadata fields.
+nullable                                        field    marks this field as allowing the "null" value.
+optional                                        field    specifies that this field is optional, if fields are required by default.
```

</details>
