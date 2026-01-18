package discover

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	coreast "github.com/lex00/wetwire-core-go/ast"
)

// DiscoveredResource represents a discovered Azure resource with metadata
type DiscoveredResource struct {
	Name         string   // Variable name
	Type         string   // Azure resource type (e.g., "Microsoft.Storage/storageAccounts")
	File         string   // Absolute path to the file
	Line         int      // Line number where the resource is declared
	Dependencies []string // Names of other resources this resource depends on
}

// azureResourceMap maps Go package paths to Azure resource types
var azureResourceMap = map[string]string{
	"storage.StorageAccount":      "Microsoft.Storage/storageAccounts",
	"compute.VirtualMachine":      "Microsoft.Compute/virtualMachines",
	"network.VirtualNetwork":      "Microsoft.Network/virtualNetworks",
	"network.NetworkInterface":    "Microsoft.Network/networkInterfaces",
	"network.Subnet":              "Microsoft.Network/subnets",
	"network.PublicIPAddress":     "Microsoft.Network/publicIPAddresses",
	"network.NetworkSecurityGroup": "Microsoft.Network/networkSecurityGroups",
	"keyvault.Vault":              "Microsoft.KeyVault/vaults",
	"sql.Server":                  "Microsoft.Sql/servers",
	"sql.Database":                "Microsoft.Sql/servers/databases",
	"web.Site":                    "Microsoft.Web/sites",
	"containerregistry.Registry":  "Microsoft.ContainerRegistry/registries",
	"aks.ManagedCluster":          "Microsoft.ContainerService/managedClusters",
}

// DiscoverResources discovers Azure resources in the given source directory
// by parsing Go AST and finding top-level variable declarations with Azure resource types.
func DiscoverResources(srcDir string) ([]DiscoveredResource, error) {
	var resources []DiscoveredResource

	// Walk through all Go files in the directory recursively
	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-Go files
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		// Parse the file
		fileResources, err := parseFile(path)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", path, err)
		}

		resources = append(resources, fileResources...)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return resources, nil
}

// parseFile parses a single Go file and extracts Azure resource declarations
func parseFile(filePath string) ([]DiscoveredResource, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var resources []DiscoveredResource
	packageImports := coreast.ExtractImports(node)

	// Visit all declarations in the file
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR {
			continue
		}

		// Process each variable specification
		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			// Process each name in the value spec
			for i, name := range valueSpec.Names {
				if name.Name == "_" {
					continue
				}

				// Check if this is an Azure resource type
				// First try the explicit type, then infer from the value
				var azureType string
				if valueSpec.Type != nil {
					azureType = getAzureResourceType(valueSpec.Type, packageImports)
				} else if i < len(valueSpec.Values) {
					azureType = inferAzureResourceType(valueSpec.Values[i], packageImports)
				}

				if azureType == "" {
					continue
				}

				// Extract dependencies from the value expression
				var dependencies []string
				if i < len(valueSpec.Values) {
					dependencies = extractDependencies(valueSpec.Values[i])
				}

				// Get the line number
				pos := fset.Position(name.Pos())

				resources = append(resources, DiscoveredResource{
					Name:         name.Name,
					Type:         azureType,
					File:         filePath,
					Line:         pos.Line,
					Dependencies: dependencies,
				})
			}
		}
	}

	return resources, nil
}


// inferAzureResourceType infers the Azure resource type from a value expression
// (e.g., from a composite literal like storage.StorageAccount{...})
func inferAzureResourceType(valueExpr ast.Expr, imports map[string]string) string {
	// Check if it's a composite literal
	if compLit, ok := valueExpr.(*ast.CompositeLit); ok {
		if compLit.Type != nil {
			return getAzureResourceType(compLit.Type, imports)
		}
	}
	return ""
}

// getAzureResourceType checks if the type expression represents an Azure resource
// and returns the Azure resource type string
func getAzureResourceType(typeExpr ast.Expr, imports map[string]string) string {
	// Use core AST utilities to extract type name and package
	typeName, pkgAlias := coreast.ExtractTypeName(typeExpr)
	if typeName == "" || pkgAlias == "" {
		return ""
	}

	key := fmt.Sprintf("%s.%s", pkgAlias, typeName)

	// Check if this is a known Azure resource type
	if azureType, ok := azureResourceMap[key]; ok {
		// Verify it's from the wetwire-azure-go package
		if importPath, exists := imports[pkgAlias]; exists {
			if strings.Contains(importPath, "wetwire-azure-go/resources") {
				return azureType
			}
		}
	}

	return ""
}

// extractDependencies extracts references to other variables from an expression
func extractDependencies(expr ast.Expr) []string {
	deps := make(map[string]bool)
	extractDependenciesRecursive(expr, deps)

	// Convert map to slice
	result := make([]string, 0, len(deps))
	for dep := range deps {
		result = append(result, dep)
	}
	return result
}

// extractDependenciesRecursive recursively extracts variable references from an expression
func extractDependenciesRecursive(expr ast.Expr, deps map[string]bool) {
	if expr == nil {
		return
	}

	switch e := expr.(type) {
	case *ast.Ident:
		// Direct variable reference
		if e.Name != "_" && !coreast.IsBuiltinIdent(e.Name) {
			deps[e.Name] = true
		}

	case *ast.UnaryExpr:
		// Handle &variable (address-of operator)
		if e.Op == token.AND {
			extractDependenciesRecursive(e.X, deps)
		}

	case *ast.CompositeLit:
		// Struct literal like storage.StorageAccount{...}
		for _, elt := range e.Elts {
			if kv, ok := elt.(*ast.KeyValueExpr); ok {
				extractDependenciesRecursive(kv.Value, deps)
			} else {
				extractDependenciesRecursive(elt, deps)
			}
		}

	case *ast.KeyValueExpr:
		// Field: value pairs
		extractDependenciesRecursive(e.Value, deps)

	case *ast.CallExpr:
		// Function calls
		for _, arg := range e.Args {
			extractDependenciesRecursive(arg, deps)
		}

	case *ast.ArrayType:
		// Array/slice type declarations
		extractDependenciesRecursive(e.Elt, deps)

	case *ast.SliceExpr:
		// Slice expressions
		extractDependenciesRecursive(e.X, deps)

	case *ast.IndexExpr:
		// Index expressions
		extractDependenciesRecursive(e.X, deps)

	case *ast.SelectorExpr:
		// Field/method selectors
		extractDependenciesRecursive(e.X, deps)

	case *ast.StarExpr:
		// Pointer types
		extractDependenciesRecursive(e.X, deps)

	case *ast.ParenExpr:
		// Parenthesized expressions
		extractDependenciesRecursive(e.X, deps)

	case *ast.BinaryExpr:
		// Binary operations
		extractDependenciesRecursive(e.X, deps)
		extractDependenciesRecursive(e.Y, deps)

	case *ast.MapType:
		// Map types
		extractDependenciesRecursive(e.Key, deps)
		extractDependenciesRecursive(e.Value, deps)
	}
}
