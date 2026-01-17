// Package main provides the MCP server implementation for wetwire-azure.
//
// The MCP server exposes wetwire-azure tools via the Model Context Protocol,
// providing wetwire_init, wetwire_lint, wetwire_build, wetwire_validate,
// wetwire_list, and wetwire_graph tools.
//
// Usage:
//
//	wetwire-azure mcp  # Runs on stdio transport
package main

import (
	"context"

	coredomain "github.com/lex00/wetwire-core-go/domain"
	"github.com/lex00/wetwire-azure-go/domain"
	"github.com/spf13/cobra"
)

// mcpCmd is the mcp subcommand.
var mcpCmd = newMCPCmd()

// newMCPCmd creates the mcp subcommand.
func newMCPCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mcp",
		Short: "Run MCP server for wetwire-azure tools",
		Long: `Run an MCP (Model Context Protocol) server that exposes wetwire-azure tools.

This command starts an MCP server on stdio transport, providing tools for:
- Initializing projects (wetwire_init)
- Building ARM templates (wetwire_build)
- Linting code (wetwire_lint)
- Validating templates (wetwire_validate)
- Listing resources (wetwire_list)
- Generating dependency graphs (wetwire_graph)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMCPServer()
		},
	}
}

// runMCPServer starts the MCP server on stdio transport.
func runMCPServer() error {
	// Use domain.BuildMCPServer for automatic MCP server generation
	server := coredomain.BuildMCPServer(&domain.AzureDomain{})

	// Start the server on stdio transport
	return server.Start(context.Background())
}
