k create namespace argocd
k apply --server-side -k contents/argocd/setup
k -n argocd rollout status deployment --timeout=300s
k -n argocd rollout status statefulset/argocd-application-controller --timeout=300s
wait_http argocd https:argocd-server:443 /healthz
