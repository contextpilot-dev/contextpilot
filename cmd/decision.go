package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jitin-nhz/contextpilot/internal/decisions"
	"github.com/spf13/cobra"
)

var (
	listDecisions   bool
	deleteDecision  int
	decisionContext string
)

var decisionCmd = &cobra.Command{
	Use:   "decision [text]",
	Short: "Log architectural decisions",
	Long: `Record architectural decisions that AI tools should know about.

Examples:
  contextpilot decision "Using Redis for sessions instead of JWT"
  contextpilot decision "Chose Prisma over Drizzle" --context "Team already knows Prisma"
  contextpilot decision --list
  contextpilot decision --delete 3

Decisions are stored in .contextpilot/decisions.md and 
automatically included in generated context files.`,
	Run: runDecision,
}

func runDecision(cmd *cobra.Command, args []string) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	mgr := decisions.New(cwd)

	// Handle delete
	if deleteDecision > 0 {
		if err := mgr.Delete(deleteDecision); err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ Deleted decision #%d\n", deleteDecision)
		return
	}

	// Handle list
	if listDecisions {
		decs, err := mgr.List()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error listing decisions: %v\n", err)
			os.Exit(1)
		}

		if len(decs) == 0 {
			fmt.Println("📋 No decisions logged yet")
			fmt.Println()
			fmt.Println("Add one with:")
			fmt.Println("  contextpilot decision \"Your decision here\"")
			return
		}

		fmt.Println("📋 Architectural Decisions")
		fmt.Println()
		
		// Print as table
		fmt.Println("┌─────┬────────────┬────────────────────────────────────────────────────────┐")
		fmt.Println("│  #  │    Date    │ Decision                                               │")
		fmt.Println("├─────┼────────────┼────────────────────────────────────────────────────────┤")
		
		for _, d := range decs {
			text := d.Text
			if len(text) > 54 {
				text = text[:51] + "..."
			}
			// Replace newlines with spaces
			text = sanitizeForTable(text)
			fmt.Printf("│ %3d │ %s │ %-54s │\n", d.ID, d.Date, text)
		}
		
		fmt.Println("└─────┴────────────┴────────────────────────────────────────────────────────┘")
		fmt.Println()
		fmt.Printf("Total: %d decision(s)\n", len(decs))
		return
	}

	// Handle add
	if len(args) == 0 {
		fmt.Println("❌ Please provide a decision to log")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  contextpilot decision \"Your decision here\"")
		fmt.Println("  contextpilot decision --list")
		fmt.Println("  contextpilot decision --delete <id>")
		return
	}

	text := args[0]
	
	// If multiple args, join them (allows unquoted input)
	if len(args) > 1 {
		text = ""
		for _, arg := range args {
			if _, err := strconv.Atoi(arg); err != nil {
				if text != "" {
					text += " "
				}
				text += arg
			}
		}
	}

	decision, err := mgr.Add(text, decisionContext)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error logging decision: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Decision #%d logged!\n", decision.ID)
	fmt.Println()
	fmt.Printf("   📝 %s\n", text)
	if decisionContext != "" {
		fmt.Printf("   📎 Context: %s\n", decisionContext)
	}
	fmt.Println()
	fmt.Println("💡 Run 'contextpilot sync' to include in context files")
}

func sanitizeForTable(s string) string {
	result := ""
	for _, c := range s {
		if c == '\n' || c == '\r' {
			result += " "
		} else {
			result += string(c)
		}
	}
	return result
}

func init() {
	rootCmd.AddCommand(decisionCmd)
	decisionCmd.Flags().BoolVarP(&listDecisions, "list", "l", false, "List all decisions")
	decisionCmd.Flags().IntVarP(&deleteDecision, "delete", "d", 0, "Delete decision by ID")
	decisionCmd.Flags().StringVarP(&decisionContext, "context", "c", "", "Add context/reasoning for the decision")
}
