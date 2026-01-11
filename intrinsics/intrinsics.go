// Package intrinsics provides ARM template function wrappers for use in resource declarations.
package intrinsics

// Intrinsic represents an ARM template intrinsic function.
// When serialized, these become ARM template expressions like "[resourceId(...)]".
type Intrinsic interface {
	// ARMExpression returns the ARM template expression for this intrinsic.
	ARMExpression() string
}

// ResourceID represents the resourceId() ARM function.
type ResourceID struct {
	ResourceType string
	ResourceName string
	// Additional segments for nested resources
	Segments []string
}

// ARMExpression returns the ARM expression for resourceId.
func (r ResourceID) ARMExpression() string {
	return "[resourceId('" + r.ResourceType + "', '" + r.ResourceName + "')]"
}

// ResourceId creates a ResourceID intrinsic.
func ResourceId(resourceType, resourceName string, segments ...string) ResourceID {
	return ResourceID{
		ResourceType: resourceType,
		ResourceName: resourceName,
		Segments:     segments,
	}
}

// Reference represents the reference() ARM function.
type Reference struct {
	ResourceName string
	APIVersion   string
	Property     string
}

// ARMExpression returns the ARM expression for reference.
func (r Reference) ARMExpression() string {
	if r.Property != "" {
		return "[reference('" + r.ResourceName + "', '" + r.APIVersion + "')." + r.Property + "]"
	}
	return "[reference('" + r.ResourceName + "', '" + r.APIVersion + "')]"
}

// Ref creates a Reference intrinsic for referencing another resource.
func Ref(resourceName, apiVersion string) Reference {
	return Reference{
		ResourceName: resourceName,
		APIVersion:   apiVersion,
	}
}

// RefProperty creates a Reference intrinsic for a specific property.
func RefProperty(resourceName, apiVersion, property string) Reference {
	return Reference{
		ResourceName: resourceName,
		APIVersion:   apiVersion,
		Property:     property,
	}
}

// Parameter represents the parameters() ARM function.
type Parameter struct {
	Name string
}

// ARMExpression returns the ARM expression for parameters.
func (p Parameter) ARMExpression() string {
	return "[parameters('" + p.Name + "')]"
}

// Parameters creates a Parameter intrinsic.
func Parameters(name string) Parameter {
	return Parameter{Name: name}
}

// Variable represents the variables() ARM function.
type Variable struct {
	Name string
}

// ARMExpression returns the ARM expression for variables.
func (v Variable) ARMExpression() string {
	return "[variables('" + v.Name + "')]"
}

// Variables creates a Variable intrinsic.
func Variables(name string) Variable {
	return Variable{Name: name}
}

// Concat represents the concat() ARM function.
type Concat struct {
	Values []any
}

// ARMExpression returns the ARM expression for concat.
func (c Concat) ARMExpression() string {
	return "[concat(...)]" // Simplified for now
}

// ResourceGroupValue represents resourceGroup() ARM function.
type ResourceGroupValue struct {
	Property string
}

// ARMExpression returns the ARM expression for resourceGroup.
func (r ResourceGroupValue) ARMExpression() string {
	if r.Property != "" {
		return "[resourceGroup()." + r.Property + "]"
	}
	return "[resourceGroup()]"
}

// ResourceGroup creates a ResourceGroupValue intrinsic.
func ResourceGroup() ResourceGroupValue {
	return ResourceGroupValue{}
}

// Subscription represents subscription() ARM function.
type Subscription struct {
	Property string
}

// ARMExpression returns the ARM expression for subscription.
func (s Subscription) ARMExpression() string {
	if s.Property != "" {
		return "[subscription()." + s.Property + "]"
	}
	return "[subscription()]"
}

// UniqueString represents uniqueString() ARM function.
type UniqueString struct {
	Values []string
}

// ARMExpression returns the ARM expression for uniqueString.
func (u UniqueString) ARMExpression() string {
	return "[uniqueString(...)]" // Simplified for now
}
