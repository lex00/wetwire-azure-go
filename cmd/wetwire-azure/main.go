package main

import (
	"fmt"
	"os"

	"github.com/lex00/wetwire-azure-go/domain"
)

// Version information set via ldflags
var version = "dev"

func main() {
	// Set domain version from ldflags
	domain.Version = version

	d := &domain.AzureDomain{}
	cmd := domain.CreateRootCommand(d)

	// Add MCP command
	cmd.AddCommand(mcpCmd)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
