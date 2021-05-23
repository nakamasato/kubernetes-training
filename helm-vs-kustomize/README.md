# Helm vs Kustomize

## Requirement

- Deploy a web application with `Deployment`.
- The web application needs to connect to `MySQL` which is given by `ConfigMap` and `Secret`.
- The `Deployment` is exposed with `Service` with `NodePort`.

## Helm

### Steps to create a chart

1. Create helm. -> would generate a folder `helm-example` and files in it.

    ```
    helm create helm-example
    ```

1. Go into the directory `helm-example` and check the generated files.

    ```
    cd helm-example
    ```

    ```
    tree
    .
    ├── Chart.yaml
    ├── charts
    ├── templates
    │   ├── NOTES.txt
    │   ├── _helpers.tpl
    │   ├── deployment.yaml
    │   ├── hpa.yaml
    │   ├── ingress.yaml
    │   ├── service.yaml
    │   ├── serviceaccount.yaml
    │   └── tests
    │       └── test-connection.yaml
    └── values.yaml

    3 directories, 10 files
    ```

1. Update templates to meet the requirements.

    1. Remove unnecessary templates.

    ```
    rm templates/hpa.yaml templates/ingress.yaml templates/serviceaccount.yaml
    ```

    1. Remove unnecessary values from values.yaml.

    1. Write your templates.

    1. You can use the following built-in objects.
        - [Built-in Objects](https://helm.sh/docs/chart_template_guide/builtin_objects/)
            - `Release`: This object describes the release itself.
                - `Release.Name`
                - `Release.Namespace`
                ...
            - `Values`: Values passed into the template from the values.yaml file.
            - `Chart`: The contents of the Chart.yaml file.
                - `Chart.Name`
                - `Chart.Version`
            - Others: `Files`, `Capabilities`, `Template`

    1. `Values`: write `values.yaml` and pass them into template yaml with `{{ Values.xxx.yyy }}`
    1. Template functions: `{{ quote .Values.favorite.drink }}` or pipelines: `{{ .Values.favorite.drink | quote }}`

1. Package a chart.

    ```
    helm package helm-example
    ```

1. Deploy the application

    ```
    helm install helm-example-0.1.0.tgz --generate-name
    ```
### Useful commands

- Dry run a wip chart.

    ```
    helm install --generate-name --debug --dry-run ./helm-example # generate random release name
    helm install test-name --debug --dry-run ./helm-example # specify release name
    ```

## Kustomize

1. Make a directory

    ```
    mkdir -p kustomize-example/{base,overlays/dev,overlays/prod} && cd kustomize-example
    ```

1. Check structure.

    ```
    tree
    .
    ├── base
    └── overlays
        ├── dev
        └── prod

    4 directories, 0 files
    ```

1. Add necessary resources to `base` folder.

    ```
    kubectl create deployment kustomize-example --image nginx --replicas=1 --dry-run=client --output yaml > kustomize-example/base/deployment.yaml
    kubectl create service nodeport kustomize-example --tcp=80:80 --dry-run=client --output yaml  > base/service.yaml
    kubectl create configmap kustomize-example --from-literal=MYSQL_HOST=mysqlhost.com --from-literal=MYSQL_USER=mysqluser --from-literal=MYSQL_PORT=3306 --dry-run=client -o yaml > base/configmap.yaml
    kubectl create secret generic kustomize-example --from-literal=MYSQL_PASSWORD=mysqlpassword --dry-run=client -o yaml > base/secret.yaml
    ```

1. Create `Namespace` `kustomize-dev` and `kustomize-prod`.

    ```
    kubectl create ns kustomize-dev --dry-run=client -o yaml > ns-kustomize-dev.yaml
    kubectl create ns kustomize-prod --dry-run=client -o yaml > ns-kustomize-prod.yaml
    kubectl apply -f ns-kustomize-dev.yaml,ns-kustomize-prod.yaml
    ```

1. Create overlays.

    1. Make each overlay same as `base`.


        - `overlays/dev/kustomization.yaml`:

            ```yaml
            namespace: kustomize-dev
            bases:
              - ../../base
            ```
        - `overlays/prod/kustomization.yaml`:

            ```yaml
            namespace: kustomize-prod
            bases:
              - ../../base
            ```

        - Check

            ```
            kubectl diff -k overlays/dev
            kustomize diff -k overlays/prod
            ```

    1. Create files to overwrite `base`.

        ```
        ```

## Example (web app with mysql)

[](example-diagram.drawio.svg)

1. Deploy dependencies.

    ```
    kubectl create ns database; kubectl apply -f dependencies/mysql.yaml -n database
    ```

1. Set up with `kustomize`

    1. Deploy `kustomize-example`

        ```
        kubectl apply -f kustomize-example/base
        ```

    1. Port forward the service.

        ```
        kubectl port-forward svc/kustomize-example 8080:80
        ```

    1. Check the application functionality.

        ```json
        curl -X POST -H "Content-Type: application/json" -d '{"name": "naka", "email": "naka@example.com"}' localhost:8080/users{"id":2,"name":"naka"}
        ```

1. Set up with `helm`

    1. Install Helm chart.
        ```
        helm install helm-example ./helm-example
        ```



## References

- https://helm.sh/docs/chart_template_guide/builtin_objects/
