# shellcheck disable=SC2154 # context, KUBECONFIG, artifacts, and k() are provided by run.sh
go test ./contents/kubernetes-operator/e2e -run '^TestControllerRuntime' -count=1 -v
