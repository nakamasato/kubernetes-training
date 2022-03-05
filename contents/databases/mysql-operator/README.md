# MySQL Operator

MySQL Version: 8.0.28
mysqlsh version: `mysqlsh   Ver 8.0.28 for macos11 on x86_64 - for MySQL 8.0.28 (MySQL Community Server (GPL))`

## Getting Started

1. Install CRDs and operator.

    ```
    kubectl apply -f https://raw.githubusercontent.com/mysql/mysql-operator/trunk/deploy/deploy-crds.yaml
    kubectl apply -f https://raw.githubusercontent.com/mysql/mysql-operator/trunk/deploy/deploy-operator.yaml
    ```

    ```
    kubectl get deployments -n mysql-operator
    NAME             READY   UP-TO-DATE   AVAILABLE   AGE
    mysql-operator   1/1     1            1           107s
    ```

1. Create Secret for root user.

    ```
    kubectl create secret generic mypwds \
            --from-literal=rootUser=root \
            --from-literal=rootHost=% \
            --from-literal=rootPassword=password
    ```

1. Create InnoDBCluster.

    ```
    kubectl apply -f https://raw.githubusercontent.com/mysql/mysql-operator/trunk/samples/sample-cluster.yaml
    ```

    check cluster:

    ```
    kubectl get innodbcluster

    NAME        STATUS   ONLINE   INSTANCES   ROUTERS   AGE
    mycluster   ONLINE   3        3           1         2m8s
    ```

    Check logs

    ```
    kubectl logs -l name=mysql-operator -n mysql-operator -f
    ```

    <details><summary>logs</summary>

    ```
    [2022-02-08 21:52:52,983] kopf.objects         [INFO    ] Creation is processed: 1 succeeded; 0 failed.
    [2022-02-08 21:52:53,117] kopf.objects         [INFO    ] POD EVENT : pod=mycluster-2 containers_ready=True deleting=False phase=Running member_info={'memberId': '67632728-8929-11ec-bee4-7e577a1faa95', 'lastTransitionTime': '2022-02-08T21:52:52Z', 'lastProbeTime': '2022-02-08T21:52:52Z', 'groupViewId': '16443571008025241:3', 'status': 'ONLINE', 'version': '8.0.28', 'role': 'SECONDARY', 'joinTime': '2022-02-08T21:52:52Z'} restarts=0 containers=['mysql=ready', 'sidecar=ready'] conditions=['mysql.oracle.com/ready=True', 'mysql.oracle.com/configured=True', 'Initialized=True', 'Ready=False', 'ContainersReady=True', 'PodScheduled=True']
    2022-02-08 21:52:53: Info: About to connect to MySQL at: mysql://mysqladmin@mycluster-0.mycluster-instances.default.svc.cluster.local:3306
    2022-02-08 21:52:53: Info: About to connect to MySQL at: mysql://mysqladmin@mycluster-0.mycluster-instances.default.svc.cluster.local:3306?connect-timeout=5000
    2022-02-08 21:52:53: Info: Group Replication 'group_name' value: 4c318de6-8929-11ec-95c2-36381ebfe6c7
    2022-02-08 21:52:53: Info: Metadata 'group_name' value: 4c318de6-8929-11ec-95c2-36381ebfe6c7
    2022-02-08 21:52:53: Info: About to connect to MySQL at: mysql://mysqladmin@mycluster-0.mycluster-instances.default.svc.cluster.local:3306?connect-timeout=5000
    2022-02-08 21:52:53: Info: About to connect to MySQL at: mysql://mysqladmin@mycluster-1.mycluster-instances.default.svc.cluster.local:3306?connect-timeout=5000
    2022-02-08 21:52:53: Info: About to connect to MySQL at: mysql://mysqladmin@mycluster-2.mycluster-instances.default.svc.cluster.local:3306?connect-timeout=5000
    [2022-02-08 21:52:53,261] kopf.objects         [INFO    ] diag instance mycluster-0 --> InstanceDiagStatus.ONLINE quorum=True
    2022-02-08 21:52:53: Info: About to connect to MySQL at: mysql://mysqladmin@mycluster-2.mycluster-instances.default.svc.cluster.local:3306
    2022-02-08 21:52:53: Info: About to connect to MySQL at: mysql://mysqladmin@mycluster-2.mycluster-instances.default.svc.cluster.local:3306?connect-timeout=5000
    2022-02-08 21:52:53: Info: Group Replication 'group_name' value: 4c318de6-8929-11ec-95c2-36381ebfe6c7
    2022-02-08 21:52:53: Info: Metadata 'group_name' value: 4c318de6-8929-11ec-95c2-36381ebfe6c7
    2022-02-08 21:52:53: Info: About to connect to MySQL at: mysql://mysqladmin@mycluster-0.mycluster-instances.default.svc.cluster.local:3306?connect-timeout=5000
    2022-02-08 21:52:53: Info: About to connect to MySQL at: mysql://mysqladmin@mycluster-1.mycluster-instances.default.svc.cluster.local:3306?connect-timeout=5000
    2022-02-08 21:52:53: Info: About to connect to MySQL at: mysql://mysqladmin@mycluster-2.mycluster-instances.default.svc.cluster.local:3306?connect-timeout=5000
    [2022-02-08 21:52:53,393] kopf.objects         [INFO    ] diag instance mycluster-2 --> InstanceDiagStatus.ONLINE quorum=True
    2022-02-08 21:52:53: Info: About to connect to MySQL at: mysql://mysqladmin@mycluster-1.mycluster-instances.default.svc.cluster.local:3306
    2022-02-08 21:52:53: Info: About to connect to MySQL at: mysql://mysqladmin@mycluster-1.mycluster-instances.default.svc.cluster.local:3306?connect-timeout=5000
    2022-02-08 21:52:53: Info: Group Replication 'group_name' value: 4c318de6-8929-11ec-95c2-36381ebfe6c7
    2022-02-08 21:52:53: Info: Metadata 'group_name' value: 4c318de6-8929-11ec-95c2-36381ebfe6c7
    2022-02-08 21:52:53: Info: About to connect to MySQL at: mysql://mysqladmin@mycluster-0.mycluster-instances.default.svc.cluster.local:3306?connect-timeout=5000
    2022-02-08 21:52:53: Info: About to connect to MySQL at: mysql://mysqladmin@mycluster-1.mycluster-instances.default.svc.cluster.local:3306?connect-timeout=5000
    2022-02-08 21:52:53: Info: About to connect to MySQL at: mysql://mysqladmin@mycluster-2.mycluster-instances.default.svc.cluster.local:3306?connect-timeout=5000
    [2022-02-08 21:52:53,509] kopf.objects         [INFO    ] diag instance mycluster-1 --> InstanceDiagStatus.ONLINE quorum=True
    [2022-02-08 21:52:53,510] kopf.objects         [INFO    ] mycluster: all={<MySQLPod mycluster-0>, <MySQLPod mycluster-2>, <MySQLPod mycluster-1>}  members={<MySQLPod mycluster-0>, <MySQLPod mycluster-2>, <MySQLPod mycluster-1>}  online={<MySQLPod mycluster-0>, <MySQLPod mycluster-2>, <MySQLPod mycluster-1>}  offline=set()  unsure=set()
    [2022-02-08 21:52:53,521] kopf.objects         [INFO    ] cluster probe: status=ClusterDiagStatus.ONLINE online=[<MySQLPod mycluster-0>, <MySQLPod mycluster-1>, <MySQLPod mycluster-2>]
    [2022-02-08 21:52:53,524] kopf.objects         [INFO    ] Handler 'on_pod_event' succeeded.
    [2022-02-08 21:52:53,655] kopf.objects         [INFO    ] POD EVENT : pod=mycluster-2 containers_ready=True deleting=False phase=Running member_info={'memberId': '67632728-8929-11ec-bee4-7e577a1faa95', 'lastTransitionTime': '2022-02-08T21:52:52Z', 'lastProbeTime': '2022-02-08T21:52:52Z', 'groupViewId': '16443571008025241:3', 'status': 'ONLINE', 'version': '8.0.28', 'role': 'SECONDARY', 'joinTime': '2022-02-08T21:52:52Z'} restarts=0 containers=['mysql=ready', 'sidecar=ready'] conditions=['mysql.oracle.com/ready=True', 'mysql.oracle.com/configured=True', 'Initialized=True', 'Ready=True', 'ContainersReady=True', 'PodScheduled=True']
    ```

    </details>

