# client-go

![](diagram.drawio.svg)

Version: [v0.37.1](https://github.com/kubernetes/client-go/releases/tag/v0.37.1)

1. [clientset](clientset): a set of clients to access Kubernetes API
1. [indexer](indexer): An indexed in-memory key-value store for objects
1. [informer](informer)
    1. indexer
    1. reflector
    1. ListerWatcher
1. [lister](lister)
    1. indexer
1. [workqueue](workqueue): manages controller work items and retries.
1. [reflector](reflector): watches a specified resource with **listerwatcher** and reflects all changes to the configured store (FIFO).
1. [listerwatcher](listerwatcher): list and watch the API server. used in **reflector**.

Commands in the component walkthroughs use the versions pinned in [go.mod](../../../go.mod). Run `go test ./contents/kubernetes-operator/client-go/...` from the repository root. Informer/cache queues update stored objects; a controller workqueue separately schedules reconciliation keys.
