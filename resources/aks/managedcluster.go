// Package aks provides Azure Kubernetes Service (AKS) resource types
package aks

// ManagedCluster represents a Microsoft.ContainerService/managedClusters resource
type ManagedCluster struct {
	// Name is the name of the managed cluster
	Name string `json:"name"`

	// Type is the resource type
	Type string `json:"type"`

	// APIVersion is the API version to use for this resource
	APIVersion string `json:"apiVersion"`

	// Location is the Azure region where the cluster will be created
	Location string `json:"location"`

	// Tags are key-value pairs to organize resources
	Tags map[string]string `json:"tags,omitempty"`

	// Properties contains the properties of the managed cluster
	Properties ManagedClusterProperties `json:"properties"`

	// Identity defines the identity configuration for the cluster
	Identity *ManagedClusterIdentity `json:"identity,omitempty"`

	// SKU defines the SKU/tier for the cluster
	SKU *ManagedClusterSKU `json:"sku,omitempty"`
}

// ManagedClusterProperties represents the properties of a managed cluster
type ManagedClusterProperties struct {
	// KubernetesVersion is the Kubernetes version
	KubernetesVersion *string `json:"kubernetesVersion,omitempty"`

	// DNSPrefix is the DNS prefix for the cluster
	DNSPrefix *string `json:"dnsPrefix,omitempty"`

	// AgentPoolProfiles are the agent pool configurations
	AgentPoolProfiles []ManagedClusterAgentPoolProfile `json:"agentPoolProfiles,omitempty"`

	// LinuxProfile is the Linux VM configuration
	LinuxProfile *ContainerServiceLinuxProfile `json:"linuxProfile,omitempty"`

	// WindowsProfile is the Windows VM configuration
	WindowsProfile *ManagedClusterWindowsProfile `json:"windowsProfile,omitempty"`

	// ServicePrincipalProfile is the service principal configuration
	ServicePrincipalProfile *ManagedClusterServicePrincipalProfile `json:"servicePrincipalProfile,omitempty"`

	// AddonProfiles are the add-on configurations
	AddonProfiles map[string]ManagedClusterAddonProfile `json:"addonProfiles,omitempty"`

	// NodeResourceGroup is the name of the resource group for cluster nodes
	NodeResourceGroup *string `json:"nodeResourceGroup,omitempty"`

	// EnableRBAC indicates whether RBAC is enabled
	EnableRBAC *bool `json:"enableRBAC,omitempty"`

	// NetworkProfile is the network configuration
	NetworkProfile *ContainerServiceNetworkProfile `json:"networkProfile,omitempty"`

	// AADProfile is the Azure AD configuration
	AADProfile *ManagedClusterAADProfile `json:"aadProfile,omitempty"`

	// AutoScalerProfile is the cluster autoscaler configuration
	AutoScalerProfile *ManagedClusterAutoScalerProfile `json:"autoScalerProfile,omitempty"`

	// APIServerAccessProfile is the API server access configuration
	APIServerAccessProfile *ManagedClusterAPIServerAccessProfile `json:"apiServerAccessProfile,omitempty"`

	// DiskEncryptionSetID is the disk encryption set ID
	DiskEncryptionSetID *string `json:"diskEncryptionSetID,omitempty"`

	// IdentityProfile is the identity profile
	IdentityProfile map[string]UserAssignedIdentity `json:"identityProfile,omitempty"`

	// EnablePodSecurityPolicy indicates whether pod security policy is enabled (deprecated)
	EnablePodSecurityPolicy *bool `json:"enablePodSecurityPolicy,omitempty"`

	// HTTPProxyConfig is the HTTP proxy configuration
	HTTPProxyConfig *ManagedClusterHTTPProxyConfig `json:"httpProxyConfig,omitempty"`

	// OIDCIssuerProfile is the OIDC issuer configuration
	OIDCIssuerProfile *ManagedClusterOIDCIssuerProfile `json:"oidcIssuerProfile,omitempty"`

	// SecurityProfile is the security configuration
	SecurityProfile *ManagedClusterSecurityProfile `json:"securityProfile,omitempty"`

	// AzureMonitorProfile is the Azure Monitor configuration
	AzureMonitorProfile *ManagedClusterAzureMonitorProfile `json:"azureMonitorProfile,omitempty"`
}

