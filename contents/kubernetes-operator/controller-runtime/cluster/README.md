# Cluster

Cluster は一つの Kubernetes クラスタへの接続と共有依存を保持する。Manager はこの機能を埋め込み、さらに Controller やリーダー選出などのライフサイクルを管理する。

| Getter | 用途 |
|---|---|
| `GetConfig()` | API サーバーへの接続設定 |
| `GetScheme()` | Go の型と GroupVersionKind の対応 |
| `GetRESTMapper()` | Kind と API resource / scope の対応 |
| `GetCache()` | Informer とキャッシュの読み取り |
| `GetClient()` | 通常はキャッシュ経由の読み取り、API への書き込み |
| `GetAPIReader()` | キャッシュを介さない読み取り |

`cluster.New(config, options)` は Scheme、HTTP client、RESTMapper、Cache、Client などを準備する。`Start(ctx)` は Cache を開始する。単独で使う場合は起動と同期を呼び出し側で管理する。

以前の `SetFields` や `NewDelegatingClient` を使う構築手順は廃止されている。現在は `client.Options.Cache` を使ってキャッシュを Client に渡す。独自 RESTMapper を作る場合、`apiutil.NewDynamicRESTMapper` には `rest.HTTPClientFor(config)` で作った HTTP client を渡す。通常は Cluster / Cache のデフォルト設定でよい。

参照: [Cluster API](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.25.1/pkg/cluster)、[構築処理](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.25.1/pkg/cluster/cluster.go)、[Client](../client)。
