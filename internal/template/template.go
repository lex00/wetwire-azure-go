// Package template provides ARM template building and generation functionality.
package template

import (
	"encoding/json"
	"fmt"

	"github.com/lex00/wetwire-azure-go/internal/discover"
)

// TemplateBuilder aggregates resources, parameters, variables, and outputs
// to build a complete ARM template.
type TemplateBuilder struct {
	resources  map[string]discover.DiscoveredResource
	parameters map[string]Parameter
	variables  map[string]interface{}
	outputs    map[string]Output
}

// Parameter represents an ARM template parameter
type Parameter struct {
	Type          string                 `json:"type"`
	DefaultValue  interface{}            `json:"defaultValue,omitempty"`
	AllowedValues []interface{}          `json:"allowedValues,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// Output represents an ARM template output
type Output struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

// ARMTemplate represents the complete ARM template structure
type ARMTemplate struct {
	Schema         string                 `json:"$schema"`
	ContentVersion string                 `json:"contentVersion"`
	Parameters     map[string]Parameter   `json:"parameters"`
	Variables      map[string]interface{} `json:"variables"`
	Resources      []ARMResource          `json:"resources"`
	Outputs        map[string]Output      `json:"outputs"`
}

// ARMResource represents a resource in the ARM template
type ARMResource struct {
	Name       string        `json:"name"`
	Type       string        `json:"type"`
	APIVersion string        `json:"apiVersion"`
	Location   string        `json:"location,omitempty"`
	DependsOn  []string      `json:"dependsOn,omitempty"`
	Properties interface{}   `json:"properties,omitempty"`
	Tags       interface{}   `json:"tags,omitempty"`
	SKU        interface{}   `json:"sku,omitempty"`
	Kind       string        `json:"kind,omitempty"`
	Identity   interface{}   `json:"identity,omitempty"`
	Zones      []string      `json:"zones,omitempty"`
	Plan       interface{}   `json:"plan,omitempty"`
}

// NewTemplateBuilder creates a new TemplateBuilder instance
func NewTemplateBuilder() *TemplateBuilder {
	return &TemplateBuilder{
		resources:  make(map[string]discover.DiscoveredResource),
		parameters: make(map[string]Parameter),
		variables:  make(map[string]interface{}),
		outputs:    make(map[string]Output),
	}
}

// AddResource adds a discovered resource to the template builder.
// Returns an error if a resource with the same name already exists.
func (tb *TemplateBuilder) AddResource(resource discover.DiscoveredResource) error {
	if _, exists := tb.resources[resource.Name]; exists {
		return fmt.Errorf("resource with name %s already exists", resource.Name)
	}
	tb.resources[resource.Name] = resource
	return nil
}

// AddParameter adds a parameter to the template.
// Returns an error if a parameter with the same name already exists.
func (tb *TemplateBuilder) AddParameter(name, paramType string, metadata map[string]interface{}) error {
	if _, exists := tb.parameters[name]; exists {
		return fmt.Errorf("parameter with name %s already exists", name)
	}

	param := Parameter{
		Type: paramType,
	}

	if metadata != nil {
		if defaultValue, ok := metadata["defaultValue"]; ok {
			param.DefaultValue = defaultValue
		}
		if allowedValues, ok := metadata["allowedValues"]; ok {
			if values, ok := allowedValues.([]interface{}); ok {
				param.AllowedValues = values
			}
		}
		if md, ok := metadata["metadata"]; ok {
			if metadataMap, ok := md.(map[string]interface{}); ok {
				param.Metadata = metadataMap
			}
		}
	}

	tb.parameters[name] = param
	return nil
}

// AddVariable adds a variable to the template.
// Returns an error if a variable with the same name already exists.
func (tb *TemplateBuilder) AddVariable(name string, value interface{}) error {
	if _, exists := tb.variables[name]; exists {
		return fmt.Errorf("variable with name %s already exists", name)
	}
	tb.variables[name] = value
	return nil
}

// AddOutput adds an output to the template.
// Returns an error if an output with the same name already exists.
func (tb *TemplateBuilder) AddOutput(name, outputType string, value interface{}) error {
	if _, exists := tb.outputs[name]; exists {
		return fmt.Errorf("output with name %s already exists", name)
	}
	tb.outputs[name] = Output{
		Type:  outputType,
		Value: value,
	}
	return nil
}

// Build executes the build pipeline and returns the ARM template JSON.
// Pipeline stages: DISCOVER → VALIDATE → ORDER → SERIALIZE → EMIT
func (tb *TemplateBuilder) Build() (string, error) {
	// DISCOVER - resources are already discovered and added via AddResource

	// VALIDATE - check references and detect cycles
	if err := tb.validateReferences(); err != nil {
		return "", fmt.Errorf("validation failed: %w", err)
	}

	// ORDER - topological sort by dependencies
	orderedResources, err := tb.topologicalSort()
	if err != nil {
		return "", fmt.Errorf("ordering failed: %w", err)
	}

	// SERIALIZE - convert to ARM JSON format
	template := tb.serialize(orderedResources)

	// EMIT - write output as JSON
	jsonBytes, err := json.MarshalIndent(template, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON serialization failed: %w", err)
	}

	return string(jsonBytes), nil
}

// validateReferences checks that all referenced resources exist and detects cycles
func (tb *TemplateBuilder) validateReferences() error {
	// Check that all dependencies exist
	for name, resource := range tb.resources {
		for _, dep := range resource.Dependencies {
			if _, exists := tb.resources[dep]; !exists {
				return fmt.Errorf("resource %s depends on non-existent resource %s", name, dep)
			}
		}
	}

	// Detect cycles using DFS
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	var hasCycle func(string) bool
	hasCycle = func(name string) bool {
		visited[name] = true
		recStack[name] = true

		resource := tb.resources[name]
		for _, dep := range resource.Dependencies {
			if !visited[dep] {
				if hasCycle(dep) {
					return true
				}
			} else if recStack[dep] {
				return true
			}
		}

		recStack[name] = false
		return false
	}

	for name := range tb.resources {
		if !visited[name] {
			if hasCycle(name) {
				return fmt.Errorf("cyclic dependency detected in resource graph")
			}
		}
	}

	return nil
}

// topologicalSort performs a topological sort on resources using Kahn's algorithm
func (tb *TemplateBuilder) topologicalSort() ([]discover.DiscoveredResource, error) {
	// Build in-degree map
	inDegree := make(map[string]int)
	for name := range tb.resources {
		inDegree[name] = 0
	}

	for _, resource := range tb.resources {
		for range resource.Dependencies {
			inDegree[resource.Name]++
		}
	}

	// Initialize queue with resources that have no dependencies
	queue := []string{}
	for name, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, name)
		}
	}

	// Process queue
	var sorted []discover.DiscoveredResource
	for len(queue) > 0 {
		// Dequeue
		current := queue[0]
		queue = queue[1:]

		sorted = append(sorted, tb.resources[current])

		// Find all resources that depend on current
		for name, resource := range tb.resources {
			for _, dep := range resource.Dependencies {
				if dep == current {
					inDegree[name]--
					if inDegree[name] == 0 {
						queue = append(queue, name)
					}
				}
			}
		}
	}

	// If we didn't process all resources, there's a cycle
	if len(sorted) != len(tb.resources) {
		return nil, fmt.Errorf("cyclic dependency detected")
	}

	return sorted, nil
}

// serialize converts the ordered resources into an ARM template structure
func (tb *TemplateBuilder) serialize(orderedResources []discover.DiscoveredResource) ARMTemplate {
	armResources := make([]ARMResource, 0, len(orderedResources))

	for _, resource := range orderedResources {
		armResource := ARMResource{
			Name:       resource.Name,
			Type:       resource.Type,
			APIVersion: getAPIVersion(resource.Type),
			Location:   "[resourceGroup().location]",
		}

		// Add dependsOn if there are dependencies
		if len(resource.Dependencies) > 0 {
			dependsOn := make([]string, 0, len(resource.Dependencies))
			for _, dep := range resource.Dependencies {
				depResource := tb.resources[dep]
				dependsOn = append(dependsOn, fmt.Sprintf("[resourceId('%s', '%s')]", depResource.Type, dep))
			}
			armResource.DependsOn = dependsOn
		}

		armResources = append(armResources, armResource)
	}

	return ARMTemplate{
		Schema:         "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		ContentVersion: "1.0.0.0",
		Parameters:     tb.parameters,
		Variables:      tb.variables,
		Resources:      armResources,
		Outputs:        tb.outputs,
	}
}

// getAPIVersion returns the appropriate API version for a given resource type
func getAPIVersion(resourceType string) string {
	apiVersions := map[string]string{
		"Microsoft.Storage/storageAccounts":          "2021-04-01",
		"Microsoft.Compute/virtualMachines":          "2021-07-01",
		"Microsoft.Network/virtualNetworks":          "2021-02-01",
		"Microsoft.Network/networkInterfaces":        "2021-02-01",
		"Microsoft.Network/publicIPAddresses":        "2021-02-01",
		"Microsoft.Network/networkSecurityGroups":    "2021-02-01",
		"Microsoft.KeyVault/vaults":                  "2021-06-01",
		"Microsoft.Sql/servers":                      "2021-02-01",
		"Microsoft.Sql/servers/databases":            "2021-02-01",
		"Microsoft.Web/sites":                        "2021-01-15",
		"Microsoft.ContainerRegistry/registries":     "2021-06-01",
		"Microsoft.ContainerService/managedClusters": "2021-05-01",
	}

	if version, ok := apiVersions[resourceType]; ok {
		return version
	}

	return "2021-04-01" // default version
}
