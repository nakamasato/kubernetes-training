# [leaderelection](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/leaderelection)

> Package leaderelection contains a constructor for a leader election resource lock. This is used to ensure that multiple copies of a controller manager can be run with only one active set of controllers, for active-passive HA.
>
> It uses built-in Kubernetes leader election APIs.

The comparison below describes two election strategies. controller-runtime v0.25.1 uses lease-based election; leader-for-life is included for comparison.

## Types

### 1. **Leader-for-life**

*In the "leader for life" approach, a specific Pod is the leader. Once established (by creating a lock record), the Pod is the leader until it is destroyed. This strategy relies on the platform establishing that the previous leader has stopped before replacing it. The leader does not need to renew a lease, consider stepping down, or do anything related to election activity once it becomes the leader.*

### 2. **Lease-based**

*Leases provide a way to indirectly observe whether the leader still exists. The leader must periodically renew its lease, usually by updating a timestamp in its lock record. If it fails to do so, it is presumed dead, and a new election takes place. If the leader is in fact still alive but unreachable, it is expected to gracefully step down. A variety of factors can cause a leader to fail at updating its lease, but continue acting as the leader before succeeding at stepping down.*


## How Manager uses leader election

```go
mgr, err := ctrl.NewManager(cfg, ctrl.Options{
    LeaderElection:          true,
    LeaderElectionID:        "example-controller.example.com",
    LeaderElectionNamespace: "default",
    LeaderElectionResourceLock: "leases",
})
if err != nil {
    return err
}
```

Replicas coordinate using the same Lease name and namespace. Grant the ServiceAccount access to get/create/update the Lease. Use different election IDs for unrelated controllers. Manager starts webhook/cache services independently of leadership and starts its LeaderElection runnable group after acquisition. If leadership is lost, the Manager reports an error; the process should stop rather than continue making changes.

LeaseDuration, RenewDeadline, and RetryPeriod default to 15s, 10s, and 2s. LeaderElectionReleaseOnCancel can shorten failover, but only enable it when all protected work has stopped before the lock is released. Lease-based election is not fencing: external systems may need a separate mechanism to reject operations from a stale leader.

The local examples leave leader election disabled to simplify single-process exercises. See [Manager](../manager/) for runnable classification and startup order.

## References

1. Go packages
    1. [client-go/tools/leaderelection](https://pkg.go.dev/k8s.io/client-go/tools/leaderelection)
    1. [controller-runtime/pkg/leaderelection](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/leaderelection): lease-based leader election
    1. [operator-framework/operator-lib/leader](https://pkg.go.dev/github.com/operator-framework/operator-lib/leader): Leader For Life
1. Readings
    1. [Operator SDK # Leader Election](https://sdk.operatorframework.io/docs/building-operators/golang/advanced-topics/#leader-election)
    1. https://d-kuro.github.io/post/kubernetes-leader-election/
    1. [Leader election in Kubernetes using client-go](https://itnext.io/leader-election-in-kubernetes-using-client-go-a19cbe7a9a85)
