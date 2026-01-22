<picture>
  <source media="(prefers-color-scheme: dark)" srcset="../../docs/wetwire-dark.svg">
  <img src="../../docs/wetwire-light.svg" width="100" height="67">
</picture>

This example demonstrates deploying a basic Azure Storage Account using wetwire-azure-go.

## Resources Created

- **StorageAccount**: A StorageV2 account with Standard_LRS replication

## Configuration

| Property | Value | Description |
|----------|-------|-------------|
| Name | mystorageaccount | Must be globally unique |
| Location | eastus | Azure region |
| SKU | Standard_LRS | Locally redundant storage |
| Kind | StorageV2 | General-purpose v2 account |

## Build

Generate the ARM template:

```bash
wetwire-azure build -o template.json
```

## Deploy

```bash
# Login to Azure
az login

# Create resource group
az group create --name storage-example-rg --location eastus

# Deploy the template
az deployment group create \
  --resource-group storage-example-rg \
  --template-file template.json
```

## Customization

### Change Storage Tier

Modify the SKU name for different redundancy:

```go
SKU: storage.SKU{
    Name: "Standard_GRS",  // Geo-redundant storage
}
```

Available SKUs:
- `Standard_LRS` - Locally redundant
- `Standard_GRS` - Geo-redundant
- `Standard_RAGRS` - Read-access geo-redundant
- `Standard_ZRS` - Zone-redundant
- `Premium_LRS` - Premium locally redundant

### Add Tags

```go
var MyStorageAccount = storage.StorageAccount{
    Name:     "mystorageaccount",
    Location: "eastus",
    Tags: map[string]string{
        "environment": "development",
        "project":     "example",
    },
    SKU: storage.SKU{
        Name: "Standard_LRS",
    },
    Kind: "StorageV2",
}
```

## Cleanup

```bash
az group delete --name storage-example-rg --yes
```
