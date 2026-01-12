# Security Best Practices Example

This example demonstrates Azure security best practices for infrastructure-as-code deployments.

## Security Principles

### 1. Defense in Depth

Multiple layers of security controls:
- Network Security Groups with explicit deny rules
- Service endpoints for private PaaS access
- Storage account network restrictions
- Encryption at rest and in transit

### 2. Least Privilege

- NSGs default to deny-all with explicit allow rules
- Storage accounts deny public access by default
- Azure AD authentication instead of shared keys

### 3. Zero Trust

- No implicit trust between network segments
- All traffic must be explicitly allowed
- Service endpoints for secure PaaS communication

## Resources

### Network Security

| Resource | Purpose |
|----------|---------|
| SecureNSG | Deny-all default with explicit allows |
| JumpboxNSG | SSH access from corporate network only |
| SecureVNet | VNet with service endpoints |

### Storage Security

| Resource | Features |
|----------|----------|
| SecureStorage | Full security configuration |
| AuditStorage | Audit logs with geo-redundancy |
| EncryptedStorage | Encryption at rest configuration |

## NSG Best Practices

```go
// Always include explicit deny rules
SecurityRules: []network.SecurityRule{
    {
        Name: "deny-all-inbound",
        Properties: network.SecurityRuleProperties{
            Priority:  4096,  // Lowest priority
            Direction: "Inbound",
            Access:    "Deny",
            Protocol:  "*",
            // ... deny all
        },
    },
}
```

## Storage Security Checklist

- [x] HTTPS-only traffic (`EnableHTTPSTrafficOnly: true`)
- [x] TLS 1.2 minimum (`MinimumTLSVersion: "TLS1_2"`)
- [x] No public blob access (`AllowBlobPublicAccess: false`)
- [x] Azure AD authentication only (`AllowSharedKeyAccess: false`)
- [x] Network restrictions (`NetworkRuleSet.DefaultAction: "Deny"`)
- [x] Service endpoints for VNet access
- [x] Geo-redundant storage for critical data

## Service Endpoints

Enable private access to Azure PaaS services:

```go
ServiceEndpoints: []network.ServiceEndpoint{
    {Service: "Microsoft.Storage", Locations: []string{"eastus"}},
    {Service: "Microsoft.Sql", Locations: []string{"eastus"}},
    {Service: "Microsoft.KeyVault", Locations: []string{"eastus"}},
},
```

## Usage

```bash
# Build ARM template
wetwire-azure build ./examples/security-best-practices

# Lint for security issues
wetwire-azure lint ./examples/security-best-practices
```

## Compliance Considerations

This example includes tags for compliance tracking:
- `security: high`
- `compliance: pci-dss`
- `data-class: confidential`

Adapt these tags to your organization's compliance requirements.

## Additional Recommendations

1. **Key Vault**: Store secrets and certificates in Azure Key Vault
2. **Managed Identities**: Use system-assigned identities for VM authentication
3. **Azure Policy**: Enforce security baselines across subscriptions
4. **Azure Defender**: Enable threat protection for storage and VMs
5. **Diagnostic Logs**: Send logs to a SIEM or Log Analytics workspace
