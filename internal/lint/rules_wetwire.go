package lint

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

// WAZ020 enforces top-level resource declarations
// Resources should be declared at package level, not inside functions
type WAZ020 struct{}

func (r *WAZ020) ID() string {
	return "WAZ020"
}

func (r *WAZ020) Description() string {
	return "Enforce top-level resource declarations"
}

func (r *WAZ020) Severity() Severity {
	return SeverityWarning
}

func (r *WAZ020) Check(file string) ([]LintResult, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var results []LintResult

	// Check for resource declarations inside function bodies
	ast.Inspect(node, func(n ast.Node) bool {
		// Look for function declarations
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		// Check the function body for variable declarations that look like resources
		if funcDecl.Body == nil {
			return true
		}

		for _, stmt := range funcDecl.Body.List {
			// Look for variable declarations
			declStmt, ok := stmt.(*ast.DeclStmt)
			if !ok {
				continue
			}

			genDecl, ok := declStmt.Decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.VAR {
				continue
			}

			for _, spec := range genDecl.Specs {
				valueSpec, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}

				// Check if this looks like a resource declaration
				for i, value := range valueSpec.Values {
					if r.isResourceDeclaration(value) {
						pos := fset.Position(valueSpec.Pos())
						name := "_"
						if i < len(valueSpec.Names) {
							name = valueSpec.Names[i].Name
						}
						results = append(results, LintResult{
							Rule:     r.ID(),
							File:     file,
							Line:     pos.Line,
							Message:  fmt.Sprintf("Resource '%s' should be declared at package level, not inside function '%s'", name, funcDecl.Name.Name),
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

// isResourceDeclaration checks if an expression looks like an Azure resource declaration
func (r *WAZ020) isResourceDeclaration(expr ast.Expr) bool {
	comp, ok := expr.(*ast.CompositeLit)
	if !ok {
		// Check for pointer to composite literal
		if unary, ok := expr.(*ast.UnaryExpr); ok && unary.Op == token.AND {
			comp, ok = unary.X.(*ast.CompositeLit)
			if !ok {
				return false
			}
		} else {
			return false
		}
	}

	// Check if the composite literal has Name and Location fields (typical Azure resources)
	hasName := false
	hasLocation := false

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
			}
		}
	}

	return hasName && hasLocation
}

// WAZ021 flags deeply nested inline structs (depth > 2)
type WAZ021 struct{}

func (r *WAZ021) ID() string {
	return "WAZ021"
}

func (r *WAZ021) Description() string {
	return "Flag deeply nested inline structs (depth > 2)"
}

func (r *WAZ021) Severity() Severity {
	return SeverityWarning
}

func (r *WAZ021) Check(file string) ([]LintResult, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var results []LintResult
	reported := make(map[token.Pos]bool) // Track reported positions to avoid duplicates

	// Find all composite literals and check their nesting depth
	ast.Inspect(node, func(n ast.Node) bool {
		comp, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}

		// Find deeply nested inline structs
		r.findDeeplyNested(comp, 0, fset, file, &results, reported)

		return true
	})

	return results, nil
}

// findDeeplyNested recursively finds inline structs that are nested too deeply
func (r *WAZ021) findDeeplyNested(expr ast.Expr, depth int, fset *token.FileSet, file string, results *[]LintResult, reported map[token.Pos]bool) {
	comp, ok := expr.(*ast.CompositeLit)
	if !ok {
		// Check for pointer to composite literal
		if unary, ok := expr.(*ast.UnaryExpr); ok && unary.Op == token.AND {
			comp, ok = unary.X.(*ast.CompositeLit)
			if !ok {
				return
			}
		} else {
			return
		}
	}

	// Check each element for nested composite literals
	for _, elt := range comp.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		// Get the nested expression (could be direct or pointer)
		var nestedComp *ast.CompositeLit
		if c, ok := kv.Value.(*ast.CompositeLit); ok {
			nestedComp = c
		} else if unary, ok := kv.Value.(*ast.UnaryExpr); ok && unary.Op == token.AND {
			if c, ok := unary.X.(*ast.CompositeLit); ok {
				nestedComp = c
			}
		}

		if nestedComp != nil {
			newDepth := depth + 1
			if newDepth > 2 {
				pos := fset.Position(nestedComp.Pos())
				if !reported[nestedComp.Pos()] {
					reported[nestedComp.Pos()] = true
					fieldName := ""
					if ident, ok := kv.Key.(*ast.Ident); ok {
						fieldName = ident.Name
					}
					*results = append(*results, LintResult{
						Rule:     r.ID(),
						File:     file,
						Line:     pos.Line,
						Message:  fmt.Sprintf("Deeply nested inline struct at field '%s' (depth %d). Extract to a named variable for better readability", fieldName, newDepth),
						Severity: r.Severity(),
					})
				}
			}
			// Continue checking for even deeper nesting
			r.findDeeplyNested(nestedComp, newDepth, fset, file, results, reported)
		}
	}
}

// WAZ022 recommends extracting SKU/Profile configs to named variables
type WAZ022 struct{}

func (r *WAZ022) ID() string {
	return "WAZ022"
}

func (r *WAZ022) Description() string {
	return "Extract SKU/Profile configs to named variables"
}

func (r *WAZ022) Severity() Severity {
	return SeverityInfo
}

func (r *WAZ022) Check(file string) ([]LintResult, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var results []LintResult

	// Configuration field names that should be extracted
	configFields := map[string]bool{
		"Sku":                true,
		"SKU":                true,
		"Profile":            true,
		"StorageProfile":     true,
		"NetworkProfile":     true,
		"OsProfile":          true,
		"OSProfile":          true,
		"HardwareProfile":    true,
		"DiagnosticsProfile": true,
		"SecurityProfile":    true,
		"IdentityProfile":    true,
	}

	ast.Inspect(node, func(n ast.Node) bool {
		kv, ok := n.(*ast.KeyValueExpr)
		if !ok {
			return true
		}

		// Check if this is a config field
		ident, ok := kv.Key.(*ast.Ident)
		if !ok || !configFields[ident.Name] {
			return true
		}

		// Check if the value is an inline composite literal (not a reference)
		var isInline bool
		switch v := kv.Value.(type) {
		case *ast.CompositeLit:
			isInline = true
		case *ast.UnaryExpr:
			if v.Op == token.AND {
				if _, ok := v.X.(*ast.CompositeLit); ok {
					isInline = true
				}
			}
		}

		if isInline {
			pos := fset.Position(kv.Pos())
			results = append(results, LintResult{
				Rule:     r.ID(),
				File:     file,
				Line:     pos.Line,
				Message:  fmt.Sprintf("Consider extracting '%s' configuration to a named variable for reusability", ident.Name),
				Severity: r.Severity(),
			})
		}

		return true
	})

	return results, nil
}
