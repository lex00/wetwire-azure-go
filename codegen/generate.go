package codegen

import (
	"bytes"
	"fmt"
	"go/format"
	"sort"
	"strings"
	"text/template"
	"unicode"
)

// StructGenerator generates Go structs from Azure ARM schemas
type StructGenerator struct {
	generatedTypes map[string]bool
}

// NewStructGenerator creates a new struct generator
func NewStructGenerator() *StructGenerator {
	return &StructGenerator{
		generatedTypes: make(map[string]bool),
	}
}

// StructField represents a field in a generated struct
type StructField struct {
	Name        string
	Type        string
	JSONTag     string
	Description string
	Required    bool
}

// GenerateStruct generates a Go struct from a resource schema
func (g *StructGenerator) GenerateStruct(structName string, schema *ResourceSchema) (string, error) {
	if g.generatedTypes == nil {
		g.generatedTypes = make(map[string]bool)
	}

	var buf bytes.Buffer

	// Add struct comment
	if schema.Description != "" {
		buf.WriteString(fmt.Sprintf("// %s represents %s\n", structName, schema.Description))
	} else {
		buf.WriteString(fmt.Sprintf("// %s represents an Azure resource\n", structName))
	}

	// Start struct definition
	buf.WriteString(fmt.Sprintf("type %s struct {\n", structName))

	// Generate fields
	fields := g.generateFields(schema)

	// Sort fields for consistent output (required first, then alphabetically)
	sort.Slice(fields, func(i, j int) bool {
		if fields[i].Required != fields[j].Required {
			return fields[i].Required
		}
		return fields[i].Name < fields[j].Name
	})

	// Write fields
	for _, field := range fields {
		if field.Description != "" {
			buf.WriteString(fmt.Sprintf("\t// %s\n", field.Description))
		}
		buf.WriteString(fmt.Sprintf("\t%s %s `json:\"%s\"`\n", field.Name, field.Type, field.JSONTag))
	}

	// Close struct
	buf.WriteString("}\n\n")

	return buf.String(), nil
}

// generateFields generates fields from a resource schema
func (g *StructGenerator) generateFields(schema *ResourceSchema) []StructField {
	var fields []StructField

	requiredMap := make(map[string]bool)
	for _, req := range schema.Required {
		requiredMap[req] = true
	}

	for propName, propSchema := range schema.Properties {
		isRequired := requiredMap[propName]
		field := StructField{
			Name:        ToGoFieldName(propName),
			Type:        ConvertToGoType(propSchema, isRequired),
			JSONTag:     g.generateJSONTag(propName, isRequired),
			Description: propSchema.Description,
			Required:    isRequired,
		}
		fields = append(fields, field)
	}

	return fields
}

// generateJSONTag generates a JSON tag for a field
func (g *StructGenerator) generateJSONTag(name string, required bool) string {
	if required {
		return name
	}
	return name + ",omitempty"
}

// ConvertToGoType converts a JSON schema type to a Go type
func ConvertToGoType(schema *PropertySchema, isRequired bool) string {
	var baseType string

	switch schema.Type {
	case "string":
		baseType = "string"
	case "integer":
		baseType = "int"
	case "number":
		baseType = "float64"
	case "boolean":
		baseType = "bool"
	case "array":
		if schema.Items != nil {
			itemType := ConvertToGoType(schema.Items, true)
			baseType = "[]" + itemType
		} else {
			baseType = "[]interface{}"
		}
		return baseType // Arrays don't need pointer wrapping
	case "object":
		if len(schema.Properties) > 0 {
			// Nested object - could generate a nested struct
			// For simplicity, using map for now
			baseType = "map[string]interface{}"
		} else {
			baseType = "map[string]interface{}"
		}
		return baseType // Maps don't need pointer wrapping
	default:
		baseType = "interface{}"
	}

	// Make optional fields pointers (except arrays and maps)
	if !isRequired && schema.Type != "array" && schema.Type != "object" {
		return "*" + baseType
	}

	return baseType
}