// ManagedClusterAgentPoolProfile represents an agent pool configuration
type ManagedClusterAgentPoolProfile struct {
	// Name is the unique name of the agent pool
	Name string `json:"name"`

	// Count is the number of agents (VMs)
	Count *int `json:"count,omitempty"`

	// VMSize is the size of agent VMs
	VMSize *string `json:"vmSize,omitempty"`

	// OSDiskSizeGB is the OS disk size in GB
	OSDiskSizeGB *int `json:"osDiskSizeGB,omitempty"`

	// OSDiskType is the OS disk type (Managed, Ephemeral)
	OSDiskType *string `json:"osDiskType,omitempty"`

	// VnetSubnetID is the VNet subnet ID
	VnetSubnetID *string `json:"vnetSubnetID,omitempty"`

	// MaxPods is the maximum number of pods per node
	MaxPods *int `json:"maxPods,omitempty"`

	// OSType is the operating system type (Linux, Windows)
	OSType *string `json:"osType,omitempty"`

	// OSSKU is the OS SKU (Ubuntu, AzureLinux, Windows2019, Windows2022)
	OSSKU *string `json:"osSKU,omitempty"`

	// MaxCount is the maximum number of nodes for auto-scaling
	MaxCount *int `json:"maxCount,omitempty"`

	// MinCount is the minimum number of nodes for auto-scaling
	MinCount *int `json:"minCount,omitempty"`

	// EnableAutoScaling enables auto-scaling
	EnableAutoScaling *bool `json:"enableAutoScaling,omitempty"`

	// Type is the agent pool type (VirtualMachineScaleSets, AvailabilitySet)
	Type *string `json:"type,omitempty"`

	// Mode is the agent pool mode (System, User)
	Mode *string `json:"mode,omitempty"`

	// AvailabilityZones are the availability zones
	AvailabilityZones []string `json:"availabilityZones,omitempty"`

	// EnableNodePublicIP enables public IP on nodes
	EnableNodePublicIP *bool `json:"enableNodePublicIP,omitempty"`

	// NodeLabels are the labels for nodes
	NodeLabels map[string]string `json:"nodeLabels,omitempty"`

	// NodeTaints are the taints for nodes
	NodeTaints []string `json:"nodeTaints,omitempty"`

	// Tags are the tags for the agent pool
	Tags map[string]string `json:"tags,omitempty"`

	// UpgradeSettings are the upgrade settings
	UpgradeSettings *AgentPoolUpgradeSettings `json:"upgradeSettings,omitempty"`

	// ScaleSetPriority is the priority (Regular, Spot)
	ScaleSetPriority *string `json:"scaleSetPriority,omitempty"`

	// ScaleSetEvictionPolicy is the eviction policy for Spot VMs (Delete, Deallocate)
	ScaleSetEvictionPolicy *string `json:"scaleSetEvictionPolicy,omitempty"`

	// SpotMaxPrice is the max price for Spot VMs
	SpotMaxPrice *float64 `json:"spotMaxPrice,omitempty"`

	// EnableFIPS enables FIPS-compliant OS
	EnableFIPS *bool `json:"enableFIPS,omitempty"`

	// EnableEncryptionAtHost enables encryption at host
	EnableEncryptionAtHost *bool `json:"enableEncryptionAtHost,omitempty"`

	// KubeletConfig is the kubelet configuration
	KubeletConfig *KubeletConfig `json:"kubeletConfig,omitempty"`

	// LinuxOSConfig is the Linux OS configuration
	LinuxOSConfig *LinuxOSConfig `json:"linuxOSConfig,omitempty"`
}

// AgentPoolUpgradeSettings represents upgrade settings for an agent pool
type AgentPoolUpgradeSettings struct {
	// MaxSurge is the maximum number of extra nodes during upgrade
	MaxSurge *string `json:"maxSurge,omitempty"`
}

// KubeletConfig represents kubelet configuration
type KubeletConfig struct {
	// CPUManagerPolicy is the CPU manager policy (none, static)
	CPUManagerPolicy *string `json:"cpuManagerPolicy,omitempty"`

	// CPUCfsQuota enables CPU CFS quota
	CPUCfsQuota *bool `json:"cpuCfsQuota,omitempty"`

	// CPUCfsQuotaPeriod is the CPU CFS quota period
	CPUCfsQuotaPeriod *string `json:"cpuCfsQuotaPeriod,omitempty"`

	// ImageGcHighThreshold is the image GC high threshold
	ImageGcHighThreshold *int `json:"imageGcHighThreshold,omitempty"`

	// ImageGcLowThreshold is the image GC low threshold
	ImageGcLowThreshold *int `json:"imageGcLowThreshold,omitempty"`

	// TopologyManagerPolicy is the topology manager policy
	TopologyManagerPolicy *string `json:"topologyManagerPolicy,omitempty"`

	// AllowedUnsafeSysctls are the allowed unsafe sysctls
	AllowedUnsafeSysctls []string `json:"allowedUnsafeSysctls,omitempty"`

	// FailSwapOn disables kubelet if swap is on
	FailSwapOn *bool `json:"failSwapOn,omitempty"`

	// ContainerLogMaxSizeMB is the max container log size in MB
	ContainerLogMaxSizeMB *int `json:"containerLogMaxSizeMB,omitempty"`

	// ContainerLogMaxFiles is the max container log files
	ContainerLogMaxFiles *int `json:"containerLogMaxFiles,omitempty"`

	// PodMaxPids is the max PIDs per pod
	PodMaxPids *int `json:"podMaxPids,omitempty"`
}

