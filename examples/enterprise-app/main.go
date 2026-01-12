// Package main demonstrates a multi-tier enterprise application deployment.
// This example shows a complete infrastructure setup with:
// - Virtual Network with multiple subnets (web, app, data tiers)
// - Network Security Groups with security rules
// - Public IP for load balancer
// - Network Interfaces for VMs
// - Virtual Machines in web and app tiers
// - Storage Account for data
package main

import (
	"github.com/lex00/wetwire-azure-go/resources/compute"
	"github.com/lex00/wetwire-azure-go/resources/network"
	"github.com/lex00/wetwire-azure-go/resources/storage"
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

// ============================================================================
// Networking Resources
// ============================================================================

// AppVNet is the main virtual network with three subnets for web, app, and data tiers.
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
					},
				},
			},
		},
	},
}

// WebNSG is the network security group for the web tier.
// Allows HTTP/HTTPS from internet, denies all other inbound traffic.
var WebNSG = network.NetworkSecurityGroup{
	Name:       "web-nsg",
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

// WebPublicIP is the public IP address for the web tier load balancer.
var WebPublicIP = network.PublicIPAddress{
	Name:       "web-pip",
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

// ============================================================================
// Compute Resources
// ============================================================================

// Helper variables for pointer fields
var (
	vmPublisher   = "Canonical"
	vmOffer       = "0001-com-ubuntu-server-jammy"
	vmSKU         = "22_04-lts-gen2"
	vmVersion     = "latest"
	caching       = "ReadWrite"
	storageType   = "Premium_LRS"
	adminUsername = "azureadmin"
	// In production, use SSH keys instead of passwords
	adminPassword = "P@ssw0rd1234!"
	isPrimary     = true
)

// WebVM is the web tier virtual machine.
var WebVM = compute.VirtualMachine{
	Name:       "web-vm-01",
	Type:       "Microsoft.Compute/virtualMachines",
	APIVersion: "2021-07-01",
	Location:   location,
	Tags:       tags,
	Properties: compute.VirtualMachineProperties{
		HardwareProfile: compute.HardwareProfile{
			VMSize: "Standard_B2s",
		},
		StorageProfile: compute.StorageProfile{
			ImageReference: &compute.ImageReference{
				Publisher: &vmPublisher,
				Offer:     &vmOffer,
				SKU:       &vmSKU,
				Version:   &vmVersion,
			},
			OSDisk: compute.OSDisk{
				CreateOption: "FromImage",
				Caching:      &caching,
				ManagedDisk: &compute.ManagedDiskParameters{
					StorageAccountType: &storageType,
				},
			},
		},
		OSProfile: &compute.OSProfile{
			ComputerName:  strPtr("webvm01"),
			AdminUsername: &adminUsername,
			AdminPassword: &adminPassword,
		},
		NetworkProfile: compute.NetworkProfile{
			NetworkInterfaces: []compute.NetworkInterfaceReference{
				{
					ID:      "/subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Network/networkInterfaces/web-nic-01",
					Primary: &isPrimary,
				},
			},
		},
	},
}

// AppVM is the application tier virtual machine.
var AppVM = compute.VirtualMachine{
	Name:       "app-vm-01",
	Type:       "Microsoft.Compute/virtualMachines",
	APIVersion: "2021-07-01",
	Location:   location,
	Tags:       tags,
	Properties: compute.VirtualMachineProperties{
		HardwareProfile: compute.HardwareProfile{
			VMSize: "Standard_D2s_v3",
		},
		StorageProfile: compute.StorageProfile{
			ImageReference: &compute.ImageReference{
				Publisher: &vmPublisher,
				Offer:     &vmOffer,
				SKU:       &vmSKU,
				Version:   &vmVersion,
			},
			OSDisk: compute.OSDisk{
				CreateOption: "FromImage",
				Caching:      &caching,
				ManagedDisk: &compute.ManagedDiskParameters{
					StorageAccountType: &storageType,
				},
			},
		},
		OSProfile: &compute.OSProfile{
			ComputerName:  strPtr("appvm01"),
			AdminUsername: &adminUsername,
			AdminPassword: &adminPassword,
		},
		NetworkProfile: compute.NetworkProfile{
			NetworkInterfaces: []compute.NetworkInterfaceReference{
				{
					ID:      "/subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Network/networkInterfaces/app-nic-01",
					Primary: &isPrimary,
				},
			},
		},
	},
}

// ============================================================================
// Storage Resources
// ============================================================================

// DataStorage is the storage account for application data.
// Configured with HTTPS-only, TLS 1.2, and restricted network access.
var DataStorage = storage.StorageAccount{
	Name:       "enterpriseappdata",
	Type:       "Microsoft.Storage/storageAccounts",
	APIVersion: "2021-04-01",
	Location:   location,
	Tags:       tags,
	Kind:       "StorageV2",
	SKU: storage.SKU{
		Name: "Standard_GRS",
	},
	Properties: &storage.StorageAccountProperties{
		EnableHTTPSTrafficOnly: boolPtr(true),
		MinimumTLSVersion:      strPtr("TLS1_2"),
		AllowBlobPublicAccess:  boolPtr(false),
		AllowSharedKeyAccess:   boolPtr(false),
		NetworkRuleSet: &storage.NetworkRuleSet{
			DefaultAction: "Deny",
			VirtualNetworkRules: []storage.VirtualNetworkRule{
				{
					ID:     "/subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Network/virtualNetworks/enterprise-vnet/subnets/data-subnet",
					Action: strPtr("Allow"),
				},
			},
		},
	},
}

// LogStorage is the storage account for diagnostic logs.
var LogStorage = storage.StorageAccount{
	Name:       "enterpriseapplogs",
	Type:       "Microsoft.Storage/storageAccounts",
	APIVersion: "2021-04-01",
	Location:   location,
	Tags:       tags,
	Kind:       "StorageV2",
	SKU: storage.SKU{
		Name: "Standard_LRS",
	},
	Properties: &storage.StorageAccountProperties{
		EnableHTTPSTrafficOnly: boolPtr(true),
		MinimumTLSVersion:      strPtr("TLS1_2"),
		AllowBlobPublicAccess:  boolPtr(false),
	},
}

// Helper functions
func boolPtr(b bool) *bool    { return &b }
func strPtr(s string) *string { return &s }

func main() {
	// wetwire-azure build discovers resources via AST parsing
	// No runtime execution is needed
}
