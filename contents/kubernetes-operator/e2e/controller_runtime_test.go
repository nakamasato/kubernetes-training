package e2e

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestControllerRuntimeReplicaSet(t *testing.T) {
	client := testClient(t)
	ns := createNamespace(t, client)
	root := repositoryRoot(t)
	binary := buildExample(t, root, "./contents/kubernetes-operator/controller-runtime/example-controller", "replicaset-controller")
	logPath, stop := startExample(t, binary)
	defer stop()

	ctx := context.Background()
	rs, err := client.AppsV1().ReplicaSets(ns).Create(ctx, &appsv1.ReplicaSet{
		ObjectMeta: metav1.ObjectMeta{Name: "counted"},
		Spec: appsv1.ReplicaSetSpec{
			Replicas: int32ptr(2),
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": "operator-e2e"}},
			Template: corev1.PodTemplateSpec{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": "operator-e2e"}}, Spec: corev1.PodSpec{Containers: []corev1.Container{{Name: "pause", Image: "registry.k8s.io/pause:3.10"}}}},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		t.Fatal(err)
	}
	waitFor(t, "ReplicaSet pod-count=2", func() (bool, error) {
		current, err := client.AppsV1().ReplicaSets(ns).Get(ctx, rs.Name, metav1.GetOptions{})
		return current.Labels["pod-count"] == "2", err
	})

	rs, err = client.AppsV1().ReplicaSets(ns).Get(ctx, rs.Name, metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
	rs.Spec.Replicas = int32ptr(1)
	if _, err = client.AppsV1().ReplicaSets(ns).Update(ctx, rs, metav1.UpdateOptions{}); err != nil {
		t.Fatal(err)
	}
	waitFor(t, "ReplicaSet pod-count=1", func() (bool, error) {
		current, err := client.AppsV1().ReplicaSets(ns).Get(ctx, rs.Name, metav1.GetOptions{})
		return current.Labels["pod-count"] == "1", err
	})

	rs, err = client.AppsV1().ReplicaSets(ns).Get(ctx, rs.Name, metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
	rs.Spec.Replicas = int32ptr(0)
	if _, err = client.AppsV1().ReplicaSets(ns).Update(ctx, rs, metav1.UpdateOptions{}); err != nil {
		t.Fatal(err)
	}
	waitFor(t, "ReplicaSet pod-count=0 after Pods are removed", func() (bool, error) {
		current, err := client.AppsV1().ReplicaSets(ns).Get(ctx, rs.Name, metav1.GetOptions{})
		return current.Labels["pod-count"] == "0", err
	})
	_ = logPath
}

func TestControllerRuntimeCache(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	_, err := client.CoreV1().Pods(metav1.NamespaceDefault).Create(ctx, &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "nginx"},
		Spec:       corev1.PodSpec{Containers: []corev1.Container{{Name: "pause", Image: "registry.k8s.io/pause:3.10"}}},
	}, metav1.CreateOptions{})
	if err != nil {
		// The cleanup from a previous interrupted run may still be completing.
		_ = client.CoreV1().Pods(metav1.NamespaceDefault).Delete(ctx, "nginx", metav1.DeleteOptions{})
		_, err = client.CoreV1().Pods(metav1.NamespaceDefault).Create(ctx, &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "nginx"}, Spec: corev1.PodSpec{Containers: []corev1.Container{{Name: "pause", Image: "registry.k8s.io/pause:3.10"}}}}, metav1.CreateOptions{})
	}
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = client.CoreV1().Pods(metav1.NamespaceDefault).Delete(ctx, "nginx", metav1.DeleteOptions{}) })
	logPath, stop := startExample(t, buildExample(t, repositoryRoot(t), "./contents/kubernetes-operator/controller-runtime/cache", "cache"))
	defer stop()
	waitForLog(t, logPath, "cache is synced")
	waitForLog(t, logPath, "cached Pod: default/nginx")
}

func TestControllerRuntimeManager(t *testing.T) {
	client := testClient(t)
	ns := createNamespace(t, client)
	logPath, stop := startExample(t, buildExample(t, repositoryRoot(t), "./contents/kubernetes-operator/controller-runtime/manager", "manager"))
	defer stop()
	waitForLog(t, logPath, "Starting Controller")
	ctx := context.Background()
	_, err := client.CoreV1().Pods(ns).Create(ctx, &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "manager-e2e"}, Spec: corev1.PodSpec{Containers: []corev1.Container{{Name: "pause", Image: "registry.k8s.io/pause:3.10"}}}}, metav1.CreateOptions{})
	if err != nil {
		t.Fatal(err)
	}
	waitForLog(t, logPath, "RunnableFunc is called")
	waitForLog(t, logPath, "podReconciler is called")
}

func buildExample(t *testing.T, root, packagePath, name string) string {
	t.Helper()
	binary := filepath.Join(t.TempDir(), name)
	cmd := exec.Command("go", "build", "-o", binary, packagePath)
	cmd.Dir = root
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build %s: %v\n%s", packagePath, err, output)
	}
	return binary
}

func startExample(t *testing.T, binary string, args ...string) (string, func()) {
	t.Helper()
	logPath := exampleLogPath(t, filepath.Base(binary)+".log")
	logFile, err := os.Create(logPath)
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(binary, args...)
	cmd.Stdout, cmd.Stderr = logFile, logFile
	if err := cmd.Start(); err != nil {
		_ = logFile.Close()
		t.Fatal(err)
	}
	return logPath, func() {
		_ = cmd.Process.Signal(os.Interrupt)
		finished := make(chan struct{})
		go func() {
			_, _ = cmd.Process.Wait()
			close(finished)
		}()
		select {
		case <-finished:
		case <-time.After(5 * time.Second):
			_ = cmd.Process.Kill()
			<-finished
		}
		_ = logFile.Close()
	}
}

func int32ptr(value int32) *int32 { return &value }
