# kubectl

## Check APIs

1. `--v=<log level>`

    ```
    kubectl get po -n kube-system --v=6
    I0301 06:49:23.060965   14466 loader.go:372] Config loaded from file:  /Users/masato-naka/.kube/config
    I0301 06:49:23.077877   14466 round_trippers.go:454] GET https://127.0.0.1:51938/api/v1/namespaces/kube-system/pods?limit=500 200 OK in 11 milliseconds
    NAME                                         READY   STATUS    RESTARTS   AGE
    coredns-558bd4d5db-j68vq                     1/1     Running   23         28d
    coredns-558bd4d5db-wrwlx                     1/1     Running   23         28d
    etcd-kind-control-plane                      1/1     Running   23         28d
    kindnet-t57rn                                1/1     Running   23         28d
    kube-apiserver-kind-control-plane            1/1     Running   23         28d
    kube-controller-manager-kind-control-plane   1/1     Running   29         28d
    kube-proxy-8bqmv                             1/1     Running   23         28d
    kube-scheduler-kind-control-plane            1/1     Running   30         28d
    ```

    ```
    kubectl get rs -n kube-system --v=6
    I0301 06:51:05.443904   14732 loader.go:372] Config loaded from file:  /Users/masato-naka/.kube/config
    I0301 06:51:05.467278   14732 round_trippers.go:454] GET https://127.0.0.1:51938/api?timeout=32s 200 OK in 22 milliseconds
    I0301 06:51:05.514229   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis?timeout=32s 200 OK in 2 milliseconds
    I0301 06:51:05.569719   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/autoscaling/v2beta1?timeout=32s 200 OK in 5 milliseconds
    I0301 06:51:05.570835   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/kubeflow.org/v1?timeout=32s 200 OK in 7 milliseconds
    I0301 06:51:05.570854   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/batch/v1beta1?timeout=32s 200 OK in 5 milliseconds
    I0301 06:51:05.573081   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/scheduling.k8s.io/v1beta1?timeout=32s 200 OK in 8 milliseconds
    I0301 06:51:05.573087   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/authorization.k8s.io/v1beta1?timeout=32s 200 OK in 7 milliseconds
    I0301 06:51:05.573119   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/extensions/v1beta1?timeout=32s 200 OK in 7 milliseconds
    I0301 06:51:05.573150   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/discovery.k8s.io/v1?timeout=32s 200 OK in 7 milliseconds
    I0301 06:51:05.573173   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/authentication.k8s.io/v1beta1?timeout=32s 200 OK in 8 milliseconds
    I0301 06:51:05.573188   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/autoscaling/v1?timeout=32s 200 OK in 8 milliseconds
    I0301 06:51:05.573175   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/authentication.k8s.io/v1?timeout=32s 200 OK in 8 milliseconds
    I0301 06:51:05.573202   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/rbac.authorization.k8s.io/v1beta1?timeout=32s 200 OK in 8 milliseconds
    I0301 06:51:05.573212   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/node.k8s.io/v1beta1?timeout=32s 200 OK in 7 milliseconds
    I0301 06:51:05.573213   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/batch/v1?timeout=32s 200 OK in 7 milliseconds
    I0301 06:51:05.573258   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/coordination.k8s.io/v1?timeout=32s 200 OK in 7 milliseconds
    I0301 06:51:05.573276   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/authorization.k8s.io/v1?timeout=32s 200 OK in 8 milliseconds
    I0301 06:51:05.573908   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/apps/v1?timeout=32s 200 OK in 8 milliseconds
    I0301 06:51:05.574005   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/admissionregistration.k8s.io/v1?timeout=32s 200 OK in 8 milliseconds
    I0301 06:51:05.574032   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/coordination.k8s.io/v1beta1?timeout=32s 200 OK in 8 milliseconds
    I0301 06:51:05.574008   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/node.k8s.io/v1?timeout=32s 200 OK in 8 milliseconds
    I0301 06:51:05.574021   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/apiregistration.k8s.io/v1?timeout=32s 200 OK in 9 milliseconds
    I0301 06:51:05.574057   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/events.k8s.io/v1?timeout=32s 200 OK in 8 milliseconds
    I0301 06:51:05.574414   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/storage.k8s.io/v1?timeout=32s 200 OK in 7 milliseconds
    I0301 06:51:05.574439   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/networking.k8s.io/v1?timeout=32s 200 OK in 5 milliseconds
    I0301 06:51:05.574446   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/discovery.k8s.io/v1beta1?timeout=32s 200 OK in 9 milliseconds
    I0301 06:51:05.574465   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/apiextensions.k8s.io/v1beta1?timeout=32s 200 OK in 8 milliseconds
    I0301 06:51:05.574467   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/policy/v1beta1?timeout=32s 200 OK in 9 milliseconds
    I0301 06:51:05.574478   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/scheduling.k8s.io/v1?timeout=32s 200 OK in 8 milliseconds
    I0301 06:51:05.574475   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/policy/v1?timeout=32s 200 OK in 6 milliseconds
    I0301 06:51:05.574526   14732 round_trippers.go:454] GET https://127.0.0.1:51938/api/v1?timeout=32s 200 OK in 10 milliseconds
    I0301 06:51:05.574537   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/events.k8s.io/v1beta1?timeout=32s 200 OK in 8 milliseconds
    I0301 06:51:05.574666   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/admissionregistration.k8s.io/v1beta1?timeout=32s 200 OK in 4 milliseconds
    I0301 06:51:05.574702   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/rbac.authorization.k8s.io/v1?timeout=32s 200 OK in 4 milliseconds
    I0301 06:51:05.574679   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/certificates.k8s.io/v1beta1?timeout=32s 200 OK in 4 milliseconds
    I0301 06:51:05.574752   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/acid.zalan.do/v1?timeout=32s 200 OK in 4 milliseconds
    I0301 06:51:05.575375   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/apiregistration.k8s.io/v1beta1?timeout=32s 200 OK in 10 milliseconds
    I0301 06:51:05.575481   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/autoscaling/v2beta2?timeout=32s 200 OK in 9 milliseconds
    I0301 06:51:05.575600   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/apiextensions.k8s.io/v1?timeout=32s 200 OK in 6 milliseconds
    I0301 06:51:05.575743   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/storage.k8s.io/v1beta1?timeout=32s 200 OK in 9 milliseconds
    I0301 06:51:05.575745   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/networking.k8s.io/v1beta1?timeout=32s 200 OK in 10 milliseconds
    I0301 06:51:05.575778   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/flowcontrol.apiserver.k8s.io/v1beta1?timeout=32s 200 OK in 4 milliseconds
    I0301 06:51:05.575799   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/certificates.k8s.io/v1?timeout=32s 200 OK in 7 milliseconds
    I0301 06:51:07.689538   14732 round_trippers.go:454] GET https://127.0.0.1:51938/apis/apps/v1/namespaces/kube-system/replicasets?limit=500 200 OK in 2 milliseconds
    NAME                 DESIRED   CURRENT   READY   AGE
    coredns-558bd4d5db   2         2         2       28d
    ```

