package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jitin-nhz/contextpilot/internal/analyzer"
	"github.com/jitin-nhz/contextpilot/internal/decisions"
	"github.com/jitin-nhz/contextpilot/internal/generator"
	"github.com/jitin-nhz/contextpilot/internal/session"
)

// JSON-RPC types
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// MCP Protocol types
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeResult struct {
	ProtocolVersion string       `json:"protocolVersion"`
	ServerInfo      ServerInfo   `json:"serverInfo"`
	Capabilities    Capabilities `json:"capabilities"`
}

type Capabilities struct {
	Tools     *ToolsCapability     `json:"tools,omitempty"`
	Resources *ResourcesCapability `json:"resources,omitempty"`
}

type ToolsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

type ResourcesCapability struct {
	Subscribe   bool `json:"subscribe,omitempty"`
	ListChanged bool `json:"listChanged,omitempty"`
}

type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema InputSchema `json:"inputSchema"`
}

type InputSchema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties,omitempty"`
	Required   []string            `json:"required,omitempty"`
}

type Property struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type Resource struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	MimeType    string `json:"mimeType,omitempty"`
}

type ResourceContent struct {
	URI      string `json:"uri"`
	MimeType string `json:"mimeType,omitempty"`
	Text     string `json:"text,omitempty"`
}

// Server handles MCP requests
type Server struct {
	rootPath string
	version  string
}

// NewServer creates a new MCP server
func NewServer(rootPath, version string) *Server {
	return &Server{
		rootPath: rootPath,
		version:  version,
	}
}

// Run starts the MCP server on stdio
func (s *Server) Run() error {
	scanner := bufio.NewScanner(os.Stdin)
	// Increase buffer size for large messages
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var req Request
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			s.sendError(nil, -32700, "Parse error")
			continue
		}

		s.handleRequest(&req)
	}

	return scanner.Err()
}

func (s *Server) handleRequest(req *Request) {
	switch req.Method {
	case "initialize":
		s.handleInitialize(req)
	case "initialized":
		// Notification, no response needed
	case "tools/list":
		s.handleToolsList(req)
	case "tools/call":
		s.handleToolsCall(req)
	case "resources/list":
		s.handleResourcesList(req)
	case "resources/read":
		s.handleResourcesRead(req)
	default:
		s.sendError(req.ID, -32601, fmt.Sprintf("Method not found: %s", req.Method))
	}
}

func (s *Server) handleInitialize(req *Request) {
	result := InitializeResult{
		ProtocolVersion: "2024-11-05",
		ServerInfo: ServerInfo{
			Name:    "contextpilot",
			Version: s.version,
		},
		Capabilities: Capabilities{
			Tools:     &ToolsCapability{},
			Resources: &ResourcesCapability{},
		},
	}
	s.sendResult(req.ID, result)
}

func (s *Server) handleToolsList(req *Request) {
	tools := []Tool{
		{
			Name:        "contextpilot_save",
			Description: "Save current work session context",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"task":  {Type: "string", Description: "What you are working on"},
					"goal":  {Type: "string", Description: "Why you are doing it"},
					"state": {Type: "string", Description: "Current progress/state"},
					"notes": {Type: "string", Description: "Additional notes"},
				},
				Required: []string{"task"},
			},
		},
		{
			Name:        "contextpilot_resume",
			Description: "Get saved session context for current branch",
			InputSchema: InputSchema{
				Type: "object",
			},
		},
		{
			Name:        "contextpilot_sync",
			Description: "Re-analyze codebase and update context files",
			InputSchema: InputSchema{
				Type: "object",
			},
		},
		{
			Name:        "contextpilot_decision",
			Description: "Log an architectural decision",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"text":    {Type: "string", Description: "The decision made"},
					"context": {Type: "string", Description: "Why this decision was made"},
				},
				Required: []string{"text"},
			},
		},
		{
			Name:        "contextpilot_score",
			Description: "Get context quality score",
			InputSchema: InputSchema{
				Type: "object",
			},
		},
	}

	s.sendResult(req.ID, map[string]interface{}{"tools": tools})
}

func (s *Server) handleToolsCall(req *Request) {
	var params struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	if err := json.Unmarshal(req.Params, &params); err != nil {
		s.sendError(req.ID, -32602, "Invalid params")
		return
	}

	var result interface{}
	var err error

	switch params.Name {
	case "contextpilot_save":
		result, err = s.toolSave(params.Arguments)
	case "contextpilot_resume":
		result, err = s.toolResume()
	case "contextpilot_sync":
		result, err = s.toolSync()
	case "contextpilot_decision":
		result, err = s.toolDecision(params.Arguments)
	case "contextpilot_score":
		result, err = s.toolScore()
	default:
		s.sendError(req.ID, -32602, fmt.Sprintf("Unknown tool: %s", params.Name))
		return
	}

	if err != nil {
		s.sendResult(req.ID, map[string]interface{}{
			"content": []map[string]string{
				{"type": "text", "text": fmt.Sprintf("Error: %v", err)},
			},
			"isError": true,
		})
		return
	}

	s.sendResult(req.ID, map[string]interface{}{
		"content": []map[string]string{
			{"type": "text", "text": fmt.Sprintf("%v", result)},
		},
	})
}

