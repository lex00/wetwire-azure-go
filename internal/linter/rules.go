package linter

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

// AllRules returns all registered lint rules
func AllRules() []Rule {
	return []Rule{
		&WAZ001{},
		&WAZ002{},
		&WAZ003{},
		&WAZ004{},
		&WAZ005{},
		&WAZ006{},
		&WAZ007{},
		&WAZ008{},
		&WAZ301{},
		&WAZ302{},
		&WAZ303{},
		&WAZ304{},
	}
}

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

// varInfo holds information about a variable for dependency analysis
type varInfo struct {
	pos  token.Pos
	deps []string
}

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

// Helper functions

func hasUpperCase(s string) bool {
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			return true
		}
	}
	return false
}

func isBuiltinType(name string) bool {
	builtins := map[string]bool{
		"bool": true, "byte": true, "complex64": true, "complex128": true,
		"error": true, "float32": true, "float64": true, "int": true,
		"int8": true, "int16": true, "int32": true, "int64": true,
		"rune": true, "string": true, "uint": true, "uint8": true,
		"uint16": true, "uint32": true, "uint64": true, "uintptr": true,
		"true": true, "false": true, "nil": true,
	}
	return builtins[name]
}

func isKeyword(name string) bool {
	keywords := map[string]bool{
		"break": true, "case": true, "chan": true, "const": true,
		"continue": true, "default": true, "defer": true, "else": true,
		"fallthrough": true, "for": true, "func": true, "go": true,
		"goto": true, "if": true, "import": true, "interface": true,
		"map": true, "package": true, "range": true, "return": true,
		"select": true, "struct": true, "switch": true, "type": true,
		"var": true,
	}
	return keywords[name]
}

// WAZ006 detects potential secrets and credentials in code
type WAZ006 struct{}

func (r *WAZ006) ID() string {
	return "WAZ006"
}

func (r *WAZ006) Description() string {
	return "Detect potential secrets and credentials in code"
}

func (r *WAZ006) Severity() Severity {
	return SeverityError
}

func (r *WAZ006) Check(file string) ([]LintResult, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var results []LintResult

	// Secret patterns to detect
	secretPatterns := []struct {
		name    string
		pattern string
	}{
		{"AWS Access Key", "AKIA[0-9A-Z]{16}"},
		{"GitHub Token", "ghp_[a-zA-Z0-9]{36}"},
		{"Azure Storage Key", "AccountKey=[a-zA-Z0-9+/=]{20,}"},
		{"Password", "(?i)password\\s*=\\s*[\"'][^\"']+[\"']"},
		{"Secret", "(?i)secret\\s*=\\s*[\"'][^\"']+[\"']"},
		{"API Key", "(?i)api[_-]?key\\s*=\\s*[\"'][^\"']+[\"']"},
	}

	ast.Inspect(node, func(n ast.Node) bool {
		lit, ok := n.(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return true
		}

		value := strings.Trim(lit.Value, `"'`)

		for _, sp := range secretPatterns {
			if r.matchesPattern(value, sp.pattern) {
				pos := fset.Position(lit.Pos())
				results = append(results, LintResult{
					Rule:     r.ID(),
					File:     file,
					Line:     pos.Line,
					Message:  fmt.Sprintf("Potential %s detected. Do not hardcode secrets in code", sp.name),
					Severity: r.Severity(),
				})
				break
			}
		}

		return true
	})

	return results, nil
}

func (r *WAZ006) matchesPattern(value, pattern string) bool {
	// Simple pattern matching for common secret formats
	switch {
	case strings.HasPrefix(pattern, "AKIA"):
		return strings.Contains(value, "AKIA") && len(value) >= 20
	case strings.HasPrefix(pattern, "ghp_"):
		return strings.HasPrefix(value, "ghp_") && len(value) >= 40
	case strings.HasPrefix(pattern, "AccountKey="):
		return strings.Contains(value, "AccountKey=")
	case strings.Contains(pattern, "password"):
		return strings.Contains(strings.ToLower(value), "password=") ||
			strings.Contains(strings.ToLower(value), "password =")
	case strings.Contains(pattern, "secret"):
		return strings.Contains(strings.ToLower(value), "secret=") ||
			strings.Contains(strings.ToLower(value), "secret =")
	case strings.Contains(pattern, "api"):
		lower := strings.ToLower(value)
		return strings.Contains(lower, "api_key=") ||
			strings.Contains(lower, "apikey=") ||
			strings.Contains(lower, "api-key=")
	}
	return false
}

// WAZ007 detects references to sensitive file paths
type WAZ007 struct{}

