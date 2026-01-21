// Package aks_golden provides a production-ready AKS cluster example.
//
// This example demonstrates best practices for deploying an AKS cluster
// with proper networking, security, and observability configurations.
package aks_golden

import (
	"github.com/lex00/wetwire-azure-go/resources/network"
)

// Location is the Azure region for all resources.
const Location = "eastus"

// Helper variables for pointer fields
var (
	privateEndpointPoliciesDisabled = "Disabled"
)

// VNet is the virtual network for the AKS cluster.
var VNet = network.VirtualNetwork{
	Name:       "aks-golden-vnet",
	Type:       "Microsoft.Network/virtualNetworks",
	APIVersion: "2021-05-01",
	Location:   Location,
	Tags: map[string]string{
		"Environment": "production",
		"ManagedBy":   "wetwire",
	},
	Properties: network.VirtualNetworkProperties{
		AddressSpace: network.AddressSpace{
			AddressPrefixes: []string{"10.0.0.0/16"},
		},
		Subnets: []network.Subnet{
			AKSSubnet,
			AppGatewaySubnet,
			PrivateEndpointSubnet,
		},
	},
}

// AKSSubnet is the subnet for AKS nodes.
// Using a /22 CIDR provides ~1000 IPs for nodes and pods with Azure CNI.
var AKSSubnet = network.Subnet{
	Name: "aks-subnet",
	Properties: network.SubnetProperties{
		AddressPrefix: "10.0.0.0/22",
		ServiceEndpoints: []network.ServiceEndpoint{
			{Service: "Microsoft.Storage", Locations: []string{Location}},
			{Service: "Microsoft.KeyVault", Locations: []string{Location}},
			{Service: "Microsoft.ContainerRegistry", Locations: []string{Location}},
		},
	},
}

// AppGatewaySubnet is the subnet for Azure Application Gateway (AGIC).
// AGIC requires its own subnet for deployment.
var AppGatewaySubnet = network.Subnet{
	Name: "appgw-subnet",
	Properties: network.SubnetProperties{
		AddressPrefix: "10.0.4.0/24",
	},
}

// PrivateEndpointSubnet is the subnet for private endpoints.
// Used for secure access to Azure PaaS services.
var PrivateEndpointSubnet = network.Subnet{
	Name: "private-endpoint-subnet",
	Properties: network.SubnetProperties{
		AddressPrefix:                  "10.0.5.0/24",
		PrivateEndpointNetworkPolicies: &privateEndpointPoliciesDisabled,
	},
}

// NSG is the network security group for AKS nodes.
var NSG = network.NetworkSecurityGroup{
	Name:       "aks-golden-nsg",
	Type:       "Microsoft.Network/networkSecurityGroups",
	APIVersion: "2021-05-01",
	Location:   Location,
	Tags: map[string]string{
		"Environment": "production",
		"ManagedBy":   "wetwire",
	},
	Properties: network.NetworkSecurityGroupProperties{
		SecurityRules: []network.SecurityRule{
			AllowHTTPS,
			AllowHTTP,
			DenyAllInbound,
		},
	},
}

// AllowHTTPS allows HTTPS traffic from the internet.
var AllowHTTPS = network.SecurityRule{
	Name: "AllowHTTPS",
	Properties: network.SecurityRuleProperties{
		Priority:                 100,
		Direction:                "Inbound",
		Access:                   "Allow",
		Protocol:                 "Tcp",
		SourcePortRange:          "*",
		DestinationPortRange:     "443",
		SourceAddressPrefix:      "Internet",
		DestinationAddressPrefix: "*",
	},
}

// AllowHTTP allows HTTP traffic from the internet.
var AllowHTTP = network.SecurityRule{
	Name: "AllowHTTP",
	Properties: network.SecurityRuleProperties{
		Priority:                 101,
		Direction:                "Inbound",
		Access:                   "Allow",
		Protocol:                 "Tcp",
		SourcePortRange:          "*",
		DestinationPortRange:     "80",
		SourceAddressPrefix:      "Internet",
		DestinationAddressPrefix: "*",
	},
}

// DenyAllInbound denies all other inbound traffic.
var DenyAllInbound = network.SecurityRule{
	Name: "DenyAllInbound",
	Properties: network.SecurityRuleProperties{
		Priority:                 4096,
		Direction:                "Inbound",
		Access:                   "Deny",
		Protocol:                 "*",
		SourcePortRange:          "*",
		DestinationPortRange:     "*",
		SourceAddressPrefix:      "*",
		DestinationAddressPrefix: "*",
	},
}
