// Package network provides Azure network resource types
package network

// VirtualNetwork represents a Microsoft.Network/virtualNetworks resource
type VirtualNetwork struct {
	// Name is the name of the virtual network
	Name string `json:"name"`

	// Type is the resource type
	Type string `json:"type"`

	// APIVersion is the API version to use for this resource
	APIVersion string `json:"apiVersion"`

	// Location is the Azure region where the resource will be created
	Location string `json:"location"`

	// Tags are key-value pairs to organize resources
	Tags map[string]string `json:"tags,omitempty"`

	// Properties contains the properties of the virtual network
	Properties VirtualNetworkProperties `json:"properties"`
}

// VirtualNetworkProperties represents the properties of a virtual network
type VirtualNetworkProperties struct {
	// AddressSpace specifies the address space
	AddressSpace AddressSpace `json:"addressSpace"`

	// Subnets specifies the subnets in the virtual network
	Subnets []Subnet `json:"subnets,omitempty"`

	// DhcpOptions specifies DHCP options
	DhcpOptions *DhcpOptions `json:"dhcpOptions,omitempty"`

	// EnableDdosProtection indicates whether DDoS protection is enabled
	EnableDdosProtection *bool `json:"enableDdosProtection,omitempty"`

	// EnableVmProtection indicates whether VM protection is enabled
	EnableVmProtection *bool `json:"enableVmProtection,omitempty"`
}

// AddressSpace represents the address space of a virtual network
type AddressSpace struct {
	// AddressPrefixes is the list of address prefixes
	AddressPrefixes []string `json:"addressPrefixes"`
}

// DhcpOptions represents DHCP options
type DhcpOptions struct {
	// DNSServers is the list of DNS servers
	DNSServers []string `json:"dnsServers,omitempty"`
}

// Subnet represents a subnet in a virtual network
type Subnet struct {
	// Name is the name of the subnet
	Name string `json:"name"`

	// Properties contains the properties of the subnet
	Properties SubnetProperties `json:"properties"`
}

// SubnetProperties represents the properties of a subnet
type SubnetProperties struct {
	// AddressPrefix is the address prefix for the subnet
	AddressPrefix string `json:"addressPrefix"`

	// NetworkSecurityGroup specifies the network security group
	NetworkSecurityGroup *SubResource `json:"networkSecurityGroup,omitempty"`

	// RouteTable specifies the route table
	RouteTable *SubResource `json:"routeTable,omitempty"`

	// ServiceEndpoints specifies service endpoints
	ServiceEndpoints []ServiceEndpoint `json:"serviceEndpoints,omitempty"`

	// PrivateEndpointNetworkPolicies specifies private endpoint network policies
	PrivateEndpointNetworkPolicies *string `json:"privateEndpointNetworkPolicies,omitempty"`

	// PrivateLinkServiceNetworkPolicies specifies private link service network policies
	PrivateLinkServiceNetworkPolicies *string `json:"privateLinkServiceNetworkPolicies,omitempty"`

	// Delegations specifies subnet delegations
	Delegations []Delegation `json:"delegations,omitempty"`
}

// ServiceEndpoint represents a service endpoint
type ServiceEndpoint struct {
	// Service is the service name
	Service string `json:"service"`

	// Locations are the locations for the service endpoint
	Locations []string `json:"locations,omitempty"`
}

// Delegation represents a subnet delegation
type Delegation struct {
	// Name is the name of the delegation
	Name string `json:"name"`

	// Properties contains the delegation properties
	Properties DelegationProperties `json:"properties"`
}

// DelegationProperties represents the properties of a delegation
type DelegationProperties struct {
	// ServiceName is the name of the service to delegate to
	ServiceName string `json:"serviceName"`
}

// NetworkInterface represents a Microsoft.Network/networkInterfaces resource
type NetworkInterface struct {
	// Name is the name of the network interface
	Name string `json:"name"`

	// Type is the resource type
	Type string `json:"type"`

	// APIVersion is the API version to use for this resource
	APIVersion string `json:"apiVersion"`

	// Location is the Azure region where the resource will be created
	Location string `json:"location"`

	// Tags are key-value pairs to organize resources
	Tags map[string]string `json:"tags,omitempty"`

	// Properties contains the properties of the network interface
	Properties NetworkInterfaceProperties `json:"properties"`
}

