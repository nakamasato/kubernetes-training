# Kubernetes the hard way

https://github.com/kelseyhightower/kubernetes-the-hard-way

## Resources

## [03 Compute Resources](https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/03-compute-resources.md)

### networks subnets

```
gcloud compute networks subnets create kubernetes \
  --network kubernetes-the-hard-way \
  --range 10.240.0.0/24
```

### firewall-rules

```
gcloud compute firewall-rules create kubernetes-the-hard-way-allow-internal \
  --allow tcp,udp,icmp \
  --network kubernetes-the-hard-way \
  --source-ranges 10.240.0.0/24,10.200.0.0/16
```

```
gcloud compute firewall-rules create kubernetes-the-hard-way-allow-external \
  --allow tcp:22,tcp:6443,icmp \
  --network kubernetes-the-hard-way \
  --source-ranges 0.0.0.0/0
```

### address

```
gcloud compute addresses create kubernetes-the-hard-way \
  --region $(gcloud config get-value compute/region)
```

### instances

controller

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

worker

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

## [Provisioning a CA and Generating TLS Certificates](https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/04-certificate-authority.md)


generate key for Kubernetes components and distribute the client and server certificates

## [Generating Kubernetes Configuration Files for Authentication](https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/05-kubernetes-configuration-files.md)

Kubeconfigs: enable Kubernetes clients to locate and authenticate to the Kubernetes API servers

- kubelet
- kube-proxy
- kube-controller-manager
- kube-scheduler
- admin

Distribute Kuberenetes configuration file
- workers
- controllers

## [Generating the Data Encryption Config and Key](https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/06-data-encryption-keys.md)


## [Bootstrapping the etcd Cluster](https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/07-bootstrapping-etcd.md)

Run the command to set up etcd on the controller-0, controller-1, and controller-2 (I used iterm2's toggle input to run exactly the same command in all of the server at the same time)

## [Bootstrapping the Kubernetes Control Plane](https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/08-bootstrapping-kubernetes-controllers.md)

```
gcloud compute http-health-checks create kubernetes \
    --description "Kubernetes Health Check" \
    --host "kubernetes.default.svc.cluster.local" \
    --request-path "/healthz"
```

```
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

## [Bootstrapping the Kubernetes Worker Nodes](https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/09-bootstrapping-kubernetes-workers.md)


## [Configuring kubectl for Remote Access](https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/10-configuring-kubectl.md)

## [Provisioning Pod Network Routes](https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/11-pod-network-routes.md)

```
for i in 0 1 2; do
  gcloud compute routes create kubernetes-route-10-200-${i}-0-24 \
    --network kubernetes-the-hard-way \
    --next-hop-address 10.240.0.2${i} \
    --destination-range 10.200.${i}.0/24
done
```

## [Deploying the DNS Cluster Add-on](https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/12-dns-addon.md)

## [Smoke Test](https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/13-smoke-test.md)

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
