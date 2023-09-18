package main

import (
	"errors"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

func main() {
	indexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
	fifo := cache.NewDeltaFIFOWithOptions(cache.DeltaFIFOOptions{
		KnownObjects:          indexer,
		EmitDeltaTypeReplaced: true,
	})
	stopCh := make(chan struct{})
	go func() {
		<-stopCh
		fifo.Close()
	}()
	// In informer, Reflector update DeltaFIFO
	// https://github.com/kubernetes/client-go/blob/ee1a5aaf793a9ace9c433f5fb26a19058ed5f37c/tools/cache/reflector.go#L460-L538
	err := fifo.Add(newPodWithoutContainer())
	if err != nil {
		fmt.Println("failed to DeltaFIFO.Add")
	}
	err = fifo.Update(newPodWithContainer())
	if err != nil {
		fmt.Println("failed to DeltaFIFO.Update")
	}
	err = fifo.Delete(newPodWithoutContainer())
	if err != nil {
		fmt.Println("failed to DeltaFIFO.Delete")
	}

	// Informer handleDeltas Pop()
	// https://github.com/kubernetes/client-go/blob/ee1a5aaf793a9ace9c433f5fb26a19058ed5f37c/tools/cache/controller.go#L182-L195
	// cannot use process (value of type func(obj interface{}) error) as "k8s.io/client-go/tools/cache".PopProcessFunc value in argument to fifo.Pop
	fifo.Pop(process)
}

func newPodWithoutContainer() *corev1.Pod {
	return &corev1.Pod{
		TypeMeta:   metav1.TypeMeta{Kind: "Pod", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "test-pod"},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{},
		},
		Status: corev1.PodStatus{},
	}
}

func newPodWithContainer() *corev1.Pod {
	return &corev1.Pod{
		TypeMeta:   metav1.TypeMeta{Kind: "Pod", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "test-pod"},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Image: "nginx",
					Name:  "nginx",
				},
			},
		},
		Status: corev1.PodStatus{},
	}
}

// https://github.com/kubernetes/client-go/blame/08f892964c345b3d94d78992b4e924cf9fa7f98a/tools/cache/fifo.go#L28
func process(obj interface{}, isInInitialList bool) error { // type PopProcessFunc func(obj interface{}, isInInitialList bool) error
	if deltas, ok := obj.(cache.Deltas); ok {
		return processDeltas(deltas)
	}
	return errors.New("object given as Process argument is not Deltas")
}

func processDeltas(
	deltas cache.Deltas,
) error {
	for _, d := range deltas {
		switch d.Type {
		case cache.Sync, cache.Replaced, cache.Added, cache.Updated:
			fmt.Println("OnUpdate or OnAdd")
		case cache.Deleted:
			fmt.Println("OnDelete")
		}
	}
	return nil
}
