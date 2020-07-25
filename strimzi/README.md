# Set up

https://strimzi.io/quickstarts/

## Strimzi Operator

prepare namespace

```
[20-07-25 13:35:56] nakamasato at Masatos-MacBook-Pro in ~/Code/MasatoNaka/kubernetes-training/strimzi on master ✘
± kubectl create namespace kafka-strimzi-18
namespace/kafka-strimzi-18 created
```

prepare strimzi operator

```
curl https://strimzi.io/install/latest\?namespace\=kafka-strimzi-18 > 0.18/strimzi-operator.yaml
```

```
[20-07-25 13:36:02] nakamasato at Masatos-MacBook-Pro in ~/Code/MasatoNaka/kubernetes-training/strimzi on master ✘
± namespace=kafka-strimzi-18

[20-07-25 13:36:19] nakamasato at Masatos-MacBook-Pro in ~/Code/MasatoNaka/kubernetes-training/strimzi on master ✘
± kubectl apply -f "https://strimzi.io/install/latest?namespace=$namespace" -n $namespace
customresourcedefinition.apiextensions.k8s.io/kafkas.kafka.strimzi.io created
rolebinding.rbac.authorization.k8s.io/strimzi-cluster-operator-entity-operator-delegation created
clusterrolebinding.rbac.authorization.k8s.io/strimzi-cluster-operator created
rolebinding.rbac.authorization.k8s.io/strimzi-cluster-operator-topic-operator-delegation created
customresourcedefinition.apiextensions.k8s.io/kafkausers.kafka.strimzi.io created
customresourcedefinition.apiextensions.k8s.io/kafkarebalances.kafka.strimzi.io created
customresourcedefinition.apiextensions.k8s.io/kafkamirrormaker2s.kafka.strimzi.io created
clusterrole.rbac.authorization.k8s.io/strimzi-entity-operator created
clusterrole.rbac.authorization.k8s.io/strimzi-cluster-operator-global created
clusterrolebinding.rbac.authorization.k8s.io/strimzi-cluster-operator-kafka-broker-delegation created
rolebinding.rbac.authorization.k8s.io/strimzi-cluster-operator created
clusterrole.rbac.authorization.k8s.io/strimzi-cluster-operator-namespaced created
clusterrole.rbac.authorization.k8s.io/strimzi-topic-operator created
serviceaccount/strimzi-cluster-operator created
clusterrole.rbac.authorization.k8s.io/strimzi-kafka-broker created
customresourcedefinition.apiextensions.k8s.io/kafkatopics.kafka.strimzi.io created
customresourcedefinition.apiextensions.k8s.io/kafkabridges.kafka.strimzi.io created
deployment.apps/strimzi-cluster-operator created
customresourcedefinition.apiextensions.k8s.io/kafkaconnectors.kafka.strimzi.io created
customresourcedefinition.apiextensions.k8s.io/kafkaconnects2is.kafka.strimzi.io created
customresourcedefinition.apiextensions.k8s.io/kafkaconnects.kafka.strimzi.io created
customresourcedefinition.apiextensions.k8s.io/kafkamirrormakers.kafka.strimzi.io created
```

## Kafka Cluster

```
kubectl apply -f https://strimzi.io/examples/latest/kafka/kafka-persistent-single.yaml -n $namespace
kafka.kafka.strimzi.io/my-cluster created
```

Ready

```
○ kubectl get pod -n $namespace

NAME                                          READY   STATUS    RESTARTS   AGE
my-cluster-entity-operator-579cdc77bc-94zth   3/3     Running   0          63s
my-cluster-kafka-0                            2/2     Running   0          2m12s
my-cluster-zookeeper-0                        1/1     Running   0          3m40s
strimzi-cluster-operator-6c9d899778-pdvdj     1/1     Running   0          10m
```

## Test with console-producer & console-consumer

producer

```
kubectl -n $namespace run kafka-producer -ti --image=strimzi/kafka:0.18.0-kafka-2.5.0 --rm=true --restart=Never -- bin/kafka-console-producer.sh --broker-list my-cluster-kafka-bootstrap:9092 --topic my-topic

If you don't see a command prompt, try pressing enter.
>test
[2020-07-25 04:48:36,949] WARN [Producer clientId=console-producer] Error while fetching metadata with correlation id 3 : {my-topic=LEADER_NOT_AVAILABLE} (org.apache.kafka.clients.NetworkClient)
[2020-07-25 04:48:37,059] WARN [Producer clientId=console-producer] Error while fetching metadata with correlation id 4 : {my-topic=LEADER_NOT_AVAILABLE} (org.apache.kafka.clients.NetworkClient)
[2020-07-25 04:48:37,180] WARN [Producer clientId=console-producer] Error while fetching metadata with correlation id 5 : {my-topic=LEADER_NOT_AVAILABLE} (org.apache.kafka.clients.NetworkClient)
[2020-07-25 04:48:37,296] WARN [Producer clientId=console-producer] Error while fetching metadata with correlation id 6 : {my-topic=LEADER_NOT_AVAILABLE} (org.apache.kafka.clients.NetworkClient)
>te
>test
>test
>test
>test
```

consumer

```
kubectl -n $namespace run kafka-consumer -ti --image=strimzi/kafka:0.18.0-kafka-2.5.0 --rm=true --restart=Never -- bin/kafka-console-consumer.sh --bootstrap-server my-cluster-kafka-bootstrap:9092 --topic my-topic --from-beginning

If you don't see a command prompt, try pressing enter.



test
te
test
test
test
test
```

topic

```
kubectl get KafkaTopic -n $namespace

NAME                                                          PARTITIONS   REPLICATION FACTOR
consumer-offsets---84e7a678d08f4bd226872e5cdd4eb527fadc1c6a   50           1
my-topic                                                      1            1
```

# Enable the Cluster Operator to watch multiple namespaces