// LinuxOSConfig represents Linux OS configuration
type LinuxOSConfig struct {
	// Sysctls are the sysctl settings
	Sysctls *SysctlConfig `json:"sysctls,omitempty"`

	// TransparentHugePageEnabled is the transparent huge page setting
	TransparentHugePageEnabled *string `json:"transparentHugePageEnabled,omitempty"`

	// TransparentHugePageDefrag is the transparent huge page defrag setting
	TransparentHugePageDefrag *string `json:"transparentHugePageDefrag,omitempty"`

	// SwapFileSizeMB is the swap file size in MB
	SwapFileSizeMB *int `json:"swapFileSizeMB,omitempty"`
}

// SysctlConfig represents sysctl configuration
type SysctlConfig struct {
	// NetCoreNetdevMaxBacklog is net.core.netdev_max_backlog
	NetCoreNetdevMaxBacklog *int `json:"netCoreNetdevMaxBacklog,omitempty"`

	// NetCoreRmemMax is net.core.rmem_max
	NetCoreRmemMax *int `json:"netCoreRmemMax,omitempty"`

	// NetCoreWmemMax is net.core.wmem_max
	NetCoreWmemMax *int `json:"netCoreWmemMax,omitempty"`

	// NetCoreSomaxconn is net.core.somaxconn
	NetCoreSomaxconn *int `json:"netCoreSomaxconn,omitempty"`

	// NetIpv4TcpMaxSynBacklog is net.ipv4.tcp_max_syn_backlog
	NetIpv4TcpMaxSynBacklog *int `json:"netIpv4TcpMaxSynBacklog,omitempty"`

	// NetIpv4TcpMaxTwBuckets is net.ipv4.tcp_max_tw_buckets
	NetIpv4TcpMaxTwBuckets *int `json:"netIpv4TcpMaxTwBuckets,omitempty"`

	// NetIpv4TcpFinTimeout is net.ipv4.tcp_fin_timeout
	NetIpv4TcpFinTimeout *int `json:"netIpv4TcpFinTimeout,omitempty"`

	// NetIpv4TcpKeepaliveTime is net.ipv4.tcp_keepalive_time
	NetIpv4TcpKeepaliveTime *int `json:"netIpv4TcpKeepaliveTime,omitempty"`

	// NetIpv4TcpKeepaliveProbes is net.ipv4.tcp_keepalive_probes
	NetIpv4TcpKeepaliveProbes *int `json:"netIpv4TcpKeepaliveProbes,omitempty"`

	// NetIpv4TcpKeepaliveIntvl is net.ipv4.tcp_keepalive_intvl
	NetIpv4TcpKeepaliveIntvl *int `json:"netIpv4TcpKeepaliveIntvl,omitempty"`

	// VMMaxMapCount is vm.max_map_count
	VMMaxMapCount *int `json:"vmMaxMapCount,omitempty"`

	// VMSwappiness is vm.swappiness
	VMSwappiness *int `json:"vmSwappiness,omitempty"`

	// VMVfsCachePressure is vm.vfs_cache_pressure
	VMVfsCachePressure *int `json:"vmVfsCachePressure,omitempty"`
}

// ContainerServiceLinuxProfile represents Linux profile configuration
type ContainerServiceLinuxProfile struct {
	// AdminUsername is the administrator username
	AdminUsername string `json:"adminUsername"`

	// SSH is the SSH configuration
	SSH ContainerServiceSshConfiguration `json:"ssh"`
}

// ContainerServiceSshConfiguration represents SSH configuration
type ContainerServiceSshConfiguration struct {
	// PublicKeys are the SSH public keys
	PublicKeys []ContainerServiceSshPublicKey `json:"publicKeys"`
}

// ContainerServiceSshPublicKey represents an SSH public key
type ContainerServiceSshPublicKey struct {
	// KeyData is the SSH public key data
	KeyData string `json:"keyData"`
}

