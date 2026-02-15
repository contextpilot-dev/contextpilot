package decisions

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Decision represents an architectural decision
type Decision struct {
	ID      int
	Date    string
	Text    string
	Context string
}

// Manager handles decision operations
type Manager struct {
	rootPath string
	filePath string
}

// New creates a new decision Manager
func New(rootPath string) *Manager {
	return &Manager{
		rootPath: rootPath,
		filePath: filepath.Join(rootPath, ".contextpilot", "decisions.md"),
	}
}

// Add adds a new decision
func (m *Manager) Add(text string, context string) (*Decision, error) {
	// Ensure .contextpilot directory exists
	dir := filepath.Dir(m.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Get next ID
	decisions, _ := m.List()
	nextID := 1
	if len(decisions) > 0 {
		nextID = decisions[len(decisions)-1].ID + 1
	}

	decision := &Decision{
		ID:      nextID,
		Date:    time.Now().Format("2006-01-02"),
		Text:    text,
		Context: context,
	}

	// Append to file
	f, err := os.OpenFile(m.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	// Check if file is empty (needs header)
	info, _ := f.Stat()
	if info.Size() == 0 {
		header := `# Architectural Decisions
# Managed by ContextPilot — https://contextpilot.dev
# Add decisions with: contextpilot decision "Your decision here"

`
		f.WriteString(header)
	}

	// Write decision
	entry := fmt.Sprintf("## [%d] %s\n**Date:** %s\n\n%s\n", 
		decision.ID, summarize(text, 60), decision.Date, text)
	if context != "" {
		entry += fmt.Sprintf("\n**Context:** %s\n", context)
	}
	entry += "\n---\n\n"

	if _, err := f.WriteString(entry); err != nil {
		return nil, fmt.Errorf("failed to write decision: %w", err)
	}

	return decision, nil
}

// List returns all decisions
func (m *Manager) List() ([]Decision, error) {
	if _, err := os.Stat(m.filePath); os.IsNotExist(err) {
		return []Decision{}, nil
	}

	f, err := os.Open(m.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	var decisions []Decision
	var current *Decision
	var textLines []string

	scanner := bufio.NewScanner(f)
	idPattern := regexp.MustCompile(`^## \[(\d+)\]`)
	datePattern := regexp.MustCompile(`^\*\*Date:\*\* (.+)$`)

	for scanner.Scan() {
		line := scanner.Text()

		// Check for new decision header
		if matches := idPattern.FindStringSubmatch(line); matches != nil {
			// Save previous decision
			if current != nil {
				current.Text = strings.TrimSpace(strings.Join(textLines, "\n"))
				decisions = append(decisions, *current)
			}

			id, _ := strconv.Atoi(matches[1])
			current = &Decision{ID: id}
			textLines = []string{}
			continue
		}

		if current == nil {
			continue
		}

		// Parse date
		if matches := datePattern.FindStringSubmatch(line); matches != nil {
			current.Date = matches[1]
			continue
		}

		// Skip separators and empty lines at start
		if line == "---" || (len(textLines) == 0 && line == "") {
			continue
		}

		// Collect text
		if !strings.HasPrefix(line, "**Context:**") {
			textLines = append(textLines, line)
		} else {
			current.Context = strings.TrimPrefix(line, "**Context:** ")
		}
	}

	// Don't forget last decision
	if current != nil {
		current.Text = strings.TrimSpace(strings.Join(textLines, "\n"))
		decisions = append(decisions, *current)
	}

	return decisions, scanner.Err()
}

// Delete removes a decision by ID
func (m *Manager) Delete(id int) error {
	decisions, err := m.List()
	if err != nil {
		return err
	}

	// Filter out the decision
	var remaining []Decision
	found := false
	for _, d := range decisions {
		if d.ID != id {
			remaining = append(remaining, d)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("decision #%d not found", id)
	}

	// Rewrite file
	return m.rewrite(remaining)
}

func (m *Manager) rewrite(decisions []Decision) error {
	f, err := os.Create(m.filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	header := `# Architectural Decisions
# Managed by ContextPilot — https://contextpilot.dev
# Add decisions with: contextpilot decision "Your decision here"

`
	f.WriteString(header)

	for _, d := range decisions {
		entry := fmt.Sprintf("## [%d] %s\n**Date:** %s\n\n%s\n",
			d.ID, summarize(d.Text, 60), d.Date, d.Text)
		if d.Context != "" {
			entry += fmt.Sprintf("\n**Context:** %s\n", d.Context)
		}
		entry += "\n---\n\n"
		f.WriteString(entry)
	}

	return nil
}

// GetForContext returns decisions formatted for inclusion in context files
func (m *Manager) GetForContext() string {
	decisions, err := m.List()
	if err != nil || len(decisions) == 0 {
		return ""
	}

	var sb strings.Builder
	for _, d := range decisions {
		sb.WriteString(fmt.Sprintf("- **%s:** %s\n", d.Date, d.Text))
	}
	return sb.String()
}

// summarize truncates text to maxLen with ellipsis
func summarize(text string, maxLen int) string {
	// Get first line only
	if idx := strings.Index(text, "\n"); idx != -1 {
		text = text[:idx]
	}
	
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen-3] + "..."
}
