# require prefix rule

## Apply Constraint Template and Constraint

```
± kubectl apply -f gatekeeper/require-prefix/k8srequiredprefix.yaml 
constrainttemplate.templates.gatekeeper.sh/k8srequiredprefixes created

± kubectl apply -f gatekeeper/require-prefix/k8srequiredprefix-ns.yaml 
k8srequiredprefixes.constraints.gatekeeper.sh/ns-must-start-with-prefix created
```

## Check

```
± kubectl create ns ns-naka                                           
Error from server ([denied by ns-must-start-with-prefix] you must provide prefix: dev, provided: ns-naka): admission webhook "validation.gatekeeper.sh" denied the request: [denied by ns-must-start-with-prefix] you must provide prefix: dev, provided: ns-naka

± kubectl create ns dev-naka
namespace/dev-naka created
```

## Clean up

```
± kubectl delete ns dev-naka
namespace "dev-naka" deleted

± kubectl delete -f gatekeeper/require-prefix/k8srequiredprefix-ns.yaml 
k8srequiredprefixes.constraints.gatekeeper.sh "ns-must-start-with-prefix" deleted

± kubectl delete -f gatekeeper/require-prefix/k8srequiredprefix.yaml    
constrainttemplate.templates.gatekeeper.sh "k8srequiredprefixes" deleted
```
