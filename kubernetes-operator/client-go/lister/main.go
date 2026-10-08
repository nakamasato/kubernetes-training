package main

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	appsv1lister "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
)

func main() {

	// Indexer (for more details, you can check ../indexer)
	fmt.Println("-------------- Indexer ----------------")
	indexer := cache.NewIndexer(
		cache.MetaNamespaceKeyFunc,
		cache.Indexers{
			cache.NamespaceIndex: cache.MetaNamespaceIndexFunc,
			"labels":             indexFuncByLabels,
		},
	)
	// Add deployment with label
	err := indexer.Add(getDeployment("deployment-with-label", map[string]string{"watch": "true"}))
	if err != nil {
		fmt.Println(err)
	}
	// Add deployment without label
	err = indexer.Add(getDeployment("deployment-without-label", map[string]string{}))
	if err != nil {
		fmt.Println(err)
	}

	// Lister
	fmt.Println("-------------- Lister ----------------")
	deploymentLister := appsv1lister.NewDeploymentLister(indexer)
	selector := labels.SelectorFromSet(labels.Set{"watch": "true"})
	ret, _ := deploymentLister.List(selector)

	for _, deploy := range ret {
		fmt.Println(deploy.Name)
	}
}

func getDeployment(name string, labels map[string]string) *appsv1.Deployment {
	replicas := int32(1)
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
			Labels:    labels,
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

func indexFuncByLabels(obj interface{}) ([]string, error) {
	meta, err := meta.Accessor(obj)
	if err != nil {
		return []string{""}, fmt.Errorf("object has no meta: %v", err)
	}
	keys := []string{}
	for key := range meta.GetLabels() {
		keys = append(keys, key)
	}
	return keys, nil
}
