// Package v1 contains ASO Network resource types for Kubernetes-native Azure infrastructure management.
//
// These types enable managing Virtual Networks, Subnets, and Network Security Groups
// using Kubernetes CRDs via Azure Service Operator (ASO).
//
// Example usage:
//
//	import (
//		networkv1 "github.com/lex00/wetwire-azure-go/resources/k8s/network/v1"
//		metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//	)
//
//	var MyVNet = networkv1.VirtualNetwork{
//		ObjectMeta: metav1.ObjectMeta{
//			Name:      "my-vnet",
//			Namespace: "aso-system",
//		},
//		Spec: networkv1.VirtualNetworkSpec{
//			Location: strPtr("eastus"),
//			Owner:    &networkv1.ResourceGroupReference{Name: "my-rg"},
//			AddressSpace: &networkv1.AddressSpace{
//				AddressPrefixes: []string{"10.0.0.0/16"},
//			},
//		},
//	}
package v1
