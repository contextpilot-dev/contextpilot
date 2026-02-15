package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Update context files after code changes",
	Long: `Analyze changes since last sync and update context files.

Uses git diff to detect:
  - New files and patterns
  - Deleted or renamed files
  - Significant code changes

Only updates relevant sections, preserving manual edits.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔄 Checking for changes since last sync...")
		// TODO: Implement git diff analysis
		fmt.Println("   └── Sync not yet implemented")
		fmt.Println("")
		fmt.Println("📝 This command will:")
		fmt.Println("   ├── Analyze git diff since last sync")
		fmt.Println("   ├── Detect new patterns, files, dependencies")
		fmt.Println("   ├── Update context files incrementally")
		fmt.Println("   └── Show summary of changes")
		fmt.Println("")
		fmt.Println("Coming in v1.0! Star us: github.com/jitin-nhz/contextpilot")
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
