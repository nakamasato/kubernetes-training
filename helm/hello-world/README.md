# Helm hello world

https://medium.com/@pablorsk/kubernetes-helm-node-hello-world-c97d20437abd

1. build

```
docker build -t hello-world .
```

1. delete container
```
docker rm $(docker ps | grep hello-world | awk '{print $1}') -f
```

1. create helm

```
helm create helloworld-chart
```

1. change `values.yaml`

1. package

```
helm package helloworld-chart --debug
```

1. apply helm `hello-world`

```
helm install helloworld-chart-0.1.0.tgz --generate-name
NAME: helloworld-chart-0-1589247182
LAST DEPLOYED: Tue May 12 10:33:03 2020
NAMESPACE: default
STATUS: deployed
REVISION: 1
NOTES:
1. Get the application URL by running these commands:
  export POD_NAME=$(kubectl get pods --namespace default -l "app.kubernetes.io/name=helloworld-chart,app.kubernetes.io/instance=helloworld-chart-0-1589247182" -o jsonpath="{.items[0].metadata.name}")
  echo "Visit http://127.0.0.1:8080 to use your application"
  kubectl --namespace default port-forward $POD_NAME 8080:80
```

1. check

```
kubectl get deployment -n default
NAME                            READY   UP-TO-DATE   AVAILABLE   AGE
helloworld-chart-0-1589247182   0/1     1            0           16s
```

1. delete

```
helm uninstall helloworld-chart-0-1589247182
release "helloworld-chart-0-1589247182" uninstalled
```

# Helm basic commands

```
helm list
NAME                            NAMESPACE       REVISION        UPDATED                                 STATUS          CHART                   APP VERSION
helloworld-chart-0-1589247182   default         1               2020-05-12 10:33:03.884258 +0900 JST    deployed        helloworld-chart-0.1.0  1.16.0
```
