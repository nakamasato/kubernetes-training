package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	mysqlv1alpha1 "github.com/nakamasato/mysql-operator/api/v1alpha1"
	"go.uber.org/zap/zapcore"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var (
	log    = logf.Log.WithName("source-examples")
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(mysqlv1alpha1.AddToScheme(scheme))
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
}

func main() {
	// Prepare log
	opts := zap.Options{
		Development: true,
		TimeEncoder: zapcore.ISO8601TimeEncoder,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
	if err := run(ctrl.SetupSignalHandler()); err != nil {
		log.Error(err, "source failed")
		os.Exit(1)
	}
}

func run(parent context.Context) (runErr error) {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	// Get a kubeconfig
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	// Let cache.New create its HTTP client and REST mapper.
	objectCache, err := cache.New(cfg, cache.Options{Scheme: scheme})
	if err != nil {
		return err
	}
	cacheErr := make(chan error, 1)
	go func() {
		cacheErr <- objectCache.Start(ctx)
		cancel()
	}()
	defer func() {
		cancel()
		if err := <-cacheErr; err != nil {
			runErr = err
		}
	}()

	// Prepare queue and eventHandler
	queue := workqueue.NewTypedRateLimitingQueue(workqueue.DefaultTypedControllerRateLimiter[WorkQueueItem]())
	defer queue.ShutDown()
	stopShutdown := context.AfterFunc(ctx, queue.ShutDown)
	defer stopShutdown()

	eventHandler := handler.TypedFuncs[client.Object, WorkQueueItem]{
		CreateFunc: func(ctx context.Context, e event.CreateEvent, q workqueue.TypedRateLimitingInterface[WorkQueueItem]) {
			log.Info("CreateFunc is called", "object", e.Object.GetName())
			q.Add(newWorkQueueItem("Create", e.Object))
		},
		UpdateFunc: func(ctx context.Context, e event.UpdateEvent, q workqueue.TypedRateLimitingInterface[WorkQueueItem]) {
			log.Info("UpdateFunc is called", "objectNew", e.ObjectNew.GetName(), "objectOld", e.ObjectOld.GetName())
			q.Add(newWorkQueueItem("Update", e.ObjectNew))
		},
		DeleteFunc: func(ctx context.Context, e event.DeleteEvent, q workqueue.TypedRateLimitingInterface[WorkQueueItem]) {
			log.Info("DeleteFunc is called", "object", e.Object.GetName())
			q.Add(newWorkQueueItem("Delete", e.Object))
		},
	}

	kindMysqlUser := source.TypedKind[client.Object](objectCache, &mysqlv1alpha1.MySQLUser{}, eventHandler)
	kindPod := source.TypedKind[client.Object](objectCache, &v1.Pod{}, eventHandler)

	// Start Source
	if err := kindMysqlUser.Start(ctx, queue); err != nil {
		return err
	}
	if err := kindPod.Start(ctx, queue); err != nil {
		return err
	}

	// Bound startup when a CRD or list/watch permission is missing.
	syncCtx, stopSync := context.WithTimeout(ctx, 30*time.Second)
	defer stopSync()
	// Wait for cache
	if err := kindMysqlUser.WaitForSync(syncCtx); err != nil {
		return err
	}
	if err := kindPod.WaitForSync(syncCtx); err != nil {
		return err
	}
	log.Info("kind is ready")

	for {
		item, shutdown := queue.Get()
		if shutdown {
			break
		}
		log.Info("got item", "item", item)
		queue.Forget(item)
		queue.Done(item)
	}
	return nil
}

type WorkQueueItem struct {
	Event     string
	Kind      string
	Namespace string
	Name      string
}

func newWorkQueueItem(eventType string, obj client.Object) WorkQueueItem {
	return WorkQueueItem{Event: eventType, Kind: fmt.Sprintf("%T", obj), Namespace: obj.GetNamespace(), Name: obj.GetName()}
}
