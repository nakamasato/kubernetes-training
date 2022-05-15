# [controller-runtime](https://pkg.go.dev/sigs.k8s.io/controller-runtime)

[controller-runtime](https://pkg.go.dev/sigs.k8s.io/controller-runtime) is a subproject of kubebuilder which provides a lot of useful tools that help develop Kubernetes Operator.

1. Manager
1. Controller
    1. Event
    1. Builder
    1. Source
    1. Handler
    1. Predicate
1. Client
1. Cache
1. Scheme
1. Webhook
1. Envtest

## Example

1. Run example controller.

    ```
    go run example-controller.go
    ```

    What this controller does:

    1. Read the ReplicaSet
    1. Read the Pods
    1. Set a Label on the ReplicaSet with the Pod count.

1. Create Deployment
    ```
    kubectl create deploy test --image=nginx
    ```
1. Check `pod-count` labels adde to the ReplicaSet
    ```
    kubectl get rs -o jsonpath='{.items[].metadata.labels}'
    {"app":"test","pod-count":"1","pod-template-hash":"8499f4f74"}
    ```
1. Clean up
    ```
    kubectl delete deploy test
    ```
