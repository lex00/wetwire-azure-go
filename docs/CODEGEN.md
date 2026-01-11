# Code Generation Pipeline

This document describes how wetwire-azure-go generates Azure resource types from ARM schemas.

## Overview

wetwire-azure-go uses a multi-stage pipeline to generate Go types from Azure ARM schemas:

```
ARM Schemas → Parse → Transform → Generate → Format → Output
```

## Pipeline Stages

### 1. Schema Fetching

ARM schemas are fetched from the official Azure repository:

```
https://github.com/Azure/azure-resource-manager-schemas
```

The schemas define:
- Resource types and their properties
- Required vs optional fields
- Allowed values (enums)
- Nested object structures

### 2. Schema Parsing

The `codegen` package parses JSON schemas:

```go
// Parse ARM schema
schema, err := codegen.ParseSchema("Microsoft.Storage/storageAccounts")
```

Key extracted information:
- Property names and types
- Nested object definitions
- Validation constraints
- API version

### 3. Type Transformation

Schemas are transformed to Go-friendly structures:

| ARM Schema Type | Go Type |
|----------------|---------|
| string | string |
| integer | int |
| number | float64 |
| boolean | bool |
| object | struct |
| array | []T |
| oneOf | interface{} |

### 4. Code Generation

Go source files are generated with:

```go
// Generated type example
type StorageAccount struct {
    Name       string                     `json:"name"`
    Location   string                     `json:"location"`
    Kind       string                     `json:"kind,omitempty"`
    SKU        SKU                        `json:"sku,omitempty"`
    Properties StorageAccountProperties  `json:"properties,omitempty"`
}
```

### 5. Formatting

Generated code is formatted with `go fmt` for consistency.

## Running Code Generation

```bash
# Generate all resource types
go generate ./codegen/...

# Generate specific resource type
go run ./codegen/cmd/generate -type Microsoft.Storage/storageAccounts
```

## Generated Files

Generated files are placed in `resources/<provider>/`:

```
resources/
├── compute/
│   ├── virtual_machine.go
│   └── disk.go
├── storage/
│   ├── storage_account.go
│   └── blob_container.go
└── network/
    ├── virtual_network.go
    └── network_interface.go
```

## Customization

### Custom Type Mappings

Some types have custom mappings for better ergonomics:

```go
// In codegen/mappings.go
var customMappings = map[string]string{
    "Microsoft.Storage/storageAccounts#/definitions/Sku": "SKU",
}
```

### Excluded Fields

Some schema fields are excluded from generation:

```go
var excludedFields = []string{
    "id",           // Computed by ARM
    "type",         // Set automatically
    "apiVersion",   // Set by package
}
```

### Manual Overrides

For complex types, manual implementations may be provided:

```go
//go:generate skip
// This type is manually implemented
type ComplexResource struct {
    // Custom implementation
}
```

## Schema Validation

Generated types are validated against:

1. **Completeness**: All required properties present
2. **Type Safety**: Go types match schema constraints
3. **Round-trip**: ARM → Go → ARM produces identical JSON

## Adding New Resource Types

1. Add schema reference to `codegen/schemas.json`
2. Run code generation
3. Review generated code
4. Add tests in `resources/<provider>/*_test.go`
5. Update documentation

## Troubleshooting

### Missing Properties

If a property is missing:
1. Check schema version
2. Verify property isn't in excluded list
3. Check for custom mapping

### Type Mismatches

If Go type doesn't match expected:
1. Review custom mappings
2. Check schema definition
3. Consider adding custom mapping

### Generation Errors

```bash
# Verbose output
go run ./codegen/cmd/generate -v

# Debug mode
go run ./codegen/cmd/generate -debug
```

## Architecture

```
codegen/
├── cmd/generate/      # Generation CLI
├── parser/            # Schema parsing
├── transformer/       # Schema to Go transformation
├── generator/         # Code generation
├── mappings.go        # Custom type mappings
└── schemas.json       # Schema sources
```

## See Also

- [Architecture](INTERNALS.md)
- [Developer Guide](DEVELOPERS.md)
- [Azure ARM Schema Repository](https://github.com/Azure/azure-resource-manager-schemas)
