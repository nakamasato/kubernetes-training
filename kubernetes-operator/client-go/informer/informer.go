package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var eventHandlerMessage = "%s is called for Pod (key: %s)\n"

func main() {
	var defaultKubeConfigPath string
	if home := homedir.HomeDir(); home != "" {
		// build kubeconfig path from $HOME dir
		defaultKubeConfigPath = filepath.Join(home, ".kube", "config")
	}

	// set kubeconfig flag
	kubeconfig := flag.String("kubeconfig", defaultKubeConfigPath, "kubeconfig config file")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatalf("Building config from flags: %v", err)
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("getting kubernetes client set: %v", err)
	}

	informerFactory := informers.NewSharedInformerFactory(kubeClient, time.Second*30)

	podInformer := informerFactory.Core().V1().Pods()
	_, err = podInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    handleAdd,
			UpdateFunc: handleUpdate,
			DeleteFunc: handleDelete,
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	ch := ctx.Done()
	informerFactory.Start(ch)
	defer informerFactory.Shutdown()

	cacheSynced := podInformer.Informer().HasSynced
	if ok := cache.WaitForCacheSync(ch, cacheSynced); !ok {
		log.Print("cache sync stopped")
		return
	}
	log.Println("cache is synced")

	<-ch
}

func handleAdd(obj interface{}) {
	key := getKeyFromObj(obj)
	if pod, ok := obj.(*v1.Pod); !ok {
		fmt.Println("couldn't convert to pod")
	} else {
		pod = pod.DeepCopy() // informer objects are shared and must be treated as read-only
		pod.SetLabels(map[string]string{"test": "test"})
		fmt.Printf("converted to Pod label: %s\n", pod.GetLabels())
	}
	log.Printf(eventHandlerMessage, "handleAdd", key)
}

func handleUpdate(old, new interface{}) {
	key := getKeyFromObj(new)
	log.Printf(eventHandlerMessage, "handleUpdate", key)
}

func handleDelete(obj interface{}) {
	key := getKeyFromObj(obj)
	log.Printf(eventHandlerMessage, "handleDelete", key)
}

func getKeyFromObj(obj interface{}) string {
	var key string
	var err error
	if key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj); err != nil {
		log.Printf("failed to get key from the cache %s\n", err.Error())
		return ""
	}
	return key
}
