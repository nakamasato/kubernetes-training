# Helm

## Version

[v3.11.2](https://github.com/helm/helm/releases/tag/v3.11.2)

## Install helm

https://helm.sh/docs/intro/install/

```
brew install helm
```

```
helm version --short
v3.11.2+g912ebc1
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

## Helm basic commands

- `helm ls`: Check releases.
- `helm template <chart_path or chart>`: Test rendering the chart locally (You can check the resulting yaml without applying)

    You can pass values via
    1. flag: `helm template chart --set .cloudSecretManagerType=gcp`
    1. value file: `helm template chart -f gcp-helm-values.yaml`

- `helm install <NAME>`: Deploy a chart. (Deploy packaged resources to the cluster.)

    There are five different ways you can express the chart you want to install:

    1. By chart reference: `helm install mymaria example/mariadb`
    2. By path to a packaged chart: `helm install mynginx ./nginx-1.2.3.tgz`
    3. By path to an unpacked chart directory: `helm install mynginx ./nginx`
    4. By absolute URL: helm install mynginx https://example.com/charts/nginx-1.2.3.tgz
    5. By chart reference and repo url: `helm install --repo https://example.com/charts/ mynginx nginx`

    You can pass values in the same way as `helm template`

- `helm uninstall <NAME>`: Remove a chart. (Remove packaged resources from the cluster.)
- `helm status <RELEASE_NAME>`

## Reference

- [Quick Start](https://helm.sh/docs/intro/quickstart/)
- [helm (github)](https://github.com/helm/helm)
