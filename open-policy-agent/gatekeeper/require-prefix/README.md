# require prefix rule


```
[20-11-03 14:00:47] masato-naka at ip-192-168-31-162 in ~/repos/MasatoNaka/kubernetes-training/open-policy-agent on add-opa-example ✘
± kubectl apply -f gatekeeper/require-prefix/k8srequiredprefix.yaml 
constrainttemplate.templates.gatekeeper.sh/k8srequiredprefixes unchanged

[20-11-03 14:00:55] masato-naka at ip-192-168-31-162 in ~/repos/MasatoNaka/kubernetes-training/open-policy-agent on add-opa-example ✘
± kubectl apply -f gatekeeper/require-prefix/k8srequiredprefix-ns.yaml 
k8srequiredprefixes.constraints.gatekeeper.sh/ns-must-start-with-ns created

[20-11-03 14:01:04] masato-naka at ip-192-168-31-162 in ~/repos/MasatoNaka/kubernetes-training/open-policy-agent on add-opa-example ✘
± kubectl create ns ns-naka                                           
namespace/ns-naka created

[20-11-03 14:01:10] masato-naka at ip-192-168-31-162 in ~/repos/MasatoNaka/kubernetes-training/open-policy-agent on add-opa-example ✘
± kubectl create ns naka   
Error from server ([denied by ns-must-start-with-ns] you must provide prefix: {"ns"}, provided: naka): admission webhook "validation.gatekeeper.sh" denied the request: [denied by ns-must-start-with-ns] you must provide prefix: {"ns"}, provided: naka
```