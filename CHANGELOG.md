# Changelog

All notable changes to wetwire-azure-go will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- CHANGELOG.md for tracking version history
- Unit tests for intrinsics package (100% coverage)
- Security lint rules:
  - WAZ006: Detect secrets and credentials (AWS keys, GitHub tokens, passwords)
  - WAZ007: Detect sensitive file paths (.env, .pem, .key)
  - WAZ008: Detect insecure defaults (HTTP, public access)

## [1.0.0] - 2026-01-11

### Added

#### CLI Commands
- `build` - Generate ARM templates from Go resource definitions
- `lint` - Check Azure resource definitions for issues (WAZ001-WAZ005 rules)
- `import` - Convert ARM JSON templates to Go code
- `validate` - Validate ARM template JSON structure against schema
- `list` - Display discovered resources in table or JSON format
- `init` - Initialize new wetwire-azure project structure
- `graph` - Generate dependency graphs in DOT or Mermaid format
- `diff` - Compare generated templates against existing files
- `watch` - Auto-rebuild on source file changes with debouncing

#### MCP Server
- `wetwire-azure-mcp` binary for Claude Code integration
- MCP tools: build, lint, import
- `--install` flag for Claude Code configuration instructions

#### Core Infrastructure
- AST-based resource discovery from Go source files
- Topological sorting for resource dependency ordering
- Cycle detection in resource dependencies
- Go struct to ARM JSON serialization
- ARM JSON to Go code import/generation
- ARM template schema validation

#### Lint Rules
- WAZ001: Location format - Use lowercase region names
- WAZ002: Direct references - Use resource references instead of resourceId()
- WAZ003: Nested configuration - Avoid deeply nested structures
- WAZ004: Duplicate names - No duplicate variable names
- WAZ005: Circular dependencies - Detect cycles in resource graphs

#### Resource Types
- `compute.VirtualMachine` - Azure Virtual Machine configuration
- `storage.StorageAccount` - Azure Storage Account configuration

#### Intrinsics
- `ResourceId()` - Generate ARM resourceId() expressions
- `Ref()` - Generate ARM reference expressions
- `RefProperty()` - Generate ARM reference with property access
- `Parameters()` - Reference ARM template parameters
- `Variables()` - Reference ARM template variables
- `ResourceGroup()` - Access resource group context
- `Subscription()` - Access subscription context
- `Concat()` - String concatenation
- `UniqueString()` - Generate unique strings

#### Documentation
- README.md - Project overview and quick start
- CLAUDE.md - AI assistant integration guide
- QUICK_START.md - 5-minute getting started guide
- INTERNALS.md - Architecture and design documentation
- EXAMPLES.md - Example project documentation
- TROUBLESHOOTING.md - Common issues and solutions
- CONTRIBUTING.md - Contribution guidelines
- docs/CLI.md - Complete CLI reference
- docs/FAQ.md - Frequently asked questions
- docs/LINT_RULES.md - Lint rule documentation

#### Examples
- `examples/storage-account/` - Simple storage account example
- `examples/virtual-machine/` - Linux VM with networking

#### Testing
- Unit tests for all core components
- Round-trip tests for ARM JSON conversion
- Integration tests for CLI commands

### Dependencies
- Go 1.23+
- github.com/stretchr/testify v1.11.1
- github.com/lex00/wetwire-core-go v1.3.0

[Unreleased]: https://github.com/lex00/wetwire-azure-go/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/lex00/wetwire-azure-go/releases/tag/v1.0.0
