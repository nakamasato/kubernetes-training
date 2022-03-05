# [Grafana Operator](https://github.com/grafana-operator/grafana-operator) (WIP)

## 1. Install operator

```
kubectl apply -k github.com/grafana-operator/grafana-operator/deploy/manifests/
```

## 3. Create Grafana

Be sure to deploy in the same namespace as the operator (`grafana-operator-system`).

1. Create Grafana

    Option 1 (simple one):

    ```
    kubectl apply -f https://raw.githubusercontent.com/grafana-operator/grafana-operator/master/deploy/examples/Grafana.yaml -n grafana-operator-system
    ```

    Option 2 (HA with Postgres for session storage)

    ```
    kubectl apply -k ha -n grafana-operator-system
    ```

1. Check status.
    ```
    kubectl get grafana example-grafana -n grafana-operator-system -o jsonpath='{.status}'
    {"message":"success","phase":"reconciling","previousServiceName":"grafana-service"}
    ```

1. port forward.

    ```
    kubectl port-forward -n grafana-operator-system svc/grafana-service 3000
    ```

1. Access to UI.

    http://localhost:3000

1. Log in with `admin`.

    Get password.

    ```
    kubectl get secret -n grafana-operator-system grafana-admin-credentials -o jsonpath='{.data.GF_SECURITY_ADMIN_PASSWORD}' | base64 -D
    9GD2tIHo-GrgTQ==%

    kubectl get secret -n grafana-operator-system grafana-admin-credentials -o jsonpath='{.data.GF_SECURITY_ADMIN_USER}' | base64 --decode
    admin%
    ```
## 2. Datasource

### 2.1. Prometheus Datasource

1. Deploy prometheus with [Prometheus Operator](../prometheus-operator). Prometheus Datasource expects `prometheus` service with port 9090.

    ```
    kubectl create -f https://raw.githubusercontent.com/prometheus-operator/prometheus-operator/master/bundle.yaml
    kubectl apply -k ../prometheus-operator
    ```

1. Create Prometheus Datasource
    ```
    kubectl apply -f datasource-prometheus.yaml -n grafana-operator-system
    ```

## 3. Dashboard

simple-dashboard

```
kubectl apply -f https://raw.githubusercontent.com/grafana-operator/grafana-operator/master/deploy/examples/dashboards/SimpleDashboard.yaml -n grafana-operator-system
```

keycloak-dashboard (data is empty)

```
kubectl apply -f https://raw.githubusercontent.com/grafana-operator/grafana-operator/master/deploy/examples/dashboards/KeycloakDashboard.yaml -n grafana-operator-system
```

dashboard from grafana (need node exporter)

```
kubectl apply -f https://raw.githubusercontent.com/grafana-operator/grafana-operator/7754cd15386ff6da1e3e7b820f8baf53e6dd9356/deploy/examples/dashboards/DashboardFromGrafana.yaml -n grafana-operator-system
```
## 4. Cleanup

```
kubectl delete --all grafana,grafanadashboard,grafanadatasource -n grafana-operator-system
kubectl delete -k github.com/grafana-operator/grafana-operator/deploy/manifests/
```

## Debug

1. Log: `kubectl logs -l control-plane=controller-manager -c manager -n grafana-operator-system -f`
