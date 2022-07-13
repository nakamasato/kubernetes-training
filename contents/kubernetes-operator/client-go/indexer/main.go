package main

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

func main() {

	// Create index
	indexer := cache.NewIndexer(
		cache.MetaNamespaceKeyFunc, // Use <namespace>/<name> as a key if <namespace> exists, otherwise <name>
		cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, // default index function that indexes based on an object's namespace
	)

	// Add object to index
	err := indexer.Add(getDeployment("test"))
	if err != nil {
		fmt.Println(err)
	}

	// List indexer
	objs := indexer.List()
	fmt.Printf("indexer.List got %d objects\n", len(objs))

	// Get object by key
	obj, exists, err := indexer.GetByKey("default/test")
	if err != nil {
		fmt.Printf("indexer.GetByKey(default/test) failed, %v\n", err)
	} else if !exists {
		fmt.Println("indexer.GetByKey(default/test) doesn't exist")
	}
	fmt.Printf("indexer.GetByKey(default/test) got result. key: %s\n", obj.(*appsv1.Deployment).Name)
}

func getDeployment(name string) *appsv1.Deployment {
	replicas := int32(1)
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
			Labels:    map[string]string{"watch": "true"},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx",
						},
					},
				},
			},
		},
	}
}
