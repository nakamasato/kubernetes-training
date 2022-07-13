package main

import (
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
