package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHandleBuild(t *testing.T) {
	// Create test directory with a Go file containing a resource
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "main.go")

	content := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{
	Name:     "mystorage",
	Location: "eastus",
}

func main() {}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Test build with valid path
	result, err := handleBuild(context.Background(), map[string]any{
		"path": tmpDir,
	})

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if !strings.Contains(result, "Successfully built ARM template") {
		t.Errorf("expected success message, got: %s", result)
	}

	if !strings.Contains(result, "Microsoft.Storage/storageAccounts") {
		t.Errorf("expected storage account in template, got: %s", result)
	}
}

func TestHandleBuildNoResources(t *testing.T) {
	// Create empty test directory
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "main.go")

	content := `package main

func main() {}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	result, err := handleBuild(context.Background(), map[string]any{
		"path": tmpDir,
	})

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if !strings.Contains(result, "No Azure resources found") {
		t.Errorf("expected 'no resources' message, got: %s", result)
	}
}

func TestHandleBuildInvalidPath(t *testing.T) {
	_, err := handleBuild(context.Background(), map[string]any{
		"path": "/nonexistent/path",
	})

	if err == nil {
		t.Fatal("expected error for nonexistent path")
	}

	if !strings.Contains(err.Error(), "path not found") {
		t.Errorf("expected 'path not found' error, got: %v", err)
	}
}

func TestHandleLint(t *testing.T) {
	// Create test directory with a Go file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "main.go")

	content := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{
	Name:     "mystorage",
	Location: "eastus",
}

func main() {}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	result, err := handleLint(context.Background(), map[string]any{
		"path": tmpDir,
	})

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Should either have no issues or report issues - both are valid lint results
	if result == "" {
		t.Error("expected non-empty result")
	}
}

func TestHandleLintInvalidPath(t *testing.T) {
	_, err := handleLint(context.Background(), map[string]any{
		"path": "/nonexistent/path",
	})

	if err == nil {
		t.Fatal("expected error for nonexistent path")
	}

	if !strings.Contains(err.Error(), "path not found") {
		t.Errorf("expected 'path not found' error, got: %v", err)
	}
}

func TestHandleImport(t *testing.T) {
	// Create test ARM template file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "template.json")

	content := `{
  "$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "resources": [
    {
      "type": "Microsoft.Storage/storageAccounts",
      "apiVersion": "2021-02-01",
      "name": "mystorage",
      "location": "eastus",
      "sku": {"name": "Standard_LRS"},
      "kind": "StorageV2"
    }
  ]
}`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	result, err := handleImport(context.Background(), map[string]any{
		"file":    testFile,
		"package": "myinfra",
	})

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if !strings.Contains(result, "Successfully imported ARM template") {
		t.Errorf("expected success message, got: %s", result)
	}

	if !strings.Contains(result, "package myinfra") {
		t.Errorf("expected package name in result, got: %s", result)
	}
}

func TestHandleImportMissingFile(t *testing.T) {
	_, err := handleImport(context.Background(), map[string]any{
		"file": "/nonexistent/template.json",
	})

	if err == nil {
		t.Fatal("expected error for missing file")
	}

	if !strings.Contains(err.Error(), "error reading file") {
		t.Errorf("expected 'error reading file' error, got: %v", err)
	}
}

func TestHandleImportMissingFileArg(t *testing.T) {
	_, err := handleImport(context.Background(), map[string]any{})

	if err == nil {
		t.Fatal("expected error for missing file argument")
	}

	if !strings.Contains(err.Error(), "file argument is required") {
		t.Errorf("expected 'file argument is required' error, got: %v", err)
	}
}

func TestHandleImportInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "template.json")

	if err := os.WriteFile(testFile, []byte("not valid json"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, err := handleImport(context.Background(), map[string]any{
		"file": testFile,
	})

	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}

	if !strings.Contains(err.Error(), "error parsing ARM template") {
		t.Errorf("expected 'error parsing' error, got: %v", err)
	}
}

func TestDefaultArguments(t *testing.T) {
	// Test that handlers use defaults when args are empty

	// Build with empty args should use "."
	_, err := handleBuild(context.Background(), map[string]any{})
	// This will fail because current dir likely has no resources,
	// but should not error on path resolution
	if err != nil && strings.Contains(err.Error(), "resolving path") {
		t.Errorf("should handle empty path arg: %v", err)
	}

	// Lint with empty args should use "."
	_, err = handleLint(context.Background(), map[string]any{})
	if err != nil && strings.Contains(err.Error(), "resolving path") {
		t.Errorf("should handle empty path arg: %v", err)
	}

	// Import should require file
	_, err = handleImport(context.Background(), map[string]any{})
	if err == nil {
		t.Error("import should require file argument")
	}
}
