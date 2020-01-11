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
	DiskSize           int64 `json:"disksize,omitempty"`
	InstanceTypes      []string `json:"instancetypes,omitempty"`
	Labels             map[string]string `json:"labels,omitempty"`
	NodegroupName      string `json:"nodegroupname"`
	// +kubebuilder:validation:Required
	NodeRole       string `json:"noderole"`
	ReleaseVersion string `json:"releaseversion,omitempty"`
	RemoteAccess   string `json:"remoteaccess,omitempty"`
	DesiredSize    int64    `json:"desiredsize,omitempty"`
	MaxSize        int64    `json:"maxsize,omitempty"`
	MinSize        int64    `json:"minsize,omitempty"`
	// +kubebuilder:validation:Required
	Subnets []string `json:"subnets"`
	// The metadata to apply to the node group
	Tags    map[string]string   `json:"tags,omitempty"`
	// The Kubernetes version to use for your managed nodes
	Version string   `json:"version,omitempty"`
	// AWS region to launch node group within, must match the region of the cluster
	Region  string   `json:"region"`
	// The Amazon EC2 SSH key that provides access for SSH communication with 
	// the worker nodes in the managed node group
	// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-key-pairs.html
	SourceSecurityGroups []string `json:"sourcesecuritygroups,omitempty"`
	// The security groups that are allowed SSH access (port 22) to the worker nodes
	Ec2SshKey	string `json:"ec2sshkey,omitempty"`
	// Use is a reference to an AWSCredentials object to use for authentication
	// +kubebuilder:validation:Required
	// +k8s:openapi-gen=false
	Use core.Ownership `json:"use"`
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
