package main

import (
	"flag"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog/v2"
)

func main() {
	var kubeconfig string
	var master string

	kubeConfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")
	flag.StringVar(&kubeconfig, "kubeconfig", kubeConfigPath, "absolute path to the kubeconfig file")
	flag.StringVar(&master, "master", "", "master url")
	flag.Parse()
	// creates the connection
	config, err := clientcmd.BuildConfigFromFlags(master, kubeconfig)
	if err != nil {
		klog.Fatal(err)
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatal(err)
	}

	podListWatcher := cache.NewListWatchFromClient(
		clientset.CoreV1().RESTClient(),
		"pods",
		v1.NamespaceDefault,
		fields.Everything(),
	)

	// List
	list, err := podListWatcher.List(metav1.ListOptions{}) // returns runtime.Object
	if err != nil {
		klog.Fatal(err)
	}

	listMetaInterface, err := meta.ListAccessor(list)
	if err != nil {
		klog.Fatal(err)
	}

	resourceVersion := listMetaInterface.GetResourceVersion()
	klog.Infof("resourceVersion: %s", resourceVersion)

	items, err := meta.ExtractList(list)
	if err != nil {
		klog.Fatal(err)
	}

	klog.Infof("items: %d", len(items))

	// Watch
	w, err := podListWatcher.Watch(metav1.ListOptions{}) // returns watch.Interface
	if err != nil {
		klog.Fatal(err)
	}
loop:
	for {
		select {
		case event, ok := <-w.ResultChan():
			if !ok {
				break loop
			}

			meta, err := meta.Accessor(event.Object)
			if err != nil {
				continue
			}
			resourceVersion := meta.GetResourceVersion()
			klog.Infof("event: %s, resourceVersion: %s", event.Type, resourceVersion)
		}
	}
}
