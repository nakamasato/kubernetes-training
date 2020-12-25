package main

deny[{"msg": msg, "details": {"missing_prefix": "nodeSelector"}}] {
  input.kind = "Deployment"
  not input.spec.template.spec.nodeSelector
  msg := "you must provide nodeSelector"
}
