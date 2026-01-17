// Package linter provides lint rules for wetwire-azure-go infrastructure code
package linter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	corelint "github.com/lex00/wetwire-core-go/lint"
)

// Severity is an alias to the core lint Severity type.
type Severity = corelint.Severity

// Severity constants from wetwire-core-go/lint.
const (
	SeverityError   = corelint.SeverityError
	SeverityWarning = corelint.SeverityWarning
	SeverityInfo    = corelint.SeverityInfo
)

// LintResult represents a single lint issue found in a file
type LintResult struct {
	// Rule is the ID of the rule that triggered (e.g., "WAZ001")
	Rule string
	// File is the absolute path to the file
	File string
	// Line is the line number where the issue was found
	Line int
	// Message is a human-readable description of the issue
	Message string
	// Severity indicates how critical the issue is
	Severity Severity
}

// String returns a formatted string representation of the lint result
func (lr LintResult) String() string {
	filename := filepath.Base(lr.File)
	return fmt.Sprintf("%s:%d: [%s] %s (%s)", filename, lr.Line, lr.Severity, lr.Message, lr.Rule)
}

// Rule defines the interface that all lint rules must implement
type Rule interface {
	// ID returns the unique identifier for this rule (e.g., "WAZ001")
	ID() string
	// Description returns a human-readable description of what this rule checks
	Description() string
	// Severity returns the severity level for violations of this rule
	Severity() Severity
	// Check analyzes a file and returns any lint results found
	Check(file string) ([]LintResult, error)
}

// FixableRule defines an interface for rules that can automatically fix issues
type FixableRule interface {
	Rule
	// CanFix returns true if this rule supports auto-fixing
	CanFix() bool
	// Fix applies auto-fix to the file and returns the fixed content
	Fix(file string) (string, error)
}

// Linter runs lint rules on Go files
type Linter struct {
	rules []Rule
}

// NewLinter creates a new linter with all default rules registered
func NewLinter() *Linter {
	l := &Linter{
		rules: []Rule{},
	}
	// Register all default rules
	for _, rule := range AllRules() {
		l.AddRule(rule)
	}
	return l
}

// AddRule adds a rule to the linter
func (l *Linter) AddRule(rule Rule) {
	l.rules = append(l.rules, rule)
}

// CheckFile runs all lint rules on a single file
func (l *Linter) CheckFile(file string) ([]LintResult, error) {
	// Verify file exists
	if _, err := os.Stat(file); err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	// Skip non-Go files
	if !strings.HasSuffix(file, ".go") {
		return nil, nil
	}

	var allResults []LintResult

	// Run all rules on the file
	for _, rule := range l.rules {
		results, err := rule.Check(file)
		if err != nil {
			return nil, fmt.Errorf("rule %s failed: %w", rule.ID(), err)
		}
		allResults = append(allResults, results...)
	}

	return allResults, nil
}

// CheckDirectory runs all lint rules on all Go files in a directory (recursively)
func (l *Linter) CheckDirectory(dir string) ([]LintResult, error) {
	// Verify directory exists
	if _, err := os.Stat(dir); err != nil {
		return nil, fmt.Errorf("directory not found: %w", err)
	}

	var allResults []LintResult

	// Walk through all Go files in the directory
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-Go files
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		// Skip test files
		if strings.HasSuffix(path, "_test.go") {
			return nil
		}

		// Check the file
		results, err := l.CheckFile(path)
		if err != nil {
			return fmt.Errorf("failed to check %s: %w", path, err)
		}

		allResults = append(allResults, results...)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return allResults, nil
}