// NetworkInterfaceProperties represents the properties of a network interface
type NetworkInterfaceProperties struct {
	// IPConfigurations specifies the IP configurations
	IPConfigurations []IPConfiguration `json:"ipConfigurations,omitempty"`

	// NetworkSecurityGroup specifies the network security group
	NetworkSecurityGroup *SubResource `json:"networkSecurityGroup,omitempty"`

	// EnableAcceleratedNetworking indicates whether accelerated networking is enabled
	EnableAcceleratedNetworking *bool `json:"enableAcceleratedNetworking,omitempty"`

	// EnableIPForwarding indicates whether IP forwarding is enabled
	EnableIPForwarding *bool `json:"enableIPForwarding,omitempty"`

	// DNSSettings specifies the DNS settings
	DNSSettings *NetworkInterfaceDNSSettings `json:"dnsSettings,omitempty"`
}

// IPConfiguration represents an IP configuration of a network interface
type IPConfiguration struct {
	// Name is the name of the IP configuration
	Name string `json:"name"`

	// Properties contains the properties of the IP configuration
	Properties IPConfigurationProperties `json:"properties"`
}

// IPConfigurationProperties represents the properties of an IP configuration
type IPConfigurationProperties struct {
	// Subnet specifies the subnet
	Subnet *SubResource `json:"subnet,omitempty"`

	// PublicIPAddress specifies the public IP address
	PublicIPAddress *SubResource `json:"publicIPAddress,omitempty"`

	// PrivateIPAddress specifies the private IP address
	PrivateIPAddress *string `json:"privateIPAddress,omitempty"`

	// PrivateIPAllocationMethod specifies the allocation method (Static or Dynamic)
	PrivateIPAllocationMethod *string `json:"privateIPAllocationMethod,omitempty"`

	// Primary indicates whether this is the primary IP configuration
	Primary *bool `json:"primary,omitempty"`
}

// NetworkInterfaceDNSSettings represents DNS settings for a network interface
type NetworkInterfaceDNSSettings struct {
	// DNSServers is the list of DNS servers
	DNSServers []string `json:"dnsServers,omitempty"`

	// InternalDNSNameLabel is the internal DNS name label
	InternalDNSNameLabel *string `json:"internalDnsNameLabel,omitempty"`
}

// PublicIPAddress represents a Microsoft.Network/publicIPAddresses resource
type PublicIPAddress struct {
	// Name is the name of the public IP address
	Name string `json:"name"`

	// Type is the resource type
	Type string `json:"type"`

	// APIVersion is the API version to use for this resource
	APIVersion string `json:"apiVersion"`

	// Location is the Azure region where the resource will be created
	Location string `json:"location"`

	// Tags are key-value pairs to organize resources
	Tags map[string]string `json:"tags,omitempty"`

	// SKU defines the SKU for the public IP address
	SKU PublicIPSKU `json:"sku"`

	// Properties contains the properties of the public IP address
	Properties PublicIPAddressProperties `json:"properties"`

	// Zones defines availability zones for the public IP address
	Zones []string `json:"zones,omitempty"`
}

// PublicIPSKU represents the SKU of a public IP address
type PublicIPSKU struct {
	// Name is the SKU name (Basic or Standard)
	Name string `json:"name"`

	// Tier is the SKU tier (Regional or Global)
	Tier *string `json:"tier,omitempty"`
}

// PublicIPAddressProperties represents the properties of a public IP address
type PublicIPAddressProperties struct {
	// PublicIPAllocationMethod specifies the allocation method (Static or Dynamic)
	PublicIPAllocationMethod string `json:"publicIPAllocationMethod"`

	// PublicIPAddressVersion specifies the IP version (IPv4 or IPv6)
	PublicIPAddressVersion *string `json:"publicIPAddressVersion,omitempty"`

	// DNSSettings specifies the DNS settings
	DNSSettings *PublicIPDNSSettings `json:"dnsSettings,omitempty"`

	// IdleTimeoutInMinutes specifies the idle timeout
	IdleTimeoutInMinutes *int `json:"idleTimeoutInMinutes,omitempty"`
}

// PublicIPDNSSettings represents DNS settings for a public IP address
type PublicIPDNSSettings struct {
	// DomainNameLabel is the DNS name label
	DomainNameLabel *string `json:"domainNameLabel,omitempty"`

	// Fqdn is the fully qualified domain name
	Fqdn *string `json:"fqdn,omitempty"`

	// ReverseFqdn is the reverse fully qualified domain name
	ReverseFqdn *string `json:"reverseFqdn,omitempty"`
}

// NetworkSecurityGroup represents a Microsoft.Network/networkSecurityGroups resource
type NetworkSecurityGroup struct {
	// Name is the name of the network security group
	Name string `json:"name"`

	// Type is the resource type
	Type string `json:"type"`

	// APIVersion is the API version to use for this resource
	APIVersion string `json:"apiVersion"`

	// Location is the Azure region where the resource will be created
	Location string `json:"location"`

	// Tags are key-value pairs to organize resources
	Tags map[string]string `json:"tags,omitempty"`

	// Properties contains the properties of the network security group
	Properties NetworkSecurityGroupProperties `json:"properties"`
}

