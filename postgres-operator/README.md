# Postgres Operator

https://github.com/zalando/postgres-operator/blob/master/docs/quickstart.md

## Install Postgres Operator

Manual

Namespace: `database`

```
kubectl apply -k operator/overlays/database

Warning: kubectl apply should be used on resource created by either kubectl create --save-config or kubectl apply
namespace/database configured
serviceaccount/postgres-operator created
clusterrole.rbac.authorization.k8s.io/postgres-operator created
clusterrole.rbac.authorization.k8s.io/postgres-pod created
clusterrolebinding.rbac.authorization.k8s.io/postgres-operator created
configmap/postgres-operator created
service/postgres-operator created
deployment.apps/postgres-operator created
```

## Create a Postgres Cluster


1. Apply

    Use official yaml

    ```
    git clone https://github.com/zalando/postgres-operator.git
    sed -i.bak 's/namespace: default/namespace: database/' manifests/*yaml
    cd postgres-operator
    kubectl create -f manifests/minimal-postgres-manifest.yaml
    ```

    Use copied

    ```
    kubectl apply -f resources/minimal-postgres-manifest.yaml

    postgresql.acid.zalan.do/acid-minimal-cluster created
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
    psql (12.2 (Ubuntu 12.2-2.pgdg18.04+1))
    Type "help" for help.

    postgres=#
    ```
