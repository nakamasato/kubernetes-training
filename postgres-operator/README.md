# Postgres Operator

https://github.com/zalando/postgres-operator/blob/master/docs/quickstart.md

Version: [v1.7.1](https://github.com/zalando/postgres-operator/releases/tag/v1.7.1)

## 1. Install Postgres Operator

Namespace: `default`

```
kubectl apply -k github.com/zalando/postgres-operator/manifests
```

or

```
helm install postgres-operator ./charts/postgres-operator
```

## 2. Deploy the Operator UI

1. Deploy

    ```
    kubectl apply -k github.com/zalando/postgres-operator/ui/manifests
    ```

    or

    ```
    helm install postgres-operator-ui ./charts/postgres-operator-ui
    ```

1. Check

    ```
    kubectl port-forward svc/postgres-operator-ui 8081:80
    ```

1. Open http://localhost:8081/

    ![](postgres-operator-ui.png)
## 3. Create a Postgres Cluster


1. Apply

    ```
    kubectl create -f https://raw.githubusercontent.com/zalando/postgres-operator/master/manifests/minimal-postgres-manifest.yaml
    ```

1. Check

    ```
    # check the deployed cluster
    kubectl get postgresql

    # check created database pods
    kubectl get pods -l application=spilo -L spilo-role

    # check created service resources
    kubectl get svc -l application=spilo -L spilo-role
    ```

1. Connect to Postgres cluster

    ```
    kubectl exec -it acid-minimal-cluster-0 -- psql -Upostgres
    psql (14.0 (Ubuntu 14.0-1.pgdg18.04+1))
    Type "help" for help.

    postgres=#
    ```

1. Check on UI

    http://localhost:8081/#/status/default/acid-minimal-cluster


    ![](postgres-operator-ui-cluster.png)

## 4. Delete cluster

```
kubectl delete -f resources/minimal-postgres-manifest.yaml
```

## 5. Remove operator

```
kubectl delete -k operator/overlays/database/
```
