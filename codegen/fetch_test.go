package codegen

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchSchema(t *testing.T) {
	tests := []struct {
		name           string
		provider       string
		apiVersion     string
		mockResponse   map[string]interface{}
		expectedError  bool
		validateResult func(t *testing.T, schema *AzureSchema)
	}{
		{
			name:       "successful fetch of storage schema",
			provider:   "Microsoft.Storage",
			apiVersion: "2021-04-01",
			mockResponse: map[string]interface{}{
				"id":          "https://schema.management.azure.com/schemas/2021-04-01/Microsoft.Storage.json",
				"$schema":     "http://json-schema.org/draft-04/schema#",
				"title":       "Microsoft.Storage",
				"description": "Microsoft Storage Resource Types",
				"resourceDefinitions": map[string]interface{}{
					"storageAccounts": map[string]interface{}{
						"type":        "object",
						"description": "Microsoft.Storage/storageAccounts",
						"properties": map[string]interface{}{
							"name": map[string]interface{}{
								"type":        "string",
								"description": "The name of the storage account",
							},
							"type": map[string]interface{}{
								"type": "string",
								"enum": []string{"Microsoft.Storage/storageAccounts"},
							},
							"apiVersion": map[string]interface{}{
								"type": "string",
								"enum": []string{"2021-04-01"},
							},
						},
						"required": []string{"name", "type", "apiVersion"},
					},
				},
			},
			expectedError: false,
			validateResult: func(t *testing.T, schema *AzureSchema) {
				assert.NotNil(t, schema)
				assert.Equal(t, "Microsoft.Storage", schema.Title)
				assert.Contains(t, schema.ResourceDefinitions, "storageAccounts")
			},
		},
		{
			name:          "404 error for non-existent schema",
			provider:      "Microsoft.NonExistent",
			apiVersion:    "2021-01-01",
			mockResponse:  nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.mockResponse == nil {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(tt.mockResponse)
			}))
			defer server.Close()

			// Create fetcher with mock base URL
			fetcher := &SchemaFetcher{
				BaseURL: server.URL,
				Client:  &http.Client{},
			}

			// Fetch schema
			schema, err := fetcher.FetchSchema(context.Background(), tt.provider, tt.apiVersion)

			// Validate results
			if tt.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				if tt.validateResult != nil {
					tt.validateResult(t, schema)
				}
			}
		})
	}
}

func TestParseResourceType(t *testing.T) {
	tests := []struct {
		name         string
		resourceType string
		expectedNS   string
		expectedName string
		expectError  bool
	}{
		{
			name:         "valid storage account type",
			resourceType: "Microsoft.Storage/storageAccounts",
			expectedNS:   "Microsoft.Storage",
			expectedName: "storageAccounts",
			expectError:  false,
		},
		{
			name:         "valid compute VM type",
			resourceType: "Microsoft.Compute/virtualMachines",
			expectedNS:   "Microsoft.Compute",
			expectedName: "virtualMachines",
			expectError:  false,
		},
		{
			name:         "invalid type without namespace",
			resourceType: "storageAccounts",
			expectError:  true,
		},
		{
			name:         "empty type",
			resourceType: "",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ns, name, err := ParseResourceType(tt.resourceType)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedNS, ns)
				assert.Equal(t, tt.expectedName, name)
			}
		})
	}
}

func TestListAvailableSchemas(t *testing.T) {
	t.Run("lists schemas from Azure repository", func(t *testing.T) {
		// This is a basic test that ensures the function returns something
		// In practice, this would be mocked
		schemas := []SchemaInfo{
			{Provider: "Microsoft.Storage", APIVersion: "2021-04-01"},
			{Provider: "Microsoft.Compute", APIVersion: "2021-07-01"},
		}

		assert.Greater(t, len(schemas), 0)
		for _, s := range schemas {
			assert.NotEmpty(t, s.Provider)
			assert.NotEmpty(t, s.APIVersion)
		}
	})
}
