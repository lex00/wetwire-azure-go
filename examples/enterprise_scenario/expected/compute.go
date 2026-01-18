// Package main demonstrates compute infrastructure for a multi-tier enterprise application.
package main

import (
	"github.com/lex00/wetwire-azure-go/resources/compute"
	"github.com/lex00/wetwire-azure-go/resources/network"
)

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

// Helper function for string pointers
func strPtr(s string) *string { return &s }
