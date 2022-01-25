# Kubernetes Operator Study Journey
# 1. Use existing operators
- [prometheus-operator](../prometheus-operator)
- [postgres-operator](../postgres-operator)
- [strimzi](../strimzi)
- [argocd](../argocd): [appcontroller.go](https://github.com/argoproj/argo-cd/blob/9025318adf367ae8f13b1a99e5c19344402b7bb9/controller/appcontroller.go)
# 2. Understand what is Kubernetes operator.

> A Kubernetes operator is an application-specific controller that extends the functionality of the Kubernetes API to create, configure, and manage instances of complex applications on behalf of a Kubernetes user.

From https://www.redhat.com/en/topics/containers/what-is-a-kubernetes-operator

> Operators are software extensions to Kubernetes that make use of custom resources to manage applications and their components. Operators follow Kubernetes principles, notably the control loop.

From [Operator Pattern in Kubernetes documentation](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)

**Operator Pattern**

![](https://github.com/cncf/tag-app-delivery/blob/eece8f7307f2970f46f100f51932db106db46968/operator-wg/whitepaper/img/02_1_operator_pattern.png?raw=true)

Three components:
1. The application or infrastructure that we want to manage.
1. A domain specific language that enables the users to specify the desired state of the application in a declarative way.
1. A controller that runs continuously:
    - reads and is aware of the state.
    - runs actions againsst the application in an automated way.
    - report the state of the application in a declarative way.

From [](https://github.com/cncf/tag-app-delivery/blob/main/operator-wg/whitepaper/Operator-WhitePaper_v1-0.md)

**Kubernetes Operator**

![](https://github.com/cncf/tag-app-delivery/blob/main/operator-wg/whitepaper/img/02_2_operator.png?raw=true)

Example:

![](diagram.drawio.svg)

1. Kubernetes Controller components.
1. How Kubernetes Controlloer works.
1. Custom Resource.


**Operator vs. Controller**

> - Controller（Custom Controller）:Custom Resourceの管理を行うController。Control Loop（Reconciliation Loop）を実行するコンポーネント
> - Operator: CRDとCustom Controllerのセット。etcd operatorやmysql operatorなどのように、特定のソフトウェアの管理を自動化するためのソフトウェア

From [実践入門Kubernetesカスタムコントローラーへの道](https://www.amazon.co.jp/%E5%AE%9F%E8%B7%B5%E5%85%A5%E9%96%80-Kubernetes%E3%82%AB%E3%82%B9%E3%82%BF%E3%83%A0%E3%82%B3%E3%83%B3%E3%83%88%E3%83%AD%E3%83%BC%E3%83%A9%E3%83%BC%E3%81%B8%E3%81%AE%E9%81%93-%E6%8A%80%E8%A1%93%E3%81%AE%E6%B3%89%E3%82%B7%E3%83%AA%E3%83%BC%E3%82%BA%EF%BC%88NextPublishing%EF%BC%89-%E7%A3%AF-%E8%B3%A2%E5%A4%A7-ebook/dp/B0851QCR81/ref=sr_1_1?__mk_ja_JP=%E3%82%AB%E3%82%BF%E3%82%AB%E3%83%8A&keywords=custom+controller&qid=1636178868&sr=8-1)

> - Controllers can act on core resources such as deployments or services, which are typically part of the Kubernetes controller manager in the control plane, or can watch and manipulate user-defined custom resources.
> - Operators are controllers that encode some operational knowledge, such as application lifecycle management, along with the custom resources defined in Chapter 4.

From [Programming Kubernetes](https://www.oreilly.com/library/view/programming-kubernetes/9781492047094/ch01.html)

> Technically, there is no difference between a typical controller and an operator. Often the difference referred to is the operational knowledge that is included in the operator. Therefore, a controller is the implementation, and the operator is the pattern of using custom controllers with CRDs and automation is what is looking to be achieved with this. As a result, a controller which spins up a pod when a custom resource is created, and the pod gets destroyed afterwards can be described as a simple controller. If the controller has operational knowledge like how to upgrade or remediate from errors, it is an operator.

From [](https://github.com/cncf/tag-app-delivery/blob/main/operator-wg/whitepaper/Operator-WhitePaper_v1-0.md#kubernetes-controllers)
# 3. Create a sample operator following a tutorial

There are several ways to create an operator. Personally, I would recommend starting from `operator-sdk` or `kubebuilder`. I studied them in the following order.

1. [operator-sdk](https://sdk.operatorframework.io/)
    - [go-based](https://sdk.operatorframework.io/docs/building-operators/golang/quickstart/): https://github.com/nakamasato/memcached-operator
    - [helm-based](https://sdk.operatorframework.io/docs/building-operators/helm/quickstart/): https://github.com/nakamasato/nginx-operator
    - [ansible-based](https://sdk.operatorframework.io/docs/building-operators/ansible/quickstart/): https://github.com/nakamasato/memcached-operator-with-ansible
1. [kubebuilder](https://book.kubebuilder.io/)
    - [Tutorial: Building CronJob](https://book.kubebuilder.io/cronjob-tutorial/cronjob-tutorial.html)
    - [つくって学ぶKubebuilder](https://zoetrope.github.io/kubebuilder-training/)
1. [sample-controller](https://github.com/kubernetes/sample-controller): https://github.com/nakamasato/foo-controller-kubebuilder

Other tools:
- [KUDO (Kubernetes Universal Declarative Operator)](https://kudo.dev/)

# 4. Understand more detail about each component

Overview of the event flow:

![](https://github.com/kubernetes/sample-controller/blob/master/docs/images/client-go-controller-interaction.jpeg?raw=true)

From https://github.com/kubernetes/sample-controller/blob/master/docs/images/client-go-controller-interaction.jpeg
1. **Reflector** lists and watches kube-apiserver and adds an object to **Delta FIFO queue**.
1. [Informer](informer) pops objects from the Delta queue one by one. (creted by a factory e.g. **SharedInformerFactory**)
    1. Informer adds an object to **indexer**
    1. Indexer stores the object and key to the **thread safe store**. (called in-memory-cache)
    1. **event handler** handles the object and enqueue object key to **workqueue**
1. **Reconciler** in the controller gets key and process the item.
1. Inside the reconciliation loop, **HandleObject** get object using the key from **indexer reference**

If you use `kubebuilder` or `operator-sdk`, you don't really need to worry about these details. However, if you understand them, it would be helpful when debugging or implementing complicated logic.

Reference
- https://adevjoe.com/post/client-go-informer/
- https://github.com/kubernetes/sample-controller/blob/master/docs/controller-client-go.md
- [kube-controller-manager源码分析（三）之 Informer机制](https://www.huweihuang.com/kubernetes-notes/code-analysis/kube-controller-manager/sharedIndexInformer.html)
- [A deep dive into Kubernetes controllers](https://engineering.bitnami.com/articles/a-deep-dive-into-kubernetes-controllers.html)
- [Introducing Operators: Putting Operational Knowledge into Software](https://web.archive.org/web/20170129131616/https://coreos.com/blog/introducing-operators.html)
- [Best practices for building Kubernetes Operators and stateful apps](https://cloud.google.com/blog/products/containers-kubernetes/best-practices-for-building-kubernetes-operators-and-stateful-apps)
# 4. Create your own operator

After creating a sample operator, you should have deeper understanding of controller. Now you can think about what kind of problem that you want to resolve by utilizing operator pattern.

To clarify a problem to resolve with a new operator, you can reference existing operators:

|operator|role|
|---|---|
|prometheus-operator|Manage Prometheus and configuration|

- [How can I have separate logic for Create, Update, and Delete events? When reconciling an object can I access its previous state?](https://sdk.operatorframework.io/docs/faqs/#how-can-i-have-separate-logic-for-create-update-and-delete-events-when-reconciling-an-object-can-i-access-its-previous-state) -> You should not have separate logic. Instead design your reconciler to be idempotent.
    - [Q: How do I have different logic in my reconciler for different types of events (e.g. create, update, delete)? in controller-runtime](https://github.com/kubernetes-sigs/controller-runtime/blob/master/FAQ.md#q-how-do-i-have-different-logic-in-my-reconciler-for-different-types-of-events-eg-create-update-delete)
- [Owners and Dependents](https://kubernetes.io/docs/concepts/overview/working-with-objects/owners-dependents/)

Tips:
- Finalizer
- Reconciliation Loop
    - [operator-sdk] Based on the return value of Reconcile() the reconcile Request may be requeued and the loop may be triggered again: ([Building a Go-based Memcached Operator using the Operator SDK](https://docs.openshift.com/container-platform/4.1/applications/operator_sdk/osdk-getting-started.html#building-memcached-operator-using-osdk_osdk-getting-started))
        ```go
        // Reconcile successful - don't requeue
        return reconcile.Result{}, nil
        // Reconcile failed due to error - requeue
        return reconcile.Result{}, err
        // Requeue for any reason other than error
        return reconcile.Result{Requeue: true}, nil
        ```
    - https://github.com/operator-framework/operator-sdk/issues/4209#issuecomment-729916367
- Testing
    - **KUbernetes Testing TooL (kuttl)**: https://kuttl.dev/ KUTTL is built to support some kubernetes integration test scenarios and is most valuable as an end-to-end (e2e) test harness.
    - **Ginkgo** (A Golang BDD Testing Framework): https://onsi.github.io/ginkgo/
    - **Gomega** (Ginkgo's preferred matcher library): https://onsi.github.io/gomega/
    - **kubetest2**: https://github.com/kubernetes-sigs/kubetest2: Kubetest2 is the framework for launching and running end-to-end tests on Kubernetes. It is intended to be the next significant iteration of kubetest.
- Managing Errors
    - https://cloud.redhat.com/blog/kubernetes-operators-best-practices
    1. Return the error in the status of the object.
    1. Generate an event describing the error.
# 5. Tools
- https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/controller/controllerutil
- https://github.com/spf13/cobra: a library for creating powerful modern CLI applications & a program to generate applications and command files.
    - Cobra is used in many Go projects such as Kubernetes, Hugo, and Github CLI to name a few. This list contains a more extensive list of projects using Cobra.

# 6. Study Golang for better codes

- [golang-standanrds/project-layout](https://github.com/golang-standards/project-layout)
- [Learn Go with tests](https://quii.gitbook.io/learn-go-with-tests/)
- [GoとDependency Injectionの現在](https://note.com/timakin/n/nc95d66a75b3d)
- [Go Blog](https://go.dev/blog)
- [Gopher Reading List](https://github.com/enocom/gopher-reading-list)
