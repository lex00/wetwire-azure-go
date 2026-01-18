// Package main demonstrates network infrastructure for a multi-tier enterprise application.
package main

import (
	"github.com/lex00/wetwire-azure-go/resources/network"
)

// Common configuration
var (
	location = "eastus"
	tags     = map[string]string{
		"environment": "production",
		"project":     "enterprise-app",
		"cost-center": "engineering",
	}
)

// AppVNet is the main virtual network with subnets for gateway, web, app, and data tiers.
var AppVNet = network.VirtualNetwork{
	Name:       "enterprise-vnet",
	Type:       "Microsoft.Network/virtualNetworks",
	APIVersion: "2021-05-01",
	Location:   location,
	Tags:       tags,
	Properties: network.VirtualNetworkProperties{
		AddressSpace: network.AddressSpace{
			AddressPrefixes: []string{"10.0.0.0/16"},
		},
		Subnets: []network.Subnet{
			{
				Name: "gateway-subnet",
				Properties: network.SubnetProperties{
					AddressPrefix: "10.0.0.0/24",
				},
			},
			{
				Name: "web-subnet",
				Properties: network.SubnetProperties{
					AddressPrefix: "10.0.1.0/24",
				},
			},
			{
				Name: "app-subnet",
				Properties: network.SubnetProperties{
					AddressPrefix: "10.0.2.0/24",
				},
			},
			{
				Name: "data-subnet",
				Properties: network.SubnetProperties{
					AddressPrefix: "10.0.3.0/24",
					ServiceEndpoints: []network.ServiceEndpoint{
						{Service: "Microsoft.Storage", Locations: []string{"eastus"}},
						{Service: "Microsoft.Sql", Locations: []string{"eastus"}},
					},
				},
			},
		},
	},
}

// GatewayNSG is the network security group for the Application Gateway subnet.
// Allows HTTP (80), HTTPS (443) from Internet and management ports for App Gateway.
var GatewayNSG = network.NetworkSecurityGroup{
	Name:       "gateway-nsg",
	Type:       "Microsoft.Network/networkSecurityGroups",
	APIVersion: "2021-05-01",
	Location:   location,
	Tags:       tags,
	Properties: network.NetworkSecurityGroupProperties{
		SecurityRules: []network.SecurityRule{
			{
				Name: "allow-http",
				Properties: network.SecurityRuleProperties{
					Priority:                 100,
					Direction:                "Inbound",
					Access:                   "Allow",
					Protocol:                 "Tcp",
					SourcePortRange:          "*",
					DestinationPortRange:     "80",
					SourceAddressPrefix:      "Internet",
					DestinationAddressPrefix: "*",
				},
			},
			{
				Name: "allow-https",
				Properties: network.SecurityRuleProperties{
					Priority:                 110,
					Direction:                "Inbound",
					Access:                   "Allow",
					Protocol:                 "Tcp",
					SourcePortRange:          "*",
					DestinationPortRange:     "443",
					SourceAddressPrefix:      "Internet",
					DestinationAddressPrefix: "*",
				},
			},
			{
				Name: "allow-appgw-mgmt",
				Properties: network.SecurityRuleProperties{
					Priority:                 120,
					Direction:                "Inbound",
					Access:                   "Allow",
					Protocol:                 "Tcp",
					SourcePortRange:          "*",
					DestinationPortRange:     "65200-65535",
					SourceAddressPrefix:      "GatewayManager",
					DestinationAddressPrefix: "*",
				},
			},
		},
	},
}

// WebNSG is the network security group for the web tier.
// Only allows traffic from gateway subnet on port 8080.
var WebNSG = network.NetworkSecurityGroup{
	Name:       "web-nsg",
	Type:       "Microsoft.Network/networkSecurityGroups",
	APIVersion: "2021-05-01",
	Location:   location,
	Tags:       tags,
	Properties: network.NetworkSecurityGroupProperties{
		SecurityRules: []network.SecurityRule{
			{
				Name: "allow-from-gateway",
				Properties: network.SecurityRuleProperties{
					Priority:                 100,
					Direction:                "Inbound",
					Access:                   "Allow",
					Protocol:                 "Tcp",
					SourcePortRange:          "*",
					DestinationPortRange:     "8080",
					SourceAddressPrefix:      "10.0.0.0/24",
					DestinationAddressPrefix: "*",
				},
			},
		},
	},
}

// AppNSG is the network security group for the app tier.
// Only allows traffic from web subnet on port 8080.
var AppNSG = network.NetworkSecurityGroup{
	Name:       "app-nsg",
	Type:       "Microsoft.Network/networkSecurityGroups",
	APIVersion: "2021-05-01",
	Location:   location,
	Tags:       tags,
	Properties: network.NetworkSecurityGroupProperties{
		SecurityRules: []network.SecurityRule{
			{
				Name: "allow-from-web",
				Properties: network.SecurityRuleProperties{
					Priority:                 100,
					Direction:                "Inbound",
					Access:                   "Allow",
					Protocol:                 "Tcp",
					SourcePortRange:          "*",
					DestinationPortRange:     "8080",
					SourceAddressPrefix:      "10.0.1.0/24",
					DestinationAddressPrefix: "*",
				},
			},
		},
	},
}

// DataNSG is the network security group for the data tier.
// Only allows SQL traffic from app subnet on port 1433.
var DataNSG = network.NetworkSecurityGroup{
	Name:       "data-nsg",
	Type:       "Microsoft.Network/networkSecurityGroups",
	APIVersion: "2021-05-01",
	Location:   location,
	Tags:       tags,
	Properties: network.NetworkSecurityGroupProperties{
		SecurityRules: []network.SecurityRule{
			{
				Name: "allow-sql-from-app",
				Properties: network.SecurityRuleProperties{
					Priority:                 100,
					Direction:                "Inbound",
					Access:                   "Allow",
					Protocol:                 "Tcp",
					SourcePortRange:          "*",
					DestinationPortRange:     "1433",
					SourceAddressPrefix:      "10.0.2.0/24",
					DestinationAddressPrefix: "*",
				},
			},
		},
	},
}

// AppGwPublicIP is the public IP address for the Application Gateway.
// Uses Standard SKU with static allocation and zone redundancy.
var AppGwPublicIP = network.PublicIPAddress{
	Name:       "appgw-pip",
	Type:       "Microsoft.Network/publicIPAddresses",
	APIVersion: "2021-05-01",
	Location:   location,
	Tags:       tags,
	SKU: network.PublicIPSKU{
		Name: "Standard",
	},
	Properties: network.PublicIPAddressProperties{
		PublicIPAllocationMethod: "Static",
	},
	Zones: []string{"1", "2", "3"},
}
