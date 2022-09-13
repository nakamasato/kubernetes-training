# client-go

Version: [v0.25.0](https://github.com/kubernetes/client-go/releases/tag/v0.25.0)

1. [clientset](clientset): a set of clients to access Kubernetes API
1. [indexer](indexer): An indexed in-memory key-value store for objects
1. [informer](informer)
    1. indexer
    1. reflector
    1. ListerWatcher
1. [lister](lister)
    1. indexer
1. [watcher](watcher)
1. [reflector](reflector): watches a specified resource with **listerwatcher** and reflects all changes to the configured store (FIFO).
1. [listerwatcher](listerwatcher): list and watch the API server. used in **reflector**.
