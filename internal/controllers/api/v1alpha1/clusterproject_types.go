/*
Copyright 2024 Timoni Authors.

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

// ClusterProjectSpec defines the desired state of ClusterProject
type ClusterProjectSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +optional
	// +kubebuilder:default={"scope":"cluster"}
	Access ClusterProjectAccess `json:"access,omitempty"`
	// +required
	Source ClusterProjectSource `json:"source,omitempty"`
	// +optional
	Watches []ClusterProjectWatch `json:"watches,omitempty"`
}

type ClusterProjectAccess struct {
	// +optional
	// +kubebuilder:validation:Enum=cluster;namespaced
	// +kubebuilder:default=cluster
	Scope string `json:"scope,omitempty"`
	// +optional
	// +kubebuilder:validation:MinItems=1
	Namespaces []string `json:"namespaces,omitempty"`
}

type ClusterProjectSource struct {
	// +required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern=`^oci://.*$`
	Repository string `json:"repository"`
	// +optional
	Tag string `json:"tag,omitempty"`
	// +optional
	Path string `json:"path,omitempty"`
}

type ClusterProjectWatch struct {
	// +required
	APIVersion string `json:"apiVersion,omitempty"`
	// +required
	Kind string `json:"kind,omitempty"`
}

// ClusterProjectStatus defines the observed state of ClusterProject
type ClusterProjectStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ClusterProject is the Schema for the clusterprojects API
type ClusterProject struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterProjectSpec   `json:"spec,omitempty"`
	Status ClusterProjectStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ClusterProjectList contains a list of ClusterProject
type ClusterProjectList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterProject `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterProject{}, &ClusterProjectList{})
}
