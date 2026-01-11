package intrinsics

import (
	"testing"
)

func TestResourceID_ARMExpression(t *testing.T) {
	tests := []struct {
		name         string
		resourceID   ResourceID
		expected     string
	}{
		{
			name: "simple storage account",
			resourceID: ResourceID{
				ResourceType: "Microsoft.Storage/storageAccounts",
				ResourceName: "mystorageaccount",
			},
			expected: "[resourceId('Microsoft.Storage/storageAccounts', 'mystorageaccount')]",
		},
		{
			name: "virtual machine",
			resourceID: ResourceID{
				ResourceType: "Microsoft.Compute/virtualMachines",
				ResourceName: "myvm",
			},
			expected: "[resourceId('Microsoft.Compute/virtualMachines', 'myvm')]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.resourceID.ARMExpression()
			if result != tt.expected {
				t.Errorf("ARMExpression() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestResourceId_Constructor(t *testing.T) {
	rid := ResourceId("Microsoft.Storage/storageAccounts", "mystorage", "segment1", "segment2")

	if rid.ResourceType != "Microsoft.Storage/storageAccounts" {
		t.Errorf("ResourceType = %q, want %q", rid.ResourceType, "Microsoft.Storage/storageAccounts")
	}
	if rid.ResourceName != "mystorage" {
		t.Errorf("ResourceName = %q, want %q", rid.ResourceName, "mystorage")
	}
	if len(rid.Segments) != 2 {
		t.Errorf("len(Segments) = %d, want 2", len(rid.Segments))
	}
}

func TestReference_ARMExpression(t *testing.T) {
	tests := []struct {
		name      string
		reference Reference
		expected  string
	}{
		{
			name: "simple reference",
			reference: Reference{
				ResourceName: "myStorage",
				APIVersion:   "2021-02-01",
			},
			expected: "[reference('myStorage', '2021-02-01')]",
		},
		{
			name: "reference with property",
			reference: Reference{
				ResourceName: "myStorage",
				APIVersion:   "2021-02-01",
				Property:     "primaryEndpoints.blob",
			},
			expected: "[reference('myStorage', '2021-02-01').primaryEndpoints.blob]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.reference.ARMExpression()
			if result != tt.expected {
				t.Errorf("ARMExpression() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestRef_Constructor(t *testing.T) {
	ref := Ref("myResource", "2021-02-01")

	if ref.ResourceName != "myResource" {
		t.Errorf("ResourceName = %q, want %q", ref.ResourceName, "myResource")
	}
	if ref.APIVersion != "2021-02-01" {
		t.Errorf("APIVersion = %q, want %q", ref.APIVersion, "2021-02-01")
	}
	if ref.Property != "" {
		t.Errorf("Property = %q, want empty string", ref.Property)
	}
}

func TestRefProperty_Constructor(t *testing.T) {
	ref := RefProperty("myResource", "2021-02-01", "primaryEndpoints.blob")

	if ref.ResourceName != "myResource" {
		t.Errorf("ResourceName = %q, want %q", ref.ResourceName, "myResource")
	}
	if ref.APIVersion != "2021-02-01" {
		t.Errorf("APIVersion = %q, want %q", ref.APIVersion, "2021-02-01")
	}
	if ref.Property != "primaryEndpoints.blob" {
		t.Errorf("Property = %q, want %q", ref.Property, "primaryEndpoints.blob")
	}
}

func TestParameter_ARMExpression(t *testing.T) {
	param := Parameter{Name: "storageAccountName"}
	expected := "[parameters('storageAccountName')]"

	result := param.ARMExpression()
	if result != expected {
		t.Errorf("ARMExpression() = %q, want %q", result, expected)
	}
}

func TestParameters_Constructor(t *testing.T) {
	param := Parameters("myParam")

	if param.Name != "myParam" {
		t.Errorf("Name = %q, want %q", param.Name, "myParam")
	}

	// Verify interface compliance
	var _ Intrinsic = param
}

func TestVariable_ARMExpression(t *testing.T) {
	v := Variable{Name: "storageAccountId"}
	expected := "[variables('storageAccountId')]"

	result := v.ARMExpression()
	if result != expected {
		t.Errorf("ARMExpression() = %q, want %q", result, expected)
	}
}

func TestVariables_Constructor(t *testing.T) {
	v := Variables("myVar")

	if v.Name != "myVar" {
		t.Errorf("Name = %q, want %q", v.Name, "myVar")
	}

	// Verify interface compliance
	var _ Intrinsic = v
}

func TestConcat_ARMExpression(t *testing.T) {
	c := Concat{Values: []any{"a", "b", "c"}}
	expected := "[concat(...)]"

	result := c.ARMExpression()
	if result != expected {
		t.Errorf("ARMExpression() = %q, want %q", result, expected)
	}

	// Verify interface compliance
	var _ Intrinsic = c
}

func TestResourceGroupValue_ARMExpression(t *testing.T) {
	tests := []struct {
		name     string
		rg       ResourceGroupValue
		expected string
	}{
		{
			name:     "simple resourceGroup",
			rg:       ResourceGroupValue{},
			expected: "[resourceGroup()]",
		},
		{
			name:     "resourceGroup with property",
			rg:       ResourceGroupValue{Property: "location"},
			expected: "[resourceGroup().location]",
		},
		{
			name:     "resourceGroup with name property",
			rg:       ResourceGroupValue{Property: "name"},
			expected: "[resourceGroup().name]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rg.ARMExpression()
			if result != tt.expected {
				t.Errorf("ARMExpression() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestResourceGroup_Constructor(t *testing.T) {
	rg := ResourceGroup()

	if rg.Property != "" {
		t.Errorf("Property = %q, want empty string", rg.Property)
	}

	// Verify interface compliance
	var _ Intrinsic = rg
}

func TestSubscription_ARMExpression(t *testing.T) {
	tests := []struct {
		name     string
		sub      Subscription
		expected string
	}{
		{
			name:     "simple subscription",
			sub:      Subscription{},
			expected: "[subscription()]",
		},
		{
			name:     "subscription with property",
			sub:      Subscription{Property: "subscriptionId"},
			expected: "[subscription().subscriptionId]",
		},
		{
			name:     "subscription with tenantId",
			sub:      Subscription{Property: "tenantId"},
			expected: "[subscription().tenantId]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.sub.ARMExpression()
			if result != tt.expected {
				t.Errorf("ARMExpression() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestUniqueString_ARMExpression(t *testing.T) {
	u := UniqueString{Values: []string{"resourceGroup().id", "storageAccountName"}}
	expected := "[uniqueString(...)]"

	result := u.ARMExpression()
	if result != expected {
		t.Errorf("ARMExpression() = %q, want %q", result, expected)
	}

	// Verify interface compliance
	var _ Intrinsic = u
}

func TestIntrinsicInterface(t *testing.T) {
	// Verify all types implement Intrinsic interface
	intrinsics := []Intrinsic{
		ResourceID{},
		Reference{},
		Parameter{},
		Variable{},
		Concat{},
		ResourceGroupValue{},
		Subscription{},
		UniqueString{},
	}

	for i, intrinsic := range intrinsics {
		expr := intrinsic.ARMExpression()
		if expr == "" {
			t.Errorf("intrinsic %d returned empty expression", i)
		}
	}
}
