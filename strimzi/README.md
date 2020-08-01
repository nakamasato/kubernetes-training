# Set up

https://strimzi.io/quickstarts/


## Download Strimzi (0.18)

Download `strimzi-0.18.0.zip`

https://github.com/strimzi/strimzi-kafka-operator/releases/tag/0.18.0

1. Put it under `base` and unzip
2. Create `kustomization.yaml`
3. Create `overlays/<younamespace>`
4. Add `rolebinding` + `clusterrolebinding` to overwrite `namespace`


## Strimzi Operator

prepare namespace

```
[20-07-25 13:35:56] nakamasato at Masatos-MacBook-Pro in ~/Code/MasatoNaka/kubernetes-training/strimzi on master ✘
± kubectl create namespace kafka-strimzi-18
namespace/kafka-strimzi-18 created
```

prepare strimzi operator


```
[20-07-25 13:36:02] nakamasato at Masatos-MacBook-Pro in ~/Code/MasatoNaka/kubernetes-training/strimzi on master ✘
± namespace=kafka-strimzi-18

[20-07-25 13:36:19] nakamasato at Masatos-MacBook-Pro in ~/Code/MasatoNaka/kubernetes-training/strimzi on master ✘
± kubectl apply -k overlays/$namespace
```

## Kafka Cluster

```
curl https://strimzi.io/examples/latest/kafka/kafka-persistent-single.yaml > overlays/$namespace/my-cluster.yaml
```

```
kubectl apply -k overlays/$namespace
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

## KafkaTopic

Example: https://github.com/strimzi/strimzi-kafka-operator/blob/master/examples/topic/kafka-topic.yaml

```
kubectl get KafkaTopic -n $namespace
NAME                                                          PARTITIONS   REPLICATION FACTOR
consumer-offsets---84e7a678d08f4bd226872e5cdd4eb527fadc1c6a   50           1
my-topic                                                      1            1
```

## KafkaUser

Example: https://github.com/strimzi/strimzi-kafka-operator/blob/master/examples/user/kafka-user.yaml

```
kubectl get KafkaUser -n $namespace
NAME      AUTHENTICATION   AUTHORIZATION
my-user   tls              simple
```

## KafkaConnect

https://strimzi.io/docs/0.16.2/full.html#deploying-kafka-connect-str

```
  annotations:
    strimzi.io/use-connector-resources: "true" # to enable connector resource
```

```
overlays/kafka-strimzi-18/connect/source/connect-source.yaml
verlays/kafka-strimzi-18/connect/source/my-source-connector.yaml
```

```
kubectl get KafkaConnect
NAME                   DESIRED REPLICAS
kafka-connect-source   2
```

```
kubectl get KafkaConnector
NAME                  AGE
my-source-connector   9m2s
```

# Enable the Cluster Operator to watch multiple namespaces

https://strimzi.io/docs/0.16.2/full.html#deploying-cluster-operator-to-watch-multiple-namespacesstr

## Change STRIMZI_NAMESPACE

TO update `STRIMZI_NAMESPACE`, add a patch yaml and include it in `kustomization.yaml` in `kafka-strimzi-18` (Reference: [Kustomize でマニフェストのフィールドを削除する](https://text.superbrothers.dev/200315-delete-field-with-kustomize/))

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: strimzi-cluster-operator
spec:
  template:
    spec:
      serviceAccountName: strimzi-cluster-operator
      containers:
      - name: strimzi-cluster-operator
        env:
        - name: STRIMZI_NAMESPACE
          value: kafka-strimzi-18,kafka-strimzi-18-staging
          valueFrom: null
```

Diff

```
kubectl diff -k overlays/kafka-strimzi-18
-  generation: 1
+  generation: 2
   labels:
     app: strimzi
   name: strimzi-cluster-operator
@@ -36,10 +36,7 @@
         - /opt/strimzi/bin/cluster_operator_run.sh
         env:
         - name: STRIMZI_NAMESPACE
-          valueFrom:
-            fieldRef:
-              apiVersion: v1
-              fieldPath: metadata.namespace
+          value: kafka-strimzi-18,kafka-strimzi-18-staging
         - name: STRIMZI_FULL_RECONCILIATION_INTERVAL_MS
```

Apply

