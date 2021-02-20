# Ingress Nginx Controller

- https://kubernetes.github.io/ingress-nginx/deploy
- https://docs.nginx.com/nginx-ingress-controller

## Prerequisite

- Kubernetes: Kubernetes for Mac
    ```
    kubectl version
    Client Version: version.Info{Major:"1", Minor:"15", GitVersion:"v1.15.5", GitCommit:"20c265fef0741dd71a66480e35bd69f18351daea", GitTreeState:"clean", BuildDate:"2019-10-15T19:16:51Z", GoVersion:"go1.12.10", Compiler:"gc", Platform:"darwin/amd64"}
    Server Version: version.Info{Major:"1", Minor:"15", GitVersion:"v1.15.5", GitCommit:"20c265fef0741dd71a66480e35bd69f18351daea", GitTreeState:"clean", BuildDate:"2019-10-15T19:07:57Z", GoVersion:"go1.12.10", Compiler:"gc", Platform:"linux/amd64"}
    ```
- Docker

    ```
    Client: Docker Engine - Community
    Version:           19.03.8
    API version:       1.40
    Go version:        go1.12.17
    Git commit:        afacb8b
    Built:             Wed Mar 11 01:21:11 2020
    OS/Arch:           darwin/amd64
    Experimental:      false

    Server: Docker Engine - Community
    Engine:
      Version:          19.03.8
      API version:      1.40 (minimum version 1.12)
      Go version:       go1.12.17
      Git commit:       afacb8b
      Built:            Wed Mar 11 01:29:16 2020
      OS/Arch:          linux/amd64
      Experimental:     false
    containerd:
      Version:          v1.2.13
      GitCommit:        7ad184331fa3e55e52b890ea95e65ba581ae3429
    runc:
      Version:          1.0.0-rc10
      GitCommit:        dc9208a3303feef5b3839f4323d9beb36df0a9dd
    docker-init:
      Version:          0.18.0
      GitCommit:        fec3683
    ```

## Get started

1. apply

    ```bash
    kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-0.32.0/deploy/static/provider/cloud/deploy.yaml
    namespace/ingress-nginx created
    serviceaccount/ingress-nginx created
    configmap/ingress-nginx-controller created
    clusterrole.rbac.authorization.k8s.io/ingress-nginx created
    clusterrolebinding.rbac.authorization.k8s.io/ingress-nginx created
    role.rbac.authorization.k8s.io/ingress-nginx created
    rolebinding.rbac.authorization.k8s.io/ingress-nginx created
    service/ingress-nginx-controller-admission created
    service/ingress-nginx-controller created
    deployment.apps/ingress-nginx-controller created
    validatingwebhookconfiguration.admissionregistration.k8s.io/ingress-nginx-admission created
    clusterrole.rbac.authorization.k8s.io/ingress-nginx-admission created
    clusterrolebinding.rbac.authorization.k8s.io/ingress-nginx-admission created
    job.batch/ingress-nginx-admission-create created
    job.batch/ingress-nginx-admission-patch created
    role.rbac.authorization.k8s.io/ingress-nginx-admission created
    rolebinding.rbac.authorization.k8s.io/ingress-nginx-admission created
    serviceaccount/ingress-nginx-admission created
    ```

1. check the pods

    ```
    kubectl get pod -n ingress-nginx
    NAME                                        READY   STATUS      RESTARTS   AGE
    ingress-nginx-admission-create-l7d5d        0/1     Completed   0          3m31s
    ingress-nginx-admission-patch-zwwmb         0/1     Completed   1          3m31s
    ingress-nginx-controller-564768b55d-v52b7   1/1     Running     0          3m41s
    ```

1. check nginx ingress controller's version

    ```
    kubectl exec -it ingress-nginx-controller-564768b55d-v52b7 -n ingress-nginx -- /nginx-ingress-controller --version
    -------------------------------------------------------------------------------
    NGINX Ingress controller
      Release:       0.32.0
      Build:         git-446845114
      Repository:    https://github.com/kubernetes/ingress-nginx
      nginx version: nginx/1.17.10

    -------------------------------------------------------------------------------
    ```

## Practice with an app

1. ReplicationController

```
kubectl create -f piyasara-pod.yaml
```

1. Service

```
piyasara-service.yaml
```

1. Ingress

```
kubectl apply -f piyasara-ingress.yaml
```

1. Check if service is working expected

```
kubectl port-forward svc/piyasara-service  8088:8761
Forwarding from 127.0.0.1:8088 -> 8761
Forwarding from [::1]:8088 -> 8761
Handling connection for 8088
```

```
curl -XPOST "http://localhost:8088/flights" -d '
{
  "Uid": "testuid",
  "FlightNo": "UL1234",
  "FlightStatus": "Done",
  "FlightFrom": "SG",
  "FlightTo": "Tokyo"
}
'
{"id":"piyasara-api1588659552418","code":201,"message":"created"}%
```

1. Check if Ingress is working

```
curl -XPOST "http://localhost/flights" -d '
{
  "Uid": "testuid",
  "FlightNo": "UL1234",
  "FlightStatus": "Done",
  "FlightFrom": "SG",
  "FlightTo": "Tokyo"
}
'
{"id":"piyasara-api1588680669155","code":201,"message":"created"}
```
