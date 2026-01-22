<picture>
  <source media="(prefers-color-scheme: dark)" srcset="docs/wetwire-dark.svg">
  <img src="docs/wetwire-light.svg" width="100" height="67" align="right">
</picture>


# wetwire-azure-go

[![CI](https://github.com/lex00/wetwire-azure-go/actions/workflows/ci.yml/badge.svg)](https://github.com/lex00/wetwire-azure-go/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/lex00/wetwire-azure-go/branch/main/graph/badge.svg)](https://codecov.io/gh/lex00/wetwire-azure-go)
[![Go Reference](https://pkg.go.dev/badge/github.com/lex00/wetwire-azure-go.svg)](https://pkg.go.dev/github.com/lex00/wetwire-azure-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/lex00/wetwire-azure-go)](https://goreportcard.com/report/github.com/lex00/wetwire-azure-go)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Go implementation of wetwire for Azure ARM/Bicep templates.

## Overview

wetwire-azure-go enables defining Azure resources using Go code, with full type safety and IDE support.

## Installation

```bash
go install github.com/lex00/wetwire-azure-go/cmd/wetwire-azure@latest
```

## Quick Start

```go
package main

import (
    "github.com/lex00/wetwire-azure-go/resources/storage"
)

// MyStorageAccount defines a basic Azure storage account.
// Resources are declared as package-level variables for AST discovery.
var MyStorageAccount = storage.StorageAccount{
    Name:     "mystorageaccount",
    Location: "eastus",
    SKU: storage.SKU{
        Name: "Standard_LRS",
    },
    Kind: "StorageV2",
}

func main() {
    // wetwire-azure build discovers resources via AST parsing
    // No runtime execution is needed
}
```

## Commands

- `wetwire-azure build` - Generate ARM/Bicep templates from Go code
- `wetwire-azure lint` - Lint Azure resource definitions
- `wetwire-azure validate` - Validate against Azure schemas
- `wetwire-azure design` - AI-assisted infrastructure design
- `wetwire-azure test` - Run synthesis tests with AI personas
- `wetwire-azure mcp` - Start MCP server for Claude Code integration

## AI-Assisted Design

Use the `design` command for interactive, AI-assisted infrastructure creation:

```bash
# No API key required - uses Claude CLI
wetwire-azure design "Create a storage account with geo-redundant storage"
```

Uses [Claude CLI](https://claude.ai/download) by default (no API key required). Falls back to Anthropic API if Claude CLI is not installed.

## Documentation

**Getting Started:**
- [Quick Start](docs/QUICK_START.md) - 5-minute tutorial
- [FAQ](docs/FAQ.md) - Common questions

**Reference:**
- [CLI Reference](docs/CLI.md) - All commands
- [Lint Rules](docs/LINT_RULES.md) - WAZ rule reference

**Advanced:**
- [Internals](docs/INTERNALS.md) - Architecture and extension points
- [Adoption Guide](docs/ADOPTION.md) - Team migration strategies
- [Examples](docs/EXAMPLES.md) - Example projects

## License

Apache 2.0
