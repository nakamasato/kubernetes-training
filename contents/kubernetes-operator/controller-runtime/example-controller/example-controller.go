package main

import (
	"context"
	"os"
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
)

func main() {
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))
	var log = ctrl.Log.WithName("builder-examples")

	manager, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{Metrics: metricsserver.Options{BindAddress: "0"}})
	if err != nil {
		log.Error(err, "could not create manager")
		os.Exit(1)
	}

	err = ctrl.
		NewControllerManagedBy(manager). // Create the Controller
		For(&appsv1.ReplicaSet{}).       // ReplicaSet is the Application API
		Owns(&corev1.Pod{}).             // ReplicaSet owns Pods created by it
		Complete(&ReplicaSetReconciler{Client: manager.GetClient()})
	if err != nil {
		log.Error(err, "could not create controller")
		os.Exit(1)
	}

	if err := manager.Start(ctrl.SetupSignalHandler()); err != nil {
		log.Error(err, "could not start manager")
		os.Exit(1)
	}
}

// ReplicaSetReconciler is a simple Controller example implementation.
type ReplicaSetReconciler struct {
	client.Client
}

// Implement the business logic:
// This function will be called when there is a change to a ReplicaSet or a Pod with an OwnerReference
// to a ReplicaSet.
//
// * Read the ReplicaSet
// * Read the Pods
// * Set a Label on the ReplicaSet with the Pod count.
func (a *ReplicaSetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	rs := &appsv1.ReplicaSet{}
	err := a.Get(ctx, req.NamespacedName, rs)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	selector, err := metav1.LabelSelectorAsSelector(rs.Spec.Selector)
	if err != nil {
		return ctrl.Result{}, err
	}
	pods := &corev1.PodList{}
	err = a.List(ctx, pods, client.InNamespace(req.Namespace), client.MatchingLabelsSelector{Selector: selector})
	if err != nil {
		return ctrl.Result{}, err
	}

	count := 0
	for i := range pods.Items {
		if metav1.IsControlledBy(&pods.Items[i], rs) {
			count++
		}
	}
	podCount := strconv.Itoa(count)
	if rs.Labels["pod-count"] == podCount {
		return ctrl.Result{}, nil
	}
	original := rs.DeepCopy()
	if rs.Labels == nil {
		rs.Labels = make(map[string]string)
	}
	rs.Labels["pod-count"] = podCount
	err = a.Patch(ctx, rs, client.MergeFrom(original))
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}
