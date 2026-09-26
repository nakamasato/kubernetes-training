# [Manager](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.25.1/pkg/manager/manager.go)

![](diagram.drawio.svg)

Manager owns controller registration, startup, shutdown, and shared cluster dependencies. [Builder](../builder/) registers controllers; [Cluster](../cluster/) supplies Client, Cache, Scheme, RESTMapper, APIReader, and recorders.

These implementation notes and the diagram target controller-runtime v0.25.1, pinned in the repository's go.mod. Use the signatures and call paths below for this version.

## types

### 1. Manager Interface

```go
type Manager interface {
    cluster.Cluster

    Add(Runnable) error

    Elected() <-chan struct{}

    AddMetricsServerExtraHandler(path string, handler http.Handler) error

    AddHealthzCheck(name string, check healthz.Checker) error

    AddReadyzCheck(name string, check healthz.Checker) error

    Start(ctx context.Context) error

    GetWebhookServer() webhook.Server

    GetLogger() logr.Logger

    GetControllerOptions() config.Controller

    GetConverterRegistry() conversion.Registry
}
```

### 2. controllerManager

The concrete Manager stores the Cluster, runnable groups, error channel, leader-election configuration, webhook server, metrics server, and shutdown state. These fields explain why Manager can coordinate resources that a standalone Reconciler cannot:

```go
type controllerManager struct {
    sync.Mutex
    started bool

    stopProcedureEngaged *int64
    errChan              chan error
    runnables            *runnables

    cluster cluster.Cluster

    recorderProvider *intrec.Provider

    resourceLock resourcelock.Interface

    leaderElectionReleaseOnCancel bool

    metricsServer metricsserver.Server

    healthProbeListener net.Listener

    readinessEndpointName string

    livenessEndpointName string

    readyzHandler *healthz.Handler

    healthzHandler *healthz.Handler

    pprofListener net.Listener

    controllerConfig config.Controller

    logger logr.Logger

    leaderElectionStopped chan struct{}

    leaderElectionCancel context.CancelFunc

    elected chan struct{}

    webhookServer webhook.Server
    webhookServerOnce sync.Once

    converterRegistry conversion.Registry

    leaderElectionID string
    leaseDuration time.Duration
    renewDeadline time.Duration
    retryPeriod time.Duration

    gracefulShutdownTimeout time.Duration

    onStoppedLeading func()

    shutdownCtx context.Context

    internalCtx    context.Context
    internalCancel context.CancelFunc

    internalProceduresStop chan struct{}
}
```

### 3. Runnable interface

Start must block while the component runs, honor context cancellation, and return errors to the Manager. A RunnableFunc adapts a function with this signature.

```go
type Runnable interface {
    Start(context.Context) error
}
```

### 4. runnables

Each group queues registrations, starts goroutines, checks readiness where applicable, reports errors, and waits for completion during shutdown.

```go
type runnables struct {
    HTTPServers    *runnableGroup
    Webhooks       *runnableGroup
    Caches         *runnableGroup
    LeaderElection *runnableGroup
    Warmup         *runnableGroup
    Others         *runnableGroup
}
```

```go
type runnableGroup struct {
    ctx    context.Context
    cancel context.CancelFunc

    start        sync.Mutex
    startOnce    sync.Once
    started      bool
    startQueue   []*readyRunnable
    startReadyCh chan *readyRunnable

    stop     sync.RWMutex
    stopOnce sync.Once
    stopped  bool

    errChan chan error

    ch chan *readyRunnable

    wg *sync.WaitGroup

    logger logr.Logger
}
```

## How `Manager` is initialized by New

### 1. Set default values for Options fields with setOptionsDefaults

```go
mgr, err := manager.New(cfg, manager.Options{})
if err != nil {
    return err
}
```

Important options and their roles:

| Option | Role/default |
| --- | --- |
| Scheme, MapperProvider, HTTPClient | Forwarded to Cluster for type mappings and transport |
| Cache / NewCache | Cache configuration / constructor |
| Client / NewClient | Client configuration / constructor |
| Metrics | `metricsserver.Options`; replaces MetricsBindAddress |
| WebhookServer | Server returned by GetWebhookServer; defaulted when omitted |
| LeaderElection | Disabled unless enabled |
| LeaseDuration / RenewDeadline / RetryPeriod | Default leadership timing: 15s / 10s / 2s |
| LeaderElectionID / Namespace | Identify the shared election lock |
| HealthProbeBindAddress | Health/readiness listener configuration |
| ReadinessEndpointName / LivenessEndpointName | `readyz` / `healthz` |
| GracefulShutdownTimeout | Time allowed for components to stop |
| BaseContext | Base context used to run managed components |

`setOptionsDefaults` supplies constructors, timing values, logging, and listener defaults. New then constructs the dependencies; it does not start reconciliation.

### 2. Initialize a controllerManager

