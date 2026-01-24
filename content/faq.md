---
title: "FAQ"
---

Frequently asked questions about wetwire-azure-go.

---

## Getting Started

<details>
<summary>How do I install wetwire-azure-go?</summary>

```bash
go install github.com/lex00/wetwire-azure-go@latest
```

</details>

<details>
<summary>How do I create a new project?</summary>

```bash
wetwire-azure init my-infrastructure
cd my-infrastructure
```

</details>

<details>
<summary>How do I build an ARM template?</summary>

```bash
wetwire-azure build ./infra > template.json
```

</details>

<details>
<summary>How do I build a Bicep template?</summary>

```bash
wetwire-azure build ./infra --format bicep > main.bicep
```

</details>

---

## Syntax

<details>
<summary>How do I reference another resource?</summary>

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

</details>

<details>
<summary>How do I use ARM template functions?</summary>

Use the intrinsics package:

```go
import . "github.com/lex00/wetwire-azure-go/intrinsics"

var MyStorage = storage.StorageAccount{
    Name: Concat([]any{ResourceGroup().Name, "-storage"}),
    Location: ResourceGroup().Location,
}
```

</details>

<details>
<summary>How do I reference resource group properties?</summary>

Use the `ResourceGroup()` function:

```go
// Resource group name
Name: ResourceGroup().Name

// Resource group location (recommended for all resources)
Location: ResourceGroup().Location

// Resource group ID
Id: ResourceGroup().Id
```

</details>

<details>
<summary>How do I create resource dependencies?</summary>

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

</details>

---

## Azure-Specific Questions

<details>
<summary>How do I handle storage account name constraints?</summary>

Storage account names must be 3-24 lowercase alphanumeric characters and globally unique:

```go
// Bad: Contains hyphens, too long, not globally unique
Name: "my-storage-account-name-that-is-too-long"

// Good: Lowercase, alphanumeric, uses uniqueString for global uniqueness
Name: Concat([]any{"storage", UniqueString(ResourceGroup().Id)})
```

</details>

<details>
<summary>How do I specify API versions for resources?</summary>

API versions are handled automatically based on the resource type schema. The build command selects appropriate API versions.

</details>

<details>
<summary>How do I use managed identities?</summary>

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

</details>

<details>
<summary>How do I create a VM with a managed disk?</summary>

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

</details>

<details>
<summary>How do I create a VNet with subnets?</summary>

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

</details>

<details>
<summary>How do I set tags on resources?</summary>

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

</details>

<details>
<summary>How do I use parameters for deployment-time values?</summary>

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

</details>

---

## CI/CD and DevOps

<details>
<summary>How do I integrate wetwire-azure with my CI/CD pipeline?</summary>

Add lint and build steps to your pipeline:

```yaml
# GitHub Actions example
steps:
  - name: Install wetwire-azure
    run: go install github.com/lex00/wetwire-azure-go@latest

  - name: Lint infrastructure code
    run: wetwire-azure lint ./infra

  - name: Build ARM template
    run: wetwire-azure build ./infra > template.json

  - name: Deploy to Azure
    run: |
      az deployment group create \
        --resource-group ${{ vars.RESOURCE_GROUP }} \
        --template-file template.json
```

For Azure DevOps:

```yaml
steps:
  - task: GoTool@0
    inputs:
      version: '1.23'

  - script: go install github.com/lex00/wetwire-azure-go@latest
    displayName: 'Install wetwire-azure'

  - script: wetwire-azure lint ./infra
    displayName: 'Lint'

  - script: wetwire-azure build ./infra > $(Build.ArtifactStagingDirectory)/template.json
    displayName: 'Build ARM template'
```

</details>

<details>
<summary>Can I import existing ARM templates?</summary>

Yes, use the import command:

```bash
# Import ARM template
wetwire-azure import template.json -o ./my-infrastructure

# Import Bicep file
wetwire-azure import main.bicep -o ./my-infrastructure

# Apply lint fixes after import
wetwire-azure lint --fix ./my-infrastructure
```

Import is best-effort. Complex templates may need manual cleanup after import.

