# etcd

## Install

```
brew install etcd
```

```
etcd --version
etcd Version: 3.5.2
Git SHA: 99018a77b
Go Version: go1.17.6
Go OS/Arch: darwin/amd64
```

## Quickstart


1. Start etcd.

    ```
    etcd
    ```
1. Set a key `greeting` and a value `Hello, etcd`.
    ```
    etcdctl put greeting "Hello, etcd"
    OK
    ```
1. Get the value for the key `greeting`
    ```
    ± etcdctl get greeting
    greeting
    Hello, etcd
    ```

## Reference
- https://etcd.io
