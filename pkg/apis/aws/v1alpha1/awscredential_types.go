package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AWSCredentialSpec defines the desired state of AWSCredential
// +k8s:openapi-gen=true
type AWSCredentialSpec struct {
	// AWS Secret Access Key
	// +kubebuilder:validation:Minimum=3
	// +kubebuilder:validation:Required
	SecretAccessKey string `json:"secret"`
	// AWS Access Key ID
	// +kubebuilder:validation:Minimum=3
	// +kubebuilder:validation:Required
	AccessKeyId string `json:"id"`
	// Account is the AWS account these credentials reside within
	// +kubebuilder:validation:Minimum=3
	// +kubebuilder:validation:Required
	AccountId string `json:"accountId"`
}

// AWSCredentialStatus defines the observed state of AWSCredential
// +k8s:openapi-gen=true
type AWSCredentialStatus struct {
	// Verified checks that the credentials are ok and valid
	Verified bool `json:"verified"`
	// Status provides a overall status
	Status string `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AWSCredential is the Schema for the awscredentials API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=awscredentials,scope=Namespaced
type AWSCredential struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AWSCredentialSpec   `json:"spec,omitempty"`
	Status AWSCredentialStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AWSCredentialList contains a list of AWSCredential
type AWSCredentialList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSCredential `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AWSCredential{}, &AWSCredentialList{})
}
