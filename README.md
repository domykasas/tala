# Tala - Terminal AI Language Assistant

![Version](https://img.shields.io/badge/version-0.0.2-blue.svg)
![Go](https://img.shields.io/badge/Go-1.24.4+-00ADD8.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)

Tala is a terminal-based AI language assistant built with Go and Bubble Tea. It provides an interactive interface for communicating with various AI providers including OpenAI, Anthropic, and Ollama, with intelligent file operations and shell command execution capabilities.

## Features

- **Multiple AI Providers**: Support for OpenAI, Anthropic Claude, and Ollama
- **AI-Integrated File Operations**: Natural language file and directory management
  - Create, read, update, delete files with simple commands
  - "Create test.txt with hello world" - and it actually creates the file!
- **Secure Shell Command Execution**: Run bash/shell commands safely
  - Built-in security filtering and timeout protection
  - Cross-platform support (Windows, Linux, macOS)
- **Copy-Paste Friendly**: No alt-screen mode - use your terminal's native copy/paste
- **Real-Time Stats**: Live response times, token counts, and session statistics
- **Local-First**: Defaults to Ollama with deepseek-r1 (no API key required)
- **Simple Configuration**: JSON-based configuration with sensible defaults
- **Terminal Native**: Works with all standard terminal features and shortcuts

## Quick Start

### Prerequisites

- Go 1.24.4 or later
- For default setup: Ollama with deepseek-r1 model

### Installation

```bash
git clone https://github.com/domykasas/tala
cd tala
go mod tidy
go build -o tala
```

### First Run

```bash
# Install and start Ollama (if using default setup)
ollama serve
ollama pull deepseek-r1

# Run Tala
./tala
```

## Configuration

Tala uses a JSON configuration file located at `~/.config/tala/config.json`.

### Default Configuration

```json
{
  "api_key": "",
  "provider": "ollama",
  "model": "deepseek-r1",
  "temperature": 0.7,
  "max_tokens": 1000,
  "system_prompt": "You are a helpful AI assistant."
}
```

### Supported Providers

- **Ollama** (default): Local AI models, no API key required
  - Models: `deepseek-r1`, `llama2`, `mistral`, `codellama`, etc.
- **OpenAI**: GPT models
  - Models: `gpt-3.5-turbo`, `gpt-4`, `gpt-4-turbo`, etc.
- **Anthropic**: Claude models
  - Models: `claude-3-sonnet`, `claude-3-haiku`, `claude-3-opus`, etc.

### Switching Providers

To use OpenAI or Anthropic, edit your config file:

```json
{
  "api_key": "your-api-key-here",
  "provider": "openai",
  "model": "gpt-4",
  "temperature": 0.7,
  "max_tokens": 1000,
  "system_prompt": "You are a helpful AI assistant."
}
```

## Usage

### Basic Commands

- **Enter**: Send message
- **Ctrl+C**: Quit application
- **Ctrl+L**: Clear screen and reset session stats
- **Backspace**: Delete characters from input

### File Operations

Tala understands natural language for file operations:

```
You: Create a file called notes.txt with "Hello World" in it
AI: ✓ Created file 'notes.txt'

You: List files in current directory  
AI: ✓ Listed directory contents

You: Run ls command
AI: ✓ Executed shell command successfully
```

### Direct Commands

You can also use direct slash commands:

- `/help` - Show available commands and file operations
- `/create <filename> <content>` - Create files directly
- `/ls` - List directory contents
- `/pwd` - Show current directory

### Copy and Paste

Tala is designed to work seamlessly with your terminal's copy-paste functionality:

- **Select text**: Use mouse or keyboard selection
- **Copy**: Use your terminal's copy shortcut (Ctrl+Shift+C, Cmd+C, etc.)
- **Paste**: Use your terminal's paste shortcut (Ctrl+Shift+V, Cmd+V, etc.)
- **Scroll back**: Use mouse wheel, PgUp/PgDn, or terminal scrollback

### Statistics

Tala displays helpful statistics:

- **Per-response**: Token count and response time for each AI response
- **Session-wide**: Total requests, total tokens, and average response time
- **Live loading**: Real-time elapsed time while AI is thinking

## Development

### Running Tests

```bash
go test ./...
```

### Project Structure

```
tala/
├── main.go              # Application entry point
├── config/              # Configuration management
│   ├── config.go        # Config struct and file operations
│   └── config_test.go   # Configuration tests
├── ai/                  # AI provider implementations
│   ├── provider.go      # Provider interface and implementations
│   ├── intent.go        # AI-powered intent detection
│   ├── tools.go         # Tool execution framework
│   └── *_test.go        # Comprehensive test suites
├── fileops/             # File system operations
│   ├── fileops.go       # CRUD operations for files/directories
│   ├── commands.go      # Direct command interface
│   └── *_test.go        # File operation tests
├── tui/                 # Terminal UI components
│   └── model.go         # Bubble Tea model implementation
├── go.mod              # Go module file
├── README.md           # This file
├── CHANGELOG.md        # Version history
├── ROADMAP.md          # Development roadmap
└── CLAUDE.md           # Development guide
```

## Troubleshooting

### Common Issues

1. **Ollama connection failed**: Ensure Ollama is running (`ollama serve`)
2. **Model not found**: Pull the model first (`ollama pull deepseek-r1`)
3. **API key errors**: Check your API key in the config file
4. **Permission errors**: Ensure config directory is writable

### Debug Information

Set environment variable for verbose output:
```bash
export DEBUG=1
./tala
```

## Contributing

We welcome contributions! Please see our [development guide](CLAUDE.md) for details on:

- Code style and conventions
- Testing requirements
- Pull request process
- Development setup

## License

MIT License - see LICENSE file for details.

## Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) for the excellent TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) for terminal styling
- [Ollama](https://github.com/ollama/ollama) for local AI model serving
- The Go community for excellent tooling and libraries

---

**Tala** - Making AI conversations as natural as terminal commands.