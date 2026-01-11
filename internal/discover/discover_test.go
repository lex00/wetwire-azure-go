package discover

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiscoverResources_EmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	resources, err := DiscoverResources(tmpDir)
	require.NoError(t, err)
	assert.Empty(t, resources)
}

func TestDiscoverResources_NonGoFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a non-Go file
	err := os.WriteFile(filepath.Join(tmpDir, "readme.txt"), []byte("test"), 0644)
	require.NoError(t, err)

	resources, err := DiscoverResources(tmpDir)
	require.NoError(t, err)
	assert.Empty(t, resources)
}

func TestDiscoverResources_SimpleStorageAccount(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var myStorage = storage.StorageAccount{
	Name:     "mystorageacct",
	Location: "eastus",
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	resources, err := DiscoverResources(tmpDir)
	require.NoError(t, err)
	require.Len(t, resources, 1)

	res := resources[0]
	assert.Equal(t, "myStorage", res.Name)
	assert.Equal(t, "Microsoft.Storage/storageAccounts", res.Type)
	assert.Contains(t, res.File, "main.go")
	assert.Equal(t, 5, res.Line)
	assert.Empty(t, res.Dependencies)
}

func TestDiscoverResources_VirtualMachineWithDependencies(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import (
	"github.com/lex00/wetwire-azure-go/resources/compute"
	"github.com/lex00/wetwire-azure-go/resources/network"
)

var myVNet = network.VirtualNetwork{
	Name:     "my-vnet",
	Location: "eastus",
}

var myNIC = network.NetworkInterface{
	Name:     "my-nic",
	Location: "eastus",
	VirtualNetwork: &myVNet,
}

var myVM = compute.VirtualMachine{
	Name:     "my-vm",
	Location: "eastus",
	NetworkInterface: &myNIC,
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	resources, err := DiscoverResources(tmpDir)
	require.NoError(t, err)
	require.Len(t, resources, 3)

	// Sort for predictable testing
	resourceMap := make(map[string]DiscoveredResource)
	for _, r := range resources {
		resourceMap[r.Name] = r
	}

	// Check VNet (no dependencies)
	vnet := resourceMap["myVNet"]
	assert.Equal(t, "Microsoft.Network/virtualNetworks", vnet.Type)
	assert.Equal(t, 8, vnet.Line)
	assert.Empty(t, vnet.Dependencies)

	// Check NIC (depends on VNet)
	nic := resourceMap["myNIC"]
	assert.Equal(t, "Microsoft.Network/networkInterfaces", nic.Type)
	assert.Equal(t, 13, nic.Line)
	assert.Contains(t, nic.Dependencies, "myVNet")

	// Check VM (depends on NIC)
	vm := resourceMap["myVM"]
	assert.Equal(t, "Microsoft.Compute/virtualMachines", vm.Type)
	assert.Equal(t, 19, vm.Line)
	assert.Contains(t, vm.Dependencies, "myNIC")
}

func TestDiscoverResources_MultipleFiles(t *testing.T) {
	tmpDir := t.TempDir()

	storage := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var storage1 = storage.StorageAccount{
	Name:     "storage1",
	Location: "eastus",
}
`
	network := `package main

import "github.com/lex00/wetwire-azure-go/resources/network"

var vnet1 = network.VirtualNetwork{
	Name:     "vnet1",
	Location: "westus",
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "storage.go"), []byte(storage), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(tmpDir, "network.go"), []byte(network), 0644)
	require.NoError(t, err)

	resources, err := DiscoverResources(tmpDir)
	require.NoError(t, err)
	assert.Len(t, resources, 2)

	names := make([]string, len(resources))
	for i, r := range resources {
		names[i] = r.Name
	}
	assert.Contains(t, names, "storage1")
	assert.Contains(t, names, "vnet1")
}

func TestDiscoverResources_NestedDirectories(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "infra")
	err := os.MkdirAll(subDir, 0755)
	require.NoError(t, err)

	code := `package infra

import "github.com/lex00/wetwire-azure-go/resources/storage"

var nestedStorage = storage.StorageAccount{
	Name:     "nested",
	Location: "eastus",
}
`
	err = os.WriteFile(filepath.Join(subDir, "storage.go"), []byte(code), 0644)
	require.NoError(t, err)

	resources, err := DiscoverResources(tmpDir)
	require.NoError(t, err)
	require.Len(t, resources, 1)
	assert.Equal(t, "nestedStorage", resources[0].Name)
}

func TestDiscoverResources_InvalidGoFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create invalid Go code
	err := os.WriteFile(filepath.Join(tmpDir, "bad.go"), []byte("this is not valid go code {{{"), 0644)
	require.NoError(t, err)

	// Should return error for unparseable code
	_, err = DiscoverResources(tmpDir)
	assert.Error(t, err)
}

func TestDiscoverResources_NonAzureResources(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

type CustomStruct struct {
	Name string
}

var myVar = CustomStruct{
	Name: "test",
}

var myString = "hello"
var myInt = 42
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	resources, err := DiscoverResources(tmpDir)
	require.NoError(t, err)
	assert.Empty(t, resources, "Should not discover non-Azure resource types")
}

func TestDiscoverResources_FunctionDeclarations(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

func createStorage() storage.StorageAccount {
	return storage.StorageAccount{
		Name:     "funcStorage",
		Location: "eastus",
	}
}

var myStorage = storage.StorageAccount{
	Name:     "varStorage",
	Location: "eastus",
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	resources, err := DiscoverResources(tmpDir)
	require.NoError(t, err)
	require.Len(t, resources, 1, "Should only discover top-level var declarations")
	assert.Equal(t, "myStorage", resources[0].Name)
}

func TestDiscoverResources_PointerReferences(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import (
	"github.com/lex00/wetwire-azure-go/resources/network"
	"github.com/lex00/wetwire-azure-go/resources/compute"
)

var subnet1 = network.Subnet{
	Name: "subnet1",
}

var subnet2 = network.Subnet{
	Name: "subnet2",
}

var vnet = network.VirtualNetwork{
	Name:    "vnet",
	Subnets: []*network.Subnet{&subnet1, &subnet2},
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	resources, err := DiscoverResources(tmpDir)
	require.NoError(t, err)
	require.Len(t, resources, 3)

	// Find vnet resource
	var vnet DiscoveredResource
	for _, r := range resources {
		if r.Name == "vnet" {
			vnet = r
			break
		}
	}

	assert.Contains(t, vnet.Dependencies, "subnet1")
	assert.Contains(t, vnet.Dependencies, "subnet2")
}

func TestDiscoverResources_ChainedDependencies(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/network"

var subnet = network.Subnet{
	Name: "subnet",
}

var vnet = network.VirtualNetwork{
	Name:   "vnet",
	Subnet: &subnet,
}

var nic = network.NetworkInterface{
	Name:           "nic",
	VirtualNetwork: &vnet,
	Subnet:         &subnet,
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	resources, err := DiscoverResources(tmpDir)
	require.NoError(t, err)
	require.Len(t, resources, 3)

	var nic DiscoveredResource
	for _, r := range resources {
		if r.Name == "nic" {
			nic = r
			break
		}
	}

	// NIC should depend on both vnet and subnet
	assert.Contains(t, nic.Dependencies, "vnet")
	assert.Contains(t, nic.Dependencies, "subnet")
}