// ManagedClusterWindowsProfile represents Windows profile configuration
type ManagedClusterWindowsProfile struct {
	// AdminUsername is the administrator username
	AdminUsername string `json:"adminUsername"`

	// AdminPassword is the administrator password
	AdminPassword *string `json:"adminPassword,omitempty"`

	// LicenseType is the license type (Windows_Server)
	LicenseType *string `json:"licenseType,omitempty"`

	// EnableCSIProxy enables CSI proxy
	EnableCSIProxy *bool `json:"enableCSIProxy,omitempty"`
}

// ManagedClusterServicePrincipalProfile represents service principal configuration
type ManagedClusterServicePrincipalProfile struct {
	// ClientID is the service principal client ID
	ClientID string `json:"clientId"`

	// Secret is the service principal secret
	Secret *string `json:"secret,omitempty"`
}

// ManagedClusterAddonProfile represents an add-on configuration
type ManagedClusterAddonProfile struct {
	// Enabled indicates whether the add-on is enabled
	Enabled bool `json:"enabled"`

	// Config is the add-on configuration
	Config map[string]string `json:"config,omitempty"`
}

// ContainerServiceNetworkProfile represents network configuration
type ContainerServiceNetworkProfile struct {
	// NetworkPlugin is the network plugin (azure, kubenet, none)
	NetworkPlugin *string `json:"networkPlugin,omitempty"`

	// NetworkPolicy is the network policy (azure, calico)
	NetworkPolicy *string `json:"networkPolicy,omitempty"`

	// NetworkMode is the network mode (transparent, bridge)
	NetworkMode *string `json:"networkMode,omitempty"`

	// PodCidr is the pod CIDR
	PodCidr *string `json:"podCidr,omitempty"`

	// ServiceCidr is the service CIDR
	ServiceCidr *string `json:"serviceCidr,omitempty"`

	// DNSServiceIP is the DNS service IP
	DNSServiceIP *string `json:"dnsServiceIP,omitempty"`

	// DockerBridgeCidr is the Docker bridge CIDR (deprecated)
	DockerBridgeCidr *string `json:"dockerBridgeCidr,omitempty"`

	// OutboundType is the outbound type (loadBalancer, userDefinedRouting, managedNATGateway, userAssignedNATGateway)
	OutboundType *string `json:"outboundType,omitempty"`

	// LoadBalancerSku is the load balancer SKU (standard, basic)
	LoadBalancerSku *string `json:"loadBalancerSku,omitempty"`

	// LoadBalancerProfile is the load balancer profile
	LoadBalancerProfile *ManagedClusterLoadBalancerProfile `json:"loadBalancerProfile,omitempty"`

	// NatGatewayProfile is the NAT gateway profile
	NatGatewayProfile *ManagedClusterNATGatewayProfile `json:"natGatewayProfile,omitempty"`

	// IPFamilies are the IP families (IPv4, IPv6)
	IPFamilies []string `json:"ipFamilies,omitempty"`
}

// ManagedClusterLoadBalancerProfile represents load balancer configuration
type ManagedClusterLoadBalancerProfile struct {
	// ManagedOutboundIPs is the managed outbound IPs configuration
	ManagedOutboundIPs *ManagedClusterLoadBalancerProfileManagedOutboundIPs `json:"managedOutboundIPs,omitempty"`

	// OutboundIPPrefixes is the outbound IP prefixes
	OutboundIPPrefixes *ManagedClusterLoadBalancerProfileOutboundIPPrefixes `json:"outboundIPPrefixes,omitempty"`

	// OutboundIPs is the outbound IPs
	OutboundIPs *ManagedClusterLoadBalancerProfileOutboundIPs `json:"outboundIPs,omitempty"`

	// AllocatedOutboundPorts is the number of allocated outbound ports
	AllocatedOutboundPorts *int `json:"allocatedOutboundPorts,omitempty"`

	// IdleTimeoutInMinutes is the idle timeout in minutes
	IdleTimeoutInMinutes *int `json:"idleTimeoutInMinutes,omitempty"`

	// EnableMultipleStandardLoadBalancers enables multiple standard load balancers
	EnableMultipleStandardLoadBalancers *bool `json:"enableMultipleStandardLoadBalancers,omitempty"`
}

// ManagedClusterLoadBalancerProfileManagedOutboundIPs represents managed outbound IPs
type ManagedClusterLoadBalancerProfileManagedOutboundIPs struct {
	// Count is the number of managed outbound IPs
	Count *int `json:"count,omitempty"`

	// CountIPv6 is the number of managed outbound IPv6 IPs
	CountIPv6 *int `json:"countIPv6,omitempty"`
}

