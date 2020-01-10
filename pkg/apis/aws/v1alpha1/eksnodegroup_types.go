package v1alpha1

import (
	core "github.com/appvia/hub-apis/pkg/apis/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EKSNodeGroupSpec defines the desired state of EKSNodeGroup
// +k8s:openapi-gen=true
type EKSNodeGroupSpec struct {
	AmiType string `json:"amitype,omitempty"`
	// +kubebuilder:validation:Required
	ClusterName        string `json:"clustername"`
	DiskSize           string `json:"disksize,omitempty"`
	ForceUpdateEnabled string `json:"forceupdateenabled,omitempty"`
	InstanceTypes      string `json:"instancetypes,omitempty"`
	Labels             string `json:"labels,omitempty"`
	NodegroupName      string `json:"nodegroupname,omitempty"`
	// +kubebuilder:validation:Required
	NodeRole       string `json:"noderole"`
	ReleaseVersion string `json:"releaseversion,omitempty"`
	RemoteAccess   string `json:"remoteaccess,omitempty"`
	ScalingConfig  string `json:"scalingconfig,omitempty"`
	// +kubebuilder:validation:Required
	Subnets string `json:"subnets"`
	Tags    string `json:"tags,omitempty"`
	Version string `json:"version,omitempty"`
}

// EKSNodeGroupStatus defines the observed state of EKSNodeGroup
// +k8s:openapi-gen=true
type EKSNodeGroupStatus struct {
	// Status provides a overall status
	Status core.Status `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EKSNodeGroup is the Schema for the eksnodegroups API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=eksnodegroups,scope=Namespaced
type EKSNodeGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EKSNodeGroupSpec   `json:"spec,omitempty"`
	Status EKSNodeGroupStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EKSNodeGroupList contains a list of EKSNodeGroup
type EKSNodeGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EKSNodeGroup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EKSNodeGroup{}, &EKSNodeGroupList{})
}
