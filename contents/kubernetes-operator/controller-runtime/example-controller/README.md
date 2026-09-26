# Example Controller

A ReplicaSet controller built with controller-runtime v0.25.1. Run commands from the repository root with a kubeconfig and permissions to read ReplicaSets/Pods and patch ReplicaSets.
## Example

1. Run example controller.

    ```
    go run ./contents/kubernetes-operator/controller-runtime/example-controller
    ```

    What this controller does:

    1. Read the ReplicaSet
    1. Read the Pods
    1. Set a Label on the ReplicaSet with the Pod count.

1. Create Deployment

    ```
    kubectl create deploy test --replicas=3 --image=nginx
    ```

    The controller counts Pods that both match the ReplicaSet selector and have its controller owner UID. It supports matchExpressions, initializes missing labels, and patches only when pod-count changes. A merge patch avoids overwriting unrelated fields; other API errors still trigger retries.

1. Check `pod-count=3` labels added to the ReplicaSet

    ```
    kubectl get rs -o jsonpath='{.items[].metadata.labels}'
    {"app":"test","pod-count":"3","pod-template-hash":"8499f4f74"}
    ```

1. Clean up
    ```
    kubectl delete deploy test
    ```

## Implementation and event flow

`For(&appsv1.ReplicaSet{})` maps ReplicaSet events to its own key. `Owns(&corev1.Pod{})` maps Pod ownerReferences to the controlling ReplicaSet. Manager supplies a cache-backed Client explicitly to ReplicaSetReconciler.

On each request:

1. Get the ReplicaSet; ignore NotFound because deletion leaves nothing to label.
2. Convert Spec.Selector with LabelSelectorAsSelector, including matchExpressions.
3. List matching Pods in the request namespace.
4. Count only Pods whose controller owner UID matches this ReplicaSet. Labels alone do not prove ownership.
5. Return immediately if pod-count already matches.
6. DeepCopy the original, allocate Labels if nil, and patch with client.MergeFrom(original).

The write produces another watch event; the no-op check makes that subsequent reconciliation settle. Reads can lag writes because the Client reads from the cache. This sample counts existing owned Pods, not only Ready Pods, and does not create or delete Pods itself.

## Tests

```sh
go test ./contents/kubernetes-operator/controller-runtime/example-controller
```

The fake-client regression tests cover missing ReplicaSets, nil labels, selector expressions, owner filtering, namespace filtering, and avoiding unchanged writes. They do not substitute for a running-cluster check of informer delivery or API-server admission.
