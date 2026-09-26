# kustomize

- Github: https://github.com/kubernetes-sigs/kustomize
- Docs: https://kustomize.io/

## Install

## Version

[v5.8.1](https://github.com/kubernetes-sigs/kustomize/releases/tag/kustomize/v5.8.1)

Mac:

```
brew install kustomize
```

## Example: environment overlays

`example/base` serves a generated ConfigMap with nginx. The `dev` and `prod`
overlays change the page content, namespace, and replica count (one/two).
Kustomize rewrites the Deployment's ConfigMap reference to the generated name.
This self-contained example runs on both arm64 and amd64 kind nodes.

From the repository root:

```sh
bash scripts/e2e/run.sh kustomize
```

The test builds and deploys both overlays, checks ready replicas, and verifies
that each Service serves the expected environment-specific content. It uses the
standalone `kustomize` binary. Install Docker, kind, kubectl and Python 3 as
described in [Testing](../TESTING.md).

To inspect an overlay without a cluster:

```sh
kustomize build contents/kustomize/example/overlays/prod
```

The older Flask/MySQL examples in `helm-vs-kustomize` are separate examples and
are not covered by this test.
