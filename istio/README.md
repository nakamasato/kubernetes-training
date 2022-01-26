# [Istio](https://istio.io/)


*A service mesh is a dedicated infrastructure layer that you can add to your applications. It allows you to transparently add capabilities like observability, traffic management, and security, without adding them to your own code. The term “service mesh” describes both the type of software you use to implement this pattern, and the security or network domain that is created when you use that software.*

## [Getting Started](https://istio.io/latest/docs/setup/getting-started/)

### Installation

```
curl -L https://istio.io/downloadIstio | sh -
cd istio-1.12.2
export PATH=$PWD/bin:$PATH
istioctl install --set profile=demo -y
```

<details><summary>Result</summary>

```
✔ Istio core installed
✔ Istiod installed
✔ Egress gateways installed
✔ Ingress gateways installed
✔ Installation complete                                                                                              Making this installation the default for injection and validation.

Thank you for installing Istio 1.12.  Please take a few minutes to tell us about your install/upgrade experience!  https://forms.gle/FegQbc9UvePd4Z9z7
```

</details>

Add a namespace label to instruct Istio to automatically inject Envoy sidecar proxies when you deploy your application later:

```
kubectl label namespace default istio-injection=enabled
```

### Deploy the Sample app

Deploy sample app.

```
kubectl apply -f samples/bookinfo/platform/kube/bookinfo.yaml
```

Envoy sider is added to all pods.

```
kubectl get po
NAME                              READY   STATUS    RESTARTS   AGE
details-v1-79f774bdb9-ctf75       2/2     Running   0          28s
productpage-v1-6b746f74dc-7zgpg   2/2     Running   0          28s
ratings-v1-b6994bb9-rw74b         2/2     Running   0          28s
reviews-v1-545db77b95-t6gkl       2/2     Running   0          28s
reviews-v2-7bf8c9648f-n9tmq       2/2     Running   0          28s
reviews-v3-84779c7bbc-tmzlr       2/2     Running   0          28s
```

Verify app is running.

```
kubectl exec "$(kubectl get pod -l app=ratings -o jsonpath='{.items[0].metadata.name}')" -c ratings -- curl -sS productpage:9080/productpage | grep -o "<title>.*</title>"

<title>Simple Bookstore App</title>
```

### Open the app to outside traffic

Istio Gateway

```
kubectl apply -f samples/bookinfo/networking/bookinfo-gateway.yaml
```

Check
```
istioctl analyze
✔ No validation issues found when analyzing namespace: default.
```

Not confirmed on a kind cluster

### Dashboard

Install [kiali](https://istio.io/latest/docs/ops/integrations/kiali/) dashboard

kubectl apply -f samples/addons

Open dashboard

```
istioctl dashboard kiali
```

kubectl port-forward svc/productpage 9080:9080

![](docs/sample-app.png)

The traffic is visualized in the graph.

![](docs/kiali.png)

### Cleanup

```
kubectl delete -f samples/addons
istioctl manifest generate --set profile=demo | kubectl delete --ignore-not-found=true -f -
istioctl tag remove default
```
