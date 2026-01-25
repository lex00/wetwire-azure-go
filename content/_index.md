---
title: "Wetwire Azure"
---

[![Go Reference](https://pkg.go.dev/badge/github.com/lex00/wetwire-azure-go.svg)](https://pkg.go.dev/github.com/lex00/wetwire-azure-go)
[![CI](https://github.com/lex00/wetwire-azure-go/actions/workflows/ci.yml/badge.svg)](https://github.com/lex00/wetwire-azure-go/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/lex00/wetwire-azure-go/graph/badge.svg)](https://codecov.io/gh/lex00/wetwire-azure-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/lex00/wetwire-azure-go)](https://goreportcard.com/report/github.com/lex00/wetwire-azure-go)
[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Semantic linting for ARM templates.

## Documentation

| Document | Description |
|----------|-------------|
| [CLI Reference]({{< relref "/cli" >}}) | Command-line interface |
| [Quick Start]({{< relref "/quick-start" >}}) | Get started in 5 minutes |
| [FAQ]({{< relref "/faq" >}}) | Frequently asked questions |

## Installation

```bash
go install github.com/lex00/wetwire-azure-go@latest
```

## Quick Example

```go
var MyStorage = storage.Account{
    Name:     "mystorageaccount",
    Location: "eastus",
    Sku:      storage.SkuStandardLRS,
}
```
