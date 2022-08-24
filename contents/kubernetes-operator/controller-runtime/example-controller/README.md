# Example Controller

First example in [controller-runtime](https://pkg.go.dev/sigs.k8s.io/controller-runtime)
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
    kubectl create deploy test --replicas=3 --image=nginx
    ```

    You would see some errors, these errors happen because replicaset is updated by replicaset controller according to new Pod's status.

    ```
    1.6565476218562e+09     ERROR   controller.replicaset   Reconciler error        {"reconciler group": "apps", "reconciler kind": "ReplicaSet", "name": "test-8499f4f74", "namespace": "default", "error": "Operation cannot be fulfilled on replicasets.apps \"test-8499f4f74\": the object has been modified; please apply your changes to the latest version and try again"}
    sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).processNextWorkItem
            /Users/nakamasato/.gvm/pkgsets/go1.17.9/global/pkg/mod/sigs.k8s.io/controller-runtime@v0.11.2/pkg/internal/controller/controller.go:266
    sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).Start.func2.2
            /Users/nakamasato/.gvm/pkgsets/go1.17.9/global/pkg/mod/sigs.k8s.io/controller-runtime@v0.11.2/pkg/internal/controller/controller.go:227
    ```

1. Check `pod-count=3` labels added to the ReplicaSet

    ```
    kubectl get rs -o jsonpath='{.items[].metadata.labels}'
    {"app":"test","pod-count":"3","pod-template-hash":"8499f4f74"}
    ```

1. Clean up
    ```
    kubectl delete deploy test
    ```
