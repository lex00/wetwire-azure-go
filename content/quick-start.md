---
title: "Quick Start"
---

Get up and running with wetwire-azure-go in 5 minutes.

## Prerequisites

- Go 1.21 or later
- Azure CLI (for deployment)

## Installation

See [README.md](../README.md#installation) for installation instructions.

## Create Your First Project

### 1. Initialize a New Project

```bash
mkdir my-azure-infra && cd my-azure-infra
wetwire-azure init
```

This creates:
- `go.mod` - Go module file
- `main.go` - Sample storage account resource
- `.gitignore` - Common ignore patterns
- `README.md` - Project documentation

### 2. Define Resources

Edit `main.go` to define your Azure resources:

```go
package main

import (
    "github.com/lex00/wetwire-azure-go/resources/storage"
)

// MyStorage defines a storage account
var MyStorage = storage.StorageAccount{
    Name:     "mystorageaccount",
    Location: "eastus",
    SKU: storage.SKU{
        Name: "Standard_LRS",
    },
    Kind: "StorageV2",
}
```

### 3. Build the ARM Template

```bash
wetwire-azure build -o template.json
```

### 4. Deploy to Azure

```bash
# Login to Azure
az login

# Create resource group
az group create --name my-rg --location eastus

# Deploy the template
az deployment group create \
  --resource-group my-rg \
  --template-file template.json
```

## Next Steps

### List Resources

See all resources discovered in your project:

```bash
wetwire-azure list
```

### Lint Your Code

Check for common issues:

```bash
wetwire-azure lint
```

### Generate Dependency Graph

Visualize resource dependencies:

```bash
wetwire-azure graph --format mermaid
```

### Import Existing Templates

Convert an existing ARM template to Go:

```bash
wetwire-azure import template.json -o infra.go
```

## Key Concepts

### Resource Declaration

Resources are Go struct literals assigned to package-level variables:

```go
var MyVM = compute.VirtualMachine{
    Name:     "my-vm",
    Location: "eastus",
    // ...
}
```

### Direct References

Reference other resources by variable name:

```go
var MyNIC = network.NetworkInterface{...}

var MyVM = compute.VirtualMachine{
    NetworkProfile: compute.NetworkProfile{
        NetworkInterfaces: []compute.NetworkInterfaceReference{
            {Id: MyNIC.Id},  // Direct reference
        },
    },
}
```

### ARM Intrinsics

Use ARM template functions via the intrinsics package:

```go
import . "github.com/lex00/wetwire-azure-go/intrinsics"

var MyStorage = storage.StorageAccount{
    Location: ResourceGroup().Location,  // Uses resourceGroup().location
}
```

## AI-Assisted Design

Let AI help create your Azure infrastructure:

```bash
# No API key required - uses Claude CLI
wetwire-azure design "Create a storage account with geo-redundant storage and a virtual network"
```

The design command creates Go code following wetwire patterns, runs linting, and builds the final ARM template.

## Learn More

- [INTERNALS.md](INTERNALS.md) - Architecture deep-dive
- [EXAMPLES.md](EXAMPLES.md) - Example projects
- [TROUBLESHOOTING.md](TROUBLESHOOTING.md) - Common issues
- [CONTRIBUTING.md](CONTRIBUTING.md) - Development guide
