package main

import (
	"context"
	"fmt"
	"log"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	clientgocache "k8s.io/client-go/tools/cache"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func main() {
	if err := run(ctrl.SetupSignalHandler()); err != nil {
		log.Fatal(err)
	}
}

func run(parent context.Context) (runErr error) {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()
	cfg, err := ctrl.GetConfig()
	if err != nil {
		return err
	}
	// cache.New creates the HTTP client and REST mapper when omitted.
	podCache, err := cache.New(cfg, cache.Options{Scheme: scheme.Scheme})
	if err != nil {
		return err
	}
	// Register before starting the cache. Get must only be used after sync.
	informer, err := podCache.GetInformer(ctx, &v1.Pod{})
	if err != nil {
		return err
	}
	if _, err := informer.AddEventHandler(clientgocache.ResourceEventHandlerFuncs{
		AddFunc:    func(obj interface{}) { printEvent("OnAdd", obj) },
		UpdateFunc: func(_, obj interface{}) { printEvent("OnUpdate", obj) },
		DeleteFunc: func(obj interface{}) { printEvent("OnDelete", obj) },
	}); err != nil {
		return err
	}
	cacheErr := make(chan error, 1)
	go func() {
		cacheErr <- podCache.Start(ctx)
		cancel()
	}()
	defer func() {
		cancel()
		if err := <-cacheErr; err != nil {
			runErr = err
		}
	}()
	if !podCache.WaitForCacheSync(ctx) {
		return fmt.Errorf("cache sync interrupted")
	}
	fmt.Println("cache is synced")
	pod := &v1.Pod{}
	if err := podCache.Get(ctx, client.ObjectKey{Namespace: "default", Name: "nginx"}, pod); err != nil {
		return err
	}
	fmt.Printf("cached Pod: %s/%s\n", pod.Namespace, pod.Name)
	<-ctx.Done()
	return nil
}

func printEvent(event string, obj interface{}) {
	key, err := clientgocache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		log.Print(err)
		return
	}
	fmt.Println(event, key)
}
