package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AWSTokenSpec defines the desired state of AWSToken
// +k8s:openapi-gen=true
type AWSTokenSpec struct {
	// AWS Secret Access Key
	// +kubebuilder:validation:Minimum=3
	// +kubebuilder:validation:Required
	SecretAccessKey string `json:"secret"`
	// AWS Access Key ID
	// +kubebuilder:validation:Minimum=12
	// +kubebuilder:validation:Maximum=12
	// +kubebuilder:validation:Required
	AccessKeyId string `json:"id"`
	// AWS Session Token
	// +kubebuilder:validation:Minimum=3
	// +kubebuilder:validation:Required
	SessionToken string `json:"token"`
	// Account is the AWS account these credentials reside within
	// +kubebuilder:validation:MinLength=12
	// +kubebuilder:validation:MaxLength=12
	// +kubebuilder:validation:Required
	AccountId string `json:"accountId"`
	// Expiration is the expiry date time of this token
	// +kubebuilder:validation:Required
	Expiration string `json:"expiration"`
}

// AWSTokenStatus defines the observed state of AWSToken
// +k8s:openapi-gen=true
type AWSTokenStatus struct {
	// Verified checks that the credentials are ok and valid
	Verified bool `json:"verified"`
	// Status provides a overall status
	Status string `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AWSToken is the Schema for the awstokens API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=awstokens,scope=Namespaced
type AWSToken struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AWSTokenSpec   `json:"spec,omitempty"`
	Status AWSTokenStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AWSTokenList contains a list of AWSToken
type AWSTokenList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSToken `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AWSToken{}, &AWSTokenList{})
}
