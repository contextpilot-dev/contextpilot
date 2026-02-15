package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jitin-nhz/contextpilot/internal/analyzer"
	"github.com/jitin-nhz/contextpilot/internal/generator"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var forceSyncFlag bool

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Update context files after code changes",
	Long: `Analyze changes since last sync and update context files.

Uses git diff to detect:
  - New files and patterns
  - Deleted or renamed files  
  - Significant code changes

Regenerates context files with latest analysis.`,
	Run: runSync,
}

type configFile struct {
	Version  int       `yaml:"version"`
	LastSync time.Time `yaml:"lastSync"`
}

func runSync(cmd *cobra.Command, args []string) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	configPath := filepath.Join(cwd, ".contextpilot", "config.yaml")

	// Check if initialized
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println("❌ ContextPilot not initialized in this directory")
		fmt.Println()
		fmt.Println("Run 'contextpilot init' first to generate context files.")
		os.Exit(1)
	}

	// Read last sync time
	var lastSync time.Time
	if data, err := os.ReadFile(configPath); err == nil {
		var cfg configFile
		if yaml.Unmarshal(data, &cfg) == nil {
			lastSync = cfg.LastSync
		}
	}

	fmt.Println("🔄 Checking for changes since last sync...")

	// Show git changes if available
	changes := getGitChanges(cwd, lastSync)
	if len(changes) > 0 {
		fmt.Printf("   ├── %d file(s) changed since last sync\n", len(changes))
		// Show up to 5 changes
		shown := 0
		for _, c := range changes {
			if shown >= 5 {
				fmt.Printf("   │  └── ... and %d more\n", len(changes)-5)
				break
			}
			prefix := "├──"
			if shown == len(changes)-1 || shown == 4 {
				prefix = "└──"
			}
			fmt.Printf("   │  %s %s\n", prefix, c)
			shown++
		}
	} else {
		fmt.Println("   ├── No git changes detected (or not a git repo)")
	}

	// Re-run analysis
	fmt.Println("   └── Re-analyzing codebase...")
	fmt.Println()

	a := analyzer.New(cwd)
	analysis, err := a.Analyze()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error analyzing codebase: %v\n", err)
		os.Exit(1)
	}

	// Sort languages
	sort.Slice(analysis.Languages, func(i, j int) bool {
		return analysis.Languages[i].FileCount > analysis.Languages[j].FileCount
	})

	// Generate updated files
	fmt.Println("📝 Updating context files...")
	gen := generator.New(analysis, cwd)
	if err := gen.GenerateAll(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error generating files: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("   ├── .cursorrules")
	fmt.Println("   ├── CLAUDE.md")
	fmt.Println("   ├── .github/copilot-instructions.md")
	fmt.Println("   └── .contextpilot/config.yaml")
	fmt.Println()
	fmt.Println("✅ Context files updated!")

	// Show summary
	if analysis.Framework != nil {
		fmt.Printf("\n📊 Current state: %s", analysis.Framework.Name)
		if len(analysis.Languages) > 0 {
			fmt.Printf(" + %s", analysis.Languages[0].Name)
		}
		fmt.Println()
	}
}

func getGitChanges(cwd string, since time.Time) []string {
	var changes []string

	// Check if git repo
	gitDir := filepath.Join(cwd, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return changes
	}

	// Get changed files
	var cmd *exec.Cmd
	if since.IsZero() {
		// No last sync, show recent changes
		cmd = exec.Command("git", "diff", "--name-only", "HEAD~10", "--", ".")
	} else {
		// Changes since last sync
		sinceStr := since.Format("2006-01-02T15:04:05")
		cmd = exec.Command("git", "log", "--since="+sinceStr, "--name-only", "--pretty=format:", "--", ".")
	}
	cmd.Dir = cwd

	output, err := cmd.Output()
	if err != nil {
		return changes
	}

	// Parse output
	seen := make(map[string]bool)
	for _, line := range strings.Split(string(output), "\n") {
		line = strings.TrimSpace(line)
		if line != "" && !seen[line] {
			// Filter out non-code files
			if isRelevantFile(line) {
				changes = append(changes, line)
				seen[line] = true
			}
		}
	}

	return changes
}

func isRelevantFile(path string) bool {
	// Skip common non-code files
	skip := []string{
		"package-lock.json", "yarn.lock", "pnpm-lock.yaml",
		"go.sum", ".DS_Store", "Thumbs.db",
	}
	for _, s := range skip {
		if strings.HasSuffix(path, s) {
			return false
		}
	}

	// Skip hidden files and directories
	parts := strings.Split(path, "/")
	for _, p := range parts {
		if strings.HasPrefix(p, ".") && p != ".github" {
			return false
		}
	}

	return true
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.Flags().BoolVarP(&forceSyncFlag, "force", "f", false, "Force sync even if no changes detected")
}
