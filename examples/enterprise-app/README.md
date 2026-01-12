# Enterprise Application Example

This example demonstrates a multi-tier enterprise application infrastructure with proper network segmentation and security controls.

## Architecture

```
                    ┌─────────────────────────────────────────┐
                    │           Azure Resource Group          │
                    │                                         │
┌──────────┐        │  ┌─────────────────────────────────┐   │
│ Internet │───────────│      Web Tier (10.0.1.0/24)     │   │
└──────────┘        │  │  ┌──────┐  ┌──────────────────┐ │   │
                    │  │  │WebNSG│  │     WebVM        │ │   │
                    │  │  └──────┘  └──────────────────┘ │   │
                    │  └─────────────────────────────────┘   │
                    │                    │                    │
                    │                    ▼                    │
                    │  ┌─────────────────────────────────┐   │
                    │  │      App Tier (10.0.2.0/24)     │   │
                    │  │  ┌──────┐  ┌──────────────────┐ │   │
                    │  │  │AppNSG│  │     AppVM        │ │   │
                    │  │  └──────┘  └──────────────────┘ │   │
                    │  └─────────────────────────────────┘   │
                    │                    │                    │
                    │                    ▼                    │
                    │  ┌─────────────────────────────────┐   │
                    │  │     Data Tier (10.0.3.0/24)     │   │
                    │  │  ┌────────────┐ ┌────────────┐  │   │
                    │  │  │ DataStorage│ │ LogStorage │  │   │
                    │  │  └────────────┘ └────────────┘  │   │
                    │  │    (Service Endpoint)           │   │
                    │  └─────────────────────────────────┘   │
                    └─────────────────────────────────────────┘
```

## Resources

| Resource | Type | Description |
|----------|------|-------------|
| AppVNet | Virtual Network | Main network with 3 subnets |
| WebNSG | Network Security Group | Controls web tier access (HTTP/HTTPS) |
| AppNSG | Network Security Group | Controls app tier access |
| WebPublicIP | Public IP Address | Zone-redundant Standard SKU |
| WebVM | Virtual Machine | Ubuntu 22.04 web server |
| AppVM | Virtual Machine | Ubuntu 22.04 application server |
| DataStorage | Storage Account | Application data with network restrictions |
| LogStorage | Storage Account | Diagnostic logs |

## Security Features

- **Network Segmentation**: Three-tier architecture with separate subnets
- **NSG Rules**: Explicit allow rules with implicit deny
- **Storage Security**:
  - HTTPS-only traffic
  - TLS 1.2 minimum
  - No public blob access
  - Service endpoint for data subnet
- **Zone Redundancy**: Public IP distributed across 3 availability zones

## Usage

```bash
# Build ARM template
wetwire-azure build ./examples/enterprise-app

# Lint the configuration
wetwire-azure lint ./examples/enterprise-app

# List discovered resources
wetwire-azure list ./examples/enterprise-app
```

## Customization

1. Update the `location` variable for your preferred Azure region
2. Replace placeholder subscription/resource group IDs in NIC references
3. Use SSH keys instead of passwords for VM authentication
4. Add additional VMs for high availability
