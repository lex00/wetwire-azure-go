// Package differ provides semantic comparison of ARM templates.
package differ

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"

	"github.com/lex00/wetwire-azure-go/internal/template"
	coredomain "github.com/lex00/wetwire-core-go/domain"
	"gopkg.in/yaml.v3"
)

// ARMDiffer implements coredomain.Differ for ARM templates.
type ARMDiffer struct{}

// Compile-time check that ARMDiffer implements Differ.
var _ coredomain.Differ = (*ARMDiffer)(nil)

// New creates a new ARM template differ.
func New() *ARMDiffer {
	return &ARMDiffer{}
}

// Diff compares two ARM templates and returns a domain DiffResult.
func (d *ARMDiffer) Diff(ctx *coredomain.Context, file1, file2 string, opts coredomain.DiffOpts) (*coredomain.DiffResult, error) {
	t1, err := loadTemplate(file1)
	if err != nil {
		return nil, fmt.Errorf("failed to load %s: %w", file1, err)
	}

	t2, err := loadTemplate(file2)
	if err != nil {
		return nil, fmt.Errorf("failed to load %s: %w", file2, err)
	}

	return compare(t1, t2, opts)
}

// loadTemplate loads an ARM template from a file (supports JSON and YAML).
func loadTemplate(path string) (*template.ARMTemplate, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var t template.ARMTemplate

	// Try JSON first
	if err := json.Unmarshal(data, &t); err != nil {
		// Try YAML
		if err := yaml.Unmarshal(data, &t); err != nil {
			return nil, fmt.Errorf("failed to parse as JSON or YAML: %w", err)
		}
	}

	return &t, nil
}

// compare compares two ARM templates and returns differences.
func compare(t1, t2 *template.ARMTemplate, opts coredomain.DiffOpts) (*coredomain.DiffResult, error) {
	result := &coredomain.DiffResult{}

	// Build resource maps by name
	res1 := buildResourceMap(t1.Resources)
	res2 := buildResourceMap(t2.Resources)

	// Find added resources (in t2 but not in t1)
	for name, r := range res2 {
		if _, exists := res1[name]; !exists {
			result.Entries = append(result.Entries, coredomain.DiffEntry{
				Resource: name,
				Type:     r.Type,
				Action:   "added",
			})
		}
	}

	// Find removed resources (in t1 but not in t2)
	for name, r := range res1 {
		if _, exists := res2[name]; !exists {
			result.Entries = append(result.Entries, coredomain.DiffEntry{
				Resource: name,
				Type:     r.Type,
				Action:   "removed",
			})
		}
	}

	// Find modified resources
	for name, r1 := range res1 {
		if r2, exists := res2[name]; exists {
			changes := compareResources(r1, r2, opts)
			if len(changes) > 0 {
				result.Entries = append(result.Entries, coredomain.DiffEntry{
					Resource: name,
					Type:     r1.Type,
					Action:   "modified",
					Changes:  changes,
				})
			}
		}
	}

	// Sort entries for consistent output (added, modified, removed)
	sort.Slice(result.Entries, func(i, j int) bool {
		if result.Entries[i].Action != result.Entries[j].Action {
			order := map[string]int{"added": 0, "modified": 1, "removed": 2}
			return order[result.Entries[i].Action] < order[result.Entries[j].Action]
		}
		return result.Entries[i].Resource < result.Entries[j].Resource
	})

	// Calculate summary
	for _, e := range result.Entries {
		switch e.Action {
		case "added":
			result.Summary.Added++
		case "removed":
			result.Summary.Removed++
		case "modified":
			result.Summary.Modified++
		}
	}
	result.Summary.Total = result.Summary.Added + result.Summary.Removed + result.Summary.Modified

	return result, nil
}

// buildResourceMap creates a map of resources by name.
func buildResourceMap(resources []template.ARMResource) map[string]template.ARMResource {
	m := make(map[string]template.ARMResource)
	for _, r := range resources {
		m[r.Name] = r
	}
	return m
}

