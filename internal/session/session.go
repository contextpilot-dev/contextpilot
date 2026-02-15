package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Session represents a work session context
type Session struct {
	ID          string    `json:"id"`
	Branch      string    `json:"branch"`
	Task        string    `json:"task"`
	Goal        string    `json:"goal,omitempty"`
	Approaches  []string  `json:"approaches,omitempty"`
	Decisions   []string  `json:"decisions,omitempty"`
	State       string    `json:"state,omitempty"`
	NextSteps   []string  `json:"nextSteps,omitempty"`
	Notes       string    `json:"notes,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Manager handles session operations
type Manager struct {
	rootPath    string
	sessionsDir string
}

// New creates a new session Manager
func New(rootPath string) *Manager {
	return &Manager{
		rootPath:    rootPath,
		sessionsDir: filepath.Join(rootPath, ".contextpilot", "sessions"),
	}
}

// Save creates or updates a session
func (m *Manager) Save(s *Session) error {
	if err := os.MkdirAll(m.sessionsDir, 0755); err != nil {
		return fmt.Errorf("failed to create sessions directory: %w", err)
	}

	// Generate ID if new
	if s.ID == "" {
		s.ID = fmt.Sprintf("%d", time.Now().UnixNano())
		s.CreatedAt = time.Now()
	}
	s.UpdatedAt = time.Now()

	// Get current branch if not set
	if s.Branch == "" {
		s.Branch = m.getCurrentBranch()
	}

	// Save to branch-specific file
	filename := fmt.Sprintf("%s.json", sanitizeBranch(s.Branch))
	filepath := filepath.Join(m.sessionsDir, filename)

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write session: %w", err)
	}

	// Also save to history
	return m.appendHistory(s)
}

// Load returns the current session for the branch
func (m *Manager) Load() (*Session, error) {
	branch := m.getCurrentBranch()
	filename := fmt.Sprintf("%s.json", sanitizeBranch(branch))
	filepath := filepath.Join(m.sessionsDir, filename)

	data, err := os.ReadFile(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // No session for this branch
		}
		return nil, fmt.Errorf("failed to read session: %w", err)
	}

	var s Session
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("failed to parse session: %w", err)
	}

	return &s, nil
}

// GeneratePrompt creates a prompt to paste into AI tools
func (m *Manager) GeneratePrompt(s *Session) string {
	if s == nil {
		return ""
	}

	prompt := "## Session Context\n\n"
	prompt += fmt.Sprintf("**Task:** %s\n", s.Task)
	
	if s.Goal != "" {
		prompt += fmt.Sprintf("**Goal:** %s\n", s.Goal)
	}

	if len(s.Approaches) > 0 {
		prompt += "\n**Approaches Tried:**\n"
		for _, a := range s.Approaches {
			prompt += fmt.Sprintf("- %s\n", a)
		}
	}

	if len(s.Decisions) > 0 {
		prompt += "\n**Decisions Made:**\n"
		for _, d := range s.Decisions {
			prompt += fmt.Sprintf("- %s\n", d)
		}
	}

	if s.State != "" {
		prompt += fmt.Sprintf("\n**Current State:** %s\n", s.State)
	}

	if len(s.NextSteps) > 0 {
		prompt += "\n**Next Steps:**\n"
		for _, n := range s.NextSteps {
			prompt += fmt.Sprintf("- %s\n", n)
		}
	}

	if s.Notes != "" {
		prompt += fmt.Sprintf("\n**Notes:** %s\n", s.Notes)
	}

	prompt += fmt.Sprintf("\n---\n*Session saved: %s*\n", s.UpdatedAt.Format("2006-01-02 15:04"))

	return prompt
}

// GetHistory returns session history for current branch
func (m *Manager) GetHistory(limit int) ([]Session, error) {
	historyFile := filepath.Join(m.sessionsDir, "history.json")
	
	data, err := os.ReadFile(historyFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []Session{}, nil
		}
		return nil, err
	}

	var history []Session
	if err := json.Unmarshal(data, &history); err != nil {
		return nil, err
	}

	// Filter by current branch
	branch := m.getCurrentBranch()
	var filtered []Session
	for _, s := range history {
		if s.Branch == branch {
			filtered = append(filtered, s)
		}
	}

	// Limit results
	if limit > 0 && len(filtered) > limit {
		filtered = filtered[len(filtered)-limit:]
	}

	return filtered, nil
}

// Clear removes the current session
func (m *Manager) Clear() error {
	branch := m.getCurrentBranch()
	filename := fmt.Sprintf("%s.json", sanitizeBranch(branch))
	filepath := filepath.Join(m.sessionsDir, filename)
	
	if err := os.Remove(filepath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (m *Manager) appendHistory(s *Session) error {
	historyFile := filepath.Join(m.sessionsDir, "history.json")
	
	var history []Session
	if data, err := os.ReadFile(historyFile); err == nil {
		json.Unmarshal(data, &history)
	}

	// Append new session
	history = append(history, *s)

	// Keep last 100 entries
	if len(history) > 100 {
		history = history[len(history)-100:]
	}

	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(historyFile, data, 0644)
}

func (m *Manager) getCurrentBranch() string {
	// Try to get git branch
	gitHead := filepath.Join(m.rootPath, ".git", "HEAD")
	if data, err := os.ReadFile(gitHead); err == nil {
		content := string(data)
		if len(content) > 16 && content[:16] == "ref: refs/heads/" {
			return content[16 : len(content)-1] // Remove "ref: refs/heads/" and newline
		}
	}
	return "main"
}

func sanitizeBranch(branch string) string {
	// Replace / with _ for filenames
	result := ""
	for _, c := range branch {
		if c == '/' || c == '\\' {
			result += "_"
		} else {
			result += string(c)
		}
	}
	return result
}
