# Kubernetes the hard way

https://github.com/kelseyhightower/kubernetes-the-hard-way

## 01 Prerequisite

https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/01-prerequisites.md

- install gcloud https://cloud.google.com/sdk/docs/quickstart#mac

```
○ ./google-cloud-sdk/bin/gcloud init
Welcome! This command will take you through the configuration of gcloud.

Settings from your current configuration [default] are:
core:
  account: xxx@gmail.com
  disable_usage_reporting: 'False'

Pick configuration to use:
 [1] Re-initialize this configuration [default] with new settings
 [2] Create a new configuration
Please enter your numeric choice:  1

Your current configuration has been set to: [default]

You can skip diagnostics next time by using the following flag:
  gcloud init --skip-diagnostics

Network diagnostic detects and fixes local network connection issues.
Checking network connection...done.
Reachability Check passed.
Network diagnostic passed (1/1 checks passed).

Choose the account you would like to use to perform operations for
this configuration:
 [1] xxx@gmail.com
 [2] Log in with a new account
Please enter your numeric choice:  1

You are logged in as: [xxx@gmail.com].

Pick cloud project to use:
 [1] xxx
 [2] xxx
 ...
 [5] Create a new project
Please enter numeric choice or text value (must exactly match list
item):  5

Enter a Project ID. Note that a Project ID CANNOT be changed later.
Project IDs must be 6-30 characters (lowercase ASCII, digits, or
hyphens) in length and start with a lowercase letter. k8s-the-hard-way-20210213
Waiting for [operations/cp.7017970540723547437] to finish...done.
Your current project has been set to: [k8s-the-hard-way-20210213].
```

```
gcloud config list
[compute]
region = us-west1
zone = us-west1-c
[core]
account = masatonaka1989@gmail.com
disable_usage_reporting = False
project = k8s-the-hard-way-20210213

Your active configuration is: [default]
```


## [03 Compute Resources](https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/03-compute-resources.md)

### Create custom VPC networks

```
gcloud compute networks create kubernetes-the-hard-way --subnet-mode custom
```


If it fails, you might need to enable Compute Engine API on console.

<details>

```
± gcloud compute networks create kubernetes-the-hard-way --subnet-mode custom
ERROR: (gcloud.compute.networks.create) Could not fetch resource:
 - Project XXXX is not found and cannot be used for API calls. If it is recently created, enable Compute Engine API by visiting https://console.developers.google.com/apis/api/compute.googleapis.com/overview?project=XXXXX then retry. If you enabled this API recently, wait a few minutes for the action to propagate to our systems and retry.
```

```
 gcloud compute networks create kubernetes-the-hard-way --subnet-mode custom
Created [https://www.googleapis.com/compute/v1/projects/k8s-the-hard-way-20210213/global/networks/kubernetes-the-hard-way].
NAME                     SUBNET_MODE  BGP_ROUTING_MODE  IPV4_RANGE  GATEWAY_IPV4
kubernetes-the-hard-way  CUSTOM       REGIONAL

Instances on this network will not be reachable until firewall rules
are created. As an example, you can allow all internal traffic between
instances as well as SSH, RDP, and ICMP by running:

$ gcloud compute firewall-rules create <FIREWALL_NAME> --network kubernetes-the-hard-way --allow tcp,udp,icmp --source-ranges <IP_RANGE>
$ gcloud compute firewall-rules create <FIREWALL_NAME> --network kubernetes-the-hard-way --allow tcp:22,tcp:3389,icmp
```

</details>

### networks subnets

```
gcloud compute networks subnets create kubernetes \
  --network kubernetes-the-hard-way \
  --range 10.240.0.0/24
```

### firewall-rules

internal:

```
gcloud compute firewall-rules create kubernetes-the-hard-way-allow-internal \
  --allow tcp,udp,icmp \
  --network kubernetes-the-hard-way \
  --source-ranges 10.240.0.0/24,10.200.0.0/16
```

external:

```
gcloud compute firewall-rules create kubernetes-the-hard-way-allow-external \
  --allow tcp:22,tcp:6443,icmp \
  --network kubernetes-the-hard-way \
  --source-ranges 0.0.0.0/0
```

