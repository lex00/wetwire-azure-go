// Package main demonstrates deploying an Azure Linux Virtual Machine.
package main

import (
	"github.com/lex00/wetwire-azure-go/resources/compute"
)

// Helper variables for pointer fields
var (
	publisher     = "Canonical"
	offer         = "0001-com-ubuntu-server-jammy"
	sku           = "22_04-lts-gen2"
	imageVersion  = "latest"
	osDiskName    = "mylinuxvm-osdisk"
	caching       = "ReadWrite"
	storageType   = "Premium_LRS"
	computerName  = "mylinuxvm"
	adminUsername = "azureuser"
	adminPassword = "P@ssw0rd1234!"
	networkIfaceID = "/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.Network/networkInterfaces/my-nic"
	isPrimary     = true
)

// LinuxVM defines a Linux virtual machine with Ubuntu 22.04.
// This example demonstrates a complete VM configuration including:
// - Hardware profile with VM size
// - Storage profile with OS image and managed disk
// - OS profile with admin credentials
// - Network profile with NIC reference
var LinuxVM = compute.VirtualMachine{
	Name:     "my-linux-vm",
	Location: "eastus",
	Properties: compute.VirtualMachineProperties{
		HardwareProfile: compute.HardwareProfile{
			VMSize: "Standard_B2s",
		},
		StorageProfile: compute.StorageProfile{
			ImageReference: &compute.ImageReference{
				Publisher: &publisher,
				Offer:     &offer,
				SKU:       &sku,
				Version:   &imageVersion,
			},
			OSDisk: compute.OSDisk{
				Name:         &osDiskName,
				CreateOption: "FromImage",
				Caching:      &caching,
				ManagedDisk: &compute.ManagedDiskParameters{
					StorageAccountType: &storageType,
				},
			},
		},
		OSProfile: &compute.OSProfile{
			ComputerName:  &computerName,
			AdminUsername: &adminUsername,
			AdminPassword: &adminPassword,
		},
		NetworkProfile: compute.NetworkProfile{
			NetworkInterfaces: []compute.NetworkInterfaceReference{
				{
					ID:      networkIfaceID,
					Primary: &isPrimary,
				},
			},
		},
	},
}

// main is required for a valid Go program but not used by wetwire-azure.
// Resources are discovered from package-level variable declarations.
func main() {
	// wetwire-azure build discovers resources via AST parsing
	// No runtime execution is needed
}
