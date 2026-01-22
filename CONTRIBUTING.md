<picture>
  <source media="(prefers-color-scheme: dark)" srcset="docs/wetwire-dark.svg">
  <img src="docs/wetwire-light.svg" width="100" height="67">
</picture>

This guide explains how to contribute to wetwire-azure-go.

## Setting Up the Development Environment

### Prerequisites

- Go 1.21 or later
- Git

### Clone and Build

```bash
git clone https://github.com/lex00/wetwire-azure-go.git
cd wetwire-azure-go

# Build the CLI
go build -o wetwire-azure ./cmd/wetwire-azure

# Verify it works
./wetwire-azure help
```

### IDE Setup

For VS Code, install the Go extension and ensure these settings:

```json
{
    "go.useLanguageServer": true,
    "go.lintTool": "golangci-lint"
}
```

## Running Tests

### Run All Tests

```bash
go test ./...
```

### Run Tests with Coverage

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Run Specific Package Tests

```bash
go test ./internal/discover/...
go test ./internal/linter/...
go test ./internal/template/...
```

### Run with Verbose Output

```bash
go test -v ./...
```

## Project Structure

```
wetwire-azure-go/
├── cmd/
│   └── wetwire-azure/      # CLI entry point
├── internal/
│   ├── discover/           # AST-based resource discovery
│   ├── serialize/          # Struct to JSON conversion
│   ├── template/           # ARM template builder
│   ├── linter/             # Lint rules
│   ├── validator/          # Template validation
│   └── importer/           # ARM to Go conversion
├── resources/              # Azure resource type definitions
│   ├── compute/
│   ├── network/
│   └── storage/
├── intrinsics/             # ARM template function wrappers
├── examples/               # Example projects
└── docs/                   # Additional documentation
```

## Adding New Resource Types

### 1. Create the Type Definition

Create a new file in the appropriate `resources/` subdirectory:

```go
// resources/network/virtualnetwork.go
package network

type VirtualNetwork struct {
    Name       string                    `json:"name"`
    Type       string                    `json:"type"`
    APIVersion string                    `json:"apiVersion"`
    Location   string                    `json:"location"`
    Tags       map[string]string         `json:"tags,omitempty"`
    Properties VirtualNetworkProperties  `json:"properties"`
}

type VirtualNetworkProperties struct {
    AddressSpace AddressSpace `json:"addressSpace"`
    Subnets      []Subnet     `json:"subnets,omitempty"`
}
// ...
```

### 2. Update the Discovery Map

Add the new type to `internal/discover/discover.go`:

```go
var azureResourceMap = map[string]string{
    "network.VirtualNetwork": "Microsoft.Network/virtualNetworks",
    // ...existing entries...
}
```

### 3. Add API Version

Update `internal/template/template.go`:

```go
func getAPIVersion(resourceType string) string {
    apiVersions := map[string]string{
        "Microsoft.Network/virtualNetworks": "2021-02-01",
        // ...existing entries...
    }
    // ...
}
```

### 4. Write Tests

Create tests for the new resource type:

```go
// resources/network/virtualnetwork_test.go
package network

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestVirtualNetworkFields(t *testing.T) {
    vnet := VirtualNetwork{
        Name:     "myvnet",
        Location: "eastus",
        // ...
    }
    assert.Equal(t, "myvnet", vnet.Name)
}
```

## Adding New Lint Rules

### 1. Create the Rule

Add a new rule in `internal/linter/rules.go`:

```go
// WAZ006 checks for...
type WAZ006 struct{}

func (r *WAZ006) ID() string { return "WAZ006" }

func (r *WAZ006) Description() string {
    return "Description of what this rule checks"
}

func (r *WAZ006) Severity() Severity { return SeverityWarning }

func (r *WAZ006) Check(file string) ([]LintResult, error) {
    // Parse file and check for issues
    fset := token.NewFileSet()
    node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
    if err != nil {
        return nil, err
    }

    var results []LintResult
    // ... inspection logic ...
    return results, nil
}
```

### 2. Register the Rule

Add it to `AllRules()` in `internal/linter/rules.go`:

```go
func AllRules() []Rule {
    return []Rule{
        &WAZ001{},
        &WAZ002{},
        // ...
        &WAZ006{},  // Add new rule
    }
}
```

### 3. Write Tests

```go
func TestWAZ006(t *testing.T) {
    rule := &WAZ006{}
    assert.Equal(t, "WAZ006", rule.ID())

    // Test with a file that should trigger the rule
    results, err := rule.Check("testdata/waz006_trigger.go")
    assert.NoError(t, err)
    assert.Len(t, results, 1)
}
```

## Code Style

### Formatting

All code must be formatted with `gofmt`:

```bash
gofmt -w .
```

### Imports

Use standard library imports first, then external packages:

```go
import (
    "fmt"
    "go/ast"

    "github.com/stretchr/testify/assert"
)
```

### Comments

- All exported types and functions must have doc comments
- Comments should be complete sentences

```go
// StorageAccount represents a Microsoft.Storage/storageAccounts resource.
// It provides type-safe access to storage account configuration.
type StorageAccount struct {
    // ...
}
```

### Error Handling

- Always check errors
- Wrap errors with context using `fmt.Errorf`:

```go
if err != nil {
    return fmt.Errorf("failed to parse file %s: %w", filename, err)
}
```

### Testing

- Use table-driven tests where appropriate
- Test both success and error cases
- Use `testify/assert` for assertions

## Pull Request Process

1. Create a feature branch from `main`
2. Make your changes
3. Run tests: `go test ./...`
4. Format code: `gofmt -w .`
5. Commit with a descriptive message
6. Push and create a pull request

### Commit Messages

Use conventional commits format:

```
feat: add VirtualNetwork resource type
fix: correct API version for storage accounts
docs: update troubleshooting guide
test: add tests for linter rule WAZ006
```

## Documentation

- Update relevant documentation when making changes
- Add examples for new features
- Keep the TROUBLESHOOTING.md updated with common issues

## Questions?

- Check existing issues and documentation
- Open an issue for discussion
