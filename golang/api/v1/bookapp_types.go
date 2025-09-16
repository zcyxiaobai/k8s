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

package v1

import (
	corvev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// BookappSpec defines the desired state of Bookapp
type BookappSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// The following markers will use OpenAPI v3 schema to validate the value
	// More info: https://book.kubebuilder.io/reference/markers/crd-validation.html

	// foo is an example field of Bookapp. Edit bookapp_types.go to remove/update
	// +optional
	//Foo *string `json:"foo,omitempty"`

	//填写自己的crd资源的字段
	Size            *int32              `json:"size"`
	Image           string              `json:"image"`
	ImagePullPolicy corvev1.PullPolicy  `json:"imagePullPolicy"`
	Port            corvev1.ServicePort `json:"port"`
	ServerName      string              `json:"serverName"`
}

// BookappStatus defines the observed state of Bookapp.
// 状态
type BookappStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// conditions represent the current state of the Bookapp resource.
	// Each condition has a unique type and reflects the status of a specific aspect of the resource.
	//
	// Standard condition types include:
	// - "Available": the resource is fully functional
	// - "Progressing": the resource is being created or updated
	// - "Degraded": the resource failed to reach or maintain its desired state
	//
	// The status of each condition is one of True, False, or Unknown.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	//定义两个字段，runing和notruning
	Runing    *int32 `json:"runing"`
	NotRuning *int32 `json:"notRuning"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// Bookapp is the Schema for the bookapps API
// +kubebuilder:printcolumn:name="Size",type="integer",JSONPath=".spec.size"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="Running",type=integer,JSONPath=`.status.runing`,description="Number of running replicas"
// +kubebuilder:printcolumn:name="NotRunning",type=integer,JSONPath=`.status.notRuning`,description="Number of not running replicas"
// +kubebuilder:printcolumn:name="Image",type=string,JSONPath=`.spec.image`,description="Container image"

type Bookapp struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of Bookapp
	// +required
	Spec BookappSpec `json:"spec"`

	// status defines the observed state of Bookapp
	// +optional
	Status BookappStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// BookappList contains a list of Bookapp
type BookappList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Bookapp `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Bookapp{}, &BookappList{})
}
