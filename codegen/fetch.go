// Package codegen provides schema fetching and Go code generation for Azure ARM resources.
package codegen

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	// DefaultSchemaBaseURL is the base URL for Azure ARM schemas
	DefaultSchemaBaseURL = "https://schema.management.azure.com/schemas"
)

// SchemaFetcher fetches Azure ARM schemas from the official repository
type SchemaFetcher struct {
	BaseURL string
	Client  *http.Client
}

// NewSchemaFetcher creates a new schema fetcher with default settings
func NewSchemaFetcher() *SchemaFetcher {
	return &SchemaFetcher{
		BaseURL: DefaultSchemaBaseURL,
		Client:  &http.Client{},
	}
}

// AzureSchema represents the top-level Azure ARM schema structure
type AzureSchema struct {
	ID                  string                        `json:"id"`
	Schema              string                        `json:"$schema"`
	Title               string                        `json:"title"`
	Description         string                        `json:"description"`
	ResourceDefinitions map[string]*ResourceSchema    `json:"resourceDefinitions"`
	Definitions         map[string]*PropertySchema    `json:"definitions"`
}

// ResourceSchema represents a resource type definition in the schema
type ResourceSchema struct {
	Type        string                     `json:"type"`
	Description string                     `json:"description"`
	Properties  map[string]*PropertySchema `json:"properties"`
	Required    []string                   `json:"required"`
	AllOf       []interface{}              `json:"allOf,omitempty"`
	OneOf       []interface{}              `json:"oneOf,omitempty"`
}

// PropertySchema represents a property within a resource schema
type PropertySchema struct {
	Type        string                     `json:"type"`
	Description string                     `json:"description"`
	Properties  map[string]*PropertySchema `json:"properties,omitempty"`
	Items       *PropertySchema            `json:"items,omitempty"`
	Enum        []interface{}              `json:"enum,omitempty"`
	Ref         string                     `json:"$ref,omitempty"`
	Default     interface{}                `json:"default,omitempty"`
	MinLength   int                        `json:"minLength,omitempty"`
	MaxLength   int                        `json:"maxLength,omitempty"`
	Pattern     string                     `json:"pattern,omitempty"`
}

// SchemaInfo contains information about an available schema
type SchemaInfo struct {
	Provider   string
	APIVersion string
	URL        string
}

// FetchSchema fetches a schema for a given provider and API version
func (f *SchemaFetcher) FetchSchema(ctx context.Context, provider, apiVersion string) (*AzureSchema, error) {
	url := fmt.Sprintf("%s/%s/%s.json", f.BaseURL, apiVersion, provider)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := f.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch schema: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("schema not found: %s (status: %d)", url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var schema AzureSchema
	if err := json.Unmarshal(body, &schema); err != nil {
		return nil, fmt.Errorf("failed to parse schema JSON: %w", err)
	}

	return &schema, nil
}

// ParseResourceType parses a resource type string into namespace and resource name
// Example: "Microsoft.Storage/storageAccounts" -> ("Microsoft.Storage", "storageAccounts")
func ParseResourceType(resourceType string) (namespace, name string, err error) {
	if resourceType == "" {
		return "", "", fmt.Errorf("resource type cannot be empty")
	}

	parts := strings.Split(resourceType, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid resource type format: %s (expected 'Namespace/ResourceType')", resourceType)
	}

	return parts[0], parts[1], nil
}

// ListAvailableSchemas returns a list of commonly available Azure schemas
// This is a curated list of common schemas; the full list can be retrieved from
// the Azure resource-manager-schemas repository
func ListAvailableSchemas() []SchemaInfo {
	return []SchemaInfo{
		{Provider: "Microsoft.Storage", APIVersion: "2021-04-01", URL: "https://schema.management.azure.com/schemas/2021-04-01/Microsoft.Storage.json"},
		{Provider: "Microsoft.Compute", APIVersion: "2021-07-01", URL: "https://schema.management.azure.com/schemas/2021-07-01/Microsoft.Compute.json"},
		{Provider: "Microsoft.Network", APIVersion: "2021-05-01", URL: "https://schema.management.azure.com/schemas/2021-05-01/Microsoft.Network.json"},
		{Provider: "Microsoft.KeyVault", APIVersion: "2021-06-01-preview", URL: "https://schema.management.azure.com/schemas/2021-06-01-preview/Microsoft.KeyVault.json"},
		{Provider: "Microsoft.Web", APIVersion: "2021-02-01", URL: "https://schema.management.azure.com/schemas/2021-02-01/Microsoft.Web.json"},
		{Provider: "Microsoft.Sql", APIVersion: "2021-05-01-preview", URL: "https://schema.management.azure.com/schemas/2021-05-01-preview/Microsoft.Sql.json"},
		{Provider: "Microsoft.ContainerService", APIVersion: "2021-07-01", URL: "https://schema.management.azure.com/schemas/2021-07-01/Microsoft.ContainerService.json"},
	}
}

// GetResourceSchema retrieves a specific resource definition from a schema
func (s *AzureSchema) GetResourceSchema(resourceName string) (*ResourceSchema, error) {
	if s.ResourceDefinitions == nil {
		return nil, fmt.Errorf("schema has no resource definitions")
	}

	resource, ok := s.ResourceDefinitions[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource %s not found in schema", resourceName)
	}

	return resource, nil
}

// ResolveReference resolves a $ref pointer within the schema
func (s *AzureSchema) ResolveReference(ref string) (*PropertySchema, error) {
	// References are typically in the format "#/definitions/DefinitionName"
	if !strings.HasPrefix(ref, "#/definitions/") {
		return nil, fmt.Errorf("unsupported reference format: %s", ref)
	}

	defName := strings.TrimPrefix(ref, "#/definitions/")
	def, ok := s.Definitions[defName]
	if !ok {
		return nil, fmt.Errorf("definition %s not found", defName)
	}

	return def, nil
}
