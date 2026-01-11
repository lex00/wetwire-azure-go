package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/lex00/wetwire-azure-go/internal/linter"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "lint":
		runLint(os.Args[2:])
	case "help", "-h", "--help":
		printUsage()
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("wetwire-azure - Azure ARM/Bicep template synthesis")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  wetwire-azure lint [path]     Lint infrastructure code")
	fmt.Println("  wetwire-azure help             Show this help message")
	fmt.Println()
	fmt.Println("Options for lint:")
	fmt.Println("  --fix                          Auto-fix issues where possible (not yet implemented)")
}

func runLint(args []string) {
	fs := flag.NewFlagSet("lint", flag.ExitOnError)
	fixFlag := fs.Bool("fix", false, "Auto-fix issues where possible")
	fs.Parse(args)

	// Default to current directory if no path provided
	path := "."
	if fs.NArg() > 0 {
		path = fs.Arg(0)
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving path: %v\n", err)
		os.Exit(1)
	}

	// Check if path exists
	info, err := os.Stat(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
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
		os.Exit(1)
	}

	// Print results
	if len(results) == 0 {
		fmt.Println("No issues found.")
		os.Exit(0)
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
		os.Exit(1)
	}
}
