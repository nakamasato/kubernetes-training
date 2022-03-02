# Kubernetes Extension

## Custom Resource Definition

- [NewCustomResourceDefinitionHandler](https://github.com/kubernetes/kubernetes/blob/ea0764452222146c47ec826977f49d7001b0ea8c/staging/src/k8s.io/apiextensions-apiserver/pkg/apiserver/customresource_handler.go#L174)
    ```go
	s.GenericAPIServer.Handler.NonGoRestfulMux.Handle("/apis", crdHandler)
	s.GenericAPIServer.Handler.NonGoRestfulMux.HandlePrefix("/apis/", crdHandler)
    ```
- Controllers
    - discoveryController
    - namingController
    - nonStructuralSchemaController
    - apiApprovalController
    - finalizingController
    - [openapicontroller](https://github.com/kubernetes/kubernetes/blob/ea0764452222146c47ec826977f49d7001b0ea8c/staging/src/k8s.io/apiextensions-apiserver/pkg/controller/openapi/controller.go#L62)
