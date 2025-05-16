/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	builderplaygroundv1alpha1 "github.com/flashbots/builder-playground-operator/api/v1alpha1"
)

// BuilderPlaygroundDeploymentReconciler reconciles a BuilderPlaygroundDeployment object
type BuilderPlaygroundDeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=builderplayground.flashbots.io,resources=builderplaygrounddeployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=builderplayground.flashbots.io,resources=builderplaygrounddeployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=builderplayground.flashbots.io,resources=builderplaygrounddeployments/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the BuilderPlaygroundDeployment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.4/pkg/reconcile
func (r *BuilderPlaygroundDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// var builderPlaygroundDeployment  builderplaygroundv1alpha1.BuilderPlaygroundDeployment
	builderPlaygroundDeployment :=  &builderplaygroundv1alpha1.BuilderPlaygroundDeployment {}
	if err := r.Get(ctx, req.NamespacedName, &builderPlaygroundDeployment); err != nil {
	log.Error(err, "Failed to fetch BuilderPlaygroundDeployment")
	return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	deploymentName := builderPlaygroundDeployment.Name + "-statefulset"

	existingSts := &appsv1.StatefulSet{}
	err = r.Get(ctx, client.ObjectKey{Name: epg.Name, Namespace: epg.Namespace}, existingSts)

	if err == nil {
		// Optionally handle updates if needed
		return ctrl.Result{}, nil
	} else if !apierrors.IsNotFound(err) {
		return ctrl.Result{}, err
	}


	sts := generateStatefulSet(builderPlaygroundDeployment, r.Scheme)
	err = ctrl.SetControllerReference(builderPlaygroundDeployment, sts, r.Scheme)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to set controller reference: %w", err)
	}

	err = r.Create(ctx, sts)
	if err != nil {
		return ctrl.Result{}, err
	}	

	return ctrl.Result{}, nil
}


func generateStatefulSet(builderPlaygroundDeployment *builderplaygroundv1alpha1.BuilderPlaygroundDeployment, scheme *runtime.Scheme) *appsv1.StatefulSet {
	labels := map[string]string{"app": builderPlaygroundDeployment.Name}
	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      builderPlaygroundDeployment.Name,
			Namespace: builderPlaygroundDeployment.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: builderPlaygroundDeployment.Name,
			Replicas:    int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{{
						Name: "artifacts-volume",
						VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{
								Path: "/home/ubuntu/my-builder-testnet-2",
								Type: hostPathTypePtr("Directory"),
							},
						},
					}},
					InitContainers: []corev1.Container{{
						Name:    "check-artifacts-dir",
						Image:   "busybox:1.36",
						Command: []string{"sh", "-c", "test -d /artifacts"},
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "artifacts-volume",
							MountPath: "/artifacts",
						}},
					}},
					// Containers: buildContainers(builderPlaygroundDeployment.Spec.Services),
				},
			},
		},
	}
}

func int32Ptr(i int32) *int32 { 
	return &i 
}

func hostPathTypePtr(t string) *corev1.HostPathType {
	hpt := corev1.HostPathType(t)
	return &hpt
}

// SetupWithManager sets up the controller with the Manager.
func (r *BuilderPlaygroundDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&builderplaygroundv1alpha1.BuilderPlaygroundDeployment{}).
		Named("builderplaygrounddeployment").
		Complete(r)
}
