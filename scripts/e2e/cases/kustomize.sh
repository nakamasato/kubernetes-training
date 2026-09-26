# shellcheck disable=SC2154 # artifacts is provided by run.sh
command -v kustomize >/dev/null
kustomize version
for environment in dev prod; do
  k create namespace "kustomize-$environment"
  kustomize build "contents/kustomize/example/overlays/$environment" | k apply -f -
  k -n "kustomize-$environment" rollout status deployment/kustomize-demo --timeout=300s
  wait_http "kustomize-$environment" kustomize-demo:80 /
  grep -qx "environment=$environment" "$artifacts/response.json"
  expected=1
  if [[ "$environment" == prod ]]; then expected=2; fi
  test "$(k -n "kustomize-$environment" get deployment kustomize-demo -o jsonpath='{.status.readyReplicas}')" = "$expected"
done
