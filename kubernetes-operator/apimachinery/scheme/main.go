package main

import (
	"fmt"

	scheme "github.com/nakamasato/kubernetes-training/contents/kubernetes-operator/apimachinery/scheme/scheme"
)

func main() {
	fmt.Println("main")
	for k, v := range scheme.Scheme.AllKnownTypes() {
		fmt.Printf("GroupVersionKind[Group:%s\tVersion:%s\tKind:%s], reflect.Type: %v\n", k.Group, k.Version, k.Kind, v)
	}
	fmt.Println("Those scheme are set in scheme/register.go")
}