1. `cluster.New` receives Scheme, MapperProvider, HTTPClient, Cache/NewCache, Client/NewClient, Logger, and EventBroadcaster options.
2. Cluster builds the shared cache, read/write client, uncached reader, RESTMapper, and recorders. See its [construction walkthrough](../cluster/#new).
3. Manager creates runnable groups, event recording, metrics/health listeners, and any leader-election resource lock.
4. Manager stores these with its internal contexts and error channel. The Cluster itself is a Runnable, classified into the Caches group when added.

### 3. Bind a Controller to the Manager

```go
if err := ctrl.NewControllerManagedBy(mgr).
    For(&appsv1.ReplicaSet{}).
    Owns(&corev1.Pod{}).
    Complete(&ReplicaSetReconciler{Client: mgr.GetClient()}); err != nil {
    return err
}
```

Builder.Build calls doController and doWatch. The controller is registered through Manager.Add, and watches are configured before startup. The [ReplicaSet sample](../example-controller/) supplies the Reconciler used here.

### 4. controllerManager.Start calls runnable group Start

1. Register the Cluster runnable and managed servers.
2. Start HTTPServers, then Webhooks. Webhooks are available before cache startup, which can require conversion webhooks.
3. Start Caches and wait for their readiness checks.
4. Start Others, which do not require leadership.
5. Start Warmup work, allowing controllers configured for warmup to prepare sources before leadership.
6. Start leader election and the LeaderElection group after acquisition. When election is disabled, start that group directly.
7. Wait for cancellation or an error and run coordinated shutdown.

Controllers normally belong to LeaderElection, not Others. Disabling election at Manager level does not change their classification; it removes the need to acquire leadership before starting that group.

## Internal process of adding a `Controller` to a `Manager`

Manager.Add checks Manager state under its lock and delegates to runnables.Add. It no longer injects dependencies using SetFields. The current classification logic is:

```go
func (r *runnables) Add(fn Runnable) error {
    switch runnable := fn.(type) {
    case *Server:
        if runnable.NeedLeaderElection() {
            return r.LeaderElection.Add(fn, nil)
        }
        return r.HTTPServers.Add(fn, nil)
    case hasCache:
        return r.Caches.Add(fn, func(ctx context.Context) bool {
            return runnable.GetCache().WaitForCacheSync(ctx)
        })
    case webhook.Server:
        return r.Webhooks.Add(fn, nil)
    case warmupRunnable, LeaderElectionRunnable:
        if warmupRunnable, ok := fn.(warmupRunnable); ok {
            if err := r.Warmup.Add(RunnableFunc(warmupRunnable.Warmup), nil); err != nil {
                return err
            }
        }

        leaderElectionRunnable, ok := fn.(LeaderElectionRunnable)
        if !ok {
            return r.LeaderElection.Add(fn, nil)
        }

        if !leaderElectionRunnable.NeedLeaderElection() {
            return r.Others.Add(fn, nil)
        }
        return r.LeaderElection.Add(fn, nil)
    default:
        return r.LeaderElection.Add(fn, nil)
    }
}
```

HTTP servers can opt into leadership. Objects exposing GetCache receive a cache synchronization readiness check. Webhook servers are their own group. Controllers exposing Warmup have warmup work registered separately; NeedLeaderElection decides their main group. A plain Runnable without a leadership interface defaults to LeaderElection.

## `Manager.GetClient()` and `GetScheme()`

These methods delegate to the embedded Cluster. Pass the returned dependencies explicitly to reconcilers and other components. GetClient uses the shared cache for eligible reads and the API server for writes; GetAPIReader bypasses the cache. Register custom API types in the Scheme before constructing components that need to resolve them.

## Example

The [sample](main.go) creates a Manager, builds Pod and Deployment controllers using `reconcile.Func`, adds a RunnableFunc, and starts the Manager with `ctrl.SetupSignalHandler()`. Every New, Complete, Add, and Start error is checked.

A Runnable must not ignore cancellation when it runs continuously:

```go
err := mgr.Add(manager.RunnableFunc(func(ctx context.Context) error {
    // Initialize component resources here.
    <-ctx.Done()
    return nil
}))
if err != nil {
    return err
}
return mgr.Start(ctrl.SetupSignalHandler())
```

## Run

From the repository root with a kubeconfig and list/watch access to Pods and Deployments:

```sh
go run ./contents/kubernetes-operator/controller-runtime/manager
```

Existing objects enqueue initial requests. The sample logs `podReconciler is called` or `deploymentReconciler is called` with the namespace/name, plus `RunnableFunc is called`. Create a test object in another terminal:

```sh
kubectl create deployment manager-demo --image=nginx
kubectl scale deployment manager-demo --replicas=2
kubectl delete deployment manager-demo
```

Each operation can produce multiple Pod/Deployment events, and duplicate queued keys can collapse. Ctrl+C cancels the Manager and waits for its managed components to stop. `created manager` is logged only after Start returns in this sample, not as a readiness signal.
