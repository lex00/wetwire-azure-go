// Package importer provides functionality to import ARM JSON templates
// and convert them to Go struct declarations for use with wetwire-azure-go.
package importer

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// ARMTemplate represents a parsed ARM template.
type ARMTemplate struct {
	Schema         string                 `json:"$schema"`
	ContentVersion string                 `json:"contentVersion"`
	Parameters     map[string]interface{} `json:"parameters"`
	Variables      map[string]interface{} `json:"variables"`
	Resources      []ARMResource          `json:"resources"`
	Outputs        map[string]interface{} `json:"outputs"`
}

// ARMResource represents a resource in an ARM template.
type ARMResource struct {
	Type       string                 `json:"type"`
	APIVersion string                 `json:"apiVersion"`
	Name       string                 `json:"name"`
	Location   string                 `json:"location"`
	Kind       string                 `json:"kind,omitempty"`
	SKU        map[string]interface{} `json:"sku,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	Tags       map[string]string      `json:"tags,omitempty"`
	DependsOn  []string               `json:"dependsOn,omitempty"`
	Identity   map[string]interface{} `json:"identity,omitempty"`
	Zones      []string               `json:"zones,omitempty"`
	Plan       map[string]interface{} `json:"plan,omitempty"`
}

// ParseARMTemplate parses an ARM JSON template from bytes.
func ParseARMTemplate(data []byte) (*ARMTemplate, error) {
	var template ARMTemplate
	if err := json.Unmarshal(data, &template); err != nil {
		return nil, fmt.Errorf("failed to parse ARM template: %w", err)
	}
	return &template, nil
}

// acronyms maps lowercase acronym patterns to their uppercase versions.
var acronyms = map[string]string{
	"api":   "API",
	"os":    "OS",
	"vm":    "VM",
	"ip":    "IP",
	"https": "HTTPS",
	"http":  "HTTP",
	"ssh":   "SSH",
	"id":    "ID",
	"sku":   "SKU",
	"nic":   "NIC",
	"gb":    "GB",
	"vnet":  "VNet",
}

// CamelToPascal converts camelCase to PascalCase with special handling
// for common Azure acronyms (VM, OS, IP, API, HTTPS, SSH, ID, SKU, NIC, GB, VNet).
func CamelToPascal(s string) string {
	if s == "" {
		return ""
	}

	// First, capitalize the first letter
	result := strings.ToUpper(string(s[0])) + s[1:]

	// Apply acronym replacements
	result = applyAcronymReplacements(result)

	return result
}

// applyAcronymReplacements applies all acronym replacements to a string.
func applyAcronymReplacements(s string) string {
	result := s

	// Sort acronyms by length (longest first) to avoid partial replacements
	type kv struct {
		lower string
		upper string
	}
	var sorted []kv
	for k, v := range acronyms {
		sorted = append(sorted, kv{k, v})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return len(sorted[i].lower) > len(sorted[j].lower)
	})

	for _, item := range sorted {
		lower := item.lower
		upper := item.upper

		// Handle start of string
		if len(result) >= len(lower) {
			prefix := strings.ToLower(result[:len(lower)])
			if prefix == lower {
				result = upper + result[len(lower):]
			}
		}

		// Handle mid-string occurrences (Title case after camel boundary)
		// e.g., "enableHttpsTrafficOnly" -> first letter caps makes "EnableHttpsTrafficOnly"
		// We need to find "Https" and replace with "HTTPS"
		titleCase := strings.ToUpper(string(lower[0])) + lower[1:]
		result = strings.ReplaceAll(result, titleCase, upper)

		// Handle occurrences at word boundaries (after another uppercase or at specific positions)
		// e.g., "NICId" -> "NICID" (after NIC, we have Id which should become ID)
		for _, item2 := range sorted {
			// Check if this acronym follows another acronym
			prefix := item2.upper
			suffix := strings.ToUpper(string(lower[0])) + lower[1:]
			if strings.Contains(result, prefix+suffix) {
				result = strings.ReplaceAll(result, prefix+suffix, prefix+upper)
			}
		}
	}

	return result
}

// ResourceTypeToPackage converts an Azure resource type to a Go package name and type name.
// For example: "Microsoft.Storage/storageAccounts" -> ("storage", "StorageAccount")
func ResourceTypeToPackage(resourceType string) (pkgName, typeName string) {
	parts := strings.Split(resourceType, "/")
	if len(parts) < 2 {
		return "", ""
	}

	// Extract provider (e.g., "Microsoft.Storage" -> "storage")
	provider := parts[0]
	providerParts := strings.Split(provider, ".")
	if len(providerParts) >= 2 {
		pkgName = strings.ToLower(providerParts[1])
	}

	// Extract resource type name
	resourceName := parts[len(parts)-1]

	// Convert plural to singular and to PascalCase
	// storageAccounts -> StorageAccount
	typeName = singularize(resourceName)
	typeName = strings.ToUpper(string(typeName[0])) + typeName[1:]

	// Handle special cases
	switch pkgName {
	case "keyvault":
		// Keep keyvault as is
	}

	return pkgName, typeName
}

// singularize converts a plural noun to singular (simple implementation).
func singularize(s string) string {
	if strings.HasSuffix(s, "ies") {
		return s[:len(s)-3] + "y"
	}
	if strings.HasSuffix(s, "ses") {
		return s[:len(s)-2]
	}
	if strings.HasSuffix(s, "s") && !strings.HasSuffix(s, "ss") {
		return s[:len(s)-1]
	}
	return s
}

// GenerateVarName generates a valid Go variable name from a resource name.
// Converts kebab-case and snake_case to PascalCase, with acronym handling.
func GenerateVarName(name string) string {
	if name == "" {
		return ""
	}

	// Split by hyphen and underscore
	parts := strings.FieldsFunc(name, func(r rune) bool {
		return r == '-' || r == '_'
	})

	var result strings.Builder
	for _, part := range parts {
		if part == "" {
			continue
		}
		// Capitalize first letter of each part
		capitalized := strings.ToUpper(string(part[0])) + part[1:]
		result.WriteString(capitalized)
	}

	// Apply acronym replacements to the final result
	return applyAcronymReplacements(result.String())
}

// ExtractDependencyName extracts the resource name from a dependsOn expression.
func ExtractDependencyName(dependsOn string) string {
	// Pattern: [resourceId('Microsoft.Type/resources', 'name')]
	resourceIdPattern := regexp.MustCompile(`resourceId\([^,]+,\s*'([^']+)'\)`)
	matches := resourceIdPattern.FindStringSubmatch(dependsOn)
	if len(matches) >= 2 {
		return matches[1]
	}

	// Pattern: Microsoft.Type/resources/name
	if !strings.HasPrefix(dependsOn, "[") {
		parts := strings.Split(dependsOn, "/")
		if len(parts) >= 3 {
			return parts[len(parts)-1]
		}
	}

	return ""
}

