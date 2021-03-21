# Autoscaler

## Overview

1. Cluster Autoscaler (CA)
1. Horizontal Pod Autoscaler (HPA)
1. Virtical Pod Autoscaler (VPA)

## Details

### CA

[cluster-autoscaler](https://github.com/kubernetes/autoscaler/tree/master/cluster-autoscaler)

### HPA

[horizontal-pod-autoscaler](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/#support-for-custom-metrics)

### VPA

[vertical-pod-autoscaler](https://github.com/kubernetes/autoscaler/tree/master/vertical-pod-autoscaler)

- [Keeping limit proportional to request](https://github.com/kubernetes/autoscaler/tree/master/vertical-pod-autoscaler#keeping-limit-proportional-to-request) Automatically set resource limit values based on limit to request ratios specified as part of the container template.

prerequisite

- Install metrics server
    ```
    kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
    ```


```
git clone https://github.com/kubernetes/autoscaler.git
cd autoscaler/vertical-pod-autoscaler
./hack/vpa-up.sh # if you get `unknown option -addext`, ./hack/vpa-up.sh on the [0.8 release branch](https://github.com/kubernetes/autoscaler/tree/vpa-release-0.8)
```

```
brew update
brew upgrade openssl
brew list openssl
```

```
echo 'export PATH=/usr/local/Cellar/openssl@1.1/1.1.1i/bin:$PATH' >> ~/.zshrc
```

```
which openssl
/usr/local/Cellar/openssl@1.1/1.1.1i/bin/openssl
openssl version
OpenSSL 1.1.1i  8 Dec 2020
```

```
./hack/vpa-up.sh
```

<details><summary>result</summary>

```
customresourcedefinition.apiextensions.k8s.io/verticalpodautoscalercheckpoints.autoscaling.k8s.io created
customresourcedefinition.apiextensions.k8s.io/verticalpodautoscalers.autoscaling.k8s.io created
clusterrole.rbac.authorization.k8s.io/system:metrics-reader created
clusterrole.rbac.authorization.k8s.io/system:vpa-actor created
clusterrole.rbac.authorization.k8s.io/system:vpa-checkpoint-actor created
clusterrole.rbac.authorization.k8s.io/system:evictioner created
clusterrolebinding.rbac.authorization.k8s.io/system:metrics-reader created
clusterrolebinding.rbac.authorization.k8s.io/system:vpa-actor created
clusterrolebinding.rbac.authorization.k8s.io/system:vpa-checkpoint-actor created
clusterrole.rbac.authorization.k8s.io/system:vpa-target-reader created
clusterrolebinding.rbac.authorization.k8s.io/system:vpa-target-reader-binding created
clusterrolebinding.rbac.authorization.k8s.io/system:vpa-evictionter-binding created
serviceaccount/vpa-admission-controller created
clusterrole.rbac.authorization.k8s.io/system:vpa-admission-controller created
clusterrolebinding.rbac.authorization.k8s.io/system:vpa-admission-controller created
clusterrole.rbac.authorization.k8s.io/system:vpa-status-reader created
clusterrolebinding.rbac.authorization.k8s.io/system:vpa-status-reader-binding created
serviceaccount/vpa-updater created
deployment.apps/vpa-updater created
serviceaccount/vpa-recommender created
deployment.apps/vpa-recommender created
Generating certs for the VPA Admission Controller in /tmp/vpa-certs.
Generating RSA private key, 2048 bit long modulus (2 primes)
.........................................................................................................................................................................................................................................................................................+++++
....+++++
e is 65537 (0x010001)
Generating RSA private key, 2048 bit long modulus (2 primes)
...+++++
....................+++++
e is 65537 (0x010001)
Signature ok
subject=CN = vpa-webhook.kube-system.svc
Getting CA Private Key
Uploading certs to the cluster.
secret/vpa-tls-certs created
Deleting /tmp/vpa-certs.
deployment.apps/vpa-admission-controller created
service/vpa-webhook created
```

</details>

```
kubectl get pod -n kube-system | grep vpa
vpa-admission-controller-664b5997b7-48ck9   1/1     Running   0          78s
vpa-recommender-5b768c88-nfj6k              1/1     Running   0          78s
vpa-updater-c9ffff655-6f7sb                 1/1     Running   0          78s
```

example

```
kubectl create -f examples/hamster.yaml
```

```
kubectl get pod -o jsonpath='{.items[].spec.containers[].resources}' | jq
{
  "requests": {
    "cpu": "100m",
    "memory": "50Mi"
  }
}
```

Five minutes later, ... nothing happened..
