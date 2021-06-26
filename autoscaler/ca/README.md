# Cluster Autoscaler

## Reference

https://eksworkshop.com/scaling/deploy_ca/

## Prerequisite
- add autoscaling policy to worker node group

## apply autoscaling yaml

```
kubectl apply -f cluster_autoscaler.yaml -n kube-system
```


## test app

```
kubectl apply -f nginx.yaml
kubectl get deployment/nginx-to-scaleout

NAME                DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
nginx-to-scaleout   1         1         1            0           8s

```

scale out to 10 replicas

```
kubectl scale --replicas=10 deployment/nginx-to-scaleout
```

-> automatically add two new nodes

```
kubectl get nodes
NAME                                            STATUS    ROLES     AGE       VERSION
ip-10-0-0-16.ap-northeast-1.compute.internal    Ready     <none>    2m        v1.12.7
ip-10-0-0-82.ap-northeast-1.compute.internal    Ready     <none>    1h        v1.12.7
ip-10-0-1-177.ap-northeast-1.compute.internal   Ready     <none>    1m        v1.12.7
ip-10-0-1-224.ap-northeast-1.compute.internal   Ready     <none>    1h        v1.12.7
```

scale in to 2 replicas

```
kubectl scale --replicas=2 deployment/nginx-to-scaleout
```

aboute ten mins later ... -> scaling down

```
I0727 15:41:44.855775       1 scale_down.go:387] ip-10-0-0-16.ap-northeast-1.compute.internal was unneeded for 10m4.589275269s
I0727 15:41:44.855814       1 scale_down.go:387] ip-10-0-1-177.ap-northeast-1.compute.internal was unneeded for 7m18.53883852s
I0727 15:41:44.937768       1 scale_down.go:594] Scale-down: removing empty node ip-10-0-0-16.ap-northeast-1.compute.internal
I0727 15:41:44.938199       1 factory.go:33] Event(v1.ObjectReference{Kind:"ConfigMap", Namespace:"kube-system", Name:"cluster-autoscaler-status", UID:"2d9cdefe-b081-11e9-be78-0a78d103117e", APIVersion:"v1", ResourceVersion:"9275", FieldPath:""}): type: 'Normal' reason: 'ScaleDownEmpty' Scale-down: removing empty node ip-10-0-0-16.ap-northeast-1.compute.internal
I0727 15:41:44.952658       1 delete.go:53] Successfully added toBeDeletedTaint on node ip-10-0-0-16.ap-northeast-1.compute.internal
I0727 15:41:45.162880       1 aws_manager.go:341] Terminating EC2 instance: i-05df3b79e52165a1d
I0727 15:41:45.163164       1 factory.go:33] Event(v1.ObjectReference{Kind:"Node", Namespace:"", Name:"ip-10-0-0-16.ap-northeast-1.compute.internal", UID:"38a0545b-b083-11e9-be78-0a78d103117e", APIVersion:"v1", ResourceVersion:"9282", FieldPath:""}): type: 'Normal' reason: 'ScaleDown' node removed by cluster autoscaler
```


## check logs

```
kubectl logs -f deployment/cluster-autoscaler -n kube-system
```

### Node group is not ready

```
W0727 15:23:39.646482       1 scale_up.go:105] Node group terraform-eks-demo is not ready for scaleup
```

max instance is too small or something