// ManagedClusterLoadBalancerProfileOutboundIPPrefixes represents outbound IP prefixes
type ManagedClusterLoadBalancerProfileOutboundIPPrefixes struct {
	// PublicIPPrefixes are the public IP prefixes
	PublicIPPrefixes []ResourceReference `json:"publicIPPrefixes,omitempty"`
}

// ManagedClusterLoadBalancerProfileOutboundIPs represents outbound IPs
type ManagedClusterLoadBalancerProfileOutboundIPs struct {
	// PublicIPs are the public IPs
	PublicIPs []ResourceReference `json:"publicIPs,omitempty"`
}

// ManagedClusterNATGatewayProfile represents NAT gateway configuration
type ManagedClusterNATGatewayProfile struct {
	// ManagedOutboundIPProfile is the managed outbound IP profile
	ManagedOutboundIPProfile *ManagedClusterManagedOutboundIPProfile `json:"managedOutboundIPProfile,omitempty"`

	// IdleTimeoutInMinutes is the idle timeout in minutes
	IdleTimeoutInMinutes *int `json:"idleTimeoutInMinutes,omitempty"`
}

// ManagedClusterManagedOutboundIPProfile represents managed outbound IP profile
type ManagedClusterManagedOutboundIPProfile struct {
	// Count is the number of managed outbound IPs
	Count *int `json:"count,omitempty"`
}

// ResourceReference represents a reference to a resource
type ResourceReference struct {
	// ID is the resource ID
	ID *string `json:"id,omitempty"`
}

// ManagedClusterAADProfile represents AAD configuration
type ManagedClusterAADProfile struct {
	// Managed indicates whether AAD integration is managed
	Managed *bool `json:"managed,omitempty"`

	// EnableAzureRBAC enables Azure RBAC for K8s authorization
	EnableAzureRBAC *bool `json:"enableAzureRBAC,omitempty"`

	// AdminGroupObjectIDs are the AAD group object IDs for cluster admins
	AdminGroupObjectIDs []string `json:"adminGroupObjectIDs,omitempty"`

	// TenantID is the AAD tenant ID
	TenantID *string `json:"tenantID,omitempty"`

	// ClientAppID is the client application ID (non-managed)
	ClientAppID *string `json:"clientAppID,omitempty"`

	// ServerAppID is the server application ID (non-managed)
	ServerAppID *string `json:"serverAppID,omitempty"`

	// ServerAppSecret is the server application secret (non-managed)
	ServerAppSecret *string `json:"serverAppSecret,omitempty"`
}

// ManagedClusterAutoScalerProfile represents cluster autoscaler configuration
type ManagedClusterAutoScalerProfile struct {
	// BalanceSimilarNodeGroups balances similar node groups
	BalanceSimilarNodeGroups *string `json:"balance-similar-node-groups,omitempty"`

	// Expander is the expander type (least-waste, most-pods, priority, random)
	Expander *string `json:"expander,omitempty"`

	// MaxEmptyBulkDelete is max empty bulk delete
	MaxEmptyBulkDelete *string `json:"max-empty-bulk-delete,omitempty"`

	// MaxGracefulTerminationSec is max graceful termination seconds
	MaxGracefulTerminationSec *string `json:"max-graceful-termination-sec,omitempty"`

	// MaxNodeProvisionTime is max node provision time
	MaxNodeProvisionTime *string `json:"max-node-provision-time,omitempty"`

	// MaxTotalUnreadyPercentage is max total unready percentage
	MaxTotalUnreadyPercentage *string `json:"max-total-unready-percentage,omitempty"`

	// NewPodScaleUpDelay is new pod scale up delay
	NewPodScaleUpDelay *string `json:"new-pod-scale-up-delay,omitempty"`

	// OkTotalUnreadyCount is ok total unready count
	OkTotalUnreadyCount *string `json:"ok-total-unready-count,omitempty"`

	// ScaleDownDelayAfterAdd is scale down delay after add
	ScaleDownDelayAfterAdd *string `json:"scale-down-delay-after-add,omitempty"`

	// ScaleDownDelayAfterDelete is scale down delay after delete
	ScaleDownDelayAfterDelete *string `json:"scale-down-delay-after-delete,omitempty"`

	// ScaleDownDelayAfterFailure is scale down delay after failure
	ScaleDownDelayAfterFailure *string `json:"scale-down-delay-after-failure,omitempty"`

	// ScaleDownUnneededTime is scale down unneeded time
	ScaleDownUnneededTime *string `json:"scale-down-unneeded-time,omitempty"`

	// ScaleDownUnreadyTime is scale down unready time
	ScaleDownUnreadyTime *string `json:"scale-down-unready-time,omitempty"`

	// ScaleDownUtilizationThreshold is scale down utilization threshold
	ScaleDownUtilizationThreshold *string `json:"scale-down-utilization-threshold,omitempty"`

	// ScanInterval is scan interval
	ScanInterval *string `json:"scan-interval,omitempty"`

	// SkipNodesWithLocalStorage skips nodes with local storage
	SkipNodesWithLocalStorage *string `json:"skip-nodes-with-local-storage,omitempty"`

	// SkipNodesWithSystemPods skips nodes with system pods
	SkipNodesWithSystemPods *string `json:"skip-nodes-with-system-pods,omitempty"`
}

