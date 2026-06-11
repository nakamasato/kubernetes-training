package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

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
	config, _ := clientcmd.BuildConfigFromFlags("", *kubeconfig)

	// get clientset for kubernetes resources
	clientset, _ := kubernetes.NewForConfig(config)

	// Get list of pod objects
	pods, _ := clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})

	// show pod object to stdout
	for i, pod := range pods.Items {
		fmt.Printf("[Pod Name %d] %s\n", i, pod.GetName())
	}
}