func (s *Server) toolSave(args json.RawMessage) (string, error) {
	var params struct {
		Task  string `json:"task"`
		Goal  string `json:"goal"`
		State string `json:"state"`
		Notes string `json:"notes"`
	}
	json.Unmarshal(args, &params)

	mgr := session.New(s.rootPath)
	sess, _ := mgr.Load()
	if sess == nil {
		sess = &session.Session{}
	}

	sess.Task = params.Task
	if params.Goal != "" {
		sess.Goal = params.Goal
	}
	if params.State != "" {
		sess.State = params.State
	}
	if params.Notes != "" {
		sess.Notes = params.Notes
	}

	if err := mgr.Save(sess); err != nil {
		return "", err
	}

	return fmt.Sprintf("Session saved: %s", params.Task), nil
}

func (s *Server) toolResume() (string, error) {
	mgr := session.New(s.rootPath)
	sess, err := mgr.Load()
	if err != nil {
		return "", err
	}
	if sess == nil {
		return "No saved session for this branch", nil
	}

	return mgr.GeneratePrompt(sess), nil
}

func (s *Server) toolSync() (string, error) {
	a := analyzer.New(s.rootPath)
	analysis, err := a.Analyze()
	if err != nil {
		return "", err
	}

	gen := generator.New(analysis, s.rootPath)
	if err := gen.GenerateAll(); err != nil {
		return "", err
	}

	return "Context files updated", nil
}

func (s *Server) toolDecision(args json.RawMessage) (string, error) {
	var params struct {
		Text    string `json:"text"`
		Context string `json:"context"`
	}
	json.Unmarshal(args, &params)

	mgr := decisions.New(s.rootPath)
	dec, err := mgr.Add(params.Text, params.Context)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Decision #%d logged: %s", dec.ID, params.Text), nil
}

func (s *Server) toolScore() (string, error) {
	// Simple score calculation
	score := 0
	files := []string{".cursorrules", "CLAUDE.md", ".github/copilot-instructions.md", ".contextpilot/config.yaml"}
	
	for _, f := range files {
		if _, err := os.Stat(filepath.Join(s.rootPath, f)); err == nil {
			score += 25
		}
	}

	return fmt.Sprintf("Context Quality Score: %d/100", score), nil
}

func (s *Server) handleResourcesList(req *Request) {
	resources := []Resource{
		{
			URI:         "contextpilot://context",
			Name:        "Project Context",
			Description: "Full project context including tech stack, conventions, and decisions",
			MimeType:    "text/markdown",
		},
		{
			URI:         "contextpilot://session",
			Name:        "Current Session",
			Description: "Current work session context",
			MimeType:    "text/markdown",
		},
	}

	s.sendResult(req.ID, map[string]interface{}{"resources": resources})
}

func (s *Server) handleResourcesRead(req *Request) {
	var params struct {
		URI string `json:"uri"`
	}
	if err := json.Unmarshal(req.Params, &params); err != nil {
		s.sendError(req.ID, -32602, "Invalid params")
		return
	}

	var content string

	switch params.URI {
	case "contextpilot://context":
		// Read CLAUDE.md or .cursorrules
		if data, err := os.ReadFile(filepath.Join(s.rootPath, "CLAUDE.md")); err == nil {
			content = string(data)
		} else if data, err := os.ReadFile(filepath.Join(s.rootPath, ".cursorrules")); err == nil {
			content = string(data)
		} else {
			content = "No context files found. Run 'contextpilot init' to generate."
		}

	case "contextpilot://session":
		mgr := session.New(s.rootPath)
		if sess, err := mgr.Load(); err == nil && sess != nil {
			content = mgr.GeneratePrompt(sess)
		} else {
			content = "No saved session for this branch."
		}

	default:
		s.sendError(req.ID, -32602, fmt.Sprintf("Unknown resource: %s", params.URI))
		return
	}

	s.sendResult(req.ID, map[string]interface{}{
		"contents": []ResourceContent{
			{URI: params.URI, MimeType: "text/markdown", Text: content},
		},
	})
}

func (s *Server) sendResult(id interface{}, result interface{}) {
	resp := Response{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
	s.send(resp)
}

func (s *Server) sendError(id interface{}, code int, message string) {
	resp := Response{
		JSONRPC: "2.0",
		ID:      id,
		Error:   &Error{Code: code, Message: message},
	}
	s.send(resp)
}

func (s *Server) send(resp Response) {
	data, _ := json.Marshal(resp)
	fmt.Println(string(data))
}