// compareResources compares two resource definitions and returns changes.
func compareResources(r1, r2 template.ARMResource, opts coredomain.DiffOpts) []string {
	var changes []string

	// Compare type
	if r1.Type != r2.Type {
		changes = append(changes, fmt.Sprintf("Type changed: %s → %s", r1.Type, r2.Type))
	}

	// Compare API version
	if r1.APIVersion != r2.APIVersion {
		changes = append(changes, fmt.Sprintf("apiVersion changed: %s → %s", r1.APIVersion, r2.APIVersion))
	}

	// Compare location
	if r1.Location != r2.Location {
		changes = append(changes, fmt.Sprintf("location changed: %s → %s", r1.Location, r2.Location))
	}

	// Compare kind
	if r1.Kind != r2.Kind {
		changes = append(changes, fmt.Sprintf("kind changed: %s → %s", r1.Kind, r2.Kind))
	}

	// Compare dependsOn
	if !equalStringSlices(r1.DependsOn, r2.DependsOn) {
		changes = append(changes, "dependsOn changed")
	}

	// Compare zones
	if !equalStringSlices(r1.Zones, r2.Zones) {
		changes = append(changes, "zones changed")
	}

	// Compare properties
	propChanges := compareProperties("properties", r1.Properties, r2.Properties, opts)
	changes = append(changes, propChanges...)

	// Compare SKU
	skuChanges := compareProperties("sku", r1.SKU, r2.SKU, opts)
	changes = append(changes, skuChanges...)

	// Compare tags
	tagChanges := compareProperties("tags", r1.Tags, r2.Tags, opts)
	changes = append(changes, tagChanges...)

	// Compare identity
	identityChanges := compareProperties("identity", r1.Identity, r2.Identity, opts)
	changes = append(changes, identityChanges...)

	// Compare plan
	planChanges := compareProperties("plan", r1.Plan, r2.Plan, opts)
	changes = append(changes, planChanges...)

	sort.Strings(changes)
	return changes
}

// compareProperties compares two property values recursively.
func compareProperties(prefix string, v1, v2 interface{}, opts coredomain.DiffOpts) []string {
	var changes []string

	// Handle nil cases
	if v1 == nil && v2 == nil {
		return nil
	}
	if v1 == nil {
		return []string{fmt.Sprintf("%s: added", prefix)}
	}
	if v2 == nil {
		return []string{fmt.Sprintf("%s: removed", prefix)}
	}

	// Use deep comparison
	if !deepEqual(v1, v2, opts) {
		// Try to provide more detail for maps
		if m1, ok := v1.(map[string]interface{}); ok {
			if m2, ok := v2.(map[string]interface{}); ok {
				return comparePropertyMaps(prefix, m1, m2, opts)
			}
		}
		return []string{fmt.Sprintf("%s: modified", prefix)}
	}

	return changes
}

// comparePropertyMaps compares two property maps and returns changes.
func comparePropertyMaps(prefix string, m1, m2 map[string]interface{}, opts coredomain.DiffOpts) []string {
	var changes []string

	// Find added/modified keys
	for key, v2 := range m2 {
		path := key
		if prefix != "" {
			path = prefix + "." + key
		}

		if v1, exists := m1[key]; exists {
			if !deepEqual(v1, v2, opts) {
				changes = append(changes, fmt.Sprintf("%s: modified", path))
			}
		} else {
			changes = append(changes, fmt.Sprintf("%s: added", path))
		}
	}

	// Find removed keys
	for key := range m1 {
		path := key
		if prefix != "" {
			path = prefix + "." + key
		}

		if _, exists := m2[key]; !exists {
			changes = append(changes, fmt.Sprintf("%s: removed", path))
		}
	}

	sort.Strings(changes)
	return changes
}

// deepEqual compares two values deeply, optionally ignoring order.
func deepEqual(a, b interface{}, opts coredomain.DiffOpts) bool {
	if opts.IgnoreOrder {
		a = normalizeValue(a)
		b = normalizeValue(b)
	}
	return reflect.DeepEqual(a, b)
}

// normalizeValue normalizes a value for comparison (e.g., sorting slices).
func normalizeValue(v interface{}) interface{} {
	switch val := v.(type) {
	case []interface{}:
		result := make([]interface{}, len(val))
		copy(result, val)
		return result
	case map[string]interface{}:
		result := make(map[string]interface{})
		for k, v := range val {
			result[k] = normalizeValue(v)
		}
		return result
	default:
		return v
	}
}

// equalStringSlices compares two string slices for equality.
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
