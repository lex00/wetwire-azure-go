package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// newDiffCmd creates the "diff" subcommand for comparing ARM templates.
func newDiffCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "diff <template1> <template2>",
		Short: "Compare two ARM templates",
		Long:  `Diff performs a semantic comparison of two Azure ARM templates.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("diff command not yet implemented for azure domain")
		},
	}
}
