# [Knative](https://knative.dev/docs)

**abstracting away the complex details and enabling developers to focus on what matters**

Enterprise-grade Serverless on your own terms.
Kubernetes-based platform to deploy and manage modern serverless workloads.

## Components

1. Serving
1. Eventing

## Concept

**[Serving Resources](https://knative.dev/docs/serving/#serving-resources)**

1. ***Service**: The `service.serving.knative.dev` resource automatically manages the whole lifecycle of your workload. It controls the creation of other objects to ensure that your app has a route, a configuration, and a new revision for each update of the service. Service can be defined to always route traffic to the latest revision or to a pinned revision.*
1. ***Route**: The `route.serving.knative.dev` resource maps a network endpoint to one or more revisions. You can manage the traffic in several ways, including fractional traffic and named routes.*
1. ***Configuration**: The `configuration.serving.knative.dev` resource maintains the desired state for your deployment. It provides a clean separation between code and configuration and follows the Twelve-Factor App methodology. Modifying a configuration creates a new revision.*
1. ***Revision**: The `revision.serving.knative.dev` resource is a point-in-time snapshot of the code and configuration for each modification made to the workload. Revisions are immutable objects and can be retained for as long as useful. Knative Serving Revisions can be automatically scaled up and down according to incoming traffic. See Configuring the Autoscaler for more information.*

![](https://raw.githubusercontent.com/knative/serving/main/docs/spec/images/object_model.png)

## [Getting Started](https://knative.dev/docs/getting-started/)

1. Install Knative CLI `kn`.

    For Mac:

    ```
    brew install kn
    ```

1. Install the Knative "Quickstart" plugin. (`kn` CLI's plugin.)

    <details><summary>※ Somehow brew fails.</summary>

    ```
    brew install knative-sandbox/kn-plugins/quickstart
    ```

    </details>


    -> Follow [this](https://github.com/knative-sandbox/kn-plugin-quickstart/blob/release-1.0/README.md#installation)

    1. Download the binary
    1. Move it to `/usr/local/bin`
        ```
        mv ~/Downloads/kn-quickstart-darwin-amd64 /usr/local/bin/kn-quickstart
        chmod +x /usr/local/bin/kn-quickstart
        ```
    1. Check plugin.
        ```
        kn plugin list
        - kn-quickstart : /usr/local/bin/kn-quickstart
        ```

    ```
    kn quickstart kind
    ```

    <details><summary>pods</summary>

    ```
    kubectl get po -A
    NAMESPACE            NAME                                            READY   STATUS    RESTARTS   AGE
    knative-eventing     eventing-controller-58875c5478-bhhcd            1/1     Running   0          56s
    knative-eventing     eventing-webhook-5968f79978-cswb4               1/1     Running   0          56s
    knative-eventing     imc-controller-86cd7b7857-ndx69                 1/1     Running   0          41s
    knative-eventing     imc-dispatcher-7fcb4b5d8c-rtjkz                 1/1     Running   0          41s
    knative-eventing     mt-broker-controller-8d979648f-nv4w4            1/1     Running   0          27s
    knative-eventing     mt-broker-filter-574dc4457f-x4lch               1/1     Running   0          27s
    knative-eventing     mt-broker-ingress-5ddd6f8b5d-llj4g              1/1     Running   0          27s
    knative-serving      activator-85bd4ddcbb-sdgkc                      1/1     Running   0          2m14s
    knative-serving      autoscaler-84fcdc5449-52rzj                     1/1     Running   0          2m14s
    knative-serving      controller-6fd5bb86df-zctsm                     1/1     Running   0          2m14s
    knative-serving      domain-mapping-74d5d688bd-k9stl                 1/1     Running   0          2m14s
    knative-serving      domainmapping-webhook-8484d5fd8b-6jb9n          1/1     Running   0          2m14s
    knative-serving      net-kourier-controller-66bc9d6697-czjbr         1/1     Running   0          97s
    knative-serving      webhook-97c648588-rwb6n                         1/1     Running   0          2m13s
    kourier-system       3scale-kourier-gateway-58856c6cc7-czpnz         1/1     Running   0          97s
    kube-system          coredns-78fcd69978-qdg6m                        1/1     Running   0          2m28s
    kube-system          coredns-78fcd69978-qj7xh                        1/1     Running   0          2m28s
    kube-system          etcd-knative-control-plane                      1/1     Running   0          2m43s
    kube-system          kindnet-h85z6                                   1/1     Running   0          2m28s
    kube-system          kube-apiserver-knative-control-plane            1/1     Running   0          2m42s
    kube-system          kube-controller-manager-knative-control-plane   1/1     Running   0          2m43s
    kube-system          kube-proxy-b9p56                                1/1     Running   0          2m28s
    kube-system          kube-scheduler-knative-control-plane            1/1     Running   0          2m43s
    local-path-storage   local-path-provisioner-85494db59d-99nt7         1/1     Running   0          2m28s
    ```

    </details>

1. Apply hello.yaml
    ```
    kubectl apply -f hello.yaml
    ```
1. Check

    Knative service:

    ```
    kn service list
    ```

    ```
    NAME    URL                                LATEST   AGE   CONDITIONS   READY     REASON
    hello   http://hello.default.example.com            36s   0 OK / 3     Unknown   RevisionMissing : Configuration "hello" is waiting for a Revision to become ready.
    ```

    ```
    kn service list
    NAME    URL                                LATEST        AGE   CONDITIONS   READY   REASON
    hello   http://hello.default.example.com   hello-world   37m   3 OK / 3     True
    ```

    Check `Hello World`.

    ```
    curl http://hello.default.127.0.0.1.sslip.io

    Hello World!
    ```

1. [Scaling to Zero](https://knative.dev/docs/getting-started/first-autoscale/#scaling-to-zero)

    > It may take up to 2 minutes for your Pods to scale down. Pinging your service again will reset this timer.

    ```
    kubectl get pod -l serving.knative.dev/service=hello
    No resources found in default namespace.
    ```

1. [Basics of Traffic Splitting](https://knative.dev/docs/getting-started/first-traffic-split/#basics-of-traffic-splitting)

    1. Creating a new Revision

        ```
        kn service update hello \
        --env TARGET=Knative \
        --revision-name=knative
        ```

        ```
        Updating Service 'hello' in namespace 'default':

          0.024s The Configuration is still working to reflect the latest desired specification.
          3.181s Traffic is not yet migrated to the latest revision.
          3.210s Ingress has not yet been reconciled.
          3.242s Waiting for load balancer to be ready
          3.419s Ready to serve.

        Service 'hello' updated to latest revision 'hello-knative' is available at URL:
        http://hello.default.127.0.0.1.sslip.io
        ```

    1. Check

        ```
        curl http://hello.default.127.0.0.1.sslip.io
        Hello Knative!
        ```

    1. Splitting Traffic

        ```
        kn revisions list

        NAME            SERVICE   TRAFFIC   TAGS   GENERATION   AGE     CONDITIONS   READY           REASON
        hello-knative   hello     100%             2            2m12s   3 OK / 4     True
        hello-world     hello                      1            120m    3 OK / 4     True
        ```

        ```
        kn service update hello \
        --traffic hello-world=50 \
        --traffic @latest=50
        ```

        ```
        kn revisions list

        NAME            SERVICE   TRAFFIC   TAGS   GENERATION   AGE    CONDITIONS   READY   REASON
        hello-knative   hello     50%              2            3m1s   3 OK / 4     True
        hello-world     hello     50%              1            121m   3 OK / 4     True
        ```

    1. Check

        ```
        ± curl http://hello.default.127.0.0.1.sslip.io
        Hello Knative!

        ± curl http://hello.default.127.0.0.1.sslip.io
        Hello Knative!

        ± curl http://hello.default.127.0.0.1.sslip.io
        Hello Knative!

        ± curl http://hello.default.127.0.0.1.sslip.io
        Hello World!
        ```

1. Clean up

    ```
    kn service delete hello
    ```
