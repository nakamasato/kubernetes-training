package main

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func main() {

	withStruct()
	withFunc()
}

type reconciler struct {
}

var _ reconcile.Reconciler = reconciler{}

func (reconciler) Reconcile(ctx context.Context, o reconcile.Request) (reconcile.Result, error) {
	fmt.Printf("reconcile is called with %s/%s\n", o.Namespace, o.Name)
	return reconcile.Result{}, nil
}

func withStruct() {
	r := reconciler{}
	res, err := r.Reconcile(context.Background(), reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "default", Name: "test"}})
	if err != nil || res.RequeueAfter != 0 {
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
	if err != nil || res.RequeueAfter != 0 {
		fmt.Printf("error: %v, res %v", err, res)
	} else {
		fmt.Println("res is expected")
	}
}
