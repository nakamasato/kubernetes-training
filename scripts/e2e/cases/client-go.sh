# shellcheck disable=SC2154 # KUBECONFIG and artifacts are provided by run.sh
go test ./contents/kubernetes-operator/e2e -run '^TestClientGo' -count=1 -v
