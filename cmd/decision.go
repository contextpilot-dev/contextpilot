package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listDecisions bool

var decisionCmd = &cobra.Command{
	Use:   "decision [text]",
	Short: "Log architectural decisions",
	Long: `Record architectural decisions that AI tools should know about.

Examples:
  contextpilot decision "Using Redis for sessions instead of JWT"
  contextpilot decision "Chose Prisma over Drizzle for team familiarity"
  contextpilot decision --list

Decisions are stored in .contextpilot/decisions.md and 
automatically included in generated context files.`,
	Run: func(cmd *cobra.Command, args []string) {
		if listDecisions {
			fmt.Println("📋 Architectural Decisions")
			fmt.Println("   └── No decisions logged yet")
			fmt.Println("")
			fmt.Println("Add one with: contextpilot decision \"Your decision here\"")
			return
		}

		if len(args) == 0 {
			fmt.Println("❌ Please provide a decision to log")
			fmt.Println("")
			fmt.Println("Usage: contextpilot decision \"Your decision here\"")
			fmt.Println("       contextpilot decision --list")
			return
		}

		fmt.Printf("📝 Decision: %s\n", args[0])
		fmt.Println("   └── Decision logging not yet implemented")
		fmt.Println("")
		fmt.Println("Coming in v1.0! Star us: github.com/jitin-nhz/contextpilot")
	},
}

func init() {
	rootCmd.AddCommand(decisionCmd)
	decisionCmd.Flags().BoolVarP(&listDecisions, "list", "l", false, "List all decisions")
}
