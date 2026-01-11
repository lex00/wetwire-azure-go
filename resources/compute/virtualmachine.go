// Package compute provides Azure compute resource types
package compute

// VirtualMachine represents a Microsoft.Compute/virtualMachines resource
type VirtualMachine struct {
	// Name is the name of the virtual machine
	Name string `json:"name"`

	// Type is the resource type
	Type string `json:"type"`

	// APIVersion is the API version to use for this resource
	APIVersion string `json:"apiVersion"`

	// Location is the Azure region where the resource will be created
	Location string `json:"location"`

	// Tags are key-value pairs to organize resources
	Tags map[string]string `json:"tags,omitempty"`

	// Properties contains the properties of the virtual machine
	Properties VirtualMachineProperties `json:"properties"`

	// Identity defines the identity configuration for the virtual machine
	Identity *Identity `json:"identity,omitempty"`

	// Zones defines availability zones for the virtual machine
	Zones []string `json:"zones,omitempty"`

	// Plan defines the marketplace image plan
	Plan *Plan `json:"plan,omitempty"`
}

// VirtualMachineProperties represents the properties of a virtual machine
type VirtualMachineProperties struct {
	// HardwareProfile specifies the hardware settings
	HardwareProfile HardwareProfile `json:"hardwareProfile"`

	// StorageProfile specifies the storage settings
	StorageProfile StorageProfile `json:"storageProfile"`

	// OSProfile specifies the operating system settings
	OSProfile *OSProfile `json:"osProfile,omitempty"`

	// NetworkProfile specifies the network interfaces
	NetworkProfile NetworkProfile `json:"networkProfile"`

	// DiagnosticsProfile specifies boot diagnostics settings
	DiagnosticsProfile *DiagnosticsProfile `json:"diagnosticsProfile,omitempty"`

	// AvailabilitySet specifies the availability set
	AvailabilitySet *SubResource `json:"availabilitySet,omitempty"`

	// LicenseType specifies the license type (Windows_Client, Windows_Server, None)
	LicenseType *string `json:"licenseType,omitempty"`

	// Priority specifies the priority (Regular, Low, Spot)
	Priority *string `json:"priority,omitempty"`

	// EvictionPolicy specifies the eviction policy for Spot VMs (Deallocate, Delete)
	EvictionPolicy *string `json:"evictionPolicy,omitempty"`

	// BillingProfile specifies billing settings
	BillingProfile *BillingProfile `json:"billingProfile,omitempty"`
}

// HardwareProfile specifies the hardware settings for a virtual machine
type HardwareProfile struct {
	// VMSize specifies the size of the virtual machine
	VMSize string `json:"vmSize"`
}

// StorageProfile specifies the storage settings for a virtual machine
type StorageProfile struct {
	// ImageReference specifies the image to use
	ImageReference *ImageReference `json:"imageReference,omitempty"`

	// OSDisk specifies the operating system disk
	OSDisk OSDisk `json:"osDisk"`

	// DataDisks specifies the data disks
	DataDisks []DataDisk `json:"dataDisks,omitempty"`
}

// ImageReference specifies an image to use
type ImageReference struct {
	// Publisher is the image publisher
	Publisher *string `json:"publisher,omitempty"`

	// Offer is the image offer
	Offer *string `json:"offer,omitempty"`

	// SKU is the image SKU
	SKU *string `json:"sku,omitempty"`

	// Version is the image version
	Version *string `json:"version,omitempty"`

	// ID is the resource ID of a custom image
	ID *string `json:"id,omitempty"`
}

// OSDisk specifies the operating system disk
type OSDisk struct {
	// Name is the disk name
	Name *string `json:"name,omitempty"`

	// CreateOption specifies how the disk should be created (FromImage, Empty, Attach)
	CreateOption string `json:"createOption"`

	// Caching specifies the caching type (None, ReadOnly, ReadWrite)
	Caching *string `json:"caching,omitempty"`

	// DiskSizeGB specifies the disk size in GB
	DiskSizeGB *int `json:"diskSizeGB,omitempty"`

	// ManagedDisk specifies managed disk parameters
	ManagedDisk *ManagedDiskParameters `json:"managedDisk,omitempty"`

	// OSType specifies the operating system type (Windows, Linux)
	OSType *string `json:"osType,omitempty"`
}

