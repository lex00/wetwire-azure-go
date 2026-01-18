Create Azure infrastructure for a multi-tier enterprise application:

**Network:**
- Virtual network (10.0.0.0/16) with subnets for gateway, web, app, and data tiers
- Network security groups for each tier with appropriate rules
- Public IP for Application Gateway

**Compute:**
- Web tier VMSS (2 instances, Standard_B2s) in web-subnet
- App tier VMSS (3 instances, Standard_D2s_v3) in app-subnet
- Application Gateway (Standard_v2) for internet-facing traffic

**Storage:**
- GRS storage account for application data (geo-redundant)
- LRS storage account for diagnostic logs

**Database:**
- Azure SQL Database (S1) with private endpoint
- SQL Server with public network access disabled

**Security:**
- Key Vault for storing secrets (SQL password, connection strings)
- Managed identities for VMSS to access Key Vault
- Network restrictions on all resources

Location: East US
Tags: environment=production, project=enterprise-app, cost-center=engineering