// ToGoFieldName converts a JSON field name to a Go field name
func ToGoFieldName(name string) string {
	// Handle common abbreviations
	abbreviations := map[string]string{
		"id":   "ID",
		"api":  "API",
		"vm":   "VM",
		"os":   "OS",
		"ip":   "IP",
		"url":  "URL",
		"uri":  "URI",
		"http": "HTTP",
		"https": "HTTPS",
		"dns":  "DNS",
		"sql":  "SQL",
		"json": "JSON",
		"xml":  "XML",
		"yaml": "YAML",
	}

	// Split on common delimiters
	words := splitWords(name)

	var result strings.Builder
	for _, word := range words {
		word = strings.ToLower(word)
		if abbr, ok := abbreviations[word]; ok {
			result.WriteString(abbr)
		} else {
			// Capitalize first letter manually instead of strings.Title
			if len(word) > 0 {
				result.WriteString(strings.ToUpper(word[:1]) + word[1:])
			}
		}
	}

	return result.String()
}

// splitWords splits a string into words based on various delimiters
func splitWords(s string) []string {
	var words []string
	var currentWord strings.Builder

	for i, r := range s {
		if r == '_' || r == '-' || r == '.' {
			if currentWord.Len() > 0 {
				words = append(words, currentWord.String())
				currentWord.Reset()
			}
		} else if unicode.IsUpper(r) && i > 0 {
			// Handle camelCase
			if currentWord.Len() > 0 {
				words = append(words, currentWord.String())
				currentWord.Reset()
			}
			currentWord.WriteRune(r)
		} else {
			currentWord.WriteRune(r)
		}
	}

	if currentWord.Len() > 0 {
		words = append(words, currentWord.String())
	}

	return words
}

// GeneratePackageFile generates a complete Go file with package declaration and imports
func (g *StructGenerator) GeneratePackageFile(packageName, structName string, schema *ResourceSchema) (string, error) {
	var buf bytes.Buffer

	// Package declaration
	buf.WriteString(fmt.Sprintf("// Package %s provides Azure %s resource types\n", packageName, packageName))
	buf.WriteString(fmt.Sprintf("package %s\n\n", packageName))

	// Generate struct
	structCode, err := g.GenerateStruct(structName, schema)
	if err != nil {
		return "", err
	}

	buf.WriteString(structCode)

	// Format the code
	return FormatGoCode(buf.String())
}

// GenerateAttrRefHelpers generates AttrRef helper code for cross-resource references
func GenerateAttrRefHelpers() string {
	return `// AttrRef represents a reference to another resource's attribute
type AttrRef struct {
	ResourceName string
	AttributePath string
}

// MarshalJSON implements json.Marshaler for AttrRef
func (a AttrRef) MarshalJSON() ([]byte, error) {
	// In ARM templates, this would be rendered as a reference() function call
	// For now, return the resource name as a placeholder
	return json.Marshal(a.ResourceName)
}

// Ref creates an AttrRef to another resource
func Ref(resourceName, attributePath string) AttrRef {
	return AttrRef{
		ResourceName: resourceName,
		AttributePath: attributePath,
	}
}
`
}

// FormatGoCode formats Go code using gofmt
func FormatGoCode(code string) (string, error) {
	formatted, err := format.Source([]byte(code))
	if err != nil {
		return "", fmt.Errorf("failed to format code: %w", err)
	}
	return string(formatted), nil
}

// GenerateResourceFile generates a complete resource file with all necessary types
func (g *StructGenerator) GenerateResourceFile(packageName, resourceName string, schema *ResourceSchema) (string, error) {
	var buf bytes.Buffer

	// Add file header
	tmpl := `// Package {{.Package}} provides Azure {{.Package}} resource types
// This file was generated from Azure ARM schemas
package {{.Package}}

`
	t := template.Must(template.New("header").Parse(tmpl))
	if err := t.Execute(&buf, map[string]string{"Package": packageName}); err != nil {
		return "", err
	}

	// Generate main struct
	structCode, err := g.GenerateStruct(resourceName, schema)
	if err != nil {
		return "", err
	}
	buf.WriteString(structCode)

	// Generate nested types if needed
	if err := g.generateNestedTypes(&buf, schema); err != nil {
		return "", err
	}

	return FormatGoCode(buf.String())
}

// generateNestedTypes generates structs for nested object properties
func (g *StructGenerator) generateNestedTypes(buf *bytes.Buffer, schema *ResourceSchema) error {
	// This is a placeholder for generating nested types
	// In a full implementation, this would recursively generate structs
	// for complex nested properties
	return nil
}

// TemplateData contains data for code generation templates
type TemplateData struct {
	Package     string
	StructName  string
	Fields      []StructField
	Description string
}
