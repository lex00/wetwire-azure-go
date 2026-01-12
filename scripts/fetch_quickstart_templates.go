// fetch_quickstart_templates.go fetches ARM templates from Azure Quickstart Templates repo
// for round-trip testing.
//
// Usage:
//
//	go run scripts/fetch_quickstart_templates.go
//
// This script downloads curated ARM templates from the Azure Quickstart Templates
// GitHub repository and saves them to testdata/azure-quickstarts/.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Template categories to fetch
var templateURLs = map[string]string{
	// Storage accounts
	"storage-account-create.json": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/quickstarts/microsoft.storage/storage-account-create/azuredeploy.json",

	// Virtual networks
	"vnet-two-subnets.json": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/quickstarts/microsoft.network/vnet-two-subnets/azuredeploy.json",

	// Network security groups
	"security-group-create.json": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/quickstarts/microsoft.network/security-group-create/azuredeploy.json",

	// Virtual machines
	"vm-simple-linux.json": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/quickstarts/microsoft.compute/vm-simple-linux/azuredeploy.json",

	// Web apps
	"webapp-basic-linux.json": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/quickstarts/microsoft.web/webapp-basic-linux/azuredeploy.json",

	// Key Vault
	"key-vault-create.json": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/quickstarts/microsoft.keyvault/key-vault-create/azuredeploy.json",

	// SQL Database
	"sql-database.json": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/quickstarts/microsoft.sql/sql-database/azuredeploy.json",

	// Cosmos DB
	"cosmosdb-free.json": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/quickstarts/microsoft.documentdb/cosmosdb-free/azuredeploy.json",

	// Public IP
	"public-ip-create.json": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/quickstarts/microsoft.network/nic-publicip-dns-vnet/azuredeploy.json",

	// Load Balancer
	"load-balancer-create.json": "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/quickstarts/microsoft.network/load-balancer-standard-create/azuredeploy.json",
}

func main() {
	outputDir := filepath.Join("testdata", "azure-quickstarts")

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Fetching %d Azure Quickstart templates...\n", len(templateURLs))

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	successCount := 0
	failCount := 0

	for filename, url := range templateURLs {
		outputPath := filepath.Join(outputDir, filename)

		fmt.Printf("  Fetching %s... ", filename)

		if err := fetchTemplate(client, url, outputPath); err != nil {
			fmt.Printf("FAILED: %v\n", err)
			failCount++
			continue
		}

		// Validate it's a valid ARM template
		if err := validateARMTemplate(outputPath); err != nil {
			fmt.Printf("INVALID: %v\n", err)
			os.Remove(outputPath)
			failCount++
			continue
		}

		fmt.Printf("OK\n")
		successCount++
	}

	fmt.Printf("\nComplete: %d succeeded, %d failed\n", successCount, failCount)

	if failCount > 0 {
		os.Exit(1)
	}
}

func fetchTemplate(client *http.Client, url, outputPath string) error {
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func validateARMTemplate(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var template map[string]interface{}
	if err := json.Unmarshal(data, &template); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	// Check for required ARM template fields
	schema, ok := template["$schema"].(string)
	if !ok || !strings.Contains(schema, "schema.management.azure.com") {
		return fmt.Errorf("missing or invalid $schema")
	}

	if _, ok := template["contentVersion"]; !ok {
		return fmt.Errorf("missing contentVersion")
	}

	if _, ok := template["resources"]; !ok {
		return fmt.Errorf("missing resources")
	}

	return nil
}
