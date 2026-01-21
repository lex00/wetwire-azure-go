// Package v1 contains ASO Container Service resource types.
package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ManagedCluster represents an ASO AKS ManagedCluster resource.
// +kubebuilder:object:root=true
type ManagedCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ManagedClusterSpec   `json:"spec,omitempty"`
	Status ManagedClusterStatus `json:"status,omitempty"`
}

// ManagedClusterSpec defines the desired state of an AKS Managed Cluster.
type ManagedClusterSpec struct {
	// Owner is the resource group owner reference.
	Owner *ResourceGroupReference `json:"owner,omitempty"`

	// AzureName is the name of the resource in Azure.
	AzureName *string `json:"azureName,omitempty"`

	// Location is the Azure region.
	Location *string `json:"location,omitempty"`

	// DNSPrefix is the DNS prefix for the cluster.
	DNSPrefix *string `json:"dnsPrefix,omitempty"`

	// KubernetesVersion is the Kubernetes version.
	KubernetesVersion *string `json:"kubernetesVersion,omitempty"`

	// AgentPoolProfiles are the agent pool configurations.
	AgentPoolProfiles []ManagedClusterAgentPoolProfile `json:"agentPoolProfiles,omitempty"`

	// Identity is the identity configuration.
	Identity *ManagedClusterIdentity `json:"identity,omitempty"`

	// NetworkProfile is the network configuration.
	NetworkProfile *ContainerServiceNetworkProfile `json:"networkProfile,omitempty"`

	// AADProfile is the Azure AD configuration.
	AADProfile *ManagedClusterAADProfile `json:"aadProfile,omitempty"`

	// APIServerAccessProfile is the API server access configuration.
	APIServerAccessProfile *ManagedClusterAPIServerAccessProfile `json:"apiServerAccessProfile,omitempty"`

	// AutoScalerProfile is the cluster autoscaler configuration.
	AutoScalerProfile *ManagedClusterAutoScalerProfile `json:"autoScalerProfile,omitempty"`

	// EnableRBAC indicates whether RBAC is enabled.
	EnableRBAC *bool `json:"enableRBAC,omitempty"`

	// NodeResourceGroup is the name of the resource group for cluster nodes.
	NodeResourceGroup *string `json:"nodeResourceGroup,omitempty"`

	// SKU is the SKU configuration.
	SKU *ManagedClusterSKU `json:"sku,omitempty"`

	// SecurityProfile is the security configuration.
	SecurityProfile *ManagedClusterSecurityProfile `json:"securityProfile,omitempty"`

	// OIDCIssuerProfile is the OIDC issuer configuration.
	OIDCIssuerProfile *ManagedClusterOIDCIssuerProfile `json:"oidcIssuerProfile,omitempty"`

	// Tags are key-value pairs.
	Tags map[string]string `json:"tags,omitempty"`
}

// ManagedClusterStatus defines the observed state of an AKS Managed Cluster.
type ManagedClusterStatus struct {
	// Conditions represent the latest available observations.
	Conditions []Condition `json:"conditions,omitempty"`

	// ID is the Azure resource ID.
	ID *string `json:"id,omitempty"`

	// Fqdn is the FQDN of the cluster.
	Fqdn *string `json:"fqdn,omitempty"`

	// ProvisioningState is the current provisioning state.
	ProvisioningState *string `json:"provisioningState,omitempty"`

	// PowerState is the current power state.
	PowerState *PowerState `json:"powerState,omitempty"`
}

// ResourceGroupReference references a Resource Group.
type ResourceGroupReference struct {
	// Name is the name of the resource group.
	Name string `json:"name,omitempty"`

	// ARMID is the Azure Resource Manager ID.
	ARMID *string `json:"armId,omitempty"`
}

