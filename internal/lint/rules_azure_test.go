package lint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestWAZ301HTTPSRequired tests the HTTPS-only requirement for storage accounts
func TestWAZ301HTTPSRequired(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		content     string
		expectIssue bool
	}{
		{
			name: "https disabled",
			content: `package main

var MyStorage = struct {
	EnableHTTPSTrafficOnly bool
}{
	EnableHTTPSTrafficOnly: false,
}
`,
			expectIssue: true,
		},
		{
			name: "https enabled",
			content: `package main

var MyStorage = struct {
	EnableHTTPSTrafficOnly bool
}{
	EnableHTTPSTrafficOnly: true,
}
`,
			expectIssue: false,
		},
		{
			name: "no https setting",
			content: `package main

var MyStorage = struct {
	Name string
}{
	Name: "test",
}
`,
			expectIssue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, "test_"+strings.ReplaceAll(tt.name, " ", "_")+".go")
			if err := os.WriteFile(testFile, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}

			rule := &WAZ301{}
			results, err := rule.Check(testFile)
			if err != nil {
				t.Fatalf("Check() error: %v", err)
			}

			if tt.expectIssue && len(results) == 0 {
				t.Error("expected lint issue but got none")
			}
			if !tt.expectIssue && len(results) > 0 {
				t.Errorf("expected no lint issues but got %d", len(results))
			}

			// Verify rule metadata
			if rule.ID() != "WAZ301" {
				t.Errorf("expected ID WAZ301, got %s", rule.ID())
			}
			if rule.Severity() != SeverityWarning {
				t.Errorf("expected SeverityWarning, got %s", rule.Severity())
			}
		})
	}
}

// TestWAZ302PermissiveNSGRules tests detection of overly permissive NSG rules
func TestWAZ302PermissiveNSGRules(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		content     string
		expectIssue bool
	}{
		{
			name: "wildcard source address",
			content: `package main

var MyRule = struct {
	SourceAddressPrefix string
}{
	SourceAddressPrefix: "0.0.0.0/0",
}
`,
			expectIssue: true,
		},
		{
			name: "any source address",
			content: `package main

var MyRule = struct {
	SourceAddressPrefix string
}{
	SourceAddressPrefix: "*",
}
`,
			expectIssue: true,
		},
		{
			name: "specific source address",
			content: `package main

var MyRule = struct {
	SourceAddressPrefix string
}{
	SourceAddressPrefix: "10.0.0.0/8",
}
`,
			expectIssue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, "test_"+strings.ReplaceAll(tt.name, " ", "_")+".go")
			if err := os.WriteFile(testFile, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}

			rule := &WAZ302{}
			results, err := rule.Check(testFile)
			if err != nil {
				t.Fatalf("Check() error: %v", err)
			}

			if tt.expectIssue && len(results) == 0 {
				t.Error("expected lint issue but got none")
			}
			if !tt.expectIssue && len(results) > 0 {
				t.Errorf("expected no lint issues but got %d", len(results))
			}

			// Verify rule metadata
			if rule.ID() != "WAZ302" {
				t.Errorf("expected ID WAZ302, got %s", rule.ID())
			}
			if rule.Severity() != SeverityWarning {
				t.Errorf("expected SeverityWarning, got %s", rule.Severity())
			}
		})
	}
}

// TestWAZ303RequireTags tests detection of resources without tags
func TestWAZ303RequireTags(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		content     string
		expectIssue bool
	}{
		{
			name: "resource without tags",
			content: `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{
	Name:     "test",
	Location: "eastus",
}
`,
			expectIssue: true,
		},
		{
			name: "resource with tags",
			content: `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{
	Name:     "test",
	Location: "eastus",
	Tags:     map[string]string{"env": "prod"},
}
`,
			expectIssue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, "test_"+strings.ReplaceAll(tt.name, " ", "_")+".go")
			if err := os.WriteFile(testFile, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}

			rule := &WAZ303{}
			results, err := rule.Check(testFile)
			if err != nil {
				t.Fatalf("Check() error: %v", err)
			}

			if tt.expectIssue && len(results) == 0 {
				t.Error("expected lint issue but got none")
			}
			if !tt.expectIssue && len(results) > 0 {
				t.Errorf("expected no lint issues but got %d", len(results))
			}

			// Verify rule metadata
			if rule.ID() != "WAZ303" {
				t.Errorf("expected ID WAZ303, got %s", rule.ID())
			}
			if rule.Severity() != SeverityWarning {
				t.Errorf("expected SeverityWarning, got %s", rule.Severity())
			}
		})
	}
}

// TestWAZ304DeprecatedAPIVersion tests detection of deprecated API versions
func TestWAZ304DeprecatedAPIVersion(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		content     string
		expectIssue bool
	}{
		{
			name: "deprecated api version",
			content: `package main

var MyResource = struct {
	APIVersion string
}{
	APIVersion: "2019-01-01",
}
`,
			expectIssue: true,
		},
		{
			name: "current api version",
			content: `package main

var MyResource = struct {
	APIVersion string
}{
	APIVersion: "2023-01-01",
}
`,
			expectIssue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, "test_"+strings.ReplaceAll(tt.name, " ", "_")+".go")
			if err := os.WriteFile(testFile, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}

			rule := &WAZ304{}
			results, err := rule.Check(testFile)
			if err != nil {
				t.Fatalf("Check() error: %v", err)
			}

			if tt.expectIssue && len(results) == 0 {
				t.Error("expected lint issue but got none")
			}
			if !tt.expectIssue && len(results) > 0 {
				t.Errorf("expected no lint issues but got %d", len(results))
			}

			// Verify rule metadata
			if rule.ID() != "WAZ304" {
				t.Errorf("expected ID WAZ304, got %s", rule.ID())
			}
			if rule.Severity() != SeverityWarning {
				t.Errorf("expected SeverityWarning, got %s", rule.Severity())
			}
		})
	}
}
