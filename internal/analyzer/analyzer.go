package analyzer

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// Analysis represents the result of analyzing a codebase
type Analysis struct {
	RootPath   string       `json:"rootPath"`
	Languages  []Language   `json:"languages"`
	Framework  *Framework   `json:"framework,omitempty"`
	Structure  Structure    `json:"structure"`
	Packages   PackageInfo  `json:"packages"`
	Patterns   Patterns     `json:"patterns"`
	Decisions  []Decision   `json:"decisions"`
}

// Language detected in the codebase
type Language struct {
	Name       string  `json:"name"`
	Extension  string  `json:"extension"`
	FileCount  int     `json:"fileCount"`
	Percentage float64 `json:"percentage"`
}

// Framework detected (Next.js, Express, FastAPI, etc.)
type Framework struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

// Structure represents project directory structure
type Structure struct {
	Type       string   `json:"type"` // "monorepo", "standard", "flat"
	SrcDir     string   `json:"srcDir,omitempty"`
	Folders    []string `json:"folders"`
	EntryPoint string   `json:"entryPoint,omitempty"`
}

// PackageInfo from package.json, go.mod, etc.
type PackageInfo struct {
	Manager      string            `json:"manager"` // npm, yarn, pnpm, go, pip
	Dependencies map[string]string `json:"dependencies,omitempty"`
	DevDeps      map[string]string `json:"devDependencies,omitempty"`
}

// Patterns detected in code
type Patterns struct {
	NamingConvention string   `json:"namingConvention"` // camelCase, snake_case, etc.
	ExportStyle      string   `json:"exportStyle"`      // named, default, mixed
	TestFramework    string   `json:"testFramework,omitempty"`
	Linter           string   `json:"linter,omitempty"`
	Formatter        string   `json:"formatter,omitempty"`
	ORM              string   `json:"orm,omitempty"`
	StateManagement  string   `json:"stateManagement,omitempty"`
	Styling          string   `json:"styling,omitempty"`
}

// Decision represents an architectural decision
type Decision struct {
	Date    string `json:"date"`
	Text    string `json:"text"`
	Context string `json:"context,omitempty"`
}

// Analyzer performs codebase analysis
type Analyzer struct {
	rootPath string
	gitIgnore []string
}

// New creates a new Analyzer for the given path
func New(rootPath string) *Analyzer {
	return &Analyzer{
		rootPath: rootPath,
		gitIgnore: []string{
			"node_modules", "vendor", ".git", "dist", "build",
			".next", "__pycache__", ".venv", "venv", ".idea",
			".vscode", "coverage", ".nyc_output",
		},
	}
}

