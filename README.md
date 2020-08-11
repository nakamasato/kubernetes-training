# kubernetes-training

# Version

- kubernetes: 1.15
- Helm: 3.2.1
- Traefik: 2.2

# Contents

- [helm](helm)
- [amazon-eks-workshop](amazon-eks-workshop)
- [ingress-nginx-controller](ingress-nginx-controller)
- [kubernetes-the-hard-way](kubernetes-the-hard-way)
- [traefik](traefik)
- [strimzi](strimzi)
- [postgres-operator](postgres-operator)

# Practice 1

- ES & Kibana

    ```
    kubectl create namespace eck
    ```

    ```
    helm repo add elastic https://helm.elastic.co
    ```

    ```
    helm install -n eck elasticsearch elastic/elasticsearch -f helm/es-config.yaml
    ```

    ```
    helm install -n eck kibana elastic/kibana -f helm/kb-config.yaml
    ```

- Kafka

    1. Update the kafka-connect-twitter with your own API token
    1. Apply Kafka

        ```
        kubectl create namespace kafka-strimzi-18
        kubectl apply -k strimzi/overlays/kafka-strimzi-18
        ```

![](docs/practice-01.drawio.svg)


```
NAMESPACE          NAME                                                             READY   STATUS    RESTARTS   AGE
eck                elasticsearch-master-0                                           1/1     Running   0          14h
eck                kibana-kibana-55f4bc96f5-7fz65                                   1/1     Running   0          14h
kafka-strimzi-18   kafka-connect-sink-connect-847cfbf66-gwtkl                       1/1     Running   0          7h27m
kafka-strimzi-18   kafka-connect-source-connect-57bf7974f7-sz8ww                    1/1     Running   0          7h27m
kafka-strimzi-18   my-cluster-entity-operator-579cdc77bc-v6rxt                      3/3     Running   5          14h
kafka-strimzi-18   my-cluster-kafka-0                                               2/2     Running   0          14h
kafka-strimzi-18   my-cluster-kafka-1                                               2/2     Running   0          14h
kafka-strimzi-18   my-cluster-kafka-2                                               2/2     Running   2          14h
kafka-strimzi-18   my-cluster-zookeeper-0                                           1/1     Running   0          14h
kafka-strimzi-18   strimzi-cluster-operator-6c9d899778-nkd9q                        1/1     Running   0          14h
kube-system        kube-dns-869d587df7-7whsm                                        3/3     Running   0          14h
kube-system        kube-dns-869d587df7-z659j                                        3/3     Running   0          14h
kube-system        kube-dns-autoscaler-645f7d66cf-r9ttj                             1/1     Running   0          14h
kube-system        kube-proxy-gke-my-gke-cluster-my-gke-cluster-nod-9dff1786-x4wz   1/1     Running   0          14h
kube-system        kube-proxy-gke-my-gke-cluster-my-gke-cluster-pre-19639e01-7jsz   1/1     Running   0          93s
kube-system        kube-proxy-gke-my-gke-cluster-my-gke-cluster-pre-19639e01-cnl2   1/1     Running   0          14h
kube-system        kube-proxy-gke-my-gke-cluster-my-gke-cluster-pre-19639e01-f6cb   1/1     Running   0          14h
kube-system        kube-proxy-gke-my-gke-cluster-my-gke-cluster-pre-19639e01-vw9d   1/1     Running   0          14h
kube-system        l7-default-backend-678889f899-fvswg                              1/1     Running   0          14h
kube-system        metrics-server-v0.3.6-7b7d6c7576-msl8x                           2/2     Running   0          14h
```

# Practice 2

- Prometheus & Grafana

    ```
    g clone https://github.com/coreos/kube-prometheus.git && kube-prometheus
    ```

    ```
    kubectl apply -f manifests/setup
    ```

    wait a few minutes

    ```
    kubectl create -f manifests
    ```

- Add strimzi monitoring

    ```
    kubectl apply -f strimzi/monitoring/prometheus-prometheus.yaml,strimzi/monitoring/prometheus-clusterRole.yaml
    ```

- Add elasticsearch monitoring

![](docs/practice-02.drawio.svg)