1. `kubectl proxy`
    ```
    kubectl proxy --port=8080
    ```

    ```
    curl localhost:8080
    ```


    <details><summary>result</summary>

    ```
    {
    "paths": [
        "/.well-known/openid-configuration",
        "/api",
        "/api/v1",
        "/apis",
        "/apis/",
        "/apis/acid.zalan.do",
        "/apis/acid.zalan.do/v1",
        "/apis/admissionregistration.k8s.io",
        "/apis/admissionregistration.k8s.io/v1",
        "/apis/admissionregistration.k8s.io/v1beta1",
        "/apis/apiextensions.k8s.io",
        "/apis/apiextensions.k8s.io/v1",
        "/apis/apiextensions.k8s.io/v1beta1",
        "/apis/apiregistration.k8s.io",
        "/apis/apiregistration.k8s.io/v1",
        "/apis/apiregistration.k8s.io/v1beta1",
        "/apis/apps",
        "/apis/apps/v1",
        "/apis/authentication.k8s.io",
        "/apis/authentication.k8s.io/v1",
        "/apis/authentication.k8s.io/v1beta1",
        "/apis/authorization.k8s.io",
        "/apis/authorization.k8s.io/v1",
        "/apis/authorization.k8s.io/v1beta1",
        "/apis/autoscaling",
        "/apis/autoscaling/v1",
        "/apis/autoscaling/v2beta1",
        "/apis/autoscaling/v2beta2",
        "/apis/batch",
        "/apis/batch/v1",
        "/apis/batch/v1beta1",
        "/apis/certificates.k8s.io",
        "/apis/certificates.k8s.io/v1",
        "/apis/certificates.k8s.io/v1beta1",
        "/apis/coordination.k8s.io",
        "/apis/coordination.k8s.io/v1",
        "/apis/coordination.k8s.io/v1beta1",
        "/apis/discovery.k8s.io",
        "/apis/discovery.k8s.io/v1",
        "/apis/discovery.k8s.io/v1beta1",
        "/apis/events.k8s.io",
        "/apis/events.k8s.io/v1",
        "/apis/events.k8s.io/v1beta1",
        "/apis/extensions",
        "/apis/extensions/v1beta1",
        "/apis/flowcontrol.apiserver.k8s.io",
        "/apis/flowcontrol.apiserver.k8s.io/v1beta1",
        "/apis/kubeflow.org",
        "/apis/kubeflow.org/v1",
        "/apis/networking.k8s.io",
        "/apis/networking.k8s.io/v1",
        "/apis/networking.k8s.io/v1beta1",
        "/apis/node.k8s.io",
        "/apis/node.k8s.io/v1",
        "/apis/node.k8s.io/v1beta1",
        "/apis/policy",
        "/apis/policy/v1",
        "/apis/policy/v1beta1",
        "/apis/rbac.authorization.k8s.io",
        "/apis/rbac.authorization.k8s.io/v1",
        "/apis/rbac.authorization.k8s.io/v1beta1",
        "/apis/scheduling.k8s.io",
        "/apis/scheduling.k8s.io/v1",
        "/apis/scheduling.k8s.io/v1beta1",
        "/apis/storage.k8s.io",
        "/apis/storage.k8s.io/v1",
        "/apis/storage.k8s.io/v1beta1",
        "/healthz",
        "/healthz/autoregister-completion",
        "/healthz/etcd",
        "/healthz/log",
        "/healthz/ping",
        "/healthz/poststarthook/aggregator-reload-proxy-client-cert",
        "/healthz/poststarthook/apiservice-openapi-controller",
        "/healthz/poststarthook/apiservice-registration-controller",
        "/healthz/poststarthook/apiservice-status-available-controller",
        "/healthz/poststarthook/bootstrap-controller",
        "/healthz/poststarthook/crd-informer-synced",
        "/healthz/poststarthook/generic-apiserver-start-informers",
        "/healthz/poststarthook/kube-apiserver-autoregistration",
        "/healthz/poststarthook/priority-and-fairness-config-consumer",
        "/healthz/poststarthook/priority-and-fairness-config-producer",
        "/healthz/poststarthook/priority-and-fairness-filter",
        "/healthz/poststarthook/rbac/bootstrap-roles",
        "/healthz/poststarthook/scheduling/bootstrap-system-priority-classes",
        "/healthz/poststarthook/start-apiextensions-controllers",
        "/healthz/poststarthook/start-apiextensions-informers",
        "/healthz/poststarthook/start-cluster-authentication-info-controller",
        "/healthz/poststarthook/start-kube-aggregator-informers",
        "/healthz/poststarthook/start-kube-apiserver-admission-initializer",
        "/livez",
        "/livez/autoregister-completion",
        "/livez/etcd",
        "/livez/log",
        "/livez/ping",
        "/livez/poststarthook/aggregator-reload-proxy-client-cert",
        "/livez/poststarthook/apiservice-openapi-controller",
        "/livez/poststarthook/apiservice-registration-controller",
        "/livez/poststarthook/apiservice-status-available-controller",
        "/livez/poststarthook/bootstrap-controller",
        "/livez/poststarthook/crd-informer-synced",
        "/livez/poststarthook/generic-apiserver-start-informers",
        "/livez/poststarthook/kube-apiserver-autoregistration",
        "/livez/poststarthook/priority-and-fairness-config-consumer",
        "/livez/poststarthook/priority-and-fairness-config-producer",
        "/livez/poststarthook/priority-and-fairness-filter",
        "/livez/poststarthook/rbac/bootstrap-roles",
        "/livez/poststarthook/scheduling/bootstrap-system-priority-classes",
        "/livez/poststarthook/start-apiextensions-controllers",
        "/livez/poststarthook/start-apiextensions-informers",
        "/livez/poststarthook/start-cluster-authentication-info-controller",
        "/livez/poststarthook/start-kube-aggregator-informers",
        "/livez/poststarthook/start-kube-apiserver-admission-initializer",
        "/logs",
        "/metrics",
        "/openapi/v2",
        "/openid/v1/jwks",
        "/readyz",
        "/readyz/autoregister-completion",
        "/readyz/etcd",
        "/readyz/informer-sync",
        "/readyz/log",
        "/readyz/ping",
        "/readyz/poststarthook/aggregator-reload-proxy-client-cert",
        "/readyz/poststarthook/apiservice-openapi-controller",
        "/readyz/poststarthook/apiservice-registration-controller",
        "/readyz/poststarthook/apiservice-status-available-controller",
        "/readyz/poststarthook/bootstrap-controller",
        "/readyz/poststarthook/crd-informer-synced",
        "/readyz/poststarthook/generic-apiserver-start-informers",
        "/readyz/poststarthook/kube-apiserver-autoregistration",
        "/readyz/poststarthook/priority-and-fairness-config-consumer",
        "/readyz/poststarthook/priority-and-fairness-config-producer",
        "/readyz/poststarthook/priority-and-fairness-filter",
        "/readyz/poststarthook/rbac/bootstrap-roles",
        "/readyz/poststarthook/scheduling/bootstrap-system-priority-classes",
        "/readyz/poststarthook/start-apiextensions-controllers",
        "/readyz/poststarthook/start-apiextensions-informers",
        "/readyz/poststarthook/start-cluster-authentication-info-controller",
        "/readyz/poststarthook/start-kube-aggregator-informers",
        "/readyz/poststarthook/start-kube-apiserver-admission-initializer",
        "/readyz/shutdown",
        "/version"
    ]
    }%
    ```

    </details>

    ```
    curl localhost:8080/api/v1  | jq -r '.resources[].name'
    ```

    <details><summary>result</summary>

    ```
    bindings
    componentstatuses
    configmaps
    endpoints
    events
    limitranges
    namespaces
    namespaces/finalize
    namespaces/status
    nodes
    nodes/proxy
    nodes/status
    persistentvolumeclaims
    persistentvolumeclaims/status
    persistentvolumes
    persistentvolumes/status
    pods
    pods/attach
    pods/binding
    pods/eviction
    pods/exec
    pods/log
    pods/portforward
    pods/proxy
    pods/status
    podtemplates
    replicationcontrollers
    replicationcontrollers/scale
    replicationcontrollers/status
    resourcequotas
    resourcequotas/status
    secrets
    serviceaccounts
    serviceaccounts/token
    services
    services/proxy
    services/status
    ```

    </details>

    ```
    curl localhost:8080/apis/apps/v1  | jq -r '.resources[].name'
    ```

    <details><summary>result</summary>

    ```
    controllerrevisions
    daemonsets
    daemonsets/status
    deployments
    deployments/scale
    deployments/status
    replicasets
    replicasets/scale
    replicasets/status
    statefulsets
    statefulsets/scale
    statefulsets/status
    ```

    </details>
