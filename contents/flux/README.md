# [Flux](https://fluxcd.io/flux)

## Version

[v2.9.5](https://github.com/fluxcd/flux2/releases/tag/v2.9.5)

## Installation

```
brew install fluxcd/tap/flux
```

Verify the CLI and install the Flux controllers into the current Kubernetes
context:

```bash
flux --version
flux check --pre
flux install
flux check
```

Remove the controllers with `flux uninstall` when finished.
