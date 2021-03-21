# HPA example


## Preparee RabbitMQ

1. RabbitMQ Operator

https://www.rabbitmq.com/kubernetes/operator/quickstart-operator.html

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