Check:

```
gcloud compute firewall-rules list --filter="network:kubernetes-the-hard-way"

NAME                                    NETWORK                  DIRECTION  PRIORITY  ALLOW                 DENY  DISABLED
kubernetes-the-hard-way-allow-external  kubernetes-the-hard-way  INGRESS    1000      tcp:22,tcp:6443,icmp        False
kubernetes-the-hard-way-allow-internal  kubernetes-the-hard-way  INGRESS    1000      tcp,udp,icmp                False

To show all fields of the firewall, please show in JSON format: --format=json
To show all fields in table format, please see the examples in --help.
```

### Static Public IP Address

```
gcloud compute addresses create kubernetes-the-hard-way \
  --region $(gcloud config get-value compute/region)
```

### instances

**Control Plane:**

```
for i in 0 1 2; do
  gcloud compute instances create controller-${i} \
    --async \
    --boot-disk-size 200GB \
    --can-ip-forward \
    --image-family ubuntu-1804-lts \
    --image-project ubuntu-os-cloud \
    --machine-type n1-standard-1 \
    --private-network-ip 10.240.0.1${i} \
    --scopes compute-rw,storage-ro,service-management,service-control,logging-write,monitoring \
    --subnet kubernetes \
    --tags kubernetes-the-hard-way,controller
done
```

**Workers**:

```
for i in 0 1 2; do
  gcloud compute instances create worker-${i} \
    --async \
    --boot-disk-size 200GB \
    --can-ip-forward \
    --image-family ubuntu-1804-lts \
    --image-project ubuntu-os-cloud \
    --machine-type n1-standard-1 \
    --metadata pod-cidr=10.200.${i}.0/24 \
    --private-network-ip 10.240.0.2${i} \
    --scopes compute-rw,storage-ro,service-management,service-control,logging-write,monitoring \
    --subnet kubernetes \
    --tags kubernetes-the-hard-way,worker
done
```

<details>

```
gcloud compute instances list --filter="tags.items=kubernetes-the-hard-way"

NAME          ZONE        MACHINE_TYPE   PREEMPTIBLE  INTERNAL_IP  EXTERNAL_IP      STATUS
controller-0  us-west1-c  e2-standard-2               10.240.0.10  35.203.130.132   RUNNING
controller-1  us-west1-c  e2-standard-2               10.240.0.11  35.197.106.92    RUNNING
controller-2  us-west1-c  e2-standard-2               10.240.0.12  104.196.237.244  RUNNING
worker-0      us-west1-c  e2-standard-2               10.240.0.20  34.82.0.74       RUNNING
worker-1      us-west1-c  e2-standard-2               10.240.0.21  35.247.104.63    RUNNING
worker-2      us-west1-c  e2-standard-2               10.240.0.22  35.230.29.250    RUNNING
```

</details>

## [04 Provisioning a CA and Generating TLS Certificates](https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/04-certificate-authority.md)


1. Provision a Certificate Authority

    1. Prepare `ca-config.json` + `ca.csr.json`
    1. `cfssl gencert -initca ca-csr.json | cfssljson -bare ca` -> generate `ca-key.pem`, `ca.csr`, `ca.pem`

1. Client and Server certificates

    1. Prepare `admin-csr.json`
    1. `cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=kubernetes admin-csr.json | cfssljson -bare admin` -> generate `admin-key.pem`, `admin.pem`, `admin.csr`

