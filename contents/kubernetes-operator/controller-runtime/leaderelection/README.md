# Leader election

複数の Manager を起動し、通常はリーダーになった一つだけが Controller を実行するための仕組み。controller-runtime は client-go の Lease ベースのリーダー選出を使う。

```go
mgr, err := ctrl.NewManager(config, ctrl.Options{
    LeaderElection:          true,
    LeaderElectionID:        "example-controller.example.com",
    LeaderElectionNamespace: "default",
})
```

同じ Controller のレプリカでは同じ ID / namespace を使い、別 Controller と衝突しない ID を選ぶ。Lease を作成・取得・更新する RBAC が必要。ローカル実行でも namespace を明示すると設定を追いやすい。

Webhook などリーダー選出を必要としない Runnable は各レプリカで起動する。`NeedLeaderElection()` を持つ Runnable はその戻り値で実行条件を指定できる。

リーダーは Lease を定期的に更新し、更新できなくなると処理を停止する。リーダー選出だけで外部システムへの厳密な排他（fencing）を保証するものではないため、調整ロジックは繰り返し実行に耐えるようにする。

参照: [Manager options](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.25.1/pkg/manager#Options)、[controller-runtime の実装](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.25.1/pkg/leaderelection/leader_election.go)、[client-go leader election](https://pkg.go.dev/k8s.io/client-go@v0.37.1/tools/leaderelection)。
