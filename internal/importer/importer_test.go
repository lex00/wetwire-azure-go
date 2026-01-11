package importer

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseARMTemplate_EmptyTemplate(t *testing.T) {
	input := `{
		"$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"resources": []
	}`

	template, err := ParseARMTemplate([]byte(input))
	require.NoError(t, err)
	assert.NotNil(t, template)
	assert.Empty(t, template.Resources)
}

func TestParseARMTemplate_SingleResource(t *testing.T) {
	input := `{
		"$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"resources": [
			{
				"type": "Microsoft.Storage/storageAccounts",
				"apiVersion": "2021-04-01",
				"name": "mystorageaccount",
				"location": "eastus",
				"kind": "StorageV2",
				"sku": {
					"name": "Standard_LRS"
				}
			}
		]
	}`

	template, err := ParseARMTemplate([]byte(input))
	require.NoError(t, err)
	require.Len(t, template.Resources, 1)

	res := template.Resources[0]
	assert.Equal(t, "Microsoft.Storage/storageAccounts", res.Type)
	assert.Equal(t, "mystorageaccount", res.Name)
	assert.Equal(t, "eastus", res.Location)
	assert.Equal(t, "StorageV2", res.Kind)
}

func TestParseARMTemplate_WithDependsOn(t *testing.T) {
	input := `{
		"$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"resources": [
			{
				"type": "Microsoft.Storage/storageAccounts",
				"apiVersion": "2021-04-01",
				"name": "mystorageaccount",
				"location": "eastus"
			},
			{
				"type": "Microsoft.Compute/virtualMachines",
				"apiVersion": "2021-07-01",
				"name": "myvm",
				"location": "eastus",
				"dependsOn": [
					"[resourceId('Microsoft.Storage/storageAccounts', 'mystorageaccount')]"
				]
			}
		]
	}`

	template, err := ParseARMTemplate([]byte(input))
	require.NoError(t, err)
	require.Len(t, template.Resources, 2)

	vmRes := template.Resources[1]
	assert.Equal(t, "myvm", vmRes.Name)
	require.Len(t, vmRes.DependsOn, 1)
	assert.Contains(t, vmRes.DependsOn[0], "mystorageaccount")
}

