# Strimzi Kafka Operator

This example uses **Strimzi 1.2.0** and Kafka 4.3.1. Strimzi 1.x only supports the `kafka.strimzi.io/v1` API; old `v1beta2` resources must be converted before upgrading.

## Install

```bash
kubectl apply -f namespace.yaml
kubectl apply -f https://github.com/strimzi/strimzi-kafka-operator/releases/download/1.2.0/strimzi-cluster-operator-1.2.0.yaml -n kafka
kubectl get pods -n kafka
```

## Create a Kafka cluster and topic

The example uses KRaft (ZooKeeper is not used):

```bash
kubectl apply -k kafka-cluster
kubectl wait --for=condition=Ready kafka/my-cluster --timeout=10m -n kafka
kubectl get kafka,kafkatopic -n kafka
```

Produce and consume with `quay.io/strimzi/kafka:0.49.0-kafka-4.3.1` using `my-cluster-kafka-bootstrap:9092`.

## Cleanup

```bash
kubectl delete -k kafka-cluster
kubectl delete -f https://github.com/strimzi/strimzi-kafka-operator/releases/download/1.2.0/strimzi-cluster-operator-1.2.0.yaml -n kafka
kubectl delete namespace kafka
```

See the [Strimzi 1.2.0 documentation](https://strimzi.io/docs/operators/1.2.0/).