// GenerateImports generates the import block for the given resource types.
func GenerateImports(resourceTypes []string) string {
	seen := make(map[string]bool)
	var imports []string

	for _, rt := range resourceTypes {
		pkgName, _ := ResourceTypeToPackage(rt)
		if pkgName != "" && !seen[pkgName] {
			seen[pkgName] = true
			imports = append(imports, fmt.Sprintf(`"github.com/lex00/wetwire-azure-go/resources/%s"`, pkgName))
		}
	}

	sort.Strings(imports)

	if len(imports) == 0 {
		return ""
	}

	return "import (\n\t" + strings.Join(imports, "\n\t") + "\n)"
}

// GenerateGoCode generates Go source code from an ARM template.
func GenerateGoCode(template *ARMTemplate, packageName string) (string, error) {
	var sb strings.Builder

	// Package declaration
	sb.WriteString(fmt.Sprintf("package %s\n\n", packageName))

	// Collect resource types for imports
	var resourceTypes []string
	for _, res := range template.Resources {
		resourceTypes = append(resourceTypes, res.Type)
	}

	// Generate imports
	imports := GenerateImports(resourceTypes)
	if imports != "" {
		sb.WriteString(imports)
		sb.WriteString("\n\n")
	}

	// Build a map of resource names for dependency resolution
	resourceMap := make(map[string]string) // ARM name -> Go var name
	for _, res := range template.Resources {
		resourceMap[res.Name] = GenerateVarName(res.Name)
	}

	// Generate each resource
	for i, res := range template.Resources {
		if i > 0 {
			sb.WriteString("\n")
		}

		code, err := generateResourceCode(res, resourceMap)
		if err != nil {
			return "", fmt.Errorf("failed to generate code for resource %s: %w", res.Name, err)
		}
		sb.WriteString(code)
	}

	return sb.String(), nil
}

// generateResourceCode generates Go code for a single ARM resource.
func generateResourceCode(res ARMResource, resourceMap map[string]string) (string, error) {
	var sb strings.Builder

	pkgName, typeName := ResourceTypeToPackage(res.Type)
	varName := GenerateVarName(res.Name)

	// Generate dependsOn comments
	if len(res.DependsOn) > 0 {
		for _, dep := range res.DependsOn {
			depName := ExtractDependencyName(dep)
			if depName != "" {
				if goVarName, ok := resourceMap[depName]; ok {
					sb.WriteString(fmt.Sprintf("// DependsOn: %s\n", goVarName))
				}
			}
		}
	}

	// Start struct declaration
	sb.WriteString(fmt.Sprintf("var %s = %s.%s{\n", varName, pkgName, typeName))

	// Add basic fields
	sb.WriteString(fmt.Sprintf("\tName:     %q,\n", res.Name))
	if res.Location != "" {
		sb.WriteString(fmt.Sprintf("\tLocation: %q,\n", res.Location))
	}
	if res.Kind != "" {
		sb.WriteString(fmt.Sprintf("\tKind:     %q,\n", res.Kind))
	}

	// Add SKU if present
	if len(res.SKU) > 0 {
		skuCode := generateStructCode(res.SKU, pkgName+".SKU", 1)
		sb.WriteString(fmt.Sprintf("\tSKU: %s,\n", skuCode))
	}

	// Add tags if present
	if len(res.Tags) > 0 {
		sb.WriteString("\tTags: map[string]string{\n")
		// Sort keys for deterministic output
		var keys []string
		for k := range res.Tags {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			sb.WriteString(fmt.Sprintf("\t\t%q: %q,\n", k, res.Tags[k]))
		}
		sb.WriteString("\t},\n")
	}

	// Add properties if present
	if len(res.Properties) > 0 {
		propsCode := generatePropertiesCode(res.Properties, pkgName, typeName, 1)
		sb.WriteString(fmt.Sprintf("\tProperties: %s,\n", propsCode))
	}

	sb.WriteString("}\n")

	return sb.String(), nil
}

