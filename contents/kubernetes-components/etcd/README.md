# etcd

## Install

```
brew install etcd
```

```
etcd --version
etcd Version: 3.5.7
Git SHA: 215b53cf3
Go Version: go1.19.5
Go OS/Arch: darwin/arm64
```

## Quickstart

1. Start etcd.

    ```
    etcd
    ```
1. Set a key `greeting` and a value `Hello, etcd`.
    ```
    etcdctl put greeting "Hello, etcd"
    OK
    ```
1. Get the value for the key `greeting`
    ```
    etcdctl get greeting
    greeting
    Hello, etcd
    ```

## Kubernetes objects in etcd

Kubernetes objects are stored under the `/registry` path

You can see all the keys under `/registry` with the following command:

```
etcdctl get /registry/ --prefix --keys-only
```

<details>

```
/registry/apiregistration.k8s.io/apiservices/v1.

/registry/apiregistration.k8s.io/apiservices/v1.admissionregistration.k8s.io

/registry/apiregistration.k8s.io/apiservices/v1.apiextensions.k8s.io

/registry/apiregistration.k8s.io/apiservices/v1.apps

/registry/apiregistration.k8s.io/apiservices/v1.authentication.k8s.io

/registry/apiregistration.k8s.io/apiservices/v1.authorization.k8s.io

/registry/apiregistration.k8s.io/apiservices/v1.autoscaling

/registry/apiregistration.k8s.io/apiservices/v1.batch

/registry/apiregistration.k8s.io/apiservices/v1.certificates.k8s.io

/registry/apiregistration.k8s.io/apiservices/v1.coordination.k8s.io

/registry/apiregistration.k8s.io/apiservices/v1.discovery.k8s.io

/registry/apiregistration.k8s.io/apiservices/v1.events.k8s.io

/registry/apiregistration.k8s.io/apiservices/v1.networking.k8s.io

/registry/apiregistration.k8s.io/apiservices/v1.node.k8s.io

/registry/apiregistration.k8s.io/apiservices/v1.policy

/registry/apiregistration.k8s.io/apiservices/v1.rbac.authorization.k8s.io

/registry/apiregistration.k8s.io/apiservices/v1.scheduling.k8s.io

/registry/apiregistration.k8s.io/apiservices/v1.storage.k8s.io

/registry/apiregistration.k8s.io/apiservices/v1beta1.storage.k8s.io

/registry/apiregistration.k8s.io/apiservices/v1beta2.flowcontrol.apiserver.k8s.io

/registry/apiregistration.k8s.io/apiservices/v1beta3.flowcontrol.apiserver.k8s.io

/registry/apiregistration.k8s.io/apiservices/v2.autoscaling

/registry/configmaps/kube-system/extension-apiserver-authentication

/registry/endpointslices/default/kubernetes

/registry/flowschemas/catch-all

/registry/flowschemas/endpoint-controller

/registry/flowschemas/exempt

/registry/flowschemas/global-default

/registry/flowschemas/kube-controller-manager

/registry/flowschemas/kube-scheduler

/registry/flowschemas/kube-system-service-accounts

/registry/flowschemas/probes

/registry/flowschemas/service-accounts

/registry/flowschemas/system-leader-election

/registry/flowschemas/system-node-high

/registry/flowschemas/system-nodes

/registry/flowschemas/workload-leader-election

/registry/leases/kube-system/kube-apiserver-c7xw4gapfj47r73yd6xe7yiqea

/registry/masterleases/192.168.10.33

/registry/namespaces/default

/registry/namespaces/kube-node-lease

/registry/namespaces/kube-public

/registry/namespaces/kube-system

/registry/pods/default/nginx

/registry/priorityclasses/system-cluster-critical

/registry/priorityclasses/system-node-critical

/registry/prioritylevelconfigurations/catch-all

/registry/prioritylevelconfigurations/exempt

/registry/prioritylevelconfigurations/global-default

/registry/prioritylevelconfigurations/leader-election

/registry/prioritylevelconfigurations/node-high

/registry/prioritylevelconfigurations/system

/registry/prioritylevelconfigurations/workload-high

/registry/prioritylevelconfigurations/workload-low

/registry/ranges/serviceips

/registry/ranges/servicenodeports

/registry/serviceaccounts/default/default

/registry/services/endpoints/default/kubernetes

/registry/services/specs/default/kubernetes
```

