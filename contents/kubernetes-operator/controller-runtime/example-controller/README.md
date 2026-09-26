# Example Controller

ReplicaSet と、その ReplicaSet を controller owner に持つ Pod を監視する。ReplicaSet に `pod-count` ラベルで所有 Pod 数を記録する。

[example-controller.go](example-controller.go) の処理:

1. ReplicaSet を取得する。削除済み（NotFound）なら正常終了する。
1. `spec.selector`（matchExpressions を含む）で同じ namespace の Pod を絞る。
1. controller owner の UID が一致する Pod を数える。他の ReplicaSet や owner のない Pod は数えない。
1. ラベルが未設定なら map を初期化し、値が変わる場合だけ merge patch する。

Pod の終了処理中でも、キャッシュ上に残っている所有 Pod は数に含む。これはラベル更新の教材で、ReplicaSet 本来の Pod 作成・削除は Kubernetes の ReplicaSet controller が担う。

## Run

リポジトリルートで実行する。ReplicaSet の get/list/watch/patch と Pod の list/watch 権限を用意する。

```sh
go run ./contents/kubernetes-operator/controller-runtime/example-controller
```

別ターミナルで Deployment を作成する。

```sh
kubectl create deployment test --replicas=3 --image=nginx
kubectl get rs -l app=test -L pod-count
kubectl scale deployment test --replicas=1
kubectl get rs -l app=test -L pod-count
```

キャッシュが追いつくと `pod-count` が 3、スケール後は 1 になる。書き込みの直後は読み取りが古い場合がある。エラーは Controller によって再試行される。

```sh
kubectl delete deployment test
```

Ctrl+C で Controller を停止する。クラスタ不要の回帰テスト:

```sh
go test ./contents/kubernetes-operator/controller-runtime/example-controller
```
