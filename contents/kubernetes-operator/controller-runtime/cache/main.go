package main

import (
	"context"
	"fmt"
	"log"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	clientgocache "k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func main() {
	// Get a kubeconfig
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	mapper, err := func(c *rest.Config) (meta.RESTMapper, error) {
		return apiutil.NewDynamicRESTMapper(c)
	}(cfg)
	if err != nil {
		log.Fatal(err)
	}

	cache, err := cache.New(cfg, cache.Options{Scheme: scheme.Scheme, Mapper: mapper}) // &informerCache{InformersMap: im}
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	pod := &v1.Pod{}
	cache.Get(ctx, client.ObjectKeyFromObject(pod), pod)

	// Start Cache
	go func() {
		if err := cache.Start(ctx); err != nil { // func (m *InformersMap) Start(ctx context.Context) error {
			log.Fatal(err)
		}
	}()
	fmt.Println("cache is started")
	if isSynced := cache.WaitForCacheSync(ctx); !isSynced {
		log.Fatal("failed to sync cache")
	}
	fmt.Println("cache is synced")

	if err := cache.Get(ctx, client.ObjectKey{
		Namespace: "default",
		Name:      "nginx",
	}, pod); err != nil {
		log.Fatal(err)
	}
	fmt.Println(pod)

	informer, err := cache.GetInformerForKind(ctx, schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Pod",
	})
	if err != nil {
		log.Fatal(err)
	}

	informer.AddEventHandler(clientgocache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			fmt.Println("OnAdd", obj)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			fmt.Println("OnUpdate", oldObj, newObj)
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Println(obj)
		},
	})

	if isSynced := informer.HasSynced(); !isSynced {
		log.Fatal("failed to sync informer")
	}
	fmt.Println("informer is synced")

	time.Sleep(30 * time.Second) // just for time to check informer
}
