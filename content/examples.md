---
title: "Examples"
---

The `examples/` directory contains example wetwire-azure-go projects demonstrating common Azure infrastructure patterns.

## Available Examples

### [storage-account](https://github.com/lex00/wetwire-azure-go/tree/main/examples/storage-account)

A minimal example showing how to define an Azure Storage Account.

**Concepts demonstrated:**
- Basic resource definition
- SKU configuration
- Resource properties

```go
var MyStorage = storage.StorageAccount{
    Name:     "mystorageaccount",
    Location: "eastus",
    SKU:      storage.SKU{Name: "Standard_LRS"},
    Kind:     "StorageV2",
}
```

### [virtual-machine](https://github.com/lex00/wetwire-azure-go/tree/main/examples/virtual-machine)

A complete example showing a Linux virtual machine with networking.

**Concepts demonstrated:**
- Multiple resource types
- Resource dependencies
- Network configuration
- OS profile setup

**Resources created:**
- Virtual Network
- Subnet
- Network Interface
- Virtual Machine

### [aks-golden](https://github.com/lex00/wetwire-azure-go/tree/main/examples/aks-golden)

Production-ready AKS cluster configuration.

### [aks-k8s](https://github.com/lex00/wetwire-azure-go/tree/main/examples/aks-k8s)

AKS cluster with Kubernetes resource integration.

### [enterprise-app](https://github.com/lex00/wetwire-azure-go/tree/main/examples/enterprise-app)

Enterprise application pattern with multiple services.

### [security-best-practices](https://github.com/lex00/wetwire-azure-go/tree/main/examples/security-best-practices)

Security-focused configuration patterns.

### [parameters-and-outputs](https://github.com/lex00/wetwire-azure-go/tree/main/examples/parameters-and-outputs)

ARM template parameters and outputs usage.

## Running Examples

```bash
# Build ARM template from example
cd examples/storage-account
wetwire-azure build .

# Build and save to file
wetwire-azure build . -o template.json

# Validate generated template
wetwire-azure validate template.json

# View dependency graph
wetwire-azure graph . --format mermaid
```

## Deploying to Azure

After generating the template:

```bash
# Create resource group
az group create --name my-rg --location eastus

# Deploy template
az deployment group create \
  --resource-group my-rg \
  --template-file template.json
```

## Example Structure

Each example follows this structure:

```
example-name/
├── main.go         # Resource definitions
├── go.mod          # Go module file
└── README.md       # Example documentation
```

## Creating Your Own

1. Initialize a new project:
   ```bash
   wetwire-azure init myproject
   ```

2. Define resources in `main.go`

3. Build and validate:
   ```bash
   wetwire-azure build .
   wetwire-azure lint .
   ```

## See Also

- [CLI](/cli/) - Command-line interface reference
- [FAQ](/faq/) - Frequently asked questions
