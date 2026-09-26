# `e2e-framework`

[`sigs.k8s.io/e2e-framework`](https://github.com/kubernetes-sigs/e2e-framework) provides Go helpers for writing Kubernetes end-to-end tests with the standard `go test` command.

## What it provides

- `env` and `envconf` for test-environment setup and teardown
- a Kubernetes client wrapper through `klient`
- resource helpers and namespace lifecycle operations
- feature-oriented test organization and bounded environment steps
- support for kubeconfig-based external clusters

It is a test harness, not a cluster provider. Use kind, a developer cluster, or another provider to supply the Kubernetes API server.

## Choosing the test type

Use a unit test when the behavior can be verified without an API server. Use `envtest` when an API server and etcd are sufficient. Use `e2e-framework` with kind or another live cluster when the test needs real list/watch delivery, built-in controllers, Pod lifecycle, or a running example process.

## Running the repository examples

From the repository root:

```sh
# Create an isolated kind cluster and run the client-go suite.
bash scripts/e2e/run.sh client-go

# Run controller-runtime examples in an isolated kind cluster.
bash scripts/e2e/run.sh controller-runtime

# Use an existing cluster. The runner does not delete it.
KUBECONFIG="$HOME/.kube/config" \
KUBE_CONTEXT=my-context \
bash scripts/e2e/run.sh client-go
```

The tests create only temporary namespaces and objects. On failure, the E2E runner stores Kubernetes state and example-process logs under `.e2e-artifacts/`.

## Test design example

A live-cluster test should make its lifecycle explicit:

1. Obtain a client from the test environment.
2. Create a unique namespace and register cleanup immediately.
3. Start the example process with the test kubeconfig.
4. Wait for a bounded readiness condition.
5. Perform API operations and assert observable events or state.
6. Stop the process and delete the namespace, including failure paths.

The repository implementation is in [`contents/kubernetes-operator/e2e`](../../kubernetes-operator/e2e/). The client-go suite verifies clientset listing, list/watch events, and informer add/update/delete delivery. The controller-runtime suite verifies cache access, manager startup/reconciliation, and ReplicaSet `pod-count` convergence while replicas scale and Pods are removed.
