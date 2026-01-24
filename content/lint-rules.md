---
title: "Lint Rules"
---
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="./wetwire-dark.svg">
  <img src="./wetwire-light.svg" width="100" height="67">
</picture>

wetwire-azure-go includes lint rules to enforce best practices and idiomatic patterns for declarative Azure ARM/Bicep infrastructure-as-code.

## Quick Start

```bash
# Lint your infrastructure code
wetwire-azure lint ./infra

# Lint with auto-fix (where supported)
wetwire-azure lint ./infra --fix

# Output in JSON format
wetwire-azure lint ./infra -f json
```

## Rule Naming Convention

All rules use the **WAZ** prefix (Wetwire AZure) followed by a 3-digit number:

- **WAZ001-099**: Type safety and constants
- **WAZ100-199**: Direct references and intrinsics
- **WAZ200-299**: Code extraction and flattening
- **WAZ300-399**: Security and best practices
- **WAZ400-499**: Azure-specific patterns

## Rule Index

| Rule | Description | Severity | Auto-fix |
|------|-------------|----------|----------|
| WAZ001 | Use location constants | warning | Yes |
| WAZ002 | Use direct references | warning | No |
| WAZ003 | Extract nested configurations | warning | No |
| WAZ004 | Detect duplicate resource names | error | No |
| WAZ005 | Detect circular dependencies | error | No |
| WAZ006 | Detect secrets and credentials | error | No |
| WAZ007 | Detect sensitive file paths | warning | No |
| WAZ008 | Detect insecure defaults | warning | No |
| WAZ301 | Require HTTPS-only for storage | warning | No |
| WAZ302 | Detect permissive NSG rules | warning | No |
| WAZ303 | Require tags on resources | warning | No |
| WAZ304 | Warn on deprecated API versions | warning | No |

## Planned Rules

The following rules are planned but not yet implemented:

### Type Safety (WAZ001-099)

- **WAZ001**: Use location constants for common regions
- **WAZ002**: Use intrinsic types for ARM template functions
- **WAZ010**: Use typed enum constants for VM sizes
- **WAZ011**: Use typed enum constants for storage SKUs
- **WAZ012**: Validate API versions are current

### Direct References (WAZ100-199)

- **WAZ100**: Use direct resource references instead of resourceId
- **WAZ101**: Avoid explicit resource ID concatenation
- **WAZ102**: Use ResourceGroup() function for resource group properties
- **WAZ103**: Use Subscription() function for subscription properties

### Code Extraction (WAZ200-299)

- **WAZ200**: Extract inline property types to named variables
- **WAZ201**: Flatten inline typed structs
- **WAZ202**: Use named var declarations (block style)
- **WAZ203**: Split files with too many resources (>20)

### Security (WAZ300-399)

**Implemented:**
- **WAZ301**: Require HTTPS-only for storage accounts
- **WAZ302**: Detect overly permissive NSG rules (0.0.0.0/0 or *)
- **WAZ303**: Require tags on Azure resources for organization
- **WAZ304**: Warn on deprecated API versions (pre-2021)

**Planned:**
- **WAZ300**: Detect hardcoded secrets and credentials
- **WAZ305**: Require encryption for storage accounts
- **WAZ306**: Require encryption for managed disks
- **WAZ307**: Require secureString for password parameters

### Azure-Specific (WAZ400-499)

- **WAZ400**: Validate storage account name constraints (3-24 chars, lowercase, alphanumeric)
- **WAZ401**: Validate VM name constraints
- **WAZ402**: Use UniqueString for globally unique names
- **WAZ403**: Require tags on resources
- **WAZ404**: Validate location values
- **WAZ405**: Detect circular dependencies
- **WAZ406**: Validate subnet address ranges within VNet range
- **WAZ407**: Require diagnostic settings for key resources

## Rule Details (Placeholders)

### WAZ001: Use Location Constants

**Description:** Use location constants for common Azure regions instead of hardcoded strings.

**Severity:** warning

**Why:** Reduces typos and improves readability.

#### Bad

```go
var MyVM = compute.VirtualMachine{
    Location: "East US",  // Incorrect format
}
```

#### Good

```go
var MyVM = compute.VirtualMachine{
    Location: "eastus",  // Correct format
}

// Or use ResourceGroup().Location
var MyVM = compute.VirtualMachine{
    Location: ResourceGroup().Location,
}
```

**Auto-fix:** Available - normalizes location strings to lowercase and removes spaces

---

### WAZ002: Use Intrinsic Types

**Description:** Use intrinsic types instead of raw `map[string]any` for ARM template functions.

**Severity:** warning

**Why:** Provides type safety and better IDE support.

#### Bad

```go
Name: map[string]any{
    "concat": []any{"storage", "account"},
}
```

