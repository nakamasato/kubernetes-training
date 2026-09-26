# controller-runtime

Kubernetes controller を構築するライブラリ。ここでは [go.mod](../../../go.mod) の **v0.25.1** を対象に、実行できるサンプルと各コンポーネントの関係を説明する。client-go は **v0.37.1**。バージョンを更新するときは両者の互換性を確認する。

## Overview

1. `ctrl.NewManager` が Cluster、Cache、Client、Scheme などの共通依存を準備する。
1. Reconciler に `mgr.GetClient()` などを明示的に渡す。
1. Builder の `For` / `Owns` / `Watches` で監視対象とイベントの対応付けを登録する。
1. `mgr.Start(ctrl.SetupSignalHandler())` でキャッシュ・Controller などを開始する。
1. Source が Informer の通知を Handler に渡し、Handler が処理対象キーをキューへ追加する。
1. Controller が同期を待って Reconciler を呼び出す。Reconciler は現在の状態を読み、必要な変更だけを書き込む。

イベントは集約されるため、Reconcile はイベント履歴ではなく現在の状態から判断する。何度呼ばれても同じ状態に収束するように実装する。

## Components

- [manager](manager): ライフサイクルと共有依存
- [cluster](cluster): クラスタへの接続、Scheme、RESTMapper
- [client](client): API の読み書きとキャッシュ
- [cache](cache): Informer の作成・同期
- [builder](builder): Controller と監視対象の構築
- [controller](controller): キュー、ワーカー、再試行
- [reconciler](reconciler): 調整ロジック
- [source](source): イベントの供給元
- [handler](handler): イベントから処理対象キーへの変換
- [log](log)、[leaderelection](leaderelection)、[webhook](webhook)
- [inject からの移行](inject): 削除された依存注入 API の置き換え

図は controller-runtime v0.25.1 / client-go v0.37.1 に合わせたもの。概念上の流れと実装の詳細を区別し、詳細は各ページのバージョン固定の参照先を確認する。

## Examples

リポジトリルートで実行する。Go のバージョンは go.mod の `go` / `toolchain` に従う。

```sh
go test ./contents/kubernetes-operator/...
go run ./contents/kubernetes-operator/controller-runtime/reconciler
go run ./contents/kubernetes-operator/controller-runtime/log
```

クラスタを使う例は有効な kubeconfig と対象リソースへの権限が必要。`KUBECONFIG` または `~/.kube/config` を使う。Manager / Cache / Source の例は Ctrl+C で終了する。

- [example-controller](example-controller): ReplicaSet に所有 Pod 数のラベルを付ける
- [manager](manager): Pod と Deployment を監視
- [cache](cache): キャッシュから Pod を取得
- [source](source): Pod / MySQLUser のイベントを観察
- [webhook](webhook): TLS サーバーで AdmissionReview を処理

参照: [v0.25.1 API](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.25.1)、[互換性](https://github.com/kubernetes-sigs/controller-runtime/tree/v0.25.1#compatibility)。

## 図

![controller-runtime の処理と構成](diagram.drawio.svg)

SVG には diagrams.net / draw.io の編集データを埋め込んでいる。再生成の定義は [generate_operator_diagrams.py](../../../scripts/generate_operator_diagrams.py) にあり、図を変更するときはこの定義を更新する。リポジトリルートで再生成・同期確認できる。

```sh
python3 scripts/generate_operator_diagrams.py
python3 scripts/generate_operator_diagrams.py --check
```
