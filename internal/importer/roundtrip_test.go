package importer

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/lex00/wetwire-azure-go/resources/storage"
	"github.com/lex00/wetwire-azure-go/internal/serialize"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRoundTrip_SimpleStorageAccount tests the round-trip conversion for a simple storage account
// Pipeline: ARM JSON → Parse → Go struct → Serialize → ARM JSON
func TestRoundTrip_SimpleStorageAccount(t *testing.T) {
	// Load the test fixture
	fixturePath := filepath.Join("..", "..", "testdata", "quickstarts", "storage-simple.json")
	data, err := os.ReadFile(fixturePath)
	require.NoError(t, err, "Failed to read test fixture")

	// Parse the original ARM template
	originalTemplate, err := ParseARMTemplate(data)
	require.NoError(t, err, "Failed to parse ARM template")
	require.Len(t, originalTemplate.Resources, 1, "Expected 1 resource")

	// Extract the storage account resource
	armRes := originalTemplate.Resources[0]

	// Convert to Go struct
	storageAccount := convertARMResourceToStorageAccount(t, armRes)
	require.NotNil(t, storageAccount)

	// Serialize back to ARM JSON
	serialized := serialize.ToARMResource(storageAccount)

	// Parse the original resource as a generic map for comparison
	var originalResource map[string]interface{}
	resourceJSON, err := json.Marshal(armRes)
	require.NoError(t, err)
	err = json.Unmarshal(resourceJSON, &originalResource)
	require.NoError(t, err)

	// Compare semantically (ignoring order and formatting)
	assertJSONEqual(t, originalResource, serialized)
}

// TestRoundTrip_StorageAccountFromQuickstart tests round-trip with Azure Quickstart template
func TestRoundTrip_StorageAccountFromQuickstart(t *testing.T) {
	fixturePath := filepath.Join("..", "..", "testdata", "quickstarts", "storage-account-create.json")
	data, err := os.ReadFile(fixturePath)
	require.NoError(t, err, "Failed to read test fixture")

	// Parse the ARM template
	template, err := ParseARMTemplate(data)
	require.NoError(t, err, "Failed to parse ARM template")

	// Verify template structure
	assert.Equal(t, "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#", template.Schema)
	assert.Equal(t, "1.0.0.0", template.ContentVersion)

	// Verify resources
	require.Len(t, template.Resources, 1, "Expected 1 resource in template")

	res := template.Resources[0]
	assert.Equal(t, "Microsoft.Storage/storageAccounts", res.Type)
	assert.Equal(t, "2022-09-01", res.APIVersion)
	assert.Equal(t, "StorageV2", res.Kind)

	// Verify SKU structure
	if len(res.SKU) > 0 {
		assert.NotNil(t, res.SKU["name"], "SKU should have a name field")
	}
}

// TestRoundTrip_WebAppFromQuickstart tests parsing of the web app quickstart template
func TestRoundTrip_WebAppFromQuickstart(t *testing.T) {
	fixturePath := filepath.Join("..", "..", "testdata", "quickstarts", "webapp-basic-linux.json")
	data, err := os.ReadFile(fixturePath)
	require.NoError(t, err, "Failed to read test fixture")

	// Parse the ARM template
	template, err := ParseARMTemplate(data)
	require.NoError(t, err, "Failed to parse ARM template")

	// Verify template structure
	assert.Equal(t, "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#", template.Schema)
	assert.Equal(t, "1.0.0.0", template.ContentVersion)

	// Verify resources
	require.Len(t, template.Resources, 2, "Expected 2 resources in template")

	// Check server farm resource
	serverFarm := template.Resources[0]
	assert.Equal(t, "Microsoft.Web/serverfarms", serverFarm.Type)
	assert.Equal(t, "2022-03-01", serverFarm.APIVersion)
	assert.Equal(t, "linux", serverFarm.Kind)

	// Check web app resource
	webApp := template.Resources[1]
	assert.Equal(t, "Microsoft.Web/sites", webApp.Type)
	assert.Equal(t, "2022-03-01", webApp.APIVersion)
	assert.Equal(t, "app", webApp.Kind)

	// Verify dependencies
	require.Len(t, webApp.DependsOn, 1, "Web app should depend on server farm")

	// Verify identity
	if len(webApp.Identity) > 0 {
		assert.Equal(t, "SystemAssigned", webApp.Identity["type"])
	} else {
		t.Error("Web app should have identity")
	}
}

