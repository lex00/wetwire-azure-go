// Package v1 contains ASO Managed Identity resource types for Kubernetes-native Azure infrastructure management.
//
// These types enable managing User Assigned Identities and Federated Identity Credentials
// using Kubernetes CRDs via Azure Service Operator (ASO).
//
// Example usage:
//
//	import (
//		identityv1 "github.com/lex00/wetwire-azure-go/resources/k8s/managedidentity/v1"
//		metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//	)
//
//	var MyIdentity = identityv1.UserAssignedIdentity{
//		ObjectMeta: metav1.ObjectMeta{
//			Name:      "my-identity",
//			Namespace: "aso-system",
//		},
//		Spec: identityv1.UserAssignedIdentitySpec{
//			Location: strPtr("eastus"),
//			Owner:    &identityv1.ResourceGroupReference{Name: "my-rg"},
//		},
//	}
package v1