</details>

<details>
<summary>How does the linter help catch errors?</summary>

The linter enforces wetwire patterns and Azure best practices:

```bash
wetwire-azure lint ./infra
```

It checks for:
- **Type safety**: Validates resource types and properties
- **Reference validity**: Ensures referenced resources exist
- **Security**: Detects hardcoded secrets and insecure configurations
- **Azure constraints**: Validates naming conventions, API versions
- **Pattern compliance**: Enforces declarative patterns

Use `--fix` to auto-fix issues where possible:

```bash
wetwire-azure lint --fix ./infra
```

</details>

<details>
<summary>What's the recommended project structure?</summary>

Organize by Azure resource category:

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

For larger projects, use subdirectories:

```
my-azure-infra/
├── go.mod
├── environments/
│   ├── dev/
│   ├── staging/
│   └── prod/
└── modules/
    ├── networking/
    ├── compute/
    └── data/
```

</details>

<details>
<summary>How do I handle linked templates?</summary>

For complex deployments with multiple ARM templates, structure your Go code in separate packages and build each independently:

```bash
# Build main template
wetwire-azure build ./infra/main > main.json

# Build linked templates
wetwire-azure build ./infra/networking > networking.json
wetwire-azure build ./infra/compute > compute.json
```

Reference linked templates using `templateLink`:

```go
var NetworkingDeployment = resources.Deployment{
    Name: "networking",
    Properties: resources.Deployment_Properties{
        Mode: "Incremental",
        TemplateLink: resources.Deployment_TemplateLink{
            Uri: "[uri(deployment().properties.templateLink.uri, 'networking.json')]",
        },
    },
}
```

</details>

---

## Lint Rules

<details>
<summary>What do the WAZ rule codes mean?</summary>

WAZ stands for "Wetwire AZure". Common rules include:

| Rule | Description |
|------|-------------|
| WAZ001 | Use location constants for common regions |
| WAZ002 | Use intrinsic types for ARM template functions |
| WAZ003 | Extract inline property types to named variables |
| WAZ004 | Use typed structs instead of `map[string]any` |
| WAZ005 | Detect duplicate resource names |

See [Lint Rules]({{< relref "/lint-rules" >}}) for complete documentation.

</details>

<details>
<summary>How do I auto-fix lint issues?</summary>

```bash
wetwire-azure lint --fix ./infra
```

</details>

<details>
<summary>How do I disable a specific lint rule?</summary>

Currently lint rules cannot be disabled individually. If you have a valid use case, file an issue.

</details>

---

## Import

<details>
<summary>How do I convert an existing ARM template?</summary>

```bash
wetwire-azure import template.json -o ./my-infrastructure
```

</details>

<details>
<summary>How do I convert a Bicep file?</summary>

```bash
wetwire-azure import main.bicep -o ./my-infrastructure
```

</details>

<details>
<summary>Import produced code that doesn't compile?</summary>

Import is best-effort. Complex templates may need manual cleanup:

1. Run `wetwire-azure lint --fix ./infra` to apply automatic fixes
2. Review and manually fix remaining issues
3. Check for unsupported ARM template functions

</details>

<details>
<summary>What ARM template features are supported by import?</summary>

- Resources (all Azure types)
- Parameters
- Variables
- Outputs
- Most ARM template functions (concat, resourceId, reference, etc.)
- Dependencies (via dependsOn)

</details>

---

## Design Mode

<details>
<summary>How do I use AI-assisted design?</summary>

```bash
export ANTHROPIC_API_KEY=your-key
wetwire-azure design "Create a Linux VM with storage account"
```

</details>

<details>
<summary>What model does design mode use?</summary>

Claude (via Anthropic API). The specific model is configured internally.

</details>

<details>
<summary>Can I use design mode without an API key?</summary>

No. Design mode requires the Anthropic API.

</details>

---

## Troubleshooting

<details>
<summary>"cannot find package" errors</summary>

Ensure your `go.mod` has the correct module path and dependencies:

```bash
go mod tidy
```

</details>

<details>
<summary>"undefined: storage" or similar import errors</summary>

