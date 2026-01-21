package differ

import (
	"os"
	"path/filepath"
	"testing"

	coredomain "github.com/lex00/wetwire-core-go/domain"
)

func TestDiff_AddedResource(t *testing.T) {
	dir := t.TempDir()

	// Template 1: single resource
	t1 := filepath.Join(dir, "template1.json")
	writeJSON(t, t1, `{
		"$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"resources": [
			{
				"name": "storage1",
				"type": "Microsoft.Storage/storageAccounts",
				"apiVersion": "2021-04-01"
			}
		]
	}`)

	// Template 2: two resources
	t2 := filepath.Join(dir, "template2.json")
	writeJSON(t, t2, `{
		"$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"resources": [
			{
				"name": "storage1",
				"type": "Microsoft.Storage/storageAccounts",
				"apiVersion": "2021-04-01"
			},
			{
				"name": "storage2",
				"type": "Microsoft.Storage/storageAccounts",
				"apiVersion": "2021-04-01"
			}
		]
	}`)

	d := New()
	result, err := d.Diff(nil, t1, t2, coredomain.DiffOpts{})
	if err != nil {
		t.Fatalf("Diff failed: %v", err)
	}

	if result.Summary.Added != 1 {
		t.Errorf("expected 1 added, got %d", result.Summary.Added)
	}
	if result.Summary.Total != 1 {
		t.Errorf("expected 1 total, got %d", result.Summary.Total)
	}
}

func TestDiff_RemovedResource(t *testing.T) {
	dir := t.TempDir()

	// Template 1: two resources
	t1 := filepath.Join(dir, "template1.json")
	writeJSON(t, t1, `{
		"$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"resources": [
			{
				"name": "storage1",
				"type": "Microsoft.Storage/storageAccounts",
				"apiVersion": "2021-04-01"
			},
			{
				"name": "storage2",
				"type": "Microsoft.Storage/storageAccounts",
				"apiVersion": "2021-04-01"
			}
		]
	}`)

	// Template 2: single resource
	t2 := filepath.Join(dir, "template2.json")
	writeJSON(t, t2, `{
		"$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"resources": [
			{
				"name": "storage1",
				"type": "Microsoft.Storage/storageAccounts",
				"apiVersion": "2021-04-01"
			}
		]
	}`)

	d := New()
	result, err := d.Diff(nil, t1, t2, coredomain.DiffOpts{})
	if err != nil {
		t.Fatalf("Diff failed: %v", err)
	}

	if result.Summary.Removed != 1 {
		t.Errorf("expected 1 removed, got %d", result.Summary.Removed)
	}
}

func TestDiff_ModifiedResource(t *testing.T) {
	dir := t.TempDir()

	// Template 1
	t1 := filepath.Join(dir, "template1.json")
	writeJSON(t, t1, `{
		"$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"resources": [
			{
				"name": "storage1",
				"type": "Microsoft.Storage/storageAccounts",
				"apiVersion": "2021-04-01",
				"location": "eastus"
			}
		]
	}`)

	// Template 2: location changed
	t2 := filepath.Join(dir, "template2.json")
	writeJSON(t, t2, `{
		"$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"resources": [
			{
				"name": "storage1",
				"type": "Microsoft.Storage/storageAccounts",
				"apiVersion": "2021-04-01",
				"location": "westus"
			}
		]
	}`)

	d := New()
	result, err := d.Diff(nil, t1, t2, coredomain.DiffOpts{})
	if err != nil {
		t.Fatalf("Diff failed: %v", err)
	}

	if result.Summary.Modified != 1 {
		t.Errorf("expected 1 modified, got %d", result.Summary.Modified)
	}

	// Check that the change is captured
	if len(result.Entries) == 0 {
		t.Fatal("expected at least one entry")
	}

	entry := result.Entries[0]
	if entry.Action != "modified" {
		t.Errorf("expected action 'modified', got %q", entry.Action)
	}
	if len(entry.Changes) == 0 {
		t.Error("expected changes to be captured")
	}
}

func TestDiff_NoDifferences(t *testing.T) {
	dir := t.TempDir()

	template := `{
		"$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"resources": [
			{
				"name": "storage1",
				"type": "Microsoft.Storage/storageAccounts",
				"apiVersion": "2021-04-01"
			}
		]
	}`

	t1 := filepath.Join(dir, "template1.json")
	t2 := filepath.Join(dir, "template2.json")
	writeJSON(t, t1, template)
	writeJSON(t, t2, template)

	d := New()
	result, err := d.Diff(nil, t1, t2, coredomain.DiffOpts{})
	if err != nil {
		t.Fatalf("Diff failed: %v", err)
	}

	if result.Summary.Total != 0 {
		t.Errorf("expected 0 differences, got %d", result.Summary.Total)
	}
}

func TestDiff_YAMLSupport(t *testing.T) {
	dir := t.TempDir()

	// YAML template 1
	t1 := filepath.Join(dir, "template1.yaml")
	writeFile(t, t1, `$schema: https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#
contentVersion: "1.0.0.0"
resources:
  - name: storage1
    type: Microsoft.Storage/storageAccounts
    apiVersion: "2021-04-01"
`)

	// YAML template 2 (same content)
	t2 := filepath.Join(dir, "template2.yaml")
	writeFile(t, t2, `$schema: https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#
contentVersion: "1.0.0.0"
resources:
  - name: storage1
    type: Microsoft.Storage/storageAccounts
    apiVersion: "2021-04-01"
`)

	d := New()
	result, err := d.Diff(nil, t1, t2, coredomain.DiffOpts{})
	if err != nil {
		t.Fatalf("Diff failed: %v", err)
	}

	if result.Summary.Total != 0 {
		t.Errorf("expected 0 differences, got %d", result.Summary.Total)
	}
}

func TestDiff_PropertyChanges(t *testing.T) {
	dir := t.TempDir()

	// Template 1 with properties
	t1 := filepath.Join(dir, "template1.json")
	writeJSON(t, t1, `{
		"$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"resources": [
			{
				"name": "storage1",
				"type": "Microsoft.Storage/storageAccounts",
				"apiVersion": "2021-04-01",
				"properties": {
					"accessTier": "Hot",
					"supportsHttpsTrafficOnly": true
				}
			}
		]
	}`)

	// Template 2: property changed
	t2 := filepath.Join(dir, "template2.json")
	writeJSON(t, t2, `{
		"$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"resources": [
			{
				"name": "storage1",
				"type": "Microsoft.Storage/storageAccounts",
				"apiVersion": "2021-04-01",
				"properties": {
					"accessTier": "Cool",
					"supportsHttpsTrafficOnly": true
				}
			}
		]
	}`)

	d := New()
	result, err := d.Diff(nil, t1, t2, coredomain.DiffOpts{})
	if err != nil {
		t.Fatalf("Diff failed: %v", err)
	}

	if result.Summary.Modified != 1 {
		t.Errorf("expected 1 modified, got %d", result.Summary.Modified)
	}
}

func writeJSON(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write %s: %v", path, err)
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write %s: %v", path, err)
	}
}
