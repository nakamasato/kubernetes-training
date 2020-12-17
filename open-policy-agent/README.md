# Open Policy Agent

## Getting Started

```
[20-09-27 19:49:20] masato-naka at PCN-537 in ~/repos/MasatoNaka/kubernetes-training/open-policy-agent/getting-started on master ✘
± curl -X PUT -H 'Content-Type:application/json' --data-binary @subordinates.json \
localhost:8181/v1/data/subordinates


[20-09-27 19:49:48] masato-naka at PCN-537 in ~/repos/MasatoNaka/kubernetes-training/open-policy-agent/getting-started on master ✘
± curl -s localhost:8181/v1/data/subordinates | jq .

{
  "result": {
    "alice": [
      "bob"
    ],
    "bob": [],
    "charlie": [
{
      "david"
    ],
    "david": []
  }
}

[20-09-27 19:51:40] masato-naka at PCN-537 in ~/repos/MasatoNaka/kubernetes-training/open-policy-agent/getting-started on master ✘
± curl -X PUT -H 'Content-Type: text/plain' --data-binary @httpapi_authz.repo \
  localhost:8181/v1/policies/httpapi_authz

{}%

[20-09-27 19:53:17] masato-naka at PCN-537 in ~/repos/MasatoNaka/kubernetes-training/open-policy-agent/getting-started on master ✘
± cat alice_to_alice.json| jq
{
  "input": {
    "method": "GET",
    "path": [
      "finance",
      "salary",
      "alice"
    ],
    "user": "alice"
  }
}

[20-09-27 19:53:18] masato-naka at PCN-537 in ~/repos/MasatoNaka/kubernetes-training/open-policy-agent/getting-started on master ✘
± curl -s -X POST -H 'Content-Type:application/json' --data-binary @alice_to_alice.json \
    localhost:8181/v1/data/httpapi/authz/allow | jq .

{
  "result": true
}

[20-09-27 19:53:52] masato-naka at PCN-537 in ~/repos/MasatoNaka/kubernetes-training/open-policy-agent/getting-started on master ✘
± curl -s -X POST -H 'Content-Type:application/json' --data-binary @alice_to_bob.json \
    localhost:8181/v1/data/httpapi/authz/allow | jq .

{
  "result": true
}

[20-09-27 19:54:34] masato-naka at PCN-537 in ~/repos/MasatoNaka/kubernetes-training/open-policy-agent/getting-started on master ✘
± curl -s -X POST -H 'Content-Type:application/json' --data-binary @alice_to_david.json \
    localhost:8181/v1/data/httpapi/authz/allow | jq .

{
  "result": false
}
```

## Gatekeeper

https://github.com/open-policy-agent/gatekeeper

1. Install

```
kubectl apply -f https://raw.githubusercontent.com/open-policy-agent/gatekeeper/master/deploy/gatekeeper.yaml
namespace/gatekeeper-system created
customresourcedefinition.apiextensions.k8s.io/configs.config.gatekeeper.sh created
customresourcedefinition.apiextensions.k8s.io/constrainttemplates.templates.gatekeeper.sh created
serviceaccount/gatekeeper-admin created
role.rbac.authorization.k8s.io/gatekeeper-manager-role created
clusterrole.rbac.authorization.k8s.io/gatekeeper-manager-role created
rolebinding.rbac.authorization.k8s.io/gatekeeper-manager-rolebinding created
clusterrolebinding.rbac.authorization.k8s.io/gatekeeper-manager-rolebinding created
secret/gatekeeper-webhook-server-cert created
service/gatekeeper-webhook-service created
deployment.apps/gatekeeper-audit created
deployment.apps/gatekeeper-controller-manager created
validatingwebhookconfiguration.admissionregistration.k8s.io/gatekeeper-validating-webhook-configuration created
```

1. Install `ConstraintTemplate` (CRD) to require `label`

```
kubectl apply -f https://raw.githubusercontent.com/open-policy-agent/gatekeeper/master/demo/basic/templates/k8srequiredlabels_template.yaml
constrainttemplate.templates.gatekeeper.sh/k8srequiredlabels created
```

