package lint

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSeverityString_Unknown tests Severity.String() with unknown value
func TestSeverityString_Unknown(t *testing.T) {
	s := Severity(99)
	assert.Equal(t, "unknown", s.String())
}

// TestLinterCheckFile_NonExistent tests CheckFile with non-existent file
func TestLinterCheckFile_NonExistent(t *testing.T) {
	linter := NewLinter()
	_, err := linter.CheckFile("/nonexistent/file.go")
	assert.Error(t, err)
}

// TestLinterCheckFile_NonGoFile tests CheckFile with non-Go file
func TestLinterCheckFile_NonGoFile(t *testing.T) {
	tmpDir := t.TempDir()
	txtFile := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(txtFile, []byte("hello"), 0644)
	require.NoError(t, err)

	linter := NewLinter()
	results, err := linter.CheckFile(txtFile)
	require.NoError(t, err)
	assert.Nil(t, results)
}

// TestLinterCheckDirectory_SkipsTestFiles tests that test files are skipped
func TestLinterCheckDirectory_SkipsTestFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test file with issues
	testFile := filepath.Join(tmpDir, "storage_test.go")
	testContent := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var TestStorage = storage.StorageAccount{
	Location: "East US",
}
`
	err := os.WriteFile(testFile, []byte(testContent), 0644)
	require.NoError(t, err)

	linter := NewLinter()
	results, err := linter.CheckDirectory(tmpDir)
	require.NoError(t, err)
	assert.Empty(t, results, "Test files should be skipped")
}

// TestLinterCheckDirectory_SubdirError tests error handling when subdirectory has issues
func TestLinterCheckDirectory_SubdirError(t *testing.T) {
	tmpDir := t.TempDir()

	// Create subdirectory
	subDir := filepath.Join(tmpDir, "subdir")
	err := os.MkdirAll(subDir, 0755)
	require.NoError(t, err)

	// Create a valid Go file in subdirectory
	subFile := filepath.Join(subDir, "storage.go")
	content := `package subdir

import "github.com/lex00/wetwire-azure-go/resources/storage"

var Storage = storage.StorageAccount{
	Location: "eastus",
}
`
	err = os.WriteFile(subFile, []byte(content), 0644)
	require.NoError(t, err)

	linter := NewLinter()
	results, err := linter.CheckDirectory(tmpDir)
	require.NoError(t, err)
	// Should not find issues for valid location
	for _, r := range results {
		if r.Rule == "WAZ001" {
			t.Error("Should not flag valid location 'eastus'")
		}
	}
}

// TestLinterCheckFile_RuleFails tests error when a rule fails
func TestLinterCheckFile_RuleFails(t *testing.T) {
	tmpDir := t.TempDir()

	// Create invalid Go file
	badFile := filepath.Join(tmpDir, "invalid.go")
	err := os.WriteFile(badFile, []byte("not valid go code {{{"), 0644)
	require.NoError(t, err)

	linter := NewLinter()
	_, err = linter.CheckFile(badFile)
	assert.Error(t, err)
}

// TestWAZ001_ARMExpression tests WAZ001 skips ARM template expressions
func TestWAZ001_ARMExpression(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{
	Location: "[resourceGroup().location]",
}
`
	testFile := filepath.Join(tmpDir, "test.go")
	err := os.WriteFile(testFile, []byte(code), 0644)
	require.NoError(t, err)

	rule := &WAZ001{}
	results, err := rule.Check(testFile)
	require.NoError(t, err)
	assert.Empty(t, results, "ARM expressions should not trigger warnings")
}

// TestWAZ001_Description tests WAZ001.Description()
func TestWAZ001_Description(t *testing.T) {
	rule := &WAZ001{}
	desc := rule.Description()
	assert.NotEmpty(t, desc)
}

// TestWAZ002_Description tests WAZ002.Description()
func TestWAZ002_Description(t *testing.T) {
	rule := &WAZ002{}
	desc := rule.Description()
	assert.NotEmpty(t, desc)
}

// TestWAZ003_Description tests WAZ003.Description()
func TestWAZ003_Description(t *testing.T) {
	rule := &WAZ003{}
	desc := rule.Description()
	assert.NotEmpty(t, desc)
}

// TestWAZ004_Description tests WAZ004.Description()
func TestWAZ004_Description(t *testing.T) {
	rule := &WAZ004{}
	desc := rule.Description()
	assert.NotEmpty(t, desc)
}

// TestWAZ005_Description tests WAZ005.Description()
func TestWAZ005_Description(t *testing.T) {
	rule := &WAZ005{}
	desc := rule.Description()
	assert.NotEmpty(t, desc)
}

// TestWAZ001_InvalidFile tests WAZ001 with invalid Go file
func TestWAZ001_InvalidFile(t *testing.T) {
	tmpDir := t.TempDir()
	badFile := filepath.Join(tmpDir, "bad.go")
	err := os.WriteFile(badFile, []byte("invalid go {{{"), 0644)
	require.NoError(t, err)

	rule := &WAZ001{}
	_, err = rule.Check(badFile)
	assert.Error(t, err)
}

