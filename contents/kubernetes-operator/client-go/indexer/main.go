package main

import (
	"fmt"
	"errors"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

func main() {
	exampleWithIndexer()
	exampleWithThreadSafeMap()
}

func exampleWithIndexer() {
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

type User struct {
	Name string
	Age int
}

func (u User) String() string {
	return fmt.Sprintf("User(Name=%s, Age=%d)", u.Name, u.Age)
}

func exampleWithThreadSafeMap() {
	// func NewIndexer(keyFunc KeyFunc, indexers Indexers) Indexer { https://github.com/kubernetes/client-go/blob/ee1a5aaf793a9ace9c433f5fb26a19058ed5f37c/tools/cache/store.go#L266-L271
	// 	return &cache{
	// 		cacheStorage: NewThreadSafeStore(indexers, Indices{}),
	// 		keyFunc:      keyFunc,
	// 	}
	// }
	// Inside NewIndexer, this is called to initialize cacheStorage
	cacheStorage := cache.NewThreadSafeStore(cache.Indexers{}, cache.Indices{})

	// example1: String
	keyStr := "keyStr"
	valStr := "valStr"
	cacheStorage.Add(keyStr, valStr)

	item, exists := cacheStorage.Get(keyStr)
	if exists {
		fmt.Printf("ThreadSafeStore.Get got %s (=valStr) for key %s\n", item, keyStr)
	}

	// example2: User
	keyUser := "userId"
	valUser := User{Name: "John", Age: 10}
	cacheStorage.Add(keyUser, valUser)
	item, exists = cacheStorage.Get(keyUser)
	if exists {
		fmt.Printf("ThreadSafeStore.Get got %s for key %s\n", item, keyUser)
	}

	// example3: Update key with different type
	cacheStorage.Add(keyStr, valUser)
	item, exists = cacheStorage.Get(keyStr)
	if exists {
		fmt.Printf("ThreadSafeStore.Get got %s for key %s\n", item, keyStr)
	}

	// example4: Add Indexer
	cacheStorage.Delete(keyStr) // need to clean up all items before adding indexers
	cacheStorage.Delete(keyUser)
	indexers := map[string]cache.IndexFunc{"nameIndexer": nameIndexer}
	cacheStorage.AddIndexers(indexers)
	cacheStorage.Add(keyUser, valUser)
	res, err := cacheStorage.ByIndex("nameIndexer", "John") // indexed by User.Name
	if err == nil {
		fmt.Printf("cacheStorage.ByIndex got %d result\n", len(res))
		for _, r := range res {
			fmt.Printf("cacheStorage.ByIndex got %s", r)
		}
	}
	res, err = cacheStorage.Index("nameIndexer", valUser) // obj -> indexed values -> convert indexed values to the keys of items -> get value from items
	if err == nil {
		fmt.Printf("cacheStorage.Index got %d result\n", len(res))
		for _, r := range res {
			fmt.Printf("cacheStorage.Index got %s", r)
		}
	}
}

func nameIndexer(obj interface{}) ([]string, error) {
	switch obj.(type) {
	case string:
		return []string{obj.(string)}, nil
	case User:
		return []string{obj.(User).Name}, nil
	default:
		return []string{}, errors.New("nameIndexer error")
	}
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
