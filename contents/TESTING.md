# kind E2E tests

Run a target from the repository root (Docker, kind, kubectl, Python 3 and curl required):

```sh
bash scripts/e2e/run.sh argocd
bash scripts/e2e/run.sh prometheus-operator
```

Each invocation creates a disposable kind cluster, verifies workload readiness and
an application endpoint, then deletes the cluster even on failure. Prometheus
Operator also checks that Prometheus discovers and successfully scrapes itself.
The versions come from the sample manifests, including the pinned Operator bundle
in `prometheus-operator/operator/kustomization.yaml`.

Every kubectl call uses an explicit `--context` and `--kubeconfig`. Both kind create
and delete use a fresh, private kubeconfig. Neither `$HOME/.kube/config` nor files
in an existing `$KUBECONFIG` are written; your current-context stays unchanged.
Tests do not reuse or delete existing clusters. Do not replace `k` with bare kubectl
inside a case.

The default node image is `kindest/node:v1.36.1`. To test another Kubernetes version:

```sh
KIND_NODE_IMAGE=kindest/node:v1.37.0 bash scripts/e2e/run.sh argocd
```

Failure diagnostics (events, pod descriptions/logs and kind logs) are saved in
`.e2e-artifacts/<target>/` and uploaded by CI. Override this with `E2E_ARTIFACTS`.
A forced process kill can prevent cleanup; use the cluster name in the run log
and an explicit disposable `--kubeconfig` when deleting a leftover cluster.

## Changed targets in CI

PRs select targets from the merge-base diff, including deleted and renamed files.
`contents/argocd/**` runs only Argo CD; `contents/prometheus-operator/**` runs only
Prometheus Operator. Changes to an individual case run that case; shared runner,
registry or workflow changes run all registered targets. Unrelated changes skip
cluster jobs. README changes under a target are included because installation
versions and instructions can live there. Manual `workflow_dispatch` runs all.
`status-check-e2e` is the stable aggregate check, including when no targets match.
Go-only examples remain covered by the separate Go workflow.

## Add coverage

1. Add `scripts/e2e/cases/<target>.sh`, using `k` for all kubectl operations.
2. Register source paths and any shared dependencies in `scripts/e2e/targets.json`.
3. Check rollout readiness and actual application behavior with bounded timeouts.
4. Run the case locally, and `python3 -m unittest discover -s scripts/e2e`.

Only registered targets have E2E coverage; this does not claim coverage of every
sample in the repository. Helm cases must pass `--kube-context "$context"` and
`--kubeconfig "$KUBECONFIG"` explicitly.
