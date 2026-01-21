// Package aks_k8s provides a production-ready AKS cluster example using ASO (Azure Service Operator).
//
// This example demonstrates the K8s-native approach to Azure infrastructure,
// where Azure resources are managed as Kubernetes CRDs via Azure Service Operator.
package aks_k8s

import (
	networkv1 "github.com/lex00/wetwire-azure-go/resources/k8s/network/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Helper variables for pointer fields
var (
	location              = "eastus"
	resourceGroupName     = "aks-k8s-rg"
	vnetName              = "aks-k8s-vnet"
	aksSubnetName         = "aks-subnet"
	appgwSubnetName       = "appgw-subnet"
	privateEndpointSubnet = "private-endpoint-subnet"
	vnetAddressPrefix     = "10.0.0.0/16"
	aksSubnetPrefix       = "10.0.0.0/22"
	appgwSubnetPrefix     = "10.0.4.0/24"
	peSubnetPrefix        = "10.0.5.0/24"
	pePoliciesDisabled    = "Disabled"
)

// VNet is the virtual network for the AKS cluster, managed via ASO.
var VNet = networkv1.VirtualNetwork{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "network.azure.com/v1",
		Kind:       "VirtualNetwork",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "aks-k8s-vnet",
		Namespace: "aso-system",
	},
	Spec: networkv1.VirtualNetworkSpec{
		Owner: &networkv1.ResourceGroupReference{
			Name: resourceGroupName,
		},
		AzureName: &vnetName,
		Location:  &location,
		AddressSpace: &networkv1.AddressSpace{
			AddressPrefixes: []string{vnetAddressPrefix},
		},
		Tags: map[string]string{
			"Environment": "production",
			"ManagedBy":   "wetwire-aso",
		},
	},
}

// AKSSubnet is the subnet for AKS nodes.
// Using a /22 CIDR provides ~1000 IPs for nodes and pods with Azure CNI.
var AKSSubnet = networkv1.Subnet{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "network.azure.com/v1",
		Kind:       "VirtualNetworksSubnet",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "aks-k8s-subnet",
		Namespace: "aso-system",
	},
	Spec: networkv1.SubnetSpec{
		Owner: &networkv1.VirtualNetworkReference{
			Name: vnetName,
		},
		AzureName:     &aksSubnetName,
		AddressPrefix: &aksSubnetPrefix,
		ServiceEndpoints: []networkv1.ServiceEndpointPropertiesFormat{
			{Service: strPtr("Microsoft.Storage"), Locations: []string{location}},
			{Service: strPtr("Microsoft.KeyVault"), Locations: []string{location}},
			{Service: strPtr("Microsoft.ContainerRegistry"), Locations: []string{location}},
		},
	},
}

// AppGatewaySubnet is the subnet for Azure Application Gateway (AGIC).
var AppGatewaySubnet = networkv1.Subnet{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "network.azure.com/v1",
		Kind:       "VirtualNetworksSubnet",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "aks-k8s-appgw-subnet",
		Namespace: "aso-system",
	},
	Spec: networkv1.SubnetSpec{
		Owner: &networkv1.VirtualNetworkReference{
			Name: vnetName,
		},
		AzureName:     &appgwSubnetName,
		AddressPrefix: &appgwSubnetPrefix,
	},
}

// PrivateEndpointSubnet is the subnet for private endpoints.
var PrivateEndpointSubnet = networkv1.Subnet{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "network.azure.com/v1",
		Kind:       "VirtualNetworksSubnet",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "aks-k8s-pe-subnet",
		Namespace: "aso-system",
	},
	Spec: networkv1.SubnetSpec{
		Owner: &networkv1.VirtualNetworkReference{
			Name: vnetName,
		},
		AzureName:                      &privateEndpointSubnet,
		AddressPrefix:                  &peSubnetPrefix,
		PrivateEndpointNetworkPolicies: &pePoliciesDisabled,
	},
}

// NSG is the network security group for AKS nodes.
var NSG = networkv1.NetworkSecurityGroup{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "network.azure.com/v1",
		Kind:       "NetworkSecurityGroup",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "aks-k8s-nsg",
		Namespace: "aso-system",
	},
	Spec: networkv1.NetworkSecurityGroupSpec{
		Owner: &networkv1.ResourceGroupReference{
			Name: resourceGroupName,
		},
		AzureName: strPtr("aks-k8s-nsg"),
		Location:  &location,
		SecurityRules: []networkv1.SecurityRule{
			{
				Name:                     strPtr("AllowHTTPS"),
				Priority:                 intPtr(100),
				Direction:                strPtr("Inbound"),
				Access:                   strPtr("Allow"),
				Protocol:                 strPtr("Tcp"),
				SourcePortRange:          strPtr("*"),
				DestinationPortRange:     strPtr("443"),
				SourceAddressPrefix:      strPtr("Internet"),
				DestinationAddressPrefix: strPtr("*"),
			},
			{
				Name:                     strPtr("AllowHTTP"),
				Priority:                 intPtr(101),
				Direction:                strPtr("Inbound"),
				Access:                   strPtr("Allow"),
				Protocol:                 strPtr("Tcp"),
				SourcePortRange:          strPtr("*"),
				DestinationPortRange:     strPtr("80"),
				SourceAddressPrefix:      strPtr("Internet"),
				DestinationAddressPrefix: strPtr("*"),
			},
			{
				Name:                     strPtr("DenyAllInbound"),
				Priority:                 intPtr(4096),
				Direction:                strPtr("Inbound"),
				Access:                   strPtr("Deny"),
				Protocol:                 strPtr("*"),
				SourcePortRange:          strPtr("*"),
				DestinationPortRange:     strPtr("*"),
				SourceAddressPrefix:      strPtr("*"),
				DestinationAddressPrefix: strPtr("*"),
			},
		},
		Tags: map[string]string{
			"Environment": "production",
			"ManagedBy":   "wetwire-aso",
		},
	},
}

// strPtr is a helper to create string pointers.
func strPtr(s string) *string {
	return &s
}

// intPtr is a helper to create int pointers.
func intPtr(i int) *int {
	return &i
}
