// Package main demonstrates a multi-tier enterprise application deployment on Azure.
// This example shows a complete infrastructure setup with network, compute, and storage resources.
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
// Network Resources
// ============================================================================

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

// ============================================================================
// Compute Resources
// ============================================================================

// Helper variables for VM configuration
var (
	vmPublisher   = "Canonical"
	vmOffer       = "0001-com-ubuntu-server-jammy"
	vmSKU         = "22_04-lts-gen2"
	vmVersion     = "latest"
	caching       = "ReadWrite"
	storageType   = "Premium_LRS"
	adminUsername = "azureadmin"
	adminPassword = "P@ssw0rd1234!" // In production, use SSH keys instead
	isPrimary     = true
)

// Web Tier Network Interfaces

// WebNIC01 is the network interface for the first web tier VM.
var WebNIC01 = network.NetworkInterface{
	Name:       "web-nic-01",
	Type:       "Microsoft.Network/networkInterfaces",
	APIVersion: "2021-05-01",
	Location:   location,
	Tags:       tags,
	Properties: network.NetworkInterfaceProperties{
		IPConfigurations: []network.IPConfiguration{
			{
				Name: "ipconfig1",
				Properties: network.IPConfigurationProperties{
					Subnet: &network.SubResource{
						ID: strPtr("/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.Network/virtualNetworks/enterprise-vnet/subnets/web-subnet"),
					},
					PrivateIPAllocationMethod: strPtr("Dynamic"),
				},
			},
		},
	},
}

// WebNIC02 is the network interface for the second web tier VM.
var WebNIC02 = network.NetworkInterface{
	Name:       "web-nic-02",
	Type:       "Microsoft.Network/networkInterfaces",
	APIVersion: "2021-05-01",
	Location:   location,
	Tags:       tags,
	Properties: network.NetworkInterfaceProperties{
		IPConfigurations: []network.IPConfiguration{
			{
				Name: "ipconfig1",
				Properties: network.IPConfigurationProperties{
					Subnet: &network.SubResource{
						ID: strPtr("/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.Network/virtualNetworks/enterprise-vnet/subnets/web-subnet"),
					},
					PrivateIPAllocationMethod: strPtr("Dynamic"),
				},
			},
		},
	},
}

// App Tier Network Interfaces

// AppNIC01 is the network interface for the first app tier VM.
var AppNIC01 = network.NetworkInterface{
	Name:       "app-nic-01",
	Type:       "Microsoft.Network/networkInterfaces",
	APIVersion: "2021-05-01",
	Location:   location,
	Tags:       tags,
	Properties: network.NetworkInterfaceProperties{
		IPConfigurations: []network.IPConfiguration{
			{
				Name: "ipconfig1",
				Properties: network.IPConfigurationProperties{
					Subnet: &network.SubResource{
						ID: strPtr("/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.Network/virtualNetworks/enterprise-vnet/subnets/app-subnet"),
					},
					PrivateIPAllocationMethod: strPtr("Dynamic"),
				},
			},
		},
	},
}

// AppNIC02 is the network interface for the second app tier VM.
var AppNIC02 = network.NetworkInterface{
	Name:       "app-nic-02",
	Type:       "Microsoft.Network/networkInterfaces",
	APIVersion: "2021-05-01",
	Location:   location,
	Tags:       tags,
	Properties: network.NetworkInterfaceProperties{
		IPConfigurations: []network.IPConfiguration{
			{
				Name: "ipconfig1",
				Properties: network.IPConfigurationProperties{
					Subnet: &network.SubResource{
						ID: strPtr("/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.Network/virtualNetworks/enterprise-vnet/subnets/app-subnet"),
					},
					PrivateIPAllocationMethod: strPtr("Dynamic"),
				},
			},
		},
	},
}

// Web Tier Virtual Machines

// WebVM01 is the first web tier virtual machine.
var WebVM01 = compute.VirtualMachine{
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
					ID:      "/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.Network/networkInterfaces/web-nic-01",
					Primary: &isPrimary,
				},
			},
		},
	},
}

// WebVM02 is the second web tier virtual machine.
var WebVM02 = compute.VirtualMachine{
	Name:       "web-vm-02",
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
			ComputerName:  strPtr("webvm02"),
			AdminUsername: &adminUsername,
			AdminPassword: &adminPassword,
		},
		NetworkProfile: compute.NetworkProfile{
			NetworkInterfaces: []compute.NetworkInterfaceReference{
				{
					ID:      "/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.Network/networkInterfaces/web-nic-02",
					Primary: &isPrimary,
				},
			},
		},
	},
}

// App Tier Virtual Machines

// AppVM01 is the first application tier virtual machine.
var AppVM01 = compute.VirtualMachine{
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
					ID:      "/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.Network/networkInterfaces/app-nic-01",
					Primary: &isPrimary,
				},
			},
		},
	},
}

// AppVM02 is the second application tier virtual machine.
var AppVM02 = compute.VirtualMachine{
	Name:       "app-vm-02",
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
			ComputerName:  strPtr("appvm02"),
			AdminUsername: &adminUsername,
			AdminPassword: &adminPassword,
		},
		NetworkProfile: compute.NetworkProfile{
			NetworkInterfaces: []compute.NetworkInterfaceReference{
				{
					ID:      "/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.Network/networkInterfaces/app-nic-02",
					Primary: &isPrimary,
				},
			},
		},
	},
}

// ============================================================================
// Storage Resources
// ============================================================================

// Helper variables for storage configuration
var (
	httpsOnly     = true
	tlsVersion    = "TLS1_2"
	noBlobPublic  = false
	noSharedKey   = false
	accessTierHot = "Hot"
	allowAction   = "Allow"
)

// DataStorage is the storage account for application data.
// Configured with geo-redundant storage, HTTPS-only, TLS 1.2, and restricted network access.
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
		AccessTier:             &accessTierHot,
		EnableHTTPSTrafficOnly: &httpsOnly,
		MinimumTLSVersion:      &tlsVersion,
		AllowBlobPublicAccess:  &noBlobPublic,
		AllowSharedKeyAccess:   &noSharedKey,
		NetworkRuleSet: &storage.NetworkRuleSet{
			DefaultAction: "Deny",
			VirtualNetworkRules: []storage.VirtualNetworkRule{
				{
					ID:     "/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.Network/virtualNetworks/enterprise-vnet/subnets/data-subnet",
					Action: &allowAction,
				},
				{
					ID:     "/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.Network/virtualNetworks/enterprise-vnet/subnets/app-subnet",
					Action: &allowAction,
				},
			},
		},
	},
}

// LogStorage is the storage account for diagnostic logs and metrics.
// Configured with locally-redundant storage, HTTPS-only, and TLS 1.2.
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
		AccessTier:             &accessTierHot,
		EnableHTTPSTrafficOnly: &httpsOnly,
		MinimumTLSVersion:      &tlsVersion,
		AllowBlobPublicAccess:  &noBlobPublic,
	},
}

// Helper function for string pointers
func strPtr(s string) *string { return &s }

func main() {
	// wetwire-azure build discovers resources via AST parsing
	// No runtime execution is needed
}
