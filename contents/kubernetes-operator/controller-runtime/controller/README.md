# Controller

Controller は Source、処理キュー、ワーカー、Reconciler を接続する。

1. `controller.New(name, mgr, options)` または [Builder](../builder) で作り、Manager に登録する。
1. `Watch(src)` で Source を登録する。Handler / Predicate は Source の構築時に渡す。
1. Manager が Controller を起動する。Controller は Source を開始し、SyncingSource の同期を待つ。
1. Handler がキーを enqueue すると、ワーカーが Reconciler を呼ぶ。
1. 結果に応じて完了、遅延再試行、エラーの rate limit 付き再試行を行う。

```go
src := source.Kind(mgr.GetCache(), &corev1.Pod{},
    &handler.TypedEnqueueRequestForObject[*corev1.Pod]{})
err := c.Watch(src)
```

`MaxConcurrentReconciles` で並列数を設定できる。同じキーのイベントはキューで集約されるため、イベントごとに必ず Reconcile が一回ずつ動くとは限らない。

通常は Manager から起動する。`NewUnmanaged` で作った Controller は呼び出し側が `Start(ctx)` と停止を管理する。

参照: [Controller API](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.25.1/pkg/controller)、[内部実装](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.25.1/pkg/internal/controller/controller.go)、[Reconciler](../reconciler)。

## 図

![controller の処理と構成](diagram.drawio.svg)
