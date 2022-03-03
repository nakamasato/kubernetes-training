# Kubernetes Extension

## Kubernetes Scheduler

[kubernetes-scheduler](kubernetes-scheduler)

## Custom Resource Definition

- [NewCustomResourceDefinitionHandler](https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apiextensions-apiserver/pkg/apiserver/customresource_handler.go) is called in [CompletedConfig.New](https://github.com/kubernetes/kubernetes/blob/16c9d59d2d646a77fa5de0532fa7c583c013b8d6/staging/src/k8s.io/apiextensions-apiserver/pkg/apiserver/apiserver.go#L133)
- [CompletedConfig.New](https://github.com/kubernetes/kubernetes/blob/16c9d59d2d646a77fa5de0532fa7c583c013b8d6/staging/src/k8s.io/apiextensions-apiserver/pkg/apiserver/apiserver.go#L133)
    1. Prepare genericServer with [completedConfig.New](https://github.com/kubernetes/kubernetes/blob/16c9d59d2d646a77fa5de0532fa7c583c013b8d6/staging/src/k8s.io/apiserver/pkg/server/config.go#L567).
    1. Initialize `CustomResourceDefinitions` with `GenericAPIServer`.
    1. Initialize `apiGroupInfo` with [genericapiserver.NewDefaultAPIGroupInfo](https://github.com/kubernetes/kubernetes/blob/16c9d59d2d646a77fa5de0532fa7c583c013b8d6/staging/src/k8s.io/apiserver/pkg/server/genericapiserver.go#L697).
    1. Install API group with `s.GenericAPIServer.InstallAPIGroup`.
    1. Initialize clientset for CRD with `crdClient, err := clientset.NewForConfig(s.GenericAPIServer.LoopbackClientConfig)`
    1. Initialize and set informer with `s.Informers = externalinformers.NewSharedInformerFactory(crdClient, 5*time.Minute)`
    1. Prepare handlers
        1. delegateHandler
        1. versionDiscoveryHandler
        1. groupDiscoveryHandler
    1. Initialize `EstablishingController`.
    1. Initialize `crdHandler` by `NewCustomResourceDefinitionHandler` with `versionDiscoveryHandler`, `groupDiscoveryHandler`, informer, `delegateHandler`, `establishingController`, etc.
    1. Set HTTP handler for GenericAPIServer with `crdHandler`.
        ```go
        s.GenericAPIServer.Handler.NonGoRestfulMux.Handle("/apis", crdHandler)
        s.GenericAPIServer.Handler.NonGoRestfulMux.HandlePrefix("/apis/", crdHandler)
        ```
    1. Initialize controllers.
        - discoveryController
        - namingController
        - nonStructuralSchemaController
        - apiApprovalController
        - finalizingController
        - [openapicontroller](https://github.com/kubernetes/kubernetes/blob/ea0764452222146c47ec826977f49d7001b0ea8c/staging/src/k8s.io/apiextensions-apiserver/pkg/controller/openapi/controller.go#L62)
    1. Set `AddPostStartHookOrDie` for `GenericAPIServer` to start informer.
    1. Set `AddPostStartHookOrDie` for `GenericAPIServer` to start controllers.
    1. Set `AddPostStartHookOrDie` for `GenericAPIServer` to wait until CRD informer is synced.