// generateStructCode generates Go struct literal code from a map.
func generateStructCode(data map[string]interface{}, structType string, indent int) string {
	var sb strings.Builder
	indentStr := strings.Repeat("\t", indent)

	sb.WriteString(structType + "{\n")

	// Sort keys for deterministic output
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := data[k]
		fieldName := CamelToPascal(k)
		sb.WriteString(indentStr + "\t" + fieldName + ": ")
		sb.WriteString(generateValueCode(v, indent+1))
		sb.WriteString(",\n")
	}

	sb.WriteString(indentStr + "}")

	return sb.String()
}

// generatePropertiesCode generates Go code for the Properties field.
func generatePropertiesCode(props map[string]interface{}, pkgName, typeName string, indent int) string {
	var sb strings.Builder
	indentStr := strings.Repeat("\t", indent)

	sb.WriteString(pkgName + "." + typeName + "Properties{\n")

	// Sort keys for deterministic output
	var keys []string
	for k := range props {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := props[k]
		fieldName := CamelToPascal(k)
		sb.WriteString(indentStr + "\t" + fieldName + ": ")
		sb.WriteString(generateNestedCode(v, pkgName, fieldName, indent+1))
		sb.WriteString(",\n")
	}

	sb.WriteString(indentStr + "}")

	return sb.String()
}

// generateNestedCode generates Go code for nested structures.
func generateNestedCode(v interface{}, pkgName, fieldName string, indent int) string {
	switch val := v.(type) {
	case map[string]interface{}:
		return generateTypedStructCode(val, pkgName, fieldName, indent)
	case []interface{}:
		return generateSliceCode(val, indent)
	default:
		return generateValueCode(v, indent)
	}
}

// generateTypedStructCode generates a typed struct literal.
func generateTypedStructCode(data map[string]interface{}, pkgName, fieldName string, indent int) string {
	var sb strings.Builder
	indentStr := strings.Repeat("\t", indent)

	sb.WriteString(pkgName + "." + fieldName + "{\n")

	// Sort keys for deterministic output
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := data[k]
		subFieldName := CamelToPascal(k)
		sb.WriteString(indentStr + "\t" + subFieldName + ": ")
		sb.WriteString(generateNestedCode(v, pkgName, subFieldName, indent+1))
		sb.WriteString(",\n")
	}

	sb.WriteString(indentStr + "}")

	return sb.String()
}

// generateSliceCode generates Go code for slice values.
func generateSliceCode(data []interface{}, indent int) string {
	var sb strings.Builder
	indentStr := strings.Repeat("\t", indent)

	sb.WriteString("[]interface{}{\n")
	for _, v := range data {
		sb.WriteString(indentStr + "\t")
		sb.WriteString(generateValueCode(v, indent+1))
		sb.WriteString(",\n")
	}
	sb.WriteString(indentStr + "}")

	return sb.String()
}

// generateValueCode generates Go code for a single value.
func generateValueCode(v interface{}, indent int) string {
	switch val := v.(type) {
	case string:
		return fmt.Sprintf("%q", val)
	case float64:
		// Check if it's an integer
		if val == float64(int(val)) {
			return fmt.Sprintf("%d", int(val))
		}
		return fmt.Sprintf("%v", val)
	case bool:
		return fmt.Sprintf("%v", val)
	case nil:
		return "nil"
	case map[string]interface{}:
		return generateMapCode(val, indent)
	case []interface{}:
		return generateSliceCode(val, indent)
	default:
		return fmt.Sprintf("%v", val)
	}
}

// generateMapCode generates Go code for a map value.
func generateMapCode(data map[string]interface{}, indent int) string {
	var sb strings.Builder
	indentStr := strings.Repeat("\t", indent)

	sb.WriteString("map[string]interface{}{\n")

	// Sort keys for deterministic output
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := data[k]
		sb.WriteString(indentStr + "\t" + fmt.Sprintf("%q", k) + ": ")
		sb.WriteString(generateValueCode(v, indent+1))
		sb.WriteString(",\n")
	}

	sb.WriteString(indentStr + "}")

	return sb.String()
}
