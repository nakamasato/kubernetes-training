# Clientset

Clientset は組み込みリソース向けの型付き client をまとめる。例: `clientset.CoreV1().Pods(namespace)`、`clientset.AppsV1().Deployments(namespace)`。

```go
cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
// err を確認してから client を作る。
clientset, err := kubernetes.NewForConfig(cfg)
// err を確認してから API を呼ぶ。
pods, err := clientset.CoreV1().Pods(metav1.NamespaceAll).List(ctx, metav1.ListOptions{})
```

各段階のエラーを必ず確認する。`List` にはキャンセルやタイムアウトを設定した context を渡す。CRD はこの組み込み Clientset に自動追加されない。生成した client、dynamic client、または controller-runtime の Client を使う。

## Run

リポジトリルートで実行する。全 namespace の Pod を list できる kubeconfig が必要。

```sh
go run ./contents/kubernetes-operator/client-go/clientset
# 接続設定を明示する場合:
go run ./contents/kubernetes-operator/client-go/clientset -kubeconfig /path/to/config
```

[podlist.go](podlist.go) は 30 秒のタイムアウトで Pod 一覧を取得し、namespace/name を表示する。リソースは変更しない。

参照: [Clientset API](https://pkg.go.dev/k8s.io/client-go@v0.37.1/kubernetes)、[client-go の例](https://github.com/kubernetes/client-go/tree/v0.37.1/examples)。
