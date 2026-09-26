# ListerWatcher

List でリソースの集合とその `metadata.resourceVersion` を取得し、Watch でその時点からの変更を受け取る。List の直後に resourceVersion を指定せず Watch すると、その間の変更を取り逃す可能性がある。

現在は context を渡せる `ListerWatcherWithContext` / `ListWatch.ListWithContext` / `WatchWithContext` を使う。context なしの `List` / `Watch` は deprecated。

```go
list, err := podListWatcher.ListWithContext(ctx, metav1.ListOptions{})
// err を確認してから ListAccessor を呼ぶ。
listMeta, err := meta.ListAccessor(list)
w, err := podListWatcher.WatchWithContext(ctx, metav1.ListOptions{
    ResourceVersion: listMeta.GetResourceVersion(),
})
// err を確認してから Stop を defer する。
defer w.Stop()
```

[main.go](main.go) に完全なエラー処理を含む例がある。`ResultChan()` の終了と `watch.Error` を扱い、Ctrl+C で context をキャンセルする。

この例は一回の List + Watch を観察する教材。接続終了時に終了する。resourceVersion の期限切れや接続断から自動復旧させたい場合は [Reflector](../reflector) / [Informer](../informer) を使う。

## Run

default namespace の Pod の list/watch 権限が必要。

```sh
go run ./contents/kubernetes-operator/client-go/listerwatcher
```

別ターミナルで変更する。

```sh
kubectl run nginx --image=nginx -n default
kubectl annotate pod nginx example=value -n default
kubectl delete pod nginx -n default
```

参照: [ListWatch API](https://pkg.go.dev/k8s.io/client-go@v0.37.1/tools/cache#ListWatch)、[実装](https://github.com/kubernetes/client-go/blob/v0.37.1/tools/cache/listwatch.go)。
