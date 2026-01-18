You generate Azure infrastructure resources using wetwire-azure-go.

## Context

**Application:** Multi-tier enterprise web application with production workload requirements

**Architecture:**
- Web tier: Public-facing layer with Application Gateway and VMSS
- App tier: Application logic layer with VMSS behind internal load balancer
- Data tier: Azure SQL Database with private endpoint
- Security: Key Vault for secrets, NSGs for network isolation

**Location:** East US

**Tags:**
- environment: production
- project: enterprise-app
- cost-center: engineering

## Output Files

Create the following Go files in the `expected/` directory:

- `expected/network.go` - Virtual network, subnets, NSGs, public IPs
- `expected/compute.go` - Virtual machine scale sets, Application Gateway
- `expected/storage.go` - Storage accounts for data and diagnostics
- `expected/database.go` - Azure SQL Database server and database
- `expected/security.go` - Key Vault and access policies

## Network Architecture

**Virtual Network: enterprise-vnet (10.0.0.0/16)**
- `gateway-subnet` (10.0.0.0/24) - Application Gateway subnet
- `web-subnet` (10.0.1.0/24) - Web tier VMSS subnet
- `app-subnet` (10.0.2.0/24) - App tier VMSS subnet
- `data-subnet` (10.0.3.0/24) - Database private endpoint subnet

**Network Security Groups:**
- `gateway-nsg` - Allow HTTP (80), HTTPS (443) from Internet, allow 65200-65535 for App Gateway management
- `web-nsg` - Allow 8080 from gateway-subnet
- `app-nsg` - Allow 8080 from web-subnet
- `data-nsg` - Allow 1433 from app-subnet

**Public IPs:**
- `appgw-pip` - Standard SKU, static allocation, zone-redundant (zones 1,2,3)

## Compute Resources

**Web Tier VMSS: web-vmss**
- VM Size: Standard_B2s
- Capacity: 2 instances
- Image: Ubuntu 22.04 LTS Gen2
- OS Disk: Premium_LRS, ReadWrite caching
- Network: web-subnet
- Load balancing: via Application Gateway backend pool

**App Tier VMSS: app-vmss**
- VM Size: Standard_D2s_v3
- Capacity: 3 instances
- Image: Ubuntu 22.04 LTS Gen2
- OS Disk: Premium_LRS, ReadWrite caching
- Network: app-subnet

**Application Gateway: enterprise-appgw**
- SKU: Standard_v2, tier Standard_v2
- Capacity: 2 instances
- Subnet: gateway-subnet
- Frontend: Public IP (appgw-pip), port 80
- Backend pool: web-vmss instances
- HTTP settings: port 8080, protocol HTTP, cookie-based affinity disabled
- Routing rule: Basic rule mapping frontend to backend pool

## Storage Resources

**Data Storage Account: enterpriseappdata**
- SKU: Standard_GRS
- Kind: StorageV2
- Access tier: Hot
- Security: HTTPS only, TLS 1.2, no public blob access, no shared key access
- Network: Deny default, allow from data-subnet via service endpoint

**Logs Storage Account: enterpriseapplogs**
- SKU: Standard_LRS
- Kind: StorageV2
- Security: HTTPS only, TLS 1.2, no public blob access
- Purpose: Diagnostic logs and metrics

## Database Resources

**Azure SQL Server: enterprise-sql-server**
- Version: 12.0
- Admin login: sqladmin
- Admin password: (reference Key Vault secret)
- Public network access: Disabled
- Minimum TLS: 1.2

**Azure SQL Database: enterprise-db**
- SKU: S1 (Standard tier)
- Max size: 250 GB
- Zone redundant: false
- Collation: SQL_Latin1_General_CP1_CI_AS

**Private Endpoint: sql-private-endpoint**
- Subnet: data-subnet
- Target: SQL Server (sqlServer)
- Private DNS zone: privatelink.database.windows.net

## Security Resources

**Key Vault: enterprise-keyvault**
- SKU: Standard
- Tenant ID: (Azure tenant)
- Access policies:
  - VMSS managed identities: Get, List secrets
- Network: Default action Deny, allow from app-subnet
- Properties: Enable RBAC authorization, soft delete enabled (90 days), purge protection enabled

**Key Vault Secrets:**
- `sql-admin-password` - SQL Server admin password
- `storage-connection-string` - Data storage account connection string

## Code Style

Use wetwire-azure-go declarative patterns:

```go
package main

import (
    "github.com/lex00/wetwire-azure-go/resources/network"
    "github.com/lex00/wetwire-azure-go/resources/compute"
    "github.com/lex00/wetwire-azure-go/resources/storage"
)

// Common configuration
var (
    location = "eastus"
    tags = map[string]string{
        "environment": "production",
        "project":     "enterprise-app",
        "cost-center": "engineering",
    }
)

// AppVNet is the main virtual network for the enterprise application.
var AppVNet = network.VirtualNetwork{
    Name:       "enterprise-vnet",
    Type:       "Microsoft.Network/virtualNetworks",
    APIVersion: "2021-05-01",
    Location:   location,
    Tags:       tags,
    Properties: network.VirtualNetworkProperties{
        AddressSpace: network.AddressSpace{
            AddressPrefixes: []string{"10.0.0.0/16"},
        },
        Subnets: []network.Subnet{
            // Define subnets inline
        },
    },
}
```

**Key points:**
- All resources as package-level `var` declarations
- Use typed structs from resources packages
- Include Name, Type, APIVersion, Location, Tags on all resources
- Add brief comments explaining each resource
- Extract nested configurations to separate variables when complex
- For pointer fields, create helper variables (e.g., `var adminUsername = "azureadmin"`)
- Use `*string`, `*bool`, `*int` for optional ARM template fields

## Validation Rules

Generated infrastructure must have:
- At least 5 Azure resources total
- At least 1 Virtual Network with 2+ subnets
- At least 1 Network Security Group
- At least 1 compute resource (VM or VMSS)
- At least 1 Storage Account
- All resources have valid ARM JSON structure (Name, Type, APIVersion, Properties)
