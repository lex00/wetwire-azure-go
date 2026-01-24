---
title: "Adoption"
---

This guide helps you migrate from other infrastructure-as-code tools to wetwire-azure-go.

## Why wetwire-azure-go?

- **Type Safety**: Full Go type checking for Azure resources
- **IDE Support**: Autocompletion and inline documentation
- **No Runtime**: Pure code generation, no runtime dependencies
- **Familiar Language**: Use Go instead of DSLs or YAML

## Migration from Terraform

### Terraform to wetwire-azure-go

**Terraform:**
```hcl
resource "azurerm_storage_account" "example" {
  name                     = "mystorageaccount"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = "eastus"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}
```

**wetwire-azure-go:**
```go
var MyStorage = storage.StorageAccount{
    Name:     "mystorageaccount",
    Location: "eastus",
    SKU: storage.SKU{
        Name: "Standard_LRS",
    },
    Kind: "StorageV2",
}
```

### Key Differences

| Terraform | wetwire-azure-go |
|-----------|------------------|
| HCL configuration | Go code |
| State management | Stateless (ARM handles it) |
| `terraform apply` | `az deployment group create` |
| Provider plugins | Direct ARM templates |

### Migration Steps

1. Export Terraform state to ARM template
2. Import using `wetwire-azure import`
3. Review and customize generated code
4. Deploy using Azure CLI or portal

## Migration from Azure CDK

### CDK to wetwire-azure-go

**Azure CDK (TypeScript):**
```typescript
const storageAccount = new storage.StorageAccount(this, 'MyStorage', {
  resourceGroupName: resourceGroup.name,
  location: 'eastus',
  kind: storage.Kind.StorageV2,
  sku: {
    name: storage.SkuName.Standard_LRS,
  },
});
```

**wetwire-azure-go:**
```go
var MyStorage = storage.StorageAccount{
    Name:     "mystorageaccount",
    Location: "eastus",
    SKU: storage.SKU{
        Name: "Standard_LRS",
    },
    Kind: "StorageV2",
}
```

### Key Differences

| Azure CDK | wetwire-azure-go |
|-----------|------------------|
| TypeScript/Python/etc | Go |
| Constructs | Direct resource types |
| Runtime synthesis | Static code generation |
| Custom constructs | Go functions/packages |

## Migration from Pulumi

### Pulumi to wetwire-azure-go

**Pulumi (Go):**
```go
storageAccount, err := storage.NewStorageAccount(ctx, "myStorage", &storage.StorageAccountArgs{
    ResourceGroupName: resourceGroup.Name,
    Location:          pulumi.String("eastus"),
    Sku: &storage.SkuArgs{
        Name: pulumi.String("Standard_LRS"),
    },
    Kind: pulumi.String("StorageV2"),
})
```

**wetwire-azure-go:**
```go
var MyStorage = storage.StorageAccount{
    Name:     "mystorageaccount",
    Location: "eastus",
    SKU: storage.SKU{
        Name: "Standard_LRS",
    },
    Kind: "StorageV2",
}
```

### Key Differences

| Pulumi | wetwire-azure-go |
|--------|------------------|
| Provider SDK | Generated ARM types |
| State management | Stateless |
| Input types | Plain Go structs |
| Async operations | Static generation |

## Migration Strategy

### Phase 1: Assessment
1. Inventory existing resources
2. Identify resource dependencies
3. Plan migration order

### Phase 2: Parallel Deployment
1. Import existing templates
2. Deploy wetwire alongside existing
3. Validate equivalence

### Phase 3: Cutover
1. Remove old IaC management
2. Establish wetwire as source of truth
3. Update CI/CD pipelines

## Best Practices

1. **Start with Non-Production**: Test migration on dev/staging first
2. **Keep History**: Maintain old IaC configs until migration is complete
3. **Test Thoroughly**: Verify deployments match expected state
4. **Document Changes**: Track any manual adjustments needed

## See Also

- [Import Workflow](IMPORT_WORKFLOW.md)
- [Quick Start](QUICK_START.md)
- [CLI Reference](CLI.md)
