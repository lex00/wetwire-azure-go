package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lex00/wetwire-azure-go/internal/discover"
	"github.com/lex00/wetwire-azure-go/internal/importer"
	"github.com/lex00/wetwire-azure-go/internal/linter"
	"github.com/lex00/wetwire-azure-go/internal/template"
)

// Exit codes
const (
	ExitSuccess         = 0
	ExitBuildError      = 1
	ExitInvalidArgument = 2
)

// osExit is the function used to exit the program. Can be overridden for testing.
var osExit = os.Exit

func main() {
	code := run(os.Args[1:])
	osExit(code)
}

// run executes the CLI and returns an exit code
func run(args []string) int {
	if len(args) < 1 {
		printUsage()
		return ExitInvalidArgument
	}

	command := args[0]

	switch command {
	case "build":
		return runBuild(args[1:])
	case "import":
		return runImport(args[1:])
	case "lint":
		return runLint(args[1:])
	case "list":
		return runList(args[1:])
	case "init":
		return runInit(args[1:])
	case "graph":
		return runGraph(args[1:])
	case "help", "-h", "--help":
		printUsage()
		return ExitSuccess
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		printUsage()
		return ExitInvalidArgument
	}
}

func printUsage() {
	fmt.Println("wetwire-azure - Azure ARM/Bicep template synthesis")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  wetwire-azure build [package-path] [flags]  Build ARM template from Go code")
	fmt.Println("  wetwire-azure import <arm-file> [flags]     Import ARM JSON to Go code")
	fmt.Println("  wetwire-azure lint [path]                   Lint infrastructure code")
	fmt.Println("  wetwire-azure list [path] [flags]           List discovered resources")
	fmt.Println("  wetwire-azure init [directory]              Initialize new wetwire project")
	fmt.Println("  wetwire-azure graph [path] [flags]          Generate resource dependency graph")
	fmt.Println("  wetwire-azure help                          Show this help message")
	fmt.Println()
	fmt.Println("Options for build:")
	fmt.Println("  -o, --output <file>       Output file path (default: stdout)")
	fmt.Println("  --format <format>         Output format: arm (default: arm)")
	fmt.Println("  --parameters-file <file>  Write parameters to separate file")
	fmt.Println()
	fmt.Println("Options for import:")
	fmt.Println("  -o, --output <file>       Output file path (default: stdout)")
	fmt.Println("  --package <name>          Package name for generated code (default: infra)")
	fmt.Println()
	fmt.Println("Options for lint:")
	fmt.Println("  --fix                     Auto-fix issues where possible (not yet implemented)")
	fmt.Println()
	fmt.Println("Options for list:")
	fmt.Println("  --format <format>         Output format: table, json (default: table)")
	fmt.Println()
	fmt.Println("Options for graph:")
	fmt.Println("  -o, --output <file>       Output file path (default: stdout)")
	fmt.Println("  --format <format>         Output format: dot, mermaid (default: dot)")
	fmt.Println()
	fmt.Println("Exit codes:")
	fmt.Println("  0  Success")
	fmt.Println("  1  Build error")
	fmt.Println("  2  Invalid arguments")
}

// runBuild executes the build command and returns an exit code
func runBuild(args []string) int {
	fs := flag.NewFlagSet("build", flag.ContinueOnError)

	var outputFile string
	var outputFileLong string
	var format string
	var parametersFile string

	fs.StringVar(&outputFile, "o", "", "Output file path")
	fs.StringVar(&outputFileLong, "output", "", "Output file path")
	fs.StringVar(&format, "format", "arm", "Output format (arm)")
	fs.StringVar(&parametersFile, "parameters-file", "", "Write parameters to separate file")

	// Custom error handling to return proper exit code
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return ExitInvalidArgument
	}

	// Use long form if short form not specified
	if outputFile == "" && outputFileLong != "" {
		outputFile = outputFileLong
	}

	// Validate format
	if format != "arm" {
		fmt.Fprintf(os.Stderr, "Error: unsupported format '%s'. Supported formats: arm\n", format)
		return ExitInvalidArgument
	}

	// Get package path (default to current directory)
	packagePath := "."
	if fs.NArg() > 0 {
		packagePath = fs.Arg(0)
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(packagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving path: %v\n", err)
		return ExitBuildError
	}

	// Check if path exists
	if _, err := os.Stat(absPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return ExitBuildError
	}

	// Build the template
	templateJSON, parametersJSON, err := buildTemplate(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Build failed: %v\n", err)
		return ExitBuildError
	}

	// Write output
	if outputFile != "" {
		if err := os.WriteFile(outputFile, []byte(templateJSON), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
			return ExitBuildError
		}
	} else {
		fmt.Println(templateJSON)
	}

	// Write parameters file if specified
	if parametersFile != "" {
		if err := os.WriteFile(parametersFile, []byte(parametersJSON), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing parameters file: %v\n", err)
			return ExitBuildError
		}
	}

	return ExitSuccess
}

