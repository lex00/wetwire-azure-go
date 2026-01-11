package serialize

import (
	"encoding/json"
	"testing"

	"github.com/lex00/wetwire-azure-go/intrinsics"
	"github.com/lex00/wetwire-azure-go/resources/compute"
	"github.com/lex00/wetwire-azure-go/resources/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSimpleStorageAccountSerialization tests basic resource serialization
func TestSimpleStorageAccountSerialization(t *testing.T) {
	sa := storage.NewStorageAccount("mystorageaccount", "eastus", "StorageV2", "Standard_LRS")

	result := ToARMResource(sa)

	assert.Equal(t, "mystorageaccount", result["name"])
	assert.Equal(t, "Microsoft.Storage/storageAccounts", result["type"])
	assert.Equal(t, "2021-04-01", result["apiVersion"])
	assert.Equal(t, "eastus", result["location"])
	assert.Equal(t, "StorageV2", result["kind"])

	sku, ok := result["sku"].(map[string]any)
	require.True(t, ok, "sku should be a map")
	assert.Equal(t, "Standard_LRS", sku["name"])
}

// TestStorageAccountWithProperties tests nested properties serialization
func TestStorageAccountWithProperties(t *testing.T) {
	httpsOnly := true
	minTLS := "TLS1_2"

	sa := storage.NewStorageAccount("mystorageaccount", "eastus", "StorageV2", "Standard_LRS")
	sa.WithHTTPSOnly(httpsOnly).WithMinTLSVersion(minTLS)

	result := ToARMResource(sa)

	props, ok := result["properties"].(map[string]any)
	require.True(t, ok, "properties should be a map")
	assert.Equal(t, true, props["supportsHttpsTrafficOnly"])
	assert.Equal(t, "TLS1_2", props["minimumTlsVersion"])
}

// TestStorageAccountWithTags tests tags serialization
func TestStorageAccountWithTags(t *testing.T) {
	tags := map[string]string{
		"Environment": "Production",
		"CostCenter":  "IT",
	}

	sa := storage.NewStorageAccount("mystorageaccount", "eastus", "StorageV2", "Standard_LRS")
	sa.WithTags(tags)

	result := ToARMResource(sa)

	resultTags, ok := result["tags"].(map[string]any)
	require.True(t, ok, "tags should be a map")
	assert.Equal(t, "Production", resultTags["Environment"])
	assert.Equal(t, "IT", resultTags["CostCenter"])
}

// TestVirtualMachineWithNestedStructs tests complex nested structure serialization
func TestVirtualMachineWithNestedStructs(t *testing.T) {
	vm := compute.NewVirtualMachine("testvm", "eastus", "Standard_DS2_v2")
	publisher := "Canonical"
	offer := "UbuntuServer"
	sku := "18.04-LTS"
	version := "latest"

	vm.WithImage(publisher, offer, sku, version)

	result := ToARMResource(vm)

	assert.Equal(t, "testvm", result["name"])
	assert.Equal(t, "Microsoft.Compute/virtualMachines", result["type"])

	props, ok := result["properties"].(map[string]any)
	require.True(t, ok, "properties should be a map")

	storageProfile, ok := props["storageProfile"].(map[string]any)
	require.True(t, ok, "storageProfile should be a map")

	imageRef, ok := storageProfile["imageReference"].(map[string]any)
	require.True(t, ok, "imageReference should be a map")
	assert.Equal(t, "Canonical", imageRef["publisher"])
	assert.Equal(t, "UbuntuServer", imageRef["offer"])
	assert.Equal(t, "18.04-LTS", imageRef["sku"])
	assert.Equal(t, "latest", imageRef["version"])
}

// TestOmitEmptyFields tests that nil pointer fields are omitted
func TestOmitEmptyFields(t *testing.T) {
	sa := storage.NewStorageAccount("mystorageaccount", "eastus", "StorageV2", "Standard_LRS")
	// Don't set any optional properties

	result := ToARMResource(sa)

	// Properties should not be present if nil
	_, hasProps := result["properties"]
	assert.False(t, hasProps, "nil properties should be omitted")

	_, hasTags := result["tags"]
	assert.False(t, hasTags, "nil tags should be omitted")
}

// TestArraySerialization tests array field serialization
func TestArraySerialization(t *testing.T) {
	vm := compute.NewVirtualMachine("testvm", "eastus", "Standard_DS2_v2")
	vm.WithNetworkInterface("/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/networkInterfaces/nic1", true)
	vm.WithNetworkInterface("/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/networkInterfaces/nic2", false)

	result := ToARMResource(vm)

	props, ok := result["properties"].(map[string]any)
	require.True(t, ok, "properties should be a map")

	netProfile, ok := props["networkProfile"].(map[string]any)
	require.True(t, ok, "networkProfile should be a map")

	nics, ok := netProfile["networkInterfaces"].([]any)
	require.True(t, ok, "networkInterfaces should be an array")
	assert.Len(t, nics, 2)

	nic1, ok := nics[0].(map[string]any)
	require.True(t, ok)
	assert.Contains(t, nic1["id"], "nic1")
	assert.Equal(t, true, nic1["primary"])

	nic2, ok := nics[1].(map[string]any)
	require.True(t, ok)
	assert.Contains(t, nic2["id"], "nic2")
	assert.Equal(t, false, nic2["primary"])
}

// TestIntrinsicResourceId tests ResourceId intrinsic serialization
func TestIntrinsicResourceId(t *testing.T) {
	resourceID := intrinsics.ResourceId("Microsoft.Storage/storageAccounts", "mystorageaccount")

	result := SerializeValue(resourceID)

	assert.Equal(t, "[resourceId('Microsoft.Storage/storageAccounts', 'mystorageaccount')]", result)
}

// TestIntrinsicReference tests Reference intrinsic serialization
func TestIntrinsicReference(t *testing.T) {
	ref := intrinsics.Ref("mystorageaccount", "2021-04-01")

	result := SerializeValue(ref)

	assert.Equal(t, "[reference('mystorageaccount', '2021-04-01')]", result)
}

// TestIntrinsicReferenceProperty tests Reference with property serialization
func TestIntrinsicReferenceProperty(t *testing.T) {
	ref := intrinsics.RefProperty("mystorageaccount", "2021-04-01", "primaryEndpoints.blob")

	result := SerializeValue(ref)

	assert.Equal(t, "[reference('mystorageaccount', '2021-04-01').primaryEndpoints.blob]", result)
}

// TestIntrinsicParameter tests Parameter intrinsic serialization
func TestIntrinsicParameter(t *testing.T) {
	param := intrinsics.Parameters("location")

	result := SerializeValue(param)

	assert.Equal(t, "[parameters('location')]", result)
}

// TestIntrinsicVariable tests Variable intrinsic serialization
func TestIntrinsicVariable(t *testing.T) {
	variable := intrinsics.Variables("storageAccountName")

	result := SerializeValue(variable)

	assert.Equal(t, "[variables('storageAccountName')]", result)
}

// TestIntrinsicResourceGroup tests ResourceGroup intrinsic serialization
func TestIntrinsicResourceGroup(t *testing.T) {
	rg := intrinsics.ResourceGroup()

	result := SerializeValue(rg)

	assert.Equal(t, "[resourceGroup()]", result)
}

// TestIntrinsicResourceGroupProperty tests ResourceGroup with property
func TestIntrinsicResourceGroupProperty(t *testing.T) {
	rg := intrinsics.ResourceGroupValue{Property: "location"}

	result := SerializeValue(rg)

	assert.Equal(t, "[resourceGroup().location]", result)
}

// TestResourceWithIntrinsics tests a resource using intrinsics in its fields
// Note: This test demonstrates how intrinsics would work when resource types support any fields
func TestResourceWithIntrinsics(t *testing.T) {
	// For now, we test intrinsics separately since resource types use string for Location
	// When resource types are updated to use any for Location, this will work directly
	locationIntrinsic := intrinsics.ResourceGroupValue{Property: "location"}

	result := SerializeValue(locationIntrinsic)

	// Location should be serialized as ARM expression
	location, ok := result.(string)
	require.True(t, ok, "location should be a string")
	assert.Equal(t, "[resourceGroup().location]", location)
}

// TestToARMTemplate tests full ARM template generation
func TestToARMTemplate(t *testing.T) {
	sa := storage.NewStorageAccount("mystorageaccount", "eastus", "StorageV2", "Standard_LRS")

	template := ToARMTemplate([]any{sa})

	assert.Equal(t, "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#", template["$schema"])
	assert.Equal(t, "1.0.0.0", template["contentVersion"])

	resources, ok := template["resources"].([]any)
	require.True(t, ok, "resources should be an array")
	require.Len(t, resources, 1)

	resource, ok := resources[0].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "mystorageaccount", resource["name"])
	assert.Equal(t, "Microsoft.Storage/storageAccounts", resource["type"])
}

