# Changelog

All notable changes to wetwire-azure-go will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Enterprise application scenario in `examples/enterprise_scenario/` demonstrating multi-tier infrastructure with network, compute, and storage resources
- Scenario configuration with beginner, intermediate, and expert persona prompts for AI-assisted generation
- Support for `LintOpts.Fix` option in domain linter (reserved for future auto-fix implementation)
- Support for `LintOpts.Disable` option to skip specific lint rules by ID (e.g., `["WAZ001", "WAZ002"]`)
- `lint.Options` struct with `DisabledRules` and `Fix` fields in internal linter
- `lint.NewLinterWithOptions()` constructor for creating linter with custom options

### Changed
- Split `internal/lint/rules.go` (1,315 lines) into category-specific files for better maintainability:
  - `rules_structure.go` - WAZ001-WAZ005 (476 lines)
  - `rules_security.go` - WAZ006-WAZ008 (244 lines)
  - `rules_wetwire.go` - WAZ020-WAZ022 (307 lines)
  - `rules_azure.go` - WAZ301-WAZ304 (251 lines)
  - `rules.go` - Main registry (22 lines)
- Removed `init()` function in `internal/kiro/config.go`, using constant for `AgentPrompt` instead
- Renamed `internal/linter` to `internal/lint` for consistency with wetwire-core-go naming conventions
- Migrated discover package to use `wetwire-core-go/ast` utilities (ExtractTypeName, IsBuiltinIdent)
- Migrated linter Severity type to use `wetwire-core-go/lint` type alias
- Updated `wetwire-core-go` to v1.16.0
- Migrated MCP server to use `domain.BuildMCPServer()` for automatic tool generation
- Updated `wetwire-core-go` to v1.13.0 for automated MCP server generation
- Replaced manual MCP tool registration with auto-generated implementation

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
- `wetwire-azure mcp` command for Model Context Protocol integration
- Auto-generated MCP tools via `domain.BuildMCPServer()`:
  - `wetwire_init` - Initialize new projects
  - `wetwire_build` - Generate ARM templates
  - `wetwire_lint` - Check code quality
  - `wetwire_validate` - Validate templates
  - `wetwire_list` - List discovered resources
  - `wetwire_graph` - Visualize dependencies

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
- github.com/lex00/wetwire-core-go v1.13.0

[Unreleased]: https://github.com/lex00/wetwire-azure-go/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/lex00/wetwire-azure-go/releases/tag/v1.0.0
