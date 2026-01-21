// Package v1 contains ASO Container Service resource types for Kubernetes-native Azure infrastructure management.
//
// These types enable managing AKS clusters using Kubernetes CRDs via Azure Service Operator (ASO).
//
// Example usage:
//
//	import (
//		aksv1 "github.com/lex00/wetwire-azure-go/resources/k8s/containerservice/v1"
//		metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//	)
//
//	var MyCluster = aksv1.ManagedCluster{
//		ObjectMeta: metav1.ObjectMeta{
//			Name:      "my-cluster",
//			Namespace: "aso-system",
//		},
//		Spec: aksv1.ManagedClusterSpec{
//			Location:          strPtr("eastus"),
//			Owner:             &aksv1.ResourceGroupReference{Name: "my-rg"},
//			KubernetesVersion: strPtr("1.29"),
//		},
//	}
package v1
