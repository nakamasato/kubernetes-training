package main

deny[msg] {
  input.kind = "Deployment"
  not input.spec.template.spec.securityContext.runAsNonRoot = true
  msg = "Containers must not run as root"
}

deny[msg] {
  input.kind = "Deployment"
  not input.spec.selector.matchLabels.app
  msg = "Containers must provide app label for pod selectors"
}

contains(array, elem) {
  array[_] = elem
}

deny[msg] {
  input.kind = "Deployment"
  not input.spec.template.spec.nodeSelector.nodegroup
  msg = "Deployment must have nodeSelector with nodegroup as a key"
}

deny[msg] {
  input.kind = "Deployment"
  not input.spec.template.spec.nodeSelector.nodegroup = "prod"
  msg = "Deployment must have nodeSelector with nodegroup as a key and prod as value"
}
