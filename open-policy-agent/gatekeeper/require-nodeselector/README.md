# require prefix rule

## Debug

```
{"level":"error","ts":1604385393.0477803,"logger":"controller","msg":"Reconciler error","controller":"constrainttemplate-controller","name":"requirenodeselectors","namespace":"","error":"customresourcedefinitions.apiextensions.k8s.io \"requirenodeselectors.constraints.gatekeeper.sh\" already exists","stacktrace":"github.com/go-logr/zapr.(*zapLogger).Error\n\t/go/src/github.com/open-policy-agent/gatekeeper/vendor/github.com/go-logr/zapr/zapr.go:128\nsigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).reconcileHandler\n\t/go/src/github.com/open-policy-agent/gatekeeper/vendor/sigs.k8s.io/controller-runtime/pkg/internal/controller/controller.go:246\nsigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).processNextWorkItem\n\t/go/src/github.com/open-policy-agent/gatekeeper/vendor/sigs.k8s.io/controller-runtime/pkg/internal/controller/controller.go:218\nsigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).worker\n\t/go/src/github.com/open-policy-agent/gatekeeper/vendor/sigs.k8s.io/controller-runtime/pkg/internal/controller/controller.go:197\nk8s.io/apimachinery/pkg/util/wait.BackoffUntil.func1\n\t/go/src/github.com/open-policy-agent/gatekeeper/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:155\nk8s.io/apimachinery/pkg/util/wait.BackoffUntil\n\t/go/src/github.com/open-policy-agent/gatekeeper/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:156\nk8s.io/apimachinery/pkg/util/wait.JitterUntil\n\t/go/src/github.com/open-policy-agent/gatekeeper/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:133\nk8s.io/apimachinery/pkg/util/wait.Until\n\t/go/src/github.com/open-policy-agent/gatekeeper/vendor/k8s.io/apimachinery/pkg/util/wait/wait.go:90"}
```

## Apply Constraint Template and Constraint

```
kubectl apply -f gatekeeper/require-nodeselector/require-nodeselector-constraint-template.yaml
kubectl apply -f gatekeeper/require-nodeselector/require-nodeselector.yaml
```

## Check

```

```

## Clean up

```
```
