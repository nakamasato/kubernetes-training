# Kubernetes Testing

Kubernetes tests should use the smallest environment that verifies the behavior under test. A fake client is useful for deterministic logic tests, while a live cluster is necessary for API watches, informer delivery, kubelet behavior, built-in controllers, and process integration.

## Testing layers

| Layer | Kubernetes components | Best for |
| --- | --- | --- |
| Unit test | None; fake objects or fake client | Pure logic, handlers, selectors, and reconciliation decisions |
| `controller-runtime/envtest` | API server and etcd | Controller integration without a full cluster; CRD/API behavior |
| KUTTL | A running Kubernetes cluster | Declarative resource state and eventual-state assertions |
| Go E2E with `e2e-framework` | A running Kubernetes cluster | Imperative workflows, polling, process control, and richer diagnostics |
| kind | Local Kubernetes cluster provider | Reproducible CI and developer clusters for live-cluster tests |

`envtest` does not start kubelet, controller-manager, scheduler, or other cluster components. It therefore cannot replace a kind-backed test when the test depends on Pod lifecycle or built-in reconciliation.

## This repository

- [e2e-framework](e2e-framework): Go-based live-cluster testing with kind or an existing `KUBECONFIG`.
- [Kubernetes operator E2E examples](../kubernetes-operator/e2e/): tests for the client-go and controller-runtime examples.

The operator E2E suites are selected independently by CI. They create temporary kind clusters in CI, and they can target a developer-provided cluster locally without deleting that cluster.
