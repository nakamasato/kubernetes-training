# Kubernetes Operator

Study Steps:

1. Use existing operators:
    - [prometheus-operator](../prometheus-operator)
    - [postgres-operator](../postgres-operator)
    - [strimzi](../strimzi)
    - rabbitmq-operator
    - [argocd](../argocd): [appcontroller.go](https://github.com/argoproj/argo-cd/blob/9025318adf367ae8f13b1a99e5c19344402b7bb9/controller/appcontroller.go)
1. Understand what is Kubernetes operator.
    1. Kubernetes Controller components.
    1. How Kubernetes Controlloer works.
    1. Custom Resource.
1. Create your own operator.
    - [sample-controller](https://github.com/kubernetes/sample-controller): https://github.com/nakamasato/foo-controller-kubebuilder
    - [operator-sdk](https://sdk.operatorframework.io/)
        - [go-based](https://sdk.operatorframework.io/docs/building-operators/golang/quickstart/): https://github.com/nakamasato/memcached-operator
        - [helm-based](https://sdk.operatorframework.io/docs/building-operators/helm/quickstart/): https://github.com/nakamasato/nginx-operator
        - [ansible-based](https://sdk.operatorframework.io/docs/building-operators/ansible/quickstart/)
        - [mysql-operator]
    - [kubebuilder](https://book.kubebuilder.io/)
        - [Tutorial: Building CronJob](https://book.kubebuilder.io/cronjob-tutorial/cronjob-tutorial.html)
1. Important topics:
    - [How can I have separate logic for Create, Update, and Delete events? When reconciling an object can I access its previous state?](https://sdk.operatorframework.io/docs/faqs/#how-can-i-have-separate-logic-for-create-update-and-delete-events-when-reconciling-an-object-can-i-access-its-previous-state) -> You should not have separate logic. Instead design your reconciler to be idempotent.
        - [Q: How do I have different logic in my reconciler for different types of events (e.g. create, update, delete)? in controller-runtime](https://github.com/kubernetes-sigs/controller-runtime/blob/master/FAQ.md#q-how-do-i-have-different-logic-in-my-reconciler-for-different-types-of-events-eg-create-update-delete)
    - [Owners and Dependents](https://kubernetes.io/docs/concepts/overview/working-with-objects/owners-dependents/)
    - Finalizer
        -
    - Reconcile Loop
        - Based on the return value of Reconcile() the reconcile Request may be requeued and the loop may be triggered again: ([Building a Go-based Memcached Operator using the Operator SDK](https://docs.openshift.com/container-platform/4.1/applications/operator_sdk/osdk-getting-started.html#building-memcached-operator-using-osdk_osdk-getting-started))
            ```go
            // Reconcile successful - don't requeue
            return reconcile.Result{}, nil
            // Reconcile failed due to error - requeue
            return reconcile.Result{}, err
            // Requeue for any reason other than error
            return reconcile.Result{Requeue: true}, nil
            ```
        - https://github.com/operator-framework/operator-sdk/issues/4209#issuecomment-729916367
    - Testing:
        - KUbernetes Testing TooL (kuttl) https://kuttl.dev/ KUTTL is built to support some kubernetes integration test scenarios and is most valuable as an end-to-end (e2e) test harness.
1. Useful tools:
    - https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/controller/controllerutil
