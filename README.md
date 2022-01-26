# kubernetes-training

# Versions

- Kubernetes: [v1.21.2](https://github.com/kubernetes/kubernetes/releases/tag/v1.21.2) (released on 2021-06-18)
- kustomize: [v4.2.0](https://github.com/kubernetes-sigs/kustomize/releases/tag/kustomize%2Fv4.2.0) (released on 2021-07-02)
- Helm: [3.6.3](https://github.com/helm/helm/releases/tag/v3.6.3) (released on 2021-07-15)
- Traefik: v2.2
- ArgoCD: [v2.1.1](https://github.com/argoproj/argo-cd/releases/tag/v2.1.1) (released on 2021-08-26)
- Prometheus-Operator: [v0.53.1](https://github.com/prometheus-operator/prometheus-operator/releases/tag/v0.53.1) (released on 2021-12-20)
- Prometheus: Latest
- Grafana: Latest
- Strimzi: [0.24.0](https://github.com/strimzi/strimzi-kafka-operator/releases/tag/0.24.0) (released on 2021-06-24)
- Kind: [v0.11.1](https://github.com/kubernetes-sigs/kind/releases/tag/v0.11.1) (released on 2021-05-28)
- Ingress Nginx Controller: [v0.48.0](https://github.com/kubernetes/ingress-nginx/releases/tag/controller-v0.48.1) (released on 2021-07-15)
- Conftest: [0.25.0](https://github.com/open-policy-agent/conftest/releases/tag/v0.25.0) (released on 2021-05-08)
- Istio: [1.12.2](https://github.com/istio/istio/releases/tag/1.12.2) (released on 2022-01-19)

# Contents

- Cluster Setup
    - [kubernetes-the-hard-way](kubernetes-the-hard-way)
    - [Kubeadm in local](kubeadm-local)
    - [kind](local-cluster/kind)
- Kubernetes Features
    - [Autoscaler HPA with custom metrics](autoscaler/hpa/custom-metrics)
    - [amazon-eks-workshop](eksworkshop)
- Kubernetes Extentions
    - [kubernetes-operator](kubernetes-operator)
    - [kubernetes-scheduler](kubernetes-scheduler)
    - plugins
- Networking
    - [traefik](traefik)
    - [ingress-nginx-controller](ingress-nginx-controller)
- Middleware (Operator)
    - [strimzi](strimzi)
    - [postgres-operator](postgres-operator)
    - [eck](eck)
- Service Proxy, Discovery, and, Mesh
    - [istio](istio)
- Monitoring
    - [Prometheus](prometheus-operator)
    - [Grafana](grafana)
    - [Loki](loki)
- Security
    - [open-policy-agent](open-policy-agent)
- Yaml Management
    - [Helm](helm)
    - [Helm vs Kustomize](helm-vs-kustomize)
- CI/CD
    - [conftest](open-policy-agent/conftest)
    - [argocd](argocd)

# Cloud Native Trail Map

- https://github.com/cncf/trailmap
- https://www.cncf.io/blog/2018/03/08/introducing-the-cloud-native-landscape-2-0-interactive-edition/

![alt text](https://github.com/cncf/trailmap/blob/master/CNCF_TrailMap_latest.png?raw=true)

## 1. CONTAINERIZATION

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

[kubeadm-local]()
#### Set up Kubernetes Cluster on GCP (kubernetes-the-hard-way)

[Kubernetes the hard way](https://github.com/kelseyhightower/kubernetes-the-hard-way)

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

### 4.1. [Prometheus](prometheus-operator)

![](prometheus-operator/prometheus.drawio.svg)

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


### TBD
- Vitess
- Rook
- etcd
- TiKV

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