1. Kubelet Client Certificates

    1. For each worker node, generate csr json and generate certificates.

        ```
        for instance in worker-0 worker-1 worker-2; do
        cat > ${instance}-csr.json <<EOF
        {
          "CN": "system:node:${instance}",
          "key": {
            "algo": "rsa",
            "size": 2048
          },
          "names": [
            {
              "C": "US",
              "L": "Portland",
              "O": "system:nodes",
              "OU": "Kubernetes The Hard Way",
              "ST": "Oregon"
            }
          ]
        }
        EOF

        EXTERNAL_IP=$(gcloud compute instances describe ${instance} \
          --format 'value(networkInterfaces[0].accessConfigs[0].natIP)')

        INTERNAL_IP=$(gcloud compute instances describe ${instance} \
          --format 'value(networkInterfaces[0].networkIP)')

        cfssl gencert \
          -ca=ca.pem \
          -ca-key=ca-key.pem \
          -config=ca-config.json \
          -hostname=${instance},${EXTERNAL_IP},${INTERNAL_IP} \
          -profile=kubernetes \
          ${instance}-csr.json | cfssljson -bare ${instance}
        done
        ```

1. The Controller Manager Client Certificate
1. The Kube Proxy Client Certificate
1. The Scheduler Client Certificate
1. The Kubernetes API Server Certificate
1. The Service Account Key Pair
1. Distribute the Client and Server Certificates

    1. Worker instance

        ```
        for instance in worker-0 worker-1 worker-2; do
          gcloud compute scp ca.pem ${instance}-key.pem ${instance}.pem ${instance}:~/
        done
        ```
    1. Controller instance

        ```
        for instance in controller-0 controller-1 controller-2; do
          gcloud compute scp ca.pem ca-key.pem kubernetes-key.pem kubernetes.pem \
            service-account-key.pem service-account.pem ${instance}:~/
        done
        ```

## [05 Generating Kubernetes Configuration Files for Authentication](https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/05-kubernetes-configuration-files.md)

Kubeconfigs: enable Kubernetes clients to locate and authenticate to the Kubernetes API servers

- kubelet (set-cluster, set-credentials, set-context for each worker)

    ```
      kubectl config set-cluster kubernetes-the-hard-way \
        --certificate-authority=ca.pem \
        --embed-certs=true \
        --server=https://${KUBERNETES_PUBLIC_ADDRESS}:6443 \
        --kubeconfig=${instance}.kubeconfig
      kubectl config set-credentials system:node:${instance} \
        --client-certificate=${instance}.pem \
        --client-key=${instance}-key.pem \
        --embed-certs=true \
        --kubeconfig=${instance}.kubeconfig

      kubectl config set-context default \
        --cluster=kubernetes-the-hard-way \
        --user=system:node:${instance} \
        --kubeconfig=${instance}.kubeconfig
    ```

- kube-proxy (set-cluster, set-credentials, set-context (generate only one))
    ```
      kubectl config set-cluster kubernetes-the-hard-way \
        --certificate-authority=ca.pem \
        --embed-certs=true \
        --server=https://${KUBERNETES_PUBLIC_ADDRESS}:6443 \
        --kubeconfig=kube-proxy.kubeconfig

      kubectl config set-credentials system:kube-proxy \
        --client-certificate=kube-proxy.pem \
        --client-key=kube-proxy-key.pem \
        --embed-certs=true \
        --kubeconfig=kube-proxy.kubeconfig

      kubectl config set-context default \
        --cluster=kubernetes-the-hard-way \
        --user=system:kube-proxy \
        --kubeconfig=kube-proxy.kubeconfig
    ```
- kube-controller-manager (almost same as kube-proxy)
- kube-scheduler (almost same as kube-proxy)
- admin (almost same as kube-proxy)

Distribute Kuberenetes configuration file
- workers (worker-x, kube-proxy)
- controllers (admin, kubec-controller-manager, kube-scheduler)

## [06 Generating the Data Encryption Config and Key](https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/06-data-encryption-keys.md)

> Kubernetes stores a variety of data including cluster state, application configurations, and secrets. Kubernetes supports the ability to encrypt cluster data at rest.

= Encrypt data in etcd.

1. Generate encryption-config.yaml with `EncryptionConfig` resource with ENCRYPTION_KEY generated in local.
1. Put the yaml file in each controler node.

## [07 Bootstrapping the etcd Cluster](https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/07-bootstrapping-etcd.md)

