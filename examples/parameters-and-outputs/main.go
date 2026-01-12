// Package main demonstrates ARM template parameters, variables, and intrinsic functions.
// This example shows how to:
// - Reference parameters in resource definitions
// - Use ARM intrinsic functions (resourceId, reference, concat, etc.)
// - Create dynamic resource names
// - Reference resource group properties
package main

import (
	"github.com/lex00/wetwire-azure-go/intrinsics"
	"github.com/lex00/wetwire-azure-go/resources/network"
	"github.com/lex00/wetwire-azure-go/resources/storage"
)

// ============================================================================
// Using Parameters and Variables
// ============================================================================

// ParameterizedStorage demonstrates using the Parameters intrinsic
// to reference ARM template parameters in resource properties.
// The location will be resolved at deployment time from the parameter.
var ParameterizedStorage = storage.StorageAccount{
	Name:       "mystorageaccount",
	Type:       "Microsoft.Storage/storageAccounts",
	APIVersion: "2021-04-01",
	// Location uses the ARM parameters() function
	// This generates: "[parameters('location')]"
	Location: intrinsics.Parameters("location").ARMExpression(),
	Kind:     "StorageV2",
	SKU: storage.SKU{
		Name: "Standard_LRS",
	},
	Properties: &storage.StorageAccountProperties{
		EnableHTTPSTrafficOnly: boolPtr(true),
		MinimumTLSVersion:      strPtr("TLS1_2"),
	},
}

// ============================================================================
// Using Resource Group Functions
// ============================================================================

// ResourceGroupLocation demonstrates using the resourceGroup() function
// to inherit the resource group's location automatically.
var ResourceGroupLocation = storage.StorageAccount{
	Name:       "rglocationstorage",
	Type:       "Microsoft.Storage/storageAccounts",
	APIVersion: "2021-04-01",
	// Uses resourceGroup().location - resolves to deployment resource group's location
	Location: intrinsics.ResourceGroup().ARMExpression(),
	Kind:     "StorageV2",
	SKU: storage.SKU{
		Name: "Standard_LRS",
	},
}

// ============================================================================
// Using ResourceId Function
// ============================================================================

// ExampleNIC demonstrates using resourceId() to reference other resources.
// This is useful for creating dependencies between resources.
var ExampleNIC = network.NetworkInterface{
	Name:       "example-nic",
	Type:       "Microsoft.Network/networkInterfaces",
	APIVersion: "2021-05-01",
	Location:   "eastus",
	Properties: network.NetworkInterfaceProperties{
		IPConfigurations: []network.IPConfiguration{
			{
				Name: "ipconfig1",
				Properties: network.IPConfigurationProperties{
					Subnet: &network.SubResource{
						// Reference subnet using resourceId function
						// Generates: "[resourceId('Microsoft.Network/virtualNetworks/subnets', 'myVnet', 'default')]"
						ID: strPtr(intrinsics.ResourceId(
							"Microsoft.Network/virtualNetworks/subnets",
							"myVnet",
							"default",
						).ARMExpression()),
					},
					PrivateIPAllocationMethod: strPtr("Dynamic"),
				},
			},
		},
	},
}

// ============================================================================
// Using Variables Function
// ============================================================================

// VariableBasedStorage shows using the variables() function
// to reference ARM template variables.
var VariableBasedStorage = storage.StorageAccount{
	Name:       "varstorage",
	Type:       "Microsoft.Storage/storageAccounts",
	APIVersion: "2021-04-01",
	// Reference a variable defined in ARM template
	Location: intrinsics.Variables("defaultLocation").ARMExpression(),
	Kind:     "StorageV2",
	SKU: storage.SKU{
		// SKU name from variable
		Name: "Standard_LRS",
	},
}

// ============================================================================
// Environment-Based Configuration Pattern
// ============================================================================

// Environment tags can be parameterized for different deployments
var envTags = map[string]string{
	"environment": "production",
	"managedBy":   "wetwire-azure",
}

// ProdStorage shows a production-ready configuration pattern
// with comprehensive settings that could be parameterized.
var ProdStorage = storage.StorageAccount{
	Name:       "prodstorage",
	Type:       "Microsoft.Storage/storageAccounts",
	APIVersion: "2021-04-01",
	Location:   "eastus",
	Tags:       envTags,
	Kind:       "StorageV2",
	SKU: storage.SKU{
		Name: "Standard_GRS", // Geo-redundant for production
	},
	Properties: &storage.StorageAccountProperties{
		EnableHTTPSTrafficOnly: boolPtr(true),
		MinimumTLSVersion:      strPtr("TLS1_2"),
		AllowBlobPublicAccess:  boolPtr(false),
		AllowSharedKeyAccess:   boolPtr(false), // Force Azure AD auth
		AccessTier:             strPtr("Hot"),
	},
}

// DevStorage shows a development configuration with relaxed settings
var DevStorage = storage.StorageAccount{
	Name:       "devstorage",
	Type:       "Microsoft.Storage/storageAccounts",
	APIVersion: "2021-04-01",
	Location:   "eastus",
	Tags: map[string]string{
		"environment": "development",
		"managedBy":   "wetwire-azure",
	},
	Kind: "StorageV2",
	SKU: storage.SKU{
		Name: "Standard_LRS", // Locally redundant for dev (cheaper)
	},
	Properties: &storage.StorageAccountProperties{
		EnableHTTPSTrafficOnly: boolPtr(true),
		MinimumTLSVersion:      strPtr("TLS1_2"),
	},
}

// Helper functions
func boolPtr(b bool) *bool    { return &b }
func strPtr(s string) *string { return &s }

func main() {
	// wetwire-azure build discovers resources via AST parsing
}
