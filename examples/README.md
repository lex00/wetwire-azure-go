# Examples

This directory contains example wetwire-azure-go projects demonstrating common Azure infrastructure patterns.

## Available Examples

### [storage-account](./storage-account/)

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

### [virtual-machine](./virtual-machine/)

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

## Attribution

These examples are designed for educational purposes to demonstrate wetwire-azure-go patterns. They are based on common Azure deployment scenarios and official Azure documentation.

## See Also

- [Quick Start Guide](../docs/QUICK_START.md)
- [CLI Reference](../docs/CLI.md)
- [Lint Rules](../docs/LINT_RULES.md)