#### Good

```go
Name: Concat([]any{"storage", "account"})
```

**Auto-fix:** Not available

---

### WAZ003: Extract Inline Property Types

**Description:** Extract inline property types to separate named variables.

**Severity:** warning

**Why:** Improves readability and makes code more scannable.

#### Bad

```go
var MyVM = compute.VirtualMachine{
    NetworkProfile: compute.VirtualMachine_NetworkProfile{
        NetworkInterfaces: []compute.VirtualMachine_NetworkInterfaceReference{
            {Id: "/subscriptions/.../networkInterfaces/mynic"},
        },
    },
}
```

#### Good

```go
var MyNICRef = compute.VirtualMachine_NetworkInterfaceReference{
    Id: MyNIC.Id,
}

var MyNetworkProfile = compute.VirtualMachine_NetworkProfile{
    NetworkInterfaces: []compute.VirtualMachine_NetworkInterfaceReference{MyNICRef},
}

var MyVM = compute.VirtualMachine{
    NetworkProfile: MyNetworkProfile,
}
```

**Auto-fix:** Not available

---

### WAZ004: Use Typed Structs

**Description:** Use typed structs instead of `map[string]any` for resource properties.

**Severity:** warning

**Why:** Provides compile-time type checking and IDE autocomplete.

#### Bad

```go
Sku: map[string]any{
    "name": "Standard_LRS",
    "tier": "Standard",
}
```

#### Good

```go
Sku: storage.StorageAccount_Sku{
    Name: "Standard_LRS",
    Tier: "Standard",
}
```

**Auto-fix:** Not available

---

### WAZ005: Detect Duplicate Resource Names

**Description:** Detect duplicate resource variable names in a file or package.

**Severity:** error

**Why:** Azure resource names must be unique within a template.

#### Bad

```go
var MyStorage = storage.StorageAccount{Name: "storage1"}
var MyStorage = storage.StorageAccount{Name: "storage2"}  // Duplicate!
```

#### Good

```go
var DataStorage = storage.StorageAccount{Name: "datastorage"}
var LogsStorage = storage.StorageAccount{Name: "logsstorage"}
```

**Auto-fix:** Not available

---

### WAZ006: Storage Account Naming Conventions

**Description:** Validate storage account names meet Azure constraints.

**Severity:** error

**Why:** Storage account names must be 3-24 lowercase alphanumeric characters and globally unique.

#### Bad

```go
Name: "My-Storage-Account"  // Contains hyphens and uppercase
Name: "st"                   // Too short
Name: "mystorageaccountwithaveryverylongname"  // Too long
```

#### Good

```go
Name: "mystorageaccount"

// Or use UniqueString for global uniqueness
Name: Concat([]any{"storage", UniqueString(ResourceGroup().Id)})
```

**Auto-fix:** Not available

---

### WAZ007: Avoid Hardcoded Secrets

**Description:** Detect hardcoded secrets, API keys, and sensitive credentials.

**Severity:** error

**Why:** Prevents accidental credential exposure in source code.

**Detected patterns:**
- Storage account keys
- Connection strings with credentials
- API keys and tokens
- Private keys
- Passwords in plaintext

#### Bad

```go
var MyConfig = Json{
    "ConnectionString": "DefaultEndpointsProtocol=https;AccountName=myaccount;AccountKey=XXXXXX",
    "ApiKey":          "secret-api-key-12345",
}
```

#### Good

```go
// Use Key Vault
var MyKeyVault = keyvault.Vault{...}
var MySecret = keyvault.Secret{...}

// Reference secret in configuration
var MyConfig = Json{
    "ConnectionString": Reference(MyKeyVault.Id, "connectionString"),
}

// Or use secureString parameters
var ApiKeyParam = Parameter{
    Type: "secureString",
}
```

**Auto-fix:** Not available

---

## Implementation Status

This is a placeholder document. Lint rules will be implemented in subsequent issues:

- [ ] Basic rule framework
- [ ] Type safety rules (WAZ001-099)
- [ ] Reference rules (WAZ100-199)
- [ ] Extraction rules (WAZ200-299)
- [ ] Security rules (WAZ300-399)
- [ ] Azure-specific rules (WAZ400-499)

## Disabling Rules

Currently, individual rules cannot be disabled. To skip linting, simply don't run `wetwire-azure lint`.

## Contributing

To add new lint rules:

1. Add the rule struct in `internal/linter/rules.go`
2. Implement the `Rule` interface: `ID()`, `Description()`, `Check()`
3. Add to `AllRules()` function
4. Add tests in `rules_test.go`
5. Document in this file

## See Also

- [CLI Reference]({{< relref "/cli" >}}) - How to run the linter
- [FAQ]({{< relref "/faq" >}}) - Common questions about lint rules
