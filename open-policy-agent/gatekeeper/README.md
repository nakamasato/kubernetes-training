# Gatekeeper

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
± kubectl apply -f gatekeeper/valid-namespace.yaml --server-dry-run
namespace/valid-namespace created (server dry run)

[20-08-05 23:35:14] nakamasato at Masatos-MacBook-Pro in ~/Code/MasatoNaka/kubernetes-training/open-policy-agent on postgres-operator ✘
± kubectl apply -f gatekeeper/invalid-namespace.yaml --server-dry-run
Error from server ([denied by ns-must-have-gk] you must provide labels: {"gatekeeper"}): error when creating "gatekeeper/invalid-namespace.yaml": admission webhook "validation.gatekeeper.sh" denied the request: [denied by ns-must-have-gk] you must provide labels: {"gatekeeper"}
```



# FAQ

1. Run on GKE

https://github.com/open-policy-agent/gatekeeper#running-on-private-gke-cluster-nodes
