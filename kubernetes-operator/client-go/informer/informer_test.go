package main

import (
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

func TestHandleAddPreservesCachedPod(t *testing.T) {
	pod := &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "example", Namespace: "default", Labels: map[string]string{"original": "value"}}}
	handleAdd(pod)
	if len(pod.Labels) != 1 || pod.Labels["original"] != "value" {
		t.Fatalf("handler mutated shared object: %v", pod.Labels)
	}
}

func TestDeletionKey(t *testing.T) {
	pod := &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "example", Namespace: "default"}}
	for _, obj := range []interface{}{pod, cache.DeletedFinalStateUnknown{Key: "default/example", Obj: pod}} {
		if key := getKeyFromObj(obj); key != "default/example" {
			t.Errorf("key = %q, want default/example", key)
		}
	}
}
