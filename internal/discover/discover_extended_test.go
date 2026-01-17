package discover

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDiscoverResources_BlankIdentifier tests that _ variables are skipped
func TestDiscoverResources_BlankIdentifier(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var _ = storage.StorageAccount{
	Name:     "skippedStorage",
	Location: "eastus",
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	resources, err := DiscoverResources(tmpDir)
	require.NoError(t, err)
	assert.Empty(t, resources, "Blank identifier should be skipped")
}

// TestDiscoverResources_ExplicitType tests resources with explicit type annotation
func TestDiscoverResources_ExplicitType(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var myStorage storage.StorageAccount = storage.StorageAccount{
	Name:     "explicitStorage",
	Location: "eastus",
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	resources, err := DiscoverResources(tmpDir)
	require.NoError(t, err)
	require.Len(t, resources, 1)
	assert.Equal(t, "myStorage", resources[0].Name)
}

// TestDiscoverResources_KeyVault tests keyvault discovery
func TestDiscoverResources_KeyVault(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/keyvault"

var myVault = keyvault.Vault{
	Name:     "mykeyvault",
	Location: "eastus",
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	resources, err := DiscoverResources(tmpDir)
	require.NoError(t, err)
	require.Len(t, resources, 1)
	assert.Equal(t, "Microsoft.KeyVault/vaults", resources[0].Type)
}

// TestDiscoverResources_SQLResources tests SQL server and database discovery
func TestDiscoverResources_SQLResources(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/sql"

var mySqlServer = sql.Server{
	Name:     "mysqlserver",
	Location: "eastus",
}

var mySqlDb = sql.Database{
	Name:     "mysqldb",
	Location: "eastus",
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	resources, err := DiscoverResources(tmpDir)
	require.NoError(t, err)
	require.Len(t, resources, 2)

	resourceMap := make(map[string]DiscoveredResource)
	for _, r := range resources {
		resourceMap[r.Name] = r
	}

	assert.Equal(t, "Microsoft.Sql/servers", resourceMap["mySqlServer"].Type)
	assert.Equal(t, "Microsoft.Sql/servers/databases", resourceMap["mySqlDb"].Type)
}

// TestDiscoverResources_WebApp tests web site discovery
func TestDiscoverResources_WebApp(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/web"

var myWebApp = web.Site{
	Name:     "mywebapp",
	Location: "eastus",
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	resources, err := DiscoverResources(tmpDir)
	require.NoError(t, err)
	require.Len(t, resources, 1)
	assert.Equal(t, "Microsoft.Web/sites", resources[0].Type)
}

// TestDiscoverResources_ContainerRegistry tests container registry discovery
func TestDiscoverResources_ContainerRegistry(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/containerregistry"

var myRegistry = containerregistry.Registry{
	Name:     "myregistry",
	Location: "eastus",
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	resources, err := DiscoverResources(tmpDir)
	require.NoError(t, err)
	require.Len(t, resources, 1)
	assert.Equal(t, "Microsoft.ContainerRegistry/registries", resources[0].Type)
}

// TestDiscoverResources_AKS tests managed cluster discovery
func TestDiscoverResources_AKS(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/aks"

var myCluster = aks.ManagedCluster{
	Name:     "myakscluster",
	Location: "eastus",
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	resources, err := DiscoverResources(tmpDir)
	require.NoError(t, err)
	require.Len(t, resources, 1)
	assert.Equal(t, "Microsoft.ContainerService/managedClusters", resources[0].Type)
}

// TestDiscoverResources_AllNetworkTypes tests all network resource types
func TestDiscoverResources_AllNetworkTypes(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/network"

var mySubnet = network.Subnet{
	Name: "mysubnet",
}

var myPublicIP = network.PublicIPAddress{
	Name:     "mypublicip",
	Location: "eastus",
}

var myNSG = network.NetworkSecurityGroup{
	Name:     "mynsg",
	Location: "eastus",
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	resources, err := DiscoverResources(tmpDir)
	require.NoError(t, err)
	require.Len(t, resources, 3)

	resourceMap := make(map[string]string)
	for _, r := range resources {
		resourceMap[r.Name] = r.Type
	}

	assert.Equal(t, "Microsoft.Network/subnets", resourceMap["mySubnet"])
	assert.Equal(t, "Microsoft.Network/publicIPAddresses", resourceMap["myPublicIP"])
	assert.Equal(t, "Microsoft.Network/networkSecurityGroups", resourceMap["myNSG"])
}

// TestDiscoverResources_NonWetwireImport tests that non-wetwire imports are ignored
func TestDiscoverResources_NonWetwireImport(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import "some/other/package/storage"

var myStorage = storage.StorageAccount{
	Name:     "mystorageaccount",
	Location: "eastus",
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	resources, err := DiscoverResources(tmpDir)
	require.NoError(t, err)
	assert.Empty(t, resources, "Non-wetwire imports should be ignored")
}

// TestDiscoverResources_AliasedImport tests imports with aliases
func TestDiscoverResources_AliasedImport(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import storage "github.com/lex00/wetwire-azure-go/resources/storage"

var myStorage = storage.StorageAccount{
	Name:     "aliasedStorage",
	Location: "eastus",
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	resources, err := DiscoverResources(tmpDir)
	require.NoError(t, err)
	require.Len(t, resources, 1)
	assert.Equal(t, "myStorage", resources[0].Name)
}

// TestDiscoverResources_CallExprDependency tests dependency extraction from function calls
func TestDiscoverResources_CallExprDependency(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var baseName = "test"

func getName(base string) string { return base + "suffix" }

var myStorage = storage.StorageAccount{
	Name:     getName(baseName),
	Location: "eastus",
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	resources, err := DiscoverResources(tmpDir)
	require.NoError(t, err)
	require.Len(t, resources, 1)
	assert.Contains(t, resources[0].Dependencies, "baseName")
}

// TestDiscoverResources_MapTypeDependency tests dependency extraction from map types
func TestDiscoverResources_MapTypeDependency(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var envValue = "prod"
var costCenter = "it"

var myStorage = storage.StorageAccount{
	Name:     "mystorage",
	Location: "eastus",
	Tags: map[string]string{
		"Environment": envValue,
		"CostCenter":  costCenter,
	},
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	resources, err := DiscoverResources(tmpDir)
	require.NoError(t, err)
	require.Len(t, resources, 1)
	assert.Contains(t, resources[0].Dependencies, "envValue")
	assert.Contains(t, resources[0].Dependencies, "costCenter")
}

// TestDiscoverResources_SliceExprDependency tests dependency extraction from slice expressions
func TestDiscoverResources_SliceExprDependency(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/network"

var subnets = []network.Subnet{
	{Name: "subnet1"},
	{Name: "subnet2"},
}

var myVNet = network.VirtualNetwork{
	Name:    "myvnet",
	Subnets: subnets[0:1],
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	resources, err := DiscoverResources(tmpDir)
	require.NoError(t, err)

	var vnet DiscoveredResource
	for _, r := range resources {
		if r.Name == "myVNet" {
			vnet = r
			break
		}
	}
	assert.Contains(t, vnet.Dependencies, "subnets")
}

// TestDiscoverResources_IndexExprDependency tests dependency extraction from index expressions
func TestDiscoverResources_IndexExprDependency(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/network"

var subnetConfigs = []string{"config1", "config2"}

var myVNet = network.VirtualNetwork{
	Name:    subnetConfigs[0],
	Location: "eastus",
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	resources, err := DiscoverResources(tmpDir)
	require.NoError(t, err)

	require.Len(t, resources, 1)
	assert.Contains(t, resources[0].Dependencies, "subnetConfigs")
}

// TestDiscoverResources_ParenExprDependency tests dependency extraction from parenthesized expressions
func TestDiscoverResources_ParenExprDependency(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var baseName = "test"
var suffix = "storage"

var myStorage = storage.StorageAccount{
	Name:     (baseName + suffix),
	Location: "eastus",
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	resources, err := DiscoverResources(tmpDir)
	require.NoError(t, err)
	require.Len(t, resources, 1)
	assert.Contains(t, resources[0].Dependencies, "baseName")
	assert.Contains(t, resources[0].Dependencies, "suffix")
}

// TestDiscoverResources_BinaryExprDependency tests dependency extraction from binary expressions
func TestDiscoverResources_BinaryExprDependency(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var prefix = "prod"
var suffix = "storage"

var myStorage = storage.StorageAccount{
	Name:     prefix + suffix,
	Location: "eastus",
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	resources, err := DiscoverResources(tmpDir)
	require.NoError(t, err)
	require.Len(t, resources, 1)
	assert.Contains(t, resources[0].Dependencies, "prefix")
	assert.Contains(t, resources[0].Dependencies, "suffix")
}

// TestDiscoverResources_StarExprDependency tests dependency extraction from pointer dereference
func TestDiscoverResources_StarExprDependency(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/network"

var subnetRef = &network.Subnet{Name: "subnet1"}

var myVNet = network.VirtualNetwork{
	Name:    (*subnetRef).Name,
	Location: "eastus",
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	resources, err := DiscoverResources(tmpDir)
	require.NoError(t, err)

	var vnet DiscoveredResource
	for _, r := range resources {
		if r.Name == "myVNet" {
			vnet = r
			break
		}
	}
	assert.Contains(t, vnet.Dependencies, "subnetRef")
}

// TestDiscoverResources_SelectorExprDependency tests dependency extraction from selectors
func TestDiscoverResources_SelectorExprDependency(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"
import "github.com/lex00/wetwire-azure-go/resources/network"

var myVNet = network.VirtualNetwork{
	Name:     "myvnet",
	Location: "eastus",
}

var myStorage = storage.StorageAccount{
	Name:     myVNet.Name,
	Location: "eastus",
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	resources, err := DiscoverResources(tmpDir)
	require.NoError(t, err)
	require.Len(t, resources, 2)

	var storageRes DiscoveredResource
	for _, r := range resources {
		if r.Name == "myStorage" {
			storageRes = r
			break
		}
	}
	assert.Contains(t, storageRes.Dependencies, "myVNet")
}

// TestDiscoverResources_NonexistentDirectory tests error handling for non-existent directories
func TestDiscoverResources_NonexistentDirectory(t *testing.T) {
	_, err := DiscoverResources("/nonexistent/path/that/does/not/exist")
	assert.Error(t, err)
}

// TestDiscoverResources_NoImports tests file with no import declarations
func TestDiscoverResources_NoImports(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

type LocalStruct struct {
	Name string
}

var myVar = LocalStruct{Name: "test"}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	resources, err := DiscoverResources(tmpDir)
	require.NoError(t, err)
	assert.Empty(t, resources)
}

// TestDiscoverResources_MultipleVarsInOneDecl tests multiple variables in one var declaration
func TestDiscoverResources_MultipleVarsInOneDecl(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var (
	storage1 = storage.StorageAccount{
		Name:     "storage1",
		Location: "eastus",
	}
	storage2 = storage.StorageAccount{
		Name:     "storage2",
		Location: "westus",
	}
)
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	resources, err := DiscoverResources(tmpDir)
	require.NoError(t, err)
	require.Len(t, resources, 2)
}

// TestExtractImports_NilPath tests extractImports edge case
func TestExtractImports_NilPath(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a valid Go file with imports
	code := `package main

import "fmt"

func main() {
	fmt.Println("hello")
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	resources, err := DiscoverResources(tmpDir)
	require.NoError(t, err)
	assert.Empty(t, resources)
}

// TestDiscoverResources_ArrayTypeDependency tests dependency extraction from array types
func TestDiscoverResources_ArrayTypeDependency(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/network"

type NetworkConfig struct {
	Subnets [5]network.Subnet
}

var config = NetworkConfig{}

var myVNet = network.VirtualNetwork{
	Name:     "myvnet",
	Location: "eastus",
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	resources, err := DiscoverResources(tmpDir)
	require.NoError(t, err)
	require.Len(t, resources, 1) // Only VirtualNetwork should be discovered
}

// TestIsBuiltinIdent tests the isBuiltinIdent helper function
func TestIsBuiltinIdent(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"bool", true},
		{"byte", true},
		{"complex64", true},
		{"complex128", true},
		{"error", true},
		{"float32", true},
		{"float64", true},
		{"int", true},
		{"int8", true},
		{"int16", true},
		{"int32", true},
		{"int64", true},
		{"rune", true},
		{"string", true},
		{"uint", true},
		{"uint8", true},
		{"uint16", true},
		{"uint32", true},
		{"uint64", true},
		{"uintptr", true},
		{"true", true},
		{"false", true},
		{"nil", true},
		// Also test builtin functions (new in coreast.IsBuiltinIdent)
		{"len", true},
		{"cap", true},
		{"make", true},
		{"new", true},
		{"append", true},
		// Non-builtins
		{"MyStruct", false},
		{"storageAccount", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isBuiltinIdent(tt.name)
			assert.Equal(t, tt.expected, result)
		})
	}
}
