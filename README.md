# kubernetes-training

# Version

- kubernetes: 1.15
- Helm: 3.2.1
- Traefik: 2.2

# Contents

- [helm](helm)
- [amazon-eks-workshop](amazon-eks-workshop)
- [ingress-nginx-controller](ingress-nginx-controller)
- [kubernetes-the-hard-way](kubernetes-the-hard-way)
- [traefik](traefik)

# All steps

## nakamasato/terraform

1. Create VPC (3 mins)

    ```
    cd aws/global/vpc
    terraform init && terraform apply
    ```

1. Create EKS cluster (12 mins)

    ```
    cd ../../eks-on-vpc/dev
    terraform init && terraform apply
    ```

    ```
    terraform output kubeconfig >> ~/.kube/config
    kubectl cluster-info
    ```

    ```
    terraform output config_map_aws_auth > config_map_aws_auth.yaml
    kubectl apply -f config_map_aws_auth.yaml
    ```

## nakamasato/k8s-deploy-test

1. Deploy Applications

    - [x] ArgoCD

        https://github.com/nakamasato/k8s-deploy-test/tree/master/argocd

        ```
        kubectl create namespace argocd
        kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
        ```

        Get login password:

        ```
        kubectl get pod -n argocd | grep argocd-server | awk '{print $1}'
        ```

        ```
        kubectl port-forward svc/argocd-server -n argocd 8080:443
        ```

        open https://localhost:8080/ user: `admin`, pass: `argocd-server-xxxxxxx-xxxxx`

    - [ ] Prometheus
    - [ ] Grafana
