# [Cert Manager](https://cert-manager.io/)

## Install Cert Manager

### Install with yaml

```
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.7.1/cert-manager.yaml
```

### Install with `cmctl`
1. Install `cmctl` following https://cert-manager.io/docs/usage/cmctl/#installation
    ```
    OS=$(go env GOOS); ARCH=$(go env GOARCH); curl -L -o cmctl.tar.gz https://github.com/jetstack/cert-manager/releases/latest/download/cmctl-$OS-$ARCH.tar.gz
    tar xzf cmctl.tar.gz
    sudo mv cmctl /usr/local/bin
    ```
    ```
    cmctl version
    Client Version: util.Version{GitVersion:"v1.6.0", GitCommit:"49914a057b39c887be0974c4657c095bd7724bc7", GitTreeState:"clean", GoVersion:"go1.17.1", Compiler:"gc", Platform:"darwin/amd64"}
    error: could not detect the cert-manager version: the cert-manager CRDs are not yet installed on the Kubernetes API server
    ```
    ```
    cmctl help
    ```
1. Install cert manager with `cmctl` (https://cert-manager.io/docs/installation/cmctl/)
    ```
    cmctl x install
    ```
    <details><summary>output</summary>
    ```
    cmctl x install
    Creating the cert-manager CRDs
    creating 6 resource(s)
    Clearing discovery cache
    beginning wait for 6 resources with timeout of 1m0s
    creating 1 resource(s)
    creating 38 resource(s)
    beginning wait for 38 resources with timeout of 5m0s
    Deployment is not ready: default/cert-manager-cainjector. 0 out of 1 expected pods are ready
    Deployment is not ready: default/cert-manager-cainjector. 0 out of 1 expected pods are ready
    Deployment is not ready: default/cert-manager-cainjector. 0 out of 1 expected pods are ready
    Deployment is not ready: default/cert-manager-cainjector. 0 out of 1 expected pods are ready
    Deployment is not ready: default/cert-manager. 0 out of 1 expected pods are ready
    Deployment is not ready: default/cert-manager. 0 out of 1 expected pods are ready
    Deployment is not ready: default/cert-manager. 0 out of 1 expected pods are ready
    Deployment is not ready: default/cert-manager. 0 out of 1 expected pods are ready
    Deployment is not ready: default/cert-manager-webhook. 0 out of 1 expected pods are ready
    Deployment is not ready: default/cert-manager-webhook. 0 out of 1 expected pods are ready
    Deployment is not ready: default/cert-manager-webhook. 0 out of 1 expected pods are ready
    Deployment is not ready: default/cert-manager-webhook. 0 out of 1 expected pods are ready
    Deployment is not ready: default/cert-manager-webhook. 0 out of 1 expected pods are ready
    Deployment is not ready: default/cert-manager-webhook. 0 out of 1 expected pods are ready
    Deployment is not ready: default/cert-manager-webhook. 0 out of 1 expected pods are ready
    Starting delete for "cert-manager-startupapicheck" ServiceAccount
    serviceaccounts "cert-manager-startupapicheck" not found
    creating 1 resource(s)
    Starting delete for "cert-manager-startupapicheck:create-cert" Role
    roles.rbac.authorization.k8s.io "cert-manager-startupapicheck:create-cert" not found
    creating 1 resource(s)
    Starting delete for "cert-manager-startupapicheck:create-cert" RoleBinding
    rolebindings.rbac.authorization.k8s.io "cert-manager-startupapicheck:create-cert" not found
    creating 1 resource(s)
    Starting delete for "cert-manager-startupapicheck" Job
    jobs.batch "cert-manager-startupapicheck" not found
    creating 1 resource(s)
    Watching for changes to Job cert-manager-startupapicheck with timeout of 5m0s
    Add/Modify event for cert-manager-startupapicheck: ADDED
    cert-manager-startupapicheck: Jobs active: 0, jobs failed: 0, jobs succeeded: 0
    Add/Modify event for cert-manager-startupapicheck: MODIFIED
    cert-manager-startupapicheck: Jobs active: 1, jobs failed: 0, jobs succeeded: 0
    Add/Modify event for cert-manager-startupapicheck: MODIFIED
    Starting delete for "cert-manager-startupapicheck" ServiceAccount
    Starting delete for "cert-manager-startupapicheck:create-cert" Role
    Starting delete for "cert-manager-startupapicheck:create-cert" RoleBinding
    Starting delete for "cert-manager-startupapicheck" Job
    NAME: cert-manager
    LAST DEPLOYED: Thu Oct 28 10:29:34 2021
    NAMESPACE: default
    STATUS: deployed
    REVISION: 1
    DESCRIPTION: Cert-manager was installed using the cert-manager CLI
    NOTES:
    cert-manager v1.6.0 has been deployed successfully!
    In order to begin issuing certificates, you will need to set up a ClusterIssuer
    or Issuer resource (for example, by creating a 'letsencrypt-staging' issuer).
    More information on the different types of issuers and how to configure them
    can be found in our documentation:
    https://cert-manager.io/docs/configuration/
    For information on how to configure cert-manager to automatically provision
    Certificates for Ingress resources, take a look at the `ingress-shim`
    documentation:
    https://cert-manager.io/docs/usage/ingress/
    ```
    ```
    kubectl get po
    NAME                                       READY   STATUS    RESTARTS   AGE
    cert-manager-6c576bddcf-jl92p              1/1     Running   0          110s
    cert-manager-cainjector-669c966b86-mx956   1/1     Running   0          110s
    cert-manager-webhook-7d6cf57d55-9qvxk      1/1     Running   0          110s
    ```
    </details>
    <details><summary>Uninstall</summary>
    ```
    cmctl x install --dry-run > cert-manager.custom.yaml
    kubectl delete -f cert-manager.custom.yaml
    ```
    </details>

##
