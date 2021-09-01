# Grafana

https://github.com/grafana/grafana
## Version

Latest

## Install

```
kubectl create ns monitoring
kubectl -k .
```

## Dashboards

Pre-configured dashboards are in `config/dashboards/dashboard.yaml`

- RabbitMQ: `config/dashboards/rabbitmq.json`

## Data sources

Configured data sources are in `config/datasources/prometheus.yaml`

- Prometheus: `http://prometheus.monitoring.svc:9090`

## Related

You may also like [Loki](../loki)
