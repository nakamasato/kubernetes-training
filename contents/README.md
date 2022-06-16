# kubernetes-training

# Versions

1. Kubernetes: [v1.21.2](https://github.com/kubernetes/kubernetes/releases/tag/v1.21.2) (released on 2021-06-18)
1. kustomize: [v4.2.0](https://github.com/kubernetes-sigs/kustomize/releases/tag/kustomize%2Fv4.2.0) (released on 2021-07-02)
1. Helm: [3.6.3](https://github.com/helm/helm/releases/tag/v3.6.3) (released on 2021-07-15)
1. Traefik: [v2.8.0](https://github.com/traefik/traefik/releases/tag/v2.8.0)
1. ArgoCD: [v2.2.3](https://github.com/argoproj/argo-cd/releases/tag/v2.2.3) (released on 2022-01-19)
1. Prometheus-Operator: [v0.53.1](https://github.com/prometheus-operator/prometheus-operator/releases/tag/v0.53.1) (released on 2021-12-20)
1. Prometheus: Latest
1. Grafana: Latest
1. Strimzi: [0.24.0](https://github.com/strimzi/strimzi-kafka-operator/releases/tag/0.24.0) (released on 2021-06-24)
1. Kind: [v0.11.1](https://github.com/kubernetes-sigs/kind/releases/tag/v0.11.1) (released on 2021-05-28)
1. Ingress Nginx Controller: [v0.48.0](https://github.com/kubernetes/ingress-nginx/releases/tag/controller-v0.48.1) (released on 2021-07-15)
1. Conftest: [0.25.0](https://github.com/open-policy-agent/conftest/releases/tag/v0.25.0) (released on 2021-05-08)
1. Istio: [1.12.2](https://github.com/istio/istio/releases/tag/1.12.2) (released on 2022-01-19)
1. PostgresOperator: [v1.7.1](https://github.com/zalando/postgres-operator/releases/tag/v1.7.1) (released on 2021-11-04)
1. Cert Manager: [v1.7.1](https://github.com/cert-manager/cert-manager/releases/tag/v1.7.1) (released on 2022-02-05)

# Contents

1. Kubernetes Cluster
    1. [kubernetes-the-hard-way](kubernetes-the-hard-way)
    1. [Kubeadm in local](kubeadm-local)
    1. [kind](local-cluster/kind)
1. [Kubernetes Features](kubernetes-features)
    1. [Autoscaler HPA with custom metrics](autoscaler/hpa/custom-metrics)
    1. [amazon-eks-workshop](eksworkshop)
1. Kubernetes Components
    1. [kubernetes-scheduler](kubernetes-components/kubernetes-scheduler)
    1. [etcd](kubernetes-components/etcd)
    1. [kube-apiserver](kubernetes-components/kube-apiserver)
    1. [kube-controller-manager](kubernetes-components/kube-controller-manager)
    1. kube-proxy
    1. kubelet
1. Kubernetes Extensions
    1. [kubernetes-operator](kubernetes-operator)
    1. [kubernetes-scheduler](kubernetes-extensions/kubernetes-scheduler)
    1. [plugins (todo)](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/)
1. Deloyment Managemet
    1. [Knative](knative)
    1. [Skaffold (todo)](https://skaffold.dev/)
1. Networking
    1. [traefik](traefik)
    1. [ingress-nginx-controller](ingress-nginx-controller)
1. Middleware (Operator)
    1. [strimzi](strimzi)
    1. [eck](eck)
1. Service Proxy, Discovery, and, Mesh
    1. [istio](istio)
1. Monitoring
    1. [Prometheus](prometheus)
    1. [Prometheus Operator](prometheus-operator)
    1. [Thanos (todo)] https://thanos.io/
    1. [Grafana](grafana)
    1. [Grafana Operator](grafana-operator)
    1. [Grafana Loki](loki)
    1. [Grafana Tempo (todo)] https://grafana.com/docs/tempo/latest/
    1. [Jaeger (todo)] https://www.jaegertracing.io/
1. Security
    1. [open-policy-agent](open-policy-agent)
    1. [Cert Manager](cert-manager)
1. Yaml Management
    1. [Helm](helm)
    1. [Helm vs Kustomize](helm-vs-kustomize)
1. CI/CD
    1. [Conftest](open-policy-agent/conftest)
    1. [ArgoCD](argocd)
    1. Kyverno https://kyverno.io/
    1. Polaris https://www.fairwinds.com/polaris
1. Machine Learning
    1. [kubeflow](https://github.com/nakamasato/kubeflow-training)
1. [Databases](databases)
    1. [Vitess] https://github.com/vitessio/vitess
    1. [TiDB] https://github.com/pingcap/tidb
    1. [TimescaleDB] https://github.com/timescale/timescaledb-kubernetes
    1. [mysql-operator](databases/mysql-operator)
    1. [postgres-operator](databases/postgres-operator)
# Cloud Native Trail Map

- https://github.com/cncf/trailmap
- https://www.cncf.io/blog/2018/03/08/introducing-the-cloud-native-landscape-2-0-interactive-edition/

![alt text](https://github.com/cncf/trailmap/blob/master/CNCF_TrailMap_latest.png?raw=true)

## 1. CONTAINERIZATION

1. [Containers 101: attach vs. exec - what's the difference?](https://iximiuz.com/en/posts/containers-101-attach-vs-exec/)
## 2. CI/CD

### 2.1 [ArgoCD](argocd)
## 3. ORCHESTRATION & APPLICATION DEFINITION

### 3.1 Kubernetes

#### Useful Commands

- DNS
    ```
    kubectl apply -f https://k8s.io/examples/admin/dns/dnsutils.yaml
    kubectl exec -i -t dnsutils -- nslookup kubernetes.default
    ```
- [Debug with ephemeral containers](https://kubernetes.io/docs/tasks/debug-application-cluster/debug-running-pod/#ephemeral-container-example) (alpha in 1.22, beta in 1.23)
    ```
    kubectl run ephemeral-demo --image=k8s.gcr.io/pause:3.1 --restart=Never
    kubectl debug -it ephemeral-demo --image=busybox --target=ephemeral-demo
    ```
- Create pod with busyboxy-curl
    ```
    kubectl run -it --rm=true busybox --image=yauritux/busybox-curl --restart=Never
    ```

#### Set up Kubernetes Cluster with kubeadm (local)

[kubeadm-local](kubeadm-local)
#### Set up Kubernetes Cluster on GCP (kubernetes-the-hard-way)

[Kubernetes The Hard Way](https://github.com/kelseyhightower/kubernetes-the-hard-way)

#### Kubernetes Components

1. [kubernetes-scheduler](kubernetes-components/kubernetes-scheduler)
1. [etcd](kubernetes-components/etcd)
1. kube-apiserver
1. kube-controller-manager
1. kube-proxy
1. kubelet

#### [More Practices of Applications on Kubernetes](PRACTICE.md)

### 3.2 [Helm](helm)

1. Create Helm chart.

    ```
    helm create <chart-name e.g. helm-example>
    ```

1. Update files under `templates` and `values.yaml`
1. Test apply.

    ```
    helm install helm-example --debug ./helm-example
    ```

1. Make a package.

    ```
    helm package helm-example
    ```

1. Create repository and set index.

    ```
    helm repo index ./ --url https://nakamasato.github.io/helm-charts-repo
    ```

1. Install a chart.

    ```
    helm repo add nakamasato https://nakamasato.github.io/helm-charts-repo
    helm repo update # update the repository info
    helm install example-from-my-repo nakamasato/helm-example
    ```

## 4. OBSERVABILITY & ANALYTICS

### 4.1. [Prometheus](prometheus)

![](prometheus/prometheus.drawio.svg)
### 4.2. [Prometheus Operator](prometheus-operator)

### TBD
- fluentd
- Jaeger
- Open Tracing

## 5. SERVICE PROXY, DISCOVERY & MESH

### 5.1. [Istio](istio)

### TBD
- envoy
- CoreDNS
- Linkerd

## 6. NETWORKING, POLICY & SECURITY

### 6.1 [Open Policy Agent](open-policy-agent)

### [gatekeeper](https://github.com/open-policy-agent/gatekeeper)

1. Install gatekeeper

    ```
    kubectl apply -f https://raw.githubusercontent.com/open-policy-agent/gatekeeper/master/deploy/gatekeeper.yaml
    ```

1. Create `ConstraintTemplate`
1. Create custom policy defined in the previous step.

### [conftest](https://github.com/open-policy-agent/conftest)

1. Write policy in `policy` directory.

    ```rego
    deny[msg] {
      input.kind = "Deployment"
      not input.spec.template.spec.nodeSelector
      msg = "Deployment must have nodeSelector"
    }
    ```

1. Write tests in the same directory.

    ```rego
    test_no_nodeSelector {
      deny["Deployment must have nodeSelector"] with input as
      {
        "kind": "Deployment",
        "spec": {
          "template": {
            "spec": {
              "containers": [
              ],
            }
          }
        }
      }
    }
    ```

1. Run test.

    ```
    conftest verify

    1 tests, 1 passed, 0 warnings, 0 failures, 0 exceptions
    ```

1. Validate a manifest file.

    ```
    conftest test manifests/valid/deployment.yaml

    1 tests, 1 passed, 0 warnings, 0 failures, 0 exceptions
    ```

### TBD
- CNI
- falco

## 7. DISTRIBUTED DATABASE & STORAGE

### 7.1. [etcd](kubernetes-components/etcd)
### TBD
- [Vitess](https://github.com/vitessio/vitess)
- Rook
- [TiDB](https://github.com/pingcap/tidb)
- [TimescaleDB](https://github.com/timescale/timescaledb-kubernetes)

## 8. STREAMING & MESSAGING

### TBD
- gRPC
- NATS
- cloudevents

## 9. CONTAINER REGISTRY & RUNTIME

### TBD
- containerd
- harbor
- cri-o

## 10. SOFTWARE DISTRIBUTION

### TBD
- TUF
- notaru
