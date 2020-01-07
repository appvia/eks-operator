package v1alpha1

import (
	core "github.com/appvia/hub-apis/pkg/apis/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EKSClusterSpec defines the desired state of EKSCluster
// +k8s:openapi-gen=true
type EKSClusterSpec struct {
	// Name the name of the EKS cluster
	// +kubebuilder:validation:MinLength=3
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// RoleArn is the role arn which provides permissions to EKS
	// +kubebuilder:validation:MinLength=10
	// +kubebuilder:validation:Required
	RoleArn string `json:"rolearn"`
	// Version is the Kubernetes version to use
	// +kubebuilder:validation:MinLength=3
	// +kubebuilder:validation:Required
	Version string `json:"version"`
	// Use is a reference to an AWSCredentials object to use for authentication
	// +kubebuilder:validation:Required
	// +k8s:openapi-gen=false
	Use core.Ownership `json:"use"`
}

// EKSClusterStatus defines the observed state of EKSCluster
// +k8s:openapi-gen=true
type EKSClusterStatus struct {
	// Status provides a overall status
	Status core.Status `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EKSCluster is the Schema for the eksclusters API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=eksclusters,scope=Namespaced
type EKSCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EKSClusterSpec   `json:"spec,omitempty"`
	Status EKSClusterStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EKSClusterList contains a list of EKSCluster
type EKSClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EKSCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EKSCluster{}, &EKSClusterList{})
}
