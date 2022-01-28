# [Istio](https://istio.io/)


*A service mesh is a dedicated infrastructure layer that you can add to your applications. It allows you to transparently add capabilities like observability, traffic management, and security, without adding them to your own code. The term “service mesh” describes both the type of software you use to implement this pattern, and the security or network domain that is created when you use that software.*


Istio uses [Envoy](https://www.envoyproxy.io/), *AN OPEN SOURCE EDGE AND SERVICE PROXY, DESIGNED FOR CLOUD-NATIVE APPLICATIONS*, proxy as its data plane.
## [Getting Started](https://istio.io/latest/docs/setup/getting-started/)

**If you test on your local cluster, pleasee use docker-desktop (or minikube).** (Not confirmed on a kind cluster)

### [Install Istio](https://istio.io/latest/docs/setup/getting-started/#bookinfo)

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

### [Deploy the sample application](https://istio.io/latest/docs/setup/getting-started/#bookinfo)

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

### [Open the app to outside traffic](https://istio.io/latest/docs/setup/getting-started/#ip)

Istio Gateway

```
kubectl apply -f samples/bookinfo/networking/bookinfo-gateway.yaml
```

Check
```
istioctl analyze
✔ No validation issues found when analyzing namespace: default.
```

Check ingress gateway

```
kubectl get svc istio-ingressgateway -n istio-system
NAME                   TYPE           CLUSTER-IP     EXTERNAL-IP   PORT(S)                                                                      AGE
istio-ingressgateway   LoadBalancer   10.103.34.38   localhost     15021:31476/TCP,80:31411/TCP,443:32714/TCP,31400:30467/TCP,15443:30550/TCP   44m
```

Set ingress ip and ports:

```
export INGRESS_HOST=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
export INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].port}')
export SECURE_INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="https")].port}')
```

Docker for Desktop:

```
export INGRESS_HOST=127.0.0.1
```

Check

```
echo "$GATEWAY_URL"
127.0.0.1:80
```

```
echo "http://$GATEWAY_URL/productpage"
http://127.0.0.1:80/productpage
```

Open http://127.0.0.1:80/productpage on your browser:

![](docs/sample-app.png)

### [View the dashboard](https://istio.io/latest/docs/setup/getting-started/#dashboard)

Install [kiali](https://istio.io/latest/docs/ops/integrations/kiali/) dashboard

```
kubectl apply -f samples/addons
kubectl rollout status deployment/kiali -n istio-system
```

Open dashboard

```
istioctl dashboard kiali
```

The traffic is visualized in the graph.

![](docs/kiali.png)

### Cleanup

```
kubectl delete -f samples/addons
istioctl manifest generate --set profile=demo | kubectl delete --ignore-not-found=true -f -
istioctl tag remove default
```

```
kubectl delete namespace istio-system
kubectl label namespace default istio-injection-
```
