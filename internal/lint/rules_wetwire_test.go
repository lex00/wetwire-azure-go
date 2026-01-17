package lint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestWAZ020TopLevelResourceDeclarations tests the top-level resource declarations rule
func TestWAZ020TopLevelResourceDeclarations(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		content     string
		expectIssue bool
	}{
		{
			name: "resource inside function",
			content: `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

func createResources() {
	var MyStorage = storage.StorageAccount{
		Name:     "mystorageaccount",
		Location: "eastus",
	}
	_ = MyStorage
}
`,
			expectIssue: true,
		},
		{
			name: "resource at package level",
			content: `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{
	Name:     "mystorageaccount",
	Location: "eastus",
}
`,
			expectIssue: false,
		},
		{
			name: "non-resource variable inside function",
			content: `package main

func doSomething() {
	var count = 10
	_ = count
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

			rule := &WAZ020{}
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
			if rule.ID() != "WAZ020" {
				t.Errorf("expected ID WAZ020, got %s", rule.ID())
			}
			if rule.Severity() != SeverityWarning {
				t.Errorf("expected SeverityWarning, got %s", rule.Severity())
			}
		})
	}
}

// TestWAZ021DeeplyNestedInlineStructs tests the deeply nested inline structs rule
func TestWAZ021DeeplyNestedInlineStructs(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		content     string
		expectIssue bool
	}{
		{
			name: "deeply nested struct (depth > 2)",
			content: `package main

import "github.com/lex00/wetwire-azure-go/resources/compute"

var MyVM = compute.VirtualMachine{
	Properties: &compute.VirtualMachineProperties{
		StorageProfile: &compute.StorageProfile{
			OsDisk: &compute.OsDisk{
				ManagedDisk: &compute.ManagedDiskParameters{
					StorageAccountType: "Premium_LRS",
				},
			},
		},
	},
}
`,
			expectIssue: true,
		},
		{
			name: "shallow nesting (depth <= 2)",
			content: `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{
	Properties: &storage.StorageAccountProperties{
		AccessTier: "Hot",
	},
}
`,
			expectIssue: false,
		},
		{
			name: "extracted nested configuration",
			content: `package main

import "github.com/lex00/wetwire-azure-go/resources/compute"

var ManagedDisk = compute.ManagedDiskParameters{
	StorageAccountType: "Premium_LRS",
}

var OsDisk = compute.OsDisk{
	ManagedDisk: &ManagedDisk,
}

var MyVM = compute.VirtualMachine{
	Properties: &compute.VirtualMachineProperties{
		StorageProfile: &compute.StorageProfile{
			OsDisk: &OsDisk,
		},
	},
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

			rule := &WAZ021{}
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
			if rule.ID() != "WAZ021" {
				t.Errorf("expected ID WAZ021, got %s", rule.ID())
			}
			if rule.Severity() != SeverityWarning {
				t.Errorf("expected SeverityWarning, got %s", rule.Severity())
			}
		})
	}
}

// TestWAZ022ExtractSKUProfileConfigs tests the SKU/Profile config extraction rule
func TestWAZ022ExtractSKUProfileConfigs(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		content     string
		expectIssue bool
	}{
		{
			name: "inline SKU config",
			content: `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{
	Name: "test",
	Sku: &storage.Sku{
		Name: "Standard_LRS",
		Tier: "Standard",
	},
}
`,
			expectIssue: true,
		},
		{
			name: "inline StorageProfile config",
			content: `package main

import "github.com/lex00/wetwire-azure-go/resources/compute"

var MyVM = compute.VirtualMachine{
	StorageProfile: &compute.StorageProfile{
		ImageReference: "Ubuntu",
	},
}
`,
			expectIssue: true,
		},
		{
			name: "extracted SKU config",
			content: `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var StandardSku = storage.Sku{
	Name: "Standard_LRS",
	Tier: "Standard",
}

var MyStorage = storage.StorageAccount{
	Name: "test",
	Sku:  &StandardSku,
}
`,
			expectIssue: false,
		},
		{
			name: "no SKU or Profile field",
			content: `package main

import "github.com/lex00/wetwire-azure-go/resources/storage"

var MyStorage = storage.StorageAccount{
	Name:     "test",
	Location: "eastus",
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

			rule := &WAZ022{}
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
			if rule.ID() != "WAZ022" {
				t.Errorf("expected ID WAZ022, got %s", rule.ID())
			}
			if rule.Severity() != SeverityInfo {
				t.Errorf("expected SeverityInfo, got %s", rule.Severity())
			}
		})
	}
}
