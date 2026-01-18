package domain

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	coredomain "github.com/lex00/wetwire-core-go/domain"
)

// TestDomainInterface verifies that AzureDomain implements the Domain interface at compile time
func TestDomainInterface(t *testing.T) {
	var _ coredomain.Domain = (*AzureDomain)(nil)
}

// TestListerInterface verifies that AzureDomain implements the ListerDomain interface at compile time
func TestListerInterface(t *testing.T) {
	var _ coredomain.ListerDomain = (*AzureDomain)(nil)
}

// TestGrapherInterface verifies that AzureDomain implements the GrapherDomain interface at compile time
func TestGrapherInterface(t *testing.T) {
	var _ coredomain.GrapherDomain = (*AzureDomain)(nil)
}

// TestBuilderInterface verifies that azureBuilder implements the Builder interface at compile time
func TestBuilderInterface(t *testing.T) {
	var _ coredomain.Builder = (*azureBuilder)(nil)
}

// TestLinterInterface verifies that azureLinter implements the Linter interface at compile time
func TestLinterInterface(t *testing.T) {
	var _ coredomain.Linter = (*azureLinter)(nil)
}

// TestInitializerInterface verifies that azureInitializer implements the Initializer interface at compile time
func TestInitializerInterface(t *testing.T) {
	var _ coredomain.Initializer = (*azureInitializer)(nil)
}

// TestValidatorInterface verifies that azureValidator implements the Validator interface at compile time
func TestValidatorInterface(t *testing.T) {
	var _ coredomain.Validator = (*azureValidator)(nil)
}

// TestListerImplInterface verifies that azureLister implements the Lister interface at compile time
func TestListerImplInterface(t *testing.T) {
	var _ coredomain.Lister = (*azureLister)(nil)
}

// TestGrapherImplInterface verifies that azureGrapher implements the Grapher interface at compile time
func TestGrapherImplInterface(t *testing.T) {
	var _ coredomain.Grapher = (*azureGrapher)(nil)
}

// TestLintOpts_Disable tests that LintOpts.Disable is respected
func TestLintOpts_Disable(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a file that triggers WAZ001 (invalid location format)
	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{
	Name:     "mystorageaccount",
	Location: "East US",
}
`
	testFile := filepath.Join(tmpDir, "test.go")
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	domain := &AzureDomain{}
	linter := domain.Linter()
	ctx := NewContext(context.Background(), tmpDir)

	// First, verify WAZ001 is triggered without disabling
	result, err := linter.Lint(ctx, testFile, LintOpts{})
	if err != nil {
		t.Fatalf("Lint() error: %v", err)
	}

	// Check that WAZ001 is found
	hasWAZ001 := false
	for _, e := range result.Errors {
		if e.Code == "WAZ001" {
			hasWAZ001 = true
			break
		}
	}
	if !hasWAZ001 {
		t.Error("Expected WAZ001 to be triggered for 'East US' location")
	}

	// Now test with WAZ001 disabled
	resultDisabled, err := linter.Lint(ctx, testFile, LintOpts{
		Disable: []string{"WAZ001"},
	})
	if err != nil {
		t.Fatalf("Lint() with Disable error: %v", err)
	}

	// Check that WAZ001 is NOT found
	for _, e := range resultDisabled.Errors {
		if e.Code == "WAZ001" {
			t.Error("WAZ001 should be disabled but was triggered")
		}
	}
}

// TestLintOpts_Fix tests that LintOpts.Fix is accepted and produces appropriate message
func TestLintOpts_Fix(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a file that triggers WAZ001
	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{
	Location: "East US",
}
`
	testFile := filepath.Join(tmpDir, "test.go")
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	domain := &AzureDomain{}
	linter := domain.Linter()
	ctx := NewContext(context.Background(), tmpDir)

	// Test with Fix enabled
	result, err := linter.Lint(ctx, testFile, LintOpts{
		Fix: true,
	})
	if err != nil {
		t.Fatalf("Lint() with Fix error: %v", err)
	}

	// Result message should indicate auto-fix not implemented
	if result.Message == "" {
		t.Error("Expected result message to be set")
	}

	// When there are errors and Fix is true, message should mention auto-fix
	if len(result.Errors) > 0 {
		expected := "lint issues found (auto-fix not yet implemented for these issues)"
		if result.Message != expected {
			t.Errorf("Expected message %q, got %q", expected, result.Message)
		}
	}
}

// TestLintOpts_DisableMultiple tests disabling multiple rules
func TestLintOpts_DisableMultiple(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a file that would trigger multiple rules
	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{
	Location: "East US",
}

var MyStorage = storage.StorageAccount{
	Location: "West US",
}
`
	testFile := filepath.Join(tmpDir, "test.go")
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	domain := &AzureDomain{}
	linter := domain.Linter()
	ctx := NewContext(context.Background(), tmpDir)

	// Test with both WAZ001 and WAZ004 disabled
	result, err := linter.Lint(ctx, testFile, LintOpts{
		Disable: []string{"WAZ001", "WAZ004"},
	})
	if err != nil {
		t.Fatalf("Lint() with multiple Disable error: %v", err)
	}

	// Check that neither WAZ001 nor WAZ004 is found
	for _, e := range result.Errors {
		if e.Code == "WAZ001" || e.Code == "WAZ004" {
			t.Errorf("Rule %s should be disabled but was triggered", e.Code)
		}
	}
}
