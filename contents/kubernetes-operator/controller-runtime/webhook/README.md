# Webhook

AdmissionReview を受け取る TLS サーバーの例。

- `/mutating`: オブジェクトに `access=granted` と `reason=not so secret` の annotation を追加する。annotations が未設定でも動作し、既存の別 annotation と未知のフィールドを保持する。
- `/validating`: サンプルとして常に拒否する。

`admission.Webhook` と `admission.HandlerFunc` を使い、`admission.PatchResponseFromRaw` で差分を作る。`webhook.NewServer(webhook.Options{...})` を `ctrl.Options.WebhookServer` に渡し、`mgr.GetWebhookServer().Register` で登録する。

## Run

リポジトリルートで開発用証明書を作成する（本番用ではない）。有効な kubeconfig も必要。

```sh
cert_dir=$(mktemp -d)
openssl req -x509 -newkey rsa:2048 -nodes -days 1 \
  -keyout "$cert_dir/tls.key" -out "$cert_dir/tls.crt" -subj '/CN=localhost'
go run ./contents/kubernetes-operator/controller-runtime/webhook -cert-dir "$cert_dir"
```

別ターミナルで直接 AdmissionReview を送信する。

```sh
curl -sk https://localhost:8443/mutating -H 'Content-Type: application/json' -d '{"apiVersion":"admission.k8s.io/v1","kind":"AdmissionReview","request":{"uid":"demo","kind":{"group":"","version":"v1","kind":"Pod"},"resource":{"group":"","version":"v1","resource":"pods"},"operation":"CREATE","object":{"apiVersion":"v1","kind":"Pod","metadata":{"name":"demo"},"spec":{"containers":[{"name":"nginx","image":"nginx"}]}}}}'
```

`response.allowed: true` と base64 の JSONPatch が返る。Ctrl+C で終了し、同じターミナルで `rm -r "$cert_dir"` する。

このサーバー起動だけでは Kubernetes の admission に登録されない。クラスタで使うには到達可能な Service、証明書、CA bundle、対象ルールを持つ MutatingWebhookConfiguration / ValidatingWebhookConfiguration が別途必要。この例の拒否 Handler を全リソースに登録しないこと。

## 型付きの検証・デフォルト設定

v0.25.1 では `admission.Validator[T]` / `admission.Defaulter[T]` を使える。各メソッドは context と対象オブジェクトを受け取り、Validator は warnings と error を返す。`admission.WithValidator` / `WithDefaulter` で Handler を構成できる。旧オブジェクト自身の `ValidateCreate()` / `Default()` の説明とは異なる。

参照: [main.go](main.go)、[Admission API](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.25.1/pkg/webhook/admission)、[Server API](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.25.1/pkg/webhook)。
