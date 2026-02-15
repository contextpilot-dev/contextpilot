package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/jitin-nhz/contextpilot/internal/analyzer"
	"github.com/jitin-nhz/contextpilot/internal/generator"
	"github.com/spf13/cobra"
)

var initTemplate string
var dryRun bool

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate context files for current project",
	Long: `Analyze your codebase and generate AI context files:
  - .cursorrules (Cursor)
  - CLAUDE.md (Claude Code)
  - .github/copilot-instructions.md (GitHub Copilot)

The generated files help AI tools understand your project's
tech stack, coding conventions, and architectural decisions.`,
	Run: runInit,
}

func runInit(cmd *cobra.Command, args []string) {
	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("🔍 Analyzing codebase...")

	// Create analyzer and run analysis
	a := analyzer.New(cwd)
	analysis, err := a.Analyze()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error analyzing codebase: %v\n", err)
		os.Exit(1)
	}

	// Sort languages by file count
	sort.Slice(analysis.Languages, func(i, j int) bool {
		return analysis.Languages[i].FileCount > analysis.Languages[j].FileCount
	})

	// Display results
	if len(analysis.Languages) > 0 {
		fmt.Println("   ├── Languages detected:")
		for i, lang := range analysis.Languages {
			prefix := "│  ├──"
			if i == len(analysis.Languages)-1 {
				prefix = "│  └──"
			}
			fmt.Printf("   %s %s (%d files, %.1f%%)\n", prefix, lang.Name, lang.FileCount, lang.Percentage)
		}
	}

	if analysis.Framework != nil {
		fmt.Printf("   ├── Framework: %s", analysis.Framework.Name)
		if analysis.Framework.Version != "" {
			fmt.Printf(" %s", analysis.Framework.Version)
		}
		fmt.Println()
	}

	if analysis.Structure.Type != "" {
		fmt.Printf("   ├── Structure: %s", analysis.Structure.Type)
		if analysis.Structure.SrcDir != "" {
			fmt.Printf(" (src: %s)", analysis.Structure.SrcDir)
		}
		fmt.Println()
	}

	if len(analysis.Structure.Folders) > 0 {
		fmt.Printf("   ├── Folders: %v\n", analysis.Structure.Folders)
	}

	// Show detected patterns
	patterns := []string{}
	if analysis.Patterns.ORM != "" {
		patterns = append(patterns, "ORM: "+analysis.Patterns.ORM)
	}
	if analysis.Patterns.TestFramework != "" {
		patterns = append(patterns, "Tests: "+analysis.Patterns.TestFramework)
	}
	if analysis.Patterns.Styling != "" {
		patterns = append(patterns, "Styling: "+analysis.Patterns.Styling)
	}
	if analysis.Patterns.StateManagement != "" {
		patterns = append(patterns, "State: "+analysis.Patterns.StateManagement)
	}
	if analysis.Patterns.Linter != "" {
		patterns = append(patterns, "Linter: "+analysis.Patterns.Linter)
	}
	if analysis.Patterns.Formatter != "" {
		patterns = append(patterns, "Formatter: "+analysis.Patterns.Formatter)
	}

	if len(patterns) > 0 {
		fmt.Println("   └── Patterns:")
		for i, p := range patterns {
			prefix := "      ├──"
			if i == len(patterns)-1 {
				prefix = "      └──"
			}
			fmt.Printf("   %s %s\n", prefix, p)
		}
	} else {
		fmt.Println("   └── Analysis complete")
	}

	fmt.Println()

	if dryRun {
		fmt.Println("🔍 Dry run - no files written")
		fmt.Println()
		fmt.Println("Would generate:")
		fmt.Println("   ├── .cursorrules")
		fmt.Println("   ├── CLAUDE.md")
		fmt.Println("   ├── .github/copilot-instructions.md")
		fmt.Println("   └── .contextpilot/config.yaml")
		return
	}

	// Generate context files
	fmt.Println("📝 Generating context files...")
	gen := generator.New(analysis, cwd)
	if err := gen.GenerateAll(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error generating files: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("   ├── .cursorrules (Cursor)")
	fmt.Println("   ├── CLAUDE.md (Claude Code)")
	fmt.Println("   ├── .github/copilot-instructions.md (GitHub Copilot)")
	fmt.Println("   └── .contextpilot/config.yaml (ContextPilot config)")
	fmt.Println()
	fmt.Println("✅ Done! Your AI tools now understand your codebase.")
	fmt.Println()
	fmt.Println("💡 Tips:")
	fmt.Println("   • Review and customize the generated files")
	fmt.Println("   • Run 'contextpilot sync' after major code changes")
	fmt.Println("   • Log decisions with 'contextpilot decision \"...\"'")
	fmt.Println()
	fmt.Println("Star us: github.com/jitin-nhz/contextpilot")
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&initTemplate, "template", "t", "", "Use a specific template (e.g., nextjs-prisma)")
	initCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview analysis without generating files")
}
