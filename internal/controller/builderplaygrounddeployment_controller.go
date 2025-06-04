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
	"fmt"
	"strings"
	"os"
	"encoding/json"
    "path/filepath"

    "sigs.k8s.io/yaml"
    "k8s.io/apimachinery/pkg/api/resource"

	"k8s.io/apimachinery/pkg/runtime"
	// "k8s.io/apimachinery/pkg/types"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	// "sigs.k8s.io/controller-runtime/pkg/controller"
	// "sigs.k8s.io/controller-runtime/pkg/manager"
	// "sigs.k8s.io/controller-runtime/pkg/reconcile"
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
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete

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
	fmt.Println("✅ start Reconcile ....")

	// var builderPlaygroundDeployment  builderplaygroundv1alpha1.BuilderPlaygroundDeployment
	builderPlaygroundDeployment :=  &builderplaygroundv1alpha1.BuilderPlaygroundDeployment {}
	var err error
	if err = r.Get(ctx, req.NamespacedName, builderPlaygroundDeployment); err != nil {
	log.Error(err, "Failed to fetch BuilderPlaygroundDeployment")
	return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// deploymentName := builderPlaygroundDeployment.Name + "-statefulset"

	existingSts := &appsv1.StatefulSet{}
	err = r.Get(ctx, client.ObjectKey{Name: builderPlaygroundDeployment.Name, Namespace: builderPlaygroundDeployment.Namespace}, existingSts)

	if err == nil {
		// Optionally handle updates if needed
		return ctrl.Result{}, nil
	} else if !apierrors.IsNotFound(err) {
		return ctrl.Result{}, err
	}

	sts := generateStatefulSetForOperator(builderPlaygroundDeployment, r.Scheme)
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

// This function geneate the StatefulSet Structured used for Operator to manage the resources
func generateStatefulSetForOperator(builderPlaygroundDeployment *builderplaygroundv1alpha1.BuilderPlaygroundDeployment, scheme *runtime.Scheme) *appsv1.StatefulSet {
	labels := map[string]string{"app": builderPlaygroundDeployment.Name}
	containers := buildContainers(builderPlaygroundDeployment)

  	fmt.Println("✅ start to geneate statefulset....")

	var volumes []corev1.Volume
	var volumeMounts []corev1.VolumeMount
	var initContainers []corev1.Container
	var volumeClaimTemplates []corev1.PersistentVolumeClaim

	volumeMounts = []corev1.VolumeMount{{
		Name:      "artifacts",
		MountPath: "/artifacts",
	}}

	if builderPlaygroundDeployment.Spec.Storage.Type == "pvc" {
		// Use PVC
		volumes = []corev1.Volume{{
			Name: "artifacts",
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: "artifacts",
				},
			},
		}}

		volumeClaimTemplates = []corev1.PersistentVolumeClaim{{
			ObjectMeta: metav1.ObjectMeta{
				Name: "artifacts",
			},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{
					corev1.ReadWriteOnce,
				},
				Resources: corev1.VolumeResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceStorage: resource.MustParse(builderPlaygroundDeployment.Spec.Storage.Size),
					},
				},
				StorageClassName: &builderPlaygroundDeployment.Spec.Storage.StorageClass,
			},
		}}

	} else {
		// Default to local-path
		volumes = []corev1.Volume{{
			Name: "artifacts",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: builderPlaygroundDeployment.Spec.Storage.Path,
					Type: hostPathTypePtr("Directory"),
				},
			},
		}}
	}

	// Init container to check directory (applies in both cases)
	initContainers = []corev1.Container{{
    	Name:  "prepare-artifacts",
    	Image: "golang:1.24",
    	Command: []string{
    		"sh", "-c",
    		`if [ "$(find /artifacts/output -mindepth 1 ! -name 'lost+found' -print -quit 2>/dev/null)" = "" ]; then
               echo "Preparing /artifacts directory..."
               cd src
               git clone https://github.com/flashbots/builder-playground.git
               cd builder-playground
               go clean -cache -modcache
               go mod tidy
               go build -o builder-playground .
               ./builder-playground cook l1 \
                 --latest-fork \
                 --use-reth-for-validation \
                 --output /artifacts/output \
                 --genesis-delay 15 \
                 --log-level debug \
                 --dry-run
               echo "Configuration files generated."
               mv /artifacts/output/* /artifacts/output/.* /artifacts/ 2>/dev/null
               echo "Data has been generated" > /artifacts/output/data-generated.txt
               echo "Init completed successfully"
            else
               echo "/artifacts/output already exists and is not empty. Skipping preparation."
            fi`,
    	},
		
    	VolumeMounts: volumeMounts,
    }}

    // initContainers = []corev1.Container{{
	// 	Name:    "check-artifacts-dir",
	// 	Image:   "busybox:1.37",
	// 	Command: []string{"sh", "-c", "test -d /artifacts"},
	// 	VolumeMounts: volumeMounts,
	// }}

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
					Volumes:        volumes,
					InitContainers: initContainers,
					Containers:     containers,
				},
			},
			VolumeClaimTemplates: volumeClaimTemplates,
		},
	}
}



