// Package main demonstrates a basic Azure Storage Account deployment.
package main

import (
	"github.com/lex00/wetwire-azure-go/resources/storage"
)

// MyStorageAccount defines a basic Azure storage account with standard configuration.
// Storage account names must be globally unique, 3-24 characters, lowercase letters and numbers only.
var MyStorageAccount = storage.StorageAccount{
	Name:     "mystorageaccount",
	Location: "eastus",
	SKU: storage.SKU{
		Name: "Standard_LRS",
	},
	Kind: "StorageV2",
}

// main is required for a valid Go program but not used by wetwire-azure.
// Resources are discovered from package-level variable declarations.
func main() {
	// wetwire-azure build discovers resources via AST parsing
	// No runtime execution is needed
}
