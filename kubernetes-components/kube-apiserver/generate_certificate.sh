# service account
openssl genrsa -out service-account-key.pem 4096
openssl req -new -x509 -days 365 -key service-account-key.pem -subj "/CN=test" -sha256 -out service-account.pem

# api-server
openssl genrsa -out ca.key 2048
openssl req -x509 -new -nodes -key ca.key -subj "/CN=test" -days 10000 -out ca.crt
openssl genrsa -out server.key 2048
openssl req -new -key server.key -out server.csr -config csr.conf
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key \
    -CAcreateserial -out server.crt -days 10000 \
    -extensions v3_ext -extfile csr.conf

# kubeconfig
kubectl config set-cluster local-apiserver \
--certificate-authority=ca.crt \
--embed-certs=true \
--server=https://127.0.0.1:6443 \
--kubeconfig=kubeconfig

kubectl config set-credentials admin \
--client-certificate=server.crt \
--client-key=server.key \
--embed-certs=true \
--kubeconfig=kubeconfig

kubectl config set-context default \
--cluster=local-apiserver \
--user=admin \
--kubeconfig=kubeconfig

kubectl config use-context default --kubeconfig=kubeconfig
