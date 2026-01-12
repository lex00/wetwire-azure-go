package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lex00/wetwire-azure-go/internal/discover"
	"github.com/lex00/wetwire-azure-go/internal/importer"
	"github.com/lex00/wetwire-azure-go/internal/linter"
	"github.com/lex00/wetwire-azure-go/internal/template"
	"github.com/lex00/wetwire-azure-go/internal/validator"
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
	case "validate":
		return runValidate(args[1:])
	case "diff":
		return runDiff(args[1:])
	case "watch":
		return runWatch(args[1:])
	case "design":
		return runDesignWrapper(args[1:])
	case "test":
		return runTestWrapper(args[1:])
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
	fmt.Println("  wetwire-azure validate <arm-file>           Validate ARM template JSON")
	fmt.Println("  wetwire-azure diff [package-path] [flags]   Compare generated vs existing template")
	fmt.Println("  wetwire-azure watch [package-path] [flags]  Watch source files and auto-rebuild")
	fmt.Println("  wetwire-azure design [flags]                AI-assisted infrastructure generation")
	fmt.Println("  wetwire-azure test [flags]                  Run persona-based testing")
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
	fmt.Println("Options for diff:")
	fmt.Println("  --against <file>          Existing template to compare (required)")
	fmt.Println("  --semantic                Use semantic comparison (ignore formatting, key order)")
	fmt.Println("  --color                   Colorized output (default: true if terminal)")
	fmt.Println()
	fmt.Println("Options for watch:")
	fmt.Println("  -o, --output <file>       Output file path")
	fmt.Println("  --interval <duration>     Polling interval (default: 500ms)")
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
		// Collect unique files with issues
		filesToFix := make(map[string]bool)
		for _, result := range results {
			filesToFix[result.File] = true
		}

		// Apply fixes using fixable rules
		fixedCount := 0
		for file := range filesToFix {
			for _, rule := range linter.AllRules() {
				if fixable, ok := rule.(linter.FixableRule); ok && fixable.CanFix() {
					// Check if this rule has issues for this file
					ruleResults, err := rule.Check(file)
					if err != nil {
						continue
					}
					if len(ruleResults) > 0 {
						fixed, err := fixable.Fix(file)
						if err != nil {
							fmt.Fprintf(os.Stderr, "Error fixing %s with %s: %v\n", file, rule.ID(), err)
							continue
						}
						if err := os.WriteFile(file, []byte(fixed), 0644); err != nil {
							fmt.Fprintf(os.Stderr, "Error writing fixed file %s: %v\n", file, err)
							continue
						}
						fixedCount++
						fmt.Printf("Fixed %s with %s\n", filepath.Base(file), rule.ID())
					}
				}
			}
		}
		if fixedCount > 0 {
			fmt.Printf("\nApplied %d fix(es)\n", fixedCount)
		} else {
			fmt.Println("\nNo auto-fixes available for these issues")
		}
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
func runValidate(args []string) int {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return ExitInvalidArgument
	}

	// Require input file
	if fs.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Error: ARM template file is required\n")
		printUsage()
		return ExitInvalidArgument
	}

	inputFile := fs.Arg(0)

	// Validate the template
	v := validator.NewValidator()
	results, err := v.ValidateFile(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Validation failed: %v\n", err)
		return ExitBuildError
	}

	// Print results
	if len(results) == 0 {
		fmt.Println("Template is valid.")
		return ExitSuccess
	}

	// Group results by severity
	errorCount := 0
	warningCount := 0
	infoCount := 0

	for _, result := range results {
		fmt.Println(result.String())
		switch result.Severity {
		case validator.SeverityError:
			errorCount++
		case validator.SeverityWarning:
			warningCount++
		case validator.SeverityInfo:
			infoCount++
		}
	}

	fmt.Println()
	fmt.Printf("Found %d error(s), %d warning(s), %d info(s)\n", errorCount, warningCount, infoCount)

	// Exit with error code if errors were found
	if errorCount > 0 {
		return ExitBuildError
	}

	return ExitSuccess
}

// ANSI color codes for diff output
const (
	colorReset = "\033[0m"
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	colorCyan  = "\033[36m"
)

