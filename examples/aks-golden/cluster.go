package aks_golden

import (
	"github.com/lex00/wetwire-azure-go/resources/aks"
)

// Helper variables for pointer fields (AKS uses pointers for optional fields)
var (
	// Cluster configuration
	kubernetesVersion = "1.29"
	dnsPrefix         = "aksgolden"
	nodeResourceGroup = "aks-golden-nodes-rg"

	// Network configuration
	networkPlugin   = "azure"
	networkPolicy   = "azure"
	serviceCidr     = "10.1.0.0/16"
	dnsServiceIP    = "10.1.0.10"
	loadBalancerSku = "standard"
	outboundType    = "loadBalancer"

	// Node pool configuration
	vmSizeSystem      = "Standard_D2s_v3"
	vmSizeApp         = "Standard_D4s_v3"
	vmSizeSpot        = "Standard_D4s_v3"
	osDiskSizeGB      = 128
	maxPods           = 110
	osType            = "Linux"
	osSKU             = "AzureLinux"
	poolTypeVMSS      = "VirtualMachineScaleSets"
	systemMode        = "System"
	userMode          = "User"
	spotPriority      = "Spot"
	spotEvictionDel   = "Delete"
	spotMaxPrice      = -1.0 // Pay up to on-demand price

	// Scaling configuration
	systemCount      = 2
	systemMin        = 2
	systemMax        = 4
	appCount         = 3
	appMin           = 2
	appMax           = 10
	spotCount        = 0
	spotMin          = 0
	spotMax          = 20
	autoScaleEnabled = true

	// SKU configuration
	skuTier = "Standard"

	// Identity configuration
	identityType = "SystemAssigned"

	// Security configuration
	rbacEnabled         = true
	aadManaged          = true
	enableAzureRBAC     = true
	workloadIdEnabled   = true
	oidcIssuerEnabled   = true
	imageCleanerEnabled = true
	imageCleanerHours   = 168 // Weekly cleanup

	// Monitoring configuration
	metricsEnabled = true

	// Upgrade configuration
	maxSurge = "33%"
)

// Cluster is the main AKS cluster resource.
var Cluster = aks.ManagedCluster{
	Name:       "aks-golden",
	Type:       "Microsoft.ContainerService/managedClusters",
	APIVersion: "2023-05-01",
	Location:   Location,
	Tags: map[string]string{
		"Environment": "production",
		"ManagedBy":   "wetwire",
	},
	Identity: &aks.ManagedClusterIdentity{
		Type: identityType,
	},
	SKU: &aks.ManagedClusterSKU{
		Tier: &skuTier,
	},
	Properties: aks.ManagedClusterProperties{
		KubernetesVersion: &kubernetesVersion,
		DNSPrefix:         &dnsPrefix,
		NodeResourceGroup: &nodeResourceGroup,
		EnableRBAC:        &rbacEnabled,
		AgentPoolProfiles: []aks.ManagedClusterAgentPoolProfile{
			SystemNodePool,
			AppNodePool,
			SpotNodePool,
		},
		NetworkProfile:   &ClusterNetworkProfile,
		AADProfile:       &ClusterAADProfile,
		SecurityProfile:  &ClusterSecurityProfile,
		AutoScalerProfile: &ClusterAutoScalerProfile,
		OIDCIssuerProfile: &aks.ManagedClusterOIDCIssuerProfile{
			Enabled: &oidcIssuerEnabled,
		},
		AzureMonitorProfile: &aks.ManagedClusterAzureMonitorProfile{
			Metrics: &aks.ManagedClusterAzureMonitorProfileMetrics{
				Enabled: &metricsEnabled,
			},
		},
	},
}

// ClusterNetworkProfile defines the network configuration for the cluster.
var ClusterNetworkProfile = aks.ContainerServiceNetworkProfile{
	NetworkPlugin:   &networkPlugin,
	NetworkPolicy:   &networkPolicy,
	ServiceCidr:     &serviceCidr,
	DNSServiceIP:    &dnsServiceIP,
	LoadBalancerSku: &loadBalancerSku,
	OutboundType:    &outboundType,
}

// ClusterAADProfile defines Azure AD integration.
var ClusterAADProfile = aks.ManagedClusterAADProfile{
	Managed:         &aadManaged,
	EnableAzureRBAC: &enableAzureRBAC,
	// AdminGroupObjectIDs would be set via parameters in production
	// AdminGroupObjectIDs: []string{"00000000-0000-0000-0000-000000000000"},
}

