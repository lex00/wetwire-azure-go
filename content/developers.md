---
title: "Developers"
---
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="./wetwire-dark.svg">
  <img src="./wetwire-light.svg" width="100" height="67">
</picture>

This guide is for contributors to wetwire-azure-go.

## Prerequisites

- Go 1.23+
- Git
- golangci-lint (for linting)

## Getting Started

```bash
# Clone the repository
git clone https://github.com/lex00/wetwire-azure-go.git
cd wetwire-azure-go

# Install dependencies
go mod download

# Run tests
go test ./...

# Run linter
golangci-lint run ./...
```

## Project Structure

```
wetwire-azure-go/
├── cmd/
│   ├── wetwire-azure/      # Main CLI
│   └── wetwire-azure-mcp/  # MCP server
├── internal/
│   ├── discover/           # AST-based resource discovery
│   ├── importer/           # ARM JSON to Go conversion
│   ├── linter/             # Lint rules engine
│   ├── serialize/          # Go to ARM JSON serialization
│   ├── template/           # ARM template building
│   └── validator/          # Schema validation
├── intrinsics/             # ARM template functions
├── resources/              # Generated resource types
│   ├── compute/
│   └── storage/
├── codegen/                # Type generation from schemas
├── examples/               # Example projects
└── docs/                   # Documentation
```

## Development Workflow

### 1. Create a Branch

```bash
git checkout -b feat/my-feature
# or
git checkout -b fix/my-bugfix
```

### 2. Make Changes

Follow these conventions:
- Write tests first (TDD)
- Keep functions small and focused
- Document exported functions

### 3. Run Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/linter/...
```

### 4. Run Linter

```bash
golangci-lint run ./...
```

### 5. Commit

Follow conventional commits:
```bash
git commit -m "feat: add new feature"
git commit -m "fix: resolve bug in linter"
git commit -m "docs: update README"
git commit -m "test: add tests for importer"
```

### 6. Push and Create PR

```bash
git push -u origin feat/my-feature
gh pr create
```

## Adding a New Lint Rule

1. Create the rule struct in `internal/linter/rules.go`:

```go
type WAZ009 struct{}

func (r *WAZ009) ID() string          { return "WAZ009" }
func (r *WAZ009) Description() string { return "Description" }
func (r *WAZ009) Severity() Severity  { return SeverityWarning }
func (r *WAZ009) Check(file string) ([]LintResult, error) {
    // Implementation
}
```

2. Register in `AllRules()`:

```go
func AllRules() []Rule {
    return []Rule{
        // ...existing rules...
        &WAZ009{},
    }
}
```

3. Add tests in `internal/linter/rules_test.go`

4. Document in `docs/LINT_RULES.md`

### Adding Auto-Fix

Implement `FixableRule` interface:

```go
func (r *WAZ009) CanFix() bool { return true }
func (r *WAZ009) Fix(file string) (string, error) {
    // Return fixed file content
}
```

## Adding a New Resource Type

1. Define struct in `resources/<provider>/<type>.go`
2. Add serialization support in `internal/serialize/`
3. Add discovery support in `internal/discover/`
4. Add import support in `internal/importer/`
5. Add examples in `examples/`

## Testing

### Unit Tests

```bash
go test ./...
```

### Coverage Targets

- Discovery: 90%+
- Serialization: 90%+
- Linting: 80%+
- CLI: 70%+

### Running Specific Tests

```bash
# By name pattern
go test -run TestWAZ001 ./...

# Verbose output
go test -v ./internal/linter/...
```

## Code Style

- Use `gofmt` for formatting (automatic with most editors)
- Follow standard Go conventions
- Keep line length reasonable (~100 chars)
- Use meaningful variable names

## Documentation

- Document all exported types and functions
- Keep examples up to date
- Update CHANGELOG.md for notable changes

## Release Process

1. Update CHANGELOG.md
2. Create release tag: `git tag v1.x.x`
3. Push tag: `git push origin v1.x.x`
4. GitHub Actions will create the release

## Getting Help

- Open an issue for bugs or features
- Check existing issues before creating new ones
- Include reproduction steps for bugs

## See Also

- [Architecture](INTERNALS.md)
- [CLI Reference](CLI.md)
- [Contributing Guidelines](../CONTRIBUTING.md)