// TestToARMTemplateMultipleResources tests ARM template with multiple resources
func TestToARMTemplateMultipleResources(t *testing.T) {
	sa := storage.NewStorageAccount("mystorageaccount", "eastus", "StorageV2", "Standard_LRS")
	vm := compute.NewVirtualMachine("myvm", "eastus", "Standard_DS2_v2")

	template := ToARMTemplate([]any{sa, vm})

	resources, ok := template["resources"].([]any)
	require.True(t, ok, "resources should be an array")
	require.Len(t, resources, 2)
}

// TestToARMTemplateJSON tests JSON marshaling of ARM template
func TestToARMTemplateJSON(t *testing.T) {
	sa := storage.NewStorageAccount("mystorageaccount", "eastus", "StorageV2", "Standard_LRS")

	jsonBytes, err := ToARMTemplateJSON([]any{sa})
	require.NoError(t, err)

	var parsed map[string]any
	err = json.Unmarshal(jsonBytes, &parsed)
	require.NoError(t, err)

	assert.Equal(t, "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#", parsed["$schema"])
	assert.Equal(t, "1.0.0.0", parsed["contentVersion"])

	resources, ok := parsed["resources"].([]any)
	require.True(t, ok)
	require.Len(t, resources, 1)
}

// TestComplexNestedStructures tests deeply nested structures
func TestComplexNestedStructures(t *testing.T) {
	allowBlob := true
	sa := &storage.StorageAccount{
		Name:       "mystorageaccount",
		Type:       "Microsoft.Storage/storageAccounts",
		APIVersion: "2021-04-01",
		Location:   "eastus",
		Kind:       "StorageV2",
		SKU: storage.SKU{
			Name: "Standard_LRS",
		},
		Properties: &storage.StorageAccountProperties{
			AllowBlobPublicAccess: &allowBlob,
			NetworkRuleSet: &storage.NetworkRuleSet{
				DefaultAction: "Deny",
				IPRules: []storage.IPRule{
					{Value: "192.168.1.0/24"},
					{Value: "10.0.0.0/8"},
				},
			},
		},
	}

	result := ToARMResource(sa)

	props, ok := result["properties"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, true, props["allowBlobPublicAccess"])

	networkAcls, ok := props["networkAcls"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "Deny", networkAcls["defaultAction"])

	ipRules, ok := networkAcls["ipRules"].([]any)
	require.True(t, ok)
	require.Len(t, ipRules, 2)

	rule1, ok := ipRules[0].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "192.168.1.0/24", rule1["value"])
}

