package cmd

import (
	"fmt"
	"os"

	"github.com/jitin-nhz/contextpilot/internal/mcp"
	"github.com/spf13/cobra"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start MCP server for AI tool integration",
	Long: `Start a Model Context Protocol (MCP) server for native integration
with AI coding tools like Claude Code and Windsurf.

The server communicates via JSON-RPC over stdio.

Add to your MCP config (claude_desktop_config.json or similar):

{
  "mcpServers": {
    "contextpilot": {
      "command": "contextpilot",
      "args": ["mcp"],
      "cwd": "/path/to/your/project"
    }
  }
}

Or with npx:
{
  "mcpServers": {
    "contextpilot": {
      "command": "npx",
      "args": ["-y", "contextpilot", "mcp"]
    }
  }
}

Available tools:
  - contextpilot_save    Save work session context
  - contextpilot_resume  Get saved session context
  - contextpilot_sync    Update context files
  - contextpilot_decision Log architectural decision
  - contextpilot_score   Get context quality score

Available resources:
  - contextpilot://context  Project context (CLAUDE.md/.cursorrules)
  - contextpilot://session  Current work session`,
	Run: runMCP,
}

func runMCP(cmd *cobra.Command, args []string) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	server := mcp.NewServer(cwd, Version)
	if err := server.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "MCP server error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(mcpCmd)
}
