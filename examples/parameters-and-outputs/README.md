# Parameters and Outputs Example

This example demonstrates how to use ARM template intrinsic functions in wetwire-azure-go resource definitions.

## Intrinsic Functions

### Parameters

Use `intrinsics.Parameters()` to reference ARM template parameters:

```go
import "github.com/lex00/wetwire-azure-go/intrinsics"

var MyResource = storage.StorageAccount{
    Location: intrinsics.Parameters("location").ARMExpression(),
    // Generates: "[parameters('location')]"
}
```

### Resource Group

Use `intrinsics.ResourceGroup()` to access resource group properties:

```go
var MyResource = storage.StorageAccount{
    // Uses the resource group's location
    Location: intrinsics.ResourceGroup().ARMExpression(),
    // Generates: "[resourceGroup()]" or "[resourceGroup().location]"
}
```

### ResourceId

Use `intrinsics.ResourceId()` to reference other resources:

```go
var MyNIC = network.NetworkInterface{
    Properties: network.NetworkInterfaceProperties{
        IPConfigurations: []network.IPConfiguration{{
            Properties: network.IPConfigurationProperties{
                Subnet: &network.SubResource{
                    ID: strPtr(intrinsics.ResourceId(
                        "Microsoft.Network/virtualNetworks/subnets",
                        "myVnet",
                        "default",
                    ).ARMExpression()),
                },
            },
        }},
    },
}
```

### Variables

Use `intrinsics.Variables()` to reference ARM template variables:

```go
var MyResource = storage.StorageAccount{
    Location: intrinsics.Variables("defaultLocation").ARMExpression(),
    // Generates: "[variables('defaultLocation')]"
}
```

## Resources

| Resource | Purpose |
|----------|---------|
| ParameterizedStorage | Shows parameter references |
| ResourceGroupLocation | Shows resourceGroup() function |
| ExampleNIC | Shows resourceId() for subnet reference |
| VariableBasedStorage | Shows variables() function |
| ProdStorage | Production configuration pattern |
| DevStorage | Development configuration pattern |

## Environment-Based Configuration

The example shows a pattern for environment-specific configurations:

- **Production**: Geo-redundant storage, Azure AD auth only
- **Development**: Local redundancy, simpler settings

## Usage

```bash
# Build ARM template
wetwire-azure build ./examples/parameters-and-outputs

# Lint the configuration
wetwire-azure lint ./examples/parameters-and-outputs
```

## Generated ARM Template

The intrinsic functions generate proper ARM template expressions:

```json
{
  "resources": [
    {
      "type": "Microsoft.Storage/storageAccounts",
      "location": "[parameters('location')]",
      ...
    }
  ]
}
```
