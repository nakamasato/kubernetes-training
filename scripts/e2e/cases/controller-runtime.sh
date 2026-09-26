# shellcheck disable=SC2154 # context, KUBECONFIG, artifacts, and k() are provided by run.sh
command -v go >/dev/null || { echo "Missing tool: go" >&2; exit 1; }

go build -o "$E2E_RUN_DIR/replicaset-controller" ./contents/kubernetes-operator/controller-runtime/example-controller
"$E2E_RUN_DIR/replicaset-controller" >"$artifacts/controller.log" 2>&1 &
E2E_BACKGROUND_PIDS+=("$!")

k create namespace operator-e2e
k -n operator-e2e apply -f - <<'EOF'
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: decoy
spec:
  replicas: 0
  selector:
    matchLabels:
      app: decoy
  template:
    metadata:
      labels:
        app: decoy
    spec:
      containers:
        - name: pause
          image: registry.k8s.io/pause:3.10
EOF

decoy_uid=$(k -n operator-e2e get rs decoy -o jsonpath='{.metadata.uid}')
k -n operator-e2e apply -f - <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: matching-but-not-owned
  labels:
    app: operator-e2e
    tier: backend
  ownerReferences:
    - apiVersion: apps/v1
      kind: ReplicaSet
      name: decoy
      uid: $decoy_uid
      controller: true
      blockOwnerDeletion: true
spec:
  containers:
    - name: pause
      image: registry.k8s.io/pause:3.10
EOF

k -n operator-e2e apply -f - <<'EOF'
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: counted
spec:
  replicas: 2
  selector:
    matchExpressions:
      - key: app
        operator: In
        values: [operator-e2e]
      - key: tier
        operator: In
        values: [backend]
  template:
    metadata:
      labels:
        app: operator-e2e
        tier: backend
    spec:
      containers:
        - name: pause
          image: registry.k8s.io/pause:3.10
EOF

wait_pod_count() {
  local expected=$1 attempt observed
  for attempt in {1..60}; do
    observed=$(k -n operator-e2e get rs counted -o 'jsonpath={.metadata.labels.pod-count}' 2>/dev/null || true)
    if [[ "$observed" == "$expected" ]]; then
      echo "ReplicaSet pod-count=$observed"
      return 0
    fi
    sleep 2
  done
  echo "Expected pod-count=$expected; got ${observed:-<unset>}" >&2
  cat "$artifacts/controller.log" >&2
  return 1
}

# The matching Pod owned by another ReplicaSet must not be counted.
wait_pod_count 2
k -n operator-e2e scale rs/counted --replicas=1
wait_pod_count 1
