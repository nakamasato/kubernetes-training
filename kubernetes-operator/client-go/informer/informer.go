package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
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
		log.Printf("Building config from flags, %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Printf("getting kubernetes client set %s\n", err.Error())
	}

	informerFactory := informers.NewSharedInformerFactory(kubeClient, time.Second*30)

	podInformer := informerFactory.Core().V1().Pods()
	podInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    handleAdd,
			UpdateFunc: handleUpdate,
			DeleteFunc: handleDelete,
		},
	)

	ch := make(chan struct{})
	informerFactory.Start(ch)

	cacheSynced := podInformer.Informer().HasSynced
	if ok := cache.WaitForCacheSync(ch, cacheSynced); !ok {
		log.Printf("cache is not synced")
	}
	log.Println("cache is synced")

	go wait.Until(run, time.Second*10, ch)
	<-ch
}

func run() {
	log.Println("run")
}

func handleAdd(obj interface{}) {
	key := getKeyFromObj(obj)
	if pod, ok := obj.(*v1.Pod); !ok {
		fmt.Println("couldn't convert to pod")
	} else {
		pod.SetLabels(map[string]string{"test": "test"}) // modify the object to trigger MutationDetector
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
	log.Printf(eventHandlerMessage, "handlDelete", key)
}

func getKeyFromObj(obj interface{}) string {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		log.Printf("failed to get key from the cache %s\n", err.Error())
		return ""
	}
	return key
}
