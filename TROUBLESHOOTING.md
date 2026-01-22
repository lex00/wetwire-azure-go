<picture>
  <source media="(prefers-color-scheme: dark)" srcset="docs/wetwire-dark.svg">
  <img src="docs/wetwire-light.svg" width="100" height="67">
</picture>

Common issues and their solutions when using wetwire-azure-go.

## Build Errors

### "discovery failed: no resources found"

**Cause:** No Azure resource types were detected in your Go files.

**Solutions:**
1. Ensure you import from `github.com/lex00/wetwire-azure-go/resources/*`
2. Use package-level `var` declarations (not inside functions)
3. Check that resource types are spelled correctly

```go
// Wrong - local variable
func main() {
    storage := storage.StorageAccount{...}  // Not discovered
}

// Correct - package-level variable
var MyStorage = storage.StorageAccount{...}  // Discovered
```

### "cyclic dependency detected"

**Cause:** Two or more resources reference each other in a loop.

**Solution:** Break the cycle by removing one reference or restructuring:

```go
// Wrong - circular dependency
var A = compute.VirtualMachine{
    Properties: compute.VirtualMachineProperties{
        DiagnosticsProfile: compute.DiagnosticsProfile{
            StorageURI: B.Properties.PrimaryEndpoints.Blob,
        },
    },
}
var B = storage.StorageAccount{
    Tags: map[string]string{"vm": A.Name},  // Creates cycle
}

// Correct - use ARM intrinsics to break cycle
var B = storage.StorageAccount{
    Tags: map[string]string{"vm": "my-vm"},  // Static value
}
```

### "resource X depends on non-existent resource Y"

**Cause:** A resource references a variable that doesn't exist or isn't recognized as a resource.

**Solution:** Check that the referenced variable is:
1. Defined in the same package
2. A valid Azure resource type
3. Spelled correctly (case-sensitive)

## Lint Errors

### WAZ001: Location format

**Message:** `Location 'East US' should use lowercase format`

**Solution:** Use lowercase Azure region names without spaces:

```go
// Wrong
Location: "East US",

// Correct
Location: "eastus",
```

### WAZ002: Direct references

**Message:** `Use direct resource references instead of resourceId()`

**Solution:** Reference resources by variable name:

```go
// Wrong
Id: resourceId("Microsoft.Network/networkInterfaces", "mynic"),

// Correct
Id: MyNIC.Id,
```

### WAZ003: Deeply nested configuration

**Message:** `Deeply nested configuration detected`

**Solution:** Extract nested structs to separate variables:

```go
// Wrong - deeply nested
var MyVM = compute.VirtualMachine{
    Properties: compute.VirtualMachineProperties{
        StorageProfile: compute.StorageProfile{
            OSDisk: compute.OSDisk{
                ManagedDisk: compute.ManagedDiskParameters{...},
            },
        },
    },
}

// Correct - extracted
var MyOSDisk = compute.OSDisk{
    ManagedDisk: compute.ManagedDiskParameters{...},
}

var MyStorageProfile = compute.StorageProfile{
    OSDisk: MyOSDisk,
}

var MyVM = compute.VirtualMachine{
    Properties: compute.VirtualMachineProperties{
        StorageProfile: MyStorageProfile,
    },
}
```

### WAZ004: Duplicate variable names

**Message:** `Duplicate variable name 'MyStorage'`

**Solution:** Use unique variable names:

```go
// Wrong
var MyStorage = storage.StorageAccount{Name: "storage1",...}
var MyStorage = storage.StorageAccount{Name: "storage2",...}  // Duplicate!

// Correct
var PrimaryStorage = storage.StorageAccount{Name: "storage1",...}
var BackupStorage = storage.StorageAccount{Name: "storage2",...}
```

### WAZ005: Circular dependencies

**Message:** `Circular dependency detected involving variable 'A'`

**Solution:** Same as "cyclic dependency detected" above.

## Import Errors

### "cannot find package"

**Cause:** The wetwire-azure-go module is not in your Go module.

**Solution:**
```bash
go mod tidy
# or explicitly add:
go get github.com/lex00/wetwire-azure-go@latest
```

### "undefined: storage.StorageAccount"

**Cause:** Missing or incorrect import statement.

**Solution:**
```go
import (
    "github.com/lex00/wetwire-azure-go/resources/storage"
)
```

### "cannot use X as type Y"

**Cause:** Type mismatch in struct field assignment.

**Solution:** Check the expected type in the resource struct definition:

```go
// Wrong - string instead of pointer
Properties: storage.StorageAccountProperties{
    AccessTier: "Hot",  // Error: cannot use string as *string
}

// Correct - use pointer
accessTier := "Hot"
Properties: storage.StorageAccountProperties{
    AccessTier: &accessTier,
}
```

## Common Mistakes

### Using Pointers in Declarations

**Wrong:**
```go
var MyStorage = &storage.StorageAccount{...}  // Pointer not supported
```

**Correct:**
```go
var MyStorage = storage.StorageAccount{...}  // Value type
```

### Function Calls in Declarations

**Wrong:**
```go
var MyStorage = storage.NewStorageAccount("name", "eastus")  // Function call
```

**Correct:**
```go
var MyStorage = storage.StorageAccount{
    Name:     "name",
    Location: "eastus",
}
```

### Missing Required Fields

**Problem:** ARM deployment fails with missing required properties.

**Solution:** Check Azure documentation for required fields:

```go
var MyStorage = storage.StorageAccount{
    Name:     "mystorageaccount",  // Required
    Location: "eastus",            // Required
    SKU: storage.SKU{
        Name: "Standard_LRS",      // Required
    },
    Kind: "StorageV2",             // Required
}
```

## Deployment Errors

### "StorageAccountAlreadyTaken"

**Cause:** Storage account names must be globally unique.

**Solution:** Use a unique name or generate one with ARM functions:

```go
import . "github.com/lex00/wetwire-azure-go/intrinsics"

var MyStorage = storage.StorageAccount{
    Name: UniqueString{Values: []string{ResourceGroup().Id}},
}
```

### "InvalidResourceReference"

**Cause:** A resource reference in the ARM template is invalid.

**Solution:** Check that all referenced resources are included in the deployment and have correct names.

## Getting Help

1. Run `wetwire-azure lint` to check for common issues
2. Use `wetwire-azure list --format json` to see discovered resources
3. Generate a dependency graph: `wetwire-azure graph`
4. Check the [INTERNALS.md](INTERNALS.md) for architecture details
