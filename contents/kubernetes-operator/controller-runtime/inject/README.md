# inject から明示的な依存の受け渡しへ

以前の `pkg/runtime/inject`、`InjectClient`、`InjectCache`、`SetFields` による自動注入は現在の API では使わない。このディレクトリは旧記事からのリンクと移行案内のために残している。

Reconciler はフィールドを定義し、構築時に必要な依存を渡す。

```go
r := &ReplicaSetReconciler{Client: mgr.GetClient()}
err := ctrl.NewControllerManagedBy(mgr).
    For(&appsv1.ReplicaSet{}).
    Owns(&corev1.Pod{}).
    Complete(r)
```

`client.Client` の埋め込みだけでは値は設定されない。Source も `source.Kind(mgr.GetCache(), object, eventHandler)` のように構築時に依存を渡す。Owner 用 Handler には `mgr.GetScheme()` と `mgr.GetRESTMapper()` を渡す。

- [動作する Controller](../example-controller/example-controller.go)
- [Source](../source)
- [依存注入の削除に関する変更](https://github.com/kubernetes-sigs/controller-runtime/pull/2120)
