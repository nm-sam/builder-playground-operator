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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// BuilderPlaygroundDeploymentStatus defines the observed state of BuilderPlaygroundDeployment.
type BuilderPlaygroundDeploymentStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// BuilderPlaygroundDeployment is the Schema for the builderplaygrounddeployments API.
type BuilderPlaygroundDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BuilderPlaygroundDeploymentSpec   `json:"spec,omitempty"`
	Status BuilderPlaygroundDeploymentStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// BuilderPlaygroundDeploymentList contains a list of BuilderPlaygroundDeployment.
type BuilderPlaygroundDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BuilderPlaygroundDeployment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BuilderPlaygroundDeployment{}, &BuilderPlaygroundDeploymentList{})
}

// BuilderPlaygroundDeploymentSpec defines the desired state of the deployment
type BuilderPlaygroundDeploymentSpec struct {
    // Recipe is the builder-playground recipe used (l1, opstack, etc)
    Recipe string `json:"recipe"`
    // Storage defines how persistent data should be stored
    Storage BuilderPlaygroundStorage `json:"storage"`
    // Network defines networking configuration (optional)
    Network *BuilderPlaygroundNetwork `json:"network,omitempty"`
    // Services is the list of services in this deployment
    Services []BuilderPlaygroundService `json:"services"`
}

// BuilderPlaygroundStorage defines storage configuration
type BuilderPlaygroundStorage struct {
    // Type is the storage type, either "local-path" or "pvc"
    Type string `json:"type"`
    // Path is the host path for local-path storage (used when type is "local-path")
    Path string `json:"path,omitempty"`
    // StorageClass is the K8s storage class (used when type is "pvc")
    StorageClass string `json:"storageClass,omitempty"`
    // Size is the storage size (used when type is "pvc")
    Size string `json:"size,omitempty"`
}

// BuilderPlaygroundNetwork defines network configuration
type BuilderPlaygroundNetwork struct {
    // Name is the name of the network
    Name string `json:"name"`
}

// BuilderPlaygroundService represents a single service in the deployment
type BuilderPlaygroundService struct {
    // Name is the service name
    Name string `json:"name"`
    // Image is the container image
    Image string `json:"image"`
    // Tag is the container image tag
    Tag string `json:"tag"`
    // Entrypoint overrides the container entrypoint
    Entrypoint string `json:"entrypoint,omitempty"`
    // Args are the container command arguments
    Args []string `json:"args,omitempty"`
    // Env defines environment variables
    Env map[string]string `json:"env,omitempty"`
    // Ports are the container ports to expose
    Ports []BuilderPlaygroundPort `json:"ports,omitempty"`
    // Dependencies defines services this service depends on
    Dependencies []BuilderPlaygroundDependency `json:"dependencies,omitempty"`
    // ReadyCheck defines how to determine service readiness
    ReadyCheck *BuilderPlaygroundReadyCheck `json:"readyCheck,omitempty"`
    // Labels are the service labels
    Labels map[string]string `json:"labels,omitempty"`
    // UseHostExecution indicates whether to run on host instead of in container
    UseHostExecution bool `json:"useHostExecution,omitempty"`
    // Volumes are the volume mounts for the service
    Volumes []BuilderPlaygroundVolume `json:"volumes,omitempty"`
}

// BuilderPlaygroundPort represents a port configuration
type BuilderPlaygroundPort struct {
    // Name is a unique identifier for this port
    Name string `json:"name"`
    // Port is the container port number
    Port int `json:"port"`
    // Protocol is either "tcp" or "udp"
    Protocol string `json:"protocol,omitempty"`
    // HostPort is the port to expose on the host (if applicable)
    HostPort int `json:"hostPort,omitempty"`
}

// BuilderPlaygroundDependency represents a service dependency
type BuilderPlaygroundDependency struct {
    // Name is the name of the dependent service
    Name string `json:"name"`
    // Condition is either "running" or "healthy"
    Condition string `json:"condition"`
}

// BuilderPlaygroundReadyCheck defines readiness checking
type BuilderPlaygroundReadyCheck struct {
    // QueryURL is the URL to query for readiness
    QueryURL string `json:"queryURL,omitempty"`
    // Test is the command to run for readiness check
    Test []string `json:"test,omitempty"`
    // Interval is the time between checks
    Interval string `json:"interval,omitempty"`
    // Timeout is the maximum time for a check
    Timeout string `json:"timeout,omitempty"`
    // Retries is the number of retry attempts
    Retries int `json:"retries,omitempty"`
    // StartPeriod is the initial delay before checks begin
    StartPeriod string `json:"startPeriod,omitempty"`
}

// BuilderPlaygroundVolume represents a volume mount
type BuilderPlaygroundVolume struct {
    // Name is the volume name
    Name string `json:"name"`
    // MountPath is the path in the container
    MountPath string `json:"mountPath"`
    // SubPath is the path within the volume (optional)
    SubPath string `json:"subPath,omitempty"`
}

