# Handler

Handler は Create / Update / Delete / Generic イベントからキューへ処理対象を追加する。通常の `EventHandler` は `TypedEventHandler[client.Object, reconcile.Request]` の別名。

| Handler | enqueue する対象 |
|---|---|
| `EnqueueRequestForObject` | イベントのオブジェクト自身 |
| `EnqueueRequestForOwner` | ownerReferences が示す親 |
| `EnqueueRequestsFromMapFunc` | 任意の対応付け関数が返すキー |
| `TypedFuncs[Object, Request]` | 独自の型とコールバックで指定する対象 |

Owner Handler には依存を明示的に渡す。

```go
h := handler.EnqueueRequestForOwner(
    mgr.GetScheme(), mgr.GetRESTMapper(), &appsv1.ReplicaSet{},
    handler.OnlyControllerOwner(),
)
```

`source.Kind` に `*corev1.Pod` を渡すなら、Handler も `TypedEnqueueRequestForObject[*corev1.Pod]` のように object 型を揃える。複数のリソースを同じ Handler で処理する例は [Source](../source) を参照。

Handler は軽い処理に留め、API の読み書きや再試行が必要な処理は Reconciler に委ねる。Predicate は Handler より前にイベントを絞り込む。

参照: [Handler API](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.25.1/pkg/handler)、[実装](https://github.com/kubernetes-sigs/controller-runtime/tree/v0.25.1/pkg/handler)。