// DataDisk specifies a data disk
type DataDisk struct {
	// Name is the disk name
	Name *string `json:"name,omitempty"`

	// Lun specifies the logical unit number
	Lun int `json:"lun"`

	// CreateOption specifies how the disk should be created (FromImage, Empty, Attach)
	CreateOption string `json:"createOption"`

	// Caching specifies the caching type (None, ReadOnly, ReadWrite)
	Caching *string `json:"caching,omitempty"`

	// DiskSizeGB specifies the disk size in GB
	DiskSizeGB *int `json:"diskSizeGB,omitempty"`

	// ManagedDisk specifies managed disk parameters
	ManagedDisk *ManagedDiskParameters `json:"managedDisk,omitempty"`
}

// ManagedDiskParameters specifies managed disk settings
type ManagedDiskParameters struct {
	// StorageAccountType specifies the storage account type (Standard_LRS, Premium_LRS, StandardSSD_LRS, UltraSSD_LRS)
	StorageAccountType *string `json:"storageAccountType,omitempty"`

	// ID is the resource ID of an existing managed disk
	ID *string `json:"id,omitempty"`

	// DiskEncryptionSet specifies the disk encryption set
	DiskEncryptionSet *SubResource `json:"diskEncryptionSet,omitempty"`
}

// OSProfile specifies the operating system settings
type OSProfile struct {
	// ComputerName is the computer name
	ComputerName *string `json:"computerName,omitempty"`

	// AdminUsername is the administrator username
	AdminUsername *string `json:"adminUsername,omitempty"`

	// AdminPassword is the administrator password
	AdminPassword *string `json:"adminPassword,omitempty"`

	// CustomData is custom data passed to the VM (base64 encoded)
	CustomData *string `json:"customData,omitempty"`

	// WindowsConfiguration specifies Windows-specific settings
	WindowsConfiguration *WindowsConfiguration `json:"windowsConfiguration,omitempty"`

	// LinuxConfiguration specifies Linux-specific settings
	LinuxConfiguration *LinuxConfiguration `json:"linuxConfiguration,omitempty"`

	// Secrets specifies certificates to install
	Secrets []VaultSecretGroup `json:"secrets,omitempty"`
}

// WindowsConfiguration specifies Windows-specific settings
type WindowsConfiguration struct {
	// ProvisionVMAgent indicates whether to provision the VM agent
	ProvisionVMAgent *bool `json:"provisionVMAgent,omitempty"`

	// EnableAutomaticUpdates indicates whether automatic updates are enabled
	EnableAutomaticUpdates *bool `json:"enableAutomaticUpdates,omitempty"`

	// TimeZone specifies the time zone
	TimeZone *string `json:"timeZone,omitempty"`

	// WinRM specifies Windows Remote Management settings
	WinRM *WinRMConfiguration `json:"winRM,omitempty"`
}

// LinuxConfiguration specifies Linux-specific settings
type LinuxConfiguration struct {
	// DisablePasswordAuthentication indicates whether password authentication is disabled
	DisablePasswordAuthentication *bool `json:"disablePasswordAuthentication,omitempty"`

	// SSH specifies SSH settings
	SSH *SSHConfiguration `json:"ssh,omitempty"`

	// ProvisionVMAgent indicates whether to provision the VM agent
	ProvisionVMAgent *bool `json:"provisionVMAgent,omitempty"`
}

// SSHConfiguration specifies SSH settings
type SSHConfiguration struct {
	// PublicKeys specifies SSH public keys
	PublicKeys []SSHPublicKey `json:"publicKeys,omitempty"`
}

// SSHPublicKey specifies an SSH public key
type SSHPublicKey struct {
	// Path is the path where the public key is stored
	Path *string `json:"path,omitempty"`

	// KeyData is the SSH public key certificate
	KeyData *string `json:"keyData,omitempty"`
}

// WinRMConfiguration specifies Windows Remote Management settings
type WinRMConfiguration struct {
	// Listeners specifies WinRM listeners
	Listeners []WinRMListener `json:"listeners,omitempty"`
}

