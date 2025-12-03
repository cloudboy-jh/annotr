# annotr

Fast local code commenting CLI. Add AI-generated comments to your code in under 2 seconds.

## Features

- **Local & Fast**: Uses Ollama with qwen2.5-coder:1.5b for 1-2 second generation
- **Zero Cost**: No API fees, runs entirely on your machine
- **Smart Context**: Tree-sitter parsing provides accurate code structure awareness
- **Beautiful UX**: Charm stack (BubbleTea, Lipgloss) for polished terminal UI
- **Multi-Provider**: Supports Ollama, Claude, OpenAI, and Groq

## Installation

```bash
curl -fsSL https://raw.githubusercontent.com/cloudboy-jh/annotr/main/install.sh | sh
```

Or with Go:

```bash
go install github.com/cloudboy-jh/annotr/cmd/annotr@latest
```

Or build from source:

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

## Configuration

Run `annotr init` to configure. It will:

1. Detect Ollama (if installed) and list available models
2. Or prompt for an API key (Claude/OpenAI/Groq)
3. Let you select a model and comment style
4. Save to `~/.annotr/config.json`

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

# Update models manifest
annotr update-models
```

## License

MIT
