# EKS Workshop

https://eksworkshop.com/

# Praparation

## eksctl

```
curl --silent --location "https://github.com/weaveworks/eksctl/releases/latest/download/eksctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp

sudo mv -v /tmp/eksctl /usr/local/bin
```

```
eksctl version
0.19.0
```

## kms

create

```
aws kms create-alias --alias-name alias/eksworkshop --target-key-id $(aws kms create-key --query KeyMetadata.Arn --output text)
```

```
export MASTER_ARN=$(aws kms describe-key --key-id alias/eksworkshop --query KeyMetadata.Arn --output text)


[20-05-15 2:38:15] nakamasato at Masatos-MacBook-Pro in ~/Code/MasatoNaka/kubernetes-training/eksworkshop on master ✘
± echo $MASTER_ARN
arn:aws:kms:ap-northeast-1:135493629466:key/3b5de2cf-b5b9-4d78-8786-b4eb52948204
```

```
export AWS_REGION=ap-northeast-1
```


# Create EKS cluster

```
cat << EOF > eksworkshop.yaml
---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: eksworkshop-eksctl
  region: ${AWS_REGION}

managedNodeGroups:
- name: nodegroup
  desiredCapacity: 3
  iam:
    withAddonPolicies:
      albIngress: true

secretsEncryption:
  keyARN: ${MASTER_ARN}
EOF
```

```
eksctl create cluster -f eksworkshop.yaml
[ℹ]  eksctl version 0.19.0
[ℹ]  using region ap-northeast-1
[ℹ]  setting availability zones to [ap-northeast-1d ap-northeast-1a ap-northeast-1c]
[ℹ]  subnets for ap-northeast-1d - public:192.168.0.0/19 private:192.168.96.0/19
[ℹ]  subnets for ap-northeast-1a - public:192.168.32.0/19 private:192.168.128.0/19
[ℹ]  subnets for ap-northeast-1c - public:192.168.64.0/19 private:192.168.160.0/19
[ℹ]  using Kubernetes version 1.15
[ℹ]  creating EKS cluster "eksworkshop-eksctl" in "ap-northeast-1" region with managed nodes
[ℹ]  1 nodegroup (nodegroup) was included (based on the include/exclude rules)
[ℹ]  will create a CloudFormation stack for cluster itself and 0 nodegroup stack(s)
[ℹ]  will create a CloudFormation stack for cluster itself and 1 managed nodegroup stack(s)
[ℹ]  if you encounter any issues, check CloudFormation console or try 'eksctl utils describe-stacks --region=ap-northeast-1 --cluster=eksworkshop-eksctl'
[ℹ]  CloudWatch logging will not be enabled for cluster "eksworkshop-eksctl" in "ap-northeast-1"
[ℹ]  you can enable it with 'eksctl utils update-cluster-logging --region=ap-northeast-1 --cluster=eksworkshop-eksctl'
[ℹ]  Kubernetes API endpoint access will use default of {publicAccess=true, privateAccess=false} for cluster "eksworkshop-eksctl" in "ap-northeast-1"
[ℹ]  2 sequential tasks: { create cluster control plane "eksworkshop-eksctl", create managed nodegroup "nodegroup" }
[ℹ]  building cluster stack "eksctl-eksworkshop-eksctl-cluster"
[ℹ]  deploying stack "eksctl-eksworkshop-eksctl-cluster"
[ℹ]  building managed nodegroup stack "eksctl-eksworkshop-eksctl-nodegroup-nodegroup"
[ℹ]  deploying stack "eksctl-eksworkshop-eksctl-nodegroup-nodegroup"
[ℹ]  waiting for the control plane availability...
[✔]  saved kubeconfig as "/Users/nakamasato/.kube/config"
[ℹ]  no tasks
[✔]  all EKS cluster resources for "eksworkshop-eksctl" have been created
[ℹ]  nodegroup "nodegroup" has 3 node(s)
[ℹ]  node "ip-192-168-21-173.ap-northeast-1.compute.internal" is ready
[ℹ]  node "ip-192-168-60-193.ap-northeast-1.compute.internal" is ready
[ℹ]  node "ip-192-168-75-58.ap-northeast-1.compute.internal" is ready
[ℹ]  waiting for at least 3 node(s) to become ready in "nodegroup"
[ℹ]  nodegroup "nodegroup" has 3 node(s)
[ℹ]  node "ip-192-168-21-173.ap-northeast-1.compute.internal" is ready
[ℹ]  node "ip-192-168-60-193.ap-northeast-1.compute.internal" is ready
[ℹ]  node "ip-192-168-75-58.ap-northeast-1.compute.internal" is ready
[ℹ]  kubectl command should work with "/Users/nakamasato/.kube/config", try 'kubectl get nodes'
[✔]  EKS cluster "eksworkshop-eksctl" in "ap-northeast-1" region is ready
```

check

```
kubectl get nodes # if we see our 3 nodes, we know we have authenticated correctly
NAME                                                STATUS   ROLES    AGE   VERSION
ip-192-168-21-173.ap-northeast-1.compute.internal   Ready    <none>   15m   v1.15.11-eks-af3caf
ip-192-168-60-193.ap-northeast-1.compute.internal   Ready    <none>   15m   v1.15.11-eks-af3caf
ip-192-168-75-58.ap-northeast-1.compute.internal    Ready    <none>   14m   v1.15.11-eks-af3caf
```

# Deploy

## [Kubernetes Dashboard](https://eksworkshop.com/beginner/040_dashboard/)

```
kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.0.0/aio/deploy/recommended.yaml
```

```
kubectl proxy --port=8080 --address=0.0.0.0 --disable-filter=true &
```

open: http://localhost:8080/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/#/overview?namespace=default

Get token and insert it

```
aws eks get-token --cluster-name eksworkshop-eksctl | jq -r '.status.token'
```

## [Example Microservice](https://eksworkshop.com/beginner/050_deploy/)

- Backend
- Frontend

## [Helm](https://eksworkshop.com/beginner/060_helm/)

# Clean-up

## eks

```
eksctl delete cluster -f eksworkshop.yaml
```

## kms

delete

```
eksctl delete cluster -f eksworkshop.yaml
```
