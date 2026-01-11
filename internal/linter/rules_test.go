package linter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestWAZ001LocationConstants tests the location constants rule
func TestWAZ001LocationConstants(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		content     string
		expectIssue bool
	}{
		{
			name: "invalid location format",
			content: `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{
	Location: "East US",
}
`,
			expectIssue: true,
		},
		{
			name: "valid location format",
			content: `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{
	Location: "eastus",
}
`,
			expectIssue: false,
		},
		{
			name: "location from function call",
			content: `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{
	Location: ResourceGroup().Location,
}
`,
			expectIssue: false,
		},
		{
			name: "no location field",
			content: `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{
	Name: "test",
}
`,
			expectIssue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, "test_"+tt.name+".go")
			if err := os.WriteFile(testFile, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}

			rule := &WAZ001{}
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
			if rule.ID() != "WAZ001" {
				t.Errorf("expected ID WAZ001, got %s", rule.ID())
			}
			if rule.Severity() != SeverityWarning {
				t.Errorf("expected SeverityWarning, got %s", rule.Severity())
			}
		})
	}
}

// TestWAZ002DirectReferences tests the direct references rule
func TestWAZ002DirectReferences(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		content     string
		expectIssue bool
	}{
		{
			name: "uses resourceId function",
			content: `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var id = resourceId("Microsoft.Storage/storageAccounts", "myaccount")
`,
			expectIssue: true,
		},
		{
			name: "uses direct reference",
			content: `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{Name: "test"}
var id = MyStorage.Id
`,
			expectIssue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, "test_"+tt.name+".go")
			if err := os.WriteFile(testFile, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}

			rule := &WAZ002{}
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
		})
	}
}

// TestWAZ003ExtractNestedConfigs tests the extract nested configurations rule
func TestWAZ003ExtractNestedConfigs(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		content     string
		expectIssue bool
	}{
		{
			name: "deeply nested configuration",
			content: `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{
	Properties: &storage.StorageAccountProperties{
		NetworkRuleSet: &storage.NetworkRuleSet{
			DefaultAction: "Deny",
			IPRules: []storage.IPRule{
				{Value: "10.0.0.1"},
			},
		},
	},
}
`,
			expectIssue: true,
		},
		{
			name: "extracted configuration",
			content: `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var NetworkRules = storage.NetworkRuleSet{
	DefaultAction: "Deny",
}

var MyStorage = storage.StorageAccount{
	Properties: &storage.StorageAccountProperties{
		NetworkRuleSet: &NetworkRules,
	},
}
`,
			expectIssue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, "test_"+tt.name+".go")
			if err := os.WriteFile(testFile, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}

			rule := &WAZ003{}
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
		})
	}
}

// TestWAZ004DuplicateResourceNames tests the duplicate resource names rule
func TestWAZ004DuplicateResourceNames(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		content     string
		expectIssue bool
	}{
		{
			name: "duplicate variable names",
			content: `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{Name: "storage1"}
var MyStorage = storage.StorageAccount{Name: "storage2"}
`,
			expectIssue: true,
		},
		{
			name: "unique variable names",
			content: `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var DataStorage = storage.StorageAccount{Name: "datastorage"}
var LogsStorage = storage.StorageAccount{Name: "logsstorage"}
`,
			expectIssue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, "test_"+tt.name+".go")
			if err := os.WriteFile(testFile, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}

			rule := &WAZ004{}
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
		})
	}
}

// TestWAZ005CircularDependencies tests the circular dependencies rule
func TestWAZ005CircularDependencies(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		content     string
		expectIssue bool
	}{
		{
			name: "circular dependency",
			content: `package main

import (
	"github.com/lex00/wetwire-azure-go/resources/storage"
	"github.com/lex00/wetwire-azure-go/resources/compute"
)

var VM = compute.VirtualMachine{
	Name: Storage.Name,
}

var Storage = storage.StorageAccount{
	Name: VM.Name,
}
`,
			expectIssue: true,
		},
		{
			name: "no circular dependency",
			content: `package main

import (
	"github.com/lex00/wetwire-azure-go/resources/storage"
	"github.com/lex00/wetwire-azure-go/resources/compute"
)

var Storage = storage.StorageAccount{
	Name: "mystorage",
}

var VM = compute.VirtualMachine{
	Name: Storage.Name,
}
`,
			expectIssue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, "test_"+tt.name+".go")
			if err := os.WriteFile(testFile, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}

			rule := &WAZ005{}
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
		})
	}
}

// TestWAZ006SecretDetection tests the secret detection rule
func TestWAZ006SecretDetection(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		content     string
		expectIssue bool
	}{
		{
			name: "AWS access key pattern",
			content: `package main

var accessKey = "AKIAIOSFODNN7EXAMPLE"
`,
			expectIssue: true,
		},
		{
			name: "GitHub personal access token",
			content: `package main

var token = "ghp_1234567890abcdefghijklmnopqrstuvwxyz"
`,
			expectIssue: true,
		},
		{
			name: "Azure storage account key pattern",
			content: `package main

var storageKey = "AccountKey=abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstu=="
`,
			expectIssue: true,
		},
		{
			name: "password in string",
			content: `package main

var config = "password=MySecretPassword123"
`,
			expectIssue: true,
		},
		{
			name: "safe configuration",
			content: `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{
	Name:     "mystorage",
	Location: "eastus",
}
`,
			expectIssue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, "test_"+tt.name+".go")
			if err := os.WriteFile(testFile, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}

			rule := &WAZ006{}
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
			if rule.ID() != "WAZ006" {
				t.Errorf("expected ID WAZ006, got %s", rule.ID())
			}
			if rule.Severity() != SeverityError {
				t.Errorf("expected SeverityError, got %s", rule.Severity())
			}
		})
	}
}

