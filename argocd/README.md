# ArgoCD

## Install

```
kubectl apply -k argocd/setup
```

## Login

```
kubectl -n argocd port-forward service/argocd-server 8080:80
```

- user: `admin`
- password: `kubectl get po -n argocd | grep argocd-server | awk '{print $1}'`

![](img/argocd.png)

## Add ArgoCD AppProject & Application


