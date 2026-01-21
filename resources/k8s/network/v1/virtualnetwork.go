// Package v1 contains ASO Network resource types.
package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VirtualNetwork represents an ASO Virtual Network resource.
// +kubebuilder:object:root=true
type VirtualNetwork struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualNetworkSpec   `json:"spec,omitempty"`
	Status VirtualNetworkStatus `json:"status,omitempty"`
}

// VirtualNetworkSpec defines the desired state of a Virtual Network.
type VirtualNetworkSpec struct {
	// Owner is the resource group owner reference.
	Owner *ResourceGroupReference `json:"owner,omitempty"`

	// AzureName is the name of the resource in Azure.
	AzureName *string `json:"azureName,omitempty"`

	// Location is the Azure region.
	Location *string `json:"location,omitempty"`

	// AddressSpace specifies the address space.
	AddressSpace *AddressSpace `json:"addressSpace,omitempty"`

	// DhcpOptions specifies DHCP options.
	DhcpOptions *DhcpOptions `json:"dhcpOptions,omitempty"`

	// EnableDdosProtection indicates whether DDoS protection is enabled.
	EnableDdosProtection *bool `json:"enableDdosProtection,omitempty"`

	// Tags are key-value pairs.
	Tags map[string]string `json:"tags,omitempty"`
}

// VirtualNetworkStatus defines the observed state of a Virtual Network.
type VirtualNetworkStatus struct {
	// Conditions represent the latest available observations.
	Conditions []Condition `json:"conditions,omitempty"`

	// ID is the Azure resource ID.
	ID *string `json:"id,omitempty"`

	// ProvisioningState is the current provisioning state.
	ProvisioningState *string `json:"provisioningState,omitempty"`
}

// ResourceGroupReference references a Resource Group.
type ResourceGroupReference struct {
	// Name is the name of the resource group.
	Name string `json:"name,omitempty"`

	// ARMID is the Azure Resource Manager ID.
	ARMID *string `json:"armId,omitempty"`
}

// AddressSpace represents the address space of a virtual network.
type AddressSpace struct {
	// AddressPrefixes is the list of address prefixes.
	AddressPrefixes []string `json:"addressPrefixes,omitempty"`
}

// DhcpOptions represents DHCP options.
type DhcpOptions struct {
	// DNSServers is the list of DNS servers.
	DNSServers []string `json:"dnsServers,omitempty"`
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

// Subnet represents an ASO Subnet resource.
// +kubebuilder:object:root=true
type Subnet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SubnetSpec   `json:"spec,omitempty"`
	Status SubnetStatus `json:"status,omitempty"`
}

// SubnetSpec defines the desired state of a Subnet.
type SubnetSpec struct {
	// Owner is the virtual network owner reference.
	Owner *VirtualNetworkReference `json:"owner,omitempty"`

	// AzureName is the name of the resource in Azure.
	AzureName *string `json:"azureName,omitempty"`

	// AddressPrefix is the address prefix for the subnet.
	AddressPrefix *string `json:"addressPrefix,omitempty"`

	// AddressPrefixes is a list of address prefixes.
	AddressPrefixes []string `json:"addressPrefixes,omitempty"`

	// NetworkSecurityGroupReference references a NetworkSecurityGroup.
	NetworkSecurityGroupReference *NetworkSecurityGroupReference `json:"networkSecurityGroupReference,omitempty"`

	// ServiceEndpoints specifies service endpoints.
	ServiceEndpoints []ServiceEndpointPropertiesFormat `json:"serviceEndpoints,omitempty"`

	// PrivateEndpointNetworkPolicies specifies private endpoint network policies.
	PrivateEndpointNetworkPolicies *string `json:"privateEndpointNetworkPolicies,omitempty"`

	// PrivateLinkServiceNetworkPolicies specifies private link service network policies.
	PrivateLinkServiceNetworkPolicies *string `json:"privateLinkServiceNetworkPolicies,omitempty"`

	// Delegations specifies subnet delegations.
	Delegations []Delegation `json:"delegations,omitempty"`
}