// ManagedClusterAgentPoolProfile represents an agent pool configuration.
type ManagedClusterAgentPoolProfile struct {
	// Name is the unique name of the agent pool.
	Name *string `json:"name,omitempty"`

	// Count is the number of agents (VMs).
	Count *int `json:"count,omitempty"`

	// VMSize is the size of agent VMs.
	VMSize *string `json:"vmSize,omitempty"`

	// OSDiskSizeGB is the OS disk size in GB.
	OSDiskSizeGB *int `json:"osDiskSizeGB,omitempty"`

	// OSDiskType is the OS disk type.
	OSDiskType *string `json:"osDiskType,omitempty"`

	// VnetSubnetReference is a reference to a Subnet.
	VnetSubnetReference *SubnetReference `json:"vnetSubnetReference,omitempty"`

	// MaxPods is the maximum number of pods per node.
	MaxPods *int `json:"maxPods,omitempty"`

	// OSType is the operating system type.
	OSType *string `json:"osType,omitempty"`

	// OSSKU is the OS SKU.
	OSSKU *string `json:"osSKU,omitempty"`

	// MinCount is the minimum number of nodes for auto-scaling.
	MinCount *int `json:"minCount,omitempty"`

	// MaxCount is the maximum number of nodes for auto-scaling.
	MaxCount *int `json:"maxCount,omitempty"`

	// EnableAutoScaling enables auto-scaling.
	EnableAutoScaling *bool `json:"enableAutoScaling,omitempty"`

	// Type is the agent pool type.
	Type *string `json:"type,omitempty"`

	// Mode is the agent pool mode (System, User).
	Mode *string `json:"mode,omitempty"`

	// AvailabilityZones are the availability zones.
	AvailabilityZones []string `json:"availabilityZones,omitempty"`

	// NodeLabels are the labels for nodes.
	NodeLabels map[string]string `json:"nodeLabels,omitempty"`

	// NodeTaints are the taints for nodes.
	NodeTaints []string `json:"nodeTaints,omitempty"`

	// ScaleSetPriority is the priority (Regular, Spot).
	ScaleSetPriority *string `json:"scaleSetPriority,omitempty"`

	// ScaleSetEvictionPolicy is the eviction policy for Spot VMs.
	ScaleSetEvictionPolicy *string `json:"scaleSetEvictionPolicy,omitempty"`

	// SpotMaxPrice is the max price for Spot VMs.
	SpotMaxPrice *float64 `json:"spotMaxPrice,omitempty"`

	// Tags are key-value pairs.
	Tags map[string]string `json:"tags,omitempty"`
}

// SubnetReference references a Subnet.
type SubnetReference struct {
	// Name is the name of the subnet resource.
	Name *string `json:"name,omitempty"`

	// ARMID is the Azure Resource Manager ID.
	ARMID *string `json:"armId,omitempty"`
}

// ManagedClusterIdentity represents identity configuration.
type ManagedClusterIdentity struct {
	// Type is the identity type (None, SystemAssigned, UserAssigned).
	Type *string `json:"type,omitempty"`

	// UserAssignedIdentities are the user-assigned identities.
	UserAssignedIdentities []UserAssignedIdentityDetails `json:"userAssignedIdentities,omitempty"`
}

// UserAssignedIdentityDetails contains user-assigned identity details.
type UserAssignedIdentityDetails struct {
	// Reference is a reference to a UserAssignedIdentity.
	Reference *UserAssignedIdentityReference `json:"reference,omitempty"`
}

// UserAssignedIdentityReference references a UserAssignedIdentity.
type UserAssignedIdentityReference struct {
	// Name is the name of the identity resource.
	Name *string `json:"name,omitempty"`

	// ARMID is the Azure Resource Manager ID.
	ARMID *string `json:"armId,omitempty"`
}

// ContainerServiceNetworkProfile represents network configuration.
type ContainerServiceNetworkProfile struct {
	// NetworkPlugin is the network plugin (azure, kubenet, none).
	NetworkPlugin *string `json:"networkPlugin,omitempty"`

	// NetworkPolicy is the network policy (azure, calico).
	NetworkPolicy *string `json:"networkPolicy,omitempty"`

	// PodCidr is the pod CIDR.
	PodCidr *string `json:"podCidr,omitempty"`

	// ServiceCidr is the service CIDR.
	ServiceCidr *string `json:"serviceCidr,omitempty"`

	// DNSServiceIP is the DNS service IP.
	DNSServiceIP *string `json:"dnsServiceIP,omitempty"`

	// OutboundType is the outbound type.
	OutboundType *string `json:"outboundType,omitempty"`

	// LoadBalancerSku is the load balancer SKU.
	LoadBalancerSku *string `json:"loadBalancerSku,omitempty"`
}

