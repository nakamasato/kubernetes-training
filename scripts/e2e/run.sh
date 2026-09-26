#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/../.."
if [[ $# != 1 ]] || [[ ! -f "scripts/e2e/cases/$1.sh" ]] || [[ "$1" == *[!a-z0-9-]* ]]; then
  echo "Usage: $0 <target>; available targets:" >&2
  python3 -c 'import json; print("\n".join(json.load(open("scripts/e2e/targets.json"))))' >&2
  exit 2
fi
for tool in kind kubectl docker python3 curl; do
  command -v "$tool" >/dev/null || { echo "Missing tool: $tool" >&2; exit 1; }
done
target=$1
# A fresh directory and random name prevent touching an existing cluster/config.
run_dir=$(mktemp -d "${TMPDIR:-/tmp}/training-e2e.XXXXXXXX")
export E2E_RUN_DIR=$run_dir
E2E_BACKGROUND_PIDS=()
provided_kubeconfig=${KUBECONFIG:-}
if [[ -n "$provided_kubeconfig" ]]; then
  context=${KUBE_CONTEXT:-$(kubectl --kubeconfig "$provided_kubeconfig" config current-context)}
  export KUBECONFIG="$provided_kubeconfig"
  external_cluster=true
else
  cluster="training-e2e-$(basename "$run_dir" | tr '[:upper:]' '[:lower:]' | tr -cd 'a-z0-9' | tail -c 8)"
  export KUBECONFIG="$run_dir/kubeconfig"
  context="kind-$cluster"
  external_cluster=false
fi
artifacts="${E2E_ARTIFACTS:-$PWD/.e2e-artifacts/$target}"
mkdir -p "$artifacts"
artifacts=$(cd "$artifacts" && pwd)
export KIND_EXPERIMENTAL_PROVIDER=docker
k() { kubectl --kubeconfig "$KUBECONFIG" --context "$context" "$@"; }
cleanup() {
  result=$?
  trap - EXIT
  set +e
  for pid in "${E2E_BACKGROUND_PIDS[@]:-}"; do
    kill "$pid" 2>/dev/null
  done
  for pid in "${E2E_BACKGROUND_PIDS[@]:-}"; do
    wait "$pid" 2>/dev/null
  done
  if (( result != 0 )); then
    k get pods -A -o wide >"$artifacts/pods.txt" 2>&1
    k get events -A --sort-by=.metadata.creationTimestamp >"$artifacts/events.txt" 2>&1
    k describe pods -A >"$artifacts/describe.txt" 2>&1
    if [[ "$external_cluster" == false ]]; then
      kind export logs "$artifacts/kind" --name "$cluster" >"$artifacts/export.txt" 2>&1
    fi
    while read -r ns pod; do
      k logs -n "$ns" "$pod" --all-containers --tail=200 >"$artifacts/$ns-$pod.log" 2>&1
    done < <(k get pods -A -o jsonpath='{range .items[*]}{.metadata.namespace}{" "}{.metadata.name}{"\n"}{end}')
  fi
  if [[ "$external_cluster" == false ]]; then
    kind delete cluster --name "$cluster" --kubeconfig "$KUBECONFIG" || result=1
  fi
  rm -rf "$run_dir"
  exit "$result"
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
if [[ "$external_cluster" == false ]]; then
  kind create cluster --name "$cluster" --kubeconfig "$KUBECONFIG" \
    --image "${KIND_NODE_IMAGE:-kindest/node:v1.36.1}" --wait 180s
fi
# kind waits for the node, but Service DNS may still be starting.
k -n kube-system rollout status deployment/coredns --timeout=180s
# API proxy exercises the Service and application without needing a curl image.
http_get() {
  local namespace=$1 service=$2 path=$3
  k --request-timeout=15s get --raw "/api/v1/namespaces/$namespace/services/$service/proxy$path"
}
wait_http() {
  local _attempt
  for _attempt in {1..60}; do
    if http_get "$@" >"$artifacts/response.json"; then
      cat "$artifacts/response.json"
      return 0
    fi
    sleep 5
  done
  return 1
}
# Cases can only reach this fresh cluster via k(), which always passes --context.
# shellcheck disable=SC1090
source "scripts/e2e/cases/$target.sh"
echo "PASS: $target"