// Analyze performs full codebase analysis
func (a *Analyzer) Analyze() (*Analysis, error) {
	analysis := &Analysis{
		RootPath:  a.rootPath,
		Languages: []Language{},
		Packages:  PackageInfo{Dependencies: make(map[string]string)},
		Patterns:  Patterns{},
		Decisions: []Decision{},
	}

	// Count files by extension
	extCount := make(map[string]int)
	totalFiles := 0

	err := filepath.Walk(a.rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}

		// Skip ignored directories
		if info.IsDir() {
			for _, ignored := range a.gitIgnore {
				if info.Name() == ignored {
					return filepath.SkipDir
				}
			}
			return nil
		}

		// Count by extension
		ext := strings.ToLower(filepath.Ext(path))
		if ext != "" && isCodeFile(ext) {
			extCount[ext]++
			totalFiles++
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Convert to Language structs
	for ext, count := range extCount {
		lang := extensionToLanguage(ext)
		if lang != "" {
			pct := float64(count) / float64(totalFiles) * 100
			analysis.Languages = append(analysis.Languages, Language{
				Name:       lang,
				Extension:  ext,
				FileCount:  count,
				Percentage: pct,
			})
		}
	}

	// Detect framework from package files
	a.detectFramework(analysis)

	// Analyze structure
	a.analyzeStructure(analysis)

	// Detect patterns
	a.detectPatterns(analysis)

	return analysis, nil
}

func (a *Analyzer) detectFramework(analysis *Analysis) {
	// Check package.json
	pkgPath := filepath.Join(a.rootPath, "package.json")
	if data, err := os.ReadFile(pkgPath); err == nil {
		var pkg struct {
			Dependencies    map[string]string `json:"dependencies"`
			DevDependencies map[string]string `json:"devDependencies"`
		}
		if json.Unmarshal(data, &pkg) == nil {
			analysis.Packages.Manager = "npm"
			analysis.Packages.Dependencies = pkg.Dependencies
			analysis.Packages.DevDeps = pkg.DevDependencies

			// Detect framework
			if _, ok := pkg.Dependencies["next"]; ok {
				analysis.Framework = &Framework{Name: "Next.js", Version: pkg.Dependencies["next"]}
			} else if _, ok := pkg.Dependencies["express"]; ok {
				analysis.Framework = &Framework{Name: "Express", Version: pkg.Dependencies["express"]}
			} else if _, ok := pkg.Dependencies["react"]; ok {
				analysis.Framework = &Framework{Name: "React", Version: pkg.Dependencies["react"]}
			} else if _, ok := pkg.Dependencies["vue"]; ok {
				analysis.Framework = &Framework{Name: "Vue.js", Version: pkg.Dependencies["vue"]}
			} else if _, ok := pkg.Dependencies["svelte"]; ok {
				analysis.Framework = &Framework{Name: "Svelte", Version: pkg.Dependencies["svelte"]}
			}

			// Detect ORM
			if _, ok := pkg.Dependencies["prisma"]; ok {
				analysis.Patterns.ORM = "Prisma"
			} else if _, ok := pkg.Dependencies["@prisma/client"]; ok {
				analysis.Patterns.ORM = "Prisma"
			} else if _, ok := pkg.Dependencies["drizzle-orm"]; ok {
				analysis.Patterns.ORM = "Drizzle"
			} else if _, ok := pkg.Dependencies["typeorm"]; ok {
				analysis.Patterns.ORM = "TypeORM"
			} else if _, ok := pkg.Dependencies["mongoose"]; ok {
				analysis.Patterns.ORM = "Mongoose"
			}

			// Detect testing
			if _, ok := pkg.DevDependencies["vitest"]; ok {
				analysis.Patterns.TestFramework = "Vitest"
			} else if _, ok := pkg.DevDependencies["jest"]; ok {
				analysis.Patterns.TestFramework = "Jest"
			} else if _, ok := pkg.DevDependencies["mocha"]; ok {
				analysis.Patterns.TestFramework = "Mocha"
			}

			// Detect styling
			if _, ok := pkg.Dependencies["tailwindcss"]; ok {
				analysis.Patterns.Styling = "Tailwind CSS"
			} else if _, ok := pkg.DevDependencies["tailwindcss"]; ok {
				analysis.Patterns.Styling = "Tailwind CSS"
			} else if _, ok := pkg.Dependencies["styled-components"]; ok {
				analysis.Patterns.Styling = "Styled Components"
			}

			// Detect state management
			if _, ok := pkg.Dependencies["zustand"]; ok {
				analysis.Patterns.StateManagement = "Zustand"
			} else if _, ok := pkg.Dependencies["@reduxjs/toolkit"]; ok {
				analysis.Patterns.StateManagement = "Redux Toolkit"
			} else if _, ok := pkg.Dependencies["jotai"]; ok {
				analysis.Patterns.StateManagement = "Jotai"
			} else if _, ok := pkg.Dependencies["recoil"]; ok {
				analysis.Patterns.StateManagement = "Recoil"
			}

			// Detect linter/formatter
			if _, ok := pkg.DevDependencies["eslint"]; ok {
				analysis.Patterns.Linter = "ESLint"
			}
			if _, ok := pkg.DevDependencies["prettier"]; ok {
				analysis.Patterns.Formatter = "Prettier"
			} else if _, ok := pkg.DevDependencies["biome"]; ok {
				analysis.Patterns.Formatter = "Biome"
			}
		}
	}

	// Check go.mod
	goModPath := filepath.Join(a.rootPath, "go.mod")
	if _, err := os.Stat(goModPath); err == nil {
		analysis.Packages.Manager = "go"
		// Could parse go.mod for dependencies
	}

	// Check pyproject.toml / requirements.txt
	pyprojectPath := filepath.Join(a.rootPath, "pyproject.toml")
	reqPath := filepath.Join(a.rootPath, "requirements.txt")
	if _, err := os.Stat(pyprojectPath); err == nil {
		analysis.Packages.Manager = "poetry/pip"
	} else if _, err := os.Stat(reqPath); err == nil {
		analysis.Packages.Manager = "pip"
	}
}

func (a *Analyzer) analyzeStructure(analysis *Analysis) {
	analysis.Structure.Type = "standard"

	// Check for common directories
	commonDirs := []string{"src", "app", "lib", "components", "pages", "api", "utils", "hooks", "services", "models", "types"}
	foundDirs := []string{}

	for _, dir := range commonDirs {
		if info, err := os.Stat(filepath.Join(a.rootPath, dir)); err == nil && info.IsDir() {
			foundDirs = append(foundDirs, dir)
		}
	}

	analysis.Structure.Folders = foundDirs

	// Detect src directory
	if contains(foundDirs, "src") {
		analysis.Structure.SrcDir = "src"
	} else if contains(foundDirs, "app") {
		analysis.Structure.SrcDir = "app"
	}

	// Check for monorepo indicators
	if _, err := os.Stat(filepath.Join(a.rootPath, "packages")); err == nil {
		analysis.Structure.Type = "monorepo"
	} else if _, err := os.Stat(filepath.Join(a.rootPath, "apps")); err == nil {
		analysis.Structure.Type = "monorepo"
	} else if _, err := os.Stat(filepath.Join(a.rootPath, "pnpm-workspace.yaml")); err == nil {
		analysis.Structure.Type = "monorepo"
	} else if _, err := os.Stat(filepath.Join(a.rootPath, "lerna.json")); err == nil {
		analysis.Structure.Type = "monorepo"
	} else if _, err := os.Stat(filepath.Join(a.rootPath, "turbo.json")); err == nil {
		analysis.Structure.Type = "monorepo"
	}

	// Detect entry point
	entryPoints := []string{"index.ts", "index.js", "main.ts", "main.js", "main.go", "main.py", "app.py"}
	for _, entry := range entryPoints {
		if _, err := os.Stat(filepath.Join(a.rootPath, entry)); err == nil {
			analysis.Structure.EntryPoint = entry
			break
		}
		if analysis.Structure.SrcDir != "" {
			if _, err := os.Stat(filepath.Join(a.rootPath, analysis.Structure.SrcDir, entry)); err == nil {
				analysis.Structure.EntryPoint = filepath.Join(analysis.Structure.SrcDir, entry)
				break
			}
		}
	}
}

func (a *Analyzer) detectPatterns(analysis *Analysis) {
	// Default patterns
	if analysis.Patterns.NamingConvention == "" {
		// Check primary language
		for _, lang := range analysis.Languages {
			if lang.Name == "Go" {
				analysis.Patterns.NamingConvention = "camelCase/PascalCase"
				analysis.Patterns.ExportStyle = "named (capitalized)"
				break
			} else if lang.Name == "Python" {
				analysis.Patterns.NamingConvention = "snake_case"
				break
			} else if lang.Name == "TypeScript" || lang.Name == "JavaScript" {
				analysis.Patterns.NamingConvention = "camelCase"
				analysis.Patterns.ExportStyle = "mixed"
				break
			}
		}
	}
}

// Helper functions

func isCodeFile(ext string) bool {
	codeExts := map[string]bool{
		".js": true, ".ts": true, ".jsx": true, ".tsx": true,
		".go": true, ".py": true, ".rb": true, ".rs": true,
		".java": true, ".kt": true, ".swift": true, ".c": true,
		".cpp": true, ".h": true, ".cs": true, ".php": true,
		".vue": true, ".svelte": true,
	}
	return codeExts[ext]
}

func extensionToLanguage(ext string) string {
	langMap := map[string]string{
		".js":     "JavaScript",
		".ts":     "TypeScript",
		".jsx":    "JavaScript (JSX)",
		".tsx":    "TypeScript (TSX)",
		".go":     "Go",
		".py":     "Python",
		".rb":     "Ruby",
		".rs":     "Rust",
		".java":   "Java",
		".kt":     "Kotlin",
		".swift":  "Swift",
		".c":      "C",
		".cpp":    "C++",
		".h":      "C/C++ Header",
		".cs":     "C#",
		".php":    "PHP",
		".vue":    "Vue",
		".svelte": "Svelte",
	}
	return langMap[ext]
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
