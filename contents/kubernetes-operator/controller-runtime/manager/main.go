package main

import (
	"context"
	"fmt"
	"os"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var (
	mgr manager.Manager
	// NB: don't call SetLogger in init(), or else you'll mess up logging in the main suite.
	log = logf.Log.WithName("manager-examples")
)

func main() {
	ctx := context.Background()

	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Errorf("unable to get kubeconfig %v\n", err)
		os.Exit(1)
	}

	mgr, err = manager.New(cfg, manager.Options{})
	if err != nil {
		fmt.Errorf("unable to set up manager %v\n", err)
		os.Exit(1)
	}

	podReconciler := reconcile.Func(func(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
		fmt.Printf("podReconciler is called for %v\n", req)
		return reconcile.Result{}, nil
	})

	deploymentReconciler := reconcile.Func(func(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
		fmt.Printf("deploymentReconciler is called for %v\n", req)
		return reconcile.Result{}, nil
	})

	ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Pod{}).
		Complete(podReconciler)

	ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		Complete(deploymentReconciler)

	err = mgr.Start(ctx)
	if err != nil {
		log.Error(err, "unable to start manager")
	}
	log.Info("created manager", "manager", mgr)
}
