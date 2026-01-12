package kiro

import (
	"context"
	"fmt"

	corekiro "github.com/lex00/wetwire-core-go/kiro"
)

// LaunchChat launches an interactive Kiro CLI chat session with the wetwire-azure agent.
// It connects stdin/stdout directly to the terminal for interactive use.
func LaunchChat(agentName, initialPrompt string) error {
	// Ensure latest config is installed
	if err := EnsureInstalledWithForce(true); err != nil {
		return fmt.Errorf("installing kiro config: %w", err)
	}

	// Use core kiro package to launch
	config := NewConfig()
	return corekiro.Launch(context.Background(), config, initialPrompt)
}

// Launch launches an interactive Kiro CLI chat session with the wetwire-azure agent.
// This is the simplified version that matches the GitLab pattern.
func Launch(initialPrompt string) error {
	// Ensure latest config is installed
	if err := EnsureInstalledWithForce(true); err != nil {
		return fmt.Errorf("installing kiro config: %w", err)
	}

	// Use core kiro package to launch
	config := NewConfig()
	return corekiro.Launch(context.Background(), config, initialPrompt)
}
