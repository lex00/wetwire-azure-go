Generate production-grade Azure ARM resources via wetwire-azure-go:

**Network (expected/network.go):**
- VNet: enterprise-vnet (10.0.0.0/16), subnets: gateway (10.0.0.0/24), web (10.0.1.0/24), app (10.0.2.0/24), data (10.0.3.0/24)
- NSGs: gateway-nsg (80/443 ingress + 65200-65535 AppGw mgmt), web-nsg (8080 from gateway), app-nsg (8080 from web), data-nsg (1433 from app)
- Public IP: appgw-pip (Standard, Static, zone-redundant 1,2,3)

**Compute (expected/compute.go):**
- VMSS: web-vmss (B2s, Ubuntu 22.04 Gen2, Premium_LRS, 2 instances, web-subnet)
- VMSS: app-vmss (D2s_v3, Ubuntu 22.04 Gen2, Premium_LRS, 3 instances, app-subnet)
- AppGw: enterprise-appgw (Standard_v2, 2 instances, gateway-subnet, frontend appgw-pip:80, backend web-vmss:8080)

**Storage (expected/storage.go):**
- SA: enterpriseappdata (Standard_GRS, StorageV2, Hot, HTTPS+TLS1.2, deny-default+allow data-subnet)
- SA: enterpriseapplogs (Standard_LRS, StorageV2, HTTPS+TLS1.2)

**Database (expected/database.go):**
- SQL Server: enterprise-sql-server (v12.0, sqladmin, no public access, TLS1.2)
- SQL DB: enterprise-db (S1, 250GB, SQL_Latin1_General_CP1_CI_AS)
- Private Endpoint: sql-private-endpoint (data-subnet, sqlServer, privatelink.database.windows.net)

**Security (expected/security.go):**
- KeyVault: enterprise-keyvault (Standard, RBAC, soft-delete 90d, purge protection, deny-default+allow app-subnet)
- Secrets: sql-admin-password, storage-connection-string
- Access policies: VMSS managed identities (Get/List secrets)

Location: eastus | Tags: environment=production, project=enterprise-app, cost-center=engineering
