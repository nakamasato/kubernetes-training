# Reflector

Reflector は指定された型のリソースを監視し、変更を Store に反映する。通常は Informer が内部で作成する。

基本の流れは List で初期状態を取得し、その resourceVersion から Watch を継続すること。v0.37.1 では設定やサーバー対応状況に応じて WatchList（初期イベント付き Watch）による初期化もあるため、常に独立した List リクエストから始まるとは限らない。

- 初期状態は Store の Replace に反映する。
- Added / Modified / Deleted を Store の Add / Update / Delete に対応させる。
- 接続断やエラーに対して再試行し、必要なら初期状態を取得し直す。
- resourceVersion は不透明な値として扱い、数値として比較・加算しない。
- resync は既知のオブジェクトを再処理する仕組みで、API の再 List と同義ではない。

単独で使う場合は `cache.NewReflector(...)` で作り、`RunWithContext(ctx)` で起動・停止を管理する。通常の Controller は [Informer](../informer) を使うと Indexer、Handler、同期もまとめて管理できる。

参照: [Reflector API](https://pkg.go.dev/k8s.io/client-go@v0.37.1/tools/cache#Reflector)、[実装](https://github.com/kubernetes/client-go/blob/v0.37.1/tools/cache/reflector.go)、[ListerWatcher](../listerwatcher)。
