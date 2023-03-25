# Ingress Nginx Controller

https://kubernetes.github.io/ingress-nginx/deploy

## Prerequisite

- Kubernetes:
    ```
    kubectl version --short
    ```

## Get started

- `namespace`: `ingress-nginx` (will be created)

1. apply

    ```bash
    kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v0.48.1/deploy/static/provider/cloud/deploy.yaml

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
    ingress-nginx-admission-create-j2mvc        0/1     Completed   0          62s
    ingress-nginx-admission-patch-vpfdn         0/1     Completed   2          62s
    ingress-nginx-controller-68649d49b8-8xld6   1/1     Running     0          62s
    ```

1. check nginx ingress controller's version

    ```
    kubectl exec -it $(kubectl get po -n ingress-nginx | grep ingress-nginx-controller | awk '{print $1}') -n ingress-nginx -- /nginx-ingress-controller --version
    -------------------------------------------------------------------------------
    NGINX Ingress controller
      Release:       v0.48.1
      Build:         30809c066cd027079cbb32dccc8a101d6fbffdcb
      Repository:    https://github.com/kubernetes/ingress-nginx
      nginx version: nginx/1.20.1

    -------------------------------------------------------------------------------
    ```

## Practice with an app

[Kubernetes ingress nginx contoller](https://medium.com/rahasak/kubernetes-ingress-nginx-contoller-fa60b8d7e5f1)

1. Deploy `ReplicaSet`, `Service` and `Ingress`.

    ```
    kubectl create -f .
    ```

1. Check `Service`

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
    {"id":"piyasara-api1628479184739","code":201,"message":"created"}%
    ```

1. Check `Ingress`

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

## [Changelogs](https://github.com/kubernetes/ingress-nginx/blob/main/Changelog.md)

- 0.40.0: Following the Ingress [extensions/v1beta1 deprecation](https://kubernetes.io/blog/2019/07/18/api-deprecations-in-1-16/), `networking.k8s.io/v1beta1` or `networking.k8s.io/v1` (Kubernetes v1.19 or higher)
- 0.25.0: Support new `networking.k8s.io/v1beta1` package (for Kubernetes cluster > v1.14.0)
