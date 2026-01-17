package lint

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLintResult(t *testing.T) {
	result := LintResult{
		Rule:     "WAZ001",
		File:     "test.go",
		Line:     10,
		Message:  "Use location constants",
		Severity: SeverityWarning,
	}

	if result.Rule != "WAZ001" {
		t.Errorf("expected Rule WAZ001, got %s", result.Rule)
	}
	if result.Severity != SeverityWarning {
		t.Errorf("expected SeverityWarning, got %s", result.Severity)
	}
}

func TestLintResultString(t *testing.T) {
	result := LintResult{
		Rule:     "WAZ001",
		File:     "/path/to/test.go",
		Line:     10,
		Message:  "Use location constants",
		Severity: SeverityWarning,
	}

	str := result.String()
	if str == "" {
		t.Error("String() should not return empty string")
	}
	// Should contain file, line, rule, and message
	if !contains(str, "test.go") || !contains(str, "10") || !contains(str, "WAZ001") {
		t.Errorf("String() missing expected components: %s", str)
	}
}

func TestNewLinter(t *testing.T) {
	linter := NewLinter()
	if linter == nil {
		t.Fatal("NewLinter() returned nil")
	}
	if len(linter.rules) == 0 {
		t.Error("NewLinter() should register default rules")
	}
}

func TestLinterAddRule(t *testing.T) {
	linter := NewLinter()
	initialCount := len(linter.rules)

	// Create a mock rule
	mockRule := &mockRule{id: "TEST001", description: "Test rule"}
	linter.AddRule(mockRule)

	if len(linter.rules) != initialCount+1 {
		t.Errorf("expected %d rules, got %d", initialCount+1, len(linter.rules))
	}
}

func TestLinterCheckFile(t *testing.T) {
	// Create a temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")
	testContent := `package main

import (
	"github.com/lex00/wetwire-azure-go/resources/storage"
)

var MyStorage = storage.StorageAccount{
	Name:     "mystorageaccount",
	Location: "East US", // Should trigger WAZ001
}
`
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatal(err)
	}

	linter := NewLinter()
	results, err := linter.CheckFile(testFile)
	if err != nil {
		t.Fatalf("CheckFile() error: %v", err)
	}

	// Should find at least one issue (location string literal)
	if len(results) == 0 {
		t.Error("expected at least one lint result")
	}
}

func TestLinterCheckDirectory(t *testing.T) {
	// Create a temporary directory with test files
	tmpDir := t.TempDir()

	testFile1 := filepath.Join(tmpDir, "file1.go")
	testContent1 := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var Storage1 = storage.StorageAccount{
	Name:     "storage1",
	Location: "East US",
}
`
	if err := os.WriteFile(testFile1, []byte(testContent1), 0644); err != nil {
		t.Fatal(err)
	}

	testFile2 := filepath.Join(tmpDir, "file2.go")
	testContent2 := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var Storage2 = storage.StorageAccount{
	Name:     "storage2",
	Location: "West US",
}
`
	if err := os.WriteFile(testFile2, []byte(testContent2), 0644); err != nil {
		t.Fatal(err)
	}

	linter := NewLinter()
	results, err := linter.CheckDirectory(tmpDir)
	if err != nil {
		t.Fatalf("CheckDirectory() error: %v", err)
	}

	// Should find issues in both files
	if len(results) == 0 {
		t.Error("expected lint results from both files")
	}
}

func TestLinterCheckDirectoryNonExistent(t *testing.T) {
	linter := NewLinter()
	_, err := linter.CheckDirectory("/nonexistent/path")
	if err == nil {
		t.Error("expected error for nonexistent directory")
	}
}

func TestSeverityString(t *testing.T) {
	tests := []struct {
		severity Severity
		expected string
	}{
		{SeverityError, "error"},
		{SeverityWarning, "warning"},
		{SeverityInfo, "info"},
	}

	for _, tt := range tests {
		if tt.severity.String() != tt.expected {
			t.Errorf("severity.String() = %s, want %s", tt.severity.String(), tt.expected)
		}
	}
}

// mockRule is a test implementation of the Rule interface
type mockRule struct {
	id          string
	description string
	results     []LintResult
}

func (m *mockRule) ID() string {
	return m.id
}

func (m *mockRule) Description() string {
	return m.description
}

func (m *mockRule) Severity() Severity {
	return SeverityWarning
}

func (m *mockRule) Check(file string) ([]LintResult, error) {
	return m.results, nil
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
