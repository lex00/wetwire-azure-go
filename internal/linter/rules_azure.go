package linter

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

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
