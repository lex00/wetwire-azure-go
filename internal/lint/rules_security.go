package lint

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

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
					"AllowBlobPublicAccess":    true,
					"PublicNetworkAccess":      true,
					"DisableLocalAuth":         false, // false is insecure
					"EnableHttpsTrafficOnly":   false, // false is insecure
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
