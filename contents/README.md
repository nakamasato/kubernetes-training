# kubernetes-training

# Versions

1. Kubernetes: [v1.25.3](https://github.com/kubernetes/kubernetes/releases/tag/v1.25.3)
1. kustomize: [v4.2.0](https://github.com/kubernetes-sigs/kustomize/releases/tag/kustomize%2Fv4.2.0) (released on 2021-07-02)
1. Helm: [v3.11.2](https://github.com/helm/helm/releases/tag/v3.11.2)
1. Traefik: [v2.9.0](https://github.com/traefik/traefik/releases/tag/v2.9.0)
1. ArgoCD: [v3.3.2](https://github.com/argoproj/argo-cd/releases/tag/v3.3.2)
1. Prometheus-Operator: [v0.43.1](https://github.com/prometheus-operator/prometheus-operator/releases/tag/v0.43.1)
1. Prometheus: [latest](https://github.com/prometheus/prometheus/releases)
1. Grafana: [latest](https://github.com/grafana/grafana/releases)
1. Strimzi: [0.24.0](https://github.com/strimzi/strimzi-kafka-operator/releases/tag/0.24.0)
1. Kind: [v0.17.0](https://github.com/kubernetes-sigs/kind/releases/tag/v0.17.0)
1. Ingress Nginx Controller: [controller-v1.7.0](https://github.com/kubernetes/ingress-nginx/releases/tag/controller-v1.7.0)
1. Conftest: [v0.25.0](https://github.com/open-policy-agent/conftest/releases/tag/v0.25.0)
1. Istio: [1.19.0](https://github.com/istio/istio/releases/tag/1.19.0)
1. PostgresOperator: [v1.7.1](https://github.com/zalando/postgres-operator/releases/tag/v1.7.1) (released on 2021-11-04)
1. Cert Manager: [v1.7.1](https://github.com/cert-manager/cert-manager/releases/tag/v1.7.1) (released on 2022-02-05)

# Contents

Contents are organized based on Cloud Native Trail Map:

- https://github.com/cncf/trailmap
- https://www.cncf.io/blog/2018/03/08/introducing-the-cloud-native-landscape-2-0-interactive-edition/

![alt text](https://github.com/cncf/trailmap/blob/master/CNCF_TrailMap_latest.png?raw=true)

## 1. CONTAINERIZATION

1. [Containers 101: attach vs. exec - what's the difference?](https://iximiuz.com/en/posts/containers-101-attach-vs-exec/)

## 2. CI/CD

1. [ArgoCD](argocd)
1. [Conftest](open-policy-agent/conftest)
1. Kyverno: https://kyverno.io/
1. Polaris: https://www.fairwinds.com/polaris

## 3. ORCHESTRATION & APPLICATION DEFINITION

1. Kubernetes
    1. Useful Commands

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
    1. Kubernetes Cluster
        1. [local cluster](local-cluster): kind, minikube, Docker Desktop
        1. [kubeadm-local](kubeadm-local): Set up Kubernetes Cluster with kubeadm (local)
        1. [Kubernetes The Hard Way](https://github.com/kelseyhightower/kubernetes-the-hard-way): Set up Kubernetes Cluster on GCP (kubernetes-the-hard-way)
    1. Kubernetes Components
        1. [kubernetes-scheduler](kubernetes-components/kubernetes-scheduler)
        1. [etcd](kubernetes-components/etcd)
        1. [kube-apiserver](kubernetes-components/kube-apiserver)
        1. [kube-controller-manager](kubernetes-components/kube-controller-manager)
        1. [kube-proxy](kubernetes-components/kube-proxy)
        1. [kubelet](kubernetes-components/kubelet)
    1. [Kubernetes Operator](kubernetes-operator)
        1. [client-go](kubernetes-operator/client-go/)
        1. [apimachinery](kubernetes-operator/apimachinery)
        1. [controller-runtime](kubernetes-operator/controller-runtime/)
    1. [More Practices of Applications on Kubernetes](PRACTICE.md)
    1. [Kubernetes Features](kubernetes-features)
        1. [Autoscaler HPA with custom metrics](autoscaler/hpa/custom-metrics)
        1. [amazon-eks-workshop](eksworkshop)
    1. Kubernetes Extensions
        1. [kubernetes-operator](kubernetes-operator)
        1. [kubernetes-scheduler](kubernetes-extensions/kubernetes-scheduler)
        1. [plugins (todo)](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/)
    1. Namespaces
        1. [hierarchical namespaces (HNC)](https://github.com/kubernetes-sigs/hierarchical-namespaces)
    1. Deloyment Managemet
        1. [Knative](knative)
        1. Skaffold: https://skaffold.dev/ (ToDo)
    1. Middleware (Operator)
        1. [strimzi](strimzi)
        1. [eck](eck)
    1. Security
        1. [Cert Manager](cert-manager)
    1. Machine Learning
        1. [kubeflow](https://github.com/nakamasato/kubeflow-training)
1. [Helm](helm)
    1. [Helm vs Kustomize](helm-vs-kustomize)

## 4. OBSERVABILITY & ANALYTICS

1. [Prometheus](prometheus)
    1. [Prometheus Operator](prometheus-operator)
1. Jaeger: https://www.jaegertracing.io/
    1. [Opentelemetry & Jaeger](https://github.com/nakamasato/golang-training/tree/main/pragmatic-cases/opentelemetry)
1. Opentelemetry (ToDo)
1. fluentd (ToDo)
1. [Thanos (todo)] https://thanos.io/
1. [Grafana](grafana)
1. [Grafana Operator](grafana-operator)
1. [Grafana Loki](loki)
1. [Grafana Tempo](tempo)

## 5. SERVICE PROXY, DISCOVERY & MESH

1. [Istio](istio)
1. [Envoy](https://github.com/nakamasato/envoy-training)
1. CoreDNS (ToDo)
1. Linkerd (ToDo)

## 6. NETWORKING, POLICY & SECURITY

1. [Open Policy Agent](open-policy-agent)
    1. [gatekeeper](open-policy-agent/README.md#gatekeeper)
    1. [conftest](open-policy-agent/README.md#conftest)
1. CNI (ToDo)
1. falco (ToDo)
1. [Kubernetes Gateway API](kubernetes-gateway-api)
    1.  Envoy Gateway
    2.  Istio
    3.  Kong
    4.  NGINX Kubernetes Gateway
    1. [traefik](traefik)
1. Ingress
    1. [ingress-nginx-controller](ingress-nginx-controller)

## 7. DISTRIBUTED DATABASE & STORAGE

1. [etcd](kubernetes-components/etcd)
1. Vitess: https://github.com/vitessio/vitess (ToDo)
1. Rook: https://rook.io/ (ToDo)
1. TiDB: https://github.com/pingcap/tidb (ToDo)
1. TimescaleDB: https://github.com/timescale/timescaledb-kubernetes (ToDo)
1. Others: [Databases](databases)
    1. [mysql-operator](databases/mysql-operator)
    1. [postgres-operator](databases/postgres-operator)
## 8. STREAMING & MESSAGING

1. gRPC: https://grpc.io/ (ToDo)
1. NATS: https://nats.io/ (ToDo)
1. cloudevents: https://cloudevents.io/ (ToDo)

## 9. CONTAINER REGISTRY & RUNTIME

1. containerd: https://containerd.io/ (ToDo)
1. harbor: https://goharbor.io/ (ToDo)
1. cri-o: https://cri-o.io/ (ToDo)

## 10. SOFTWARE DISTRIBUTION

1. The Update Framework: https://theupdateframework.io/ (ToDo)
1. Notary: https://notaryproject.dev/ (ToDo)
