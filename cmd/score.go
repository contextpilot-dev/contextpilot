package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jitin-nhz/contextpilot/internal/analyzer"
	"github.com/jitin-nhz/contextpilot/internal/decisions"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var scoreCmd = &cobra.Command{
	Use:   "score",
	Short: "Check your context quality score",
	Long: `Analyze your context files and provide a quality score.

Scores based on:
  - Completeness (tech stack, conventions, decisions)
  - Freshness (how recently updated vs code changes)
  - Specificity (generic vs project-specific content)

Provides actionable suggestions for improvement.`,
	Run: runScore,
}

type scoreResult struct {
	total       int
	completeness int
	freshness   int
	decisions   int
	issues      []string
	suggestions []string
}

func runScore(cmd *cobra.Command, args []string) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		os.Exit(1)
	}

	configPath := filepath.Join(cwd, ".contextpilot", "config.yaml")

	// Check if initialized
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println("📊 Context Quality Score: N/A")
		fmt.Println()
		fmt.Println("❌ No context files found")
		fmt.Println()
		fmt.Println("Run 'contextpilot init' to generate context files.")
		os.Exit(0)
	}

	result := calculateScore(cwd)

	// Display score
	emoji := "🟢"
	if result.total < 50 {
		emoji = "🔴"
	} else if result.total < 75 {
		emoji = "🟡"
	}

	fmt.Printf("📊 Context Quality Score: %s %d/100\n", emoji, result.total)
	fmt.Println()

	// Breakdown
	fmt.Println("┌────────────────────┬───────┬─────────────────────────────────┐")
	fmt.Println("│ Category           │ Score │ Status                          │")
	fmt.Println("├────────────────────┼───────┼─────────────────────────────────┤")
	fmt.Printf("│ Completeness       │ %2d/40 │ %-31s │\n", result.completeness, getStatus(result.completeness, 40))
	fmt.Printf("│ Freshness          │ %2d/30 │ %-31s │\n", result.freshness, getStatus(result.freshness, 30))
	fmt.Printf("│ Decisions          │ %2d/30 │ %-31s │\n", result.decisions, getStatus(result.decisions, 30))
	fmt.Println("└────────────────────┴───────┴─────────────────────────────────┘")
	fmt.Println()

	// Issues
	if len(result.issues) > 0 {
		fmt.Println("⚠️  Issues:")
		for _, issue := range result.issues {
			fmt.Printf("   • %s\n", issue)
		}
		fmt.Println()
	}

	// Suggestions
	if len(result.suggestions) > 0 {
		fmt.Println("💡 Suggestions:")
		for _, sug := range result.suggestions {
			fmt.Printf("   • %s\n", sug)
		}
		fmt.Println()
	}

	if result.total >= 75 {
		fmt.Println("🎉 Great job! Your context files are in good shape.")
	}
}

func calculateScore(cwd string) scoreResult {
	result := scoreResult{
		issues:      []string{},
		suggestions: []string{},
	}

	// Check file existence (completeness)
	files := []struct {
		path   string
		points int
		name   string
	}{
		{".cursorrules", 10, ".cursorrules"},
		{"CLAUDE.md", 10, "CLAUDE.md"},
		{".github/copilot-instructions.md", 10, "copilot-instructions.md"},
		{".contextpilot/config.yaml", 10, "config.yaml"},
	}

	for _, f := range files {
		if _, err := os.Stat(filepath.Join(cwd, f.path)); err == nil {
			result.completeness += f.points
		} else {
			result.issues = append(result.issues, fmt.Sprintf("Missing: %s", f.name))
		}
	}

	// Check analysis completeness
	a := analyzer.New(cwd)
	analysis, err := a.Analyze()
	if err == nil {
		if analysis.Framework != nil {
			// Framework detected is good
		} else {
			result.suggestions = append(result.suggestions, "Add framework detection (create package.json or go.mod)")
		}
	}

	// Check freshness
	configPath := filepath.Join(cwd, ".contextpilot", "config.yaml")
	if data, err := os.ReadFile(configPath); err == nil {
		var cfg struct {
			LastSync time.Time `yaml:"lastSync"`
		}
		if yaml.Unmarshal(data, &cfg) == nil && !cfg.LastSync.IsZero() {
			daysSinceSync := int(time.Since(cfg.LastSync).Hours() / 24)
			if daysSinceSync == 0 {
				result.freshness = 30 // Synced today
			} else if daysSinceSync <= 7 {
				result.freshness = 25 // Synced this week
			} else if daysSinceSync <= 30 {
				result.freshness = 15 // Synced this month
				result.suggestions = append(result.suggestions, "Run 'contextpilot sync' — last sync was over a week ago")
			} else {
				result.freshness = 5 // Stale
				result.issues = append(result.issues, fmt.Sprintf("Context files stale (%d days since sync)", daysSinceSync))
			}
		}
	}

	// Check decisions
	decMgr := decisions.New(cwd)
	decs, _ := decMgr.List()
	decCount := len(decs)

	if decCount == 0 {
		result.decisions = 5
		result.suggestions = append(result.suggestions, "Add architectural decisions with 'contextpilot decision \"...\"'")
	} else if decCount < 3 {
		result.decisions = 15
		result.suggestions = append(result.suggestions, fmt.Sprintf("Add more decisions (currently %d, aim for 5+)", decCount))
	} else if decCount < 5 {
		result.decisions = 22
	} else {
		result.decisions = 30 // 5+ decisions is great
	}

	result.total = result.completeness + result.freshness + result.decisions
	return result
}

func getStatus(score, max int) string {
	pct := float64(score) / float64(max) * 100
	if pct >= 80 {
		return "✅ Excellent"
	} else if pct >= 60 {
		return "👍 Good"
	} else if pct >= 40 {
		return "⚠️  Needs improvement"
	}
	return "❌ Poor"
}

func init() {
	rootCmd.AddCommand(scoreCmd)
}