// func generateStatefulSetForOperator(builderPlaygroundDeployment *builderplaygroundv1alpha1.BuilderPlaygroundDeployment, scheme *runtime.Scheme) *appsv1.StatefulSet {
// 	labels := map[string]string{"app": builderPlaygroundDeployment.Name}

// 	containers := buildContainers(builderPlaygroundDeployment)


// 	return &appsv1.StatefulSet{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      builderPlaygroundDeployment.Name,
// 			Namespace: builderPlaygroundDeployment.Namespace,
// 		},
// 		Spec: appsv1.StatefulSetSpec{
// 			ServiceName: builderPlaygroundDeployment.Name,
// 			Replicas:    int32Ptr(1),
// 			Selector: &metav1.LabelSelector{
// 				MatchLabels: labels,
// 			},
// 			Template: corev1.PodTemplateSpec{
// 				ObjectMeta: metav1.ObjectMeta{
// 					Labels: labels,
// 				},
// 				Spec: corev1.PodSpec{
// 					Volumes: []corev1.Volume{{
// 						Name: "artifacts",
// 						VolumeSource: corev1.VolumeSource{
// 							HostPath: &corev1.HostPathVolumeSource{
// 								Path: builderPlaygroundDeployment.Spec.Storage.Path,
// 								Type: hostPathTypePtr("Directory"),
// 							},
// 						},
// 					}},
// 					InitContainers: []corev1.Container{{
// 						Name:    "check-artifacts-dir",
// 						Image:   "busybox:1.37",
// 						Command: []string{"sh", "-c", "test -d /artifacts"},
// 						VolumeMounts: []corev1.VolumeMount{{
// 							Name:      "artifacts",
// 							MountPath: "/artifacts",
// 						}},
// 					}},					
// 					Containers: containers,
// 				},
// 			},
// 		},
// 	}
// }

// This is used for StatefulSet YAML file
func generateStatefulSetManifest(builderPlaygroundDeployment *builderplaygroundv1alpha1.BuilderPlaygroundDeployment) *appsv1.StatefulSet {
	labels := map[string]string{"app": builderPlaygroundDeployment.Name}
	containers := buildContainers(builderPlaygroundDeployment)

	return &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "StatefulSet",
		},
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
						Name: "artifacts",
						VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{
								Path: builderPlaygroundDeployment.Spec.Storage.Path,
								Type: hostPathTypePtr("Directory"),
							},
						},
					}},
					InitContainers: []corev1.Container{{
						Name:    "check-artifacts-dir",
						Image:   "busybox:1.37",
						Command: []string{"sh", "-c", "test -d /artifacts"},
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "artifacts",
							MountPath: "/artifacts",
						}},
					}},
					Containers: containers,
				},
			},
		},
	}
}


