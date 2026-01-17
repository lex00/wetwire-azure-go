package domain

import (
	"testing"

	coredomain "github.com/lex00/wetwire-core-go/domain"
)

// TestDomainInterface verifies that AzureDomain implements the Domain interface at compile time
func TestDomainInterface(t *testing.T) {
	var _ coredomain.Domain = (*AzureDomain)(nil)
}

// TestListerInterface verifies that AzureDomain implements the ListerDomain interface at compile time
func TestListerInterface(t *testing.T) {
	var _ coredomain.ListerDomain = (*AzureDomain)(nil)
}

// TestGrapherInterface verifies that AzureDomain implements the GrapherDomain interface at compile time
func TestGrapherInterface(t *testing.T) {
	var _ coredomain.GrapherDomain = (*AzureDomain)(nil)
}

// TestBuilderInterface verifies that azureBuilder implements the Builder interface at compile time
func TestBuilderInterface(t *testing.T) {
	var _ coredomain.Builder = (*azureBuilder)(nil)
}

// TestLinterInterface verifies that azureLinter implements the Linter interface at compile time
func TestLinterInterface(t *testing.T) {
	var _ coredomain.Linter = (*azureLinter)(nil)
}

// TestInitializerInterface verifies that azureInitializer implements the Initializer interface at compile time
func TestInitializerInterface(t *testing.T) {
	var _ coredomain.Initializer = (*azureInitializer)(nil)
}

// TestValidatorInterface verifies that azureValidator implements the Validator interface at compile time
func TestValidatorInterface(t *testing.T) {
	var _ coredomain.Validator = (*azureValidator)(nil)
}

// TestListerImplInterface verifies that azureLister implements the Lister interface at compile time
func TestListerImplInterface(t *testing.T) {
	var _ coredomain.Lister = (*azureLister)(nil)
}

// TestGrapherImplInterface verifies that azureGrapher implements the Grapher interface at compile time
func TestGrapherImplInterface(t *testing.T) {
	var _ coredomain.Grapher = (*azureGrapher)(nil)
}
