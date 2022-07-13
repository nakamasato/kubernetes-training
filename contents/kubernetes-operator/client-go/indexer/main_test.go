package main

import (
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/tools/cache"
)

func TestIndexer(t *testing.T) {
	indexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
	deploymentName := "test"
	namespace := "default"
	namespaceName := namespace + "/" + deploymentName
	err := indexer.Add(getDeployment(deploymentName))
	if err != nil {
		t.Errorf("indexer.Add failed %v\n", err)
	}

	objs := indexer.List()
	if len(objs) != 1 {
		t.Errorf("want 1 got %d\n", len(objs))
	}

	obj, exists, err := indexer.GetByKey(namespaceName)
	if err != nil {
		t.Errorf("failed to GetByKey, %v\n", err)
	} else if !exists {
		t.Error("GetByKey doesn't exist")
	} else if obj.(*appsv1.Deployment).Name != deploymentName {
		t.Errorf("indexer.GetByKey want '%s', got %s\n", namespaceName, obj.(*appsv1.Deployment).Name)
	}
}
