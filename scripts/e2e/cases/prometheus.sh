# shellcheck disable=SC2154 # artifacts is provided by run.sh
k create namespace monitoring
k apply -k contents/prometheus
k -n monitoring rollout status statefulset/prometheus --timeout=300s
wait_http monitoring prometheus:9090 /-/ready
for _attempt in {1..30}; do
  http_get monitoring prometheus:9090 '/api/v1/query?query=up' >"$artifacts/response.json"
  if python3 -c 'import json,sys; r=json.load(sys.stdin); sys.exit(not any(x["metric"].get("job")=="prometheus" and x["value"][1]=="1" for x in r["data"]["result"]))' <"$artifacts/response.json"; then break; fi
  sleep 5
done
python3 -c 'import json,sys; r=json.load(sys.stdin); assert any(x["metric"].get("job")=="prometheus" and x["value"][1]=="1" for x in r["data"]["result"])' <"$artifacts/response.json"