// ManagedClusterAPIServerAccessProfile represents API server access configuration
type ManagedClusterAPIServerAccessProfile struct {
	// AuthorizedIPRanges are the authorized IP ranges
	AuthorizedIPRanges []string `json:"authorizedIPRanges,omitempty"`

	// EnablePrivateCluster enables private cluster
	EnablePrivateCluster *bool `json:"enablePrivateCluster,omitempty"`

	// PrivateDNSZone is the private DNS zone
	PrivateDNSZone *string `json:"privateDNSZone,omitempty"`

	// EnablePrivateClusterPublicFQDN enables public FQDN for private cluster
	EnablePrivateClusterPublicFQDN *bool `json:"enablePrivateClusterPublicFQDN,omitempty"`

	// DisableRunCommand disables run command
	DisableRunCommand *bool `json:"disableRunCommand,omitempty"`
}

// ManagedClusterHTTPProxyConfig represents HTTP proxy configuration
type ManagedClusterHTTPProxyConfig struct {
	// HTTPProxy is the HTTP proxy URL
	HTTPProxy *string `json:"httpProxy,omitempty"`

	// HTTPSProxy is the HTTPS proxy URL
	HTTPSProxy *string `json:"httpsProxy,omitempty"`

	// NoProxy are the no-proxy addresses
	NoProxy []string `json:"noProxy,omitempty"`

	// TrustedCa is the trusted CA certificate
	TrustedCa *string `json:"trustedCa,omitempty"`
}

// ManagedClusterOIDCIssuerProfile represents OIDC issuer configuration
type ManagedClusterOIDCIssuerProfile struct {
	// Enabled enables OIDC issuer
	Enabled *bool `json:"enabled,omitempty"`
}

// ManagedClusterSecurityProfile represents security configuration
type ManagedClusterSecurityProfile struct {
	// Defender is the Microsoft Defender configuration
	Defender *ManagedClusterSecurityProfileDefender `json:"defender,omitempty"`

	// AzureKeyVaultKms is the Azure Key Vault KMS configuration
	AzureKeyVaultKms *AzureKeyVaultKms `json:"azureKeyVaultKms,omitempty"`

	// WorkloadIdentity is the workload identity configuration
	WorkloadIdentity *ManagedClusterSecurityProfileWorkloadIdentity `json:"workloadIdentity,omitempty"`

	// ImageCleaner is the image cleaner configuration
	ImageCleaner *ManagedClusterSecurityProfileImageCleaner `json:"imageCleaner,omitempty"`
}

// ManagedClusterSecurityProfileDefender represents Defender configuration
type ManagedClusterSecurityProfileDefender struct {
	// LogAnalyticsWorkspaceResourceId is the Log Analytics workspace resource ID
	LogAnalyticsWorkspaceResourceId *string `json:"logAnalyticsWorkspaceResourceId,omitempty"`

	// SecurityMonitoring is the security monitoring configuration
	SecurityMonitoring *ManagedClusterSecurityProfileDefenderSecurityMonitoring `json:"securityMonitoring,omitempty"`
}

// ManagedClusterSecurityProfileDefenderSecurityMonitoring represents security monitoring configuration
type ManagedClusterSecurityProfileDefenderSecurityMonitoring struct {
	// Enabled enables security monitoring
	Enabled *bool `json:"enabled,omitempty"`
}

// AzureKeyVaultKms represents Azure Key Vault KMS configuration
type AzureKeyVaultKms struct {
	// Enabled enables Azure Key Vault KMS
	Enabled *bool `json:"enabled,omitempty"`

	// KeyId is the Key Vault key ID
	KeyId *string `json:"keyId,omitempty"`

	// KeyVaultNetworkAccess is the key vault network access (Public, Private)
	KeyVaultNetworkAccess *string `json:"keyVaultNetworkAccess,omitempty"`

	// KeyVaultResourceId is the Key Vault resource ID
	KeyVaultResourceId *string `json:"keyVaultResourceId,omitempty"`
}

