// Package serialize provides utilities for converting Go structs to ARM template JSON format.
package serialize

import (
	"encoding/json"
	"reflect"

	"github.com/lex00/wetwire-azure-go/intrinsics"
)

// ToARMResource converts a Go struct resource to a map suitable for ARM template JSON.
// It respects JSON struct tags and handles nested structures, arrays, and ARM intrinsics.
func ToARMResource(resource any) map[string]any {
	return structToMap(reflect.ValueOf(resource))
}

// ToARMTemplate creates a complete ARM template with the given resources.
func ToARMTemplate(resources []any) map[string]any {
	armResources := make([]any, 0, len(resources))
	for _, res := range resources {
		armResources = append(armResources, ToARMResource(res))
	}

	return map[string]any{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"resources":      armResources,
	}
}

// ToARMTemplateJSON converts resources to a complete ARM template JSON.
func ToARMTemplateJSON(resources []any) ([]byte, error) {
	template := ToARMTemplate(resources)
	return json.MarshalIndent(template, "", "  ")
}

// SerializeValue serializes a single value, handling intrinsics specially.
func SerializeValue(v any) any {
	// Check if it's an intrinsic first
	if intrinsic, ok := v.(intrinsics.Intrinsic); ok {
		return intrinsic.ARMExpression()
	}

	return convertValue(reflect.ValueOf(v))
}

// structToMap converts a struct to a map using JSON tags.
func structToMap(v reflect.Value) map[string]any {
	// Handle pointer
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil
	}

	result := make(map[string]any)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Get JSON tag name
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		// Parse json tag (handle "name,omitempty")
		key, omitEmpty := parseJSONTag(jsonTag)

		// Check if we should omit empty values
		if omitEmpty && isZeroValue(fieldValue) {
			continue
		}

		// Convert the value
		value := convertValue(fieldValue)

		// Skip nil values for omitempty
		if omitEmpty && value == nil {
			continue
		}

		// Always skip nil values (from nil slices, maps, pointers)
		if value == nil {
			continue
		}

		// Skip empty slices (even if not marked omitempty, empty arrays should be omitted)
		if slice, ok := value.([]any); ok && len(slice) == 0 {
			continue
		}

		result[key] = value
	}

	return result
}

// convertValue converts a reflect.Value to an appropriate type for JSON serialization.
func convertValue(v reflect.Value) any {
	// Check for intrinsic types first
	if v.CanInterface() {
		if intrinsic, ok := v.Interface().(intrinsics.Intrinsic); ok {
			return intrinsic.ARMExpression()
		}
	}

	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return nil
		}
		return convertValue(v.Elem())

	case reflect.Struct:
		// Check again after dereferencing for intrinsics
		if v.CanInterface() {
			if intrinsic, ok := v.Interface().(intrinsics.Intrinsic); ok {
				return intrinsic.ARMExpression()
			}
		}
		return structToMap(v)

	case reflect.Slice, reflect.Array:
		if v.IsNil() || v.Len() == 0 {
			return nil
		}
		result := make([]any, v.Len())
		for i := 0; i < v.Len(); i++ {
			result[i] = convertValue(v.Index(i))
		}
		return result

	case reflect.Map:
		if v.IsNil() {
			return nil
		}
		result := make(map[string]any)
		for _, key := range v.MapKeys() {
			keyStr := key.String()
			result[keyStr] = convertValue(v.MapIndex(key))
		}
		return result

	case reflect.Interface:
		if v.IsNil() {
			return nil
		}
		return convertValue(v.Elem())

	default:
		if v.CanInterface() {
			return v.Interface()
		}
		return nil
	}
}

// isZeroValue checks if a value is the zero value for its type.
func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.Slice, reflect.Map, reflect.Array:
		return v.IsNil() || v.Len() == 0
	case reflect.Struct:
		return v.IsZero()
	default:
		return false
	}
}

// parseJSONTag parses a JSON tag and returns the key name and omitempty flag.
func parseJSONTag(tag string) (string, bool) {
	if tag == "" || tag == "-" {
		return "", false
	}

	// Handle "name,omitempty" format
	for i, c := range tag {
		if c == ',' {
			return tag[:i], tag[i:] == ",omitempty"
		}
	}

	return tag, false
}
