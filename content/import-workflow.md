---
title: "Import Workflow"
---

This document describes the workflow for importing existing ARM templates into wetwire-azure-go.

## Overview

The `wetwire-azure import` command converts ARM JSON templates to Go code, enabling you to manage existing infrastructure with wetwire.

## Basic Import

```bash
# Import a single ARM template
wetwire-azure import template.json

# Import with custom output file
wetwire-azure import template.json -o resources.go

# Import with custom package name
wetwire-azure import template.json --package myinfra
```

## Workflow Steps

### 1. Export ARM Template

Export your existing ARM template from the Azure portal or CLI:

```bash
# Using Azure CLI
az group export --name my-resource-group > template.json

# Or download from Azure portal
# Resource Group -> Export template -> Download
```

### 2. Import to Go

Run the import command:

```bash
wetwire-azure import template.json -o infra/main.go --package infra
```

### 3. Review Generated Code

The importer will:
- Convert resources to typed Go structs
- Preserve resource names and properties
- Generate appropriate imports
- Handle dependencies between resources

### 4. Lint and Fix

Run the linter to ensure code quality:

```bash
# Check for issues
wetwire-azure lint ./infra

# Auto-fix where possible
wetwire-azure lint ./infra --fix
```

### 5. Build and Validate

Verify the round-trip works correctly:

```bash
# Generate ARM template from Go code
wetwire-azure build ./infra -o generated.json

# Validate against Azure schema
wetwire-azure validate generated.json

# Compare with original (should be semantically equivalent)
wetwire-azure diff generated.json --against template.json
```

## Supported Resource Types

The importer supports:
- `Microsoft.Storage/storageAccounts`
- `Microsoft.Compute/virtualMachines`
- `Microsoft.Network/virtualNetworks`
- `Microsoft.Network/networkInterfaces`

## Handling Unsupported Resources

For unsupported resource types, the importer will:
1. Generate a placeholder struct
2. Add a comment indicating manual review needed
3. Preserve the raw JSON for reference

## Best Practices

1. **Start Small**: Import one resource group at a time
2. **Review Generated Code**: Always review and understand the generated code
3. **Run Tests**: Ensure the generated templates deploy successfully
4. **Incremental Migration**: Gradually move resources to wetwire management

## Troubleshooting

### Invalid JSON
If the import fails with JSON parsing errors, validate your template:
```bash
jq . template.json
```

### Missing Dependencies
If resources have dependencies not in the template, you may need to:
1. Export the dependent resources
2. Import them first
3. Update references in your code

## See Also

- [CLI Reference](CLI.md)
- [FAQ](FAQ.md)
- [Adoption Guide](ADOPTION.md)
