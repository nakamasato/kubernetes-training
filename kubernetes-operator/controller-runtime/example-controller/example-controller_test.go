package main

import (
	"context"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
)

func TestReconcilePodCount(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	if err := appsv1.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	if err := corev1.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	rs := &appsv1.ReplicaSet{
		ObjectMeta: metav1.ObjectMeta{Name: "example", Namespace: "default", UID: "owner"},
		Spec: appsv1.ReplicaSetSpec{
			Selector: &metav1.LabelSelector{MatchExpressions: []metav1.LabelSelectorRequirement{
				{Key: "app", Operator: metav1.LabelSelectorOpIn, Values: []string{"demo"}},
			}},
		},
	}
	pod := func(name, namespace, label string, owner types.UID) *corev1.Pod {
		p := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace, Labels: map[string]string{"app": label}}}
		if owner != "" {
			ref := metav1.NewControllerRef(rs, appsv1.SchemeGroupVersion.WithKind("ReplicaSet"))
			ref.UID = owner
			p.OwnerReferences = []metav1.OwnerReference{*ref}
		}
		return p
	}
	owned := pod("owned", "default", "demo", rs.UID)
	patches := 0
	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(rs, owned,
		pod("unowned", "default", "demo", ""),
		pod("other-owner", "default", "demo", "different"),
		pod("other-namespace", "other", "demo", rs.UID),
		pod("not-selected", "default", "other", rs.UID),
	).WithInterceptorFuncs(interceptor.Funcs{
		Patch: func(ctx context.Context, c client.WithWatch, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
			patches++
			return c.Patch(ctx, obj, patch, opts...)
		},
	}).Build()
	r := &ReplicaSetReconciler{Client: c}
	req := ctrl.Request{NamespacedName: client.ObjectKeyFromObject(rs)}
	check := func(want string) {
		t.Helper()
		if _, err := r.Reconcile(ctx, req); err != nil {
			t.Fatal(err)
		}
		got := &appsv1.ReplicaSet{}
		if err := c.Get(ctx, req.NamespacedName, got); err != nil {
			t.Fatal(err)
		}
		if got.Labels["pod-count"] != want {
			t.Fatalf("pod-count = %q, want %q", got.Labels["pod-count"], want)
		}
	}
	check("1") // nil labels, matchExpressions, ownership and namespace filtering
	check("1")
	if patches != 1 {
		t.Fatalf("unchanged count triggered another patch: %d", patches)
	}
	if err := c.Delete(ctx, owned); err != nil {
		t.Fatal(err)
	}
	check("0")
	if err := c.Delete(ctx, rs); err != nil {
		t.Fatal(err)
	}
	if _, err := r.Reconcile(ctx, req); err != nil {
		t.Fatalf("deleted ReplicaSet should not retry: %v", err)
	}
}
