# Helm

## Version

`v3.2.4` (https://github.com/helm/helm/releases/tag/v3.2.4)


## Install helm

https://helm.sh/docs/intro/install/

```
brew install helm
```

## Update helm version

```
brew upgrade helm
```

## Usage

### Install

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

### Customize

```
helm show values elastic/elasticsearch > helm/es-config.yaml
helm install -n eck elasticsearch elastic/elasticsearch -f helm/es-config.yaml
```

## Upgrade

```
helm upgrade elasticsearch elastic/elasticsearch -n eck -f helm/es-config.yaml
```

### Uninstall

```
helm uninstall elasticsearch elastic/elasticsearch
```