1. Connect to MySQL InnoDB Cluster.

    ```
    kubectl port-forward service/mycluster mysql
    ```

    ※ If you don't have `mysqlsh`, you can install from [Chapter 2 Installing MySQL Shell](https://dev.mysql.com/doc/mysql-shell/8.0/en/mysql-shell-install.html).

    ```
    mysqlsh -h127.0.0.1 -P6446 -uroot -p
    ```

    ```
    \sql
    show databases;
    +-------------------------------+
    | Database                      |
    +-------------------------------+
    | information_schema            |
    | mysql                         |
    | mysql_innodb_cluster_metadata |
    | performance_schema            |
    | sys                           |
    +-------------------------------+
    5 rows in set (0.0034 sec)
    ```

1. Check cluster status

    ```
    mysqlsh -h127.0.0.1 -P6446 -uroot -p --cluster
    ```

    ```js
    MySQL  127.0.0.1:6446 ssl  JS > cluster.status()
    {
        "clusterName": "mycluster",
        "defaultReplicaSet": {
            "name": "default",
            "primary": "mycluster-0.mycluster-instances.default.svc.cluster.local:3306",
            "ssl": "REQUIRED",
            "status": "OK",
            "statusText": "Cluster is ONLINE and can tolerate up to ONE failure.",
            "topology": {
                "mycluster-0.mycluster-instances.default.svc.cluster.local:3306": {
                    "address": "mycluster-0.mycluster-instances.default.svc.cluster.local:3306",
                    "memberRole": "PRIMARY",
                    "memberState": "(MISSING)",
                    "mode": "n/a",
                    "readReplicas": {},
                    "role": "HA",
                    "shellConnectError": "MySQL Error 2005: Could not open connection to 'mycluster-0.mycluster-instances.default.svc.cluster.local:3306': Unknown MySQL server host 'mycluster-0.mycluster-instances.default.svc.cluster.local' (8)",
                    "status": "ONLINE",
                    "version": "8.0.28"
                },
                "mycluster-1.mycluster-instances.default.svc.cluster.local:3306": {
                    "address": "mycluster-1.mycluster-instances.default.svc.cluster.local:3306",
                    "memberRole": "SECONDARY",
                    "memberState": "(MISSING)",
                    "mode": "n/a",
                    "readReplicas": {},
                    "role": "HA",
                    "shellConnectError": "MySQL Error 2005: Could not open connection to 'mycluster-1.mycluster-instances.default.svc.cluster.local:3306': Unknown MySQL server host 'mycluster-1.mycluster-instances.default.svc.cluster.local' (8)",
                    "status": "ONLINE",
                    "version": "8.0.28"
                },
                "mycluster-2.mycluster-instances.default.svc.cluster.local:3306": {
                    "address": "mycluster-2.mycluster-instances.default.svc.cluster.local:3306",
                    "memberRole": "SECONDARY",
                    "memberState": "(MISSING)",
                    "mode": "n/a",
                    "readReplicas": {},
                    "role": "HA",
                    "shellConnectError": "MySQL Error 2005: Could not open connection to 'mycluster-2.mycluster-instances.default.svc.cluster.local:3306': Unknown MySQL server host 'mycluster-2.mycluster-instances.default.svc.cluster.local' (8)",
                    "status": "ONLINE",
                    "version": "8.0.28"
                }
            },
            "topologyMode": "Single-Primary"
        },
        "groupInformationSourceMember": "mycluster-0.mycluster-instances.default.svc.cluster.local:3306"
    }
    ```

1. Clean up.

    ```
    kubectl delete -f https://raw.githubusercontent.com/mysql/mysql-operator/trunk/deploy/deploy-crds.yaml
    kubectl delete -f https://raw.githubusercontent.com/mysql/mysql-operator/trunk/deploy/deploy-operator.yaml
    ```

## References

- https://speakerdeck.com/oracle4engineer/mysql-technology-cafe-number-13-oracle-mysql-operator-for-kubernetes-purebiyuririsuban-jie-shuo
- [Chapter 2 Installing MySQL Shell](https://dev.mysql.com/doc/mysql-shell/8.0/en/mysql-shell-install.html)
- [3.1 MySQL Shell Commands](https://dev.mysql.com/doc/mysql-shell/8.0/en/mysql-shell-commands.html)
