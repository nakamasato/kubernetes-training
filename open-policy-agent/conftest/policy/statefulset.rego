package main

warn[msg] {
  input.kind = "StatefulSet"
  not input.spec.template.spec.securityContext.runAsNonRoot = true
  msg = "Containers should not run as root"
}

warn[msg] {
  input.kind = "StatefulSet"
  not input.spec.selector.matchLabels.app
  msg = "Containers should provide app label for pod selectors"
}

contains(array, elem) {
  array[_] = elem
}

deny[msg] {
  input.kind = "StatefulSet"
  not input.spec.template.spec.nodeSelector
  msg = "StatefulSet must have nodeSelector"
}

deny[msg] {
  input.kind = "StatefulSet"
  not input.spec.template.spec.nodeSelector.nodegroup
  msg = "StatefulSet must have nodeSelector with nodegroup as a key"
}

deny[msg] {
  input.kind = "StatefulSet"
  input.spec.template.spec.nodeSelector.nodegroup
  not contains(["dev", "prod"], input.spec.template.spec.nodeSelector.nodegroup)
  msg = "StatefulSet must have nodeSelector with nodegroup as a key and prod or dev as value"
}

deny[msg] {
  input.kind = "StatefulSet"
  input.metadata.namespace = "prod"
  input.spec.template.spec.nodeSelector.nodegroup
  not input.spec.template.spec.nodeSelector.nodegroup = "prod"
  msg = "nodegroup must be prod in prod namespace"
}

deny[msg] {
  input.kind = "StatefulSet"
  not input.metadata.namespace = "prod"
  input.spec.template.spec.nodeSelector.nodegroup
  not input.spec.template.spec.nodeSelector.nodegroup = "dev"
  msg = "nodegroup must be dev in non-prod namespace"
}
