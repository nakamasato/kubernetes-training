# Cache

controller-runtime の Cache は client-go の Informer を管理し、`client.Reader` として `Get` / `List` を提供する。書き込みは Client を通じて API サーバーへ行う。

## 起動と同期

1. `cache.New(config, cache.Options{Scheme: scheme})` で作成する。HTTP client と RESTMapper は省略すると作成される。
1. `GetInformer(ctx, &corev1.Pod{})` で監視対象を登録する。
1. `AddEventHandler` で Handler を登録し、返されたエラーを処理する。
1. 別 goroutine で `Start(ctx)` を実行する。
1. `WaitForCacheSync(ctx)` が成功してから `Get` / `List` する。
1. context をキャンセルして停止する。

開始前に空の Pod を `Get` して Informer 作成を促す旧手順は使わない。開始前の読み取りは `ErrCacheNotStarted` になり、開始後でも初回 Informer 同期を待つ場合がある。

Manager を使うと通常はこれらのライフサイクルを Manager / Controller が管理する。`cache.Options.DefaultNamespaces`、`ByObject` などで監視範囲を制限できる。`ReaderFailOnMissingInformer` を有効にすると未登録の型の読み取りで Informer を暗黙に作らずエラーにできる。

## Run

リポジトリルートで実行する。サンプルは全 namespace の Pod を監視し、`default/nginx` をキャッシュから取得する。

```sh
kubectl run nginx --image=nginx -n default
go run ./contents/kubernetes-operator/controller-runtime/cache
```

別ターミナルで Pod を変更・削除すると `OnUpdate` / `OnDelete` が表示される。Ctrl+C で終了する。

```sh
kubectl delete pod nginx -n default --ignore-not-found
```

参照: [main.go](main.go)、[Cache API](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.25.1/pkg/cache)、[Informer Cache](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.25.1/pkg/cache/informer_cache.go)。