// buildTemplate discovers resources and builds the ARM template
func buildTemplate(srcDir string) (string, string, error) {
	// DISCOVER - Find resources in the package
	resources, err := discover.DiscoverResources(srcDir)
	if err != nil {
		return "", "", fmt.Errorf("discovery failed: %w", err)
	}

	// BUILD - Create template builder and add resources
	builder := template.NewTemplateBuilder()
	for _, res := range resources {
		if err := builder.AddResource(res); err != nil {
			return "", "", fmt.Errorf("failed to add resource %s: %w", res.Name, err)
		}
	}

	// SERIALIZE - Build the ARM template JSON
	templateJSON, err := builder.Build()
	if err != nil {
		return "", "", fmt.Errorf("template build failed: %w", err)
	}

	// Generate parameters file content
	parametersJSON := generateParametersFile()

	return templateJSON, parametersJSON, nil
}

// generateParametersFile creates an ARM template parameters file structure
func generateParametersFile() string {
	params := map[string]interface{}{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentParameters.json#",
		"contentVersion": "1.0.0.0",
		"parameters":     map[string]interface{}{},
	}

	jsonBytes, _ := json.MarshalIndent(params, "", "  ")
	return string(jsonBytes)
}

// runImport executes the import command and returns an exit code
func runImport(args []string) int {
	fs := flag.NewFlagSet("import", flag.ContinueOnError)

	var outputFile string
	var outputFileLong string
	var packageName string

	fs.StringVar(&outputFile, "o", "", "Output file path")
	fs.StringVar(&outputFileLong, "output", "", "Output file path")
	fs.StringVar(&packageName, "package", "infra", "Package name for generated code")

	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return ExitInvalidArgument
	}

	// Use long form if short form not specified
	if outputFile == "" && outputFileLong != "" {
		outputFile = outputFileLong
	}

	// Require input file
	if fs.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Error: ARM template file is required\n")
		printUsage()
		return ExitInvalidArgument
	}

	inputFile := fs.Arg(0)

	// Read input file
	data, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		return ExitBuildError
	}

	// Parse ARM template
	template, err := importer.ParseARMTemplate(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing ARM template: %v\n", err)
		return ExitBuildError
	}

	// Generate Go code
	goCode, err := importer.GenerateGoCode(template, packageName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating Go code: %v\n", err)
		return ExitBuildError
	}

	// Write output
	if outputFile != "" {
		if err := os.WriteFile(outputFile, []byte(goCode), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
			return ExitBuildError
		}
	} else {
		fmt.Println(goCode)
	}

	return ExitSuccess
}

// runLint executes the lint command and returns an exit code
func runLint(args []string) int {
	fs := flag.NewFlagSet("lint", flag.ContinueOnError)
	fixFlag := fs.Bool("fix", false, "Auto-fix issues where possible")
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return ExitInvalidArgument
	}

	// Default to current directory if no path provided
	path := "."
	if fs.NArg() > 0 {
		path = fs.Arg(0)
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving path: %v\n", err)
		return ExitBuildError
	}

	// Check if path exists
	info, err := os.Stat(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return ExitBuildError
	}

	// Create linter
	l := linter.NewLinter()

	var results []linter.LintResult

	// Check file or directory
	if info.IsDir() {
		results, err = l.CheckDirectory(absPath)
	} else {
		results, err = l.CheckFile(absPath)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Linting failed: %v\n", err)
		return ExitBuildError
	}

	// Print results
	if len(results) == 0 {
		fmt.Println("No issues found.")
		return ExitSuccess
	}

	// Group results by severity
	errorCount := 0
	warningCount := 0

	for _, result := range results {
		fmt.Println(result.String())
		switch result.Severity {
		case linter.SeverityError:
			errorCount++
		case linter.SeverityWarning:
			warningCount++
		}
	}

	fmt.Println()
	fmt.Printf("Found %d error(s) and %d warning(s)\n", errorCount, warningCount)

	if *fixFlag {
		fmt.Println("Note: Auto-fix is not yet implemented")
	}

	// Exit with error code if issues were found
	if errorCount > 0 || warningCount > 0 {
		return ExitBuildError
	}

	return ExitSuccess
}