// TestJSONStructTags tests that JSON struct tags are respected
func TestJSONStructTags(t *testing.T) {
	sa := storage.NewStorageAccount("mystorageaccount", "eastus", "StorageV2", "Standard_LRS")

	result := ToARMResource(sa)

	// Verify that fields use json tag names
	assert.Contains(t, result, "name")
	assert.Contains(t, result, "type")
	assert.Contains(t, result, "apiVersion")
	assert.Contains(t, result, "location")
	assert.Contains(t, result, "kind")
	assert.Contains(t, result, "sku")
}

// TestPointerFields tests serialization of pointer fields
func TestPointerFields(t *testing.T) {
	tier := "Standard"
	sa := &storage.StorageAccount{
		Name:       "mystorageaccount",
		Type:       "Microsoft.Storage/storageAccounts",
		APIVersion: "2021-04-01",
		Location:   "eastus",
		Kind:       "StorageV2",
		SKU: storage.SKU{
			Name: "Standard_LRS",
			Tier: &tier,
		},
	}

	result := ToARMResource(sa)

	sku, ok := result["sku"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "Standard_LRS", sku["name"])
	assert.Equal(t, "Standard", sku["tier"])
}

// TestEmptyArraysOmitted tests that empty arrays are omitted
func TestEmptyArraysOmitted(t *testing.T) {
	vm := compute.NewVirtualMachine("testvm", "eastus", "Standard_DS2_v2")
	// NetworkInterfaces is initialized as empty slice

	result := ToARMResource(vm)

	props, ok := result["properties"].(map[string]any)
	require.True(t, ok)

	netProfile, ok := props["networkProfile"].(map[string]any)
	require.True(t, ok)

	// Empty array should be omitted
	_, hasNics := netProfile["networkInterfaces"]
	assert.False(t, hasNics, "empty networkInterfaces array should be omitted")
}

