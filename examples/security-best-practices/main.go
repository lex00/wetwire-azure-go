// Package main demonstrates Azure security best practices.
// This example shows:
// - Network Security Groups with deny-all default rules
// - Storage accounts with encryption and access controls
// - Service endpoints for secure PaaS access
// - Defense in depth with multiple security layers
package main

import (
	"github.com/lex00/wetwire-azure-go/resources/network"
	"github.com/lex00/wetwire-azure-go/resources/storage"
)

// Common configuration
var (
	location = "eastus"
	tags     = map[string]string{
		"environment":  "production",
		"security":     "high",
		"compliance":   "pci-dss",
		"data-class":   "confidential",
	}
)

// ============================================================================
// Network Security - Defense in Depth
// ============================================================================

// SecureNSG implements a deny-all default with explicit allow rules.
// This follows the principle of least privilege for network access.
var SecureNSG = network.NetworkSecurityGroup{
	Name:       "secure-nsg",
	Type:       "Microsoft.Network/networkSecurityGroups",
	APIVersion: "2021-05-01",
	Location:   location,
	Tags:       tags,
	Properties: network.NetworkSecurityGroupProperties{
		SecurityRules: []network.SecurityRule{
			// Explicit deny rules at lowest priority
			{
				Name: "deny-all-inbound",
				Properties: network.SecurityRuleProperties{
					Priority:                 4096,
					Direction:                "Inbound",
					Access:                   "Deny",
					Protocol:                 "*",
					SourcePortRange:          "*",
					DestinationPortRange:     "*",
					SourceAddressPrefix:      "*",
					DestinationAddressPrefix: "*",
					Description:              strPtr("Deny all inbound traffic by default"),
				},
			},
			{
				Name: "deny-all-outbound",
				Properties: network.SecurityRuleProperties{
					Priority:                 4096,
					Direction:                "Outbound",
					Access:                   "Deny",
					Protocol:                 "*",
					SourcePortRange:          "*",
					DestinationPortRange:     "*",
					SourceAddressPrefix:      "*",
					DestinationAddressPrefix: "*",
					Description:              strPtr("Deny all outbound traffic by default"),
				},
			},
			// Allow only necessary outbound traffic
			{
				Name: "allow-azure-outbound",
				Properties: network.SecurityRuleProperties{
					Priority:                 100,
					Direction:                "Outbound",
					Access:                   "Allow",
					Protocol:                 "Tcp",
					SourcePortRange:          "*",
					DestinationPortRange:     "443",
					SourceAddressPrefix:      "VirtualNetwork",
					DestinationAddressPrefix: "AzureCloud",
					Description:              strPtr("Allow HTTPS to Azure services"),
				},
			},
			// Allow internal VNet communication
			{
				Name: "allow-vnet-internal",
				Properties: network.SecurityRuleProperties{
					Priority:                 200,
					Direction:                "Inbound",
					Access:                   "Allow",
					Protocol:                 "*",
					SourcePortRange:          "*",
					DestinationPortRange:     "*",
					SourceAddressPrefix:      "VirtualNetwork",
					DestinationAddressPrefix: "VirtualNetwork",
					Description:              strPtr("Allow internal VNet traffic"),
				},
			},
		},
	},
}

// JumpboxNSG allows SSH only from specific IP ranges (e.g., corporate network).
var JumpboxNSG = network.NetworkSecurityGroup{
	Name:       "jumpbox-nsg",
	Type:       "Microsoft.Network/networkSecurityGroups",
	APIVersion: "2021-05-01",
	Location:   location,
	Tags:       tags,
	Properties: network.NetworkSecurityGroupProperties{
		SecurityRules: []network.SecurityRule{
			{
				Name: "allow-ssh-from-corp",
				Properties: network.SecurityRuleProperties{
					Priority:                 100,
					Direction:                "Inbound",
					Access:                   "Allow",
					Protocol:                 "Tcp",
					SourcePortRange:          "*",
					DestinationPortRange:     "22",
					SourceAddressPrefix:      "203.0.113.0/24", // Replace with corporate IP range
					DestinationAddressPrefix: "*",
					Description:              strPtr("Allow SSH from corporate network only"),
				},
			},
			{
				Name: "deny-all-inbound",
				Properties: network.SecurityRuleProperties{
					Priority:                 4096,
					Direction:                "Inbound",
					Access:                   "Deny",
					Protocol:                 "*",
					SourcePortRange:          "*",
					DestinationPortRange:     "*",
					SourceAddressPrefix:      "*",
					DestinationAddressPrefix: "*",
					Description:              strPtr("Deny all other inbound traffic"),
				},
			},
		},
	},
}

// ============================================================================
// Virtual Network with Service Endpoints
// ============================================================================

