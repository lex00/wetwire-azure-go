# wetwire-azure (Go)

Generate Azure ARM/Bicep templates from Go resource declarations.

## Syntax Principles

All resources are Go struct literals. No function calls, no pointers, no registration.

### Resource Declaration

Resources are declared as package-level variables:

```go
var DataStorage = storage.StorageAccount{
    Name:     "mystorageaccount",
    Location: "eastus",
    Sku: storage.StorageAccount_Sku{
        Name: "Standard_LRS",
    },
}
```

### Direct References

Reference other resources directly by variable name:

```go
var MyVM = compute.VirtualMachine{
    Name:     "my-vm",
    Location: "eastus",
    NetworkProfile: compute.VirtualMachine_NetworkProfile{
        NetworkInterfaces: []compute.VirtualMachine_NetworkInterfaceReference{
            {Id: MyNIC.Id},  // Direct reference
        },
    },
}
```

### Nested Types

Extract nested configurations to separate variables:

```go
var MyNICConfig = network.NetworkInterface_IpConfiguration{
    Name:   "ipconfig1",
    Subnet: network.NetworkInterface_Subnet{
        Id: MySubnet.Id,
    },
}

var MyNIC = network.NetworkInterface{
    Name:                 "my-nic",
    Location:            "eastus",
    IpConfigurations:    []network.NetworkInterface_IpConfiguration{MyNICConfig},
}
```

### Azure Intrinsics

For ARM template functions, use the intrinsics package:

```go
import (
    . "github.com/lex00/wetwire-azure-go/intrinsics"
    "github.com/lex00/wetwire-azure-go/resources/storage"
)

var MyStorage = storage.StorageAccount{
    Name: Concat([]any{ResourceGroup().Name, "-storage"}),
    Location: ResourceGroup().Location,
}
```

**intrinsics provides:**
- `Concat`, `UniqueString`, `Guid`, `Uri` - String manipulation
- `ResourceGroup()`, `Subscription()`, `Deployment()` - Context functions
- `ResourceId`, `Reference` - Resource references
- `Parameters`, `Variables` - Template parameters

## Package Structure

```
wetwire-azure-go/
├── resources/           # Generated Azure resource types
│   ├── compute/        # Virtual machines, disks, etc.
│   ├── network/        # VNets, NICs, load balancers
│   ├── storage/        # Storage accounts, blobs
│   ├── keyvault/       # Key vaults, secrets
│   ├── webapp/         # App services, plans
│   └── k8s/            # K8s CRD types (via ASO) for kubectl deployments
│       ├── containerservice/ # ManagedCluster (AKS) CRDs
│       ├── network/          # VirtualNetwork, Subnet CRDs
│       └── managedidentity/  # UserAssignedIdentity CRDs
├── intrinsics/         # ARM template functions
├── internal/
│   ├── discover/       # AST-based resource discovery
│   ├── template/       # ARM template builder
│   ├── lint/           # Lint rules (WAZ001-WAZ999)
│   └── importer/       # JSON/Bicep to Go conversion
└── cmd/wetwire-azure/  # CLI application
```

## Lint Rules

Uses the `WAZ` prefix (Wetwire AZure). See [docs/LINT_RULES.md](docs/LINT_RULES.md) for the complete rule reference.

## Key Principles

1. **Flat variables** - Extract all nested structs into named variables
2. **No pointers** - Never use `&` or `*` in declarations
3. **Direct references** - Variables reference each other by name
4. **Struct literals only** - No function calls in declarations

## Build

```bash
wetwire-azure build ./infra > template.json
# or
wetwire-azure build ./infra --format bicep > template.bicep
```

## K8s-Native Deployments (ASO)

The `resources/k8s/` directory contains Azure Service Operator (ASO) types for deploying Azure resources via `kubectl apply`. This provides a Kubernetes-native alternative to ARM templates.

### Using ASO Types

```go
import (
    aksv1 "github.com/lex00/wetwire-azure-go/resources/k8s/containerservice/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Cluster = aksv1.ManagedCluster{
    TypeMeta: metav1.TypeMeta{
        APIVersion: "containerservice.azure.com/v1",
        Kind:       "ManagedCluster",
    },
    ObjectMeta: metav1.ObjectMeta{
        Name:      "my-cluster",
        Namespace: "aso-system",
    },
    Spec: aksv1.ManagedClusterSpec{
        Location:          "eastus",
        KubernetesVersion: "1.28",
    },
}
```

