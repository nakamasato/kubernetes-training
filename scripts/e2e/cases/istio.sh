# shellcheck shell=bash
# shellcheck disable=SC2154 # run.sh provides the private cluster and directories.
version=$(python3 - <<'PY'
import re
from pathlib import Path
matches = set(re.findall(r'^    ISTIO_VERSION=(\d+\.\d+\.\d+)$', Path('contents/istio/README.md').read_text(), re.M))
if len(matches) != 1:
    raise SystemExit('Expected one pinned ISTIO_VERSION in the Istio README')
print(matches.pop())
PY
)
case "$(uname -s)" in
  Darwin) platform=osx ;;
  Linux) platform=linux ;;
  *) echo 'Unsupported Istio platform' >&2; exit 1 ;;
esac
case "$(uname -m)" in
  arm64|aarch64) arch=arm64 ;;
  x86_64) arch=amd64 ;;
  *) echo 'Unsupported Istio architecture' >&2; exit 1 ;;
esac
archive="istio-$version-$platform-$arch.tar.gz"
release="https://github.com/istio/istio/releases/download/$version"
curl -fsSL --retry 3 --max-time 180 "$release/$archive" -o "$run_dir/$archive"
curl -fsSL --retry 3 --max-time 60 "$release/$archive.sha256" -o "$run_dir/$archive.sha256"
python3 - "$run_dir/$archive" <<'PY'
import hashlib
from pathlib import Path
import sys
archive = Path(sys.argv[1])
expected = Path(str(archive) + '.sha256').read_text().split()[0]
assert hashlib.sha256(archive.read_bytes()).hexdigest() == expected, 'Istio archive checksum mismatch'
PY
tar -xzf "$run_dir/$archive" -C "$run_dir"
istio_dir="$run_dir/istio-$version"
i() { "$istio_dir/bin/istioctl" --kubeconfig "$KUBECONFIG" --context "$context" "$@"; }
i version --remote=false
i install --set profile=demo --skip-confirmation --readiness-timeout=300s
for deployment in istiod istio-ingressgateway istio-egressgateway; do
  k -n istio-system rollout status "deployment/$deployment" --timeout=300s
done
k label namespace default istio-injection=enabled
samples="$istio_dir/samples/bookinfo"
k apply -f "$samples/platform/kube/bookinfo.yaml"
for deployment in details-v1 ratings-v1 reviews-v1 reviews-v2 reviews-v3 productpage-v1; do
  k rollout status "deployment/$deployment" --timeout=300s
done
k get pods -n default -o json >"$artifacts/bookinfo-pods.json"
python3 - "$artifacts/bookinfo-pods.json" <<'PY'
import json
import sys
pods = json.load(open(sys.argv[1]))['items']
assert len(pods) == 6, f'Expected six Bookinfo pods, got {len(pods)}'
for pod in pods:
    # Istio may inject a native sidecar as a restartable init container.
    containers = pod['spec']['containers'] + pod['spec'].get('initContainers', [])
    assert any(c['name'] == 'istio-proxy' for c in containers), pod['metadata']['name']
PY
k apply -f "$samples/networking/bookinfo-gateway.yaml"
k apply -f contents/istio/destination-rule-all.yaml
k apply -f "$samples/networking/virtual-service-all-v1.yaml"
i analyze -n default
wait_http istio-system istio-ingressgateway:80 /productpage
grep -q '<title>Simple Bookstore App</title>' "$artifacts/response.json"
# Requests originate inside the mesh so the caller's sidecar applies the routes.
check_reviews() {
  local user=$1 expected=$2 _attempt consecutive=0
  for _attempt in {1..30}; do
    if k exec deployment/ratings-v1 -c ratings -- \
      curl -fsS --max-time 10 -H "end-user: $user" http://reviews:9080/reviews/0 \
      >"$artifacts/reviews-$user.json" && \
      python3 - "$artifacts/reviews-$user.json" "$expected" <<'PY'
import json
import sys
response = json.load(open(sys.argv[1]))
reviews = response['reviews']
assert len(reviews) == 2, response
if sys.argv[2] == 'v1':
    assert all('rating' not in r for r in reviews), response
else:
    assert all(r['rating']['color'] == 'black' and r['rating']['stars'] > 0 for r in reviews), response
PY
    then
      consecutive=$((consecutive + 1))
      if (( consecutive == 5 )); then
        return 0
      fi
    else
      consecutive=0
    fi
    sleep 2
  done
  return 1
}
check_reviews anonymous v1
k apply -f "$samples/networking/virtual-service-reviews-test-v2.yaml"
check_reviews jason v2
check_reviews anonymous v1
