// Package main implements the MCP server for wetwire-azure.
// This allows Claude Code to interact with wetwire-azure through tool calls.
package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lex00/wetwire-azure-go/internal/discover"
	"github.com/lex00/wetwire-azure-go/internal/importer"
	"github.com/lex00/wetwire-azure-go/internal/linter"
	"github.com/lex00/wetwire-azure-go/internal/template"
	"github.com/lex00/wetwire-core-go/mcp"
)

const (
	serverName    = "wetwire-azure"
	serverVersion = "1.0.0"
)

func main() {
	// Check for --install flag to show configuration
	if len(os.Args) > 1 && os.Args[1] == "--install" {
		fmt.Println(mcp.GetInstallInstructions(serverName, "wetwire-azure-mcp"))
		return
	}

	// Create MCP server
	server := mcp.NewServer(mcp.Config{
		Name:    serverName,
		Version: serverVersion,
	})

	// Register tools
	registerBuildTool(server)
	registerLintTool(server)
	registerImportTool(server)

	// Start the server
	if err := server.Start(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}

// registerBuildTool registers the build tool with the MCP server.
func registerBuildTool(server *mcp.Server) {
	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "Path to the Go package containing Azure resource definitions (defaults to current directory)",
			},
		},
	}

	server.RegisterToolWithSchema("build", "Build ARM template from Go resource definitions", handleBuild, schema)
}

// registerLintTool registers the lint tool with the MCP server.
func registerLintTool(server *mcp.Server) {
	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "Path to file or directory to lint (defaults to current directory)",
			},
		},
	}

	server.RegisterToolWithSchema("lint", "Check Azure resource definitions for issues", handleLint, schema)
}

// registerImportTool registers the import tool with the MCP server.
func registerImportTool(server *mcp.Server) {
	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"file": map[string]any{
				"type":        "string",
				"description": "Path to ARM template JSON file to import",
			},
			"package": map[string]any{
				"type":        "string",
				"description": "Package name for generated Go code (defaults to 'infra')",
			},
		},
		"required": []string{"file"},
	}

	server.RegisterToolWithSchema("import", "Convert ARM template JSON to Go code", handleImport, schema)
}

// handleBuild processes the build tool invocation.
func handleBuild(ctx context.Context, args map[string]any) (string, error) {
	// Get path argument
	path := "."
	if p, ok := args["path"].(string); ok && p != "" {
		path = p
	}

	// Resolve to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("error resolving path: %w", err)
	}

	// Check if path exists
	if _, err := os.Stat(absPath); err != nil {
		return "", fmt.Errorf("path not found: %w", err)
	}

	// Discover resources
	resources, err := discover.DiscoverResources(absPath)
	if err != nil {
		return "", fmt.Errorf("discovery failed: %w", err)
	}

	if len(resources) == 0 {
		return "No Azure resources found in the specified path.", nil
	}

	// Build template
	builder := template.NewTemplateBuilder()
	for _, res := range resources {
		if err := builder.AddResource(res); err != nil {
			return "", fmt.Errorf("failed to add resource %s: %w", res.Name, err)
		}
	}

	templateJSON, err := builder.Build()
	if err != nil {
		return "", fmt.Errorf("template build failed: %w", err)
	}

	return fmt.Sprintf("Successfully built ARM template with %d resource(s):\n\n%s", len(resources), templateJSON), nil
}

// handleLint processes the lint tool invocation.
func handleLint(ctx context.Context, args map[string]any) (string, error) {
	// Get path argument
	path := "."
	if p, ok := args["path"].(string); ok && p != "" {
		path = p
	}

	// Resolve to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("error resolving path: %w", err)
	}

	// Check if path exists
	info, err := os.Stat(absPath)
	if err != nil {
		return "", fmt.Errorf("path not found: %w", err)
	}

	// Create linter and check
	l := linter.NewLinter()
	var results []linter.LintResult

	if info.IsDir() {
		results, err = l.CheckDirectory(absPath)
	} else {
		results, err = l.CheckFile(absPath)
	}

	if err != nil {
		return "", fmt.Errorf("linting failed: %w", err)
	}

	// Format results
	if len(results) == 0 {
		return "No issues found.", nil
	}

	var sb strings.Builder
	errorCount := 0
	warningCount := 0

	for _, result := range results {
		sb.WriteString(result.String())
		sb.WriteString("\n")
		switch result.Severity {
		case linter.SeverityError:
			errorCount++
		case linter.SeverityWarning:
			warningCount++
		}
	}

	sb.WriteString(fmt.Sprintf("\nFound %d error(s) and %d warning(s)", errorCount, warningCount))

	return sb.String(), nil
}

// handleImport processes the import tool invocation.
func handleImport(ctx context.Context, args map[string]any) (string, error) {
	// Get file argument (required)
	file, ok := args["file"].(string)
	if !ok || file == "" {
		return "", fmt.Errorf("file argument is required")
	}

	// Get package argument
	packageName := "infra"
	if p, ok := args["package"].(string); ok && p != "" {
		packageName = p
	}

	// Read input file
	data, err := os.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}

	// Parse ARM template
	tmpl, err := importer.ParseARMTemplate(data)
	if err != nil {
		return "", fmt.Errorf("error parsing ARM template: %w", err)
	}

	// Generate Go code
	goCode, err := importer.GenerateGoCode(tmpl, packageName)
	if err != nil {
		return "", fmt.Errorf("error generating Go code: %w", err)
	}

	resourceCount := 0
	if tmpl.Resources != nil {
		resourceCount = len(tmpl.Resources)
	}

	return fmt.Sprintf("Successfully imported ARM template with %d resource(s):\n\n%s", resourceCount, goCode), nil
}
