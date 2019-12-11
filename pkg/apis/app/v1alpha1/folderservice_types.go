package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// FolderServiceSpec defines the desired state of FolderService
// +k8s:openapi-gen=true
type FolderServiceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	// Sleep string `json:"sleep"`

	// CertificateSecret is the reference to the secret where certificates are stored.
	UserName        string          `json:"userName"`
	UserSecret      UserSecret      `json:"userSecret"`
	PlatformSecrets PlatformSecrets `json:"platformSecrets"`
}

// UserSecret defines the type of Usersecret struct
type UserSecret struct {
	Name string `json:"name"`
}

// FolderServiceStatus defines the observed state of FolderService
// +k8s:openapi-gen=true
type FolderServiceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	// Sleep string `json:"sleep"`
	SetupComplete bool `json:"setupComplete"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FolderService is the Schema for the folderservices API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=folderservices,scope=Namespaced
type FolderService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FolderServiceSpec   `json:"spec,omitempty"`
	Status FolderServiceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FolderServiceList contains a list of FolderService
type FolderServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FolderService `json:"items"`
}

// PlatformSecrets defines the secrets to be used by various clouds.
type PlatformSecrets struct {
	AWS *AWSPlatformSecrets `json:"aws"`
}

// AWSPlatformSecrets contains secrets for clusters on the AWS platform.
type AWSPlatformSecrets struct {
	// Credentials refers to a secret that contains the AWS account access
	// credentials.
	Credentials corev1.LocalObjectReference `json:"credentials"`
}

func init() {
	SchemeBuilder.Register(&FolderService{}, &FolderServiceList{})
}