// ManagedClusterSecurityProfileWorkloadIdentity represents workload identity configuration
type ManagedClusterSecurityProfileWorkloadIdentity struct {
	// Enabled enables workload identity
	Enabled *bool `json:"enabled,omitempty"`
}

// ManagedClusterSecurityProfileImageCleaner represents image cleaner configuration
type ManagedClusterSecurityProfileImageCleaner struct {
	// Enabled enables image cleaner
	Enabled *bool `json:"enabled,omitempty"`

	// IntervalHours is the interval in hours
	IntervalHours *int `json:"intervalHours,omitempty"`
}

// ManagedClusterAzureMonitorProfile represents Azure Monitor configuration
type ManagedClusterAzureMonitorProfile struct {
	// Metrics is the metrics configuration
	Metrics *ManagedClusterAzureMonitorProfileMetrics `json:"metrics,omitempty"`
}

// ManagedClusterAzureMonitorProfileMetrics represents metrics configuration
type ManagedClusterAzureMonitorProfileMetrics struct {
	// Enabled enables metrics collection
	Enabled *bool `json:"enabled,omitempty"`

	// KubeStateMetrics is the kube-state-metrics configuration
	KubeStateMetrics *ManagedClusterAzureMonitorProfileKubeStateMetrics `json:"kubeStateMetrics,omitempty"`
}

// ManagedClusterAzureMonitorProfileKubeStateMetrics represents kube-state-metrics configuration
type ManagedClusterAzureMonitorProfileKubeStateMetrics struct {
	// MetricLabelsAllowlist is the metric labels allowlist
	MetricLabelsAllowlist *string `json:"metricLabelsAllowlist,omitempty"`

	// MetricAnnotationsAllowList is the metric annotations allowlist
	MetricAnnotationsAllowList *string `json:"metricAnnotationsAllowList,omitempty"`
}

// ManagedClusterIdentity represents identity configuration
type ManagedClusterIdentity struct {
	// Type is the identity type (SystemAssigned, UserAssigned, None)
	Type string `json:"type"`

	// UserAssignedIdentities are the user-assigned identities
	UserAssignedIdentities map[string]UserAssignedIdentity `json:"userAssignedIdentities,omitempty"`
}

// UserAssignedIdentity represents a user-assigned identity
type UserAssignedIdentity struct {
	// ClientID is the client ID
	ClientID *string `json:"clientId,omitempty"`

	// ObjectID is the object ID
	ObjectID *string `json:"objectId,omitempty"`

	// ResourceID is the resource ID
	ResourceID *string `json:"resourceId,omitempty"`
}

// ManagedClusterSKU represents SKU configuration
type ManagedClusterSKU struct {
	// Name is the SKU name (Base)
	Name *string `json:"name,omitempty"`

	// Tier is the SKU tier (Free, Standard, Premium)
	Tier *string `json:"tier,omitempty"`
}

// NewManagedCluster creates a new managed cluster with required fields
func NewManagedCluster(name, location, dnsPrefix string) *ManagedCluster {
	return &ManagedCluster{
		Name:       name,
		Type:       "Microsoft.ContainerService/managedClusters",
		APIVersion: "2023-05-01",
		Location:   location,
		Properties: ManagedClusterProperties{
			DNSPrefix: &dnsPrefix,
		},
	}
}

// WithTags adds tags to the cluster
func (m *ManagedCluster) WithTags(tags map[string]string) *ManagedCluster {
	m.Tags = tags
	return m
}

// WithKubernetesVersion sets the Kubernetes version
func (m *ManagedCluster) WithKubernetesVersion(version string) *ManagedCluster {
	m.Properties.KubernetesVersion = &version
	return m
}

// WithSystemAssignedIdentity configures system-assigned managed identity
func (m *ManagedCluster) WithSystemAssignedIdentity() *ManagedCluster {
	m.Identity = &ManagedClusterIdentity{
		Type: "SystemAssigned",
	}
	return m
}

// WithUserAssignedIdentity configures user-assigned managed identity
func (m *ManagedCluster) WithUserAssignedIdentity(identityID string) *ManagedCluster {
	m.Identity = &ManagedClusterIdentity{
		Type: "UserAssigned",
		UserAssignedIdentities: map[string]UserAssignedIdentity{
			identityID: {},
		},
	}
	return m
}