func (r *WAZ007) ID() string {
	return "WAZ007"
}

func (r *WAZ007) Description() string {
	return "Detect references to sensitive file paths"
}

func (r *WAZ007) Severity() Severity {
	return SeverityWarning
}

func (r *WAZ007) Check(file string) ([]LintResult, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var results []LintResult

	// Sensitive file extensions and names
	sensitivePatterns := []string{
		".env",
		".pem",
		".key",
		".p12",
		".pfx",
		".crt",
		"id_rsa",
		"id_dsa",
		"id_ecdsa",
		"credentials.json",
		"serviceaccount.json",
	}

	ast.Inspect(node, func(n ast.Node) bool {
		lit, ok := n.(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return true
		}

		value := strings.Trim(lit.Value, `"'`)

		for _, pattern := range sensitivePatterns {
			if strings.HasSuffix(value, pattern) || strings.Contains(value, pattern) {
				pos := fset.Position(lit.Pos())
				results = append(results, LintResult{
					Rule:     r.ID(),
					File:     file,
					Line:     pos.Line,
					Message:  fmt.Sprintf("Reference to sensitive file pattern '%s' detected. Ensure secrets are not committed", pattern),
					Severity: r.Severity(),
				})
				break
			}
		}

		return true
	})

	return results, nil
}

// WAZ008 detects insecure default configurations
type WAZ008 struct{}

func (r *WAZ008) ID() string {
	return "WAZ008"
}

func (r *WAZ008) Description() string {
	return "Detect insecure default configurations"
}

func (r *WAZ008) Severity() Severity {
	return SeverityWarning
}

func (r *WAZ008) Check(file string) ([]LintResult, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var results []LintResult

	ast.Inspect(node, func(n ast.Node) bool {
		switch expr := n.(type) {
		case *ast.BasicLit:
			if expr.Kind == token.STRING {
				value := strings.Trim(expr.Value, `"'`)

				// Check for HTTP URLs (should use HTTPS)
				if strings.HasPrefix(value, "http://") && !strings.HasPrefix(value, "http://localhost") &&
					!strings.HasPrefix(value, "http://127.0.0.1") {
					pos := fset.Position(expr.Pos())
					results = append(results, LintResult{
						Rule:     r.ID(),
						File:     file,
						Line:     pos.Line,
						Message:  "Use HTTPS instead of HTTP for secure communication",
						Severity: r.Severity(),
					})
				}
			}

		case *ast.KeyValueExpr:
			// Check for insecure boolean settings
			if ident, ok := expr.Key.(*ast.Ident); ok {
				insecureSettings := map[string]bool{
					"AllowBlobPublicAccess":   true,
					"PublicNetworkAccess":     true,
					"DisableLocalAuth":        false, // false is insecure
					"EnableHttpsTrafficOnly":  false, // false is insecure
					"SupportsHttpsTrafficOnly": false,
				}

				if checkValue, exists := insecureSettings[ident.Name]; exists {
					if lit, ok := expr.Value.(*ast.Ident); ok {
						boolValue := lit.Name == "true"
						if boolValue == checkValue {
							pos := fset.Position(expr.Pos())
							results = append(results, LintResult{
								Rule:     r.ID(),
								File:     file,
								Line:     pos.Line,
								Message:  fmt.Sprintf("Insecure setting: %s=%v may expose resources to security risks", ident.Name, boolValue),
								Severity: r.Severity(),
							})
						}
					}
				}
			}
		}

		return true
	})

	return results, nil
}

// WAZ301 checks that HTTPS-only is enabled for storage accounts
type WAZ301 struct{}

func (r *WAZ301) ID() string {
	return "WAZ301"
}

func (r *WAZ301) Description() string {
	return "Require HTTPS-only for storage accounts"
}

func (r *WAZ301) Severity() Severity {
	return SeverityWarning
}

func (r *WAZ301) Check(file string) ([]LintResult, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var results []LintResult

	ast.Inspect(node, func(n ast.Node) bool {
		kv, ok := n.(*ast.KeyValueExpr)
		if !ok {
			return true
		}

		// Check for EnableHTTPSTrafficOnly or SupportsHttpsTrafficOnly set to false
		if ident, ok := kv.Key.(*ast.Ident); ok {
			if ident.Name == "EnableHTTPSTrafficOnly" || ident.Name == "SupportsHttpsTrafficOnly" {
				if lit, ok := kv.Value.(*ast.Ident); ok && lit.Name == "false" {
					pos := fset.Position(kv.Pos())
					results = append(results, LintResult{
						Rule:     r.ID(),
						File:     file,
						Line:     pos.Line,
						Message:  "HTTPS-only should be enabled for storage accounts. Set to true for secure communication",
						Severity: r.Severity(),
					})
				}
			}
		}

		return true
	})

	return results, nil
}

