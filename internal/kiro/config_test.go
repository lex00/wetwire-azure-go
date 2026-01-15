package kiro

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewConfig_IncludesWorkDir(t *testing.T) {
	// Test that NewConfig includes the working directory
	// This is required for proper Kiro provider execution context
	// See: https://github.com/lex00/wetwire-core-go/pull/72

	origWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)

	// Create a temporary directory and change to it
	tmpDir := t.TempDir()
	testDir := filepath.Join(tmpDir, "testproject")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(testDir); err != nil {
		t.Fatal(err)
	}

	config := NewConfig()

	if config.WorkDir == "" {
		t.Error("Config.WorkDir should not be empty")
	}

	// WorkDir should be the current working directory
	expectedDir, _ := os.Getwd()
	if config.WorkDir != expectedDir {
		t.Errorf("Config.WorkDir = %q, want %q", config.WorkDir, expectedDir)
	}

	// Verify other required fields are set
	if config.AgentName == "" {
		t.Error("Config.AgentName should not be empty")
	}
	if config.MCPCommand == "" {
		t.Error("Config.MCPCommand should not be empty")
	}
}