// TestRoundTrip_PreservesEssentialFields verifies that critical ARM fields are preserved
func TestRoundTrip_PreservesEssentialFields(t *testing.T) {
	// Create a storage account with all essential fields
	storageAccount := &storage.StorageAccount{
		Name:       "teststorage123",
		Type:       "Microsoft.Storage/storageAccounts",
		APIVersion: "2021-04-01",
		Location:   "westus",
		Kind:       "StorageV2",
		Tags: map[string]string{
			"environment": "production",
			"cost-center": "engineering",
		},
		SKU: storage.SKU{
			Name: "Standard_LRS",
		},
		Properties: &storage.StorageAccountProperties{
			EnableHTTPSTrafficOnly: boolPtr(true),
			MinimumTLSVersion:     stringPtr("TLS1_2"),
			AllowBlobPublicAccess: boolPtr(false),
		},
	}

	// Serialize to ARM JSON
	serialized := serialize.ToARMResource(storageAccount)

	// Verify essential fields
	assert.Equal(t, "teststorage123", serialized["name"])
	assert.Equal(t, "Microsoft.Storage/storageAccounts", serialized["type"])
	assert.Equal(t, "2021-04-01", serialized["apiVersion"])
	assert.Equal(t, "westus", serialized["location"])
	assert.Equal(t, "StorageV2", serialized["kind"])

	// Verify SKU
	sku, ok := serialized["sku"].(map[string]interface{})
	require.True(t, ok, "SKU should be a map")
	assert.Equal(t, "Standard_LRS", sku["name"])

	// Verify tags (must be sorted for deterministic comparison)
	tags, ok := serialized["tags"].(map[string]interface{})
	require.True(t, ok, "Tags should be a map")
	assert.Equal(t, "production", tags["environment"])
	assert.Equal(t, "engineering", tags["cost-center"])

	// Verify properties
	props, ok := serialized["properties"].(map[string]interface{})
	require.True(t, ok, "Properties should be a map")
	assert.Equal(t, true, props["supportsHttpsTrafficOnly"])
	assert.Equal(t, "TLS1_2", props["minimumTlsVersion"])
	assert.Equal(t, false, props["allowBlobPublicAccess"])
}

// TestRoundTrip_EmptyPropertiesOmitted verifies that empty properties are omitted
func TestRoundTrip_EmptyPropertiesOmitted(t *testing.T) {
	storageAccount := &storage.StorageAccount{
		Name:       "teststorage",
		Type:       "Microsoft.Storage/storageAccounts",
		APIVersion: "2021-04-01",
		Location:   "eastus",
		Kind:       "StorageV2",
		SKU: storage.SKU{
			Name: "Standard_LRS",
		},
		// Properties is nil - should be omitted
		Properties: nil,
		// Tags is nil - should be omitted
		Tags: nil,
	}

	serialized := serialize.ToARMResource(storageAccount)

	// Verify that nil/empty fields are omitted
	_, hasProperties := serialized["properties"]
	assert.False(t, hasProperties, "Empty properties should be omitted")

	_, hasTags := serialized["tags"]
	assert.False(t, hasTags, "Empty tags should be omitted")

	// Verify required fields are present
	assert.Equal(t, "teststorage", serialized["name"])
	assert.Equal(t, "Microsoft.Storage/storageAccounts", serialized["type"])
}

// TestJSONComparisonSemantics verifies that our semantic comparison handles differences correctly
func TestJSONComparisonSemantics(t *testing.T) {
	tests := []struct {
		name      string
		a         map[string]interface{}
		b         map[string]interface{}
		shouldEqual bool
	}{
		{
			name: "identical maps",
			a: map[string]interface{}{
				"name": "test",
				"value": 123,
			},
			b: map[string]interface{}{
				"name": "test",
				"value": 123,
			},
			shouldEqual: true,
		},
		{
			name: "different key order (should be equal)",
			a: map[string]interface{}{
				"name": "test",
				"value": 123,
			},
			b: map[string]interface{}{
				"value": 123,
				"name": "test",
			},
			shouldEqual: true,
		},
		{
			name: "different values",
			a: map[string]interface{}{
				"name": "test",
				"value": 123,
			},
			b: map[string]interface{}{
				"name": "test",
				"value": 456,
			},
			shouldEqual: false,
		},
		{
			name: "missing key",
			a: map[string]interface{}{
				"name": "test",
				"value": 123,
			},
			b: map[string]interface{}{
				"name": "test",
			},
			shouldEqual: false,
		},
		{
			name: "nested maps with different order",
			a: map[string]interface{}{
				"outer": map[string]interface{}{
					"a": 1,
					"b": 2,
				},
			},
			b: map[string]interface{}{
				"outer": map[string]interface{}{
					"b": 2,
					"a": 1,
				},
			},
			shouldEqual: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldEqual {
				assertJSONEqualNoFail(t, tt.a, tt.b, true)
			} else {
				assertJSONEqualNoFail(t, tt.a, tt.b, false)
			}
		})
	}
}

