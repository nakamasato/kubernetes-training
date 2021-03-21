# HPA example


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
