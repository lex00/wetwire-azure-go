package main

import (
	"os"

	"github.com/lex00/wetwire-azure-go/domain"
)

func main() {
	d := &domain.AzureDomain{}
	cmd := domain.CreateRootCommand(d)

	// Add MCP command
	cmd.AddCommand(mcpCmd)

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
