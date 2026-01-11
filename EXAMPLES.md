# Examples

This document indexes the example projects demonstrating wetwire-azure-go usage.

## Example Projects

| Example | Description | Resources |
|---------|-------------|-----------|
| [storage-account](examples/storage-account/) | Basic storage account deployment | StorageAccount |
| [virtual-machine](examples/virtual-machine/) | Linux VM with managed disk | VirtualMachine |

## Running Examples

### Build an Example

```bash
cd examples/storage-account
wetwire-azure build -o template.json
```

### List Resources

```bash
wetwire-azure list ./examples/storage-account
```

### Lint an Example

```bash
wetwire-azure lint ./examples/storage-account
```

### Deploy to Azure

```bash
# Login to Azure
az login

# Create resource group
az group create --name my-rg --location eastus

# Deploy
az deployment group create \
  --resource-group my-rg \
  --template-file template.json
```

## Example: Storage Account

Location: `examples/storage-account/`

A minimal example showing a single storage account:

```go
package main

import (
    "github.com/lex00/wetwire-azure-go/resources/storage"
)

var MyStorageAccount = storage.StorageAccount{
    Name:     "mystorageaccount",
    Location: "eastus",
    SKU: storage.SKU{
        Name: "Standard_LRS",
    },
    Kind: "StorageV2",
}
```

## Example: Virtual Machine

Location: `examples/virtual-machine/`

A Linux VM deployment with Ubuntu image:

```go
package main

import (
    "github.com/lex00/wetwire-azure-go/resources/compute"
)

var LinuxVM = compute.VirtualMachine{
    Name:     "my-linux-vm",
    Location: "eastus",
    Properties: compute.VirtualMachineProperties{
        HardwareProfile: compute.HardwareProfile{
            VMSize: "Standard_B2s",
        },
        StorageProfile: compute.StorageProfile{
            ImageReference: &compute.ImageReference{...},
            OSDisk: compute.OSDisk{...},
        },
        OSProfile: &compute.OSProfile{...},
        NetworkProfile: compute.NetworkProfile{...},
    },
}
```

## Creating Your Own Example

1. Create a new directory under `examples/`
2. Add a `main.go` with resource declarations
3. Add a `README.md` explaining the example
4. Test with `wetwire-azure build`

### Template

```go
package main

import (
    "github.com/lex00/wetwire-azure-go/resources/storage"
    // Add other imports as needed
)

// Describe your resource
var MyResource = storage.StorageAccount{
    Name:     "uniquename",
    Location: "eastus",
    // ...
}
```

## Best Practices Demonstrated

### Flat Structure

Examples show extracting nested configurations:

```go
// Instead of deeply nested structs, extract to separate variables
var MyOSDisk = compute.OSDisk{...}
var MyStorageProfile = compute.StorageProfile{OSDisk: MyOSDisk}
var MyVM = compute.VirtualMachine{StorageProfile: MyStorageProfile}
```

### Direct References

Examples use direct variable references:

```go
var MyNIC = network.NetworkInterface{...}
var MyVM = compute.VirtualMachine{
    NetworkProfile: compute.NetworkProfile{
        NetworkInterfaces: []compute.NetworkInterfaceReference{
            {Id: MyNIC.Id},  // Direct reference
        },
    },
}
```

### Consistent Naming

Examples follow naming conventions:

```go
// Resource variables use PascalCase
var MyStorageAccount = storage.StorageAccount{...}

// Azure resource names use lowercase
var MyStorageAccount = storage.StorageAccount{
    Name: "mystorageaccount",  // lowercase
}
```
