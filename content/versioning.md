---
title: "Versioning"
---
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="./wetwire-dark.svg">
  <img src="./wetwire-light.svg" width="100" height="67">
</picture>

This document describes the versioning policy for wetwire-azure-go.

## Semantic Versioning

wetwire-azure-go follows [Semantic Versioning 2.0.0](https://semver.org/):

```
MAJOR.MINOR.PATCH
```

- **MAJOR**: Breaking changes to public API
- **MINOR**: New features, backwards compatible
- **PATCH**: Bug fixes, backwards compatible

## Public API

The following are considered part of the public API:

### CLI Commands
- Command names and subcommands
- Required flags
- Exit codes
- Output format (JSON schema)

### Go Packages
- Exported types in `resources/*`
- Exported functions in `intrinsics`
- Exported types in `internal/linter` (Rule, LintResult, Severity)

### ARM Template Output
- Schema version
- Resource structure
- Expression format

## Stability Guarantees

### Stable (v1.x.x)
- No breaking changes without major version bump
- Deprecated features marked and maintained for one minor version
- Security fixes backported to supported versions

### Pre-release (v0.x.x)
- API may change between minor versions
- No long-term support guarantees
- Use with caution in production

## Go Version Support

| wetwire-azure-go | Minimum Go Version |
|------------------|-------------------|
| v1.x.x           | Go 1.23           |
| v2.x.x (future)  | TBD               |

We support the two most recent Go versions.

## Azure API Versions

Resource types include their Azure API version:

```go
var MyStorage = storage.StorageAccount{
    // Uses API version defined in storage package
}
```

### API Version Updates

- Patch releases may update API versions for bug fixes
- Minor releases may add support for new API versions
- Resource packages document supported API versions

## Dependency Policy

### Direct Dependencies

| Dependency | Policy |
|------------|--------|
| Go standard library | Follow Go version support |
| testify | Latest minor version |

### Updating Dependencies

```bash
# Update all dependencies
go get -u ./...

# Update specific dependency
go get -u github.com/stretchr/testify

# Tidy module
go mod tidy
```

## Breaking Changes

Before introducing a breaking change:

1. Mark existing API as deprecated
2. Provide migration path in CHANGELOG
3. Maintain deprecated API for one minor release
4. Remove in next major version

### Deprecation Notice Format

```go
// Deprecated: Use NewFunction instead. Will be removed in v2.0.0.
func OldFunction() {}
```

## Release Schedule

- **Patch releases**: As needed for bug fixes
- **Minor releases**: Monthly (if new features available)
- **Major releases**: As needed for breaking changes

## Version Checking

Check your installed version:

```bash
wetwire-azure version
```

Check latest available:

```bash
go list -m -versions github.com/lex00/wetwire-azure-go
```

## Migration Guides

Major version migration guides are provided in:
- CHANGELOG.md (summary)
- Dedicated migration doc (detailed)

## See Also

- [Changelog](../CHANGELOG.md)
- [Developer Guide](DEVELOPERS.md)
- [Go Modules Reference](https://go.dev/ref/mod)