// runList executes the list command and returns an exit code
func runList(args []string) int {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)

	var format string
	fs.StringVar(&format, "format", "table", "Output format (table, json)")

	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return ExitInvalidArgument
	}

	// Validate format
	if format != "table" && format != "json" {
		fmt.Fprintf(os.Stderr, "Error: unsupported format '%s'. Supported formats: table, json\n", format)
		return ExitInvalidArgument
	}

	// Default to current directory if no path provided
	path := "."
	if fs.NArg() > 0 {
		path = fs.Arg(0)
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving path: %v\n", err)
		return ExitBuildError
	}

	// Check if path exists
	if _, err := os.Stat(absPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return ExitBuildError
	}

	// Discover resources
	resources, err := discover.DiscoverResources(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Discovery failed: %v\n", err)
		return ExitBuildError
	}

	// Output based on format
	if format == "json" {
		return outputListJSON(resources)
	}
	return outputListTable(resources)
}

// outputListJSON outputs resources in JSON format
func outputListJSON(resources []discover.DiscoveredResource) int {
	type jsonResource struct {
		Name         string   `json:"name"`
		Type         string   `json:"type"`
		File         string   `json:"file"`
		Line         int      `json:"line"`
		Dependencies []string `json:"dependencies,omitempty"`
	}

	jsonResources := make([]jsonResource, len(resources))
	for i, res := range resources {
		jsonResources[i] = jsonResource{
			Name:         res.Name,
			Type:         res.Type,
			File:         res.File,
			Line:         res.Line,
			Dependencies: res.Dependencies,
		}
	}

	output, err := json.MarshalIndent(jsonResources, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating JSON: %v\n", err)
		return ExitBuildError
	}

	fmt.Println(string(output))
	return ExitSuccess
}

// outputListTable outputs resources in table format
func outputListTable(resources []discover.DiscoveredResource) int {
	if len(resources) == 0 {
		fmt.Println("NAME  TYPE  FILE  LINE")
		fmt.Println("----  ----  ----  ----")
		fmt.Println("No resources found.")
		return ExitSuccess
	}

	// Calculate column widths
	maxName := len("NAME")
	maxType := len("TYPE")
	maxFile := len("FILE")

	for _, res := range resources {
		if len(res.Name) > maxName {
			maxName = len(res.Name)
		}
		if len(res.Type) > maxType {
			maxType = len(res.Type)
		}
		fileName := filepath.Base(res.File)
		if len(fileName) > maxFile {
			maxFile = len(fileName)
		}
	}

	// Print header
	fmt.Printf("%-*s  %-*s  %-*s  LINE\n", maxName, "NAME", maxType, "TYPE", maxFile, "FILE")
	fmt.Printf("%s  %s  %s  ----\n",
		strings.Repeat("-", maxName),
		strings.Repeat("-", maxType),
		strings.Repeat("-", maxFile))

	// Print resources
	for _, res := range resources {
		fileName := filepath.Base(res.File)
		fmt.Printf("%-*s  %-*s  %-*s  %d\n",
			maxName, res.Name,
			maxType, res.Type,
			maxFile, fileName,
			res.Line)
	}

	return ExitSuccess
}

