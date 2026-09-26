# kind E2E tests

Run a target from the repository root (Docker, kind, kubectl, Python 3 and curl required):

```sh
bash scripts/e2e/run.sh argocd
bash scripts/e2e/run.sh prometheus-operator
bash scripts/e2e/run.sh prometheus
bash scripts/e2e/run.sh grafana
bash scripts/e2e/run.sh eck
bash scripts/e2e/run.sh helm       # also requires Helm
bash scripts/e2e/run.sh kustomize  # also requires standalone Kustomize
```

Each invocation creates a disposable kind cluster, verifies workload readiness and
an application endpoint, then deletes the cluster even on failure. Prometheus
Operator also checks that Prometheus discovers and successfully scrapes itself.
Run local targets sequentially to avoid exhausting Docker Desktop resources.
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

Prometheus tests readiness and a successful self-scrape query. Grafana tests startup
with the repository's dashboard/datasource provisioning and database health; it
does not provision the optional MySQL or RabbitMQ dashboard backends.

Helm tests the hello-world chart's install, connection hook, upgrade with optional
Ingress/HPA resources, HTTP response and uninstall. Kustomize tests the portable
`contents/kustomize/example` dev/prod overlays, content and replicas. These tests
do not cover `helm-vs-kustomize`'s older Flask/MySQL example or the packaged Helm
archive. CI installs Helm and Kustomize versions from their README version links,
so version-only README updates exercise the new binaries. Local runs use the
binaries on `PATH` and print their versions; use the documented versions for parity.

## Changed targets in CI

PRs select targets from the merge-base diff, including deleted and renamed files.
`contents/argocd/**` runs only Argo CD; `contents/prometheus-operator/**` runs only
Prometheus Operator. `contents/prometheus/**` and `contents/grafana/**` select
the corresponding standalone sample. The hello-world chart and Helm README
select `helm`; `contents/kustomize/**` selects `kustomize`. Changes to an individual case run that case; shared runner,
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

ECK tests the pinned operator and Elasticsearch/Kibana fresh installation, an
authenticated document write/read over the Elasticsearch Service, and Kibana API
health. It needs roughly 4 GiB of additional free cluster memory. The historical
Elastic Helm/Filebeat examples and upgrades of existing data are not covered.
