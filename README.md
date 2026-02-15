# ContextPilot

**Make every AI tool understand your codebase.**

ContextPilot generates and maintains AI context files (`.cursorrules`, `CLAUDE.md`, `copilot-instructions.md`) so your AI coding tools actually understand your project.

## The Problem

AI coding assistants are powerful but stateless. They don't know:
- Why you chose Postgres over MongoDB
- Your team's coding conventions beyond linting
- Architectural patterns specific to your codebase

Context files fix this, but:
- They go stale within weeks
- New team members don't know they exist
- Managing files for each tool is tedious

## The Solution

```bash
$ contextpilot init

🔍 Analyzing codebase...
   ├── Detected: TypeScript, React, Next.js
   ├── Found: 47 components, 12 API routes, Prisma ORM
   └── Patterns: Feature-based folders, custom hooks

📝 Generating context files...
   ├── .cursorrules (Cursor)
   ├── CLAUDE.md (Claude Code)
   └── .github/copilot-instructions.md (Copilot)

✅ Done! Your AI tools now understand your codebase.
```

## Installation

```bash
# Coming soon!
brew install contextpilot
# or
npm install -g contextpilot
# or
curl -fsSL https://contextpilot.dev/install.sh | sh
```

## Commands

| Command | Description |
|---------|-------------|
| `contextpilot init` | Generate context files for current project |
| `contextpilot sync` | Update context files after code changes |
| `contextpilot decision "..."` | Log architectural decisions |
| `contextpilot score` | Check your context quality |

## Quick Start

```bash
# 1. Initialize in your project
cd my-project
contextpilot init

# 2. Review generated files
cat .cursorrules
cat CLAUDE.md

# 3. Log decisions as you make them
contextpilot decision "Using Zustand for state management"

# 4. Keep context fresh
contextpilot sync  # Run after major changes
```

## Supported AI Tools

- **Cursor** → `.cursorrules`
- **Claude Code** → `CLAUDE.md`
- **GitHub Copilot** → `.github/copilot-instructions.md`
- **OpenClaw** → `CLAUDE.md` (same format)

## Roadmap

- [x] CLI scaffold
- [ ] Basic codebase analysis (tree-sitter)
- [ ] Context file generation
- [ ] Template library
- [ ] Quality scoring
- [ ] Team sync (cloud)
- [ ] VS Code extension
- [ ] Learning mode (APO-inspired optimization)

## Development

```bash
# Build
go build -o contextpilot .

# Run
./contextpilot --help

# Test
go test ./...
```

## License

MIT

---

**Built by [@jitin-nhz](https://github.com/jitin-nhz)** | [contextpilot.dev](https://contextpilot.dev)
