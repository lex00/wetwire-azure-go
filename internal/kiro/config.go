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
// The full prompt is embedded in configs/wetwire-azure-runner.json.
const AgentPrompt = "You are an Azure infrastructure expert using wetwire-azure to create ARM templates."

// MCPCommand is the command to run the MCP server.
const MCPCommand = "wetwire-azure"

// MCPArgs are the arguments to pass to the MCP command.
var MCPArgs = []string{"mcp"}

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
