package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var defaultKubeConfigPath string
	if home := homedir.HomeDir(); home != "" {
		// build kubeconfig path from $HOME dir
		defaultKubeConfigPath = filepath.Join(home, ".kube", "config")
	}

	// set kubeconfig flag
	kubeconfig := flag.String("kubeconfig", defaultKubeConfigPath, "kubeconfig config file")
	flag.Parse()

	// retrieve kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	// get clientset for kubernetes resources
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	// Get list of pod objects
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pods, err := clientset.CoreV1().Pods(metav1.NamespaceAll).List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	// show pod object to stdout
	for i, pod := range pods.Items {
		fmt.Printf("[Pod %d] %s/%s\n", i, pod.Namespace, pod.Name)
	}
}