func TestCamelToPascal(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"name", "Name"},
		{"location", "Location"},
		{"vmSize", "VMSize"},
		{"apiVersion", "APIVersion"},
		{"storageAccountType", "StorageAccountType"},
		{"osDisk", "OSDisk"},
		{"osProfile", "OSProfile"},
		{"ipAddress", "IPAddress"},
		{"httpsOnly", "HTTPSOnly"},
		{"enableHttpsTrafficOnly", "EnableHTTPSTrafficOnly"},
		{"vnetSubnetId", "VNetSubnetID"},
		{"diskSizeGB", "DiskSizeGB"},
		{"sshPublicKey", "SSHPublicKey"},
		{"skuName", "SKUName"},
		{"nicId", "NICID"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := CamelToPascal(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestResourceTypeToPackage(t *testing.T) {
	tests := []struct {
		resourceType string
		pkgName      string
		typeName     string
	}{
		{
			resourceType: "Microsoft.Storage/storageAccounts",
			pkgName:      "storage",
			typeName:     "StorageAccount",
		},
		{
			resourceType: "Microsoft.Compute/virtualMachines",
			pkgName:      "compute",
			typeName:     "VirtualMachine",
		},
		{
			resourceType: "Microsoft.Network/virtualNetworks",
			pkgName:      "network",
			typeName:     "VirtualNetwork",
		},
		{
			resourceType: "Microsoft.Network/networkInterfaces",
			pkgName:      "network",
			typeName:     "NetworkInterface",
		},
		{
			resourceType: "Microsoft.KeyVault/vaults",
			pkgName:      "keyvault",
			typeName:     "Vault",
		},
	}

	for _, tt := range tests {
		t.Run(tt.resourceType, func(t *testing.T) {
			pkgName, typeName := ResourceTypeToPackage(tt.resourceType)
			assert.Equal(t, tt.pkgName, pkgName)
			assert.Equal(t, tt.typeName, typeName)
		})
	}
}

func TestGenerateVarName(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"mystorageaccount", "Mystorageaccount"},
		{"my-storage-account", "MyStorageAccount"},
		{"my_storage_account", "MyStorageAccount"},
		{"my-vm-01", "MyVM01"},
		{"MyVM", "MyVM"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateVarName(tt.name)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateGoCode_SingleStorageAccount(t *testing.T) {
	input := `{
		"$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"resources": [
			{
				"type": "Microsoft.Storage/storageAccounts",
				"apiVersion": "2021-04-01",
				"name": "mystorageaccount",
				"location": "eastus",
				"kind": "StorageV2",
				"sku": {
					"name": "Standard_LRS"
				}
			}
		]
	}`

	template, err := ParseARMTemplate([]byte(input))
	require.NoError(t, err)

	code, err := GenerateGoCode(template, "infra")
	require.NoError(t, err)

	// Verify package declaration
	assert.Contains(t, code, "package infra")

	// Verify import
	assert.Contains(t, code, `"github.com/lex00/wetwire-azure-go/resources/storage"`)

	// Verify variable declaration
	assert.Contains(t, code, "var Mystorageaccount = storage.StorageAccount{")
	assert.Contains(t, code, `Name:     "mystorageaccount"`)
	assert.Contains(t, code, `Location: "eastus"`)
	assert.Contains(t, code, `Kind:     "StorageV2"`)
}

func TestGenerateGoCode_VirtualMachine(t *testing.T) {
	input := `{
		"$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"resources": [
			{
				"type": "Microsoft.Compute/virtualMachines",
				"apiVersion": "2021-07-01",
				"name": "my-vm",
				"location": "eastus",
				"properties": {
					"hardwareProfile": {
						"vmSize": "Standard_DS1_v2"
					}
				}
			}
		]
	}`

	template, err := ParseARMTemplate([]byte(input))
	require.NoError(t, err)

	code, err := GenerateGoCode(template, "infra")
	require.NoError(t, err)

	// Verify import
	assert.Contains(t, code, `"github.com/lex00/wetwire-azure-go/resources/compute"`)

	// Verify variable
	assert.Contains(t, code, "var MyVM = compute.VirtualMachine{")
	assert.Contains(t, code, `Name:     "my-vm"`)
	assert.Contains(t, code, `Location: "eastus"`)
}

func TestGenerateGoCode_WithDependsOn(t *testing.T) {
	input := `{
		"$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"resources": [
			{
				"type": "Microsoft.Storage/storageAccounts",
				"apiVersion": "2021-04-01",
				"name": "my-storage",
				"location": "eastus"
			},
			{
				"type": "Microsoft.Compute/virtualMachines",
				"apiVersion": "2021-07-01",
				"name": "my-vm",
				"location": "eastus",
				"dependsOn": [
					"[resourceId('Microsoft.Storage/storageAccounts', 'my-storage')]"
				]
			}
		]
	}`

	template, err := ParseARMTemplate([]byte(input))
	require.NoError(t, err)

	code, err := GenerateGoCode(template, "infra")
	require.NoError(t, err)

	// The storage account should be generated
	assert.Contains(t, code, "var MyStorage = storage.StorageAccount{")

	// The VM should reference the storage account
	// dependsOn in ARM becomes a Go comment with reference
	assert.Contains(t, code, "var MyVM = compute.VirtualMachine{")
	assert.Contains(t, code, "// DependsOn: MyStorage")
}

func TestGenerateGoCode_MultipleResources(t *testing.T) {
	input := `{
		"$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"resources": [
			{
				"type": "Microsoft.Storage/storageAccounts",
				"apiVersion": "2021-04-01",
				"name": "mystorageaccount",
				"location": "eastus"
			},
			{
				"type": "Microsoft.Compute/virtualMachines",
				"apiVersion": "2021-07-01",
				"name": "myvm",
				"location": "westus"
			}
		]
	}`

	template, err := ParseARMTemplate([]byte(input))
	require.NoError(t, err)

	code, err := GenerateGoCode(template, "infra")
	require.NoError(t, err)

	// Verify both imports are present
	assert.Contains(t, code, `"github.com/lex00/wetwire-azure-go/resources/storage"`)
	assert.Contains(t, code, `"github.com/lex00/wetwire-azure-go/resources/compute"`)

	// Verify both resources
	assert.Contains(t, code, "var Mystorageaccount = storage.StorageAccount{")
	assert.Contains(t, code, "var Myvm = compute.VirtualMachine{")
}

func TestGenerateGoCode_NestedProperties(t *testing.T) {
	input := `{
		"$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"resources": [
			{
				"type": "Microsoft.Compute/virtualMachines",
				"apiVersion": "2021-07-01",
				"name": "my-vm",
				"location": "eastus",
				"properties": {
					"hardwareProfile": {
						"vmSize": "Standard_DS1_v2"
					},
					"storageProfile": {
						"osDisk": {
							"createOption": "FromImage",
							"diskSizeGB": 128
						}
					}
				}
			}
		]
	}`

	template, err := ParseARMTemplate([]byte(input))
	require.NoError(t, err)

	code, err := GenerateGoCode(template, "infra")
	require.NoError(t, err)

	// Verify nested properties are generated
	assert.Contains(t, code, "Properties:")
	assert.Contains(t, code, "HardwareProfile:")
	assert.Contains(t, code, `VMSize: "Standard_DS1_v2"`)
	assert.Contains(t, code, "StorageProfile:")
	assert.Contains(t, code, "OSDisk:")
	assert.Contains(t, code, `CreateOption: "FromImage"`)
}

func TestGenerateGoCode_WithTags(t *testing.T) {
	input := `{
		"$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"resources": [
			{
				"type": "Microsoft.Storage/storageAccounts",
				"apiVersion": "2021-04-01",
				"name": "mystorageaccount",
				"location": "eastus",
				"tags": {
					"environment": "production",
					"team": "platform"
				}
			}
		]
	}`

	template, err := ParseARMTemplate([]byte(input))
	require.NoError(t, err)

	code, err := GenerateGoCode(template, "infra")
	require.NoError(t, err)

	// Verify tags are generated
	assert.Contains(t, code, "Tags: map[string]string{")
	assert.Contains(t, code, `"environment": "production"`)
	assert.Contains(t, code, `"team": "platform"`)
}

func TestParseARMTemplate_InvalidJSON(t *testing.T) {
	input := `{invalid json}`

	_, err := ParseARMTemplate([]byte(input))
	assert.Error(t, err)
}

func TestExtractDependencyName(t *testing.T) {
	tests := []struct {
		dependsOn string
		expected  string
	}{
		{
			dependsOn: "[resourceId('Microsoft.Storage/storageAccounts', 'mystorageaccount')]",
			expected:  "mystorageaccount",
		},
		{
			dependsOn: "[resourceId('Microsoft.Network/virtualNetworks', 'my-vnet')]",
			expected:  "my-vnet",
		},
		{
			dependsOn: "Microsoft.Storage/storageAccounts/mystorageaccount",
			expected:  "mystorageaccount",
		},
		{
			dependsOn: "[concat('Microsoft.Storage/storageAccounts/', variables('storageAccountName'))]",
			expected:  "", // Cannot extract from concat expressions
		},
	}

	for _, tt := range tests {
		t.Run(tt.dependsOn, func(t *testing.T) {
			result := ExtractDependencyName(tt.dependsOn)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateGoCode_SKU(t *testing.T) {
	input := `{
		"$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"resources": [
			{
				"type": "Microsoft.Storage/storageAccounts",
				"apiVersion": "2021-04-01",
				"name": "mystorageaccount",
				"location": "eastus",
				"kind": "StorageV2",
				"sku": {
					"name": "Standard_LRS",
					"tier": "Standard"
				}
			}
		]
	}`

	template, err := ParseARMTemplate([]byte(input))
	require.NoError(t, err)

	code, err := GenerateGoCode(template, "infra")
	require.NoError(t, err)

	// Verify SKU is generated
	assert.Contains(t, code, "SKU: storage.SKU{")
	assert.Contains(t, code, `Name: "Standard_LRS"`)
}

func TestGenerateImports(t *testing.T) {
	resourceTypes := []string{
		"Microsoft.Storage/storageAccounts",
		"Microsoft.Compute/virtualMachines",
		"Microsoft.Network/virtualNetworks",
	}

	imports := GenerateImports(resourceTypes)

	assert.Contains(t, imports, `"github.com/lex00/wetwire-azure-go/resources/storage"`)
	assert.Contains(t, imports, `"github.com/lex00/wetwire-azure-go/resources/compute"`)
	assert.Contains(t, imports, `"github.com/lex00/wetwire-azure-go/resources/network"`)

	// Should not have duplicate imports
	count := strings.Count(imports, "storage")
	assert.Equal(t, 1, count)
}