// SecureVNet demonstrates a network with service endpoints for secure PaaS access.
var SecureVNet = network.VirtualNetwork{
	Name:       "secure-vnet",
	Type:       "Microsoft.Network/virtualNetworks",
	APIVersion: "2021-05-01",
	Location:   location,
	Tags:       tags,
	Properties: network.VirtualNetworkProperties{
		AddressSpace: network.AddressSpace{
			AddressPrefixes: []string{"10.0.0.0/16"},
		},
		Subnets: []network.Subnet{
			{
				Name: "private-subnet",
				Properties: network.SubnetProperties{
					AddressPrefix: "10.0.1.0/24",
					// Service endpoints enable private access to PaaS services
					ServiceEndpoints: []network.ServiceEndpoint{
						{Service: "Microsoft.Storage", Locations: []string{"eastus", "westus"}},
						{Service: "Microsoft.Sql", Locations: []string{"eastus"}},
						{Service: "Microsoft.KeyVault", Locations: []string{"eastus"}},
					},
					// Disable private endpoint network policies for Private Link
					PrivateEndpointNetworkPolicies: strPtr("Disabled"),
				},
			},
			{
				Name: "public-subnet",
				Properties: network.SubnetProperties{
					AddressPrefix: "10.0.2.0/24",
				},
			},
		},
		// Enable DDoS protection for internet-facing resources
		EnableDdosProtection: boolPtr(true),
	},
}

// ============================================================================
// Storage Security Best Practices
// ============================================================================

// SecureStorage demonstrates a storage account with all security features enabled.
var SecureStorage = storage.StorageAccount{
	Name:       "securestorage",
	Type:       "Microsoft.Storage/storageAccounts",
	APIVersion: "2021-04-01",
	Location:   location,
	Tags:       tags,
	Kind:       "StorageV2",
	SKU: storage.SKU{
		Name: "Standard_GRS", // Geo-redundant for data protection
	},
	Properties: &storage.StorageAccountProperties{
		// Encryption settings
		EnableHTTPSTrafficOnly: boolPtr(true),  // Force HTTPS
		MinimumTLSVersion:      strPtr("TLS1_2"), // Modern TLS only

		// Access controls
		AllowBlobPublicAccess: boolPtr(false), // No anonymous access
		AllowSharedKeyAccess:  boolPtr(false), // Force Azure AD auth

		// Network restrictions
		NetworkRuleSet: &storage.NetworkRuleSet{
			DefaultAction: "Deny", // Deny by default
			Bypass:        strPtr("AzureServices"),
			VirtualNetworkRules: []storage.VirtualNetworkRule{
				{
					// Allow access only from specific subnet
					ID:     "/subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Network/virtualNetworks/secure-vnet/subnets/private-subnet",
					Action: strPtr("Allow"),
				},
			},
			IPRules: []storage.IPRule{
				{
					// Allow from specific trusted IP (e.g., CI/CD)
					Value:  "203.0.113.50",
					Action: strPtr("Allow"),
				},
			},
		},
	},
}

// AuditStorage is a separate storage account for audit logs.
// Uses immutable storage to prevent log tampering.
var AuditStorage = storage.StorageAccount{
	Name:       "auditlogstorage",
	Type:       "Microsoft.Storage/storageAccounts",
	APIVersion: "2021-04-01",
	Location:   location,
	Tags: map[string]string{
		"environment": "production",
		"purpose":     "audit-logs",
		"retention":   "7-years",
	},
	Kind: "StorageV2",
	SKU: storage.SKU{
		Name: "Standard_RAGRS", // Read-access geo-redundant for DR
	},
	Properties: &storage.StorageAccountProperties{
		EnableHTTPSTrafficOnly: boolPtr(true),
		MinimumTLSVersion:      strPtr("TLS1_2"),
		AllowBlobPublicAccess:  boolPtr(false),
		AllowSharedKeyAccess:   boolPtr(false),
		// Note: Immutable storage policies are configured at container level
	},
}

// ============================================================================
// Encryption at Rest Configuration
// ============================================================================

// EncryptedStorage shows customer-managed key configuration pattern.
// Note: Requires Key Vault setup (not shown in this example).
var EncryptedStorage = storage.StorageAccount{
	Name:       "encryptedstorage",
	Type:       "Microsoft.Storage/storageAccounts",
	APIVersion: "2021-04-01",
	Location:   location,
	Tags:       tags,
	Kind:       "StorageV2",
	SKU: storage.SKU{
		Name: "Standard_LRS",
	},
	Properties: &storage.StorageAccountProperties{
		EnableHTTPSTrafficOnly: boolPtr(true),
		MinimumTLSVersion:      strPtr("TLS1_2"),
		AllowBlobPublicAccess:  boolPtr(false),
		// Note: Customer-managed keys require additional KeyVault configuration
		// This demonstrates the security baseline
		Encryption: &storage.Encryption{
			KeySource: "Microsoft.Storage", // Platform-managed keys
			Services: &storage.EncryptionServices{
				Blob: &storage.EncryptionService{
					Enabled: true,
					KeyType: strPtr("Account"),
				},
				File: &storage.EncryptionService{
					Enabled: true,
					KeyType: strPtr("Account"),
				},
			},
		},
	},
}

// Helper functions
func boolPtr(b bool) *bool    { return &b }
func strPtr(s string) *string { return &s }

func main() {
	// wetwire-azure build discovers resources via AST parsing
}
