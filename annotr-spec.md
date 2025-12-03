# annotr - Fast Local Code Commenting CLI

**Ship Date:** December 2024  
**Status:** Active Development  
**Goal:** Sub-2-second AI-generated comments for code files using local inference

## Overview

annotr is a terminal-native CLI tool that automatically adds intelligent, contextual comments to code files using local LLM inference via Ollama. Built with Go and BubbleTea for a fast, beautiful TUI experience.

## Core Value Proposition

- **Local & Fast**: Uses Ollama with qwen2.5-coder:1.5b for 1-2 second generation
- **Zero Cost**: No API fees, runs entirely on your machine
- **Smart Context**: Tree-sitter parsing provides accurate code structure awareness
- **Beautiful UX**: Charm stack (BubbleTea, Lipgloss, Bubbles) for polished terminal UI
- **Language Agnostic**: Supports any language with tree-sitter grammar

## Technical Stack

### Core Technologies
- **Language**: Go
- **TUI Framework**: BubbleTea (Charm) - for init screen only
- **UI Components**: Bubbles (Charm component library)
- **Styling**: Lipgloss (Charm styling library)
- **Parser**: tree-sitter (code structure analysis)
- **LLM Providers**: 
  - Ollama (local, default)
  - Claude (Anthropic API)
  - OpenAI (GPT models)
  - Groq (fast cloud inference)
- **Default Model**: qwen2.5-coder:1.5b (via Ollama)

### Why This Stack?
- **Go**: Fast compilation, single binary distribution, excellent concurrency
- **BubbleTea**: Elm-inspired architecture, production-ready TUI framework (only for init)
- **Charm Stack**: Battle-tested by tools like Glow, VHS, Soft Serve
- **tree-sitter**: Industry standard for code parsing (used by GitHub, Neovim)
- **OpenAI-compatible APIs**: One client, multiple providers (Ollama, Claude, OpenAI, Groq)
- **Config approach**: Matches churn - familiar, proven, simple

## Features

### MVP (December Release)

#### First-Time Setup (`annotr init`)
- Detect Ollama installation and list available models
- If no Ollama: prompt for API key (Claude/OpenAI/Groq)
- Select model from available options
- Choose comment style (line/block/docstring)
- Save configuration to `~/.annotr/config.json`
- Show success and exit

#### Single File Mode
```bash
annotr file.go
```
- Parse file with tree-sitter
- Identify functions, classes, complex blocks
- Generate contextual comments via LLM
- Write comments directly to file
- Show success message: "Enjoy your comments! 😉"
- Exit

#### Folder Mode
```bash
annotr ./src
```
- Recursively scan directory for supported files
- For each file: prompt "Process [filename]? (y/n)"
- If yes: process and show success
- If no: skip to next file
- Exit when done

