# shellcheck shell=bash
# shellcheck disable=SC2016 # Variables expand inside the Elasticsearch container.
# shellcheck disable=SC2154 # artifacts is provided by run.sh.
k apply --server-side -k contents/eck/operator
k -n elastic-system rollout status statefulset/elastic-operator --timeout=300s
k create namespace eck
k apply -f contents/eck/elasticsearch.yaml
k apply -f contents/eck/kibana.yaml
k -n eck wait --for=jsonpath='{.status.availableNodes}'=1 elasticsearch/quickstart --timeout=600s
k -n eck wait --for=jsonpath='{.status.availableNodes}'=1 kibana/quickstart --timeout=600s
# Pass generated credentials/certificates over stdin, never print them in logs.
password=$(k -n eck get secret quickstart-es-elastic-user -o jsonpath='{.data.elastic}' | python3 -c 'import base64,sys; print(base64.b64decode(sys.stdin.read()).decode())')
eck_request() {
  local service=$1 port=$2 secret=$3 path=$4
  shift 4
  {
    printf '%s\n' "$password"
    k -n eck get secret "$secret" -o jsonpath='{.data.tls\.crt}' | python3 -c 'import base64,sys; sys.stdout.buffer.write(base64.b64decode(sys.stdin.read()))'
  } | k -n eck exec -i quickstart-es-default-0 -c elasticsearch -- bash -ec '
    read -r password
    ca=$(mktemp)
    trap '\''rm -f "$ca"'\'' EXIT
    cat > "$ca"
    curl --fail --silent --show-error --max-time 30 --cacert "$ca" -u "elastic:$password" "$@"
  ' bash "https://$service.eck.svc:$port$path" "$@"
}
eck_request quickstart-es-http 9200 quickstart-es-http-certs-public '/training-e2e/_doc/1?refresh=true' \
  -X PUT -H 'Content-Type: application/json' -d '{"message":"service-write-read-ok"}' >"$artifacts/write.json"
eck_request quickstart-es-http 9200 quickstart-es-http-certs-public /training-e2e/_doc/1 >"$artifacts/read.json"
python3 -c 'import json,sys; r=json.load(sys.stdin); assert r["found"] and r["_source"]["message"] == "service-write-read-ok"' <"$artifacts/read.json"
eck_request quickstart-kb-http 5601 quickstart-kb-http-certs-public /api/status >"$artifacts/kibana.json"
python3 -c 'import json,sys; assert json.load(sys.stdin)["status"]["overall"]["level"] == "available"' <"$artifacts/kibana.json"
