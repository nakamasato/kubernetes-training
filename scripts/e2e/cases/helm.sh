# shellcheck disable=SC2154 # context and KUBECONFIG are provided by run.sh
command -v helm >/dev/null
helm version --short
h() { helm --kubeconfig "$KUBECONFIG" --kube-context "$context" "$@"; }
chart=contents/helm/hello-world/helloworld-chart
h lint "$chart"
h install hello "$chart" --namespace helm-e2e --create-namespace --wait --timeout 300s
h test hello --namespace helm-e2e --timeout 120s
# Exercise optional resources too: deprecated APIs can hide behind false defaults.
h upgrade hello "$chart" --namespace helm-e2e --wait --timeout 300s \
  --set ingress.enabled=true --set autoscaling.enabled=true
k -n helm-e2e get ingress/hello-helloworld-chart hpa/hello-helloworld-chart
wait_http helm-e2e hello-helloworld-chart:80 /
h uninstall hello --namespace helm-e2e --wait --timeout 120s
