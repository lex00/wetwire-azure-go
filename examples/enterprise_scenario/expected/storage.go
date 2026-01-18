// Package main demonstrates storage infrastructure for a multi-tier enterprise application.
package main

import (
	"github.com/lex00/wetwire-azure-go/resources/storage"
)

// Helper variables for storage configuration
var (
	httpsOnly     = true
	tlsVersion    = "TLS1_2"
	noBlobPublic  = false
	noSharedKey   = false
	accessTierHot = "Hot"
	allowAction   = "Allow"
)

// DataStorage is the storage account for application data.
// Configured with geo-redundant storage, HTTPS-only, TLS 1.2, and restricted network access.
var DataStorage = storage.StorageAccount{
	Name:       "enterpriseappdata",
	Type:       "Microsoft.Storage/storageAccounts",
	APIVersion: "2021-04-01",
	Location:   location,
	Tags:       tags,
	Kind:       "StorageV2",
	SKU: storage.SKU{
		Name: "Standard_GRS",
	},
	Properties: &storage.StorageAccountProperties{
		AccessTier:             &accessTierHot,
		EnableHTTPSTrafficOnly: &httpsOnly,
		MinimumTLSVersion:      &tlsVersion,
		AllowBlobPublicAccess:  &noBlobPublic,
		AllowSharedKeyAccess:   &noSharedKey,
		NetworkRuleSet: &storage.NetworkRuleSet{
			DefaultAction: "Deny",
			VirtualNetworkRules: []storage.VirtualNetworkRule{
				{
					ID:     "/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.Network/virtualNetworks/enterprise-vnet/subnets/data-subnet",
					Action: &allowAction,
				},
				{
					ID:     "/subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.Network/virtualNetworks/enterprise-vnet/subnets/app-subnet",
					Action: &allowAction,
				},
			},
		},
	},
}

// LogStorage is the storage account for diagnostic logs and metrics.
// Configured with locally-redundant storage, HTTPS-only, and TLS 1.2.
var LogStorage = storage.StorageAccount{
	Name:       "enterpriseapplogs",
	Type:       "Microsoft.Storage/storageAccounts",
	APIVersion: "2021-04-01",
	Location:   location,
	Tags:       tags,
	Kind:       "StorageV2",
	SKU: storage.SKU{
		Name: "Standard_LRS",
	},
	Properties: &storage.StorageAccountProperties{
		AccessTier:             &accessTierHot,
		EnableHTTPSTrafficOnly: &httpsOnly,
		MinimumTLSVersion:      &tlsVersion,
		AllowBlobPublicAccess:  &noBlobPublic,
	},
}