### When to Use ASO vs ARM

| Approach | Use When |
|----------|----------|
| **ARM** (`resources/compute/`, etc.) | Traditional IaC, Azure-native tooling, existing ARM pipelines |
| **ASO** (`resources/k8s/`) | GitOps workflows, Kubernetes-centric teams, unified K8s API |

See `examples/aks-golden/` for ARM approach and `examples/aks-k8s/` for ASO approach.

## Azure-Specific Patterns

### Resource Naming

Azure resource names have specific constraints:

```go
// Storage accounts: 3-24 lowercase alphanumeric
var DataStorage = storage.StorageAccount{
    Name: "mystorageaccount123",  // Must be globally unique
}

// Virtual machines: 1-64 chars (Windows), 1-15 chars (Linux)
var LinuxVM = compute.VirtualMachine{
    Name: "webserver01",
}
```

### Location Management

Use consistent location references:

```go
// Option 1: Direct string (for single-region deployments)
var WebVM = compute.VirtualMachine{
    Location: "eastus",
}

// Option 2: Resource group location (recommended)
var WebVM = compute.VirtualMachine{
    Location: ResourceGroup().Location,
}

// Option 3: Parameter (for multi-region deployments)
var LocationParam = Parameter{Type: "string", DefaultValue: "eastus"}
var WebVM = compute.VirtualMachine{
    Location: LocationParam,
}
```

### Managed Identities

```go
var MyVMIdentity = compute.VirtualMachine_Identity{
    Type: "SystemAssigned",
}

var MyVM = compute.VirtualMachine{
    Identity: MyVMIdentity,
}

// Reference the managed identity in role assignments
var MyRoleAssignment = authorization.RoleAssignment{
    PrincipalId: MyVM.Identity.PrincipalId,
}
```

### Dependencies

Azure ARM templates use implicit dependencies via resource references:

```go
// VM depends on NIC (implicit via Id reference)
var MyVM = compute.VirtualMachine{
    NetworkProfile: compute.VirtualMachine_NetworkProfile{
        NetworkInterfaces: []compute.VirtualMachine_NetworkInterfaceReference{
            {Id: MyNIC.Id},  // Creates dependency: MyVM -> MyNIC
        },
    },
}
```

## Project Structure

```
my-azure-infra/
├── go.mod
├── network.go         # VNets, subnets, NSGs
├── compute.go         # VMs, scale sets
├── storage.go         # Storage accounts, disks
└── security.go        # Key vaults, managed identities
```

## When Editing This Repository

### Repository Organization

This is a domain package implementing wetwire for Azure:
- `resources/` contains generated Azure resource type definitions
- `intrinsics/` contains ARM template function wrappers
- `internal/discover/` implements AST-based resource discovery
- `internal/template/` builds ARM JSON or Bicep output
- `internal/lint/` enforces declarative patterns (WAZ rules)
- `cmd/wetwire-azure/` is the CLI entry point

### Key Files

- `/Users/alex/Documents/checkouts/wetwire/docs/WETWIRE_SPEC.md` - Core wetwire philosophy and requirements
- `/Users/alex/Documents/checkouts/wetwire/docs/DOCUMENTATION_GUIDE.md` - Documentation standards

### Common Commands

```bash
# Build example project
wetwire-azure build ./examples/basic-vm

# Run linter
wetwire-azure lint ./examples/basic-vm

# Import existing ARM template
wetwire-azure import template.json -o ./output

# Validate generated template
wetwire-azure validate ./infra

# Run tests
go test ./...
```

### Adding New Resource Types

Azure resource types should be generated from Azure Resource Manager schemas:

1. Update schema definitions in code generation pipeline
2. Regenerate resource types: `go generate ./resources/...`
3. Add tests for new resource types
4. Update documentation

### Lint Rule Development

Lint rules follow the WAZ prefix (Wetwire AZure). See [docs/LINT_RULES.md](docs/LINT_RULES.md) for the complete rule reference with category ranges WAZ001-WAZ499.
