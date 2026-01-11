package template

import (
	"encoding/json"
	"testing"

	"github.com/lex00/wetwire-azure-go/internal/discover"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTemplateBuilder(t *testing.T) {
	builder := NewTemplateBuilder()

	assert.NotNil(t, builder)
	assert.NotNil(t, builder.resources)
	assert.NotNil(t, builder.parameters)
	assert.NotNil(t, builder.variables)
	assert.NotNil(t, builder.outputs)
}

func TestAddResource(t *testing.T) {
	tests := []struct {
		name     string
		resource discover.DiscoveredResource
		wantErr  bool
	}{
		{
			name: "valid resource",
			resource: discover.DiscoveredResource{
				Name: "myStorage",
				Type: "Microsoft.Storage/storageAccounts",
				File: "/path/to/file.go",
				Line: 10,
			},
			wantErr: false,
		},
		{
			name: "duplicate resource name",
			resource: discover.DiscoveredResource{
				Name: "myStorage",
				Type: "Microsoft.Storage/storageAccounts",
				File: "/path/to/file2.go",
				Line: 20,
			},
			wantErr: true,
		},
	}

	builder := NewTemplateBuilder()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := builder.AddResource(tt.resource)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAddParameter(t *testing.T) {
	tests := []struct {
		name      string
		paramName string
		paramType string
		wantErr   bool
	}{
		{
			name:      "valid parameter",
			paramName: "location",
			paramType: "string",
			wantErr:   false,
		},
		{
			name:      "duplicate parameter",
			paramName: "location",
			paramType: "string",
			wantErr:   true,
		},
	}

	builder := NewTemplateBuilder()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := builder.AddParameter(tt.paramName, tt.paramType, nil)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAddVariable(t *testing.T) {
	tests := []struct {
		name     string
		varName  string
		varValue interface{}
		wantErr  bool
	}{
		{
			name:     "valid variable",
			varName:  "storageAccountName",
			varValue: "[concat('storage', uniqueString(resourceGroup().id))]",
			wantErr:  false,
		},
		{
			name:     "duplicate variable",
			varName:  "storageAccountName",
			varValue: "something else",
			wantErr:  true,
		},
	}

	builder := NewTemplateBuilder()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := builder.AddVariable(tt.varName, tt.varValue)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAddOutput(t *testing.T) {
	tests := []struct {
		name       string
		outputName string
		outputType string
		value      interface{}
		wantErr    bool
	}{
		{
			name:       "valid output",
			outputName: "storageAccountId",
			outputType: "string",
			value:      "[resourceId('Microsoft.Storage/storageAccounts', 'myStorage')]",
			wantErr:    false,
		},
		{
			name:       "duplicate output",
			outputName: "storageAccountId",
			outputType: "string",
			value:      "something",
			wantErr:    true,
		},
	}

	builder := NewTemplateBuilder()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := builder.AddOutput(tt.outputName, tt.outputType, tt.value)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBuild_EmptyTemplate(t *testing.T) {
	builder := NewTemplateBuilder()

	result, err := builder.Build()
	require.NoError(t, err)

	// Parse the JSON to verify structure
	var template map[string]interface{}
	err = json.Unmarshal([]byte(result), &template)
	require.NoError(t, err)

	// Verify required fields
	assert.Equal(t, "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#", template["$schema"])
	assert.Equal(t, "1.0.0.0", template["contentVersion"])
	assert.NotNil(t, template["resources"])
	assert.NotNil(t, template["parameters"])
	assert.NotNil(t, template["variables"])
	assert.NotNil(t, template["outputs"])
}

func TestBuild_WithResources(t *testing.T) {
	builder := NewTemplateBuilder()

	// Add resources
	err := builder.AddResource(discover.DiscoveredResource{
		Name: "myStorage",
		Type: "Microsoft.Storage/storageAccounts",
		File: "/path/to/file.go",
		Line: 10,
	})
	require.NoError(t, err)

	result, err := builder.Build()
	require.NoError(t, err)

	// Parse the JSON
	var template map[string]interface{}
	err = json.Unmarshal([]byte(result), &template)
	require.NoError(t, err)

	// Verify resources
	resources := template["resources"].([]interface{})
	assert.Len(t, resources, 1)

	resource := resources[0].(map[string]interface{})
	assert.Equal(t, "myStorage", resource["name"])
	assert.Equal(t, "Microsoft.Storage/storageAccounts", resource["type"])
}

func TestBuild_WithDependencies(t *testing.T) {
	builder := NewTemplateBuilder()

	// Add storage account (no dependencies)
	err := builder.AddResource(discover.DiscoveredResource{
		Name: "myStorage",
		Type: "Microsoft.Storage/storageAccounts",
		File: "/path/to/file.go",
		Line: 10,
	})
	require.NoError(t, err)

	// Add VM that depends on storage
	err = builder.AddResource(discover.DiscoveredResource{
		Name:         "myVM",
		Type:         "Microsoft.Compute/virtualMachines",
		File:         "/path/to/file.go",
		Line:         20,
		Dependencies: []string{"myStorage"},
	})
	require.NoError(t, err)

	result, err := builder.Build()
	require.NoError(t, err)

	// Parse the JSON
	var template map[string]interface{}
	err = json.Unmarshal([]byte(result), &template)
	require.NoError(t, err)

	// Verify resources are ordered correctly (storage before VM)
	resources := template["resources"].([]interface{})
	assert.Len(t, resources, 2)

	// First resource should be storage (no dependencies)
	firstResource := resources[0].(map[string]interface{})
	assert.Equal(t, "myStorage", firstResource["name"])

	// Second resource should be VM (depends on storage)
	secondResource := resources[1].(map[string]interface{})
	assert.Equal(t, "myVM", secondResource["name"])

	// Verify dependsOn is present
	dependsOn := secondResource["dependsOn"].([]interface{})
	assert.Len(t, dependsOn, 1)
	assert.Contains(t, dependsOn[0].(string), "myStorage")
}

func TestBuild_CyclicDependency(t *testing.T) {
	builder := NewTemplateBuilder()

	// Add resource A that depends on B
	err := builder.AddResource(discover.DiscoveredResource{
		Name:         "resourceA",
		Type:         "Microsoft.Storage/storageAccounts",
		File:         "/path/to/file.go",
		Line:         10,
		Dependencies: []string{"resourceB"},
	})
	require.NoError(t, err)

	// Add resource B that depends on A (creating a cycle)
	err = builder.AddResource(discover.DiscoveredResource{
		Name:         "resourceB",
		Type:         "Microsoft.Compute/virtualMachines",
		File:         "/path/to/file.go",
		Line:         20,
		Dependencies: []string{"resourceA"},
	})
	require.NoError(t, err)

	// Build should fail due to cyclic dependency
	_, err = builder.Build()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cyclic dependency")
}

func TestBuild_MissingDependency(t *testing.T) {
	builder := NewTemplateBuilder()

	// Add resource that depends on non-existent resource
	err := builder.AddResource(discover.DiscoveredResource{
		Name:         "myVM",
		Type:         "Microsoft.Compute/virtualMachines",
		File:         "/path/to/file.go",
		Line:         20,
		Dependencies: []string{"nonExistentStorage"},
	})
	require.NoError(t, err)

	// Build should fail due to missing dependency
	_, err = builder.Build()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nonExistentStorage")
}

func TestBuild_ComplexDependencyGraph(t *testing.T) {
	builder := NewTemplateBuilder()

	// Create a complex dependency graph:
	// A (no deps)
	// B depends on A
	// C depends on A
	// D depends on B and C

	resources := []discover.DiscoveredResource{
		{
			Name: "resourceA",
			Type: "Microsoft.Storage/storageAccounts",
			File: "/path/to/file.go",
			Line: 10,
		},
		{
			Name:         "resourceB",
			Type:         "Microsoft.Compute/virtualMachines",
			File:         "/path/to/file.go",
			Line:         20,
			Dependencies: []string{"resourceA"},
		},
		{
			Name:         "resourceC",
			Type:         "Microsoft.Network/virtualNetworks",
			File:         "/path/to/file.go",
			Line:         30,
			Dependencies: []string{"resourceA"},
		},
		{
			Name:         "resourceD",
			Type:         "Microsoft.Network/networkInterfaces",
			File:         "/path/to/file.go",
			Line:         40,
			Dependencies: []string{"resourceB", "resourceC"},
		},
	}

	for _, res := range resources {
		err := builder.AddResource(res)
		require.NoError(t, err)
	}

	result, err := builder.Build()
	require.NoError(t, err)

	// Parse the JSON
	var template map[string]interface{}
	err = json.Unmarshal([]byte(result), &template)
	require.NoError(t, err)

	// Verify all resources are present
	templateResources := template["resources"].([]interface{})
	assert.Len(t, templateResources, 4)

	// Build a map of resource positions
	positions := make(map[string]int)
	for i, res := range templateResources {
		resMap := res.(map[string]interface{})
		positions[resMap["name"].(string)] = i
	}

	// Verify topological ordering
	// A must come before B and C
	assert.Less(t, positions["resourceA"], positions["resourceB"])
	assert.Less(t, positions["resourceA"], positions["resourceC"])

	// B and C must come before D
	assert.Less(t, positions["resourceB"], positions["resourceD"])
	assert.Less(t, positions["resourceC"], positions["resourceD"])
}

func TestBuild_WithParametersVariablesOutputs(t *testing.T) {
	builder := NewTemplateBuilder()

	// Add parameter
	err := builder.AddParameter("location", "string", map[string]interface{}{
		"defaultValue": "eastus",
	})
	require.NoError(t, err)

	// Add variable
	err = builder.AddVariable("storageAccountName", "[concat('storage', uniqueString(resourceGroup().id))]")
	require.NoError(t, err)

	// Add resource
	err = builder.AddResource(discover.DiscoveredResource{
		Name: "myStorage",
		Type: "Microsoft.Storage/storageAccounts",
		File: "/path/to/file.go",
		Line: 10,
	})
	require.NoError(t, err)

	// Add output
	err = builder.AddOutput("storageId", "string", "[resourceId('Microsoft.Storage/storageAccounts', 'myStorage')]")
	require.NoError(t, err)

	result, err := builder.Build()
	require.NoError(t, err)

	// Parse the JSON
	var template map[string]interface{}
	err = json.Unmarshal([]byte(result), &template)
	require.NoError(t, err)

	// Verify parameters
	parameters := template["parameters"].(map[string]interface{})
	assert.Contains(t, parameters, "location")
	locationParam := parameters["location"].(map[string]interface{})
	assert.Equal(t, "string", locationParam["type"])
	assert.Equal(t, "eastus", locationParam["defaultValue"])

	// Verify variables
	variables := template["variables"].(map[string]interface{})
	assert.Contains(t, variables, "storageAccountName")
	assert.Equal(t, "[concat('storage', uniqueString(resourceGroup().id))]", variables["storageAccountName"])

	// Verify outputs
	outputs := template["outputs"].(map[string]interface{})
	assert.Contains(t, outputs, "storageId")
	storageIdOutput := outputs["storageId"].(map[string]interface{})
	assert.Equal(t, "string", storageIdOutput["type"])
	assert.Equal(t, "[resourceId('Microsoft.Storage/storageAccounts', 'myStorage')]", storageIdOutput["value"])
}

func TestValidateReferences(t *testing.T) {
	tests := []struct {
		name      string
		resources []discover.DiscoveredResource
		wantErr   bool
		errMsg    string
	}{
		{
			name: "valid references",
			resources: []discover.DiscoveredResource{
				{Name: "resourceA", Type: "Microsoft.Storage/storageAccounts"},
				{Name: "resourceB", Type: "Microsoft.Compute/virtualMachines", Dependencies: []string{"resourceA"}},
			},
			wantErr: false,
		},
		{
			name: "missing reference",
			resources: []discover.DiscoveredResource{
				{Name: "resourceA", Type: "Microsoft.Storage/storageAccounts"},
				{Name: "resourceB", Type: "Microsoft.Compute/virtualMachines", Dependencies: []string{"nonExistent"}},
			},
			wantErr: true,
			errMsg:  "nonExistent",
		},
		{
			name: "cyclic dependency",
			resources: []discover.DiscoveredResource{
				{Name: "resourceA", Type: "Microsoft.Storage/storageAccounts", Dependencies: []string{"resourceB"}},
				{Name: "resourceB", Type: "Microsoft.Compute/virtualMachines", Dependencies: []string{"resourceA"}},
			},
			wantErr: true,
			errMsg:  "cyclic",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewTemplateBuilder()
			for _, res := range tt.resources {
				err := builder.AddResource(res)
				require.NoError(t, err)
			}

			err := builder.validateReferences()
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTopologicalSort(t *testing.T) {
	tests := []struct {
		name      string
		resources []discover.DiscoveredResource
		wantOrder []string
		wantErr   bool
	}{
		{
			name: "simple linear dependency",
			resources: []discover.DiscoveredResource{
				{Name: "C", Dependencies: []string{"B"}},
				{Name: "B", Dependencies: []string{"A"}},
				{Name: "A"},
			},
			wantOrder: []string{"A", "B", "C"},
			wantErr:   false,
		},
		{
			name: "multiple independent resources",
			resources: []discover.DiscoveredResource{
				{Name: "A"},
				{Name: "B"},
				{Name: "C"},
			},
			wantOrder: nil, // order may vary for independent resources
			wantErr:   false,
		},
		{
			name: "diamond dependency",
			resources: []discover.DiscoveredResource{
				{Name: "D", Dependencies: []string{"B", "C"}},
				{Name: "B", Dependencies: []string{"A"}},
				{Name: "C", Dependencies: []string{"A"}},
				{Name: "A"},
			},
			// A must come first, D must come last, B and C can be in any order
			wantOrder: nil, // we'll verify constraints instead
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewTemplateBuilder()
			for _, res := range tt.resources {
				err := builder.AddResource(res)
				require.NoError(t, err)
			}

			sorted, err := builder.topologicalSort()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Len(t, sorted, len(tt.resources))

			if tt.wantOrder != nil {
				// Verify exact order
				for i, expectedName := range tt.wantOrder {
					assert.Equal(t, expectedName, sorted[i].Name)
				}
			} else {
				// Verify dependency constraints
				positions := make(map[string]int)
				for i, res := range sorted {
					positions[res.Name] = i
				}

				for _, res := range sorted {
					for _, dep := range res.Dependencies {
						assert.Less(t, positions[dep], positions[res.Name],
							"dependency %s should come before %s", dep, res.Name)
					}
				}
			}
		})
	}
}
