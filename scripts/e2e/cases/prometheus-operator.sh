# shellcheck disable=SC2154 # artifacts is provided by run.sh
k create namespace monitoring
k apply --server-side -k contents/prometheus-operator/operator
k wait --for=condition=Established crd/prometheuses.monitoring.coreos.com crd/servicemonitors.monitoring.coreos.com --timeout=120s
k rollout status deployment/prometheus-operator --timeout=300s
k apply --server-side -k contents/prometheus-operator
# The operator creates the StatefulSet asynchronously.
for _attempt in {1..60}; do
  if k -n monitoring get statefulset prometheus-prometheus >/dev/null 2>&1; then break; fi
  sleep 5
done
k -n monitoring rollout status statefulset/prometheus-prometheus --timeout=300s
wait_http monitoring prometheus:9090 /-/ready
# Check discovery and scraping, not just that the process is running.
for _attempt in {1..60}; do
  http_get monitoring prometheus:9090 '/api/v1/query?query=up' >"$artifacts/response.json"
  if python3 -c 'import json,sys; r=json.load(sys.stdin); assert any(x["metric"].get("service")=="prometheus" and x["value"][1]=="1" for x in r["data"]["result"])' <"$artifacts/response.json"; then break; fi
  sleep 5
done
python3 -c 'import json,sys; r=json.load(sys.stdin); assert any(x["metric"].get("service")=="prometheus" and x["value"][1]=="1" for x in r["data"]["result"])' <"$artifacts/response.json"
