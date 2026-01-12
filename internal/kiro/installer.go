// Package kiro provides Kiro CLI integration for wetwire-azure.
//
// This package handles:
//   - Auto-installation of Kiro agent configuration
//   - Project-level MCP configuration
//   - Launching Kiro CLI chat sessions
//
// It builds on the infrastructure from github.com/lex00/wetwire-core-go/kiro.
package kiro

import (
	"fmt"

	corekiro "github.com/lex00/wetwire-core-go/kiro"
)

// Config wraps corekiro.Config with Azure-specific defaults.
type Config = corekiro.Config

// EnsureInstalled checks if Kiro configs are installed and installs them if needed.
// It installs:
//   - ~/.kiro/agents/wetwire-azure-runner.json (user-level agent config)
//   - .kiro/mcp.json (project-level MCP config)
//
// Existing files are not overwritten unless force is true.
func EnsureInstalled() error {
	return EnsureInstalledWithForce(false)
}

// EnsureInstalledWithForce installs Kiro configs, optionally overwriting existing ones.
// When force is true, configs are always reinstalled to ensure latest prompt is used.
func EnsureInstalledWithForce(force bool) error {
	// Use core kiro Install function
	config := NewConfig()
	if err := corekiro.Install(config); err != nil {
		return fmt.Errorf("installing kiro config: %w", err)
	}
	return nil
}