// TestWAZ007SensitivePaths tests the sensitive file paths rule
func TestWAZ007SensitivePaths(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		content     string
		expectIssue bool
	}{
		{
			name: "env file path",
			content: `package main

var configPath = ".env"
`,
			expectIssue: true,
		},
		{
			name: "private key path",
			content: `package main

var keyPath = "/path/to/private.key"
`,
			expectIssue: true,
		},
		{
			name: "pem file path",
			content: `package main

var certPath = "server.pem"
`,
			expectIssue: true,
		},
		{
			name: "safe path",
			content: `package main

var configPath = "config.yaml"
`,
			expectIssue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, "test_"+tt.name+".go")
			if err := os.WriteFile(testFile, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}

			rule := &WAZ007{}
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
		})
	}
}

// TestWAZ008InsecureDefaults tests the insecure defaults rule
func TestWAZ008InsecureDefaults(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		content     string
		expectIssue bool
	}{
		{
			name: "http protocol",
			content: `package main

var endpoint = "http://storage.azure.com"
`,
			expectIssue: true,
		},
		{
			name: "public access enabled",
			content: `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{
	Properties: &storage.StorageAccountProperties{
		AllowBlobPublicAccess: true,
	},
}
`,
			expectIssue: true,
		},
		{
			name: "https protocol",
			content: `package main

var endpoint = "https://storage.azure.com"
`,
			expectIssue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, "test_"+tt.name+".go")
			if err := os.WriteFile(testFile, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}

			rule := &WAZ008{}
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
		})
	}
}

// TestAllRules verifies that all rules are registered
func TestAllRules(t *testing.T) {
	rules := AllRules()

	if len(rules) < 8 {
		t.Errorf("expected at least 8 rules, got %d", len(rules))
	}

	// Check for specific rules
	ruleIDs := make(map[string]bool)
	for _, rule := range rules {
		ruleIDs[rule.ID()] = true
	}

	expectedRules := []string{"WAZ001", "WAZ002", "WAZ003", "WAZ004", "WAZ005", "WAZ006", "WAZ007", "WAZ008"}
	for _, id := range expectedRules {
		if !ruleIDs[id] {
			t.Errorf("expected rule %s to be registered", id)
		}
	}
}

// TestWAZ001_AutoFix tests the location constants rule auto-fix capability
func TestWAZ001_AutoFix(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name           string
		content        string
		expectedResult string
	}{
		{
			name: "fix uppercase location",
			content: `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{
	Location: "East US",
}
`,
			expectedResult: `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{
	Location: "eastus",
}
`,
		},
		{
			name: "fix mixed case location",
			content: `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{
	Location: "West Europe",
}
`,
			expectedResult: `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{
	Location: "westeurope",
}
`,
		},
		{
			name: "preserve already lowercase location",
			content: `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{
	Location: "eastus",
}
`,
			expectedResult: `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{
	Location: "eastus",
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, "test_"+tt.name+".go")
			if err := os.WriteFile(testFile, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}

			rule := &WAZ001{}
			results, err := rule.Check(testFile)
			if err != nil {
				t.Fatalf("Check() error: %v", err)
			}

			// If there are issues, apply fixes
			if len(results) > 0 {
				// Verify rule supports fixing
				if !rule.CanFix() {
					t.Fatal("WAZ001 should support auto-fixing")
				}

				fixed, err := rule.Fix(testFile)
				if err != nil {
					t.Fatalf("Fix() error: %v", err)
				}

				if fixed != tt.expectedResult {
					t.Errorf("Fix() result mismatch.\nExpected:\n%s\nGot:\n%s", tt.expectedResult, fixed)
				}
			} else {
				// No issues, content should be unchanged
				content, err := os.ReadFile(testFile)
				if err != nil {
					t.Fatalf("ReadFile() error: %v", err)
				}
				if string(content) != tt.expectedResult {
					t.Errorf("Content mismatch.\nExpected:\n%s\nGot:\n%s", tt.expectedResult, string(content))
				}
			}
		})
	}
}

// TestFixableRules tests that fixable rules implement the FixableRule interface
func TestFixableRules(t *testing.T) {
	// WAZ001 should be fixable
	waz001 := &WAZ001{}
	if !waz001.CanFix() {
		t.Error("WAZ001 should support auto-fixing")
	}

	// Verify Fix method returns valid content
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")
	content := `package main

var MyStorage = struct{ Location string }{
	Location: "East US",
}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	fixed, err := waz001.Fix(testFile)
	if err != nil {
		t.Fatalf("Fix() error: %v", err)
	}

	if !strings.Contains(fixed, `"eastus"`) {
		t.Errorf("Fix() should normalize location to eastus, got: %s", fixed)
	}
}
