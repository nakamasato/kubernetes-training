# HPA example

## Overview

- RabbitMQ Producer (Java app)
- RabbitMQ
- RabbitMQ Consumer (Java app)
- Prometheus -> http://localhost:30900
- Grafana -> http://localhost:32000

![](diagram.drawio.svg)


## Deploy RabbitMQ with operator

https://www.rabbitmq.com/kubernetes/operator/quickstart-operator.html

1. RabbitMQ Operator

    ```
    kubectl apply -f "https://github.com/rabbitmq/cluster-operator/releases/latest/download/cluster-operator.yml"
    ```

1. Create a RabbitMQ cluster

    ```
    kubectl apply -f rabbitmq
    ```

1. Get username and password

    ```
    username="$(kubectl get secret rabbitmq-default-user -o jsonpath='{.data.username}' | base64 --decode)"
    echo "username: $username"
    password="$(kubectl get secret rabbitmq-default-user -o jsonpath='{.data.password}' | base64 --decode)"
    echo "password: $password"
    ```

1. port-forward

    ```
    kubectl port-forward "service/rabbitmq" 15672
    ```

    Open: http://localhost:15672/ and use the username and password got in the previous step.

Metrics:

> As of 3.8.0, RabbitMQ ships with built-in Prometheus & Grafana support.
> Support for Prometheus metric collector ships in the rabbitmq_prometheus plugin. The plugin exposes all RabbitMQ metrics on a dedicated TCP port, in Prometheus text format.

Check if `rabbitmq_prometheus` plugin is enabled.

```
kubectl exec -it rabbitmq-server-0 -- rabbitmq-plugins list | grep prometheus
[E*] rabbitmq_prometheus               3.8.12
```

## Deploy producer

Create `rabbitmq-producer` `CronJob` (run hourly)

```
kubectl apply -f rabbitmq-producer
```

If you want to run a job manually, you can run the following command after creating `CronJob`

```
kubectl create job --from=cronjob/rabbitmq-producer rabbitmq-producer-$(date '+%s')
```

## Deploy consumer

Create `rabbitmq-consumer` `Deployment`

```
kubectl apply -f rabbitmq-consumer
```

## Deploy Prometheus

References:
- https://github.com/prometheus-operator/prometheus-operator#quickstart
- https://github.com/prometheus-operator/prometheus-operator/blob/master/Documentation/user-guides/getting-started.md

Steps:

1. Create prometheus operator

    ```
    kubectl apply -f https://raw.githubusercontent.com/prometheus-operator/prometheus-operator/master/bundle.yaml
    ```

1. Prometheus

    ```
    kubectl apply -f ../../../prometheus-operator
    ```

1. Check UI at http://localhost:30900

    You can check [targets](http://localhost:30900/targets)

    ![](prometheus-target.png)


Monitoring RabbitMQ: https://www.rabbitmq.com/kubernetes/operator/operator-monitoring.html

We cannot use `ServiceMonitor` for RabbitMQ as RabbitMQ service doesn't have prometheus port (15692). We need to use `PodMonitor` as is recommended in the documentation.


## Deploy Grafana

https://devopscube.com/setup-grafana-kubernetes/

```
kubectl apply -f grafana
```

log in to http://localhost:32000 with `admin` for both username and password

import dashboard https://grafana.com/grafana/dashboards/10991

![](grafana-dashboard-for-rabbitmq.png)

## HPA with custom metrics

TBD
## Clean up

```
for component in grafana rabbitmq rabbitmq-consumer rabbitmq-producer; do
    kubectl delete -f $component
done
kubectl delete -f "https://github.com/rabbitmq/cluster-operator/releases/latest/download/cluster-operator.yml"
kubectl delete -f https://raw.githubusercontent.com/prometheus-operator/prometheus-operator/master/bundle.yaml
```

## References
- [Prometheus ServiceMonitor vs PodMonitor](https://github.com/prometheus-operator/prometheus-operator/issues/3119)
- https://qiita.com/Kameneko/items/071c2a064775badd939e
    > ただし、1点注意が必要で、これはPodのラベルではなく、Service…更に正しく言えばEndpointsのラベルを指定する必要があります。
- https://grafana.com/docs/grafana-cloud/quickstart/prometheus_operator/
- [Troubleshooting ServiceMonitor Changes](https://github.com/prometheus-operator/prometheus-operator/blob/master/Documentation/troubleshooting.md)