// ClusterSecurityProfile defines security settings.
var ClusterSecurityProfile = aks.ManagedClusterSecurityProfile{
	WorkloadIdentity: &aks.ManagedClusterSecurityProfileWorkloadIdentity{
		Enabled: &workloadIdEnabled,
	},
	ImageCleaner: &aks.ManagedClusterSecurityProfileImageCleaner{
		Enabled:       &imageCleanerEnabled,
		IntervalHours: &imageCleanerHours,
	},
}

// ClusterAutoScalerProfile defines cluster autoscaler settings.
var ClusterAutoScalerProfile = aks.ManagedClusterAutoScalerProfile{
	Expander:                  strPtr("least-waste"),
	ScanInterval:              strPtr("10s"),
	ScaleDownDelayAfterAdd:    strPtr("10m"),
	ScaleDownUnneededTime:     strPtr("10m"),
	ScaleDownUtilizationThreshold: strPtr("0.5"),
	MaxGracefulTerminationSec: strPtr("600"),
}

// SystemNodePool runs critical system workloads like CoreDNS.
var SystemNodePool = aks.ManagedClusterAgentPoolProfile{
	Name:              "system",
	VMSize:            &vmSizeSystem,
	Count:             &systemCount,
	MinCount:          &systemMin,
	MaxCount:          &systemMax,
	EnableAutoScaling: &autoScaleEnabled,
	OSDiskSizeGB:      &osDiskSizeGB,
	MaxPods:           &maxPods,
	OSType:            &osType,
	OSSKU:             &osSKU,
	Type:              &poolTypeVMSS,
	Mode:              &systemMode,
	AvailabilityZones: []string{"1", "2", "3"},
	NodeLabels: map[string]string{
		"nodepool-type": "system",
	},
	NodeTaints: []string{"CriticalAddonsOnly=true:NoSchedule"},
	UpgradeSettings: &aks.AgentPoolUpgradeSettings{
		MaxSurge: &maxSurge,
	},
	Tags: map[string]string{
		"Environment": "production",
		"NodePool":    "system",
	},
}

// AppNodePool runs general application workloads.
var AppNodePool = aks.ManagedClusterAgentPoolProfile{
	Name:              "app",
	VMSize:            &vmSizeApp,
	Count:             &appCount,
	MinCount:          &appMin,
	MaxCount:          &appMax,
	EnableAutoScaling: &autoScaleEnabled,
	OSDiskSizeGB:      &osDiskSizeGB,
	MaxPods:           &maxPods,
	OSType:            &osType,
	OSSKU:             &osSKU,
	Type:              &poolTypeVMSS,
	Mode:              &userMode,
	AvailabilityZones: []string{"1", "2", "3"},
	NodeLabels: map[string]string{
		"nodepool-type": "application",
	},
	UpgradeSettings: &aks.AgentPoolUpgradeSettings{
		MaxSurge: &maxSurge,
	},
	Tags: map[string]string{
		"Environment": "production",
		"NodePool":    "application",
	},
}

// SpotNodePool runs fault-tolerant workloads on spot instances.
var SpotNodePool = aks.ManagedClusterAgentPoolProfile{
	Name:                   "spot",
	VMSize:                 &vmSizeSpot,
	Count:                  &spotCount,
	MinCount:               &spotMin,
	MaxCount:               &spotMax,
	EnableAutoScaling:      &autoScaleEnabled,
	OSDiskSizeGB:           &osDiskSizeGB,
	MaxPods:                &maxPods,
	OSType:                 &osType,
	OSSKU:                  &osSKU,
	Type:                   &poolTypeVMSS,
	Mode:                   &userMode,
	AvailabilityZones:      []string{"1", "2", "3"},
	ScaleSetPriority:       &spotPriority,
	ScaleSetEvictionPolicy: &spotEvictionDel,
	SpotMaxPrice:           &spotMaxPrice,
	NodeLabels: map[string]string{
		"nodepool-type":                    "spot",
		"kubernetes.azure.com/scalesetpriority": "spot",
	},
	NodeTaints: []string{"kubernetes.azure.com/scalesetpriority=spot:NoSchedule"},
	UpgradeSettings: &aks.AgentPoolUpgradeSettings{
		MaxSurge: &maxSurge,
	},
	Tags: map[string]string{
		"Environment": "production",
		"NodePool":    "spot",
	},
}

// strPtr is a helper to create string pointers for autoscaler profile.
func strPtr(s string) *string {
	return &s
}
