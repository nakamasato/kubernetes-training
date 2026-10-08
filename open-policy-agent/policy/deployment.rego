package main

violation[msg] {
  input.kind = "Deployment"
  not input.spec.template.spec.securityContext.runAsNonRoot = true
  msg = "Containers must not run as root"
}

violation[msg] {
  input.kind = "Deployment"
  not input.spec.selector.matchLabels.app
  msg = "Containers must provide app label for pod selectors"
}

violation[msg] {
  input.kind = "Deployment"
  not input.spec.template.spec.nodeSelector
  msg = "Deployment must have nodeSelector"
}