// SubnetStatus defines the observed state of a Subnet.
type SubnetStatus struct {
	// Conditions represent the latest available observations.
	Conditions []Condition `json:"conditions,omitempty"`

	// ID is the Azure resource ID.
	ID *string `json:"id,omitempty"`

	// ProvisioningState is the current provisioning state.
	ProvisioningState *string `json:"provisioningState,omitempty"`
}

// VirtualNetworkReference references a Virtual Network.
type VirtualNetworkReference struct {
	// Name is the name of the virtual network resource.
	Name string `json:"name,omitempty"`

	// ARMID is the Azure Resource Manager ID.
	ARMID *string `json:"armId,omitempty"`
}

// NetworkSecurityGroupReference references a NetworkSecurityGroup.
type NetworkSecurityGroupReference struct {
	// Name is the name of the NSG resource.
	Name *string `json:"name,omitempty"`

	// ARMID is the Azure Resource Manager ID.
	ARMID *string `json:"armId,omitempty"`
}

// ServiceEndpointPropertiesFormat represents a service endpoint.
type ServiceEndpointPropertiesFormat struct {
	// Service is the service name.
	Service *string `json:"service,omitempty"`

	// Locations are the locations for the service endpoint.
	Locations []string `json:"locations,omitempty"`
}

// Delegation represents a subnet delegation.
type Delegation struct {
	// Name is the name of the delegation.
	Name *string `json:"name,omitempty"`

	// ServiceName is the name of the service to delegate to.
	ServiceName *string `json:"serviceName,omitempty"`
}

// NetworkSecurityGroup represents an ASO Network Security Group resource.
// +kubebuilder:object:root=true
type NetworkSecurityGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NetworkSecurityGroupSpec   `json:"spec,omitempty"`
	Status NetworkSecurityGroupStatus `json:"status,omitempty"`
}

// NetworkSecurityGroupSpec defines the desired state of a Network Security Group.
type NetworkSecurityGroupSpec struct {
	// Owner is the resource group owner reference.
	Owner *ResourceGroupReference `json:"owner,omitempty"`

	// AzureName is the name of the resource in Azure.
	AzureName *string `json:"azureName,omitempty"`

	// Location is the Azure region.
	Location *string `json:"location,omitempty"`

	// SecurityRules specifies the security rules.
	SecurityRules []SecurityRule `json:"securityRules,omitempty"`

	// Tags are key-value pairs.
	Tags map[string]string `json:"tags,omitempty"`
}

// NetworkSecurityGroupStatus defines the observed state of a Network Security Group.
type NetworkSecurityGroupStatus struct {
	// Conditions represent the latest available observations.
	Conditions []Condition `json:"conditions,omitempty"`

	// ID is the Azure resource ID.
	ID *string `json:"id,omitempty"`

	// ProvisioningState is the current provisioning state.
	ProvisioningState *string `json:"provisioningState,omitempty"`
}

// SecurityRule represents a security rule in a network security group.
type SecurityRule struct {
	// Name is the name of the security rule.
	Name *string `json:"name,omitempty"`

	// Priority is the priority of the rule (100-4096).
	Priority *int `json:"priority,omitempty"`

	// Direction is the direction of the rule (Inbound or Outbound).
	Direction *string `json:"direction,omitempty"`

	// Access is the access type (Allow or Deny).
	Access *string `json:"access,omitempty"`

	// Protocol is the protocol (* , Tcp, Udp, Icmp).
	Protocol *string `json:"protocol,omitempty"`

	// SourcePortRange is the source port range.
	SourcePortRange *string `json:"sourcePortRange,omitempty"`

	// DestinationPortRange is the destination port range.
	DestinationPortRange *string `json:"destinationPortRange,omitempty"`

	// SourceAddressPrefix is the source address prefix.
	SourceAddressPrefix *string `json:"sourceAddressPrefix,omitempty"`

	// DestinationAddressPrefix is the destination address prefix.
	DestinationAddressPrefix *string `json:"destinationAddressPrefix,omitempty"`

	// Description is a description of the rule.
	Description *string `json:"description,omitempty"`
}
