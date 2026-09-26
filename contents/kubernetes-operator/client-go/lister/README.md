# [Lister](https://pkg.go.dev/k8s.io/client-go@v0.37.1/listers/apps/v1#DeploymentLister)

## Overview

***Lister*** is any object that knows how to perform an initial list.

Interface:

```go
type Lister interface {
    // List should return a list type object; the Items field will be extracted, and the
    // ResourceVersion field will be used to start the watch in the right place.
    List(options metav1.ListOptions) (runtime.Object, error)
}
```

## Usage

1. Prepare Indexer (dependency). Read [indexer](../indexer) for more details.

    ```go
    indexer := cache.NewIndexer(
        cache.MetaNamespaceKeyFunc,
        cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
    )
    ```

    Add object to indexer with/without labels (which will be used as selector when listing with Lister below).

    ```go
    // Add deployment with label
    err := indexer.Add(getDeployment("deployment-with-label", map[string]string{"watch": "true"}))
    if err != nil {
        t.Errorf("indexer.Add failed %v\n", err)
    }
    // Add deployment without label
    err = indexer.Add(getDeployment("deployment-without-label", map[string]string{}))
    if err != nil {
        t.Errorf("indexer.Add failed %v\n", err)
    }
    ```

1. Create a Lister for a target resource.

    For a built-in resource:

    ```go
    import (
        appsv1lister "k8s.io/client-go/listers/apps/v1"
    )
    ```
    ```go
    deploymentLister := appsv1lister.NewDeploymentLister(indexer)
    ```

    For a custom resource: Need to generate lister with [code-generator](https://github.com/kubernetes/code-generator)


1. List object with selector.

    ```go
    selector := labels.SelectorFromSet(labels.Set{"watch": "true"})
    ret, _ := deploymentLister.List(selector)

    for _, deploy := range ret {
        fmt.Println(deploy.Name)
    }
    ```

    Objects are got from the indexers.

## Cache ownership and indexes

A Lister reads the Indexer rather than the API server. Wait for informer synchronization before using it, and `DeepCopy()` returned objects before modification. The sample's label selector filters candidates; adding a custom `labels` index does not make a generated Lister use that index automatically. A namespace-scoped List can use the namespace index.

Run `go run ./contents/kubernetes-operator/client-go/lister` or `go test ./contents/kubernetes-operator/client-go/lister` from the repository root without a cluster. See the [generated DeploymentLister](https://github.com/kubernetes/client-go/blob/v0.37.1/listers/apps/v1/deployment.go) and [ListAll / ListAllByNamespace](https://github.com/kubernetes/client-go/blob/v0.37.1/tools/cache/listers.go).