Add the missing import statement:

```go
import "github.com/lex00/wetwire-azure-go/resources/storage"
```

</details>

<details>
<summary>Build produces empty template</summary>

Check that:
1. Resources are declared as package-level `var` statements
2. Resources have the correct type (e.g., `storage.StorageAccount`)
3. The package path is correct in the build command

</details>

<details>
<summary>Circular dependency detected</summary>

Resources cannot have circular references. Review the dependency graph:

```bash
wetwire-azure graph ./infra | dot -Tpng -o deps.png
```

Break the cycle by restructuring resources or using parameters.

</details>

<details>
<summary>"unknown resource type" error</summary>

The resource type may be misspelled or not supported. Check:
1. Correct package import (e.g., `compute` not `vm`)
2. Correct type name (e.g., `VirtualMachine` not `VM`)

</details>

<details>
<summary>Storage account name validation errors</summary>

Azure storage account names must be:
- 3-24 characters long
- Lowercase letters and numbers only
- Globally unique across all of Azure

Use `UniqueString()` for global uniqueness:

```go
Name: Concat([]any{"storage", UniqueString(ResourceGroup().Id)})
```

</details>

<details>
<summary>ARM template deployment fails with "resource not found"</summary>

Check resource dependencies. Resources must be deployed in order. Use the graph command to visualize:

```bash
wetwire-azure graph ./infra
```

</details>

<details>
<summary>Bicep conversion issues</summary>

Some complex Go patterns may not convert cleanly to Bicep:
1. Generate ARM JSON first and validate
2. Use Azure CLI to convert ARM to Bicep if needed: `az bicep decompile --file template.json`

</details>

<details>
<summary>Import generates code that doesn't compile</summary>

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

</details>

<details>
<summary>"ANTHROPIC_API_KEY not set" error</summary>

Design and test commands require an API key:

```bash
export ANTHROPIC_API_KEY="sk-ant-..."
wetwire-azure design "Create a storage account"
```

</details>

<details>
<summary>Lint reports issues but --fix doesn't help</summary>

Some lint rules are advisory and don't have auto-fixes:
- WAZ005 (duplicate names) - rename resources manually
- Resource-specific validation errors - fix property values manually

</details>

<details>
<summary>Build succeeds but Azure deployment fails</summary>

The template is syntactically valid but may have semantic issues:
1. Check Azure RBAC permissions for deployment
2. Verify resource names are unique in the resource group
3. Review Azure deployment error messages in portal or CLI

</details>

<details>
<summary>How do I handle resources in different resource groups?</summary>

ARM templates deploy to a single resource group by default. For multi-resource-group deployments:
1. Use nested deployments
2. Deploy separate templates to each resource group
3. Use Azure Blueprints for complex multi-resource-group infrastructure

</details>

<details>
<summary>How do I reference existing resources not in my template?</summary>

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

</details>

---

## Best Practices

<details>
<summary>Should I use ARM JSON or Bicep output?</summary>

**Use ARM JSON when:**
- You need maximum compatibility
- You're using Azure DevOps or CI/CD expecting JSON
- You want to validate with ARM template tools

**Use Bicep when:**
- You prefer more readable output
- You're using Bicep tooling and workflows
- You want better diff visualization in Git

Both formats represent the same infrastructure and can be converted between each other.

</details>

<details>
<summary>How should I organize my Azure infrastructure code?</summary>

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

</details>

<details>
<summary>Should I use parameters or hardcoded values?</summary>

**Use parameters for:**
- Values that change between environments (location, size, SKU)
- Sensitive values (admin passwords, connection strings)
- Values set at deployment time

**Use hardcoded values for:**
- Resource names (generated with uniqueString)
- Configuration that's environment-independent
- Tags and metadata

</details>

<details>
<summary>How do I handle secrets and credentials?</summary>

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

</details>

---

## See Also

- [CLI Reference]({{< relref "/cli" >}})
- [Lint Rules]({{< relref "/lint-rules" >}})
- [Azure ARM Template Reference](https://docs.microsoft.com/azure/templates/)