// ManagedClusterAADProfile represents Azure AD configuration.
type ManagedClusterAADProfile struct {
	// Managed indicates whether AAD integration is managed.
	Managed *bool `json:"managed,omitempty"`

	// EnableAzureRBAC enables Azure RBAC for K8s authorization.
	EnableAzureRBAC *bool `json:"enableAzureRBAC,omitempty"`

	// AdminGroupObjectIDs are the AAD group object IDs for cluster admins.
	AdminGroupObjectIDs []string `json:"adminGroupObjectIDs,omitempty"`

	// TenantID is the AAD tenant ID.
	TenantID *string `json:"tenantID,omitempty"`
}

// ManagedClusterAPIServerAccessProfile represents API server access configuration.
type ManagedClusterAPIServerAccessProfile struct {
	// AuthorizedIPRanges are the authorized IP ranges.
	AuthorizedIPRanges []string `json:"authorizedIPRanges,omitempty"`

	// EnablePrivateCluster enables private cluster.
	EnablePrivateCluster *bool `json:"enablePrivateCluster,omitempty"`

	// PrivateDNSZone is the private DNS zone.
	PrivateDNSZone *string `json:"privateDNSZone,omitempty"`
}

// ManagedClusterAutoScalerProfile represents cluster autoscaler configuration.
type ManagedClusterAutoScalerProfile struct {
	// Expander is the expander type.
	Expander *string `json:"expander,omitempty"`

	// ScanInterval is the scan interval.
	ScanInterval *string `json:"scanInterval,omitempty"`

	// ScaleDownDelayAfterAdd is scale down delay after add.
	ScaleDownDelayAfterAdd *string `json:"scaleDownDelayAfterAdd,omitempty"`

	// ScaleDownUnneededTime is scale down unneeded time.
	ScaleDownUnneededTime *string `json:"scaleDownUnneededTime,omitempty"`

	// ScaleDownUtilizationThreshold is scale down utilization threshold.
	ScaleDownUtilizationThreshold *string `json:"scaleDownUtilizationThreshold,omitempty"`
}

// ManagedClusterSKU represents SKU configuration.
type ManagedClusterSKU struct {
	// Name is the SKU name.
	Name *string `json:"name,omitempty"`

	// Tier is the SKU tier (Free, Standard, Premium).
	Tier *string `json:"tier,omitempty"`
}

// ManagedClusterSecurityProfile represents security configuration.
type ManagedClusterSecurityProfile struct {
	// WorkloadIdentity is the workload identity configuration.
	WorkloadIdentity *ManagedClusterSecurityProfileWorkloadIdentity `json:"workloadIdentity,omitempty"`

	// ImageCleaner is the image cleaner configuration.
	ImageCleaner *ManagedClusterSecurityProfileImageCleaner `json:"imageCleaner,omitempty"`
}

// ManagedClusterSecurityProfileWorkloadIdentity represents workload identity configuration.
type ManagedClusterSecurityProfileWorkloadIdentity struct {
	// Enabled enables workload identity.
	Enabled *bool `json:"enabled,omitempty"`
}

// ManagedClusterSecurityProfileImageCleaner represents image cleaner configuration.
type ManagedClusterSecurityProfileImageCleaner struct {
	// Enabled enables image cleaner.
	Enabled *bool `json:"enabled,omitempty"`

	// IntervalHours is the interval in hours.
	IntervalHours *int `json:"intervalHours,omitempty"`
}

// ManagedClusterOIDCIssuerProfile represents OIDC issuer configuration.
type ManagedClusterOIDCIssuerProfile struct {
	// Enabled enables OIDC issuer.
	Enabled *bool `json:"enabled,omitempty"`
}

// PowerState represents the power state of the cluster.
type PowerState struct {
	// Code is the power state code.
	Code *string `json:"code,omitempty"`
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
