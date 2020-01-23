package v1alpha1

import (
	core "github.com/appvia/hub-apis/pkg/apis/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EKSBearerTokenSpec defines the desired state of EKSBearerToken
// +k8s:openapi-gen=true
type EKSBearerTokenSpec struct {
	Token string `json:"token"`
}

// EKSBearerTokenStatus defines the observed state of EKSBearerToken
// +k8s:openapi-gen=true
type EKSBearerTokenStatus struct {
	// Status provides a overall status
	Status core.Status `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EKSBearerToken is the Schema for the eksbearertokens API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=eksbearertokens,scope=Namespaced
type EKSBearerToken struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EKSBearerTokenSpec   `json:"spec,omitempty"`
	Status EKSBearerTokenStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EKSBearerTokenList contains a list of EKSBearerToken
type EKSBearerTokenList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EKSBearerToken `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EKSBearerToken{}, &EKSBearerTokenList{})
}