// WinRMListener specifies a WinRM listener
type WinRMListener struct {
	// Protocol is the listener protocol (Http, Https)
	Protocol *string `json:"protocol,omitempty"`

	// CertificateURL is the certificate URL
	CertificateURL *string `json:"certificateUrl,omitempty"`
}

// VaultSecretGroup specifies a group of certificates from Key Vault
type VaultSecretGroup struct {
	// SourceVault specifies the Key Vault
	SourceVault *SubResource `json:"sourceVault,omitempty"`

	// VaultCertificates specifies the certificates
	VaultCertificates []VaultCertificate `json:"vaultCertificates,omitempty"`
}

// VaultCertificate specifies a certificate from Key Vault
type VaultCertificate struct {
	// CertificateURL is the certificate URL
	CertificateURL *string `json:"certificateUrl,omitempty"`

	// CertificateStore is the certificate store (Windows only)
	CertificateStore *string `json:"certificateStore,omitempty"`
}

// NetworkProfile specifies network settings
type NetworkProfile struct {
	// NetworkInterfaces specifies the network interfaces
	NetworkInterfaces []NetworkInterfaceReference `json:"networkInterfaces"`
}

// NetworkInterfaceReference represents a reference to a network interface
type NetworkInterfaceReference struct {
	// ID is the resource ID of the network interface
	ID string `json:"id"`

	// Primary indicates whether this is the primary network interface
	Primary *bool `json:"primary,omitempty"`
}

// DiagnosticsProfile specifies boot diagnostics settings
type DiagnosticsProfile struct {
	// BootDiagnostics specifies boot diagnostics settings
	BootDiagnostics *BootDiagnostics `json:"bootDiagnostics,omitempty"`
}

// BootDiagnostics specifies boot diagnostics settings
type BootDiagnostics struct {
	// Enabled indicates whether boot diagnostics are enabled
	Enabled *bool `json:"enabled,omitempty"`

	// StorageURI is the storage account URI for boot diagnostics
	StorageURI *string `json:"storageUri,omitempty"`
}

// BillingProfile specifies billing settings
type BillingProfile struct {
	// MaxPrice is the maximum price for a Spot VM
	MaxPrice *float64 `json:"maxPrice,omitempty"`
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

// Plan represents a marketplace image plan
type Plan struct {
	// Name is the plan name
	Name *string `json:"name,omitempty"`

	// Publisher is the plan publisher
	Publisher *string `json:"publisher,omitempty"`

	// Product is the plan product
	Product *string `json:"product,omitempty"`
}

// SubResource represents a reference to another resource
type SubResource struct {
	// ID is the resource ID
	ID *string `json:"id,omitempty"`
}

// NewVirtualMachine creates a new virtual machine with required fields
func NewVirtualMachine(name, location, vmSize string) *VirtualMachine {
	return &VirtualMachine{
		Name:       name,
		Type:       "Microsoft.Compute/virtualMachines",
		APIVersion: "2021-07-01",
		Location:   location,
		Properties: VirtualMachineProperties{
			HardwareProfile: HardwareProfile{
				VMSize: vmSize,
			},
			StorageProfile: StorageProfile{
				OSDisk: OSDisk{
					CreateOption: "FromImage",
				},
			},
			NetworkProfile: NetworkProfile{
				NetworkInterfaces: []NetworkInterfaceReference{},
			},
		},
	}
}

// WithTags adds tags to the virtual machine
func (vm *VirtualMachine) WithTags(tags map[string]string) *VirtualMachine {
	vm.Tags = tags
	return vm
}

// WithImage sets the OS image
func (vm *VirtualMachine) WithImage(publisher, offer, sku, version string) *VirtualMachine {
	vm.Properties.StorageProfile.ImageReference = &ImageReference{
		Publisher: &publisher,
		Offer:     &offer,
		SKU:       &sku,
		Version:   &version,
	}
	return vm
}

// WithNetworkInterface adds a network interface
func (vm *VirtualMachine) WithNetworkInterface(id string, primary bool) *VirtualMachine {
	vm.Properties.NetworkProfile.NetworkInterfaces = append(
		vm.Properties.NetworkProfile.NetworkInterfaces,
		NetworkInterfaceReference{
			ID:      id,
			Primary: &primary,
		},
	)
	return vm
}