// TestWAZ002_InvalidFile tests WAZ002 with invalid Go file
func TestWAZ002_InvalidFile(t *testing.T) {
	tmpDir := t.TempDir()
	badFile := filepath.Join(tmpDir, "bad.go")
	err := os.WriteFile(badFile, []byte("invalid go {{{"), 0644)
	require.NoError(t, err)

	rule := &WAZ002{}
	_, err = rule.Check(badFile)
	assert.Error(t, err)
}

// TestWAZ003_InvalidFile tests WAZ003 with invalid Go file
func TestWAZ003_InvalidFile(t *testing.T) {
	tmpDir := t.TempDir()
	badFile := filepath.Join(tmpDir, "bad.go")
	err := os.WriteFile(badFile, []byte("invalid go {{{"), 0644)
	require.NoError(t, err)

	rule := &WAZ003{}
	_, err = rule.Check(badFile)
	assert.Error(t, err)
}

// TestWAZ004_InvalidFile tests WAZ004 with invalid Go file
func TestWAZ004_InvalidFile(t *testing.T) {
	tmpDir := t.TempDir()
	badFile := filepath.Join(tmpDir, "bad.go")
	err := os.WriteFile(badFile, []byte("invalid go {{{"), 0644)
	require.NoError(t, err)

	rule := &WAZ004{}
	_, err = rule.Check(badFile)
	assert.Error(t, err)
}

// TestWAZ005_InvalidFile tests WAZ005 with invalid Go file
func TestWAZ005_InvalidFile(t *testing.T) {
	tmpDir := t.TempDir()
	badFile := filepath.Join(tmpDir, "bad.go")
	err := os.WriteFile(badFile, []byte("invalid go {{{"), 0644)
	require.NoError(t, err)

	rule := &WAZ005{}
	_, err = rule.Check(badFile)
	assert.Error(t, err)
}

// TestWAZ004_SkipsBlankIdentifier tests WAZ004 skips blank identifiers
func TestWAZ004_SkipsBlankIdentifier(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var _ = storage.StorageAccount{Name: "storage1"}
var _ = storage.StorageAccount{Name: "storage2"}
`
	testFile := filepath.Join(tmpDir, "test.go")
	err := os.WriteFile(testFile, []byte(code), 0644)
	require.NoError(t, err)

	rule := &WAZ004{}
	results, err := rule.Check(testFile)
	require.NoError(t, err)
	assert.Empty(t, results, "Blank identifiers should not be checked for duplicates")
}

// TestWAZ005_SelfReference tests WAZ005 detects self-reference
func TestWAZ005_SelfReference(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var Storage = storage.StorageAccount{
	Name: Storage.Name,
}
`
	testFile := filepath.Join(tmpDir, "test.go")
	err := os.WriteFile(testFile, []byte(code), 0644)
	require.NoError(t, err)

	rule := &WAZ005{}
	results, err := rule.Check(testFile)
	require.NoError(t, err)
	assert.NotEmpty(t, results, "Self-reference should be detected")
}

// TestWAZ005_TransitiveDependency tests WAZ005 detects transitive circular dependency
func TestWAZ005_TransitiveDependency(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var A = storage.StorageAccount{Name: B.Name}
var B = storage.StorageAccount{Name: C.Name}
var C = storage.StorageAccount{Name: A.Name}
`
	testFile := filepath.Join(tmpDir, "test.go")
	err := os.WriteFile(testFile, []byte(code), 0644)
	require.NoError(t, err)

	rule := &WAZ005{}
	results, err := rule.Check(testFile)
	require.NoError(t, err)
	assert.NotEmpty(t, results, "Transitive circular dependency should be detected")
}

// TestWAZ005_NonExistentDependency tests WAZ005 handles non-existent dependencies
func TestWAZ005_NonExistentDependency(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var Storage = storage.StorageAccount{
	Name: NonExistentVar.Name,
}
`
	testFile := filepath.Join(tmpDir, "test.go")
	err := os.WriteFile(testFile, []byte(code), 0644)
	require.NoError(t, err)

	rule := &WAZ005{}
	_, err = rule.Check(testFile)
	require.NoError(t, err)
	// Should not crash on non-existent dependency
}

