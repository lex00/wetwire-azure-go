---
title: "Faq"
---
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="./wetwire-dark.svg">
  <img src="./wetwire-light.svg" width="100" height="67">
</picture>

This FAQ covers questions specific to the Go implementation of wetwire for Azure ARM/Bicep templates. For general wetwire questions, see the [central FAQ](https://github.com/lex00/wetwire/blob/main/docs/FAQ.md).

---

## Getting Started

### How do I install wetwire-azure-go?

See [README.md](../README.md#installation) for installation instructions.

### How do I create a new project?

```bash
wetwire-azure init my-infrastructure
cd my-infrastructure
```

### How do I build an ARM template?

```bash
wetwire-azure build ./infra > template.json
```

### How do I build a Bicep template?

```bash
wetwire-azure build ./infra --format bicep > main.bicep
```

---

## Syntax

### How do I reference another resource?

Use direct variable references:

```go
var MyVM = compute.VirtualMachine{
    NetworkProfile: compute.VirtualMachine_NetworkProfile{
        NetworkInterfaces: []compute.VirtualMachine_NetworkInterfaceReference{
            {Id: MyNIC.Id},  // Direct reference
        },
    },
}
```

### How do I use ARM template functions?

Use the intrinsics package:

```go
import . "github.com/lex00/wetwire-azure-go/intrinsics"

var MyStorage = storage.StorageAccount{
    Name: Concat([]any{ResourceGroup().Name, "-storage"}),
    Location: ResourceGroup().Location,
}
```

### How do I reference resource group properties?

Use the `ResourceGroup()` function:

```go
// Resource group name
Name: ResourceGroup().Name

// Resource group location (recommended for all resources)
Location: ResourceGroup().Location

// Resource group ID
Id: ResourceGroup().Id
```

### How do I create resource dependencies?

Dependencies are implicit via resource references:

```go
// VM depends on NIC (implicit via Id reference)
var MyVM = compute.VirtualMachine{
    NetworkProfile: compute.VirtualMachine_NetworkProfile{
        NetworkInterfaces: []compute.VirtualMachine_NetworkInterfaceReference{
            {Id: MyNIC.Id},  // Creates dependency
        },
    },
}
```

---

## Azure-Specific Questions

### How do I handle storage account name constraints?

Storage account names must be 3-24 lowercase alphanumeric characters and globally unique:

```go
// Bad: Contains hyphens, too long, not globally unique
Name: "my-storage-account-name-that-is-too-long"

// Good: Lowercase, alphanumeric, uses uniqueString for global uniqueness
Name: Concat([]any{"storage", UniqueString(ResourceGroup().Id)})
```

### How do I specify API versions for resources?

API versions are handled automatically based on the resource type schema. The build command selects appropriate API versions.

### How do I use managed identities?

```go
var MyVMIdentity = compute.VirtualMachine_Identity{
    Type: "SystemAssigned",
}

var MyVM = compute.VirtualMachine{
    Identity: MyVMIdentity,
}

// Reference in role assignments
var MyRoleAssignment = authorization.RoleAssignment{
    PrincipalId: MyVM.Identity.PrincipalId,
    RoleDefinitionId: "/subscriptions/.../roleDefinitions/...",
}
```

### How do I create a VM with a managed disk?

```go
var MyDisk = compute.Disk{
    Name:     "my-vm-disk",
    Location: "eastus",
    Sku: compute.Disk_Sku{
        Name: "Premium_LRS",
    },
    DiskSizeGB: 128,
}

var MyVM = compute.VirtualMachine{
    StorageProfile: compute.VirtualMachine_StorageProfile{
        OsDisk: compute.VirtualMachine_OSDisk{
            ManagedDisk: compute.VirtualMachine_ManagedDiskParameters{
                Id: MyDisk.Id,
            },
        },
    },
}
```

### How do I create a VNet with subnets?

```go
var MySubnet = network.VirtualNetwork_Subnet{
    Name: "default",
    AddressPrefix: "10.0.1.0/24",
}

var MyVNet = network.VirtualNetwork{
    Name:     "my-vnet",
    Location: "eastus",
    AddressSpace: network.VirtualNetwork_AddressSpace{
        AddressPrefixes: []string{"10.0.0.0/16"},
    },
    Subnets: []network.VirtualNetwork_Subnet{MySubnet},
}
```

### How do I set tags on resources?

```go
var MyStorage = storage.StorageAccount{
    Name:     "mystorageaccount",
    Location: "eastus",
    Tags: map[string]string{
        "Environment": "Production",
        "Department":  "Engineering",
    },
}
```

### How do I use parameters for deployment-time values?

```go
var LocationParam = Parameter{
    Type:         "string",
    DefaultValue: "eastus",
    AllowedValues: []string{"eastus", "westus", "centralus"},
}

var MyStorage = storage.StorageAccount{
    Location: LocationParam,
}
```

---

## Lint Rules

### What do the WAZ rule codes mean?

WAZ stands for "Wetwire AZure". Common rules include:

| Rule | Description |
|------|-------------|
| WAZ001 | Use location constants for common regions |
| WAZ002 | Use intrinsic types for ARM template functions |
| WAZ003 | Extract inline property types to named variables |
| WAZ004 | Use typed structs instead of `map[string]any` |
| WAZ005 | Detect duplicate resource names |

See [LINT_RULES.md](LINT_RULES.md) for complete documentation.

### How do I auto-fix lint issues?

```bash
wetwire-azure lint --fix ./infra
```

### How do I disable a specific lint rule?

Currently lint rules cannot be disabled individually. If you have a valid use case, file an issue.

---

## Import

### How do I convert an existing ARM template?

```bash
wetwire-azure import template.json -o ./my-infrastructure
```

### How do I convert a Bicep file?

```bash
wetwire-azure import main.bicep -o ./my-infrastructure
```

### Import produced code that doesn't compile?

Import is best-effort. Complex templates may need manual cleanup:

1. Run `wetwire-azure lint --fix ./infra` to apply automatic fixes
2. Review and manually fix remaining issues
3. Check for unsupported ARM template functions

### What ARM template features are supported by import?

- Resources (all Azure types)
- Parameters
- Variables
- Outputs
- Most ARM template functions (concat, resourceId, reference, etc.)
- Dependencies (via dependsOn)

---

## Design Mode

### How do I use AI-assisted design?

```bash
export ANTHROPIC_API_KEY=your-key
wetwire-azure design "Create a Linux VM with storage account"
```

### What model does design mode use?

Claude (via Anthropic API). The specific model is configured in wetwire-core-go.

### Can I use design mode without an API key?

No. Design mode requires the Anthropic API.

---

## Troubleshooting

### "cannot find package" errors

Ensure your `go.mod` has the correct module path and dependencies:

```bash
go mod tidy
```

### "undefined: storage" or similar import errors

Add the missing import statement:

```go
import "github.com/lex00/wetwire-azure-go/resources/storage"
```

### Build produces empty template

Check that:
1. Resources are declared as package-level `var` statements
2. Resources have the correct type (e.g., `storage.StorageAccount`)
3. The package path is correct in the build command

### Circular dependency detected

Resources cannot have circular references. Review the dependency graph:

```bash
wetwire-azure graph ./infra | dot -Tpng -o deps.png
```

Break the cycle by restructuring resources or using parameters.

### "unknown resource type" error

The resource type may be misspelled or not supported. Check:
1. Correct package import (e.g., `compute` not `vm`)
2. Correct type name (e.g., `VirtualMachine` not `VM`)

### Storage account name validation errors

Azure storage account names must be:
- 3-24 characters long
- Lowercase letters and numbers only
- Globally unique across all of Azure

Use `UniqueString()` for global uniqueness:

```go
Name: Concat([]any{"storage", UniqueString(ResourceGroup().Id)})
```

### ARM template deployment fails with "resource not found"

Check resource dependencies. Resources must be deployed in order. Use the graph command to visualize:

```bash
wetwire-azure graph ./infra
```

### Bicep conversion issues

Some complex Go patterns may not convert cleanly to Bicep:
1. Generate ARM JSON first and validate
2. Use Azure CLI to convert ARM to Bicep if needed: `az bicep decompile --file template.json`

### Import generates code that doesn't compile

This is expected for complex templates. Fix with:

```bash
# Apply automatic fixes
wetwire-azure lint --fix ./infra

# Then manually fix remaining issues
```

Common import issues:
- Forward references (resource used before declaration)
- Complex nested ARM template functions
- Custom resource types

### "ANTHROPIC_API_KEY not set" error

Design and test commands require an API key:

```bash
export ANTHROPIC_API_KEY="sk-ant-..."
wetwire-azure design "Create a storage account"
```

### Lint reports issues but --fix doesn't help

Some lint rules are advisory and don't have auto-fixes:
- WAZ005 (duplicate names) - rename resources manually
- Resource-specific validation errors - fix property values manually

### Build succeeds but Azure deployment fails

The template is syntactically valid but may have semantic issues:
1. Check Azure RBAC permissions for deployment
2. Verify resource names are unique in the resource group
3. Review Azure deployment error messages in portal or CLI

### How do I handle resources in different resource groups?

ARM templates deploy to a single resource group by default. For multi-resource-group deployments:
1. Use nested deployments
2. Deploy separate templates to each resource group
3. Use Azure Blueprints for complex multi-resource-group infrastructure

### How do I reference existing resources not in my template?

Use the `Reference()` function with explicit resource IDs:

```go
// Reference existing storage account
ExistingStorageRef := Reference(
    ResourceId(
        "Microsoft.Storage/storageAccounts",
        "existing-storage-account",
    ),
)

// Use properties
ConnectionString: ExistingStorageRef.primaryEndpoints.blob
```

---

## Best Practices

### Should I use ARM JSON or Bicep output?

**Use ARM JSON when:**
- You need maximum compatibility
- You're using Azure DevOps or CI/CD expecting JSON
- You want to validate with ARM template tools

**Use Bicep when:**
- You prefer more readable output
- You're using Bicep tooling and workflows
- You want better diff visualization in Git

Both formats represent the same infrastructure and can be converted between each other.

### How should I organize my Azure infrastructure code?

Recommended structure:

```
my-azure-infra/
├── go.mod
├── network.go         # VNets, subnets, NSGs, route tables
├── compute.go         # VMs, scale sets, availability sets
├── storage.go         # Storage accounts, disks, file shares
├── database.go        # SQL databases, Cosmos DB
├── security.go        # Key vaults, managed identities
└── monitoring.go      # Log Analytics, Application Insights
```

### Should I use parameters or hardcoded values?

**Use parameters for:**
- Values that change between environments (location, size, SKU)
- Sensitive values (admin passwords, connection strings)
- Values set at deployment time

**Use hardcoded values for:**
- Resource names (generated with uniqueString)
- Configuration that's environment-independent
- Tags and metadata

### How do I handle secrets and credentials?

Never hardcode secrets. Use:

1. **Azure Key Vault** - Store secrets in Key Vault, reference in template
2. **Managed Identities** - Use for service-to-service authentication
3. **Parameters with secureString** - For deployment-time secrets

```go
var AdminPasswordParam = Parameter{
    Type: "secureString",
}

var MyVM = compute.VirtualMachine{
    OsProfile: compute.VirtualMachine_OSProfile{
        AdminPassword: AdminPasswordParam,
    },
}
```

---

## See Also

- [Wetwire Specification](https://github.com/lex00/wetwire/blob/main/docs/WETWIRE_SPEC.md)
- [CLI Reference](CLI.md)
- [Lint Rules](LINT_RULES.md)
- [Azure ARM Template Reference](https://docs.microsoft.com/azure/templates/)
