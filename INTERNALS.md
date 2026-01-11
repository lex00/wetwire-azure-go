# Architecture Internals

This document explains how wetwire-azure-go transforms Go code into ARM templates.

## Build Pipeline Overview

```
Go Source Files
      |
      v
  DISCOVER  ─────>  AST Parsing
      |
      v
  VALIDATE  ─────>  Reference Check + Cycle Detection
      |
      v
    ORDER   ─────>  Topological Sort
      |
      v
 SERIALIZE  ─────>  Struct to JSON Conversion
      |
      v
    EMIT    ─────>  ARM Template JSON
```

## Phase 1: Discovery

**Location:** `internal/discover/discover.go`

The discovery phase parses Go source files using `go/ast` to find Azure resource declarations.

### How Resources Are Found

1. Walk all `.go` files in the target directory
2. Parse each file into an AST
3. Find top-level `var` declarations
4. Check if the type is from `wetwire-azure-go/resources/*`
5. Map Go types to Azure resource types

```go
// Resource type mapping
var azureResourceMap = map[string]string{
    "storage.StorageAccount": "Microsoft.Storage/storageAccounts",
    "compute.VirtualMachine": "Microsoft.Compute/virtualMachines",
    // ...
}
```

### Dependency Extraction

For each resource, the discovery phase extracts dependencies by recursively walking the value expression:

```go
var MyVM = compute.VirtualMachine{
    NetworkProfile: compute.NetworkProfile{
        NetworkInterfaces: []compute.NetworkInterfaceReference{
            {Id: MyNIC.Id},  // MyNIC is extracted as a dependency
        },
    },
}
```

### DiscoveredResource Structure

```go
type DiscoveredResource struct {
    Name         string   // Variable name (e.g., "MyVM")
    Type         string   // Azure type (e.g., "Microsoft.Compute/virtualMachines")
    File         string   // Source file path
    Line         int      // Line number
    Dependencies []string // Referenced variable names
}
```

## Phase 2: Serialization

**Location:** `internal/serialize/serialize.go`

Converts Go structs to ARM-compatible JSON maps.

### Key Behaviors

1. **JSON Tags**: Respects `json:"name"` struct tags for field names
2. **Omitempty**: Skips zero values when `json:",omitempty"` is present
3. **Nested Structs**: Recursively converts nested structures
4. **Intrinsics**: Detects `intrinsics.Intrinsic` interface and calls `ARMExpression()`

### Intrinsic Handling

When the serializer encounters a value implementing `Intrinsic`:

```go
if intrinsic, ok := v.Interface().(intrinsics.Intrinsic); ok {
    return intrinsic.ARMExpression()  // Returns "[resourceGroup().location]"
}
```

## Phase 3: Template Building

**Location:** `internal/template/template.go`

Orchestrates the build process and produces the final ARM template.

### TemplateBuilder

```go
type TemplateBuilder struct {
    resources  map[string]DiscoveredResource
    parameters map[string]Parameter
    variables  map[string]interface{}
    outputs    map[string]Output
}
```

### Build Pipeline

1. **Validate References**: Ensure all referenced resources exist
2. **Detect Cycles**: Use DFS to find circular dependencies
3. **Topological Sort**: Order resources so dependencies come first
4. **Serialize**: Convert to ARM JSON format

### Cycle Detection

Uses depth-first search with a recursion stack:

```go
func (tb *TemplateBuilder) validateReferences() error {
    visited := make(map[string]bool)
    recStack := make(map[string]bool)

    var hasCycle func(string) bool
    hasCycle = func(name string) bool {
        visited[name] = true
        recStack[name] = true

        for _, dep := range resource.Dependencies {
            if !visited[dep] {
                if hasCycle(dep) { return true }
            } else if recStack[dep] {
                return true  // Cycle found
            }
        }
        recStack[name] = false
        return false
    }
    // ...
}
```

### Topological Sort

Uses Kahn's algorithm for stable ordering:

1. Calculate in-degree for each resource
2. Start with resources having no dependencies
3. Process queue, reducing in-degrees
4. Repeat until all resources are ordered

## CLI Architecture

**Location:** `cmd/wetwire-azure/main.go`

### Command Structure

```
wetwire-azure
├── build     Generate ARM template
├── import    Convert ARM JSON to Go
├── lint      Check for issues
├── list      Show discovered resources
├── init      Create new project
├── graph     Generate dependency graph
├── validate  Validate ARM template
└── help      Show usage
```

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Build error |
| 2 | Invalid arguments |

## Linting System

**Location:** `internal/linter/`

### Rule Interface

```go
type Rule interface {
    ID() string                            // e.g., "WAZ001"
    Description() string                   // Human-readable description
    Severity() Severity                    // Error, Warning, Info
    Check(file string) ([]LintResult, error)
}
```

### Implemented Rules

| Rule | Severity | Description |
|------|----------|-------------|
| WAZ001 | Warning | Use lowercase location format |
| WAZ002 | Warning | Use direct references, not resourceId() |
| WAZ003 | Warning | Avoid deeply nested configurations |
| WAZ004 | Error | No duplicate variable names |
| WAZ005 | Error | No circular dependencies |

## Intrinsics System

**Location:** `intrinsics/intrinsics.go`

### Intrinsic Interface

```go
type Intrinsic interface {
    ARMExpression() string
}
```

### Available Intrinsics

| Function | Go Usage | ARM Output |
|----------|----------|------------|
| ResourceGroup() | `ResourceGroup().Location` | `[resourceGroup().location]` |
| Parameters() | `Parameters("vmName")` | `[parameters('vmName')]` |
| Variables() | `Variables("prefix")` | `[variables('prefix')]` |
| ResourceId() | `ResourceId(type, name)` | `[resourceId(type, name)]` |
| Ref() | `Ref(name, api)` | `[reference(name, api)]` |

## Resource Types

**Location:** `resources/`

### Package Organization

```
resources/
├── compute/        # VirtualMachine, etc.
├── network/        # VNet, NIC, etc.
├── storage/        # StorageAccount, etc.
└── keyvault/       # Vault, etc.
```

### Type Structure

Each resource type follows this pattern:

```go
type StorageAccount struct {
    Name       string                     `json:"name"`
    Type       string                     `json:"type"`
    APIVersion string                     `json:"apiVersion"`
    Location   string                     `json:"location"`
    Properties *StorageAccountProperties  `json:"properties,omitempty"`
    // ...
}
```

## Import System

**Location:** `internal/importer/importer.go`

Converts ARM JSON templates back to Go code:

1. Parse ARM JSON structure
2. Extract resources, parameters, variables
3. Generate Go variable declarations
4. Format with gofmt