// WAZ302 detects overly permissive NSG rules
type WAZ302 struct{}

func (r *WAZ302) ID() string {
	return "WAZ302"
}

func (r *WAZ302) Description() string {
	return "Detect overly permissive NSG rules (0.0.0.0/0 or *)"
}

func (r *WAZ302) Severity() Severity {
	return SeverityWarning
}

func (r *WAZ302) Check(file string) ([]LintResult, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var results []LintResult

	ast.Inspect(node, func(n ast.Node) bool {
		kv, ok := n.(*ast.KeyValueExpr)
		if !ok {
			return true
		}

		// Check for SourceAddressPrefix with wildcard values
		if ident, ok := kv.Key.(*ast.Ident); ok {
			if ident.Name == "SourceAddressPrefix" || ident.Name == "DestinationAddressPrefix" {
				if lit, ok := kv.Value.(*ast.BasicLit); ok && lit.Kind == token.STRING {
					value := strings.Trim(lit.Value, `"'`)
					if value == "*" || value == "0.0.0.0/0" || value == "::/0" {
						pos := fset.Position(kv.Pos())
						results = append(results, LintResult{
							Rule:     r.ID(),
							File:     file,
							Line:     pos.Line,
							Message:  fmt.Sprintf("Overly permissive NSG rule: %s='%s'. Restrict to specific IP ranges", ident.Name, value),
							Severity: r.Severity(),
						})
					}
				}
			}
		}

		return true
	})

	return results, nil
}

// WAZ303 checks that resources have tags
type WAZ303 struct{}

func (r *WAZ303) ID() string {
	return "WAZ303"
}

func (r *WAZ303) Description() string {
	return "Require tags on Azure resources"
}

func (r *WAZ303) Severity() Severity {
	return SeverityWarning
}

func (r *WAZ303) Check(file string) ([]LintResult, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var results []LintResult

	// Look for Azure resource struct literals
	ast.Inspect(node, func(n ast.Node) bool {
		comp, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}

		// Check if this is an Azure resource (has Name and Location fields)
		hasName := false
		hasLocation := false
		hasTags := false

		for _, elt := range comp.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}

			if ident, ok := kv.Key.(*ast.Ident); ok {
				switch ident.Name {
				case "Name":
					hasName = true
				case "Location":
					hasLocation = true
				case "Tags":
					hasTags = true
				}
			}
		}

		// If it looks like an Azure resource (has Name and Location) but no Tags
		if hasName && hasLocation && !hasTags {
			pos := fset.Position(comp.Pos())
			results = append(results, LintResult{
				Rule:     r.ID(),
				File:     file,
				Line:     pos.Line,
				Message:  "Azure resource should have Tags for organization and cost management",
				Severity: r.Severity(),
			})
		}

		return true
	})

	return results, nil
}

// WAZ304 checks for deprecated API versions
type WAZ304 struct{}

func (r *WAZ304) ID() string {
	return "WAZ304"
}

func (r *WAZ304) Description() string {
	return "Warn on deprecated API versions"
}

func (r *WAZ304) Severity() Severity {
	return SeverityWarning
}

func (r *WAZ304) Check(file string) ([]LintResult, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var results []LintResult

	// Minimum recommended year for API versions
	minYear := 2021

	ast.Inspect(node, func(n ast.Node) bool {
		kv, ok := n.(*ast.KeyValueExpr)
		if !ok {
			return true
		}

		// Check for APIVersion field
		if ident, ok := kv.Key.(*ast.Ident); ok && ident.Name == "APIVersion" {
			if lit, ok := kv.Value.(*ast.BasicLit); ok && lit.Kind == token.STRING {
				value := strings.Trim(lit.Value, `"'`)
				// Extract year from API version (format: YYYY-MM-DD or YYYY-MM-DD-preview)
				if len(value) >= 4 {
					yearStr := value[:4]
					var year int
					if _, err := fmt.Sscanf(yearStr, "%d", &year); err == nil {
						if year < minYear {
							pos := fset.Position(kv.Pos())
							results = append(results, LintResult{
								Rule:     r.ID(),
								File:     file,
								Line:     pos.Line,
								Message:  fmt.Sprintf("API version '%s' may be deprecated. Consider using a newer version (2021 or later)", value),
								Severity: r.Severity(),
							})
						}
					}
				}
			}
		}

		return true
	})

	return results, nil
}
