# shellcheck disable=SC2154 # artifacts is provided by run.sh
k create namespace monitoring
k apply -k contents/grafana
k -n monitoring rollout status deployment/grafana --timeout=300s
wait_http monitoring grafana:3000 /api/health
python3 -c 'import json,sys; assert json.load(sys.stdin)["database"] == "ok"' <"$artifacts/response.json"
