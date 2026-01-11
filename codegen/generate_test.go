package codegen

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateStruct(t *testing.T) {
	tests := []struct {
		name           string
		schema         *ResourceSchema
		structName     string
		expectedFields []string
		expectedTags   []string
	}{
		{
			name: "generates struct with basic fields",
			schema: &ResourceSchema{
				Type:        "object",
				Description: "Storage account resource",
				Properties: map[string]*PropertySchema{
					"name": {
						Type:        "string",
						Description: "The name of the storage account",
					},
					"location": {
						Type:        "string",
						Description: "The location of the resource",
					},
					"tags": {
						Type:        "object",
						Description: "Resource tags",
					},
				},
				Required: []string{"name", "location"},
			},
			structName: "StorageAccount",
			expectedFields: []string{
				"Name",
				"Location",
				"Tags",
			},
			expectedTags: []string{
				`json:"name"`,
				`json:"location"`,
				`json:"tags,omitempty"`,
			},
		},
		{
			name: "generates struct with nested properties",
			schema: &ResourceSchema{
				Type: "object",
				Properties: map[string]*PropertySchema{
					"properties": {
						Type: "object",
						Properties: map[string]*PropertySchema{
							"accountType": {
								Type: "string",
							},
							"encryption": {
								Type: "object",
								Properties: map[string]*PropertySchema{
									"keySource": {
										Type: "string",
									},
								},
							},
						},
					},
				},
			},
			structName: "StorageAccount",
			expectedFields: []string{
				"Properties",
			},
		},
		{
			name: "generates struct with array fields",
			schema: &ResourceSchema{
				Type: "object",
				Properties: map[string]*PropertySchema{
					"ipRules": {
						Type: "array",
						Items: &PropertySchema{
							Type: "object",
							Properties: map[string]*PropertySchema{
								"value": {
									Type: "string",
								},
							},
						},
					},
				},
			},
			structName: "NetworkRuleSet",
			expectedFields: []string{
				"IPRules",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generator := &StructGenerator{}
			code, err := generator.GenerateStruct(tt.structName, tt.schema)

			require.NoError(t, err)
			assert.NotEmpty(t, code)

			// Verify struct name
			assert.Contains(t, code, "type "+tt.structName+" struct")

			// Verify expected fields
			for _, field := range tt.expectedFields {
				assert.Contains(t, code, field)
			}

			// Verify expected tags
			for _, tag := range tt.expectedTags {
				assert.Contains(t, code, tag)
			}
		})
	}
}

func TestConvertToGoType(t *testing.T) {
	tests := []struct {
		name       string
		schema     *PropertySchema
		expected   string
		isPointer  bool
		isRequired bool
	}{
		{
			name:       "string type",
			schema:     &PropertySchema{Type: "string"},
			expected:   "string",
			isPointer:  false,
			isRequired: true,
		},
		{
			name:       "optional string type",
			schema:     &PropertySchema{Type: "string"},
			expected:   "*string",
			isPointer:  true,
			isRequired: false,
		},
		{
			name:       "integer type",
			schema:     &PropertySchema{Type: "integer"},
			expected:   "int",
			isPointer:  false,
			isRequired: true,
		},
		{
			name:       "boolean type",
			schema:     &PropertySchema{Type: "boolean"},
			expected:   "bool",
			isPointer:  false,
			isRequired: true,
		},
		{
			name: "array of strings",
			schema: &PropertySchema{
				Type:  "array",
				Items: &PropertySchema{Type: "string"},
			},
			expected:   "[]string",
			isPointer:  false,
			isRequired: true,
		},
		{
			name: "array of objects",
			schema: &PropertySchema{
				Type: "array",
				Items: &PropertySchema{
					Type: "object",
				},
			},
			expected:   "[]map[string]interface{}",
			isPointer:  false,
			isRequired: true,
		},
		{
			name:       "object type",
			schema:     &PropertySchema{Type: "object"},
			expected:   "map[string]interface{}",
			isPointer:  false,
			isRequired: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertToGoType(tt.schema, tt.isRequired)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGeneratePackageFile(t *testing.T) {
	t.Run("generates complete Go file with package and imports", func(t *testing.T) {
		schema := &ResourceSchema{
			Type:        "object",
			Description: "Storage account resource",
			Properties: map[string]*PropertySchema{
				"name": {
					Type:        "string",
					Description: "The name of the storage account",
				},
			},
			Required: []string{"name"},
		}

		generator := &StructGenerator{}
		code, err := generator.GeneratePackageFile("storage", "StorageAccount", schema)

		require.NoError(t, err)
		assert.NotEmpty(t, code)

		// Verify package declaration
		assert.Contains(t, code, "package storage")

		// Verify struct generation
		assert.Contains(t, code, "type StorageAccount struct")
	})
}

func TestToGoFieldName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "camelCase to PascalCase",
			input:    "accountType",
			expected: "AccountType",
		},
		{
			name:     "snake_case to PascalCase",
			input:    "storage_account",
			expected: "StorageAccount",
		},
		{
			name:     "already PascalCase",
			input:    "Name",
			expected: "Name",
		},
		{
			name:     "single word lowercase",
			input:    "tags",
			expected: "Tags",
		},
		{
			name:     "with numbers",
			input:    "vmSize",
			expected: "VMSize",
		},
		{
			name:     "API field name",
			input:    "apiVersion",
			expected: "APIVersion",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToGoFieldName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateAttrRef(t *testing.T) {
	t.Run("generates AttrRef helper functions", func(t *testing.T) {
		code := GenerateAttrRefHelpers()

		assert.NotEmpty(t, code)
		assert.Contains(t, code, "type AttrRef struct")
		assert.Contains(t, code, "func (a AttrRef) MarshalJSON()")
	})
}

func TestGenerateResourceWithReferences(t *testing.T) {
	t.Run("includes AttrRef fields for cross-resource references", func(t *testing.T) {
		schema := &ResourceSchema{
			Type: "object",
			Properties: map[string]*PropertySchema{
				"name": {
					Type: "string",
				},
				"subnetId": {
					Type:        "string",
					Description: "Resource ID reference",
				},
			},
			Required: []string{"name"},
		}

		generator := &StructGenerator{}
		code, err := generator.GenerateStruct("NetworkInterface", schema)

		require.NoError(t, err)

		// The generator should recognize *Id fields as potential references
		// and could optionally support AttrRef
		assert.Contains(t, code, "SubnetID")
	})
}

func TestFormatGoCode(t *testing.T) {
	t.Run("formats generated code with gofmt", func(t *testing.T) {
		unformatted := `package storage
type StorageAccount struct {
Name string
Location string
}`

		formatted, err := FormatGoCode(unformatted)
		require.NoError(t, err)

		// Should have proper formatting
		assert.Contains(t, formatted, "package storage")
		assert.Contains(t, formatted, "type StorageAccount struct")

		// Should not have trailing spaces
		lines := strings.Split(formatted, "\n")
		for _, line := range lines {
			assert.Equal(t, strings.TrimRight(line, " "), line)
		}
	})

	t.Run("returns error for invalid Go code", func(t *testing.T) {
		invalid := `package storage
type StorageAccount struct {
Name string
// missing closing brace
`

		_, err := FormatGoCode(invalid)
		assert.Error(t, err)
	})
}
