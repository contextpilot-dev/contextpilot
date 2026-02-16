<p align="center">
  <img src="assets/logo/contextpilot-logo-main-2k.png" alt="ContextPilot Logo" width="200">
</p>

<h1 align="center">ContextPilot</h1>

<p align="center">
  <strong>Make every AI tool understand your codebase.</strong>
</p>

<p align="center">
  <a href="https://github.com/contextpilot-dev/contextpilot/releases"><img src="https://img.shields.io/github/v/release/contextpilot-dev/contextpilot?style=flat-square" alt="Release"></a>
  <a href="https://github.com/contextpilot-dev/contextpilot/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue?style=flat-square" alt="License"></a>
  <a href="https://contextpilot.dev"><img src="https://img.shields.io/badge/website-contextpilot.dev-orange?style=flat-square" alt="Website"></a>
</p>

---

ContextPilot generates and maintains AI context files (`.cursorrules`, `CLAUDE.md`, `copilot-instructions.md`) and tracks your work sessions so AI coding tools actually understand your project.

## The Problem

AI coding assistants are powerful but have two critical flaws:

1. **They don't know your codebase** — Your tech stack, conventions, patterns, architectural decisions
2. **They forget your sessions** — Every new chat starts from zero, re-explaining everything

## The Solution

```bash
# Generate context files from your codebase
$ contextpilot init

🔍 Analyzing codebase...
   ├── Framework: Next.js ^14.2.0
   ├── Languages: TypeScript, JavaScript
   ├── ORM: Prisma
   └── Styling: Tailwind CSS

📝 Generating context files...
   ├── .cursorrules (Cursor)
   ├── CLAUDE.md (Claude Code)
   └── .github/copilot-instructions.md (Copilot)

✅ Done! Your AI tools now understand your codebase.
```

```bash
# Save your work session
$ contextpilot save "Refactoring auth to use JWT"

# Resume in any editor, any machine
$ contextpilot resume
✅ Session context copied to clipboard!
```

## Installation

### Quick Install (recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/contextpilot-dev/contextpilot/main/scripts/install.sh | sh
```

### Homebrew (macOS/Linux)

```bash
brew tap contextpilot-dev/tap
brew install contextpilot
```

### npm

```bash
npm install -g contextpilot
```

### Go

```bash
go install github.com/contextpilot-dev/contextpilot@latest
```

### Manual Download

Download the latest release for your platform from [GitHub Releases](https://github.com/contextpilot-dev/contextpilot/releases).

## Commands

### Codebase Context

| Command | Description |
|---------|-------------|
| `contextpilot init` | Analyze codebase and generate context files |
| `contextpilot sync` | Update context files after code changes |
| `contextpilot decision "..."` | Log architectural decisions |
| `contextpilot score` | Check your context quality score |

### Session Context

| Command | Description |
|---------|-------------|
| `contextpilot save "task"` | Save current work session |
| `contextpilot resume` | Restore session and copy to clipboard |

### Integration

| Command | Description |
|---------|-------------|
| `contextpilot mcp` | Start MCP server for AI tool integration |

## Quick Start

```bash
# 1. Initialize in your project
cd my-project
contextpilot init

# 2. Log decisions as you make them
contextpilot decision "Using Zustand for state management"

# 3. Save your work session
contextpilot save "Building user dashboard" --state "Table done, working on charts"

# 4. Resume anytime (copies to clipboard)
contextpilot resume

# 5. Check context quality
contextpilot score
```

## MCP Server Integration

ContextPilot includes a Model Context Protocol (MCP) server for native integration with Claude Code, Windsurf, and other AI tools.

Add to your MCP config:

```json
{
  "mcpServers": {
    "contextpilot": {
      "command": "contextpilot",
      "args": ["mcp"],
      "cwd": "/path/to/your/project"
    }
  }
}
```

**Available MCP Tools:**
- `contextpilot_save` — Save work session
- `contextpilot_resume` — Get saved session
- `contextpilot_sync` — Update context files
- `contextpilot_decision` — Log decision
- `contextpilot_score` — Get quality score

**Available MCP Resources:**
- `contextpilot://context` — Project context (CLAUDE.md)
- `contextpilot://session` — Current work session

## Supported AI Tools

| Tool | Context File |
|------|--------------|
| Cursor | `.cursorrules` |
| Claude Code | `CLAUDE.md` |
| GitHub Copilot | `.github/copilot-instructions.md` |
| OpenClaw | `CLAUDE.md` |
| Windsurf | MCP server |

## What Gets Detected

- **Languages:** TypeScript, JavaScript, Python, Go, Rust, and more
- **Frameworks:** Next.js, React, Vue, Express, FastAPI, etc.
- **ORMs:** Prisma, Drizzle, TypeORM, Mongoose
- **Testing:** Vitest, Jest, Mocha, pytest
- **Styling:** Tailwind, Styled Components
- **State:** Zustand, Redux, Jotai
- **Tooling:** ESLint, Prettier, Biome

## Roadmap

- [x] CLI with init, sync, decision, score
- [x] Session save/resume
- [x] MCP server
- [ ] VS Code extension
- [ ] Team handoffs
- [ ] Watch mode (auto-sync)
- [ ] Git hooks
- [ ] Template library
- [ ] Learning mode (APO-inspired)

## Development

```bash
# Clone
git clone https://github.com/contextpilot-dev/contextpilot.git
cd contextpilot

# Build
go build -o contextpilot .

# Test
./contextpilot init --dry-run
```

## License

MIT

---

<p align="center">
  <strong>Built by <a href="https://github.com/contextpilot-dev">ContextPilot Team</a></strong>
  <br>
  <a href="https://contextpilot.dev">contextpilot.dev</a>
</p>