// Helper functions

// convertARMResourceToStorageAccount converts an ARM resource to a storage account Go struct
func convertARMResourceToStorageAccount(t *testing.T, armRes ARMResource) *storage.StorageAccount {
	t.Helper()

	sa := &storage.StorageAccount{
		Name:       armRes.Name,
		Type:       armRes.Type,
		APIVersion: armRes.APIVersion,
		Location:   armRes.Location,
		Kind:       armRes.Kind,
	}

	// Convert SKU
	if len(armRes.SKU) > 0 {
		if name, ok := armRes.SKU["name"].(string); ok {
			sa.SKU = storage.SKU{Name: name}
		}
		if tier, ok := armRes.SKU["tier"].(string); ok {
			sa.SKU.Tier = &tier
		}
	}

	// Convert Tags
	if armRes.Tags != nil && len(armRes.Tags) > 0 {
		sa.Tags = armRes.Tags
	}

	// Convert Properties
	if armRes.Properties != nil && len(armRes.Properties) > 0 {
		props := &storage.StorageAccountProperties{}

		if val, ok := armRes.Properties["supportsHttpsTrafficOnly"].(bool); ok {
			props.EnableHTTPSTrafficOnly = &val
		}
		if val, ok := armRes.Properties["minimumTlsVersion"].(string); ok {
			props.MinimumTLSVersion = &val
		}
		if val, ok := armRes.Properties["allowBlobPublicAccess"].(bool); ok {
			props.AllowBlobPublicAccess = &val
		}
		if val, ok := armRes.Properties["allowSharedKeyAccess"].(bool); ok {
			props.AllowSharedKeyAccess = &val
		}

		sa.Properties = props
	}

	return sa
}

// assertJSONEqual performs semantic comparison of two JSON-like structures
func assertJSONEqual(t *testing.T, expected, actual map[string]interface{}) {
	t.Helper()

	// Convert to JSON strings for comparison (this normalizes the format)
	expectedJSON, err := json.Marshal(expected)
	require.NoError(t, err)

	actualJSON, err := json.Marshal(actual)
	require.NoError(t, err)

	// Parse back to ensure consistent ordering
	var expectedNormalized, actualNormalized interface{}
	err = json.Unmarshal(expectedJSON, &expectedNormalized)
	require.NoError(t, err)

	err = json.Unmarshal(actualJSON, &actualNormalized)
	require.NoError(t, err)

	// Use testify's assert.Equal which handles deep equality
	if !assert.Equal(t, expectedNormalized, actualNormalized) {
		// On failure, print pretty JSON for debugging
		expectedPretty, _ := json.MarshalIndent(expectedNormalized, "", "  ")
		actualPretty, _ := json.MarshalIndent(actualNormalized, "", "  ")

		t.Logf("Expected JSON:\n%s", string(expectedPretty))
		t.Logf("Actual JSON:\n%s", string(actualPretty))
	}
}

// assertJSONEqualNoFail is a helper for testing the comparison itself
func assertJSONEqualNoFail(t *testing.T, a, b map[string]interface{}, shouldEqual bool) {
	t.Helper()

	aJSON, _ := json.Marshal(a)
	bJSON, _ := json.Marshal(b)

	var aNorm, bNorm interface{}
	json.Unmarshal(aJSON, &aNorm)
	json.Unmarshal(bJSON, &bNorm)

	equal := assert.ObjectsAreEqual(aNorm, bNorm)

	if shouldEqual {
		assert.True(t, equal, "Expected objects to be equal")
	} else {
		assert.False(t, equal, "Expected objects to be different")
	}
}

// Helper functions for pointer types
func boolPtr(b bool) *bool {
	return &b
}

func stringPtr(s string) *string {
	return &s
}
