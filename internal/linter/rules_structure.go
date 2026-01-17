package linter

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

// WAZ001 checks for use of typed constants instead of string literals for locations
type WAZ001 struct{}

func (r *WAZ001) ID() string {
	return "WAZ001"
}

func (r *WAZ001) Description() string {
	return "Use location constants instead of string literals for locations"
}

func (r *WAZ001) Severity() Severity {
	return SeverityWarning
}

func (r *WAZ001) Check(file string) ([]LintResult, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var results []LintResult

	ast.Inspect(node, func(n ast.Node) bool {
		// Look for key-value pairs in struct literals
		kv, ok := n.(*ast.KeyValueExpr)
		if !ok {
			return true
		}

		// Check if the key is "Location"
		if ident, ok := kv.Key.(*ast.Ident); ok && ident.Name == "Location" {
			// Check if the value is a string literal
			if lit, ok := kv.Value.(*ast.BasicLit); ok && lit.Kind == token.STRING {
				locationValue := strings.Trim(lit.Value, `"`)
				// Skip ARM template expressions (e.g., "[resourceGroup().location]")
				if strings.HasPrefix(locationValue, "[") && strings.HasSuffix(locationValue, "]") {
					return true
				}
				// Check for invalid location formats (spaces, uppercase)
				if strings.Contains(locationValue, " ") || hasUpperCase(locationValue) {
					pos := fset.Position(lit.Pos())
					results = append(results, LintResult{
						Rule:     r.ID(),
						File:     file,
						Line:     pos.Line,
						Message:  fmt.Sprintf("Location '%s' should use lowercase format without spaces (e.g., 'eastus' not 'East US')", locationValue),
						Severity: r.Severity(),
					})
				}
			}
		}

		return true
	})

	return results, nil
}

// CanFix returns true since WAZ001 supports auto-fixing
func (r *WAZ001) CanFix() bool {
	return true
}

// Fix applies auto-fix to normalize location string literals
func (r *WAZ001) Fix(file string) (string, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return "", err
	}

	// Read the original file to preserve formatting
	content, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}

	result := string(content)

	ast.Inspect(node, func(n ast.Node) bool {
		// Look for key-value pairs in struct literals
		kv, ok := n.(*ast.KeyValueExpr)
		if !ok {
			return true
		}

		// Check if the key is "Location"
		if ident, ok := kv.Key.(*ast.Ident); ok && ident.Name == "Location" {
			// Check if the value is a string literal
			if lit, ok := kv.Value.(*ast.BasicLit); ok && lit.Kind == token.STRING {
				locationValue := strings.Trim(lit.Value, `"`)
				// Skip ARM template expressions (e.g., "[resourceGroup().location]")
				if strings.HasPrefix(locationValue, "[") && strings.HasSuffix(locationValue, "]") {
					return true
				}
				// Check for invalid location formats (spaces, uppercase)
				if strings.Contains(locationValue, " ") || hasUpperCase(locationValue) {
					// Normalize location: lowercase and remove spaces
					normalized := strings.ToLower(strings.ReplaceAll(locationValue, " ", ""))
					// Replace in the result string
					oldLit := fmt.Sprintf(`"%s"`, locationValue)
					newLit := fmt.Sprintf(`"%s"`, normalized)
					result = strings.Replace(result, oldLit, newLit, 1)
				}
			}
		}

		return true
	})

	return result, nil
}

// WAZ002 checks for use of direct references instead of explicit resourceId() calls
type WAZ002 struct{}

func (r *WAZ002) ID() string {
	return "WAZ002"
}

func (r *WAZ002) Description() string {
	return "Use direct references instead of explicit resourceId() calls"
}

func (r *WAZ002) Severity() Severity {
	return SeverityWarning
}

func (r *WAZ002) Check(file string) ([]LintResult, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var results []LintResult

	ast.Inspect(node, func(n ast.Node) bool {
		// Look for function calls named "resourceId"
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if ident, ok := call.Fun.(*ast.Ident); ok && ident.Name == "resourceId" {
			pos := fset.Position(call.Pos())
			results = append(results, LintResult{
				Rule:     r.ID(),
				File:     file,
				Line:     pos.Line,
				Message:  "Use direct resource references (e.g., MyResource.Id) instead of resourceId() function",
				Severity: r.Severity(),
			})
		}

		return true
	})

	return results, nil
}

// WAZ003 checks for extracting nested configurations to named variables
type WAZ003 struct{}

func (r *WAZ003) ID() string {
	return "WAZ003"
}

func (r *WAZ003) Description() string {
	return "Extract nested configurations to named variables"
}

func (r *WAZ003) Severity() Severity {
	return SeverityWarning
}

func (r *WAZ003) Check(file string) ([]LintResult, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var results []LintResult

	ast.Inspect(node, func(n ast.Node) bool {
		// Look for composite literals (struct literals)
		comp, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}

		// Check nesting depth of composite literals
		depth := r.getNestingDepth(comp)
		if depth >= 3 {
			pos := fset.Position(comp.Pos())
			results = append(results, LintResult{
				Rule:     r.ID(),
				File:     file,
				Line:     pos.Line,
				Message:  "Deeply nested configuration detected. Consider extracting to named variables for better readability",
				Severity: r.Severity(),
			})
		}

		return true
	})

	return results, nil
}

