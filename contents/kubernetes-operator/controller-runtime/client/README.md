# Client

`client.Client` は Kubernetes API を Go のオブジェクトとして読み書きするためのインターフェース。

| 操作 | Manager のデフォルト Client |
|---|---|
| `Get` / `List` | 通常は Cache から取得 |
| `Create` / `Update` / `Patch` / `Delete` | API サーバーへ送信 |
| `Status().Update` / `Status().Patch` | status サブリソースへ送信 |
| `mgr.GetAPIReader().Get` / `List` | API サーバーから直接取得 |

書き込み直後にキャッシュが更新済みとは限らない。Reconcile は読み取りが遅れても再実行によって収束するようにする。キャッシュ内のフィールド検索には `mgr.GetFieldIndexer().IndexField` で対応する index を登録する。

`client.New(config, client.Options{})` で単独に作った Client は、Cache を指定しなければ直接 API を読む。Manager の Client と同じ挙動とは限らない。現在のキャッシュ設定は `client.Options.Cache` / `client.CacheOptions`（`Reader`、`DisableFor`、`Unstructured`）。旧 `NewDelegatingClient` / `delegatingClient` の構築例は使わない。

## 更新例

```go
before := obj.DeepCopy()
if obj.Labels == nil {
    obj.Labels = map[string]string{}
}
obj.Labels["example"] = "value"
err := c.Patch(ctx, obj, client.MergeFrom(before))
```

`MergeFrom` は変更差分の JSON merge patch を作る。読み取った値との競合検出が必要なら `MergeFromWithOptions(before, client.MergeFromWithOptimisticLock{})` を使う。毎回同じ書き込みを行うと不要なイベントが発生するので、値が変わるときだけ書き込む。

Get の NotFound は削除済みなら正常終了にできる。例: `return ctrl.Result{}, client.IgnoreNotFound(err)`。

参照: [Client API](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.25.1/pkg/client)、[実装](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.25.1/pkg/client/client.go)、[動作する例](../example-controller)。

## 図

![client の処理と構成](diagram.drawio.svg)
