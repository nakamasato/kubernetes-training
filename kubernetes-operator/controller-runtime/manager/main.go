package main

import (
	"context"
	"flag"
	"os"

	"go.uber.org/zap/zapcore"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var (
	mgr manager.Manager
	// NB: don't call SetLogger in init(), or else you'll mess up logging in the main suite.
	log = logf.Log.WithName("manager-examples")
)

func main() {
	// Prepare log
	opts := zap.Options{
		Development: true,
		TimeEncoder: zapcore.ISO8601TimeEncoder,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	// Get a kubeconfig
	cfg, err := config.GetConfig()
	if err != nil {
		log.Error(err, "unable to get kubeconfig")
		os.Exit(1)
	}

	// Create a Manager
	mgr, err = manager.New(cfg, manager.Options{})
	if err != nil {
		log.Error(err, "unable to set up manager")
		os.Exit(1)
	}

	// Create Reconcilers
	podReconciler := reconcile.Func(func(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
		log.Info("podReconciler is called", "req", req)
		return reconcile.Result{}, nil
	})

	deploymentReconciler := reconcile.Func(func(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
		log.Info("deploymentReconciler is called", "req", req)
		return reconcile.Result{}, nil
	})

	// Create Controller with Manager
	ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Pod{}).
		Complete(podReconciler)

	ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		Complete(deploymentReconciler)

	// Add raw RunnableFunc to Manager
	mgr.Add(manager.RunnableFunc(func(ctx context.Context) error {
		log.Info("RunnableFunc is called")
		return nil
	}))

	// Start the Manager
	ctx := context.Background()
	err = mgr.Start(ctx)
	if err != nil {
		log.Error(err, "unable to start manager")
	}
	log.Info("created manager", "manager", mgr)
}
