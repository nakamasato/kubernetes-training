# 3.2 Helm

1. Create Helm chart.

    ```
    helm create <chart-name e.g. helm-example>
    ```

1. Update files under `templates` and `values.yaml`
1. Test apply.

    ```
    helm install helm-example --debug ./helm-example
    ```

1. Make a package.

    ```
    helm package helm-example
    ```

1. Create repository and set index.

    ```
    helm repo index ./ --url https://nakamasato.github.io/helm-charts-repo
    ```

1. Install a chart.

    ```
    helm repo add nakamasato https://nakamasato.github.io/helm-charts-repo
    helm repo update # update the repository info
    helm install example-from-my-repo nakamasato/helm-example
    ```
