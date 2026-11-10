package e2e

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/e2e-framework/klient"
)

func TestClientGoClientset(t *testing.T) {
	client := testClient(t)
	ns := createNamespace(t, client)
	name := "clientset-e2e"
	_, err := client.CoreV1().Pods(ns).Create(context.Background(), &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Spec:       corev1.PodSpec{Containers: []corev1.Container{{Name: "pause", Image: "registry.k8s.io/pause:3.10"}}},
	}, metav1.CreateOptions{})
	if err != nil {
		t.Fatal(err)
	}
	output := runExample(t, buildExample(t, repositoryRoot(t), "./contents/kubernetes-operator/client-go/clientset", "clientset"), "-kubeconfig", kubeconfig(t))
	if !strings.Contains(output, ns+"/"+name) {
		t.Fatalf("clientset output does not contain %s:\n%s", ns+"/"+name, output)
	}
}

func TestClientGoListerWatcher(t *testing.T) {
	client := testClient(t)
	logPath, stop := startExample(t, buildExample(t, repositoryRoot(t), "./contents/kubernetes-operator/client-go/listerwatcher", "listerwatcher"), "-kubeconfig", kubeconfig(t))
	defer stop()
	ctx := context.Background()
	name := "listerwatcher-e2e"
	ns := metav1.NamespaceDefault
	_ = client.CoreV1().Pods(ns).Delete(ctx, name, metav1.DeleteOptions{})
	waitForLog(t, logPath, "resourceVersion:")
	_, err := client.CoreV1().Pods(ns).Create(ctx, &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: name}, Spec: corev1.PodSpec{Containers: []corev1.Container{{Name: "pause", Image: "registry.k8s.io/pause:3.10"}}}}, metav1.CreateOptions{})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = client.CoreV1().Pods(ns).Delete(ctx, name, metav1.DeleteOptions{}) })
	waitForLog(t, logPath, "event: ADDED")
	pod, err := client.CoreV1().Pods(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
	pod.Labels = map[string]string{"changed": "true"}
	if _, err = client.CoreV1().Pods(ns).Update(ctx, pod, metav1.UpdateOptions{}); err != nil {
		t.Fatal(err)
	}
	waitForLog(t, logPath, "event: MODIFIED")
	if err := client.CoreV1().Pods(ns).Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
		t.Fatal(err)
	}
	waitForLog(t, logPath, "event: DELETED")
}

func TestClientGoInformerCRUD(t *testing.T) {
	ctx := context.Background()
	client := testClient(t)
	ns := createNamespace(t, client)

	root := repositoryRoot(t)
	binary := filepath.Join(t.TempDir(), "informer")
	run := exec.Command("go", "build", "-o", binary, "./contents/kubernetes-operator/client-go/informer")
	run.Dir = root
	if output, err := run.CombinedOutput(); err != nil {
		t.Fatalf("build informer: %v\n%s", err, output)
	}

	logPath := exampleLogPath(t, "informer.log")
	logFile, err := os.Create(logPath)
	if err != nil {
		t.Fatal(err)
	}
	process := exec.Command(binary, "-kubeconfig", kubeconfig(t))
	process.Stdout = logFile
	process.Stderr = logFile
	if err := process.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = process.Process.Signal(os.Interrupt)
		_ = process.Wait()
		_ = logFile.Close()
	}()

	name := "informer-e2e"
	key := ns + "/" + name
	_, err = client.CoreV1().Pods(ns).Create(ctx, &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Spec:       corev1.PodSpec{Containers: []corev1.Container{{Name: "pause", Image: "registry.k8s.io/pause:3.10"}}},
	}, metav1.CreateOptions{})
	if err != nil {
		t.Fatal(err)
	}
	waitForLog(t, logPath, "cache is synced")
	waitForLog(t, logPath, "handleAdd is called for Pod (key: "+key+")")

	pod, err := client.CoreV1().Pods(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
	pod.Labels = map[string]string{"changed": "true"}
	if _, err = client.CoreV1().Pods(ns).Update(ctx, pod, metav1.UpdateOptions{}); err != nil {
		t.Fatal(err)
	}
	waitForLog(t, logPath, "handleUpdate is called for Pod (key: "+key+")")
	if err := client.CoreV1().Pods(ns).Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
		t.Fatal(err)
	}
	waitForLog(t, logPath, "handleDelete is called for Pod (key: "+key+")")
}

func testClient(t *testing.T) kubernetes.Interface {
	t.Helper()
	frameworkClient, err := klient.NewWithKubeConfigFile(kubeconfig(t))
	if err != nil {
		t.Fatal(err)
	}
	client, err := kubernetes.NewForConfig(frameworkClient.RESTConfig())
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func kubeconfig(t *testing.T) string {
	t.Helper()
	if value := os.Getenv("KUBECONFIG"); value != "" {
		return value
	}
	t.Skip("KUBECONFIG is not set; live-cluster E2E is run by an E2E target")
	return ""
}

func createNamespace(t *testing.T, client kubernetes.Interface) string {
	t.Helper()
	ns, err := client.CoreV1().Namespaces().Create(context.Background(), &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{GenerateName: "operator-e2e-"}}, metav1.CreateOptions{})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = client.CoreV1().Namespaces().Delete(context.Background(), ns.Name, metav1.DeleteOptions{}) })
	return ns.Name
}

func repositoryRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", "..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return root
}

func exampleLogPath(t *testing.T, name string) string {
	t.Helper()
	if dir := os.Getenv("E2E_ARTIFACTS"); dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		return filepath.Join(dir, name)
	}
	return filepath.Join(t.TempDir(), name)
}

func waitForLog(t *testing.T, path, want string) {
	t.Helper()
	deadline := time.Now().Add(90 * time.Second)
	for time.Now().Before(deadline) {
		data, _ := os.ReadFile(path)
		if strings.Contains(string(data), want) {
			return
		}
		time.Sleep(500 * time.Millisecond)
	}
	data, _ := os.ReadFile(path)
	t.Fatalf("timed out waiting for %q in %s\n%s", want, path, data)
}

func waitFor(t *testing.T, description string, condition func() (bool, error)) {
	t.Helper()
	deadline := time.Now().Add(90 * time.Second)
	var lastErr error
	for time.Now().Before(deadline) {
		ok, err := condition()
		if ok && err == nil {
			return
		}
		lastErr = err
		time.Sleep(500 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %s: %v", description, lastErr)
}

func runExample(t *testing.T, binary string, args ...string) string {
	t.Helper()
	cmd := exec.Command(binary, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run %s: %v\n%s", binary, err, output)
	}
	return string(output)
}
