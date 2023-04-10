# [kubelet](https://github.com/kubernetes/kubernetes/tree/master/pkg/kubelet)

## Main roles

1. [configmap](https://github.com/kubernetes/kubernetes/tree/master/pkg/kubelet/configmap)
1. [secret](https://github.com/kubernetes/kubernetes/tree/master/pkg/kubelet/secret)
1. Create container
1. Garbage Collection (image, container)

## Implementation

1. [Run](https://github.com/kubernetes/kubernetes/blob/ad18954259eae3db51bac2274ed4ca7304b923c4/pkg/kubelet/kubelet.go#L1509-L1605)

## Mode

[standalone-kubelet-tutorial](https://github.com/kelseyhightower/standalone-kubelet-tutorial)