// runDiff executes the diff command and returns an exit code
func runDiff(args []string) int {
	fs := flag.NewFlagSet("diff", flag.ContinueOnError)

	var againstFile string
	var semantic bool
	var color bool

	fs.StringVar(&againstFile, "against", "", "Existing template to compare")
	fs.BoolVar(&semantic, "semantic", false, "Use semantic comparison")
	fs.BoolVar(&color, "color", isTerminal(), "Colorized output")

	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return ExitInvalidArgument
	}

	// Require --against flag
	if againstFile == "" {
		fmt.Fprintf(os.Stderr, "Error: --against flag is required\n")
		printUsage()
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

	// Check if source path exists
	if _, err := os.Stat(absPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return ExitBuildError
	}

	// Read existing template
	existingData, err := os.ReadFile(againstFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading against file: %v\n", err)
		return ExitBuildError
	}

	// Validate existing template is valid JSON
	var existingJSON interface{}
	if err := json.Unmarshal(existingData, &existingJSON); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing against file JSON: %v\n", err)
		return ExitBuildError
	}

	// Build the new template
	newTemplateJSON, _, err := buildTemplate(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Build failed: %v\n", err)
		return ExitBuildError
	}

	// Parse new template
	var newJSON interface{}
	if err := json.Unmarshal([]byte(newTemplateJSON), &newJSON); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing generated template: %v\n", err)
		return ExitBuildError
	}

	// Perform diff
	if semantic {
		return performSemanticDiff(existingJSON, newJSON, color)
	}
	return performTextDiff(string(existingData), newTemplateJSON, color)
}

// isTerminal checks if stdout is a terminal
func isTerminal() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

// performTextDiff performs a line-by-line text diff
func performTextDiff(existing, generated string, useColor bool) int {
	existingLines := strings.Split(strings.TrimSpace(existing), "\n")
	generatedLines := strings.Split(strings.TrimSpace(generated), "\n")

	// Simple line-by-line diff using LCS algorithm
	diff := computeDiff(existingLines, generatedLines)

	if len(diff) == 0 {
		fmt.Println("No differences found.")
		return ExitSuccess
	}

	// Print diff
	for _, line := range diff {
		if useColor {
			switch {
			case strings.HasPrefix(line, "+"):
				fmt.Printf("%s%s%s\n", colorGreen, line, colorReset)
			case strings.HasPrefix(line, "-"):
				fmt.Printf("%s%s%s\n", colorRed, line, colorReset)
			case strings.HasPrefix(line, "@"):
				fmt.Printf("%s%s%s\n", colorCyan, line, colorReset)
			default:
				fmt.Println(line)
			}
		} else {
			fmt.Println(line)
		}
	}

	return ExitSuccess
}

// computeDiff computes a unified diff between two sets of lines
func computeDiff(a, b []string) []string {
	// Use a simple LCS-based diff algorithm
	lcs := longestCommonSubsequence(a, b)
	var result []string

	i, j := 0, 0
	for _, match := range lcs {
		// Lines removed from a
		for i < match.aIndex {
			result = append(result, "- "+a[i])
			i++
		}
		// Lines added to b
		for j < match.bIndex {
			result = append(result, "+ "+b[j])
			j++
		}
		// Common line
		result = append(result, "  "+a[i])
		i++
		j++
	}

	// Remaining lines in a (removed)
	for i < len(a) {
		result = append(result, "- "+a[i])
		i++
	}
	// Remaining lines in b (added)
	for j < len(b) {
		result = append(result, "+ "+b[j])
		j++
	}

	// Filter out common lines if only showing changes
	var changes []string
	for _, line := range result {
		if strings.HasPrefix(line, "+") || strings.HasPrefix(line, "-") {
			changes = append(changes, line)
		}
	}

	return changes
}

// lcsMatch represents a matching line in LCS
type lcsMatch struct {
	aIndex int
	bIndex int
}

// longestCommonSubsequence finds the LCS of two string slices
func longestCommonSubsequence(a, b []string) []lcsMatch {
	m, n := len(a), len(b)
	if m == 0 || n == 0 {
		return nil
	}

	// Build DP table
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				if dp[i-1][j] > dp[i][j-1] {
					dp[i][j] = dp[i-1][j]
				} else {
					dp[i][j] = dp[i][j-1]
				}
			}
		}
	}

	// Backtrack to find LCS
	var result []lcsMatch
	i, j := m, n
	for i > 0 && j > 0 {
		if a[i-1] == b[j-1] {
			result = append([]lcsMatch{{aIndex: i - 1, bIndex: j - 1}}, result...)
			i--
			j--
		} else if dp[i-1][j] > dp[i][j-1] {
			i--
		} else {
			j--
		}
	}

	return result
}

// performSemanticDiff performs a semantic comparison (ignoring formatting and key order)
func performSemanticDiff(existing, generated interface{}, useColor bool) int {
	// Normalize both JSON structures
	normalizedExisting := normalizeJSON(existing)
	normalizedGenerated := normalizeJSON(generated)

	// Compare normalized JSON
	existingStr, _ := json.MarshalIndent(normalizedExisting, "", "  ")
	generatedStr, _ := json.MarshalIndent(normalizedGenerated, "", "  ")

	if string(existingStr) == string(generatedStr) {
		fmt.Println("No differences found (semantic comparison).")
		return ExitSuccess
	}

	// Show the diff with normalized versions
	return performTextDiff(string(existingStr), string(generatedStr), useColor)
}