### Configuration Files
**Config Location**: `~/.annotr/config.json` (matching churn's approach)

```json
{
  "version": "1.0.0",
  "apiKeys": {
    "anthropic": "sk-ant-...",
    "openai": "sk-...",
    "groq": "gsk_..."
  },
  "defaultProvider": "ollama",
  "defaultModel": "qwen2.5-coder:1.5b",
  "commentStyle": "line"
}
```

**Models Manifest**: `~/.annotr/models.json` (updatable via `annotr update-models`)

```json
{
  "version": "1.0.0",
  "lastUpdated": "2024-12-02",
  "providers": {
    "anthropic": {
      "apiKeyPattern": "sk-ant-",
      "endpoint": "https://api.anthropic.com/v1/messages",
      "models": [
        {
          "id": "claude-sonnet-4-20250514",
          "name": "Claude Sonnet 4",
          "contextWindow": 200000
        },
        {
          "id": "claude-sonnet-4-5-20250929",
          "name": "Claude Sonnet 4.5",
          "contextWindow": 200000
        }
      ]
    },
    "openai": {
      "apiKeyPattern": "sk-",
      "endpoint": "https://api.openai.com/v1/chat/completions",
      "models": [
        {
          "id": "gpt-4o",
          "name": "GPT-4o",
          "contextWindow": 128000
        },
        {
          "id": "gpt-4o-mini",
          "name": "GPT-4o Mini",
          "contextWindow": 128000
        }
      ]
    },
    "groq": {
      "apiKeyPattern": "gsk_",
      "endpoint": "https://api.groq.com/openai/v1/chat/completions",
      "models": [
        {
          "id": "llama-3.3-70b-versatile",
          "name": "Llama 3.3 70B",
          "contextWindow": 32768
        }
      ]
    },
    "ollama": {
      "endpoint": "http://localhost:11434/v1",
      "requiresKey": false
    }
  }
}
```

## Architecture

### Command Flow

```
User Input → CLI Parser → Action Router
                              ↓
                    ┌─────────┴─────────┐
                    ↓                   ↓
              annotr init        annotr file.go
                    ↓                   ↓
              TUI Config          Parse File
              Screen                   ↓
                    ↓              tree-sitter
              Save Config             ↓
                    ↓              Extract Context
                  Exit                 ↓
                                  Call LLM API
                                       ↓
                                  Parse Response
                                       ↓
                                  Write Comments
                                       ↓
                                  Show Success
                                       ↓
                                     Exit
```

### Key Components

#### 1. CLI Interface (`cmd/`)
- `root.go` - Main command setup
- `init.go` - First-time configuration with TUI
- `run.go` - File/folder processing
- `update.go` - Update models manifest

#### 2. Parser (`internal/parser/`)
- `treesitter.go` - tree-sitter integration
- `extractor.go` - AST node extraction
- `context.go` - Context window building

#### 3. LLM Client (`internal/llm/`)
- `client.go` - Universal API client (OpenAI-compatible)
- `ollama.go` - Ollama-specific handling
- `anthropic.go` - Claude API adapter
- `prompt.go` - Prompt engineering
- `response.go` - Response parsing

#### 4. Config UI (`internal/ui/`)
- `init.go` - BubbleTea init screen
- `models.go` - Model selection component
- `success.go` - Success message component
- `styles.go` - Lipgloss styling

#### 5. Config (`internal/config/`)
- `config.go` - Load/save config.json
- `models.go` - Models manifest handling
- `detect.go` - Ollama detection
- `validate.go` - API key validation

#### 6. File Operations (`internal/fileops/`)
- `reader.go` - File reading
- `writer.go` - Safe file writing
- `scanner.go` - Directory scanning

## Init Flow Detail

The `annotr init` command is the only place where we use a TUI. It's a one-time setup that configures the tool.

### Detection Phase
```
┌─ annotr init ─────────────────────────────┐
│                                           │
│  Checking for Ollama...                   │
│  ✓ Found at localhost:11434               │
│  ✓ 12 models available                    │
│                                           │
└───────────────────────────────────────────┘
```

### If Ollama Found - Model Selection
```
┌─ Select Model ────────────────────────────┐
│                                           │
│  ○ qwen2.5-coder:1.5b                     │
│    └─ 1.5B params, fast, good quality    │
│                                           │
│  ○ qwen2.5-coder:7b                       │
│    └─ 7B params, slower, better quality  │
│                                           │
│  ○ deepseek-coder:6.7b                    │
│    └─ 6.7B params, balanced              │
│                                           │
│  [Skip and use API key instead]           │
│                                           │
└───────────────────────────────────────────┘
```

### If No Ollama - Provider Selection
```
┌─ Select Provider ─────────────────────────┐
│                                           │
│  ○ Claude (Anthropic)                     │
│  ○ OpenAI                                 │
│  ○ Groq                                   │
│                                           │
│  [Install Ollama for local/free option]   │
│                                           │
└───────────────────────────────────────────┘
```

### API Key Entry
```
┌─ Configure Claude ────────────────────────┐
│                                           │
│  API Key: sk-ant-**********************   │
│                                           │
│  Models available:                        │
│  ○ claude-sonnet-4-20250514               │
│  ○ claude-sonnet-4-5-20250929             │
│                                           │
└───────────────────────────────────────────┘
```

### Comment Style Selection
```
┌─ Comment Style ───────────────────────────┐
│                                           │
│  ○ Line comments                          │
│    // comment or # comment                │
│                                           │
│  ○ Block comments                         │
│    /* comment */ or """ comment """       │
│                                           │
│  ○ Language-specific                      │
│    Auto-detect best format per language   │
│                                           │
└───────────────────────────────────────────┘
```

### Confirmation & Save
```
┌─ Configuration Complete ──────────────────┐
│                                           │
│  ✓ Provider: Ollama                       │
│  ✓ Model: qwen2.5-coder:1.5b             │
│  ✓ Style: Line comments                   │
│                                           │
│  Saved to: ~/.annotr/config.json          │
│                                           │
│  Ready to use! Try:                       │
│  $ annotr file.go                         │
│                                           │
└───────────────────────────────────────────┘
```

## Success Messages

After processing files, annotr shows a simple success message and exits.

### Single File
```bash
$ annotr main.go

Processing main.go...
✓ Added 12 comments

Enjoy your comments! 😉
```

### Folder Mode
```bash
$ annotr ./src

Process main.go? (y/n): y
✓ Added 12 comments to main.go

Process utils.go? (y/n): y
✓ Added 8 comments to utils.go

Process config.go? (y/n): n
Skipped config.go

Done! Commented 2 of 3 files.
Enjoy your comments! 😉
```

## Prompt Engineering

### Context Window Strategy
For each commentable code block:

```
Language: {language}
File: {filename}

Code Context:
{surrounding_code}

Target Block:
{code_to_comment}

Task: Generate a concise, accurate comment for the target block.
Style: {style_preference}
Length: {max_length} characters

Comment:
```

### Comment Targets (via tree-sitter)
- Function declarations
- Class definitions
- Complex conditional blocks
- Loop constructs
- Error handling sections
- Public API surfaces

## Performance Goals

- **Single File**: < 2 seconds end-to-end
- **10 Files**: < 15 seconds
- **100 Files**: < 2 minutes (with progress)
- **Memory**: < 50MB base usage
- **Binary Size**: < 15MB

## Installation

### Via Go Install
```bash
go install github.com/yourusername/annotr@latest
```

### Via Homebrew (Post-Launch)
```bash
brew install annotr
```

### First-Time Setup
```bash
# Run init to configure
annotr init

# If you have Ollama:
# → Detects installed models
# → Select one and you're done

# If you don't have Ollama:
# → Enter API key (Claude/OpenAI/Groq)
# → Select model from provider
# → You're ready to go
```

### Optional: Install Ollama (for local/free)
```bash
# macOS/Linux
curl -fsSL https://ollama.com/install.sh | sh

# Pull recommended model
ollama pull qwen2.5-coder:1.5b

# Then run annotr init
```

## Usage Examples

### Basic File Comment
```bash
# Add comments to a single file
annotr file.go
# → Auto-processes and writes comments
# → Shows: "Enjoy your comments! 😉"

# No flags needed - it just works
```

### Folder Processing
```bash
# Process all files in directory
annotr ./src

# For each file:
# → "Process main.go? (y/n)"
# → If yes: processes and shows success
# → If no: skips to next file
```

### First Time Setup
```bash
# Initialize configuration
annotr init

# Interactive TUI shows:
# 1. Detect Ollama or prompt for API key
# 2. Select model from available options
# 3. Choose comment style (line/block/docstring)
# 4. Save and confirm
```

## Development Roadmap

### Phase 1: MVP (December 2024) ✓
- [x] Name finalized: annotr
- [ ] Basic CLI structure (cobra)
- [ ] Config system matching churn's approach
- [ ] tree-sitter integration for Go/TypeScript/Python
- [ ] LLM client (OpenAI-compatible)
- [ ] Provider adapters (Ollama, Claude, OpenAI, Groq)
- [ ] Init TUI with BubbleTea
- [ ] Single file processing
- [ ] Folder mode with confirmation prompts
- [ ] Models manifest system

### Phase 2: Polish (January 2025)
- [ ] More language support (Rust, JavaScript, Java, C++)
- [ ] Error handling and edge cases
- [ ] Unit tests for core components
- [ ] Integration tests
- [ ] Performance optimization
- [ ] Better success messages/output
- [ ] Update models manifest command

### Phase 3: Enhancement (February 2025)
- [ ] Custom prompt templates
- [ ] Style presets (JSDoc, GoDoc, PyDoc, RustDoc)
- [ ] Git integration (only uncommitted files)
- [ ] CI/CD mode (--ci flag for non-interactive)
- [ ] Batch operations improvements
- [ ] Support for more providers (Mistral, Cohere, etc.)

### Phase 4: Community (March 2025)
- [ ] Homebrew formula
- [ ] VSCode extension (optional companion)
- [ ] GitHub Action
- [ ] Documentation site
- [ ] Community prompt library
- [ ] Contributing guide

## Differentiation

### vs. GitHub Copilot
- **Local-first option**: Run completely offline with Ollama
- **Zero cost option**: No subscription required
- **Batch mode**: Process entire codebases at once
- **Customizable**: Full control over prompts/models/providers

### vs. Aider/Other AI Coders
- **Focused scope**: Just commenting, does it exceptionally well
- **Non-invasive**: Doesn't change code logic
- **Fast**: Optimized for single task
- **Simple**: No agent complexity

### vs. IDE Extensions
- **Universal**: Works with any editor
- **Scriptable**: CI/CD integration
- **Portable**: Single binary
- **Terminal-native**: Fits developer workflow
- **Provider agnostic**: Use local or cloud models

## Technical Challenges & Solutions

### Challenge 1: Fast Inference
**Problem**: Need sub-2-second response times  
**Solution**: 
- Default to smallest effective local model (qwen2.5-coder:1.5b)
- Keep Ollama warm with health checks
- Cloud APIs (Claude/OpenAI/Groq) are also fast
- Stream responses for perceived speed
- Cache common patterns

### Challenge 2: Accurate Context
**Problem**: Comments need understanding of surrounding code  
**Solution**:
- tree-sitter provides perfect AST
- Include parent scopes in context
- Add import statements
- Function signatures of callers

### Challenge 3: File Safety
**Problem**: Don't corrupt user's code  
**Solution**:
- Atomic writes (write to temp, rename)
- Validate syntax before applying
- User relies on git for undo (git diff/revert)
- Keep it simple - trust the version control

### Challenge 4: Cross-Platform
**Problem**: Works on Mac/Linux/Windows  
**Solution**:
- Go's excellent cross-compilation
- Use filepath package for paths
- Test on all platforms in CI
- Single static binary (no dependencies)

## Project Structure

```
annotr/
├── cmd/
│   └── annotr/
│       └── main.go          # Entry point
├── internal/
│   ├── cli/                 # CLI commands
│   │   ├── root.go
│   │   ├── init.go          # Init configuration
│   │   ├── run.go           # Run on file/folder
│   │   └── update.go        # Update models manifest
│   ├── config/              # Configuration
│   │   ├── config.go        # Load/save config.json
│   │   ├── models.go        # Models manifest
│   │   ├── detect.go        # Ollama detection
│   │   └── validate.go      # API key validation
│   ├── parser/              # tree-sitter integration
│   │   ├── parser.go
│   │   ├── extractor.go
│   │   └── languages/
│   ├── llm/                 # LLM clients
│   │   ├── client.go        # Universal client
│   │   ├── ollama.go        # Ollama adapter
│   │   ├── anthropic.go     # Claude adapter
│   │   ├── openai.go        # OpenAI adapter
│   │   ├── groq.go          # Groq adapter
│   │   ├── prompt.go        # Prompt templates
│   │   └── response.go      # Response parsing
│   ├── ui/                  # BubbleTea UI (init only)
│   │   ├── init.go          # Init screen
│   │   ├── models.go        # Model selector
│   │   ├── success.go       # Success message
│   │   └── styles.go        # Lipgloss styles
│   └── fileops/             # File operations
│       ├── reader.go
│       ├── writer.go
│       └── scanner.go       # Directory scanning
├── pkg/                     # Public API (if needed)
├── test/
│   ├── fixtures/
│   └── integration/
├── docs/
│   ├── installation.md
│   └── usage.md
├── .github/
│   └── workflows/
│       ├── test.yml
│       └── release.yml
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Marketing & Launch

### Target Audience
1. **Individual Developers**: Want documented code without manual effort
2. **Code Reviewers**: Need comments before reviewing PRs
3. **Legacy Codebases**: Adding docs to undocumented code
4. **Open Source**: Improving project accessibility
5. **Teams**: Enforcing documentation standards

### Launch Strategy
1. **Dev.to Article**: "I built a CLI that comments your code in 2 seconds"
2. **Hacker News**: Post with live demo GIF
3. **Reddit**: r/golang, r/programming, r/coding
4. **Twitter/X**: Terminal demo video
5. **GitHub**: Polish README with badges, screenshots

### Demo Content
- Animated GIF showing before/after
- Comparison table vs. alternatives
- Real-world example (commenting an open-source project)
- Performance benchmarks

## Success Metrics

### Launch Week
- 100+ GitHub stars
- 50+ Homebrew installs
- Featured on Hacker News front page
- 5+ blog mentions

### First Month
- 500+ stars
- 200+ weekly installs
- 10+ contributors
- 5+ language support requests

### Long Term
- 2000+ stars
- VSCode extension with 10k+ installs
- Community prompt library
- Referenced in "awesome" lists

## License & Distribution

- **License**: MIT
- **Repository**: github.com/yourusername/annotr
- **Package Managers**: 
  - npm (via npx wrapper)
  - cargo (binary crate)
  - PyPI (binary wrapper)
  - Homebrew tap

## Open Questions

1. **Comment Placement**: Above function or inline with code? Configurable?
2. **Existing Comments**: Update, skip, or merge with existing?
3. **Format Standards**: Auto-detect language-specific formats (JSDoc, GoDoc, etc.)?
4. **Model Updates**: How often to refresh models.json? Auto-update check?
5. **Folder Mode**: Should we add a `--yes` flag to skip prompts and auto-process all?
6. **Error Recovery**: If LLM call fails, retry or skip file?

## Resources

### Documentation Links
- BubbleTea: https://github.com/charmbracelet/bubbletea
- tree-sitter: https://tree-sitter.github.io/tree-sitter/
- Ollama: https://ollama.ai/
- Charm Libraries: https://charm.sh/

### Inspiration
- `glow` - Markdown viewer (Charm)
- `lazygit` - Git TUI
- `k9s` - Kubernetes TUI
- `fx` - JSON viewer

### Competition Research
- GitHub Copilot
- Tabnine
- Aider
- Cursor
- Continue.dev

---

**Next Steps:**
1. Initialize Go module
2. Set up basic CLI with cobra/urfave
3. Integrate tree-sitter for Go files
4. Build Ollama client
5. Create simple TUI prototype
6. Test end-to-end flow
7. Iterate on UX
8. Launch! 🚀
