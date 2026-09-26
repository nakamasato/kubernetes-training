# Grafana Operator

This example runs two Grafana 13.2.2 replicas against a shared PostgreSQL
[18.6](https://www.postgresql.org/docs/release/18.6/) database (previously 17),
managed by Grafana Operator v5.25.0.

## Install

Run from the repository root. The pinned release bundle creates the `grafana`
namespace and installs the operator there:

```sh
kubectl apply --server-side -k contents/grafana-operator/operator
kubectl -n grafana rollout status deployment/grafana-operator-controller-manager --timeout=300s
kubectl apply -k contents/grafana-operator/ha -n grafana
kubectl -n grafana rollout status deployment/postgres --timeout=300s
kubectl -n grafana wait --for=create deployment/example-grafana-deployment --timeout=180s
kubectl -n grafana rollout status deployment/example-grafana-deployment --timeout=360s
```

The Grafana custom resource uses the v5 API: configuration values are strings,
and replica overrides live under `spec.deployment.spec`. PostgreSQL 18 uses
`/var/lib/postgresql/18/docker` for data with a volume at `/var/lib/postgresql`,
following the [official image](https://hub.docker.com/_/postgres).

Database credentials, anonymous viewer access and `emptyDir` storage are for
local training. Two Grafana replicas demonstrate a shared database; PostgreSQL
itself is a single ephemeral instance. This is not a data migration procedure
for an existing PostgreSQL 17 database.

## Access

```sh
kubectl -n grafana port-forward service/example-grafana-service 3000
```

Open http://localhost:3000 or check http://localhost:3000/api/health. Database
health should be `ok`.

## E2E

```sh
bash scripts/e2e/run.sh grafana-operator
```

The disposable kind test checks both Grafana replicas, Service HTTP/database
health, and Grafana's migration table in PostgreSQL through the database Service.
It cleans up the cluster on completion or failure.

## Optional integrations

The older `datasource-prometheus.yaml`, `dashboard.yaml` and `prometheus/`
examples are outside this E2E target and still need migration to the v5
integration APIs. Use the operator's
[versioned examples](https://github.com/grafana/grafana-operator/tree/v5.25.0/examples)
for current dashboard and datasource configuration.

## Cleanup

```sh
kubectl delete -k contents/grafana-operator/ha -n grafana
kubectl delete -k contents/grafana-operator/operator
```
