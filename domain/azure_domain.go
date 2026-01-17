// Package domain provides the AzureDomain implementation for wetwire-core-go.
package domain

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	coredomain "github.com/lex00/wetwire-core-go/domain"
	"github.com/lex00/wetwire-azure-go/internal/discover"
	"github.com/lex00/wetwire-azure-go/internal/lint"
	"github.com/lex00/wetwire-azure-go/internal/template"
)

// AzureDomain implements the Domain interface for Azure infrastructure.
type AzureDomain struct{}

// Compile-time checks
var (
	_ coredomain.Domain        = (*AzureDomain)(nil)
	_ coredomain.ListerDomain  = (*AzureDomain)(nil)
	_ coredomain.GrapherDomain = (*AzureDomain)(nil)
)

// Name returns "azure"
func (d *AzureDomain) Name() string {
	return "azure"
}

// Version returns the current version
func (d *AzureDomain) Version() string {
	return Version
}

// Builder returns the Azure builder implementation
func (d *AzureDomain) Builder() coredomain.Builder {
	return &azureBuilder{}
}

// Linter returns the Azure linter implementation
func (d *AzureDomain) Linter() coredomain.Linter {
	return &azureLinter{}
}

// Initializer returns the Azure initializer implementation
func (d *AzureDomain) Initializer() coredomain.Initializer {
	return &azureInitializer{}
}

// Validator returns the Azure validator implementation
func (d *AzureDomain) Validator() coredomain.Validator {
	return &azureValidator{}
}

// Lister returns the Azure lister implementation
func (d *AzureDomain) Lister() coredomain.Lister {
	return &azureLister{}
}

// Grapher returns the Azure grapher implementation
func (d *AzureDomain) Grapher() coredomain.Grapher {
	return &azureGrapher{}
}

// azureBuilder implements domain.Builder
type azureBuilder struct{}

func (b *azureBuilder) Build(ctx *Context, path string, opts BuildOpts) (*Result, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve path: %w", err)
	}

	// Discover all resources
	resources, err := discover.DiscoverResources(absPath)
	if err != nil {
		return nil, fmt.Errorf("discovery failed: %w", err)
	}

	if len(resources) == 0 {
		return NewErrorResult("no resources found", Error{
			Path:    absPath,
			Message: "no Azure resources found",
		}), nil
	}

	// Build template
	builder := template.NewTemplateBuilder()
	for _, res := range resources {
		if err := builder.AddResource(res); err != nil {
			return nil, fmt.Errorf("failed to add resource %s: %w", res.Name, err)
		}
	}

	// Generate ARM template JSON
	templateJSON, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("template build failed: %w", err)
	}

	// Handle output file
	if !opts.DryRun && opts.Output != "" {
		if err := os.WriteFile(opts.Output, []byte(templateJSON), 0644); err != nil {
			return nil, fmt.Errorf("write output: %w", err)
		}
		return NewResult(fmt.Sprintf("Wrote %s", opts.Output)), nil
	}

	return NewResultWithData("Build completed", templateJSON), nil
}

// azureLinter implements domain.Linter
type azureLinter struct{}

func (l *azureLinter) Lint(ctx *Context, path string, opts LintOpts) (*Result, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve path: %w", err)
	}

	// Create linter
	azureLint := lint.NewLinter()

	var results []lint.LintResult

	// Check file or directory
	info, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("stat path: %w", err)
	}

	if info.IsDir() {
		results, err = azureLint.CheckDirectory(absPath)
	} else {
		results, err = azureLint.CheckFile(absPath)
	}

	if err != nil {
		return nil, fmt.Errorf("linting failed: %w", err)
	}

	if len(results) == 0 {
		return NewResult("No lint issues found"), nil
	}

	// Convert to domain errors
	errs := make([]Error, 0, len(results))
	for _, r := range results {
		errs = append(errs, Error{
			Path:     r.File,
			Line:     r.Line,
			Severity: r.Severity.String(),
			Message:  r.Message,
			Code:     r.Rule,
		})
	}

	return NewErrorResultMultiple("lint issues found", errs), nil
}

// azureInitializer implements domain.Initializer
type azureInitializer struct{}

