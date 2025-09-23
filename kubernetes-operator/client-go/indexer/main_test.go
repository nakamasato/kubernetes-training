package main

import (
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/tools/cache"
)

func TestIndexer(t *testing.T) {
	// Given
	indexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
	deploymentName := "test"
	namespace := "default"
	namespaceName := namespace + "/" + deploymentName

	// When: Add
	err := indexer.Add(getDeployment(deploymentName))
	if err != nil {
		t.Errorf("indexer.Add failed %v\n", err)
	}

	// Expect the length of objects is one
	objs := indexer.List()
	if len(objs) != 1 {
		t.Errorf("want 1 got %d\n", len(objs))
	}

	// enable to get the object by <namespace>/<name>
	obj, exists, err := indexer.GetByKey(namespaceName)
	if err != nil {
		t.Errorf("failed to GetByKey, %v\n", err)
	} else if !exists {
		t.Error("GetByKey doesn't exist")
	} else if obj.(*appsv1.Deployment).Name != deploymentName {
		t.Errorf("indexer.GetByKey want '%s', got %s\n", namespaceName, obj.(*appsv1.Deployment).Name)
	}
}

func TestThreadSafeMap(t *testing.T) {
	cacheStorage := cache.NewThreadSafeStore(cache.Indexers{}, cache.Indices{})

	// example1: String
	keyStr := "keyStr"
	valStr := "valStr"
	cacheStorage.Add(keyStr, valStr)

	item, exists := cacheStorage.Get(keyStr)
	if !exists {
		t.Errorf("not exist %s", keyStr)
	}
	if item != valStr {
		t.Errorf("want %s, got %s", valStr, item)
	}

	// example2: User
	keyUser := "userId"
	valUser := User{Name: "John", Age: 10}
	cacheStorage.Add(keyUser, valUser)
	item, exists = cacheStorage.Get(keyUser)
	if !exists {
		t.Errorf("not exist %s", keyUser)
	}
	if item != valUser {
		t.Errorf("want %s, got %s", valUser, item)
	}

	// example3: Update key with different type
	cacheStorage.Add(keyStr, valUser)
	item, exists = cacheStorage.Get(keyStr)
	if !exists {
		t.Errorf("not exist %s", keyStr)
	}
	if item != valUser {
		t.Errorf("want %s, got %s", valUser, item)
	}

	// example4: Add Indexer
	cacheStorage.Delete(keyStr) // need to clean up all items before adding indexers
	cacheStorage.Delete(keyUser)
	indexers := map[string]cache.IndexFunc{"nameIndexer": nameIndexer}
	cacheStorage.AddIndexers(indexers)
	cacheStorage.Add(keyUser, valUser)
	res, err := cacheStorage.ByIndex("nameIndexer", "John") // indexed by User.Name
	if err != nil {
		t.Errorf("err: %v", err)
	}
	if len(res) != 1 {
		t.Errorf("length should be 1, got %d", len(res))
	}
	if res[0] != valUser {
		t.Errorf("want %s, got %s", valUser, res[0])
	}

	res, err = cacheStorage.Index("nameIndexer", valUser) // obj -> indexed values -> convert indexed values to the keys of items -> get value from items
	if err != nil {
		t.Errorf("err: %v", err)
	}
	if len(res) != 1 {
		t.Errorf("length should be 1, got %d", len(res))
	}
	if res[0] != valUser {
		t.Errorf("want %s, got %s", valUser, res[0])
	}
}
