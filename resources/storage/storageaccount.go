// Package storage provides Azure storage resource types
package storage

// StorageAccount represents a Microsoft.Storage/storageAccounts resource
type StorageAccount struct {
	// Name is the name of the storage account (3-24 characters, lowercase letters and numbers only)
	Name string `json:"name"`

	// Type is the resource type
	Type string `json:"type"`

	// APIVersion is the API version to use for this resource
	APIVersion string `json:"apiVersion"`

	// Location is the Azure region where the resource will be created
	Location string `json:"location"`

	// Tags are key-value pairs to organize resources
	Tags map[string]string `json:"tags,omitempty"`

	// Kind is the kind of storage account (Storage, StorageV2, BlobStorage, FileStorage, BlockBlobStorage)
	Kind string `json:"kind"`

	// SKU defines the SKU/pricing tier for the storage account
	SKU SKU `json:"sku"`

	// Properties contains the properties of the storage account
	Properties *StorageAccountProperties `json:"properties,omitempty"`

	// Identity defines the identity configuration for the storage account
	Identity *Identity `json:"identity,omitempty"`
}

// SKU represents the SKU of a storage account
type SKU struct {
	// Name is the SKU name (Standard_LRS, Standard_GRS, Standard_RAGRS, Standard_ZRS, Premium_LRS, Premium_ZRS)
	Name string `json:"name"`

	// Tier is the SKU tier (Standard or Premium)
	Tier *string `json:"tier,omitempty"`
}

// StorageAccountProperties represents the properties of a storage account
type StorageAccountProperties struct {
	// AccessTier defines the access tier (Hot or Cool)
	AccessTier *string `json:"accessTier,omitempty"`

	// AllowBlobPublicAccess indicates whether public access is allowed to all blobs
	AllowBlobPublicAccess *bool `json:"allowBlobPublicAccess,omitempty"`

	// AllowSharedKeyAccess indicates whether the storage account permits requests via Shared Key
	AllowSharedKeyAccess *bool `json:"allowSharedKeyAccess,omitempty"`

	// EnableHTTPSTrafficOnly indicates whether only HTTPS traffic is allowed
	EnableHTTPSTrafficOnly *bool `json:"supportsHttpsTrafficOnly,omitempty"`

	// MinimumTLSVersion sets the minimum TLS version
	MinimumTLSVersion *string `json:"minimumTlsVersion,omitempty"`

	// NetworkRuleSet defines network access rules
	NetworkRuleSet *NetworkRuleSet `json:"networkAcls,omitempty"`

	// Encryption defines the encryption settings
	Encryption *Encryption `json:"encryption,omitempty"`

	// IsHnsEnabled indicates whether hierarchical namespace is enabled (for Data Lake Gen2)
	IsHnsEnabled *bool `json:"isHnsEnabled,omitempty"`

	// LargeFileSharesState indicates whether large file shares are enabled
	LargeFileSharesState *string `json:"largeFileSharesState,omitempty"`
}

// NetworkRuleSet represents network access control rules
type NetworkRuleSet struct {
	// DefaultAction specifies the default action when no rule matches (Allow or Deny)
	DefaultAction string `json:"defaultAction"`

	// Bypass specifies which services can bypass the network rules
	Bypass *string `json:"bypass,omitempty"`

	// IPRules defines IP access control rules
	IPRules []IPRule `json:"ipRules,omitempty"`

	// VirtualNetworkRules defines virtual network access control rules
	VirtualNetworkRules []VirtualNetworkRule `json:"virtualNetworkRules,omitempty"`
}

// IPRule represents an IP access control rule
type IPRule struct {
	// Value is the IP address or CIDR range
	Value string `json:"value"`

	// Action is the action to take (Allow)
	Action *string `json:"action,omitempty"`
}

// VirtualNetworkRule represents a virtual network access control rule
type VirtualNetworkRule struct {
	// ID is the resource ID of the subnet
	ID string `json:"id"`

	// Action is the action to take (Allow)
	Action *string `json:"action,omitempty"`

	// State is the state of the rule
	State *string `json:"state,omitempty"`
}

// Encryption represents encryption settings
type Encryption struct {
	// KeySource indicates the encryption key source (Microsoft.Storage or Microsoft.Keyvault)
	KeySource string `json:"keySource"`

	// Services defines which services are encrypted
	Services *EncryptionServices `json:"services,omitempty"`

	// KeyVaultProperties defines Key Vault properties for customer-managed keys
	KeyVaultProperties *KeyVaultProperties `json:"keyvaultproperties,omitempty"`
}

// EncryptionServices represents encryption settings for storage services
type EncryptionServices struct {
	// Blob defines blob service encryption
	Blob *EncryptionService `json:"blob,omitempty"`

	// File defines file service encryption
	File *EncryptionService `json:"file,omitempty"`

	// Table defines table service encryption
	Table *EncryptionService `json:"table,omitempty"`

	// Queue defines queue service encryption
	Queue *EncryptionService `json:"queue,omitempty"`
}

// EncryptionService represents encryption settings for a specific service
type EncryptionService struct {
	// Enabled indicates whether encryption is enabled
	Enabled bool `json:"enabled"`

	// KeyType specifies the encryption key type (Account or Service)
	KeyType *string `json:"keyType,omitempty"`
}

// KeyVaultProperties represents Key Vault properties for customer-managed keys
type KeyVaultProperties struct {
	// KeyName is the name of the Key Vault key
	KeyName *string `json:"keyname,omitempty"`

	// KeyVersion is the version of the Key Vault key
	KeyVersion *string `json:"keyversion,omitempty"`

	// KeyVaultURI is the URI of the Key Vault
	KeyVaultURI *string `json:"keyvaulturi,omitempty"`
}

// Identity represents the identity configuration
type Identity struct {
	// Type is the identity type (SystemAssigned, UserAssigned, SystemAssigned,UserAssigned, None)
	Type string `json:"type"`

	// UserAssignedIdentities contains user-assigned managed identities
	UserAssignedIdentities map[string]UserAssignedIdentity `json:"userAssignedIdentities,omitempty"`
}

// UserAssignedIdentity represents a user-assigned managed identity
type UserAssignedIdentity struct {
	// ClientID is the client ID of the identity
	ClientID *string `json:"clientId,omitempty"`

	// PrincipalID is the principal ID of the identity
	PrincipalID *string `json:"principalId,omitempty"`
}

// NewStorageAccount creates a new storage account with required fields
func NewStorageAccount(name, location, kind, skuName string) *StorageAccount {
	return &StorageAccount{
		Name:       name,
		Type:       "Microsoft.Storage/storageAccounts",
		APIVersion: "2021-04-01",
		Location:   location,
		Kind:       kind,
		SKU: SKU{
			Name: skuName,
		},
	}
}

// WithTags adds tags to the storage account
func (s *StorageAccount) WithTags(tags map[string]string) *StorageAccount {
	s.Tags = tags
	return s
}

// WithHTTPSOnly enables HTTPS-only traffic
func (s *StorageAccount) WithHTTPSOnly(enabled bool) *StorageAccount {
	if s.Properties == nil {
		s.Properties = &StorageAccountProperties{}
	}
	s.Properties.EnableHTTPSTrafficOnly = &enabled
	return s
}

// WithMinTLSVersion sets the minimum TLS version
func (s *StorageAccount) WithMinTLSVersion(version string) *StorageAccount {
	if s.Properties == nil {
		s.Properties = &StorageAccountProperties{}
	}
	s.Properties.MinimumTLSVersion = &version
	return s
}