// WithAgentPool adds an agent pool profile
func (m *ManagedCluster) WithAgentPool(pool ManagedClusterAgentPoolProfile) *ManagedCluster {
	m.Properties.AgentPoolProfiles = append(m.Properties.AgentPoolProfiles, pool)
	return m
}

// WithNetworkProfile sets the network profile
func (m *ManagedCluster) WithNetworkProfile(profile ContainerServiceNetworkProfile) *ManagedCluster {
	m.Properties.NetworkProfile = &profile
	return m
}

// WithAADProfile sets the AAD profile
func (m *ManagedCluster) WithAADProfile(profile ManagedClusterAADProfile) *ManagedCluster {
	m.Properties.AADProfile = &profile
	return m
}

// WithPrivateCluster enables private cluster
func (m *ManagedCluster) WithPrivateCluster() *ManagedCluster {
	if m.Properties.APIServerAccessProfile == nil {
		m.Properties.APIServerAccessProfile = &ManagedClusterAPIServerAccessProfile{}
	}
	enabled := true
	m.Properties.APIServerAccessProfile.EnablePrivateCluster = &enabled
	return m
}

// WithRBAC enables RBAC
func (m *ManagedCluster) WithRBAC() *ManagedCluster {
	enabled := true
	m.Properties.EnableRBAC = &enabled
	return m
}

// WithAddon enables an add-on
func (m *ManagedCluster) WithAddon(name string, enabled bool, config map[string]string) *ManagedCluster {
	if m.Properties.AddonProfiles == nil {
		m.Properties.AddonProfiles = make(map[string]ManagedClusterAddonProfile)
	}
	m.Properties.AddonProfiles[name] = ManagedClusterAddonProfile{
		Enabled: enabled,
		Config:  config,
	}
	return m
}

// WithWorkloadIdentity enables workload identity
func (m *ManagedCluster) WithWorkloadIdentity() *ManagedCluster {
	if m.Properties.SecurityProfile == nil {
		m.Properties.SecurityProfile = &ManagedClusterSecurityProfile{}
	}
	enabled := true
	m.Properties.SecurityProfile.WorkloadIdentity = &ManagedClusterSecurityProfileWorkloadIdentity{
		Enabled: &enabled,
	}
	if m.Properties.OIDCIssuerProfile == nil {
		m.Properties.OIDCIssuerProfile = &ManagedClusterOIDCIssuerProfile{
			Enabled: &enabled,
		}
	}
	return m
}

// WithStandardTier sets the cluster to Standard tier
func (m *ManagedCluster) WithStandardTier() *ManagedCluster {
	tier := "Standard"
	m.SKU = &ManagedClusterSKU{
		Tier: &tier,
	}
	return m
}

// NewAgentPool creates a new agent pool profile
func NewAgentPool(name, vmSize string, count int) ManagedClusterAgentPoolProfile {
	return ManagedClusterAgentPoolProfile{
		Name:   name,
		VMSize: &vmSize,
		Count:  &count,
	}
}

// WithOSDiskSizeGB sets the OS disk size
func (p ManagedClusterAgentPoolProfile) WithOSDiskSizeGB(size int) ManagedClusterAgentPoolProfile {
	p.OSDiskSizeGB = &size
	return p
}

// WithAutoScaling enables auto-scaling
func (p ManagedClusterAgentPoolProfile) WithAutoScaling(min, max int) ManagedClusterAgentPoolProfile {
	enabled := true
	p.EnableAutoScaling = &enabled
	p.MinCount = &min
	p.MaxCount = &max
	return p
}

// WithMode sets the pool mode
func (p ManagedClusterAgentPoolProfile) WithMode(mode string) ManagedClusterAgentPoolProfile {
	p.Mode = &mode
	return p
}

// WithVnetSubnet sets the VNet subnet ID
func (p ManagedClusterAgentPoolProfile) WithVnetSubnet(subnetID string) ManagedClusterAgentPoolProfile {
	p.VnetSubnetID = &subnetID
	return p
}

// WithAvailabilityZones sets the availability zones
func (p ManagedClusterAgentPoolProfile) WithAvailabilityZones(zones ...string) ManagedClusterAgentPoolProfile {
	p.AvailabilityZones = zones
	return p
}

// AsSpot configures the pool as spot instances
func (p ManagedClusterAgentPoolProfile) AsSpot(evictionPolicy string, maxPrice float64) ManagedClusterAgentPoolProfile {
	priority := "Spot"
	p.ScaleSetPriority = &priority
	p.ScaleSetEvictionPolicy = &evictionPolicy
	p.SpotMaxPrice = &maxPrice
	return p
}
