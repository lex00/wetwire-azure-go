package kiro

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestInstall_HasToolsArray(t *testing.T) {
	// Test that the generated config includes tools array
	// Required for kiro to enable MCP tool usage
	// See: https://github.com/aws/amazon-q-developer-cli/issues/2640

	tmpDir := t.TempDir()
	projectDir := filepath.Join(tmpDir, "project")
	homeDir := filepath.Join(tmpDir, "home")

	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(homeDir, 0755); err != nil {
		t.Fatal(err)
	}

	origHome := os.Getenv("HOME")
	os.Setenv("HOME", homeDir)
	defer os.Setenv("HOME", origHome)

	origWd, _ := os.Getwd()
	os.Chdir(projectDir)
	defer os.Chdir(origWd)

	if err := EnsureInstalledWithForce(true); err != nil {
		t.Fatalf("EnsureInstalledWithForce failed: %v", err)
	}

	agentPath := filepath.Join(homeDir, ".kiro", "agents", AgentName+".json")
	data, err := os.ReadFile(agentPath)
	if err != nil {
		t.Fatalf("failed to read agent config: %v", err)
	}

	var agent map[string]any
	if err := json.Unmarshal(data, &agent); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	tools, ok := agent["tools"].([]any)
	if !ok {
		t.Fatal("agent config must have 'tools' array - required for kiro MCP tool usage")
	}

	if len(tools) == 0 {
		t.Error("tools array must not be empty")
	}

	if len(tools) > 0 {
		tool, ok := tools[0].(string)
		if !ok || len(tool) == 0 || tool[0] != '@' {
			t.Errorf("tools should use @server_name format, got: %v", tools)
		}
	}
}
