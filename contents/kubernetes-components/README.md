# Kubernetes Components

![](https://raw.githubusercontent.com/kubernetes/website/main/static/images/docs/components-of-kubernetes.svg)

*Ref: https://kubernetes.io/docs/concepts/overview/components*

- [etcd](etcd)
- [kubernetes-scheduler](kubernetes-scheduler)
- [kube-apiserver](kube-apiserver)
- [cloud-controller-manager](cloud-controller-manager)
- [kube-controller-manager](kube-controller-manager)
- [kube-proxy](kube-proxy)
- [kubelet](kubelet)
- [kubectl](kubectl)

## Build Kubernetes

```
git clone https://github.com/kubernetes/kubernetes
cd kubernetes
```

With `Golang`:
```
make
```

## Tips

### Change log level of Kubernetes component (kind)

1. Create `kind` cluster.
    ```
    kind create cluster
    ```
1. Create role and rolebinding and attach it to the `default` service account:

    ```bash
    cat << EOT | kubectl apply -f -
    apiVersion: rbac.authorization.k8s.io/v1
    kind: ClusterRole
    metadata:
      name: edit-debug-flags-v
    rules:
      # only for kubelet
    - apiGroups:
      - ""
      resources:
      - nodes/proxy
      verbs:
      - update
      # enough for other component
    - nonResourceURLs:
      - /debug/flags/v
      verbs:
      - put
    ---
    apiVersion: rbac.authorization.k8s.io/v1
    kind: ClusterRoleBinding
    metadata:
      name: edit-debug-flags-v
    roleRef:
      apiGroup: rbac.authorization.k8s.io
      kind: ClusterRole
      name: edit-debug-flags-v
    subjects:
    - kind: ServiceAccount
      name: default
      namespace: default
    EOT
    ```
1. Generate token
    ```
    TOKEN=$(kubectl create token default)
    ```
1. Change Log Level
    1. API Server
        ```
        APISERVER=$(kubectl config view -o jsonpath="{.clusters[?(@.name==\"kind-kind\")].cluster.server}")
        curl -s -X PUT -d '5' $APISERVER/debug/flags/v --header "Authorization: Bearer $TOKEN" -k
        ```
    1. kube-scheduler
        ```
        kubectl -n kube-system port-forward kube-scheduler-kind-control-plane 10259:10259
        ```

        ```
        curl -s -X PUT -d '5' https://localhost:10259/debug/flags/v --header "Authorization: Bearer $TOKEN" -k
        ```
    1. kubelet
        ```
        docker exec kind-control-plane curl -s -X PUT -d '5' https://localhost:10250/debug/flags/v --header "Authorization: Bearer $TOKEN" -k
        ```
        You might see the following warning:
        ```
        Warning: resource configmaps/kube-proxy is missing the kubectl.kubernetes.io/last-applied-configuration annotation which is required by kubectl apply. kubectl apply should only be used on resources created declaratively by either kubectl create --save-config or kubectl apply. The missing annotation will be patched automatically.
        ```

        restart kube-proxy

        ```
        kubectl -n kube-system rollout restart daemonset/kube-proxy
        ```

        port-forward:

        ```
        kubectl port-forward -n kube-system $(kubectl get pod -n kube-system -l k8s-app=kube-proxy -o jsonpath='{.items[0].metadata.name}') 10249:10249
        ```

        ```
        curl -s -XPUT -d '5' http://localhost:10249/debug/flags/v
        ```
    1. kube-controller-manager (enabled by [kubernetes/kubernetes#104571](https://github.com/kubernetes/kubernetes/pull/104571)) 10257

        ```
        kubectl -n kube-system port-forward kube-controller-manager-kind-control-plane 10257:10257
        ```

        ```
        curl -s -X PUT -d '5' https://localhost:10257/debug/flags/v --header "Authorization: Bearer $TOKEN" -k
        ```
## References
1. [99% to 99.9% SLO: High Performance Kubernetes Control Plane at Pinterest](https://medium.com/pinterest-engineering/99-to-99-9-slo-high-performance-kubernetes-control-plane-at-pinterest-894bc8a964f9)
1. [Kubernetesの主要コンポーネントのログレベルを動的に変更する](https://qiita.com/everpeace/items/a12d378c47c3ae30602f) + [](https://zaki-hmkc.hatenablog.com/entry/2022/07/27/002213)
