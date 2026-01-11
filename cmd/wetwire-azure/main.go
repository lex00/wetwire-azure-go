package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/lex00/wetwire-azure-go/internal/discover"
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
	case "lint":
		return runLint(args[1:])
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
	fmt.Println("  wetwire-azure lint [path]                   Lint infrastructure code")
	fmt.Println("  wetwire-azure help                          Show this help message")
	fmt.Println()
	fmt.Println("Options for build:")
	fmt.Println("  -o, --output <file>       Output file path (default: stdout)")
	fmt.Println("  --format <format>         Output format: arm (default: arm)")
	fmt.Println("  --parameters-file <file>  Write parameters to separate file")
	fmt.Println()
	fmt.Println("Options for lint:")
	fmt.Println("  --fix                     Auto-fix issues where possible (not yet implemented)")
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