</details>

You can also check a specific object e.g. nginx Pod

```
etcdctl get /registry/pods/default/nginx
```

<details>

```
/registry/pods/default/nginx
k8s

v1Pod�
�
nginxdefault"*$a77f3131-9ce0-4319-a7c2-ea859df720212����Z

runnginx��

kubectl-runUpdatev����FieldsV1:�
�{"f:metadata":{"f:labels":{".":{},"f:run":{}}},"f:spec":{"f:containers":{"k:{\"name\":\"nginx\"}":{".":{},"f:image":{},"f:imagePullPolicy":{},"f:name":{},"f:resources":{},"f:terminationMessagePath":{},"f:terminationMessagePolicy":{}}},"f:dnsPolicy":{},"f:enableServiceLinks":{},"f:restartPolicy":{},"f:schedulerName":{},"f:securityContext":{},"f:terminationGracePeriodSeconds":{}}}B�
�
kube-api-access-2f85qk�h
"

�token
(&

kube-root-ca.crt
ca.crtca.crt
)'
%
        namespace
v1metadata.namespace��
nginxnginx*BJL
kube-api-access-2f85q-/var/run/secrets/kubernetes.io/serviceaccount"2j/dev/termination-logrAlways����FileAlways 2
                                                        ClusterFirstBdefaultJdefaultRX`hr���default-scheduler�6
node.kubernetes.io/not-readyExists"     NoExecute(��8
node.kubernetes.io/unreachableExists"   NoExecute(�����PreemptLowerPriority
Pending"*2J
BestEffortZ"
```

</details>

However, the data is encoded. You can decode with https://github.com/jpbetz/auger

```
git clone https://github.com/jpbetz/auger && cd auger
go build -o anger main.go
```

```
etcdctl get /registry/pods/default/nginx | ~/repos/jpbetz/auger/anger decode
```

<details>

```
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: "2023-03-25T00:21:26Z"
  labels:
    run: nginx
  name: nginx
  namespace: default
  uid: a77f3131-9ce0-4319-a7c2-ea859df72021
spec:
  containers:
  - image: nginx
    imagePullPolicy: Always
    name: nginx
    resources: {}
    terminationMessagePath: /dev/termination-log
    terminationMessagePolicy: File
    volumeMounts:
    - mountPath: /var/run/secrets/kubernetes.io/serviceaccount
      name: kube-api-access-2f85q
      readOnly: true
  dnsPolicy: ClusterFirst
  priority: 0
  restartPolicy: Always
  schedulerName: default-scheduler
  securityContext: {}
  serviceAccount: default
  serviceAccountName: default
  terminationGracePeriodSeconds: 30
  tolerations:
  - effect: NoExecute
    key: node.kubernetes.io/not-ready
    operator: Exists
    tolerationSeconds: 300
  - effect: NoExecute
    key: node.kubernetes.io/unreachable
    operator: Exists
    tolerationSeconds: 300
  volumes:
  - name: kube-api-access-2f85q
    projected:
      defaultMode: 420
      sources:
      - {}
      - configMap:
          items:
          - key: ca.crt
            path: ca.crt
          name: kube-root-ca.crt
      - downwardAPI:
          items:
          - fieldRef:
              apiVersion: v1
              fieldPath: metadata.namespace
            path: namespace
status:
  phase: Pending
  qosClass: BestEffort
```

</details>


For more details about encoding in Kubernetes, you can read [Kubernetes Protobuf encoding](https://kubernetes.io/docs/reference/using-api/api-concepts/#protobuf-encoding)

## Reference
- https://etcd.io
- https://technekey.com/check-whats-inside-the-etcd-database-in-kubernetes/
