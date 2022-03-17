# Kubernetes Scheduler

## [Scheduling Framework](https://kubernetes.io/docs/concepts/scheduling-eviction/scheduling-framework/)

![](https://raw.githubusercontent.com/kubernetes/website/main/static/images/docs/scheduling-framework-extensions.png)


## [Extension Points](https://kubernetes.io/docs/reference/scheduling/config/#extension-points)

1. queueSort
1. preFilter
1. filter
1. postFilter
1. preScore
1. score
1. reserve
1. permit
1. preBind
1. bind
1. postBind
1. multiPoint

## Default Plugins

https://kubernetes.io/docs/reference/scheduling/config/#scheduling-plugins

## Custom Scheduler
- [kube-batch](https://github.com/kubernetes-sigs/kube-batch)
- [Volcano](https://github.com/volcano-sh/volcano)
- [random-scheduler](random-scheduler): Simplest scheduler to start with.
- [mini-kube-scheduler](https://github.com/nakamasato/mini-kube-scheduler): You can create a Kubernetes scheduler from zero.
    - https://speakerdeck.com/sanposhiho/zi-zuo-sitexue-bukubernetes-schedulerru-men
    - https://event.cloudnativedays.jp/cndt2021/talks/1184
    - https://github.com/sanposhiho/mini-kube-scheduler
