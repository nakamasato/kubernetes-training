# Helm hello world

1. Create helm.

    ```
    helm create helloworld-chart
    ```

    By default,
    - `Deployment`, `HPA`, `Ingress`, `Service` and `ServiceAccount` are generated under `templates`.
    - `values.yaml` has `replicaCount`, `image`, and other configuration fields for the resources.

1. (Optional) Change `values.yaml` (you can change image name, tag or any fields)
1. Run `helm package` command and it'll generate `helloworld-chart-0.1.0.tgz`.

    ```
    helm package helloworld-chart
    Successfully packaged chart and saved it to: /Users/masato-naka/repos/nakamasato/kubernetes-training/helm/hello-world/helloworld-chart-0.1.0.tgz
    ```

1. Apply helm `hello-world`.

    ```
    helm install helloworld-chart-0.1.0.tgz --generate-name
    NAME: helloworld-chart-0-1589247182
    LAST DEPLOYED: Tue May 12 10:33:03 2020
    NAMESPACE: default
    STATUS: deployed
    REVISION: 1
    NOTES:
    1. Get the application URL by running these commands:
      export POD_NAME=$(kubectl get pods --namespace default -l "app.kubernetes.io/name=helloworld-chart,app.kubernetes.io/instance=helloworld-chart-0-1589247182" -o jsonpath="{.items[0].metadata.name}")
      echo "Visit http://127.0.0.1:8080 to use your application"
      kubectl --namespace default port-forward $POD_NAME 8080:80
    ```

1. Check `helm list` (status should be `deployed`).

    ```
    helm list
    NAME                            NAMESPACE       REVISION        UPDATED                                 STATUS          CHART                        APP VERSION
    helloworld-chart-0-1621414497   default         1               2021-05-19 17:54:59.111853 +0900 JST    deployed        helloworld-chart-0.1.0       1.16.0
    ```

1. Check the created resources.

    1. `Deployment`

        ```
        kubectl get deployment -n default
        NAME                            READY   UP-TO-DATE   AVAILABLE   AGE
        helloworld-chart-0-1589247182   0/1     1            0           16s
        ```

    1. `Service`

        ```
        kubectl get svc
        NAME                            TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)   AGE
        helloworld-chart-0-1621414497   ClusterIP   10.97.169.140   <none>        80/TCP    7m58s
        kubernetes                      ClusterIP   10.96.0.1       <none>        443/TCP   12d
        ```

    1. `ServiceAccount`

        ```
        kubectl get sa
        NAME                            SECRETS   AGE
        default                         1         12d
        helloworld-chart-0-1621414497   1         8m59s
        ```

    1. `Ingress` and `HPA` are not deployed because of `enable: false` in `values.yaml`.

1. Uninstall the chart.

    ```
    helm uninstall helloworld-chart-0-1589247182
    release "helloworld-chart-0-1589247182" uninstalled
    ```

# Reference
- https://artifacthub.io/packages/search?kind=0
- https://helm.sh/docs/intro/quickstart/
- https://github.com/helm/helm
