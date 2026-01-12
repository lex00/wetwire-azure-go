# wetwire-azure-go

[![CI](https://github.com/lex00/wetwire-azure-go/actions/workflows/ci.yml/badge.svg)](https://github.com/lex00/wetwire-azure-go/actions/workflows/ci.yml)
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
- `wetwire-azure test` - Run synthesis tests with AI personas

## License

Apache 2.0