// runInit executes the init command and returns an exit code
func runInit(args []string) int {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return ExitInvalidArgument
	}

	// Default to current directory if no path provided
	targetDir := "."
	if fs.NArg() > 0 {
		targetDir = fs.Arg(0)
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(targetDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving path: %v\n", err)
		return ExitBuildError
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(absPath, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directory: %v\n", err)
		return ExitBuildError
	}

	// Check if go.mod already exists
	goModPath := filepath.Join(absPath, "go.mod")
	if _, err := os.Stat(goModPath); err == nil {
		fmt.Fprintf(os.Stderr, "Error: go.mod already exists in %s\n", absPath)
		return ExitBuildError
	}

	// Get module name from directory name
	moduleName := filepath.Base(absPath)

	// Create go.mod
	goModContent := fmt.Sprintf(`module %s

go 1.21

require github.com/lex00/wetwire-azure-go v0.1.0
`, moduleName)

	if err := os.WriteFile(goModPath, []byte(goModContent), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing go.mod: %v\n", err)
		return ExitBuildError
	}

	// Create main.go
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
	mainGoPath := filepath.Join(absPath, "main.go")
	if err := os.WriteFile(mainGoPath, []byte(mainGoContent), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing main.go: %v\n", err)
		return ExitBuildError
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
	gitignorePath := filepath.Join(absPath, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing .gitignore: %v\n", err)
		return ExitBuildError
	}

	// Create README.md
	readmeContent := fmt.Sprintf(`# %s

Azure infrastructure as code using wetwire-azure-go.

## Getting Started

1. Install dependencies:
   `+"```bash"+`
   go mod download
   `+"```"+`

2. Build the ARM template:
   `+"```bash"+`
   wetwire-azure build -o template.json
   `+"```"+`

3. Deploy to Azure:
   `+"```bash"+`
   az deployment group create \
     --resource-group <your-rg> \
     --template-file template.json
   `+"```"+`

## Available Commands

- `+"`wetwire-azure build`"+` - Generate ARM template from Go code
- `+"`wetwire-azure list`"+` - List all discovered resources
- `+"`wetwire-azure graph`"+` - Generate dependency graph
- `+"`wetwire-azure lint`"+` - Lint infrastructure code

## Learn More

- [wetwire-azure-go Documentation](https://github.com/lex00/wetwire-azure-go)
- [Azure ARM Templates](https://docs.microsoft.com/en-us/azure/azure-resource-manager/templates/)
`, moduleName)

	readmePath := filepath.Join(absPath, "README.md")
	if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing README.md: %v\n", err)
		return ExitBuildError
	}

	fmt.Printf("Initialized wetwire-azure project in %s\n", absPath)
	fmt.Println("\nCreated files:")
	fmt.Println("  - go.mod")
	fmt.Println("  - main.go")
	fmt.Println("  - .gitignore")
	fmt.Println("  - README.md")
	fmt.Println("\nNext steps:")
	fmt.Println("  1. cd", absPath)
	fmt.Println("  2. go mod download")
	fmt.Println("  3. wetwire-azure build")

	return ExitSuccess
}

// runGraph executes the graph command and returns an exit code
func runGraph(args []string) int {
	fs := flag.NewFlagSet("graph", flag.ContinueOnError)

	var outputFile string
	var outputFileLong string
	var format string

	fs.StringVar(&outputFile, "o", "", "Output file path")
	fs.StringVar(&outputFileLong, "output", "", "Output file path")
	fs.StringVar(&format, "format", "dot", "Output format (dot, mermaid)")

	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return ExitInvalidArgument
	}

	// Use long form if short form not specified
	if outputFile == "" && outputFileLong != "" {
		outputFile = outputFileLong
	}

	// Validate format
	if format != "dot" && format != "mermaid" {
		fmt.Fprintf(os.Stderr, "Error: unsupported format '%s'. Supported formats: dot, mermaid\n", format)
		return ExitInvalidArgument
	}

	// Default to current directory if no path provided
	path := "."
	if fs.NArg() > 0 {
		path = fs.Arg(0)
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving path: %v\n", err)
		return ExitBuildError
	}

	// Check if path exists
	if _, err := os.Stat(absPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return ExitBuildError
	}

	// Discover resources
	resources, err := discover.DiscoverResources(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Discovery failed: %v\n", err)
		return ExitBuildError
	}

	// Generate graph
	var graphOutput string
	if format == "mermaid" {
		graphOutput = generateMermaidGraph(resources)
	} else {
		graphOutput = generateDOTGraph(resources)
	}

	// Write output
	if outputFile != "" {
		if err := os.WriteFile(outputFile, []byte(graphOutput), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
			return ExitBuildError
		}
		fmt.Printf("Graph written to %s\n", outputFile)
	} else {
		fmt.Println(graphOutput)
	}

	return ExitSuccess
}

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

// runValidate executes the validate command and returns an exit code

// runValidate executes the validate command and returns an exit code
