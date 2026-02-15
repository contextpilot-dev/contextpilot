package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/jitin-nhz/contextpilot/internal/session"
	"github.com/spf13/cobra"
)

var (
	resumeNoCopy bool
	resumeFormat string
)

var resumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Restore session context and copy to clipboard",
	Long: `Generate a prompt from your saved session and copy it to clipboard.

Paste this prompt into Cursor, Claude Code, ChatGPT, or any AI tool
to restore your working context.

Examples:
  contextpilot resume           # Copy to clipboard
  contextpilot resume --no-copy # Just print, don't copy
  contextpilot resume --format markdown`,
	Run: runResume,
}

func runResume(cmd *cobra.Command, args []string) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		os.Exit(1)
	}

	mgr := session.New(cwd)
	s, err := mgr.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error loading session: %v\n", err)
		os.Exit(1)
	}

	if s == nil {
		fmt.Println("📋 No saved session for this branch")
		fmt.Println()
		fmt.Println("Save one with: contextpilot save \"Your task description\"")
		return
	}

	// Generate prompt
	prompt := mgr.GeneratePrompt(s)

	// Copy to clipboard (unless --no-copy)
	if !resumeNoCopy {
		if err := copyToClipboard(prompt); err != nil {
			fmt.Fprintf(os.Stderr, "⚠️  Could not copy to clipboard: %v\n", err)
			fmt.Println()
			resumeNoCopy = true // Fall back to printing
		} else {
			fmt.Println("✅ Session context copied to clipboard!")
			fmt.Println()
			fmt.Println("Paste into Cursor, Claude Code, or ChatGPT to resume.")
			fmt.Println()
		}
	}

	// Print preview or full content
	if resumeNoCopy {
		fmt.Println("📋 Session Context:")
		fmt.Println(repeatStr("─", 50))
		fmt.Println(prompt)
		fmt.Println(repeatStr("─", 50))
	} else {
		// Show preview
		fmt.Printf("📝 Task: %s\n", s.Task)
		if s.State != "" {
			fmt.Printf("📍 State: %s\n", s.State)
		}
		if len(s.NextSteps) > 0 {
			fmt.Printf("➡️  Next: %s\n", s.NextSteps[0])
			if len(s.NextSteps) > 1 {
				fmt.Printf("   (+%d more steps)\n", len(s.NextSteps)-1)
			}
		}
	}
}

func copyToClipboard(text string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "linux":
		// Try xclip first, then xsel
		if _, err := exec.LookPath("xclip"); err == nil {
			cmd = exec.Command("xclip", "-selection", "clipboard")
		} else if _, err := exec.LookPath("xsel"); err == nil {
			cmd = exec.Command("xsel", "--clipboard", "--input")
		} else {
			return fmt.Errorf("no clipboard tool found (install xclip or xsel)")
		}
	case "windows":
		cmd = exec.Command("clip")
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	stdin.Write([]byte(text))
	stdin.Close()

	return cmd.Wait()
}

// Helper for string repeat
func repeatStr(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}

func init() {
	rootCmd.AddCommand(resumeCmd)
	resumeCmd.Flags().BoolVar(&resumeNoCopy, "no-copy", false, "Print instead of copying to clipboard")
	resumeCmd.Flags().StringVar(&resumeFormat, "format", "markdown", "Output format (markdown, plain)")
}
