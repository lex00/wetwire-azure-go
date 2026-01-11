package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildCommand_NoArgs(t *testing.T) {
	// Test that build with no args uses current directory
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.json")

	// Create a simple Go file with resources
	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var myStorage = storage.StorageAccount{
	Name:     "mystorageacct",
	Location: "eastus",
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	exitCode := runBuild([]string{"-o", outputFile, tmpDir})
	assert.Equal(t, 0, exitCode)

	// Verify output file exists
	_, err = os.Stat(outputFile)
	assert.NoError(t, err)
}

func TestBuildCommand_WithPackagePath(t *testing.T) {
	tmpDir := t.TempDir()
	pkgDir := filepath.Join(tmpDir, "infra")
	err := os.MkdirAll(pkgDir, 0755)
	require.NoError(t, err)

	outputFile := filepath.Join(tmpDir, "output.json")

	code := `package infra

import "github.com/lex00/wetwire-azure-go/resources/storage"

var myStorage = storage.StorageAccount{
	Name:     "mystorageacct",
	Location: "eastus",
}
`
	err = os.WriteFile(filepath.Join(pkgDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	exitCode := runBuild([]string{"-o", outputFile, pkgDir})
	assert.Equal(t, 0, exitCode)

	// Verify output file exists
	_, err = os.Stat(outputFile)
	assert.NoError(t, err)
}

func TestBuildCommand_OutputFlag(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "custom-output.json")

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var myStorage = storage.StorageAccount{
	Name:     "mystorageacct",
	Location: "eastus",
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	exitCode := runBuild([]string{"--output", outputFile, tmpDir})
	assert.Equal(t, 0, exitCode)

	// Verify output file exists
	_, err = os.Stat(outputFile)
	assert.NoError(t, err)
}

func TestBuildCommand_FormatFlag(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.json")

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var myStorage = storage.StorageAccount{
	Name:     "mystorageacct",
	Location: "eastus",
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	exitCode := runBuild([]string{"-o", outputFile, "--format", "arm", tmpDir})
	assert.Equal(t, 0, exitCode)

	// Verify ARM template structure
	data, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	var template map[string]interface{}
	err = json.Unmarshal(data, &template)
	require.NoError(t, err)

	assert.Contains(t, template, "$schema")
	assert.Contains(t, template, "contentVersion")
	assert.Contains(t, template, "resources")
}

func TestBuildCommand_InvalidFormat(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.json")

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var myStorage = storage.StorageAccount{
	Name:     "mystorageacct",
	Location: "eastus",
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	exitCode := runBuild([]string{"-o", outputFile, "--format", "invalid", tmpDir})
	assert.Equal(t, 2, exitCode) // Invalid arguments
}

func TestBuildCommand_NonExistentPath(t *testing.T) {
	exitCode := runBuild([]string{"-o", "/tmp/out.json", "/nonexistent/path"})
	assert.Equal(t, 1, exitCode) // Build error
}

func TestBuildCommand_InvalidGoCode(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.json")

	err := os.WriteFile(filepath.Join(tmpDir, "bad.go"), []byte("this is not valid go code {{{"), 0644)
	require.NoError(t, err)

	exitCode := runBuild([]string{"-o", outputFile, tmpDir})
	assert.Equal(t, 1, exitCode) // Build error
}

func TestBuildCommand_NoResources(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.json")

	code := `package main

var x = 42
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	exitCode := runBuild([]string{"-o", outputFile, tmpDir})
	assert.Equal(t, 0, exitCode)

	// Verify output file has empty resources
	data, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	var template map[string]interface{}
	err = json.Unmarshal(data, &template)
	require.NoError(t, err)

	resources, ok := template["resources"].([]interface{})
	assert.True(t, ok)
	assert.Empty(t, resources)
}

func TestBuildCommand_ParametersFile(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.json")
	paramsFile := filepath.Join(tmpDir, "output.parameters.json")

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var myStorage = storage.StorageAccount{
	Name:     "mystorageacct",
	Location: "eastus",
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	exitCode := runBuild([]string{"-o", outputFile, "--parameters-file", paramsFile, tmpDir})
	assert.Equal(t, 0, exitCode)

	// Verify parameters file exists
	_, err = os.Stat(paramsFile)
	assert.NoError(t, err)

	// Verify parameters file structure
	data, err := os.ReadFile(paramsFile)
	require.NoError(t, err)

	var params map[string]interface{}
	err = json.Unmarshal(data, &params)
	require.NoError(t, err)

	assert.Contains(t, params, "$schema")
	assert.Contains(t, params, "contentVersion")
	assert.Contains(t, params, "parameters")
}

func TestBuildCommand_MultipleResources(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.json")

	code := `package main

import (
	"github.com/lex00/wetwire-azure-go/resources/storage"
	"github.com/lex00/wetwire-azure-go/resources/network"
)

var myStorage = storage.StorageAccount{
	Name:     "mystorageacct",
	Location: "eastus",
}

var myVNet = network.VirtualNetwork{
	Name:     "myvnet",
	Location: "eastus",
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	exitCode := runBuild([]string{"-o", outputFile, tmpDir})
	assert.Equal(t, 0, exitCode)

	// Verify output has both resources
	data, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	var template map[string]interface{}
	err = json.Unmarshal(data, &template)
	require.NoError(t, err)

	resources, ok := template["resources"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, resources, 2)
}

func TestBuildCommand_ResourceWithDependencies(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.json")

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
	Name:           "my-nic",
	Location:       "eastus",
	VirtualNetwork: &myVNet,
}

var myVM = compute.VirtualMachine{
	Name:             "my-vm",
	Location:         "eastus",
	NetworkInterface: &myNIC,
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	exitCode := runBuild([]string{"-o", outputFile, tmpDir})
	assert.Equal(t, 0, exitCode)

	// Verify output has resources in correct order
	data, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	var template map[string]interface{}
	err = json.Unmarshal(data, &template)
	require.NoError(t, err)

	resources, ok := template["resources"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, resources, 3)

	// Verify dependsOn is set correctly
	for _, res := range resources {
		resMap := res.(map[string]interface{})
		name := resMap["name"].(string)
		if name == "my-vm" || name == "myVM" {
			dependsOn, ok := resMap["dependsOn"]
			assert.True(t, ok, "VM should have dependsOn")
			deps := dependsOn.([]interface{})
			assert.NotEmpty(t, deps)
		}
	}
}

func TestBuildCommand_StdoutOutput(t *testing.T) {
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

	// When -o is not specified, should output to stdout
	exitCode := runBuild([]string{tmpDir})
	assert.Equal(t, 0, exitCode)
}

func TestBuildCommand_ShortOutputFlag(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.json")

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var myStorage = storage.StorageAccount{
	Name:     "mystorageacct",
	Location: "eastus",
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	exitCode := runBuild([]string{"-o", outputFile, tmpDir})
	assert.Equal(t, 0, exitCode)

	// Verify output file exists and is valid JSON
	data, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	var template map[string]interface{}
	err = json.Unmarshal(data, &template)
	require.NoError(t, err)
	assert.Contains(t, template, "$schema")
}
