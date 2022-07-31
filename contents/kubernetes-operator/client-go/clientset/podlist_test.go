package main_test

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var clientset *kubernetes.Clientset

func TestMain(m *testing.M) {
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
	clientset, _ = kubernetes.NewForConfig(config)

	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestPodList(t *testing.T) {
	// Get list of pod objects
	pods, _ := clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})

	// show pod object to stdout
	if len(pods.Items) < 1 {
		t.Errorf("There should be at least one Pod runnning (pod: %d)", len(pods.Items))
	}
}