```
○ kubectl get ConstraintTemplate

NAME                AGE
k8srequiredlabels   45s
```

1. Create `Constraint`

```
kubectl apply -f https://raw.githubusercontent.com/open-policy-agent/gatekeeper/master/demo/basic/constraints/all_ns_must_have_gatekeeper.yaml
k8srequiredlabels.constraints.gatekeeper.sh/ns-must-have-gk created
```

```
kubectl get K8sRequiredLabels

NAME              AGE
ns-must-have-gk   75s
```

1. Check

```
[20-08-05 23:34:43] nakamasato at Masatos-MacBook-Pro in ~/Code/MasatoNaka/kubernetes-training/open-policy-agent on postgres-operator ✘
± kubectl apply -f gatekeeper/examples/valid-namespace.yaml --dry-run=server
namespace/valid-namespace created (server dry run)

[20-08-05 23:35:14] nakamasato at Masatos-MacBook-Pro in ~/Code/MasatoNaka/kubernetes-training/open-policy-agent on postgres-operator ✘
± kubectl apply -f gatekeeper/examples/invalid-namespace.yaml --dry-run=server
Error from server ([denied by ns-must-have-gk] you must provide labels: {"gatekeeper"}): error when creating "gatekeeper/invalid-namespace.yaml": admission webhook "validation.gatekeeper.sh" denied the request: [denied by ns-must-have-gk] you must provide labels: {"gatekeeper"}
```

## Example 1

```
kubectl apply -f gatekeeper/require-labels/k8srequiredlabels.yaml
kubectl apply -f gatekeeper/require-labels/k8srequiredlabels-ns.yaml
kubectl create ns naka
Error from server ([denied by ns-must-have-gk] you must provide labels: {"gatekeeper"}
[denied by ns-must-have-hr] you must provide labels: {"hr"}): admission webhook "validation.gatekeeper.sh" denied the request: [denied by ns-must-have-gk] you must provide labels: {"gatekeeper"}
[denied by ns-must-have-hr] you must provide labels: {"hr"}
```

## Conftest

https://github.com/open-policy-agent/conftest

# FAQ

1. Run on GKE

https://github.com/open-policy-agent/gatekeeper#running-on-private-gke-cluster-nodes

# Study steps

- [[Youtube] Deep Dive: Open Policy Agent - Torin Sandall & Tim Hinrichs, Styra (2019/05/23)](https://www.youtube.com/watch?v=n94_FNhuzy4&feature=youtu.be)
- [[Kubernetes Blog] A Guide to Kubernetes Admission Controllers (2020/03/21)](https://kubernetes.io/blog/2019/03/21/a-guide-to-kubernetes-admission-controllers/)
  - `ValidatingAdmissionWebhooks`
  - `MutatingAdmissionWebhooks`

  > We will examine these two admission controllers closely, as they do not implement any policy decision logic themselves. Instead, the respective action is obtained from a REST endpoint (a webhook) of a service running inside the cluster. This approach decouples the admission controller logic from the Kubernetes API server, thus allowing users to implement custom logic to be executed whenever resources are created, updated, or deleted in a Kubernetes cluster.

- [EKS Enables Support for Kubernetes Dynamic Admission Controllers (2018/10/12)](https://aws.amazon.com/about-aws/whats-new/2018/10/amazon-eks-enables-support-for-kubernetes-dynamic-admission-cont/)

  - 1.17 ([Platform versions](https://docs.aws.amazon.com/eks/latest/userguide/platform-versions.html)): `NamespaceLifecycle, LimitRanger, ServiceAccount, DefaultStorageClass, ResourceQuota, DefaultTolerationSeconds, NodeRestriction, MutatingAdmissionWebhook, ValidatingAdmissionWebhook, PodSecurityPolicy, TaintNodesByCondition, Priority, StorageObjectInUseProtection, PersistentVolumeClaimResize`

- [Dynamic Admission Control](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/)
- [Integrating Open Policy Agent (OPA) With Kubernetes](https://www.magalix.com/blog/integrating-open-policy-agent-opa-with-kubernetes-a-deep-dive-tutorial)
