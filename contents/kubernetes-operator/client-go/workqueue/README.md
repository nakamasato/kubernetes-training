# Workqueue

Controller が処理するキーを保持するキュー。同じキーの重複追加を集約し、処理中に再追加されたキーは `Done` 後にもう一度処理できる。

現在は型付き API を使う。

```go
q := workqueue.NewTypedRateLimitingQueue(
    workqueue.DefaultTypedControllerRateLimiter[types.NamespacedName](),
)
defer q.ShutDown()
q.Add(types.NamespacedName{Namespace: "default", Name: "example"})
```

ワーカーで一つ処理する例（`sync` はそのキーの最新状態を読み調整する関数）:

```go
key, shutdown := q.Get()
if shutdown {
    return false
}
defer q.Done(key)
if err := sync(ctx, key); err != nil {
    q.AddRateLimited(key)
} else {
    q.Forget(key)
}
return true
```

- `Done` は取得した項目の処理完了を通知する。成功・失敗のどちらでも必要。
- `Forget` は失敗回数など rate limiter の状態をリセットする。`Done` の代わりにはならない。
- `AddAfter` は一定時間後、`AddRateLimited` は rate limiter に従って再追加する。
- 空のキューの `Get` は待機する。context をキャンセルするだけでは解除されないので、停止時に `ShutDown` する。

[Source の実行例](../../controller-runtime/source) は型付きキューを使い、context のキャンセルで `ShutDown` する。通常の controller-runtime Controller はキューと再試行を内部で管理するので、Reconciler は Result / error を返す。

参照: [Workqueue API](https://pkg.go.dev/k8s.io/client-go@v0.37.1/util/workqueue)、[公式サンプル](https://github.com/kubernetes/client-go/blob/v0.37.1/examples/workqueue/main.go)。
