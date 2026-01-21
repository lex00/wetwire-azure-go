// Package v1 contains ASO Managed Identity resource types.
package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// UserAssignedIdentity represents an ASO User Assigned Identity resource.
// +kubebuilder:object:root=true
type UserAssignedIdentity struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   UserAssignedIdentitySpec   `json:"spec,omitempty"`
	Status UserAssignedIdentityStatus `json:"status,omitempty"`
}

// UserAssignedIdentitySpec defines the desired state of a User Assigned Identity.
type UserAssignedIdentitySpec struct {
	// Owner is the resource group owner reference.
	Owner *ResourceGroupReference `json:"owner,omitempty"`

	// AzureName is the name of the resource in Azure.
	AzureName *string `json:"azureName,omitempty"`

	// Location is the Azure region.
	Location *string `json:"location,omitempty"`

	// Tags are key-value pairs.
	Tags map[string]string `json:"tags,omitempty"`
}

// UserAssignedIdentityStatus defines the observed state of a User Assigned Identity.
type UserAssignedIdentityStatus struct {
	// Conditions represent the latest available observations.
	Conditions []Condition `json:"conditions,omitempty"`

	// ID is the Azure resource ID.
	ID *string `json:"id,omitempty"`

	// ClientId is the client ID of the identity.
	ClientId *string `json:"clientId,omitempty"`

	// PrincipalId is the principal ID of the identity.
	PrincipalId *string `json:"principalId,omitempty"`

	// TenantId is the tenant ID of the identity.
	TenantId *string `json:"tenantId,omitempty"`
}

// ResourceGroupReference references a Resource Group.
type ResourceGroupReference struct {
	// Name is the name of the resource group.
	Name string `json:"name,omitempty"`

	// ARMID is the Azure Resource Manager ID.
	ARMID *string `json:"armId,omitempty"`
}

// Condition represents a condition.
type Condition struct {
	// Type is the type of condition.
	Type string `json:"type,omitempty"`

	// Status is the status of the condition.
	Status string `json:"status,omitempty"`

	// LastTransitionTime is when the condition last transitioned.
	LastTransitionTime *metav1.Time `json:"lastTransitionTime,omitempty"`

	// Message is a human-readable message.
	Message *string `json:"message,omitempty"`

	// Reason is a brief reason for the condition.
	Reason *string `json:"reason,omitempty"`
}

// FederatedIdentityCredential represents an ASO Federated Identity Credential resource.
// Used for workload identity federation with Kubernetes.
// +kubebuilder:object:root=true
type FederatedIdentityCredential struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FederatedIdentityCredentialSpec   `json:"spec,omitempty"`
	Status FederatedIdentityCredentialStatus `json:"status,omitempty"`
}

// FederatedIdentityCredentialSpec defines the desired state of a Federated Identity Credential.
type FederatedIdentityCredentialSpec struct {
	// Owner is the user assigned identity owner reference.
	Owner *UserAssignedIdentityReference `json:"owner,omitempty"`

	// AzureName is the name of the resource in Azure.
	AzureName *string `json:"azureName,omitempty"`

	// Audiences are the token audiences.
	Audiences []string `json:"audiences,omitempty"`

	// Issuer is the OIDC issuer URL.
	Issuer *string `json:"issuer,omitempty"`

	// Subject is the subject claim (e.g., system:serviceaccount:namespace:sa-name).
	Subject *string `json:"subject,omitempty"`
}

// FederatedIdentityCredentialStatus defines the observed state of a Federated Identity Credential.
type FederatedIdentityCredentialStatus struct {
	// Conditions represent the latest available observations.
	Conditions []Condition `json:"conditions,omitempty"`

	// ID is the Azure resource ID.
	ID *string `json:"id,omitempty"`
}

// UserAssignedIdentityReference references a UserAssignedIdentity.
type UserAssignedIdentityReference struct {
	// Name is the name of the identity resource.
	Name string `json:"name,omitempty"`

	// ARMID is the Azure Resource Manager ID.
	ARMID *string `json:"armId,omitempty"`
}
