# Kubernetes Components

## etcd

## kubernetes-scheduler

## kube-proxy

Modes:
- user space mode (legacy)
- iptable mode (default)
- ipvs mode

Reference:
- [kube-proxy詳細](https://ichi.pro/k-8-s-kube-proxy-no-shosai-3791464738960)
- [Comparing kube-proxy modes: iptables or IPVS?](https://www.tigera.io/blog/comparing-kube-proxy-modes-iptables-or-ipvs/)
-

## Cloud Controller Manager

[cloud-provider](https://github.com/kubernetes/cloud-provider/tree/master/controllers):
1. node_controller
1. nodelifecycle_controller
1. route_controller
1. service_controller

References:
- `cloud-controller-manager` moved to `k8s.io/cloud-provider` in [PR#95748](https://github.com/kubernetes/kubernetes/pull/95740) to resolve [Issue#29](https://github.com/kubernetes/cloud-provider/issues/29)
- https://bells17.medium.com/kubernetes-ccm-d4d3c71ba523
