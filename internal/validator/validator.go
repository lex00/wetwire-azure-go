// Package validator provides ARM template validation functionality.
package validator

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Severity represents the severity level of a validation result.
type Severity int

const (
	SeverityInfo Severity = iota
	SeverityWarning
	SeverityError
)

// String returns the string representation of the severity.
func (s Severity) String() string {
	switch s {
	case SeverityError:
		return "error"
	case SeverityWarning:
		return "warning"
	case SeverityInfo:
		return "info"
	default:
		return "unknown"
	}
}

// ValidationResult represents a single validation finding.
type ValidationResult struct {
	Severity Severity
	Message  string
	Field    string
}

// String returns a formatted string representation of the validation result.
func (r ValidationResult) String() string {
	if r.Field != "" {
		return fmt.Sprintf("[%s] %s: %s", r.Severity.String(), r.Field, r.Message)
	}
	return fmt.Sprintf("[%s] %s", r.Severity.String(), r.Message)
}

// Validator validates ARM templates.
type Validator struct{}

// NewValidator creates a new Validator instance.
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateTemplate validates ARM template JSON and returns validation results.
func (v *Validator) ValidateTemplate(data []byte) ([]ValidationResult, error) {
	var template map[string]interface{}
	if err := json.Unmarshal(data, &template); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	var results []ValidationResult

	// Check for required fields
	if _, ok := template["$schema"]; !ok {
		results = append(results, ValidationResult{
			Severity: SeverityError,
			Field:    "$schema",
			Message:  "missing required field",
		})
	} else {
		// Validate schema URL
		schema, ok := template["$schema"].(string)
		if ok && !isValidSchemaURL(schema) {
			results = append(results, ValidationResult{
				Severity: SeverityWarning,
				Field:    "$schema",
				Message:  "unrecognized schema URL",
			})
		}
	}

	if _, ok := template["contentVersion"]; !ok {
		results = append(results, ValidationResult{
			Severity: SeverityError,
			Field:    "contentVersion",
			Message:  "missing required field",
		})
	}

	if _, ok := template["resources"]; !ok {
		results = append(results, ValidationResult{
			Severity: SeverityError,
			Field:    "resources",
			Message:  "missing required field",
		})
	} else {
		// Validate resources array
		resources, ok := template["resources"].([]interface{})
		if ok {
			for i, res := range resources {
				resResults := v.validateResource(res, i)
				results = append(results, resResults...)
			}
		}
	}

	return results, nil
}

// validateResource validates a single resource in the template.
func (v *Validator) validateResource(res interface{}, index int) []ValidationResult {
	var results []ValidationResult

	resMap, ok := res.(map[string]interface{})
	if !ok {
		results = append(results, ValidationResult{
			Severity: SeverityError,
			Field:    fmt.Sprintf("resources[%d]", index),
			Message:  "resource must be an object",
		})
		return results
	}

	// Check for required resource fields
	if _, ok := resMap["type"]; !ok {
		results = append(results, ValidationResult{
			Severity: SeverityError,
			Field:    fmt.Sprintf("resources[%d].type", index),
			Message:  "missing required field",
		})
	}

	if _, ok := resMap["name"]; !ok {
		results = append(results, ValidationResult{
			Severity: SeverityError,
			Field:    fmt.Sprintf("resources[%d].name", index),
			Message:  "missing required field",
		})
	}

	if _, ok := resMap["apiVersion"]; !ok {
		results = append(results, ValidationResult{
			Severity: SeverityError,
			Field:    fmt.Sprintf("resources[%d].apiVersion", index),
			Message:  "missing required field",
		})
	}

	return results
}

// ValidateFile reads and validates an ARM template file.
func (v *Validator) ValidateFile(path string) ([]ValidationResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return v.ValidateTemplate(data)
}

// isValidSchemaURL checks if the schema URL is a known ARM template schema.
func isValidSchemaURL(url string) bool {
	validPrefixes := []string{
		"https://schema.management.azure.com/schemas/",
		"http://schema.management.azure.com/schemas/",
	}

	for _, prefix := range validPrefixes {
		if strings.HasPrefix(url, prefix) {
			return true
		}
	}

	return false
}
