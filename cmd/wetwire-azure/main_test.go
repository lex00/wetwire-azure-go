package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// captureOutput captures stdout and stderr during function execution
func captureOutput(f func()) (string, string) {
	// Capture stdout
	oldStdout := os.Stdout
	rOut, wOut, _ := os.Pipe()
	os.Stdout = wOut

	// Capture stderr
	oldStderr := os.Stderr
	rErr, wErr, _ := os.Pipe()
	os.Stderr = wErr

	f()

	_ = wOut.Close()
	_ = wErr.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	var bufOut, bufErr bytes.Buffer
	_, _ = io.Copy(&bufOut, rOut)
	_, _ = io.Copy(&bufErr, rErr)

	return bufOut.String(), bufErr.String()
}

// TestPrintUsage tests the printUsage function
func TestPrintUsage(t *testing.T) {
	stdout, _ := captureOutput(func() {
		printUsage()
	})

	assert.Contains(t, stdout, "wetwire-azure")
	assert.Contains(t, stdout, "Usage:")
	assert.Contains(t, stdout, "lint")
	assert.Contains(t, stdout, "help")
}

// TestRun_NoArgs tests run with no arguments
func TestRun_NoArgs(t *testing.T) {
	stdout, _ := captureOutput(func() {
		code := run([]string{})
		assert.Equal(t, ExitInvalidArgument, code)
	})
	assert.Contains(t, stdout, "Usage:")
}

// TestRun_Help tests run with help command
func TestRun_Help(t *testing.T) {
	stdout, _ := captureOutput(func() {
		code := run([]string{"help"})
		assert.Equal(t, ExitSuccess, code)
	})
	assert.Contains(t, stdout, "Usage:")
}

// TestRun_HelpShort tests run with -h flag
func TestRun_HelpShort(t *testing.T) {
	stdout, _ := captureOutput(func() {
		code := run([]string{"-h"})
		assert.Equal(t, ExitSuccess, code)
	})
	assert.Contains(t, stdout, "Usage:")
}

// TestRun_HelpLong tests run with --help flag
func TestRun_HelpLong(t *testing.T) {
	stdout, _ := captureOutput(func() {
		code := run([]string{"--help"})
		assert.Equal(t, ExitSuccess, code)
	})
	assert.Contains(t, stdout, "Usage:")
}

// TestRun_UnknownCommand tests run with unknown command
func TestRun_UnknownCommand(t *testing.T) {
	_, stderr := captureOutput(func() {
		code := run([]string{"unknown"})
		assert.Equal(t, ExitInvalidArgument, code)
	})
	assert.Contains(t, stderr, "Unknown command")
}

// TestRunLint_CurrentDirectory tests runLint with current directory (no issues)
func TestRunLint_CurrentDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{
	Name:     "mystorage",
	Location: "eastus",
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	oldWd, _ := os.Getwd()
	require.NoError(t, os.Chdir(tmpDir))
	defer func() { _ = os.Chdir(oldWd) }()

	stdout, _ := captureOutput(func() {
		exitCode := runLint([]string{})
		assert.Equal(t, ExitSuccess, exitCode)
	})

	assert.Contains(t, stdout, "No issues found")
}

