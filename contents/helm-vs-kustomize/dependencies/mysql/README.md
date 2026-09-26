# MySQL dependency

This fresh-install sample uses MySQL 8.4.11 LTS (previously 5.6).
See the [official image](https://hub.docker.com/_/mysql) for supported tags.
Run from the repository root:

```sh
kubectl create namespace database
kubectl apply -k contents/helm-vs-kustomize/dependencies/mysql
kubectl -n database rollout status deployment/mysql --timeout=300s
bash scripts/e2e/run.sh mysql
```

The startup/readiness probes authenticate as the sample user and query the sample
database. E2E also verifies SQL writes and reads through `mysql.database.svc`.
Credentials are for local training only and data is ephemeral. This is not an
in-place upgrade procedure for an existing 5.6 database. The Flask examples are
outside this test; older clients may need updates for MySQL 8.4 authentication.
