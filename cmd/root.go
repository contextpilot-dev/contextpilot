package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	Version = "0.1.0-dev"
	Commit  = "none"
	Date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "contextpilot",
	Short: "Make every AI tool understand your codebase",
	Long: `ContextPilot generates and maintains AI context files
(.cursorrules, CLAUDE.md, copilot-instructions.md) so your
AI coding tools actually understand your codebase.

Quick start:
  contextpilot init      Generate context files for current project
  contextpilot sync      Update context files after code changes
  contextpilot decision  Log architectural decisions
  contextpilot score     Check your context quality`,
	Version: fmt.Sprintf("%s (commit: %s, built: %s)", Version, Commit, Date),
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate(`ContextPilot {{.Version}}
`)
}
