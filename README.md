![annotr logo](annotr-logo.png)

# annotr

Fast local code commenting CLI. Add AI-generated comments to your code in under 2 seconds.

## Features

- **Local or API**: Run Ollama locally or use Claude, OpenAI, or Groq
- **Fast by Default**: Local models can generate in ~1-2 seconds
- **Smart Context**: Tree-sitter parsing provides accurate code structure awareness
- **Beautiful UX**: Charm stack (BubbleTea, Lipgloss) for polished terminal UI
- **No Lock-In**: Switch providers anytime with `annotr model`

## Installation

Homebrew (macOS/Linux):

```bash
brew install cloudboy-jh/homebrew-annotr/annotr
```

Scoop (Windows):

```bash
scoop bucket add annotr https://github.com/cloudboy-jh/scoop-annotr
scoop install annotr
```

Install script:

```bash
curl -fsSL https://raw.githubusercontent.com/cloudboy-jh/annotr/main/install.sh | sh
```

Go install:

```bash
go install github.com/cloudboy-jh/annotr/cmd/annotr@latest
```

Build from source:

```bash
git clone https://github.com/cloudboy-jh/annotr.git
cd annotr
make install
```

## Quick Start

```bash
# First-time setup
annotr init

# Add comments to a file
annotr file.go

# Process all files in a directory
annotr ./src
```

Choose local or API providers during setup, and switch anytime with `annotr model`.

## Configuration

Run `annotr init` to configure. It will:

1. Detect Ollama (if installed) and list available models
2. Prompt for an API key if you choose Claude/OpenAI/Groq
3. Let you select a provider, model, and comment style
4. Save to `~/.annotr/config.json`

### Use API providers (no local model required)

If you prefer cloud providers, just run `annotr init` and select Claude, OpenAI, or Groq. You can switch later with `annotr model`.

### Recommended: Install Ollama (free, local)

```bash
curl -fsSL https://ollama.com/install.sh | sh
ollama pull qwen2.5-coder:1.5b
annotr init
```

## Supported Languages

- Go
- Python
- JavaScript
- TypeScript

## Usage

```bash
# Single file
annotr main.go
# → Adds comments, shows: "Enjoy your comments! ;)"

# Directory (interactive)
annotr ./src
# → Prompts for each file: "Process main.go? (y/n)"

# Remove comments from a file or directory
annotr clear file.go
annotr clear ./src

# Switch provider or model
annotr model

# Update models manifest
annotr update-models
```

## License

MIT