Run the command to set up etcd on the controller-0, controller-1, and controller-2 (I used iterm2's toggle input to run exactly the same command in all of the server at the same time)

In each controller node,

1. Download and Install the etcd Binaries

    `wget -q --show-progress --https-only --timestamping "https://github.com/etcd-io/etcd/releases/download/v3.4.10/etcd-v3.4.10-linux-amd64.tar.gz"`

    ```
    tar -xvf etcd-v3.4.10-linux-amd64.tar.gz
    sudo mv etcd-v3.4.10-linux-amd64/etcd* /usr/local/bin/
    ```
1. Configure the etcd Server (`/etc/etcd/`)

    ```
    sudo mkdir -p /etc/etcd /var/lib/etcd
    sudo chmod 700 /var/lib/etcd
    sudo cp ca.pem kubernetes-key.pem kubernetes.pem /etc/etcd/
    ```

    1. `INTERNAL_IP`

    1. Set `ETCD_NAME`

    1. Set `/etc/systemd/system/etcd.service`

1. Start etcd

    ```
    sudo systemctl daemon-reload
    sudo systemctl enable etcd
    sudo systemctl start etcd
    ```


## [08 Bootstrapping the Kubernetes Control Plane](https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/08-bootstrapping-kubernetes-controllers.md)

The following commands are run in all the controller nodes

1. Create dir `/etc/kubernetes/config`

1. Download `kube-api-server`, `kube-controller-manager`, `kube-scheduler`, and `kubectl` and move to `/usr/local/bin`

```
wget -q --show-progress --https-only --timestamping \
  "https://storage.googleapis.com/kubernetes-release/release/v1.18.6/bin/linux/amd64/kube-apiserver" \
  "https://storage.googleapis.com/kubernetes-release/release/v1.18.6/bin/linux/amd64/kube-controller-manager" \
  "https://storage.googleapis.com/kubernetes-release/release/v1.18.6/bin/linux/amd64/kube-scheduler" \
  "https://storage.googleapis.com/kubernetes-release/release/v1.18.6/bin/linux/amd64/kubectl"
```

1. Configure the Kubernetes API Server

    1. Prepare `/etc/systemd/system/kube-apiserver.service`


1. Configure the Kubernetes Controller Manager

    1. Prepare `/etc/systemd/system/kube-controller-manager.service`

1. Configure the Kubernetes Scheduler


    1. `/etc/kubernetes/config/kube-scheduler.yaml`

        ```yaml
        apiVersion: kubescheduler.config.k8s.io/v1alpha1
        kind: KubeSchedulerConfiguration
        clientConnection:
          kubeconfig: "/var/lib/kubernetes/kube-scheduler.kubeconfig"
        leaderElection:
          leaderElect: true
        ```

    1. Prepare `/etc/systemd/system/kube-scheduler.service`

        ```
        ExecStart=/usr/local/bin/kube-scheduler \\
          --config=/etc/kubernetes/config/kube-scheduler.yaml \\
          --v=2
        ```

1. Start `kube-apiserver`, `kube-controller-manager`, `kube-scheduler`

1. Enable HTTP Health Checks

    1. Install nginx

    1. Configuration of `kubernetes.default.svc.cluster.local`

        ```
        server {
          listen      80;
          server_name kubernetes.default.svc.cluster.local;

          location /healthz {
            proxy_pass                    https://127.0.0.1:6443/healthz;
            proxy_ssl_trusted_certificate /var/lib/kubernetes/ca.pem;
          }
        }
        ```

    1. Restart nginx

1. Verification

    ```
    kubectl get componentstatuses --kubeconfig admin.kubeconfig
    NAME                 STATUS    MESSAGE             ERROR
    scheduler            Healthy   ok
    controller-manager   Healthy   ok
    etcd-1               Healthy   {"health":"true"}
    etcd-0               Healthy   {"health":"true"}
    etcd-2               Healthy   {"health":"true"}
    ```

    ```
    curl -H "Host: kubernetes.default.svc.cluster.local" -i http://127.0.0.1/healthz
    HTTP/1.1 200 OK
    Server: nginx/1.18.0 (Ubuntu)
    Date: Sat, 13 Feb 2021 08:45:57 GMT
    Content-Type: text/plain; charset=utf-8
    Content-Length: 2
    Connection: keep-alive
    Cache-Control: no-cache, private
    X-Content-Type-Options: nosniff
    ```

1. RBAC for Kubelet Authorization (only on one controller)

    ClusterRole

    ```bash
    cat <<EOF | kubectl apply --kubeconfig admin.kubeconfig -f -
    apiVersion: rbac.authorization.k8s.io/v1beta1
    kind: ClusterRole
    metadata:
      annotations:
        rbac.authorization.kubernetes.io/autoupdate: "true"
      labels:
        kubernetes.io/bootstrapping: rbac-defaults
      name: system:kube-apiserver-to-kubelet
    rules:
      - apiGroups:
          - ""
        resources:
          - nodes/proxy
          - nodes/stats
          - nodes/log
          - nodes/spec
          - nodes/metrics
        verbs:
          - "*"
    EOF
    ```

    ClusterRoleBinding

    ```bash
    cat <<EOF | kubectl apply --kubeconfig admin.kubeconfig -f -
    apiVersion: rbac.authorization.k8s.io/v1beta1
    kind: ClusterRoleBinding
    metadata:
      name: system:kube-apiserver
      namespace: ""
    roleRef:
      apiGroup: rbac.authorization.k8s.io
      kind: ClusterRole
      name: system:kube-apiserver-to-kubelet
    subjects:
      - apiGroup: rbac.authorization.k8s.io
        kind: User
        name: kubernetes
    EOF
    ```

1. The Kubernetes Frontend Load Balancer

    1. Provision a Network Load Balancer


        ```
        KUBERNETES_PUBLIC_ADDRESS=$(gcloud compute addresses describe kubernetes-the-hard-way \
            --region $(gcloud config get-value compute/region) \
            --format 'value(address)')

          gcloud compute http-health-checks create kubernetes \
            --description "Kubernetes Health Check" \
            --host "kubernetes.default.svc.cluster.local" \
            --request-path "/healthz"

          gcloud compute firewall-rules create kubernetes-the-hard-way-allow-health-check \
            --network kubernetes-the-hard-way \
            --source-ranges 209.85.152.0/22,209.85.204.0/22,35.191.0.0/16 \
            --allow tcp

          gcloud compute target-pools create kubernetes-target-pool \
            --http-health-check kubernetes

          gcloud compute target-pools add-instances kubernetes-target-pool \
          --instances controller-0,controller-1,controller-2

          gcloud compute forwarding-rules create kubernetes-forwarding-rule \
            --address ${KUBERNETES_PUBLIC_ADDRESS} \
            --ports 6443 \
            --region $(gcloud config get-value compute/region) \
            --target-pool kubernetes-target-pool
        ```

    1. Verification

        ```
        curl --cacert ca.pem https://${KUBERNETES_PUBLIC_ADDRESS}:6443/version

        {
          "major": "1",
          "minor": "18",
          "gitVersion": "v1.18.6",
          "gitCommit": "dff82dc0de47299ab66c83c626e08b245ab19037",
          "gitTreeState": "clean",
          "buildDate": "2020-07-15T16:51:04Z",
          "goVersion": "go1.13.9",
          "compiler": "gc",
          "platform": "linux/amd64"
        }%
        ```

## [Bootstrapping the Kubernetes Worker Nodes](https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/09-bootstrapping-kubernetes-workers.md)


Login to all the worker nodes `gcloud compute ssh worker-0`

1. Install `socat`, `conntrack`, `ipset`

    ```
    sudo apt-get update && sudo apt-get -y install socat conntrack ipset
    ```

1. Download and Install Worker Binaries (`crictl`, `runc`, `cni-plugins-linx`, `containerd`, `kubectl`, `kube-proxy`, `kubelet`)

    ```
    wget -q --show-progress --https-only --timestamping \
      https://github.com/kubernetes-sigs/cri-tools/releases/download/v1.18.0/crictl-v1.18.0-linux-amd64.tar.gz \
      https://github.com/opencontainers/runc/releases/download/v1.0.0-rc91/runc.amd64 \
      https://github.com/containernetworking/plugins/releases/download/v0.8.6/cni-plugins-linux-amd64-v0.8.6.tgz \
      https://github.com/containerd/containerd/releases/download/v1.3.6/containerd-1.3.6-linux-amd64.tar.gz \
      https://storage.googleapis.com/kubernetes-release/release/v1.18.6/bin/linux/amd64/kubectl \
      https://storage.googleapis.com/kubernetes-release/release/v1.18.6/bin/linux/amd64/kube-proxy \
      https://storage.googleapis.com/kubernetes-release/release/v1.18.6/bin/linux/amd64/kubelet
    ```

1. Make dir
1. Install binaries

    ```
      mkdir containerd
      tar -xvf crictl-v1.18.0-linux-amd64.tar.gz
      tar -xvf containerd-1.3.6-linux-amd64.tar.gz -C containerd
      sudo tar -xvf cni-plugins-linux-amd64-v0.8.6.tgz -C /opt/cni/bin/
      sudo mv runc.amd64 runc
      chmod +x crictl kubectl kube-proxy kubelet runc
      sudo mv crictl kubectl kube-proxy kubelet runc /usr/local/bin/
      sudo mv containerd/bin/* /bin/
    ```

1. Configure CNI Networking

    Set Pod_Cidr: `10.200.0.0/24`, `10.200.1.0/24`, `10.200.2.0/24`

    `bridge` network configuration file:

    ```
    cat <<EOF | sudo tee /etc/cni/net.d/10-bridge.conf
    > {
    >     "cniVersion": "0.3.1",
    >     "name": "bridge",
    >     "type": "bridge",
    >     "bridge": "cnio0",
    >     "isGateway": true,
    >     "ipMasq": true,
    >     "ipam": {
    >         "type": "host-local",
    >         "ranges": [
    >           [{"subnet": "${POD_CIDR}"}]
    >         ],
    >         "routes": [{"dst": "0.0.0.0/0"}]
    >     }
    > }
    > EOF
    ```

    `loopback` network configuration file:

    ```
    cat <<EOF | sudo tee /etc/cni/net.d/99-loopback.conf
    {
        "cniVersion": "0.3.1",
        "name": "lo",
        "type": "loopback"
    }
    EOF
    ```

1. Configure containerd `/etc/systemd/system/containerd.service`

1. Configure the Kubelet

    1. `/var/lib/kubelet/kubelet-config.yaml` with `KubeletConfiguration`
    1. `/etc/systemd/system/kubelet.service`

1. Configure the Kubernetes Proxy

    1. `/var/lib/kube-proxy/kube-proxy-config.yaml` with `KubeProxyConfiguration`
    1. `/etc/systemd/system/kube-proxy.service`

1. Start worker services (`containerd`, `kubelet`, `kube-proxy`)

    ```
    sudo systemctl daemon-reload
    sudo systemctl enable containerd kubelet kube-proxy
    sudo systemctl start containerd kubelet kube-proxy
    ```

1. Verification

    ```
    gcloud compute ssh controller-0 \
      --command "kubectl get nodes --kubeconfig admin.kubeconfig"
    NAME       STATUS   ROLES    AGE   VERSION
    worker-0   Ready    <none>   49s   v1.18.6
    worker-1   Ready    <none>   49s   v1.18.6
    worker-2   Ready    <none>   49s   v1.18.6
    ```

## [10 Configuring kubectl for Remote Access](https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/10-configuring-kubectl.md)


1. set kubeconfig for `admin` user using pem

    ```
      KUBERNETES_PUBLIC_ADDRESS=$(gcloud compute addresses describe kubernetes-the-hard-way \
        --region $(gcloud config get-value compute/region) \
        --format 'value(address)')

      kubectl config set-cluster kubernetes-the-hard-way \
        --certificate-authority=ca.pem \
        --embed-certs=true \
        --server=https://${KUBERNETES_PUBLIC_ADDRESS}:6443

      kubectl config set-credentials admin \
        --client-certificate=admin.pem \
        --client-key=admin-key.pem

      kubectl config set-context kubernetes-the-hard-way \
        --cluster=kubernetes-the-hard-way \
        --user=admin

      kubectl config use-context kubernetes-the-hard-way
    ```

1. verification

    ```
    kubectl get componentstatuses

    NAME                 STATUS    MESSAGE             ERROR
    scheduler            Healthy   ok
    controller-manager   Healthy   ok
    etcd-1               Healthy   {"health":"true"}
    etcd-2               Healthy   {"health":"true"}
    etcd-0               Healthy   {"health":"true"}
    ```

    ```
    kubectl get nodes

    NAME       STATUS   ROLES    AGE     VERSION
    worker-0   Ready    <none>   4m11s   v1.18.6
    worker-1   Ready    <none>   4m11s   v1.18.6
    worker-2   Ready    <none>   4m11s   v1.18.6
    ```

## [11 Provisioning Pod Network Routes](https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/11-pod-network-routes.md)

Pods scheduled to a node receive an IP address from the node's Pod CIDR range. At this point pods can not communicate with other pods running on different nodes due to missing network routes.

1. Routing table

    ```
    for instance in worker-0 worker-1 worker-2; do
      gcloud compute instances describe ${instance} \
        --format 'value[separator=" "](networkInterfaces[0].networkIP,metadata.items[0].value)'
    done
    10.240.0.20 10.200.0.0/24
    10.240.0.21 10.200.1.0/24
    10.240.0.22 10.200.2.0/24
    ```

1. Route

    ```
    for i in 0 1 2; do
      gcloud compute routes create kubernetes-route-10-200-${i}-0-24 \
        --network kubernetes-the-hard-way \
        --next-hop-address 10.240.0.2${i} \
        --destination-range 10.200.${i}.0/24
    done
    ```

1. Check

    ```
    gcloud compute routes list --filter "network: kubernetes-the-hard-way"

    NAME                            NETWORK                  DEST_RANGE     NEXT_HOP                  PRIORITY
    default-route-46c0079825e1f0e3  kubernetes-the-hard-way  0.0.0.0/0      default-internet-gateway  1000
    default-route-66490234317b35fc  kubernetes-the-hard-way  10.240.0.0/24  kubernetes-the-hard-way   0
    kubernetes-route-10-200-0-0-24  kubernetes-the-hard-way  10.200.0.0/24  10.240.0.20               1000
    kubernetes-route-10-200-1-0-24  kubernetes-the-hard-way  10.200.1.0/24  10.240.0.21               1000
    kubernetes-route-10-200-2-0-24  kubernetes-the-hard-way  10.200.2.0/24  10.240.0.22               1000
    ```
## [12 Deploying the DNS Cluster Add-on](https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/12-dns-addon.md)

In this lab you will deploy the DNS add-on which provides DNS based service discovery, backed by CoreDNS, to applications running inside the Kubernetes cluster.

```
kubectl apply -f https://storage.googleapis.com/kubernetes-the-hard-way/coredns-1.7.0.yaml
```

```
kubectl run busybox --image=busybox:1.28 --command -- sleep 3600
POD_NAME=$(kubectl get pods -l run=busybox -o jsonpath="{.items[0].metadata.name}")
kubectl exec -ti $POD_NAME -- nslookup kubernetes

Server:    10.32.0.10
Address 1: 10.32.0.10 kube-dns.kube-system.svc.cluster.local

Name:      kubernetes
Address 1: 10.32.0.1 kubernetes.default.svc.cluster.local
```


## [13 Smoke Test](https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/13-smoke-test.md)

allow remote access to the `nginx` node port:

```
gcloud compute firewall-rules create kubernetes-the-hard-way-allow-nginx-service \
  --allow=tcp:${NODE_PORT} \
  --network kubernetes-the-hard-way
```

## [Cleaning Up](https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/14-cleanup.md)

# Spent time: 3 hours!

2020-05-03 10:21 Created a project `naka-kubernetes-the-hard-way`
2020-05-03 13:12 Completed Clean Up
