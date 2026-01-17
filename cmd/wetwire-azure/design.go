package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// newDesignCmd creates the "design" subcommand for AI-assisted infrastructure design.
func newDesignCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "design [prompt]",
		Short: "AI-assisted infrastructure design",
		Long:  `Start an interactive AI-assisted session to design Azure ARM templates.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("design command not yet implemented for azure domain")
		},
	}
}
