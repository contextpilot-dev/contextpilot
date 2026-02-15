package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jitin-nhz/contextpilot/internal/session"
	"github.com/spf13/cobra"
)

var (
	saveTask      string
	saveGoal      string
	saveState     string
	saveNotes     string
	saveQuick     bool
)

var saveCmd = &cobra.Command{
	Use:   "save [task]",
	Short: "Save current work session context",
	Long: `Save your current work session context for later resumption.

Captures: task, goal, approaches tried, decisions, current state, next steps.

Examples:
  contextpilot save "Refactoring payment service"
  contextpilot save --task "Auth migration" --state "JWT implemented, testing SSO"
  contextpilot save  # Interactive mode

The session is scoped to your current git branch.`,
	Run: runSave,
}

func runSave(cmd *cobra.Command, args []string) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		os.Exit(1)
	}

	mgr := session.New(cwd)

	// Load existing session or create new
	s, _ := mgr.Load()
	if s == nil {
		s = &session.Session{}
	}

	// Get task from args or flag
	if len(args) > 0 {
		s.Task = strings.Join(args, " ")
	} else if saveTask != "" {
		s.Task = saveTask
	}

	// Apply flags
	if saveGoal != "" {
		s.Goal = saveGoal
	}
	if saveState != "" {
		s.State = saveState
	}
	if saveNotes != "" {
		s.Notes = saveNotes
	}

	// Interactive mode if no task provided
	if s.Task == "" && !saveQuick {
		s = interactiveSession(s)
	} else if s.Task == "" {
		fmt.Println("❌ Please provide a task description")
		fmt.Println()
		fmt.Println("Usage: contextpilot save \"Your task description\"")
		fmt.Println("   or: contextpilot save  # for interactive mode")
		os.Exit(1)
	}

	// Save session
	if err := mgr.Save(s); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error saving session: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Session saved!")
	fmt.Println()
	fmt.Printf("   📝 Task: %s\n", s.Task)
	if s.Goal != "" {
		fmt.Printf("   🎯 Goal: %s\n", s.Goal)
	}
	if s.State != "" {
		fmt.Printf("   📍 State: %s\n", s.State)
	}
	if len(s.Approaches) > 0 {
		fmt.Printf("   🔄 Approaches: %d logged\n", len(s.Approaches))
	}
	if len(s.NextSteps) > 0 {
		fmt.Printf("   ➡️  Next steps: %d items\n", len(s.NextSteps))
	}
	fmt.Println()
	fmt.Println("💡 Run 'contextpilot resume' to restore this context")
}

func interactiveSession(s *session.Session) *session.Session {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("📝 Save Session Context")
	fmt.Println("(Press Enter to skip optional fields)")
	fmt.Println()

	// Task (required)
	if s.Task == "" {
		fmt.Print("Task (what are you working on?): ")
		s.Task = readLine(reader)
		if s.Task == "" {
			fmt.Println("❌ Task is required")
			os.Exit(1)
		}
	} else {
		fmt.Printf("Task [%s]: ", s.Task)
		if input := readLine(reader); input != "" {
			s.Task = input
		}
	}

	// Goal
	fmt.Print("Goal (why?): ")
	if input := readLine(reader); input != "" {
		s.Goal = input
	}

	// Approaches
	fmt.Println("Approaches tried (one per line, empty line to finish):")
	for {
		fmt.Print("  - ")
		input := readLine(reader)
		if input == "" {
			break
		}
		s.Approaches = append(s.Approaches, input)
	}

	// Current state
	fmt.Print("Current state (where did you leave off?): ")
	if input := readLine(reader); input != "" {
		s.State = input
	}

	// Next steps
	fmt.Println("Next steps (one per line, empty line to finish):")
	for {
		fmt.Print("  - ")
		input := readLine(reader)
		if input == "" {
			break
		}
		s.NextSteps = append(s.NextSteps, input)
	}

	// Notes
	fmt.Print("Notes (anything else?): ")
	if input := readLine(reader); input != "" {
		s.Notes = input
	}

	return s
}

func readLine(reader *bufio.Reader) string {
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func init() {
	rootCmd.AddCommand(saveCmd)
	saveCmd.Flags().StringVarP(&saveTask, "task", "t", "", "Task description")
	saveCmd.Flags().StringVarP(&saveGoal, "goal", "g", "", "Goal/purpose")
	saveCmd.Flags().StringVarP(&saveState, "state", "s", "", "Current state")
	saveCmd.Flags().StringVarP(&saveNotes, "notes", "n", "", "Additional notes")
	saveCmd.Flags().BoolVarP(&saveQuick, "quick", "q", false, "Quick save (skip interactive)")
}
