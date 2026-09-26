# shellcheck shell=bash
# shellcheck disable=SC2154 # artifacts is provided by run.sh.
# shellcheck disable=SC2016 # Variables expand inside the PostgreSQL container.
k apply --server-side -k contents/grafana-operator/operator
k -n grafana rollout status deployment/grafana-operator-controller-manager --timeout=300s
k apply -k contents/grafana-operator/ha -n grafana
k -n grafana rollout status deployment/postgres --timeout=300s
# The operator creates the Deployment asynchronously.
k -n grafana wait --for=create deployment/example-grafana-deployment --timeout=180s
k -n grafana rollout status deployment/example-grafana-deployment --timeout=360s
test "$(k -n grafana get deployment/example-grafana-deployment -o jsonpath='{.status.readyReplicas}')" = 2
wait_http grafana example-grafana-service:3000 /api/health
python3 -c 'import json,sys; assert json.load(sys.stdin)["database"] == "ok"' <"$artifacts/response.json"
# Grafana must actually initialize its schema in PostgreSQL, not fall back to SQLite.
result=$(k -n grafana exec deployment/postgres -- sh -ec '
  export PGPASSWORD="$POSTGRES_PASSWORD" PGCONNECT_TIMEOUT=10
  psql -h postgres.grafana.svc -U "$POSTGRES_USER" -d "$POSTGRES_DB" -v ON_ERROR_STOP=1 -Atc "SELECT count(*) > 0 FROM migration_log;"
')
test "$result" = t
