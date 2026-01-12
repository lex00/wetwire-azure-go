// Package storage provides Azure storage resource types
package storage

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStorageAccount(t *testing.T) {
	sa := NewStorageAccount("mystorageacct", "eastus", "StorageV2", "Standard_LRS")

	assert.Equal(t, "mystorageacct", sa.Name)
	assert.Equal(t, "Microsoft.Storage/storageAccounts", sa.Type)
	assert.Equal(t, "2021-04-01", sa.APIVersion)
	assert.Equal(t, "eastus", sa.Location)
	assert.Equal(t, "StorageV2", sa.Kind)
	assert.Equal(t, "Standard_LRS", sa.SKU.Name)
}

func TestStorageAccount_WithTags(t *testing.T) {
	sa := NewStorageAccount("mystorageacct", "eastus", "StorageV2", "Standard_LRS").
		WithTags(map[string]string{"env": "prod", "team": "platform"})

	assert.Equal(t, "prod", sa.Tags["env"])
	assert.Equal(t, "platform", sa.Tags["team"])
}

func TestStorageAccount_WithHTTPSOnly(t *testing.T) {
	sa := NewStorageAccount("mystorageacct", "eastus", "StorageV2", "Standard_LRS").
		WithHTTPSOnly(true)

	require.NotNil(t, sa.Properties)
	require.NotNil(t, sa.Properties.EnableHTTPSTrafficOnly)
	assert.True(t, *sa.Properties.EnableHTTPSTrafficOnly)
}

func TestStorageAccount_WithMinTLSVersion(t *testing.T) {
	sa := NewStorageAccount("mystorageacct", "eastus", "StorageV2", "Standard_LRS").
		WithMinTLSVersion("TLS1_2")

	require.NotNil(t, sa.Properties)
	require.NotNil(t, sa.Properties.MinimumTLSVersion)
	assert.Equal(t, "TLS1_2", *sa.Properties.MinimumTLSVersion)
}

func TestStorageAccount_ChainedBuilders(t *testing.T) {
	sa := NewStorageAccount("mystorageacct", "eastus", "StorageV2", "Standard_LRS").
		WithTags(map[string]string{"env": "prod"}).
		WithHTTPSOnly(true).
		WithMinTLSVersion("TLS1_2")

	assert.Equal(t, "mystorageacct", sa.Name)
	assert.Equal(t, "prod", sa.Tags["env"])
	require.NotNil(t, sa.Properties)
	assert.True(t, *sa.Properties.EnableHTTPSTrafficOnly)
	assert.Equal(t, "TLS1_2", *sa.Properties.MinimumTLSVersion)
}

func TestStorageAccount_JSON(t *testing.T) {
	sa := NewStorageAccount("mystorageacct", "eastus", "StorageV2", "Standard_LRS").
		WithTags(map[string]string{"env": "prod"})

	data, err := json.Marshal(sa)
	require.NoError(t, err)

	var result map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &result))

	assert.Equal(t, "mystorageacct", result["name"])
	assert.Equal(t, "Microsoft.Storage/storageAccounts", result["type"])
	assert.Equal(t, "2021-04-01", result["apiVersion"])
	assert.Equal(t, "eastus", result["location"])
	assert.Equal(t, "StorageV2", result["kind"])

	sku := result["sku"].(map[string]interface{})
	assert.Equal(t, "Standard_LRS", sku["name"])

	tags := result["tags"].(map[string]interface{})
	assert.Equal(t, "prod", tags["env"])
}

func TestStorageAccount_JSONWithProperties(t *testing.T) {
	sa := NewStorageAccount("mystorageacct", "eastus", "StorageV2", "Standard_LRS").
		WithHTTPSOnly(true).
		WithMinTLSVersion("TLS1_2")

	data, err := json.Marshal(sa)
	require.NoError(t, err)

	var result map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &result))

	props := result["properties"].(map[string]interface{})
	assert.Equal(t, true, props["supportsHttpsTrafficOnly"])
	assert.Equal(t, "TLS1_2", props["minimumTlsVersion"])
}

func TestSKU(t *testing.T) {
	tier := "Standard"
	sku := SKU{
		Name: "Standard_LRS",
		Tier: &tier,
	}

	data, err := json.Marshal(sku)
	require.NoError(t, err)

	var result map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &result))

	assert.Equal(t, "Standard_LRS", result["name"])
	assert.Equal(t, "Standard", result["tier"])
}

func TestNetworkRuleSet(t *testing.T) {
	bypass := "AzureServices"
	nrs := NetworkRuleSet{
		DefaultAction: "Deny",
		Bypass:        &bypass,
		IPRules: []IPRule{
			{Value: "10.0.0.0/8"},
		},
	}

	data, err := json.Marshal(nrs)
	require.NoError(t, err)

	var result map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &result))

	assert.Equal(t, "Deny", result["defaultAction"])
	assert.Equal(t, "AzureServices", result["bypass"])

	rules := result["ipRules"].([]interface{})
	require.Len(t, rules, 1)
	rule := rules[0].(map[string]interface{})
	assert.Equal(t, "10.0.0.0/8", rule["value"])
}

func TestEncryption(t *testing.T) {
	enc := Encryption{
		KeySource: "Microsoft.Storage",
		Services: &EncryptionServices{
			Blob: &EncryptionService{Enabled: true},
			File: &EncryptionService{Enabled: true},
		},
	}

	data, err := json.Marshal(enc)
	require.NoError(t, err)

	var result map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &result))

	assert.Equal(t, "Microsoft.Storage", result["keySource"])

	services := result["services"].(map[string]interface{})
	blob := services["blob"].(map[string]interface{})
	assert.Equal(t, true, blob["enabled"])
}

func TestIdentity(t *testing.T) {
	id := Identity{
		Type: "SystemAssigned",
	}

	data, err := json.Marshal(id)
	require.NoError(t, err)

	var result map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &result))

	assert.Equal(t, "SystemAssigned", result["type"])
}
