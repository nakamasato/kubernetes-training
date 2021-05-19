# Helm

## Version

`v3.5.4` (https://github.com/helm/helm/releases/tag/v3.5.4)


## Install helm

https://helm.sh/docs/intro/install/

```
brew install helm
```

```
helm version --short
v3.5.4+g1b5edb6
```

## Update helm version

```
brew upgrade helm
```

## Usage

### 1. Install

```
helm repo add elastic https://helm.elastic.co
```

```
helm install elasticsearch elastic/elasticsearch
NAME: elasticsearch
LAST DEPLOYED: Sat Aug  8 17:23:21 2020
NAMESPACE: default
STATUS: deployed
REVISION: 1
NOTES:
1. Watch all cluster members come up.
  $ kubectl get pods --namespace=default -l app=elasticsearch-master -w
2. Test cluster health using Helm test.
  $ helm test elasticsearch --cleanup
```

### 2. Customize

```
helm show values elastic/elasticsearch > helm/es-config.yaml
helm install -n eck elasticsearch elastic/elasticsearch -f helm/es-config.yaml
```

### 3. Upgrade

```
helm upgrade elasticsearch elastic/elasticsearch -n eck -f helm/es-config.yaml
```

### 4. Uninstall

```
helm uninstall elasticsearch elastic/elasticsearch
```

## Reference

- [Quick Start](https://helm.sh/docs/intro/quickstart/)
- [helm (github)](https://github.com/helm/helm)