// TestWAZ005_ExtractDependencies_AllTypes tests dependency extraction for various expression types
func TestWAZ005_ExtractDependencies_AllTypes(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var dep1 = "test"
var dep2 = "test2"

var Storage = storage.StorageAccount{
	Name:     dep1 + dep2,
	Location: func() string { return dep1 }(),
}
`
	testFile := filepath.Join(tmpDir, "test.go")
	err := os.WriteFile(testFile, []byte(code), 0644)
	require.NoError(t, err)

	rule := &WAZ005{}
	results, err := rule.Check(testFile)
	require.NoError(t, err)
	// No circular dependency, should pass
	assert.Empty(t, results)
}

// TestWAZ003_ShallowNesting tests WAZ003 allows shallow nesting
func TestWAZ003_ShallowNesting(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{
	Properties: &storage.StorageAccountProperties{
		MinimumTlsVersion: "TLS1_2",
	},
}
`
	testFile := filepath.Join(tmpDir, "test.go")
	err := os.WriteFile(testFile, []byte(code), 0644)
	require.NoError(t, err)

	rule := &WAZ003{}
	results, err := rule.Check(testFile)
	require.NoError(t, err)
	assert.Empty(t, results, "Shallow nesting should not trigger warning")
}

// TestWAZ002_NoResourceIdCalls tests WAZ002 with no resourceId calls
func TestWAZ002_NoResourceIdCalls(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var Storage = storage.StorageAccount{
	Name:     "mystorage",
	Location: "eastus",
}
`
	testFile := filepath.Join(tmpDir, "test.go")
	err := os.WriteFile(testFile, []byte(code), 0644)
	require.NoError(t, err)

	rule := &WAZ002{}
	results, err := rule.Check(testFile)
	require.NoError(t, err)
	assert.Empty(t, results)
}

// TestHasUpperCase tests the hasUpperCase helper function
func TestHasUpperCase(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"eastus", false},
		{"East US", true},
		{"EASTUS", true},
		{"", false},
		{"123", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := hasUpperCase(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestIsBuiltinType_Linter tests isBuiltinType in linter package
func TestIsBuiltinType_Linter(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"string", true},
		{"int", true},
		{"bool", true},
		{"MyType", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isBuiltinType(tt.name)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestIsKeyword tests isKeyword helper function
func TestIsKeyword(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"if", true},
		{"for", true},
		{"func", true},
		{"var", true},
		{"return", true},
		{"myVar", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isKeyword(tt.name)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestLintResultFields tests LintResult fields
func TestLintResultFields(t *testing.T) {
	result := LintResult{
		Rule:     "WAZ001",
		File:     "/path/to/file.go",
		Line:     42,
		Message:  "Test message",
		Severity: SeverityError,
	}

	assert.Equal(t, "WAZ001", result.Rule)
	assert.Equal(t, "/path/to/file.go", result.File)
	assert.Equal(t, 42, result.Line)
	assert.Equal(t, "Test message", result.Message)
	assert.Equal(t, SeverityError, result.Severity)
}

// TestAllRules_RuleInterfaces tests all rules implement Rule interface properly
func TestAllRules_RuleInterfaces(t *testing.T) {
	rules := AllRules()

	for _, rule := range rules {
		t.Run(rule.ID(), func(t *testing.T) {
			// Test ID is not empty
			assert.NotEmpty(t, rule.ID())

			// Test Description is not empty
			assert.NotEmpty(t, rule.Description())

			// Test Severity is valid
			sev := rule.Severity()
			assert.True(t, sev == SeverityError || sev == SeverityWarning || sev == SeverityInfo)
		})
	}
}

// TestWAZ005_SkipsBlankIdentifier tests WAZ005 skips blank identifiers
func TestWAZ005_SkipsBlankIdentifier(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var _ = storage.StorageAccount{Name: "storage1"}
`
	testFile := filepath.Join(tmpDir, "test.go")
	err := os.WriteFile(testFile, []byte(code), 0644)
	require.NoError(t, err)

	rule := &WAZ005{}
	results, err := rule.Check(testFile)
	require.NoError(t, err)
	assert.Empty(t, results)
}

// TestWAZ003_UnaryExprNesting tests WAZ003 with unary expression nesting
func TestWAZ003_UnaryExprNesting(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{
	Properties: &storage.StorageAccountProperties{
		NetworkRuleSet: &storage.NetworkRuleSet{
			DefaultAction: "Deny",
			IPRules: []storage.IPRule{
				{Value: "10.0.0.1"},
			},
		},
	},
}
`
	testFile := filepath.Join(tmpDir, "test.go")
	err := os.WriteFile(testFile, []byte(code), 0644)
	require.NoError(t, err)

	rule := &WAZ003{}
	results, err := rule.Check(testFile)
	require.NoError(t, err)
	// Should detect deep nesting
	assert.NotEmpty(t, results)
}

// TestWAZ005_ExtractDependencies_ParenExpr tests dependency extraction with parenthesized expressions
func TestWAZ005_ExtractDependencies_ParenExpr(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var dep1 = "test"

var Storage = storage.StorageAccount{
	Name: (dep1),
}
`
	testFile := filepath.Join(tmpDir, "test.go")
	err := os.WriteFile(testFile, []byte(code), 0644)
	require.NoError(t, err)

	rule := &WAZ005{}
	results, err := rule.Check(testFile)
	require.NoError(t, err)
	assert.Empty(t, results) // No circular dependency
}
