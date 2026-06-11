package main

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func main() {

	withStruct()
	withFunc()
}

type reconciler struct {
}

func (reconciler) reconcile(ctx context.Context, o reconcile.Request) (reconcile.Result, error) {
	fmt.Printf("reconcile is called with %s/%s\n", o.Namespace, o.Name)
	return reconcile.Result{}, nil
}

func withStruct() {
	r := reconciler{}
	res, err := r.reconcile(context.Background(), reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "default", Name: "test"}})
	if err != nil || res.Requeue || res.RequeueAfter != time.Duration(0) {
		fmt.Printf("error: %v, res %v", err, res)
	} else {
		fmt.Println("res is expected")
	}
}

func withFunc() {
	r := reconcile.Func(func(ctx context.Context, o reconcile.Request) (reconcile.Result, error) {
		fmt.Printf("reconcile is called with %s/%s\n", o.Namespace, o.Name)
		return reconcile.Result{}, nil
	}) // implements reconcile.Reconciler interface

	res, err := r.Reconcile(context.Background(), reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "default", Name: "test"}})
	if err != nil || res.Requeue || res.RequeueAfter != time.Duration(0) {
		fmt.Printf("error: %v, res %v", err, res)
	} else {
		fmt.Println("res is expected")
	}
}
