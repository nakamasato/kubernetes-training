# kube-apiserver

## Run kube-apiserver in local

1. Build Kubernetes binary (ref: [Build Kubernetes](../README.md#build-kubernetes)).
1. Run `etcd`. (ref: [etcd](../etcd/))

    version:
    ```
    etcd --version
    etcd Version: 3.5.2
    Git SHA: 99018a77b
    Go Version: go1.17.6
    Go OS/Arch: darwin/amd64
    ```

    start:

    ```
    etcd
    ```
1. Create certificates for `service-account`.
    ```
    openssl version
    LibreSSL 2.8.3
    ```

    ```
    openssl genrsa -out service-account-key.pem 4096
    ```

    <details>

    ```
    Generating RSA private key, 4096 bit long modulus
    .....................................................................++
    ................................................++
    e is 65537 (0x10001)
    ```

    </details>

    ```
    openssl req -new -x509 -days 365 -key service-account-key.pem -sha256 -out service-account.pem
    ```

    <details>

    ```
    You are about to be asked to enter information that will be incorporated
    into your certificate request.
    What you are about to enter is what is called a Distinguished Name or a DN.
    There are quite a few fields but you can leave some blank
    For some fields there will be a default value,
    If you enter '.', the field will be left blank.
    -----
    Country Name (2 letter code) []:JP
    State or Province Name (full name) []:Tokyo
    Locality Name (eg, city) []:Kita
    Organization Name (eg, company) []:Test
    Organizational Unit Name (eg, section) []:Test
    Common Name (eg, fully qualified host name) []:Test
    Email Address []:masatonaka1989@gmail.com
    ```

    </details>

1. Create certificate for tls.

    1. Generate a `ca.key` with 2048bit:
        ```
        openssl genrsa -out ca.key 2048
        ```
    1. According to the `ca.key` generate a `ca.crt` (use -days to set the certificate effective time):
        ```
        openssl req -x509 -new -nodes -key ca.key -subj "/CN=127.0.0.1" -days 10000 -out ca.crt
        ```
    1. `server.key`
        ```
        openssl genrsa -out server.key 2048
        ```
    1. `csr.conf`
    1. generate certificate signing request (`server.csr`)
        ```
        openssl req -new -key server.key -out server.csr -config csr.conf
        ```
    1. generate server certificate `server.crt` using `ca.key`, `ca.crt` and `server.csr`.
        ```
        openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key \
        -CAcreateserial -out server.crt -days 10000 \
        -extensions v3_ext -extfile csr.conf
        ```


1. Run the built binary.

    ```
    PATH_TO_KUBERNETES_DIR=~/repos/kubernetes/kubernetes
    ```

    ```
    ${PATH_TO_KUBERNETES_DIR}/_output/bin/kube-apiserver --version
    Kubernetes v1.23.4-dirty
    ```

    ```
    ${PATH_TO_KUBERNETES_DIR}/_output/bin/kube-apiserver --etcd-servers http://localhost:2379 \
    --service-account-key-file=service-account-key.pem \
    --service-account-signing-key-file=service-account-key.pem \
    --service-account-issuer=api \
    --tls-cert-file=server.crt \
    --tls-private-key-file=server.key \
    --client-ca-file=ca.crt
    ```

1. Configure `admin.kubeconfig`.

    (I'm too lazy to generate crt and key for kubectl. So used the same one as server here.)

    ```
    kubectl config set-cluster local-apiserver \
    --certificate-authority=ca.crt \
    --embed-certs=true \
    --server=https://127.0.0.1:6443 \
    --kubeconfig=admin.kubeconfig

    kubectl config set-credentials admin \
    --client-certificate=server.crt \
    --client-key=server.key \
    --embed-certs=true \
    --kubeconfig=admin.kubeconfig

    kubectl config set-context default \
    --cluster=local-apiserver \
    --user=admin \
    --kubeconfig=admin.kubeconfig

    kubectl config use-context default --kubeconfig=admin.kubeconfig
    ```
1. Check component status. (only `etcd` is healthy.)
    ```
    kubectl get componentstatuses --kubeconfig admin.kubeconfig
    Warning: v1 ComponentStatus is deprecated in v1.19+
    NAME                 STATUS      MESSAGE                                                                                        ERROR
    controller-manager   Unhealthy   Get "https://127.0.0.1:10257/healthz": dial tcp 127.0.0.1:10257: connect: connection refused
    scheduler            Unhealthy   Get "https://127.0.0.1:10259/healthz": dial tcp 127.0.0.1:10259: connect: connection refused
    etcd-0               Healthy     {"health":"true","reason":""}
    ```

## Errors

###  mkdir /var/run/kubernetes: permission denied

```
E0302 06:40:09.767084   37385 run.go:74] "command failed" err="error creating self-signed certificates: mkdir /var/run/kubernetes: permission denied"
```

Run
```
sudo mkdir /var/run/kubernetes
chown -R `whoami` /var/run/kubernetes
```

### service-account-issuer is a required flag, --service-account-signing-key-file and --service-account-issuer are required flags

```
E0302 07:14:46.234431   79468 run.go:74] "command failed" err="[service-account-issuer is a required flag, --service-account-signing-key-file and --service-account-issuer are required flags]"
```

`BoundServiceAccountTokenVolume` is now GA from 1.22.

## References
- https://kubernetes.io/docs/reference/command-line-tools-reference/feature-gates/
- https://github.com/kelseyhightower/kubernetes-the-hard-way/issues/626
- https://headtonirvana.hatenablog.com/entry/2021/10/11/Kubernetes_The_Hard_Way_On_VirtualBox_6%E6%97%A5%E7%9B%AE
- https://kubernetes.io/docs/tasks/administer-cluster/certificates/
- https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/ca96371e4d2d2176e8b2c3f5b656b5d92973479e/docs/05-kubernetes-configuration-files.md
