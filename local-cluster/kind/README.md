# kind

## Description

kind is a tool for running local Kubernetes clusters using Docker container “nodes”.

## Prerequisite

- go (1.11+)
- docker installed

## Quick Start

```
GO111MODULE="on" go get sigs.k8s.io/kind@v0.10.0 && kind create cluster
```

## Installation

https://kind.sigs.k8s.io/docs/user/quick-start/#installation

On Mac:

```
brew install kind
```

```
kind version
kind v0.10.0 go1.15.7 darwin/amd64
```

## Configure a cluster

```
kind create cluster --config <kind-config>.yaml
```

Example 1: enable alpha feature https://kind.sigs.k8s.io/docs/user/configuration/#feature-gates (https://kubernetes.io/docs/reference/command-line-tools-reference/feature-gates/)

```
kind create cluster --config cluster-with-alpha-feature.yaml
Creating cluster "kind" ...
 ✓ Ensuring node image (kindest/node:v1.20.2) 🖼
 ✓ Preparing nodes 📦
 ✓ Writing configuration 📜
 ✓ Starting control-plane 🕹️
 ✓ Installing CNI 🔌
 ✓ Installing StorageClass 💾
Set kubectl context to "kind-kind"
You can now use your cluster with:

kubectl cluster-info --context kind-kind

Have a nice day! 👋
```

## Why kind?

- kind supports multi-node (including HA) clusters
- kind supports building Kubernetes release builds from source
    support for make / bash / docker, or bazel, in addition to pre-published builds
- kind supports Linux, macOS and Windows
- kind is a CNCF certified conformant Kubernetes installer

## Official Document

https://kind.sigs.k8s.io/