// TestIntrinsicInNestedField tests intrinsic serialization
// Note: This demonstrates how intrinsics would work when resource types support any fields
func TestIntrinsicInNestedField(t *testing.T) {
	nicID := intrinsics.ResourceId("Microsoft.Network/networkInterfaces", "mynic")

	result := SerializeValue(nicID)

	id, ok := result.(string)
	require.True(t, ok)
	assert.Equal(t, "[resourceId('Microsoft.Network/networkInterfaces', 'mynic')]", id)
}

// TestStructWithAnyFieldAndIntrinsic tests a custom struct with any fields containing intrinsics
func TestStructWithAnyFieldAndIntrinsic(t *testing.T) {
	// Define a test struct that uses any for fields that can contain intrinsics
	type FlexibleResource struct {
		Name     string `json:"name"`
		Location any    `json:"location"`
		Tags     any    `json:"tags,omitempty"`
	}

	// Create a resource with an intrinsic in the Location field
	resource := FlexibleResource{
		Name:     "testresource",
		Location: intrinsics.ResourceGroupValue{Property: "location"},
	}

	result := ToARMResource(resource)

	assert.Equal(t, "testresource", result["name"])

	location, ok := result["location"].(string)
	require.True(t, ok, "location should be a string ARM expression")
	assert.Equal(t, "[resourceGroup().location]", location)
}

// TestStructWithIntrinsicInMap tests intrinsics within map fields
func TestStructWithIntrinsicInMap(t *testing.T) {
	type ResourceWithMap struct {
		Name       string         `json:"name"`
		Properties map[string]any `json:"properties"`
	}

	resource := ResourceWithMap{
		Name: "testresource",
		Properties: map[string]any{
			"location":      intrinsics.ResourceGroupValue{Property: "location"},
			"staticValue":   "eastus",
			"parameterized": intrinsics.Parameters("environmentName"),
		},
	}

	result := ToARMResource(resource)

	props, ok := result["properties"].(map[string]any)
	require.True(t, ok)

	assert.Equal(t, "[resourceGroup().location]", props["location"])
	assert.Equal(t, "eastus", props["staticValue"])
	assert.Equal(t, "[parameters('environmentName')]", props["parameterized"])
}

// TestStructWithIntrinsicInSlice tests intrinsics within slice fields
func TestStructWithIntrinsicInSlice(t *testing.T) {
	type ResourceWithSlice struct {
		Name string `json:"name"`
		IDs  []any  `json:"ids"`
	}

	resource := ResourceWithSlice{
		Name: "testresource",
		IDs: []any{
			intrinsics.ResourceId("Microsoft.Network/virtualNetworks", "vnet1"),
			intrinsics.ResourceId("Microsoft.Network/virtualNetworks", "vnet2"),
			"/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet3",
		},
	}

	result := ToARMResource(resource)

	ids, ok := result["ids"].([]any)
	require.True(t, ok)
	require.Len(t, ids, 3)

	assert.Equal(t, "[resourceId('Microsoft.Network/virtualNetworks', 'vnet1')]", ids[0])
	assert.Equal(t, "[resourceId('Microsoft.Network/virtualNetworks', 'vnet2')]", ids[1])
	assert.Equal(t, "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet3", ids[2])
}
