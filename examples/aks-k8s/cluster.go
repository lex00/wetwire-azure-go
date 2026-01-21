package aks_k8s

import (
	aksv1 "github.com/lex00/wetwire-azure-go/resources/k8s/containerservice/v1"
	identityv1 "github.com/lex00/wetwire-azure-go/resources/k8s/managedidentity/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Helper variables for pointer fields
var (
	clusterName       = "aks-k8s-cluster"
	kubernetesVersion = "1.29"
	dnsPrefix         = "aksk8s"
	nodeResourceGroup = "aks-k8s-nodes-rg"

	// Network configuration
	networkPlugin   = "azure"
	networkPolicy   = "azure"
	serviceCidr     = "10.1.0.0/16"
	dnsServiceIP    = "10.1.0.10"
	loadBalancerSku = "standard"
	outboundType    = "loadBalancer"

	// Node pool configuration
	vmSizeSystem = "Standard_D2s_v3"
	vmSizeApp    = "Standard_D4s_v3"
	vmSizeSpot   = "Standard_D4s_v3"
	osDiskSizeGB = 128
	maxPods      = 110
	osType       = "Linux"
	osSKU        = "AzureLinux"
	poolTypeVMSS = "VirtualMachineScaleSets"
	systemMode   = "System"
	userMode     = "User"
	spotPriority = "Spot"
	spotEviction = "Delete"
	spotMaxPrice = -1.0

	// Scaling configuration
	systemCount = 2
	systemMin   = 2
	systemMax   = 4
	appCount    = 3
	appMin      = 2
	appMax      = 10
	spotCount   = 0
	spotMin     = 0
	spotMax     = 20

	// Identity and security
	identityType      = "SystemAssigned"
	skuTier           = "Standard"
	rbacEnabled       = true
	aadManaged        = true
	enableAzureRBAC   = true
	workloadIdEnabled = true
	oidcEnabled       = true
	imageCleanEnabled = true
	imageCleanHours   = 168

	// Autoscaler configuration
	autoscalerExpander     = "least-waste"
	autoscalerScanInterval = "10s"
	scaleDownDelay         = "10m"
	scaleDownUnneeded      = "10m"
	scaleDownThreshold     = "0.5"

	// Upgrade configuration
	maxSurge = "33%"
)

// ClusterIdentity is the managed identity for the AKS cluster.
var ClusterIdentity = identityv1.UserAssignedIdentity{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "managedidentity.azure.com/v1",
		Kind:       "UserAssignedIdentity",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "aks-k8s-cluster-identity",
		Namespace: "aso-system",
	},
	Spec: identityv1.UserAssignedIdentitySpec{
		Owner: &identityv1.ResourceGroupReference{
			Name: resourceGroupName,
		},
		AzureName: strPtr("aks-k8s-cluster-identity"),
		Location:  &location,
		Tags: map[string]string{
			"Environment": "production",
			"ManagedBy":   "wetwire-aso",
		},
	},
}

// Cluster is the main AKS cluster resource, managed via ASO.
var Cluster = aksv1.ManagedCluster{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "containerservice.azure.com/v1",
		Kind:       "ManagedCluster",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "aks-k8s-cluster",
		Namespace: "aso-system",
	},
	Spec: aksv1.ManagedClusterSpec{
		Owner: &aksv1.ResourceGroupReference{
			Name: resourceGroupName,
		},
		AzureName:         &clusterName,
		Location:          &location,
		KubernetesVersion: &kubernetesVersion,
		DNSPrefix:         &dnsPrefix,
		NodeResourceGroup: &nodeResourceGroup,
		EnableRBAC:        &rbacEnabled,
		Identity: &aksv1.ManagedClusterIdentity{
			Type: &identityType,
		},
		SKU: &aksv1.ManagedClusterSKU{
			Tier: &skuTier,
		},
		NetworkProfile: &aksv1.ContainerServiceNetworkProfile{
			NetworkPlugin:   &networkPlugin,
			NetworkPolicy:   &networkPolicy,
			ServiceCidr:     &serviceCidr,
			DNSServiceIP:    &dnsServiceIP,
			LoadBalancerSku: &loadBalancerSku,
			OutboundType:    &outboundType,
		},
		AADProfile: &aksv1.ManagedClusterAADProfile{
			Managed:         &aadManaged,
			EnableAzureRBAC: &enableAzureRBAC,
		},
		SecurityProfile: &aksv1.ManagedClusterSecurityProfile{
			WorkloadIdentity: &aksv1.ManagedClusterSecurityProfileWorkloadIdentity{
				Enabled: &workloadIdEnabled,
			},
			ImageCleaner: &aksv1.ManagedClusterSecurityProfileImageCleaner{
				Enabled:       &imageCleanEnabled,
				IntervalHours: &imageCleanHours,
			},
		},
		OIDCIssuerProfile: &aksv1.ManagedClusterOIDCIssuerProfile{
			Enabled: &oidcEnabled,
		},
		AutoScalerProfile: &aksv1.ManagedClusterAutoScalerProfile{
			Expander:                      &autoscalerExpander,
			ScanInterval:                  &autoscalerScanInterval,
			ScaleDownDelayAfterAdd:        &scaleDownDelay,
			ScaleDownUnneededTime:         &scaleDownUnneeded,
			ScaleDownUtilizationThreshold: &scaleDownThreshold,
		},
		AgentPoolProfiles: []aksv1.ManagedClusterAgentPoolProfile{
			SystemNodePool,
			AppNodePool,
			SpotNodePool,
		},
		Tags: map[string]string{
			"Environment": "production",
			"ManagedBy":   "wetwire-aso",
		},
	},
}

// SystemNodePool runs critical system workloads like CoreDNS.
var SystemNodePool = aksv1.ManagedClusterAgentPoolProfile{
	Name:              &systemPoolName,
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
	VnetSubnetReference: &aksv1.SubnetReference{
		Name: &aksSubnetName,
	},
	NodeLabels: map[string]string{
		"nodepool-type": "system",
	},
	NodeTaints: []string{"CriticalAddonsOnly=true:NoSchedule"},
	Tags: map[string]string{
		"Environment": "production",
		"NodePool":    "system",
	},
}

// AppNodePool runs general application workloads.
var AppNodePool = aksv1.ManagedClusterAgentPoolProfile{
	Name:              &appPoolName,
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
	VnetSubnetReference: &aksv1.SubnetReference{
		Name: &aksSubnetName,
	},
	NodeLabels: map[string]string{
		"nodepool-type": "application",
	},
	Tags: map[string]string{
		"Environment": "production",
		"NodePool":    "application",
	},
}

// SpotNodePool runs fault-tolerant workloads on spot instances.
var SpotNodePool = aksv1.ManagedClusterAgentPoolProfile{
	Name:                   &spotPoolName,
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
	ScaleSetEvictionPolicy: &spotEviction,
	SpotMaxPrice:           &spotMaxPrice,
	VnetSubnetReference: &aksv1.SubnetReference{
		Name: &aksSubnetName,
	},
	NodeLabels: map[string]string{
		"nodepool-type":                          "spot",
		"kubernetes.azure.com/scalesetpriority": "spot",
	},
	NodeTaints: []string{"kubernetes.azure.com/scalesetpriority=spot:NoSchedule"},
	Tags: map[string]string{
		"Environment": "production",
		"NodePool":    "spot",
	},
}

// Additional helper variables for node pool names
var (
	systemPoolName   = "system"
	appPoolName      = "app"
	spotPoolName     = "spot"
	autoScaleEnabled = true
)
