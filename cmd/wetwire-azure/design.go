// Command design provides AI-assisted infrastructure design.
package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/lex00/wetwire-azure-go/internal/kiro"
	"github.com/lex00/wetwire-core-go/agent/agents"
	"github.com/lex00/wetwire-core-go/agent/orchestrator"
	"github.com/lex00/wetwire-core-go/agent/results"
	"github.com/spf13/cobra"
)

// newDesignCmd creates the design subcommand.
func newDesignCmd() *cobra.Command {
	var prompt string
	var outputDir string
	var maxLintCycles int
	var stream bool
	var provider string

	cmd := &cobra.Command{
		Use:   "design",
		Short: "AI-assisted Azure infrastructure generation",
		Long: `Design uses AI to generate Azure ARM templates based on
natural language descriptions.

The AI agent will:
1. Ask clarifying questions if needed
2. Generate Go code using wetwire-azure patterns
3. Run lint and fix any issues
4. Build the final ARM template

Examples:
  wetwire-azure design --prompt "Create a storage account with geo-redundancy"
  wetwire-azure design --output-dir ./infra --prompt "VM with network and storage"
  wetwire-azure design --provider kiro --prompt "Create a virtual network"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if prompt == "" {
				return fmt.Errorf("--prompt flag is required")
			}

			// Handle kiro provider
			if provider == "kiro" {
				return runDesignKiro(prompt)
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// Handle interrupt
			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				<-sigCh
				fmt.Println("\nInterrupted, cleaning up...")
				cancel()
			}()

			// Create session for tracking
			session := results.NewSession("human", "design")

			// Create human developer (reads from stdin)
			reader := bufio.NewReader(os.Stdin)
			developer := orchestrator.NewHumanDeveloper(func() (string, error) {
				return reader.ReadString('\n')
			})

			// Create stream handler if streaming enabled
			var streamHandler agents.StreamHandler
			if stream {
				streamHandler = func(text string) {
					fmt.Print(text)
				}
			}

			// Create runner agent with Azure domain config
			runner, err := agents.NewRunnerAgent(agents.RunnerConfig{
				Domain:        DefaultAzureDomain(),
				WorkDir:       outputDir,
				MaxLintCycles: maxLintCycles,
				Session:       session,
				Developer:     developer,
				StreamHandler: streamHandler,
			})
			if err != nil {
				return fmt.Errorf("creating runner: %w", err)
			}

			fmt.Println("Starting AI-assisted design session...")
			fmt.Println("The AI will ask questions and generate infrastructure code.")
			fmt.Println("Press Ctrl+C to stop.")
			fmt.Println()

			// Run the agent
			if err := runner.Run(ctx, prompt); err != nil {
				return fmt.Errorf("design session failed: %w", err)
			}

			// Print summary
			fmt.Println("\n--- Session Summary ---")
			fmt.Printf("Generated files: %d\n", len(runner.GetGeneratedFiles()))
			for _, f := range runner.GetGeneratedFiles() {
				fmt.Printf("  - %s\n", f)
			}
			fmt.Printf("Lint cycles: %d\n", runner.GetLintCycles())
			fmt.Printf("Lint passed: %v\n", runner.LintPassed())

			return nil
		},
	}

	cmd.Flags().StringVar(&prompt, "prompt", "", "Natural language description of the infrastructure to generate")
	cmd.Flags().StringVarP(&outputDir, "output-dir", "o", ".", "Output directory for generated code")
	cmd.Flags().IntVar(&maxLintCycles, "max-lint-cycles", 3, "Maximum lint/fix cycles")
	cmd.Flags().BoolVar(&stream, "stream", true, "Stream AI responses")
	cmd.Flags().StringVar(&provider, "provider", "core", "AI provider to use (core, kiro)")

	return cmd
}

// runDesignKiro runs the design command using the Kiro provider.
func runDesignKiro(prompt string) error {
	fmt.Println("Launching Kiro design session...")
	return kiro.Launch(prompt)
}
