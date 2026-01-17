package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// newWatchCmd creates the "watch" subcommand for auto-rebuilding on file changes.
func newWatchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "watch [path]",
		Short: "Auto-rebuild on source file changes",
		Long:  `Watch monitors source files for changes and automatically rebuilds.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("watch command not yet implemented for azure domain")
		},
	}
}
