# Source

Source はイベントの供給元。現在は Handler と Predicate を構築時に渡す。`Source` は `TypedSource[reconcile.Request]` の別名で、`Start(ctx, queue)` を持つ。同期が必要な `TypedSyncingSource` は `WaitForSync(ctx)` も持つ。

```go
src := source.Kind(mgr.GetCache(), &corev1.Pod{},
    &handler.TypedEnqueueRequestForObject[*corev1.Pod]{})
err := c.Watch(src)
```

- `Kind(cache, object, handler, predicates...)`: Cache から Informer を取得してイベントを受け取る。
- `TypedKind[Object, Request](...)`: 独自の comparable な request 型をキューに入れる。
- `Channel(events, handler, ...)`: クラスタ外からの GenericEvent を受け取る。イベントを channel に送る処理は利用側で実装する。
- `Informer`: 既存 Informer を使う低レベルの Source。

通常は Controller が Source の開始と同期を管理する。このサンプルは仕組みを観察するために手動で Cache、Source、キューを開始する。

## Run

リポジトリルートから実行する。Pod と MySQLUser の list/watch 権限が必要。MySQLUser の Go 型の Scheme 登録と、クラスタへの CRD インストールは別々に必要。

```sh
kubectl apply -f https://raw.githubusercontent.com/nakamasato/mysql-operator/v0.4.3/config/crd/bases/mysql.nakamasato.com_mysqlusers.yaml
go run ./contents/kubernetes-operator/controller-runtime/source
```

別ターミナルでイベントを発生させる。

```sh
kubectl run nginx --image=nginx -n default
kubectl delete pod nginx -n default
```

`WorkQueueItem` はイベント種別、Go の型、namespace、名前を記録する。同名の別 namespace / 別 kind のオブジェクトを区別できる。これは観察用で、通常の Controller では namespace/name のキーから最新状態を読む。

同期は 30 秒でタイムアウトする。CRD 未登録や権限不足ならログを確認する。Ctrl+C で Cache を停止し、キューを ShutDown して `Get()` の待機を解除する。サンプルはキュー取得後に `Forget` と `Done` を呼ぶ。

必要に応じて学習用クラスタの CRD を削除する（その CRD のリソースも削除される）。

```sh
kubectl delete crd mysqlusers.mysql.nakamasato.com
```

参照: [main.go](main.go)、[Source API](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.25.1/pkg/source)、[Kind の実装](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.25.1/pkg/internal/source/kind.go)。

## 図

![source の処理と構成](diagram.drawio.svg)

![Source から処理キューへのイベントの流れ](dataflow.drawio.svg)
