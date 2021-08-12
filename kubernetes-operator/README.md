# Kubernetes Operator

Study Steps:

1. Use existing operators:
    - [prometheus-operator](../prometheus-operator)
    - [postgres-operator](../postgres-operator)
    - [strimzi](../strimzi)
    - rabbitmq-operator
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
    - [kubebuilder](https://book.kubebuilder.io/)
        - [Tutorial: Building CronJob](https://book.kubebuilder.io/cronjob-tutorial/cronjob-tutorial.html)
