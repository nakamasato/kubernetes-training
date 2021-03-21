# HPA example

## Overview

- RabbitMQ Producer (Java app)
- RabbitMQ
- RabbitMQ Consumer (Java app)
- Prometheus -> localhost:30900
- Grafana -> localhost:32000

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

`localhost:61830/metrics`

> As of 3.8.0, RabbitMQ ships with built-in Prometheus & Grafana support.

> Support for Prometheus metric collector ships in the rabbitmq_prometheus plugin. The plugin exposes all RabbitMQ metrics on a dedicated TCP port, in Prometheus text format.

## Deploy producer

1. Create `rabbitmq-producer` `CronJob` (run hourly)

```
kubectl apply -f rabbitmq-producer
```

If you want to run a job manually, you can run the following command after creating `CronJob`

```
kubectl create job --from=cronjob/rabbitmq-producer rabbitmq-producer-$(date '+%s')
```

## Deploy consumer

1. Create `rabbitmq-consumer` `Deployment`

```
kubectl apply -f rabbitmq-consumer
```

## Deploy Prometheus

https://github.com/prometheus-operator/prometheus-operator#quickstart

1. Create prometheus operator

```
kubectl apply -f https://raw.githubusercontent.com/prometheus-operator/prometheus-operator/master/bundle.yaml
```

Monitoring RabbitMQ

https://www.rabbitmq.com/kubernetes/operator/operator-monitoring.html


Issues
- [ ] Prometheus cannot scrape metrics from RabbitMQ

## Deploy Grafana

https://devopscube.com/setup-grafana-kubernetes/


## Clean up

```
for component in prometheus grafana rabbitmq rabbitmq-consumer rabbitmq-producer; do
    kubectl delete -f $component
done
kubectl delete -f "https://github.com/rabbitmq/cluster-operator/releases/latest/download/cluster-operator.yml"
kubectl delete -f https://raw.githubusercontent.com/prometheus-operator/prometheus-operator/master/bundle.yaml
```
