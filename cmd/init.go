package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var initTemplate string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate context files for current project",
	Long: `Analyze your codebase and generate AI context files:
  - .cursorrules (Cursor)
  - CLAUDE.md (Claude Code)
  - .github/copilot-instructions.md (GitHub Copilot)

The generated files help AI tools understand your project's
tech stack, coding conventions, and architectural decisions.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔍 Analyzing codebase...")
		// TODO: Implement codebase analysis
		fmt.Println("   └── Analysis not yet implemented")
		fmt.Println("")
		fmt.Println("📝 This command will:")
		fmt.Println("   ├── Detect tech stack (frameworks, languages, tools)")
		fmt.Println("   ├── Analyze project structure and patterns")
		fmt.Println("   ├── Generate .cursorrules, CLAUDE.md, copilot-instructions.md")
		fmt.Println("   └── Create .contextpilot/config.yaml")
		fmt.Println("")
		fmt.Println("Coming in v1.0! Star us: github.com/jitin-nhz/contextpilot")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&initTemplate, "template", "t", "", "Use a specific template (e.g., nextjs-prisma)")
}
