package lint

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
		&WAZ020{},
		&WAZ021{},
		&WAZ022{},
		&WAZ301{},
		&WAZ302{},
		&WAZ303{},
		&WAZ304{},
	}
}