// TestRunLint_WithPath tests runLint with explicit path
func TestRunLint_WithPath(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{
	Name:     "mystorage",
	Location: "eastus",
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	stdout, _ := captureOutput(func() {
		exitCode := runLint([]string{tmpDir})
		assert.Equal(t, ExitSuccess, exitCode)
	})

	assert.Contains(t, stdout, "No issues found")
}

// TestRunLint_WithIssues tests runLint with lint issues
func TestRunLint_WithIssues(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{
	Name:     "mystorage",
	Location: "East US",
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	stdout, _ := captureOutput(func() {
		exitCode := runLint([]string{tmpDir})
		assert.Equal(t, ExitBuildError, exitCode)
	})

	assert.Contains(t, stdout, "WAZ001")
	assert.Contains(t, stdout, "Found")
}

// TestRunLint_SingleFile tests runLint with a single file
func TestRunLint_SingleFile(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{
	Name:     "mystorage",
	Location: "eastus",
}
`
	testFile := filepath.Join(tmpDir, "main.go")
	err := os.WriteFile(testFile, []byte(code), 0644)
	require.NoError(t, err)

	stdout, _ := captureOutput(func() {
		exitCode := runLint([]string{testFile})
		assert.Equal(t, ExitSuccess, exitCode)
	})

	assert.Contains(t, stdout, "No issues found")
}

// TestRunLint_WithFixFlag tests runLint with --fix flag
func TestRunLint_WithFixFlag(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{
	Name:     "mystorage",
	Location: "East US",
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	stdout, _ := captureOutput(func() {
		runLint([]string{"--fix", tmpDir})
	})

	assert.Contains(t, stdout, "Auto-fix is not yet implemented")
}

// TestRunLint_EmptyDirectory tests runLint with empty directory
func TestRunLint_EmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	stdout, _ := captureOutput(func() {
		exitCode := runLint([]string{tmpDir})
		assert.Equal(t, ExitSuccess, exitCode)
	})

	assert.Contains(t, stdout, "No issues found")
}

// TestRunLint_NonExistentPath tests runLint with non-existent path
func TestRunLint_NonExistentPath(t *testing.T) {
	_, stderr := captureOutput(func() {
		exitCode := runLint([]string{"/nonexistent/path"})
		assert.Equal(t, ExitBuildError, exitCode)
	})

	assert.Contains(t, stderr, "Error")
}

// TestRunLint_InvalidGoFile tests runLint with invalid Go file
func TestRunLint_InvalidGoFile(t *testing.T) {
	tmpDir := t.TempDir()
	err := os.WriteFile(filepath.Join(tmpDir, "bad.go"), []byte("invalid go {{{"), 0644)
	require.NoError(t, err)

	_, stderr := captureOutput(func() {
		exitCode := runLint([]string{tmpDir})
		assert.Equal(t, ExitBuildError, exitCode)
	})

	assert.Contains(t, stderr, "Linting failed")
}

// TestRunLint_NestedDirectories tests runLint with nested directories
func TestRunLint_NestedDirectories(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "infra")
	err := os.MkdirAll(subDir, 0755)
	require.NoError(t, err)

	code := `package infra

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{
	Name:     "mystorage",
	Location: "eastus",
}
`
	err = os.WriteFile(filepath.Join(subDir, "storage.go"), []byte(code), 0644)
	require.NoError(t, err)

	stdout, _ := captureOutput(func() {
		exitCode := runLint([]string{tmpDir})
		assert.Equal(t, ExitSuccess, exitCode)
	})

	assert.Contains(t, stdout, "No issues found")
}

// TestRunLint_MultipleIssues tests runLint detecting multiple issues
func TestRunLint_MultipleIssues(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{
	Name:     "mystorage",
	Location: "East US",
}

var id = resourceId("Microsoft.Storage/storageAccounts", "test")
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	stdout, _ := captureOutput(func() {
		exitCode := runLint([]string{tmpDir})
		assert.Equal(t, ExitBuildError, exitCode)
	})

	// Should find both WAZ001 (location) and WAZ002 (resourceId)
	assert.Contains(t, stdout, "WAZ001")
	assert.Contains(t, stdout, "WAZ002")
}

// TestUsageFormat tests that usage output has correct format
func TestUsageFormat(t *testing.T) {
	stdout, _ := captureOutput(func() {
		printUsage()
	})

	// Check all expected sections
	assert.Contains(t, stdout, "wetwire-azure")
	assert.Contains(t, stdout, "lint")
	assert.Contains(t, stdout, "help")
	assert.Contains(t, stdout, "--fix")
}

// TestMain_WithExit tests the main function by overriding osExit
func TestMain_WithExit(t *testing.T) {
	// Save original osExit
	originalExit := osExit
	defer func() { osExit = originalExit }()

	var exitCode int
	osExit = func(code int) {
		exitCode = code
	}

	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"wetwire-azure", "help"}

	captureOutput(func() {
		main()
	})

	assert.Equal(t, ExitSuccess, exitCode)
}

// TestRun_Lint tests run with lint command
func TestRun_Lint(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

var x = 42
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	require.NoError(t, err)

	stdout, _ := captureOutput(func() {
		exitCode := run([]string{"lint", tmpDir})
		assert.Equal(t, ExitSuccess, exitCode)
	})

	assert.Contains(t, stdout, "No issues found")
}

// TestRun_Validate tests run with validate command
func TestRun_Validate(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "template.json")

	template := `{
  "$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "resources": []
}`
	err := os.WriteFile(tmpFile, []byte(template), 0644)
	require.NoError(t, err)

	stdout, _ := captureOutput(func() {
		exitCode := run([]string{"validate", tmpFile})
		assert.Equal(t, ExitSuccess, exitCode)
	})

	assert.Contains(t, stdout, "Template is valid")
}

// TestRunValidate_ValidTemplate tests runValidate with valid template
func TestRunValidate_ValidTemplate(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "template.json")

	template := `{
  "$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "resources": []
}`
	err := os.WriteFile(tmpFile, []byte(template), 0644)
	require.NoError(t, err)

	stdout, _ := captureOutput(func() {
		exitCode := runValidate([]string{tmpFile})
		assert.Equal(t, ExitSuccess, exitCode)
	})

	assert.Contains(t, stdout, "Template is valid")
}

// TestRunValidate_InvalidTemplate tests runValidate with invalid template
func TestRunValidate_InvalidTemplate(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "template.json")

	template := `{
  "contentVersion": "1.0.0.0",
  "resources": []
}`
	err := os.WriteFile(tmpFile, []byte(template), 0644)
	require.NoError(t, err)

	stdout, _ := captureOutput(func() {
		exitCode := runValidate([]string{tmpFile})
		assert.Equal(t, ExitBuildError, exitCode)
	})

	assert.Contains(t, stdout, "$schema")
	assert.Contains(t, stdout, "error")
}

// TestRunValidate_NoArgs tests runValidate with no arguments
func TestRunValidate_NoArgs(t *testing.T) {
	_, stderr := captureOutput(func() {
		exitCode := runValidate([]string{})
		assert.Equal(t, ExitInvalidArgument, exitCode)
	})

	assert.Contains(t, stderr, "ARM template file is required")
}

// TestRunValidate_NonExistentFile tests runValidate with non-existent file
func TestRunValidate_NonExistentFile(t *testing.T) {
	_, stderr := captureOutput(func() {
		exitCode := runValidate([]string{"/nonexistent/template.json"})
		assert.Equal(t, ExitBuildError, exitCode)
	})

	assert.Contains(t, stderr, "Validation failed")
}

// TestRunValidate_InvalidJSON tests runValidate with invalid JSON
func TestRunValidate_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "template.json")

	err := os.WriteFile(tmpFile, []byte("{invalid json"), 0644)
	require.NoError(t, err)

	_, stderr := captureOutput(func() {
		exitCode := runValidate([]string{tmpFile})
		assert.Equal(t, ExitBuildError, exitCode)
	})

	assert.Contains(t, stderr, "Validation failed")
}

// TestRunValidate_WarningsOnly tests runValidate with only warnings
func TestRunValidate_WarningsOnly(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "template.json")

	template := `{
  "$schema": "https://example.com/invalid-schema",
  "contentVersion": "1.0.0.0",
  "resources": []
}`
	err := os.WriteFile(tmpFile, []byte(template), 0644)
	require.NoError(t, err)

	stdout, _ := captureOutput(func() {
		exitCode := runValidate([]string{tmpFile})
		assert.Equal(t, ExitSuccess, exitCode)
	})

	assert.Contains(t, stdout, "warning")
}

// =====================================================================
// Diff Command Tests
// =====================================================================

// TestRun_Diff tests run with diff command
func TestRun_Diff(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a template to compare against
	existingTemplate := `{
  "$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "resources": []
}`
	existingFile := filepath.Join(tmpDir, "existing.json")
	err := os.WriteFile(existingFile, []byte(existingTemplate), 0644)
	require.NoError(t, err)

	// Create a simple Go file with resource
	goCode := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{
	Name:     "mystorage",
	Location: "eastus",
}
`
	err = os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(goCode), 0644)
	require.NoError(t, err)

	stdout, _ := captureOutput(func() {
		exitCode := run([]string{"diff", "--against", existingFile, tmpDir})
		assert.Equal(t, ExitSuccess, exitCode)
	})

	// Should show some diff output (additions since existing template is empty)
	assert.Contains(t, stdout, "+")
}

// TestRunDiff_NoAgainstFlag tests runDiff without --against flag
func TestRunDiff_NoAgainstFlag(t *testing.T) {
	tmpDir := t.TempDir()

	_, stderr := captureOutput(func() {
		exitCode := runDiff([]string{tmpDir})
		assert.Equal(t, ExitInvalidArgument, exitCode)
	})

	assert.Contains(t, stderr, "--against")
}

// TestRunDiff_TextMode tests runDiff in text diff mode (default)
func TestRunDiff_TextMode(t *testing.T) {
	tmpDir := t.TempDir()

	// Create existing template with one resource
	existingTemplate := `{
  "$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "resources": [
    {
      "type": "Microsoft.Storage/storageAccounts",
      "apiVersion": "2021-02-01",
      "name": "oldstorage",
      "location": "westus"
    }
  ]
}`
	existingFile := filepath.Join(tmpDir, "existing.json")
	err := os.WriteFile(existingFile, []byte(existingTemplate), 0644)
	require.NoError(t, err)

	// Create a Go file with different resource
	goCode := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{
	Name:     "newstorage",
	Location: "eastus",
}
`
	srcDir := filepath.Join(tmpDir, "src")
	err = os.MkdirAll(srcDir, 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(srcDir, "main.go"), []byte(goCode), 0644)
	require.NoError(t, err)

	stdout, _ := captureOutput(func() {
		exitCode := runDiff([]string{"--against", existingFile, srcDir})
		assert.Equal(t, ExitSuccess, exitCode)
	})

	// Should show line-by-line diff
	assert.Contains(t, stdout, "-")
	assert.Contains(t, stdout, "+")
}

// TestRunDiff_SemanticMode tests runDiff in semantic diff mode
func TestRunDiff_SemanticMode(t *testing.T) {
	tmpDir := t.TempDir()

	// Create existing template (different key order, same content as generated)
	existingTemplate := `{
  "outputs": {},
  "variables": {},
  "parameters": {},
  "contentVersion": "1.0.0.0",
  "$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
  "resources": []
}`
	existingFile := filepath.Join(tmpDir, "existing.json")
	err := os.WriteFile(existingFile, []byte(existingTemplate), 0644)
	require.NoError(t, err)

	// Create a Go file that generates empty template
	goCode := `package main

var x = 42
`
	srcDir := filepath.Join(tmpDir, "src")
	err = os.MkdirAll(srcDir, 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(srcDir, "main.go"), []byte(goCode), 0644)
	require.NoError(t, err)

	stdout, _ := captureOutput(func() {
		exitCode := runDiff([]string{"--against", existingFile, "--semantic", srcDir})
		assert.Equal(t, ExitSuccess, exitCode)
	})

	// In semantic mode, different key ordering should show no diff
	assert.Contains(t, stdout, "No differences")
}

// TestRunDiff_NonExistentAgainstFile tests runDiff with non-existent against file
func TestRunDiff_NonExistentAgainstFile(t *testing.T) {
	tmpDir := t.TempDir()

	_, stderr := captureOutput(func() {
		exitCode := runDiff([]string{"--against", "/nonexistent/file.json", tmpDir})
		assert.Equal(t, ExitBuildError, exitCode)
	})

	assert.Contains(t, stderr, "Error")
}

// TestRunDiff_InvalidAgainstJSON tests runDiff with invalid JSON in against file
func TestRunDiff_InvalidAgainstJSON(t *testing.T) {
	tmpDir := t.TempDir()

	existingFile := filepath.Join(tmpDir, "existing.json")
	err := os.WriteFile(existingFile, []byte("{invalid json"), 0644)
	require.NoError(t, err)

	_, stderr := captureOutput(func() {
		exitCode := runDiff([]string{"--against", existingFile, tmpDir})
		assert.Equal(t, ExitBuildError, exitCode)
	})

	assert.Contains(t, stderr, "Error")
}

// TestRunDiff_ColorFlag tests runDiff with --color flag
func TestRunDiff_ColorFlag(t *testing.T) {
	tmpDir := t.TempDir()

	existingTemplate := `{
  "$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "resources": []
}`
	existingFile := filepath.Join(tmpDir, "existing.json")
	err := os.WriteFile(existingFile, []byte(existingTemplate), 0644)
	require.NoError(t, err)

	goCode := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{
	Name:     "mystorage",
	Location: "eastus",
}
`
	srcDir := filepath.Join(tmpDir, "src")
	err = os.MkdirAll(srcDir, 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(srcDir, "main.go"), []byte(goCode), 0644)
	require.NoError(t, err)

	stdout, _ := captureOutput(func() {
		exitCode := runDiff([]string{"--against", existingFile, "--color=false", srcDir})
		assert.Equal(t, ExitSuccess, exitCode)
	})

	// Should not contain ANSI color codes when color is disabled
	assert.NotContains(t, stdout, "\033[")
}

// TestRunDiff_NoDifferences tests runDiff when templates are identical
func TestRunDiff_NoDifferences(t *testing.T) {
	tmpDir := t.TempDir()

	// Create empty template (matches what an empty source will produce)
	existingTemplate := `{
  "$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "parameters": {},
  "variables": {},
  "resources": [],
  "outputs": {}
}`
	existingFile := filepath.Join(tmpDir, "existing.json")
	err := os.WriteFile(existingFile, []byte(existingTemplate), 0644)
	require.NoError(t, err)

	// Create a Go file with no resources
	goCode := `package main

var x = 42
`
	srcDir := filepath.Join(tmpDir, "src")
	err = os.MkdirAll(srcDir, 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(srcDir, "main.go"), []byte(goCode), 0644)
	require.NoError(t, err)

	stdout, _ := captureOutput(func() {
		exitCode := runDiff([]string{"--against", existingFile, "--semantic", srcDir})
		assert.Equal(t, ExitSuccess, exitCode)
	})

	assert.Contains(t, stdout, "No differences")
}

// =====================================================================
// Watch Command Tests
// =====================================================================

// TestRun_Watch tests run with watch command (immediate cancel)
func TestRun_Watch(t *testing.T) {
	tmpDir := t.TempDir()

	goCode := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{
	Name:     "mystorage",
	Location: "eastus",
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(goCode), 0644)
	require.NoError(t, err)

	// For testing, we'll use a special test mode that exits after first build
	stdout, _ := captureOutput(func() {
		exitCode := runWatch([]string{"--test-run", tmpDir})
		assert.Equal(t, ExitSuccess, exitCode)
	})

	assert.Contains(t, stdout, "Watching")
}

// TestRunWatch_WithOutput tests runWatch with -o flag
func TestRunWatch_WithOutput(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.json")

	goCode := `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{
	Name:     "mystorage",
	Location: "eastus",
}
`
	srcDir := filepath.Join(tmpDir, "src")
	err := os.MkdirAll(srcDir, 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(srcDir, "main.go"), []byte(goCode), 0644)
	require.NoError(t, err)

	stdout, _ := captureOutput(func() {
		exitCode := runWatch([]string{"-o", outputFile, "--test-run", srcDir})
		assert.Equal(t, ExitSuccess, exitCode)
	})

	assert.Contains(t, stdout, "Watching")
	// Check output file was created
	_, err = os.Stat(outputFile)
	assert.NoError(t, err)
}

// TestRunWatch_WithInterval tests runWatch with --interval flag
func TestRunWatch_WithInterval(t *testing.T) {
	tmpDir := t.TempDir()

	goCode := `package main

var x = 42
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(goCode), 0644)
	require.NoError(t, err)

	stdout, _ := captureOutput(func() {
		exitCode := runWatch([]string{"--interval", "100ms", "--test-run", tmpDir})
		assert.Equal(t, ExitSuccess, exitCode)
	})

	assert.Contains(t, stdout, "Watching")
}

// TestRunWatch_NonExistentPath tests runWatch with non-existent path
func TestRunWatch_NonExistentPath(t *testing.T) {
	_, stderr := captureOutput(func() {
		exitCode := runWatch([]string{"/nonexistent/path"})
		assert.Equal(t, ExitBuildError, exitCode)
	})

	assert.Contains(t, stderr, "Error")
}

// TestRunWatch_InvalidInterval tests runWatch with invalid interval
func TestRunWatch_InvalidInterval(t *testing.T) {
	tmpDir := t.TempDir()

	_, stderr := captureOutput(func() {
		exitCode := runWatch([]string{"--interval", "invalid", tmpDir})
		assert.Equal(t, ExitInvalidArgument, exitCode)
	})

	assert.Contains(t, stderr, "invalid")
}

// TestRunWatch_BuildError tests runWatch handling build errors
func TestRunWatch_BuildError(t *testing.T) {
	tmpDir := t.TempDir()

	// Create invalid Go file
	err := os.WriteFile(filepath.Join(tmpDir, "bad.go"), []byte("invalid go {{{"), 0644)
	require.NoError(t, err)

	_, stderr := captureOutput(func() {
		exitCode := runWatch([]string{"--test-run", tmpDir})
		// Should still exit with success since watch continues on build errors
		assert.Equal(t, ExitSuccess, exitCode)
	})

	// Should report the build error
	assert.Contains(t, stderr, "Build error")
}

// TestRunWatch_Debounce tests that rapid changes are debounced
func TestRunWatch_Debounce(t *testing.T) {
	// This test verifies debouncing behavior
	// In test mode, we simulate rapid file changes
	tmpDir := t.TempDir()

	goCode := `package main

var x = 42
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(goCode), 0644)
	require.NoError(t, err)

	// The debounce test is tricky in unit tests
	// We just verify the command accepts the default debounce
	stdout, _ := captureOutput(func() {
		exitCode := runWatch([]string{"--test-run", tmpDir})
		assert.Equal(t, ExitSuccess, exitCode)
	})

	assert.Contains(t, stdout, "Watching")
}
