/*
Copyright 2021 Red Hat, Inc.

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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ComponentSrcType describes the type of
// the src for the Component.
// Only one of the following location type may be specified.
// +kubebuilder:validation:Enum=Git;Image
type ComponentSrcType string

const (
	GitComponentSrcType   ComponentSrcType = "Git"
	ImageComponentSrcType ComponentSrcType = "Image"
)

type GitSource struct {
	// If importing from git, the repository to create the component from
	URL string `json:"url"`

	// Secret containing a Person Access Token to clone a sample, if using private repository
	Secret string `json:"secret,omitempty"`

	// If specified, the devfile at the URL will be used for the component.
	DevfileURL string `json:"devfileUrl,omitempty"`
}

type ImageSource struct {
	// If importing from git, container image to create the component from
	ContainerImage string `json:"containerImage"`
}

// ComponentSource describes the Component source
type ComponentSource struct {
	ComponentSourceUnion `json:",inline"`
}

// +union
type ComponentSourceUnion struct {
	// Git Source for a Component
	GitSource *GitSource `json:"git,omitempty"`

	// Image Source for a Component
	ImageSource *ImageSource `json:"image,omitempty"`
}

// Build describes the various build artifacts associated with a given component
type Build struct {
	// The container image that is created during the component build.
	ContainerImage string `json:"containerImage"`
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ComponentSpec defines the desired state of Component
type ComponentSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// ComponentName is name of the component to be added to the HASApplication
	ComponentName string `json:"componentName"`

	// Application to add the component to
	Application string `json:"application"`

	// Source describes the Component source
	Source ComponentSource `json:"source"`

	// A relative path inside the git repo containing the component
	Context string `json:"context,omitempty"`

	// Compute Resources required by this component
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`

	// The number of replicas to deploy the component with
	Replicas int `json:"replicas,omitempty"`

	// The port to expose the component over
	TargetPort int `json:"targetPort,omitempty"`

	// The route to expose the component with
	Route string `json:"route,omitempty"`

	// An array of environment variables to add to the component
	Env []corev1.EnvVar `json:"env,omitempty"`

	// The build artifacts associated with the component
	Build Build `json:"build,omitempty"`
}

// ComponentStatus defines the observed state of Component
type ComponentStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Condition about the Component CR
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// ContainerImage stores the associated built container image for the component
	ContainerImage string `json:"containerImage,omitempty"`

	// The devfile model for the Component CR
	Devfile string `json:"devfile,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Component is the Schema for the components API
// +kubebuilder:resource:path=components,shortName=hascmp;hc;comp
type Component struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ComponentSpec   `json:"spec,omitempty"`
	Status ComponentStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ComponentList contains a list of Component
type ComponentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Component `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Component{}, &ComponentList{})
}