```
kubectl apply -k overlays/kafka-strimzi-18
customresourcedefinition.apiextensions.k8s.io/kafkabridges.kafka.strimzi.io unchanged
customresourcedefinition.apiextensions.k8s.io/kafkaconnectors.kafka.strimzi.io unchanged
customresourcedefinition.apiextensions.k8s.io/kafkaconnects.kafka.strimzi.io unchanged
customresourcedefinition.apiextensions.k8s.io/kafkaconnects2is.kafka.strimzi.io unchanged
customresourcedefinition.apiextensions.k8s.io/kafkamirrormaker2s.kafka.strimzi.io unchanged
customresourcedefinition.apiextensions.k8s.io/kafkamirrormakers.kafka.strimzi.io unchanged
customresourcedefinition.apiextensions.k8s.io/kafkarebalances.kafka.strimzi.io unchanged
customresourcedefinition.apiextensions.k8s.io/kafkas.kafka.strimzi.io unchanged
customresourcedefinition.apiextensions.k8s.io/kafkatopics.kafka.strimzi.io unchanged
customresourcedefinition.apiextensions.k8s.io/kafkausers.kafka.strimzi.io unchanged
serviceaccount/strimzi-cluster-operator unchanged
clusterrole.rbac.authorization.k8s.io/strimzi-cluster-operator-global unchanged
clusterrole.rbac.authorization.k8s.io/strimzi-cluster-operator-namespaced unchanged
clusterrole.rbac.authorization.k8s.io/strimzi-entity-operator unchanged
clusterrole.rbac.authorization.k8s.io/strimzi-kafka-broker unchanged
clusterrole.rbac.authorization.k8s.io/strimzi-topic-operator unchanged
rolebinding.rbac.authorization.k8s.io/strimzi-cluster-operator-entity-operator-delegation unchanged
rolebinding.rbac.authorization.k8s.io/strimzi-cluster-operator-topic-operator-delegation unchanged
rolebinding.rbac.authorization.k8s.io/strimzi-cluster-operator unchanged
clusterrolebinding.rbac.authorization.k8s.io/strimzi-cluster-operator-kafka-broker-delegation unchanged
clusterrolebinding.rbac.authorization.k8s.io/strimzi-cluster-operator unchanged
deployment.apps/strimzi-cluster-operator configured
kafka.kafka.strimzi.io/my-cluster unchanged
```

## RoleBinding & ClusterRoleBinding

Copy `RoleBinding` from `kafka-strimzi-18`

```
mkdir -p overlays/kafka-strimzi-18-staging/strimzi-0.18.0/install/cluster-operator
cp overlays/kafka-strimzi-18/strimzi-0.18.0/install/cluster-operator/*-RoleBinding*yaml overlays/kafka-strimzi-18-staging/strimzi-0.18.0/install/cluster-operator
```

Apply

```
kubectl apply -k overlays/kafka-strimzi-18-staging
rolebinding.rbac.authorization.k8s.io/strimzi-cluster-operator-entity-operator-delegation created
rolebinding.rbac.authorization.k8s.io/strimzi-cluster-operator-topic-operator-delegation created
rolebinding.rbac.authorization.k8s.io/strimzi-cluster-operator created
```

## Kafka Cluster

Prepare `my-cluster.yaml` and add it to `kustomization.yaml`

```
cp overlays/kafka-strimzi-18/my-cluster.yaml overlays/kafka-strimzi-18-staging
```

Apply

```
kubectl apply -k overlays/kafka-strimzi-18-staging
rolebinding.rbac.authorization.k8s.io/strimzi-cluster-operator-entity-operator-delegation unchanged
rolebinding.rbac.authorization.k8s.io/strimzi-cluster-operator-topic-operator-delegation unchanged
rolebinding.rbac.authorization.k8s.io/strimzi-cluster-operator unchanged
kafka.kafka.strimzi.io/my-cluster created
```

Check

```
kubectl get pod -n $namespace-staging
NAME                                         READY   STATUS    RESTARTS   AGE
my-cluster-entity-operator-fd45b849f-9vk62   3/3     Running   0          59s
my-cluster-kafka-0                           2/2     Running   0          2m23s
my-cluster-zookeeper-0                       1/1     Running   0          3m19s
```

# Authentication & Authorization

[Listerner authentication](https://strimzi.io/docs/operators/master/using.html#assembly-kafka-broker-listener-authentication-deployment-configuration-kafka)

- Mutual TLS authentication
    - The client supports authenticaton using mutual TLS authentication
    - It is necessary to ue the TLS certificates rather than passwords
    - You can reconfigure and restart client applications periodically so that they do not use expired certificates
- SCRAM-SHA(Salted Challenge Response Authenticatoin Mechanism) authentication
    - Support `SCRAM-SHA-512` only.
    - The client supports authentication using SCRAM-SHA-512
    - It is necessary to use passwords rather than the TLS certificates
    - Authentication for unencrypted communication is required
- no `authentication` property -> not authenticate


# reference

- Custom image for KafkaConnect: https://strimzi.io/docs/operators/0.18.0/using.html#creating-new-image-from-base-str
- https://github.com/nakamasato/kafka-connect
