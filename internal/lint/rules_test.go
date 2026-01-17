package lint

import (
	"testing"
)

// TestAllRules verifies that all rules are registered
func TestAllRules(t *testing.T) {
	rules := AllRules()

	if len(rules) < 15 {
		t.Errorf("expected at least 15 rules, got %d", len(rules))
	}

	// Check for specific rules
	ruleIDs := make(map[string]bool)
	for _, rule := range rules {
		ruleIDs[rule.ID()] = true
	}

	expectedRules := []string{"WAZ001", "WAZ002", "WAZ003", "WAZ004", "WAZ005", "WAZ006", "WAZ007", "WAZ008", "WAZ020", "WAZ021", "WAZ022", "WAZ301", "WAZ302", "WAZ303", "WAZ304"}
	for _, id := range expectedRules {
		if !ruleIDs[id] {
			t.Errorf("expected rule %s to be registered", id)
		}
	}
}