func (r *WAZ003) getNestingDepth(expr ast.Expr) int {
	comp, ok := expr.(*ast.CompositeLit)
	if !ok {
		return 0
	}

	maxDepth := 0
	for _, elt := range comp.Elts {
		var depth int
		if kv, ok := elt.(*ast.KeyValueExpr); ok {
			// Check if value is a composite literal or unary expr with composite literal
			if compLit, ok := kv.Value.(*ast.CompositeLit); ok {
				depth = 1 + r.getNestingDepth(compLit)
			} else if unary, ok := kv.Value.(*ast.UnaryExpr); ok {
				if compLit, ok := unary.X.(*ast.CompositeLit); ok {
					depth = 1 + r.getNestingDepth(compLit)
				}
			}
		} else if compLit, ok := elt.(*ast.CompositeLit); ok {
			depth = 1 + r.getNestingDepth(compLit)
		}

		if depth > maxDepth {
			maxDepth = depth
		}
	}

	return maxDepth
}

// WAZ004 checks for duplicate resource names
type WAZ004 struct{}

func (r *WAZ004) ID() string {
	return "WAZ004"
}

func (r *WAZ004) Description() string {
	return "Flag duplicate resource names"
}

func (r *WAZ004) Severity() Severity {
	return SeverityError
}

func (r *WAZ004) Check(file string) ([]LintResult, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var results []LintResult
	varNames := make(map[string]token.Pos)

	// Visit all top-level variable declarations
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR {
			continue
		}

		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			for _, name := range valueSpec.Names {
				if name.Name == "_" {
					continue
				}

				// Check if we've seen this name before
				if firstPos, exists := varNames[name.Name]; exists {
					pos := fset.Position(name.Pos())
					firstPosLine := fset.Position(firstPos).Line
					results = append(results, LintResult{
						Rule:     r.ID(),
						File:     file,
						Line:     pos.Line,
						Message:  fmt.Sprintf("Duplicate variable name '%s' (first declared at line %d)", name.Name, firstPosLine),
						Severity: r.Severity(),
					})
				} else {
					varNames[name.Name] = name.Pos()
				}
			}
		}
	}

	return results, nil
}

// WAZ005 checks for circular dependencies in resource references
type WAZ005 struct{}

func (r *WAZ005) ID() string {
	return "WAZ005"
}

func (r *WAZ005) Description() string {
	return "Detect circular dependencies in resource references"
}

func (r *WAZ005) Severity() Severity {
	return SeverityError
}

func (r *WAZ005) Check(file string) ([]LintResult, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var results []LintResult

	// Build a dependency graph
	vars := make(map[string]*varInfo)

	// First pass: collect all variables and their dependencies
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR {
			continue
		}

		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			for i, name := range valueSpec.Names {
				if name.Name == "_" {
					continue
				}

				deps := []string{}
				if i < len(valueSpec.Values) {
					deps = r.extractDependencies(valueSpec.Values[i])
				}

				vars[name.Name] = &varInfo{
					pos:  name.Pos(),
					deps: deps,
				}
			}
		}
	}

	// Second pass: check for circular dependencies
	for varName := range vars {
		if r.hasCircularDependency(varName, vars, make(map[string]bool), []string{}) {
			pos := fset.Position(vars[varName].pos)
			results = append(results, LintResult{
				Rule:     r.ID(),
				File:     file,
				Line:     pos.Line,
				Message:  fmt.Sprintf("Circular dependency detected involving variable '%s'", varName),
				Severity: r.Severity(),
			})
		}
	}

	return results, nil
}

func (r *WAZ005) extractDependencies(expr ast.Expr) []string {
	deps := make(map[string]bool)
	r.extractDependenciesRecursive(expr, deps)

	result := make([]string, 0, len(deps))
	for dep := range deps {
		result = append(result, dep)
	}
	return result
}

func (r *WAZ005) extractDependenciesRecursive(expr ast.Expr, deps map[string]bool) {
	if expr == nil {
		return
	}

	switch e := expr.(type) {
	case *ast.Ident:
		// Direct variable reference
		if e.Name != "_" && !isBuiltinType(e.Name) && !isKeyword(e.Name) {
			deps[e.Name] = true
		}

	case *ast.CompositeLit:
		for _, elt := range e.Elts {
			if kv, ok := elt.(*ast.KeyValueExpr); ok {
				r.extractDependenciesRecursive(kv.Value, deps)
			} else {
				r.extractDependenciesRecursive(elt, deps)
			}
		}

	case *ast.SelectorExpr:
		// For expressions like MyResource.Id, we want MyResource
		r.extractDependenciesRecursive(e.X, deps)

	case *ast.CallExpr:
		for _, arg := range e.Args {
			r.extractDependenciesRecursive(arg, deps)
		}

	case *ast.UnaryExpr:
		r.extractDependenciesRecursive(e.X, deps)

	case *ast.BinaryExpr:
		r.extractDependenciesRecursive(e.X, deps)
		r.extractDependenciesRecursive(e.Y, deps)

	case *ast.ParenExpr:
		r.extractDependenciesRecursive(e.X, deps)
	}
}

func (r *WAZ005) hasCircularDependency(varName string, vars map[string]*varInfo, visited map[string]bool, path []string) bool {
	// Check if we've already seen this variable in the current path
	for _, p := range path {
		if p == varName {
			return true
		}
	}

	// Check if we've already visited this variable
	if visited[varName] {
		return false
	}

	visited[varName] = true
	path = append(path, varName)

	// Check dependencies
	varData, exists := vars[varName]
	if !exists {
		return false
	}

	for _, dep := range varData.deps {
		if r.hasCircularDependency(dep, vars, visited, path) {
			return true
		}
	}

	return false
}