func (i *azureInitializer) Init(ctx *Context, path string, opts InitOpts) (*Result, error) {
	// Use opts.Path if provided, otherwise fall back to path argument
	targetPath := opts.Path
	if targetPath == "" || targetPath == "." {
		targetPath = path
	}

	// Handle scenario initialization
	if opts.Scenario {
		name := opts.Name
		if name == "" {
			name = filepath.Base(targetPath)
		}

		description := opts.Description
		if description == "" {
			description = "Azure ARM template scenario"
		}

		// Use core's scenario scaffolding
		scenario := coredomain.ScaffoldScenario(name, description, "azure")
		created, err := coredomain.WriteScenario(targetPath, scenario)
		if err != nil {
			return NewErrorResult(err.Error(), Error{
				Path:    targetPath,
				Message: err.Error(),
			}), nil
		}

		return NewResultWithData(fmt.Sprintf("Created scenario in %s", targetPath), map[string]interface{}{
			"created": created,
		}), nil
	}

	// Create directory
	if err := os.MkdirAll(targetPath, 0755); err != nil {
		return nil, fmt.Errorf("create directory: %w", err)
	}

	// Check if go.mod already exists
	goModPath := filepath.Join(targetPath, "go.mod")
	if _, err := os.Stat(goModPath); err == nil {
		return NewErrorResult("go.mod already exists", Error{
			Path:    goModPath,
			Message: "go.mod already exists in this directory",
		}), nil
	}

	// Get module name from directory name
	moduleName := filepath.Base(path)

	// Create go.mod
	goModContent := fmt.Sprintf(`module %s

go 1.23.0

require github.com/lex00/wetwire-azure-go v1.3.1
`, moduleName)

	if err := os.WriteFile(goModPath, []byte(goModContent), 0644); err != nil {
		return nil, fmt.Errorf("write go.mod: %w", err)
	}

	// Create example main.go
	mainGoContent := `package main

import (
	"github.com/lex00/wetwire-azure-go/resources/storage"
)

// Example storage account resource
var MyStorage = storage.StorageAccount{
	Name:     "mystorageaccount",
	Location: "eastus",
	SKU: storage.SKU{
		Name: "Standard_LRS",
	},
	Kind: "StorageV2",
	Properties: storage.StorageAccountProperties{
		AccessTier: "Hot",
	},
}
`
	mainGoPath := filepath.Join(targetPath, "main.go")
	if err := os.WriteFile(mainGoPath, []byte(mainGoContent), 0644); err != nil {
		return nil, fmt.Errorf("write main.go: %w", err)
	}

	// Create .gitignore
	gitignoreContent := `# Build outputs
*.json
*.bicep

# Go build artifacts
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test
*.out

# IDE
.vscode/
.idea/
*.swp
*.swo
*~
`
	gitignorePath := filepath.Join(targetPath, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
		return nil, fmt.Errorf("write .gitignore: %w", err)
	}

	return NewResult(fmt.Sprintf("Initialized wetwire-azure project in %s", targetPath)), nil
}

// azureValidator implements domain.Validator
type azureValidator struct{}

func (v *azureValidator) Validate(ctx *Context, path string, opts ValidateOpts) (*Result, error) {
	// For Azure, validation is the same as linting
	linter := &azureLinter{}
	return linter.Lint(ctx, path, LintOpts{})
}

// azureLister implements domain.Lister
type azureLister struct{}

func (l *azureLister) List(ctx *Context, path string, opts ListOpts) (*Result, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve path: %w", err)
	}

	// Discover all resources
	resources, err := discover.DiscoverResources(absPath)
	if err != nil {
		return nil, fmt.Errorf("discovery failed: %w", err)
	}

	// Build list
	list := make([]map[string]string, 0, len(resources))
	for _, res := range resources {
		list = append(list, map[string]string{
			"name": res.Name,
			"type": res.Type,
			"file": res.File,
			"line": fmt.Sprintf("%d", res.Line),
		})
	}

	return NewResultWithData(fmt.Sprintf("Discovered %d resources", len(list)), list), nil
}

// azureGrapher implements domain.Grapher
type azureGrapher struct{}

func (g *azureGrapher) Graph(ctx *Context, path string, opts GraphOpts) (*Result, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve path: %w", err)
	}

	// Discover all resources
	resources, err := discover.DiscoverResources(absPath)
	if err != nil {
		return nil, fmt.Errorf("discovery failed: %w", err)
	}

	// Generate graph
	var graph string
	switch opts.Format {
	case "dot", "":
		graph = generateDOTGraph(resources)
	case "mermaid":
		graph = generateMermaidGraph(resources)
	default:
		return nil, fmt.Errorf("unknown format: %s", opts.Format)
	}

	return NewResultWithData("Graph generated", graph), nil
}

// Helper functions

// generateDOTGraph generates a Graphviz DOT format graph
func generateDOTGraph(resources []discover.DiscoveredResource) string {
	var sb strings.Builder

	sb.WriteString("digraph \"Azure Resources\" {\n")
	sb.WriteString("  rankdir=TB;\n")
	sb.WriteString("  node [shape=box, style=rounded];\n")
	sb.WriteString("\n")

	// Add nodes
	for _, res := range resources {
		// Escape quotes in labels
		label := fmt.Sprintf("%s\\n%s", res.Name, res.Type)
		sb.WriteString(fmt.Sprintf("  \"%s\" [label=\"%s\"];\n", res.Name, label))
	}

	// Add edges (dependencies)
	sb.WriteString("\n")
	for _, res := range resources {
		for _, dep := range res.Dependencies {
			// Check if dependency is a resource
			if isResource(dep, resources) {
				sb.WriteString(fmt.Sprintf("  \"%s\" -> \"%s\";\n", res.Name, dep))
			}
		}
	}

	sb.WriteString("}\n")
	return sb.String()
}

// generateMermaidGraph generates a Mermaid format graph
func generateMermaidGraph(resources []discover.DiscoveredResource) string {
	var sb strings.Builder

	sb.WriteString("graph TD\n")

	// Add nodes
	for _, res := range resources {
		// Sanitize for Mermaid (replace spaces and special chars)
		label := fmt.Sprintf("%s<br/>%s", res.Name, res.Type)
		sb.WriteString(fmt.Sprintf("  %s[\"%s\"]\n", res.Name, label))
	}

	// Add edges (dependencies)
	for _, res := range resources {
		for _, dep := range res.Dependencies {
			// Check if dependency is a resource
			if isResource(dep, resources) {
				sb.WriteString(fmt.Sprintf("  %s --> %s\n", res.Name, dep))
			}
		}
	}

	return sb.String()
}

// isResource checks if a name corresponds to a discovered resource
func isResource(name string, resources []discover.DiscoveredResource) bool {
	for _, res := range resources {
		if res.Name == name {
			return true
		}
	}
	return false
}
