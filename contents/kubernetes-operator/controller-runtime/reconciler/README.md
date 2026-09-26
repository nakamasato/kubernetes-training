# Reconciler

Reconciler は現在の状態を読み、望ましい状態に近づける処理を実装する。通常の `reconcile.Reconciler` は `TypedReconciler[Request]` の別名で、次のメソッドを持つ。

```go
Reconcile(context.Context, reconcile.Request) (reconcile.Result, error)
```

`Request` は `types.NamespacedName` を持つ。作成・更新・削除といったイベント種別や古いオブジェクトは含まない。現在の状態を取得し、繰り返し実行しても不要な変更を行わないようにする。

| 戻り値 | 動作 |
|---|---|
| `reconcile.Result{}, nil` | 今回の処理を完了。新しいイベントでは再度呼ばれる |
| `reconcile.Result{RequeueAfter: time.Minute}, nil` | 指定時間後に再実行 |
| `reconcile.Result{}, err` | 通常のエラーは rate limit 付きで再試行 |
| `reconcile.Result{}, reconcile.TerminalError(err)` | エラーを記録するが、そのエラーを理由に再試行しない |

`Result.Requeue` は deprecated。明示的な待ち時間には `RequeueAfter` を使う。エラーを返す場合、Result は無視される。

## Run

[main.go](main.go) は struct と `reconcile.Func` の二通りで実装する。struct は大文字の `Reconcile` メソッドが必要で、コンパイル時にインターフェースへの適合を確認している。クラスタは不要。

```sh
go run ./contents/kubernetes-operator/controller-runtime/reconciler
```

参照: [Reconcile API](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.25.1/pkg/reconcile)、[実際の Controller](../example-controller)。
