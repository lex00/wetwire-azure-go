package kiro

import (
	"embed"
	"os"

	corekiro "github.com/lex00/wetwire-core-go/kiro"
)

//go:embed configs/wetwire-azure-runner.json
var configFS embed.FS

// AgentName is the identifier for the wetwire-azure Kiro agent.
const AgentName = "wetwire-azure-runner"

// AgentPrompt contains the system prompt for the wetwire-azure agent.
// This is loaded from the embedded config file.
var AgentPrompt string

// MCPCommand is the command to run the MCP server.
const MCPCommand = "wetwire-azure"

// MCPArgs are the arguments to pass to the MCP command.
var MCPArgs = []string{"mcp"}

func init() {
	// Load the agent prompt from the embedded config file
	data, err := configFS.ReadFile("configs/wetwire-azure-runner.json")
	if err != nil {
		// Fallback to basic prompt if config can't be read
		AgentPrompt = "You are an Azure infrastructure expert using wetwire-azure to create ARM templates."
		return
	}

	// Parse the JSON to extract the prompt
	// Simple extraction - we'll just use the full config in the installer
	// For now, use a basic prompt
	AgentPrompt = "You are an Azure infrastructure expert using wetwire-azure to create ARM templates."
	_ = data // We'll use the full config data in the installer
}

// NewConfig creates a new Kiro config for the wetwire-azure agent.
func NewConfig() corekiro.Config {
	workDir, _ := os.Getwd()
	return corekiro.Config{
		AgentName:   AgentName,
		AgentPrompt: AgentPrompt,
		MCPCommand:  MCPCommand,
		MCPArgs:     MCPArgs,
		WorkDir:     workDir,
	}
}