// NetworkSecurityGroupProperties represents the properties of a network security group
type NetworkSecurityGroupProperties struct {
	// SecurityRules specifies the security rules
	SecurityRules []SecurityRule `json:"securityRules,omitempty"`
}

// SecurityRule represents a security rule in a network security group
type SecurityRule struct {
	// Name is the name of the security rule
	Name string `json:"name"`

	// Properties contains the properties of the security rule
	Properties SecurityRuleProperties `json:"properties"`
}

// SecurityRuleProperties represents the properties of a security rule
type SecurityRuleProperties struct {
	// Priority is the priority of the rule (100-4096)
	Priority int `json:"priority"`

	// Direction is the direction of the rule (Inbound or Outbound)
	Direction string `json:"direction"`

	// Access is the access type (Allow or Deny)
	Access string `json:"access"`

	// Protocol is the protocol (* , Tcp, Udp, Icmp, Esp, Ah)
	Protocol string `json:"protocol"`

	// SourcePortRange is the source port range
	SourcePortRange string `json:"sourcePortRange,omitempty"`

	// DestinationPortRange is the destination port range
	DestinationPortRange string `json:"destinationPortRange,omitempty"`

	// SourceAddressPrefix is the source address prefix
	SourceAddressPrefix string `json:"sourceAddressPrefix,omitempty"`

	// DestinationAddressPrefix is the destination address prefix
	DestinationAddressPrefix string `json:"destinationAddressPrefix,omitempty"`

	// SourcePortRanges is the list of source port ranges
	SourcePortRanges []string `json:"sourcePortRanges,omitempty"`

	// DestinationPortRanges is the list of destination port ranges
	DestinationPortRanges []string `json:"destinationPortRanges,omitempty"`

	// SourceAddressPrefixes is the list of source address prefixes
	SourceAddressPrefixes []string `json:"sourceAddressPrefixes,omitempty"`

	// DestinationAddressPrefixes is the list of destination address prefixes
	DestinationAddressPrefixes []string `json:"destinationAddressPrefixes,omitempty"`

	// Description is a description of the rule
	Description *string `json:"description,omitempty"`
}

// SubResource represents a reference to another resource
type SubResource struct {
	// ID is the resource ID
	ID *string `json:"id,omitempty"`
}

// NewVirtualNetwork creates a new virtual network with required fields
func NewVirtualNetwork(name, location string, addressPrefixes []string) *VirtualNetwork {
	return &VirtualNetwork{
		Name:       name,
		Type:       "Microsoft.Network/virtualNetworks",
		APIVersion: "2021-05-01",
		Location:   location,
		Properties: VirtualNetworkProperties{
			AddressSpace: AddressSpace{
				AddressPrefixes: addressPrefixes,
			},
		},
	}
}

// WithTags adds tags to the virtual network
func (v *VirtualNetwork) WithTags(tags map[string]string) *VirtualNetwork {
	v.Tags = tags
	return v
}

// WithSubnet adds a subnet to the virtual network
func (v *VirtualNetwork) WithSubnet(name, addressPrefix string) *VirtualNetwork {
	v.Properties.Subnets = append(v.Properties.Subnets, Subnet{
		Name: name,
		Properties: SubnetProperties{
			AddressPrefix: addressPrefix,
		},
	})
	return v
}

// WithDNSServers sets the DNS servers for the virtual network
func (v *VirtualNetwork) WithDNSServers(dnsServers []string) *VirtualNetwork {
	v.Properties.DhcpOptions = &DhcpOptions{
		DNSServers: dnsServers,
	}
	return v
}

// NewSubnet creates a new subnet with required fields
func NewSubnet(name, addressPrefix string) *Subnet {
	return &Subnet{
		Name: name,
		Properties: SubnetProperties{
			AddressPrefix: addressPrefix,
		},
	}
}

// WithNSG attaches a network security group to the subnet
func (s *Subnet) WithNSG(nsgID string) *Subnet {
	s.Properties.NetworkSecurityGroup = &SubResource{ID: &nsgID}
	return s
}

// WithServiceEndpoint adds a service endpoint to the subnet
func (s *Subnet) WithServiceEndpoint(service string, locations []string) *Subnet {
	s.Properties.ServiceEndpoints = append(s.Properties.ServiceEndpoints, ServiceEndpoint{
		Service:   service,
		Locations: locations,
	})
	return s
}

