package main

import (
	scheme "github.com/nakamasato/kubernetes-training/contents/kubernetes-operator/apimachinery/scheme/scheme"
	"testing"
)

func TestScheme(t *testing.T) {
	typesNum := len(scheme.Scheme.AllKnownTypes())
	if typesNum != 31 {
		t.Errorf("AllKnownTypes want 31, got %d\n", typesNum) // appsv1 has 31 types
	}
}
