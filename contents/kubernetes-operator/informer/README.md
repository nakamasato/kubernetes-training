# informer

***Informer*** monitors the changes of target resource. An informer is created for each of the target resources if you need to handle multiple resources (e.g. podInformer, deploymentInformer). Informer stores the

## Overview
1. Initialize clientset with `.kube/config`
1. Create an informer factory with the following line.
    ```go
    informerFactory := informers.NewSharedInformerFactory(kubeClient, time.Second*30)
    ```
    The second argument specifies ***ResyncPeriod***, which defines the interval of resync (*The resync operation consists of delivering to the handler an update notification for every object in the informer's local cache*). For more detail, please read [NewSharedInformer](https://pkg.go.dev/k8s.io/client-go@v0.23.1/tools/cache#NewSharedInformer)
1. Create an informer for Pods, which watches Pod's changes.
    ```go
    podInformer := informerFactory.Core().V1().Pods()
    ```
1. Add event handlers (`AddFunc`, `UpdateFunc`, and `DeleteFunc`) to the pod informer.
    ```go
    podInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    handleAdd,
			UpdateFunc: handleUpdate,
			DeleteFunc: handleDelete,
		},
	)
    ```

    `handleAdd`, `handleUpdate`, and `handleDelete` define custom logic for each event. In this example, just print `"handleXXX is called"`
1. Create a stop channel and start the factory.
    ```go
    ch := make(chan struct{}) // stop channel
	informerFactory.Start(ch)
    ```
1. Wait until the cache is synced.
    ```go
    cacheSynced := podInformer.Informer().HasSynced
	if ok := cache.WaitForCacheSync(ch, cacheSynced); !ok {
		log.Printf("cache is not synced")
	}
	log.Println("cache is synced")
    ```
1. Run `run` function every 10 seconds
    ```go
    go wait.Until(run, time.Second*10, ch)
	<-ch
    ```

## Run and check
1. Run
    ```
    go run informer.go
    ```

1. All Pods are synced in the cache.

    ```
    2021/12/21 09:05:08 handleAdd is called for Pod (key: local-path-storage/local-path-provisioner-547f784dff-lhwfk)
    2021/12/21 09:05:08 handleAdd is called for Pod (key: kube-system/kube-scheduler-kind-control-plane)
    2021/12/21 09:05:08 handleAdd is called for Pod (key: kube-system/etcd-kind-control-plane)
    2021/12/21 09:05:08 handleAdd is called for Pod (key: kube-system/kube-apiserver-kind-control-plane)
    2021/12/21 09:05:08 handleAdd is called for Pod (key: kube-system/kindnet-nzc7p)
    2021/12/21 09:05:08 handleAdd is called for Pod (key: kube-system/coredns-558bd4d5db-b4wjg)
    2021/12/21 09:05:08 handleAdd is called for Pod (key: kube-system/kube-controller-manager-kind-control-plane)
    2021/12/21 09:05:08 handleAdd is called for Pod (key: kube-system/kube-proxy-vrcbc)
    2021/12/21 09:05:08 handleAdd is called for Pod (key: kube-system/coredns-558bd4d5db-8q78s)
    2021/12/21 09:05:08 handleAdd is called for Pod (key: default/foo-sample-688594b488-782kw)
    2021/12/21 09:05:08 cache is synced
    2021/12/21 09:05:08 run
    ```
1. Create a `Pod` with name `nginx`.
    ```
    kubectl run nginx --image=nginx
    ```
1. Handlers are called by the events of the created `Pod`.
    ```
    2021/12/21 09:05:18 run
    2021/12/21 09:05:20 handleAdd is called for Pod (key: default/nginx)
    2021/12/21 09:05:20 handleUpdate is called for Pod (key: default/nginx)
    2021/12/21 09:05:20 handleUpdate is called for Pod (key: default/nginx)
    ```
1. Delete the `Pod`
    ```
    kubectl delete po nginx
    ```
1. Handlers are called by the events of the Pod deletion.
    ```
    2021/12/21 09:05:29 handleUpdate is called for Pod (key: default/nginx)
    2021/12/21 09:05:30 handleUpdate is called for Pod (key: default/nginx)
    2021/12/21 09:05:31 handleUpdate is called for Pod (key: default/nginx)
    2021/12/21 09:05:31 handleUpdate is called for Pod (key: default/nginx)
    2021/12/21 09:05:31 handlDelete is called for Pod (key: default/nginx)
    ```
1. `run` function is called every 10 seconds.
    ```
    2021/12/21 09:26:08 run
    2021/12/21 09:26:18 run
    2021/12/21 09:26:28 run
    ```
1. The cached is resynced every 30 seconds.

    ```
    2021/12/21 09:27:08 handleUpdate is called for Pod (key: local-path-storage/local-path-provisioner-547f784dff-lhwfk)
    2021/12/21 09:27:08 handleUpdate is called for Pod (key: kube-system/kube-apiserver-kind-control-plane)
    2021/12/21 09:27:08 handleUpdate is called for Pod (key: kube-system/coredns-558bd4d5db-b4wjg)
    2021/12/21 09:27:08 handleUpdate is called for Pod (key: kube-system/kube-controller-manager-kind-control-plane)
    2021/12/21 09:27:08 handleUpdate is called for Pod (key: kube-system/coredns-558bd4d5db-8q78s)
    2021/12/21 09:27:08 handleUpdate is called for Pod (key: default/foo-sample-688594b488-782kw)
    2021/12/21 09:27:08 handleUpdate is called for Pod (key: kube-system/kube-scheduler-kind-control-plane)
    2021/12/21 09:27:08 handleUpdate is called for Pod (key: kube-system/etcd-kind-control-plane)
    2021/12/21 09:27:08 handleUpdate is called for Pod (key: kube-system/kindnet-nzc7p)
    2021/12/21 09:27:08 handleUpdate is called for Pod (key: kube-system/kube-proxy-vrcbc)
    ```

# reference
- https://adevjoe.com/post/client-go-informer/
- https://www.huweihuang.com/kubernetes-notes/code-analysis/kube-controller-manager/sharedIndexInformer.html
- https://yangxikun.com/kubernetes/2020/03/05/informer-lister.html
