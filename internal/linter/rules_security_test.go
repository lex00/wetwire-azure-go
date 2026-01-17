package linter

import (
	"os"
	"path/filepath"
	"testing"
)

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
