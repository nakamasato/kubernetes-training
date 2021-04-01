# Horizontal Pod Autoscaler

## Overview

### What HPA does

Automatically scales the number of Pods in a replication controller, deployment, replica set or stateful set based on observed CPU utilization or with custom metrics.

### API Object

- API Group: `autoscaling`
- CPU autoscaling: `autoscaling/v1`
- memory & custom metrics: `autoscaling/v2beta2`

### Kubectl

```
kubectl autoscale rs foo --min=2 --max=5 --cpu-percent=80
```

## Examples

### CPU

1. Apply an apache application

    ```
    kubectl apply -f https://k8s.io/examples/application/php-apache.yaml
    ```

1. Set autoscale by kubectl

    ```
    kubectl autoscale deployment php-apache --cpu-percent=50 --min=1 --max=10
    ```

    ```
    kubectl get hpa
    NAME         REFERENCE               TARGETS         MINPODS   MAXPODS   REPLICAS   AGE
    php-apache   Deployment/php-apache   <unknown>/50%   1         10        0          9s
    ```

    <details>

    ```
    apiVersion: autoscaling/v1
    kind: HorizontalPodAutoscaler
    metadata:
      name: php-apache
      namespace: default
    spec:
      maxReplicas: 10
      minReplicas: 1
      scaleTargetRef:
        apiVersion: apps/v1
        kind: Deployment
        name: php-apache
      targetCPUUtilizationPercentage: 50
    ```

    </details>

1. Increase the load

    ```
    kubectl run -i --tty load-generator --rm --image=busybox --restart=Never -- /bin/sh -c "while sleep 0.01; do wget -q -O- http://php-apache; done"
    ```

    ```
    kubectl get hpa

    NAME         REFERENCE               TARGETS   MINPODS   MAXPODS   REPLICAS   AGE
    php-apache   Deployment/php-apache   76%/50%   1         10        7          4m10s
    ```
