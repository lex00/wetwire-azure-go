package validator

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateTemplate_ValidTemplate(t *testing.T) {
	template := map[string]interface{}{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"resources":      []interface{}{},
	}

	jsonBytes, _ := json.Marshal(template)
	validator := NewValidator()
	results, err := validator.ValidateTemplate(jsonBytes)

	if err != nil {
		t.Fatalf("ValidateTemplate failed: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("Expected no validation errors, got %d", len(results))
	}
}

func TestValidateTemplate_InvalidJSON(t *testing.T) {
	validator := NewValidator()
	_, err := validator.ValidateTemplate([]byte("{invalid"))

	if err == nil {
		t.Fatal("Expected error for invalid JSON")
	}
}

func TestValidateTemplate_MissingSchema(t *testing.T) {
	template := map[string]interface{}{
		"contentVersion": "1.0.0.0",
		"resources":      []interface{}{},
	}

	jsonBytes, _ := json.Marshal(template)
	validator := NewValidator()
	results, _ := validator.ValidateTemplate(jsonBytes)

	if len(results) == 0 {
		t.Error("Expected validation errors for missing $schema")
	}

	found := false
	for _, r := range results {
		if r.Field == "$schema" && r.Severity == SeverityError {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected error for missing $schema field")
	}
}

func TestValidateTemplate_MissingContentVersion(t *testing.T) {
	template := map[string]interface{}{
		"$schema":   "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"resources": []interface{}{},
	}

	jsonBytes, _ := json.Marshal(template)
	validator := NewValidator()
	results, _ := validator.ValidateTemplate(jsonBytes)

	found := false
	for _, r := range results {
		if r.Field == "contentVersion" && r.Severity == SeverityError {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected error for missing contentVersion field")
	}
}

func TestValidateTemplate_MissingResources(t *testing.T) {
	template := map[string]interface{}{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
	}

	jsonBytes, _ := json.Marshal(template)
	validator := NewValidator()
	results, _ := validator.ValidateTemplate(jsonBytes)

	found := false
	for _, r := range results {
		if r.Field == "resources" && r.Severity == SeverityError {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected error for missing resources field")
	}
}

func TestValidateTemplate_InvalidSchemaURL(t *testing.T) {
	template := map[string]interface{}{
		"$schema":        "https://example.com/invalid-schema",
		"contentVersion": "1.0.0.0",
		"resources":      []interface{}{},
	}

	jsonBytes, _ := json.Marshal(template)
	validator := NewValidator()
	results, _ := validator.ValidateTemplate(jsonBytes)

	found := false
	for _, r := range results {
		if r.Field == "$schema" && r.Severity == SeverityWarning {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected warning for invalid schema URL")
	}
}

func TestValidateTemplate_ResourceWithoutType(t *testing.T) {
	template := map[string]interface{}{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"resources": []interface{}{
			map[string]interface{}{
				"name":       "myResource",
				"apiVersion": "2021-01-01",
			},
		},
	}

	jsonBytes, _ := json.Marshal(template)
	validator := NewValidator()
	results, _ := validator.ValidateTemplate(jsonBytes)

	found := false
	for _, r := range results {
		if strings.Contains(r.Field, "type") && r.Severity == SeverityError {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected error for resource missing type")
	}
}

func TestValidateFile_ValidFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "template.json")

	template := map[string]interface{}{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"resources":      []interface{}{},
	}

	jsonBytes, _ := json.MarshalIndent(template, "", "  ")
	if err := os.WriteFile(tmpFile, jsonBytes, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	validator := NewValidator()
	results, err := validator.ValidateFile(tmpFile)

	if err != nil {
		t.Fatalf("ValidateFile failed: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("Expected no validation errors, got %d", len(results))
	}
}

func TestValidateFile_NonExistentFile(t *testing.T) {
	validator := NewValidator()
	_, err := validator.ValidateFile("/nonexistent/file.json")

	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestValidationResult_String(t *testing.T) {
	result := ValidationResult{
		Severity: SeverityError,
		Message:  "Test message",
		Field:    "testField",
	}

	str := result.String()
	if !strings.Contains(str, "error") || !strings.Contains(str, "Test message") {
		t.Error("String representation missing expected content")
	}
}

func TestValidationResult_StringNoField(t *testing.T) {
	result := ValidationResult{
		Severity: SeverityWarning,
		Message:  "Test message",
	}

	str := result.String()
	if !strings.Contains(str, "warning") || !strings.Contains(str, "Test message") {
		t.Error("String representation missing expected content")
	}
	if strings.Contains(str, ":") && strings.Count(str, ":") > 1 {
		t.Error("String should not have field separator when field is empty")
	}
}

func TestSeverity_String(t *testing.T) {
	tests := []struct {
		severity Severity
		expected string
	}{
		{SeverityError, "error"},
		{SeverityWarning, "warning"},
		{SeverityInfo, "info"},
		{Severity(99), "unknown"},
	}

	for _, tt := range tests {
		if got := tt.severity.String(); got != tt.expected {
			t.Errorf("Severity.String() = %v, want %v", got, tt.expected)
		}
	}
}