func buildContainers(deploy *builderplaygroundv1alpha1.BuilderPlaygroundDeployment) []corev1.Container {
	var containers []corev1.Container

	for _, svc := range deploy.Spec.Services {
		container := corev1.Container{
			Name:    svc.Name,
			Image:   svc.Image + ":" + svc.Tag,
			Command: []string{svc.Entrypoint},
			Args:    svc.Args,
		}

		// Handle ports
		for _, p := range svc.Ports {
			container.Ports = append(container.Ports, corev1.ContainerPort{
				Name:          p.Name,
				ContainerPort: int32(p.Port),
				Protocol:      corev1.Protocol(strings.ToUpper(p.Protocol)), // tcp/udp
			})
		}

		// Handle env vars
		for k, v := range svc.Env {
			container.Env = append(container.Env, corev1.EnvVar{
				Name:  k,
				Value: v,
			})
		}

		// Handle volume mounts
		for _, vol := range svc.Volumes {
			container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
				Name:      vol.Name,
				MountPath: vol.MountPath,
			})
		}

		containers = append(containers, container)
	}

	return containers
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

// Generate CR file and StatefulSet YAML files 
func GenerateCRAndStatefulSet(manifestPath, outputDir string, builderConfigDir string) error {	
	ProcessFileForArgs(manifestPath, outputDir)
	ProcessFileForPorts(outputDir, outputDir)
	ProcessFileForURL(outputDir, outputDir)

	fmt.Printf("✅ Reading manifest from: %s\n", manifestPath)
	fmt.Printf("✅ Generating Kubernetes manifests into: %s\n", outputDir)
	fmt.Printf("✅ Builder Playground configuration files location: %s\n", builderConfigDir)

	// Reading a json file and generating K8S CR files
	jsonData, err := os.ReadFile(outputDir + "/" + "processed.json")
	if err != nil {
		fmt.Printf("✅ Reading manifest from: %v\n", err)
	}
	
	// Load Json file content to BuilderPlaygroundDeploymentSpec
	var spec builderplaygroundv1alpha1.BuilderPlaygroundDeploymentSpec
	if err := json.Unmarshal(jsonData, &spec); err != nil {
		fmt.Printf("✅ Reading manifest from: %v\n", err)
	}

	var storage builderplaygroundv1alpha1.BuilderPlaygroundStorage
	storage.Type = "local-path"
	storage.Path = builderConfigDir	
    spec.Storage = storage
	
	// Generate other part of the CR
	builderPlaygroundDeployment := &builderplaygroundv1alpha1.BuilderPlaygroundDeployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "builderplayground.flashbots.io/v1alpha1",
			Kind:       "BuilderPlaygroundDeployment",
	    },
		ObjectMeta: metav1.ObjectMeta{
			Name:      "builder-playground-sts",
			Namespace: "dev-test", // change as needed
		},
	}

	builderPlaygroundDeployment.Spec = spec

	// Marshal to YAML
	yamlData, err := yaml.Marshal(builderPlaygroundDeployment)
	if err != nil {
		return fmt.Errorf("error marshaling YAML: %w", err)
	}

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("error creating output directory: %w", err)
	}

	// Write CR YAML to file
	outputPath := filepath.Join(outputDir, "CR-BuilderPlaygroundDeployment.yaml")
	if err := os.WriteFile(outputPath, yamlData, 0644); err != nil {
		return fmt.Errorf("error writing YAML file: %w", err)
	}

	fmt.Printf("✅ Manifest written to: %s\n", outputPath)

	sts := generateStatefulSetManifest(builderPlaygroundDeployment)
	yamlData, err = yaml.Marshal(sts)
	outputPath = filepath.Join(outputDir, "BuilderPlaygroundStatefulSet.yaml")
	if err := os.WriteFile(outputPath, yamlData, 0644); err != nil {
		return fmt.Errorf("error writing YAML file: %w", err)
	}

	fmt.Printf("✅ Manifest written to: %s\n", outputPath)

	return nil
}
