package main

import (
	"fmt"
	"testing"

	"k8s.io/apimachinery/pkg/labels"
	appsv1lister "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
)

func TestLister(t *testing.T) {
	// Prepare Indexer
	indexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
	// Add deployment with label
	err := indexer.Add(getDeployment("deployment-with-label", map[string]string{"watch": "true"}))
	if err != nil {
		t.Errorf("indexer.Add failed %v\n", err)
	}
	// Add deployment without label
	err = indexer.Add(getDeployment("deployment-without-label", map[string]string{}))
	if err != nil {
		t.Errorf("indexer.Add failed %v\n", err)
	}

	// Prepare Lister
	deploymentLister := appsv1lister.NewDeploymentLister(indexer)
	selector := labels.SelectorFromSet(labels.Set{"watch": "true"})
	ret, _ := deploymentLister.List(selector)

	if len(ret) != 1 {
		t.Errorf("want %d, got %d", 1, len(ret))
	}
}

func TestListerWithIndex(t *testing.T) {
	// Prepare Indexer
	indexer := cache.NewIndexer(
		cache.MetaNamespaceKeyFunc,
		cache.Indexers{
			cache.NamespaceIndex: cache.MetaNamespaceIndexFunc,
			"labels":             indexFuncByLabels,
		},
	)

	// Add deployment with label
	NUM_OF_DEPLOYMENT := 1000
	for i := 0; i < NUM_OF_DEPLOYMENT; i++ {
		err := indexer.Add(getDeployment(fmt.Sprintf("deployment-with-label-%d", i), map[string]string{"watch": "true", "id": fmt.Sprintf("%d", i)}))
		if err != nil {
			t.Errorf("indexer.Add failed %v\n", err)
		}
	}
	// Add deployment without label
	for i := 0; i < NUM_OF_DEPLOYMENT; i++ {
		err := indexer.Add(getDeployment(fmt.Sprintf("deployment-without-label-%d", i), map[string]string{}))
		if err != nil {
			t.Errorf("indexer.Add failed %v\n", err)
		}
	}

	// Prepare Lister
	deploymentListerWithIndex := appsv1lister.NewDeploymentLister(indexer)
	ret, _ := deploymentListerWithIndex.List(labels.SelectorFromSet(labels.Set{"id": "1"}))

	if expected, actual := 1, len(ret); expected != actual {
		t.Errorf("want %d, got %d", expected, actual)
	}

	ret, _ = deploymentListerWithIndex.List(labels.SelectorFromSet(labels.Set{"watch": "true"}))

	if expected, actual := NUM_OF_DEPLOYMENT, len(ret); expected != actual {
		t.Errorf("want %d, got %d", expected, actual)
	}
}

func TestListerWithoutIndex(t *testing.T) {
	// Prepare Indexer
	indexer := cache.NewIndexer(
		cache.MetaNamespaceKeyFunc,
		cache.Indexers{},
	)

	// Add deployment with label
	NUM_OF_DEPLOYMENT := 1000
	for i := 0; i < NUM_OF_DEPLOYMENT; i++ {
		err := indexer.Add(getDeployment(fmt.Sprintf("deployment-with-label-%d", i), map[string]string{"watch": "true", "id": fmt.Sprintf("%d", i)}))
		if err != nil {
			t.Errorf("indexer.Add failed %v\n", err)
		}
	}
	// Add deployment without label
	for i := 0; i < NUM_OF_DEPLOYMENT; i++ {
		err := indexer.Add(getDeployment(fmt.Sprintf("deployment-without-label-%d", i), map[string]string{}))
		if err != nil {
			t.Errorf("indexer.Add failed %v\n", err)
		}
	}
	// Prepare Lister
	deploymentListerWithoutIndex := appsv1lister.NewDeploymentLister(indexer)
	ret, _ := deploymentListerWithoutIndex.List(labels.SelectorFromSet(labels.Set{"id": "0"}))

	if expected, actual := 1, len(ret); expected != actual {
		t.Errorf("want %d, got %d", expected, actual)
	}

	ret, _ = deploymentListerWithoutIndex.List(labels.SelectorFromSet(labels.Set{"watch": "true"}))

	if expected, actual := NUM_OF_DEPLOYMENT, len(ret); expected != actual {
		t.Errorf("want %d, got %d", expected, actual)
	}
}