// normalizeJSON recursively normalizes a JSON structure (sorts map keys)
func normalizeJSON(v interface{}) interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		for k, v := range val {
			result[k] = normalizeJSON(v)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(val))
		for i, v := range val {
			result[i] = normalizeJSON(v)
		}
		return result
	default:
		return v
	}
}

// runWatch executes the watch command and returns an exit code
func runWatch(args []string) int {
	fs := flag.NewFlagSet("watch", flag.ContinueOnError)

	var outputFile string
	var outputFileLong string
	var interval string
	var testRun bool

	fs.StringVar(&outputFile, "o", "", "Output file path")
	fs.StringVar(&outputFileLong, "output", "", "Output file path")
	fs.StringVar(&interval, "interval", "500ms", "Polling interval")
	fs.BoolVar(&testRun, "test-run", false, "Run once and exit (for testing)")

	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return ExitInvalidArgument
	}

	// Use long form if short form not specified
	if outputFile == "" && outputFileLong != "" {
		outputFile = outputFileLong
	}

	// Parse interval
	pollInterval, err := time.ParseDuration(interval)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid interval '%s': %v\n", interval, err)
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

	// Check if source path exists
	if _, err := os.Stat(absPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return ExitBuildError
	}

	fmt.Printf("Watching %s for changes (interval: %s)\n", absPath, pollInterval)

	// Initial build
	if err := doBuild(absPath, outputFile); err != nil {
		fmt.Fprintf(os.Stderr, "Build error: %v\n", err)
	} else {
		fmt.Println("Build successful")
	}

	// For test mode, exit after initial build
	if testRun {
		return ExitSuccess
	}

	// Watch loop with debouncing
	return watchLoop(absPath, outputFile, pollInterval)
}

// doBuild performs a build and writes output
func doBuild(srcDir, outputFile string) error {
	templateJSON, _, err := buildTemplate(srcDir)
	if err != nil {
		return err
	}

	if outputFile != "" {
		if err := os.WriteFile(outputFile, []byte(templateJSON), 0644); err != nil {
			return fmt.Errorf("error writing output file: %w", err)
		}
	}

	return nil
}

// watchLoop runs the main watch loop with polling and debouncing
func watchLoop(srcDir, outputFile string, pollInterval time.Duration) int {
	const debounceDelay = 300 * time.Millisecond

	// Get initial file states
	lastModTimes := getFileModTimes(srcDir)
	var debounceTimer *time.Timer

	for {
		time.Sleep(pollInterval)

		// Check for changes
		currentModTimes := getFileModTimes(srcDir)
		if hasChanges(lastModTimes, currentModTimes) {
			lastModTimes = currentModTimes

			// Debounce: cancel previous timer and start a new one
			if debounceTimer != nil {
				debounceTimer.Stop()
			}

			debounceTimer = time.AfterFunc(debounceDelay, func() {
				fmt.Println("Changes detected, rebuilding...")
				if err := doBuild(srcDir, outputFile); err != nil {
					fmt.Fprintf(os.Stderr, "Build error: %v\n", err)
				} else {
					fmt.Println("Build successful")
				}
			})
		}
	}
}

// getFileModTimes returns a map of file paths to modification times
func getFileModTimes(dir string) map[string]time.Time {
	result := make(map[string]time.Time)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip files we can't access
		}
		if info.IsDir() {
			return nil
		}
		// Only watch Go files
		if strings.HasSuffix(path, ".go") {
			result[path] = info.ModTime()
		}
		return nil
	})
	if err != nil {
		return result
	}

	return result
}

// hasChanges checks if any files have been modified, added, or removed
func hasChanges(old, new map[string]time.Time) bool {
	// Check for new or modified files
	for path, newTime := range new {
		oldTime, exists := old[path]
		if !exists || !oldTime.Equal(newTime) {
			return true
		}
	}

	// Check for removed files
	for path := range old {
		if _, exists := new[path]; !exists {
			return true
		}
	}

	return false
}

// runTest executes the test command for persona-based testing

// runTestWrapper wraps the new cobra-based test command for compatibility with the flag-based CLI
func runTestWrapper(args []string) int {
	cmd := newTestCmd()
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		return ExitBuildError
	}
	return ExitSuccess
}

// runDesignWrapper wraps the new cobra-based design command for compatibility with the flag-based CLI
func runDesignWrapper(args []string) int {
	cmd := newDesignCmd()
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		return ExitBuildError
	}
	return ExitSuccess
}