// NewNetworkInterface creates a new network interface with required fields
func NewNetworkInterface(name, location string) *NetworkInterface {
	return &NetworkInterface{
		Name:       name,
		Type:       "Microsoft.Network/networkInterfaces",
		APIVersion: "2021-05-01",
		Location:   location,
		Properties: NetworkInterfaceProperties{},
	}
}

// WithTags adds tags to the network interface
func (n *NetworkInterface) WithTags(tags map[string]string) *NetworkInterface {
	n.Tags = tags
	return n
}

// WithIPConfiguration adds an IP configuration to the network interface
func (n *NetworkInterface) WithIPConfiguration(name, subnetID string, primary bool) *NetworkInterface {
	n.Properties.IPConfigurations = append(n.Properties.IPConfigurations, IPConfiguration{
		Name: name,
		Properties: IPConfigurationProperties{
			Subnet:  &SubResource{ID: &subnetID},
			Primary: &primary,
		},
	})
	return n
}

// WithPublicIP attaches a public IP to an existing IP configuration
func (n *NetworkInterface) WithPublicIP(ipConfigName, publicIPID string) *NetworkInterface {
	for i := range n.Properties.IPConfigurations {
		if n.Properties.IPConfigurations[i].Name == ipConfigName {
			n.Properties.IPConfigurations[i].Properties.PublicIPAddress = &SubResource{ID: &publicIPID}
			break
		}
	}
	return n
}

// WithNSG attaches a network security group to the network interface
func (n *NetworkInterface) WithNSG(nsgID string) *NetworkInterface {
	n.Properties.NetworkSecurityGroup = &SubResource{ID: &nsgID}
	return n
}

// WithAcceleratedNetworking enables or disables accelerated networking
func (n *NetworkInterface) WithAcceleratedNetworking(enabled bool) *NetworkInterface {
	n.Properties.EnableAcceleratedNetworking = &enabled
	return n
}

// NewPublicIPAddress creates a new public IP address with required fields
func NewPublicIPAddress(name, location, allocationMethod, skuName string) *PublicIPAddress {
	return &PublicIPAddress{
		Name:       name,
		Type:       "Microsoft.Network/publicIPAddresses",
		APIVersion: "2021-05-01",
		Location:   location,
		SKU: PublicIPSKU{
			Name: skuName,
		},
		Properties: PublicIPAddressProperties{
			PublicIPAllocationMethod: allocationMethod,
		},
	}
}

// WithTags adds tags to the public IP address
func (p *PublicIPAddress) WithTags(tags map[string]string) *PublicIPAddress {
	p.Tags = tags
	return p
}

// WithDNSLabel sets the DNS label for the public IP address
func (p *PublicIPAddress) WithDNSLabel(label string) *PublicIPAddress {
	p.Properties.DNSSettings = &PublicIPDNSSettings{
		DomainNameLabel: &label,
	}
	return p
}

// WithIdleTimeout sets the idle timeout for the public IP address
func (p *PublicIPAddress) WithIdleTimeout(minutes int) *PublicIPAddress {
	p.Properties.IdleTimeoutInMinutes = &minutes
	return p
}

// WithZones sets the availability zones for the public IP address
func (p *PublicIPAddress) WithZones(zones []string) *PublicIPAddress {
	p.Zones = zones
	return p
}

// NewNetworkSecurityGroup creates a new network security group with required fields
func NewNetworkSecurityGroup(name, location string) *NetworkSecurityGroup {
	return &NetworkSecurityGroup{
		Name:       name,
		Type:       "Microsoft.Network/networkSecurityGroups",
		APIVersion: "2021-05-01",
		Location:   location,
		Properties: NetworkSecurityGroupProperties{},
	}
}

// WithTags adds tags to the network security group
func (n *NetworkSecurityGroup) WithTags(tags map[string]string) *NetworkSecurityGroup {
	n.Tags = tags
	return n
}

// WithRule adds a security rule to the network security group
func (n *NetworkSecurityGroup) WithRule(name string, priority int, direction, access, protocol, srcPort, dstPort, srcAddr, dstAddr string) *NetworkSecurityGroup {
	n.Properties.SecurityRules = append(n.Properties.SecurityRules, SecurityRule{
		Name: name,
		Properties: SecurityRuleProperties{
			Priority:                 priority,
			Direction:                direction,
			Access:                   access,
			Protocol:                 protocol,
			SourcePortRange:          srcPort,
			DestinationPortRange:     dstPort,
			SourceAddressPrefix:      srcAddr,
			DestinationAddressPrefix: dstAddr,
		},
	})
	return n
}
