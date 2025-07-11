# Tala - Terminal AI Language Assistant

[![Version](https://img.shields.io/badge/version-1.0.11-blue.svg)](https://github.com/domykasas/tala/releases/latest)
[![Go](https://img.shields.io/badge/Go-1.24.4+-00ADD8.svg)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](https://github.com/domykasas/tala?tab=MIT-1-ov-file#MIT-1-ov-file)

**📥 [Download Latest Release](https://github.com/domykasas/tala/releases/latest)**

Tala is a terminal-based AI language assistant built with Go. It provides a colorful, interactive interface for communicating with various AI providers including OpenAI, Anthropic, and Ollama, with intelligent file operations and shell command execution capabilities.

**🎯 Available in multiple formats**: DEB, RPM, Snap, AppImage, DMG, Windows executables, and standalone binaries for Linux, Windows, macOS, and FreeBSD.

## Features

- **Multiple AI Providers**: Support for OpenAI, Anthropic Claude, and Ollama
- **AI-Integrated File Operations**: Natural language file and directory management
  - Create, read, update, delete files with simple commands
  - "Create test.txt with hello world" - and it actually creates the file!
- **Secure Shell Command Execution**: Run bash/shell commands safely
  - Built-in security filtering and timeout protection
  - Cross-platform support (Windows, Linux, macOS)
- **Copy-Paste Friendly**: No alt-screen mode - use your terminal's native copy/paste
- **Colorful Interface**: ANSI color support with semantic color coding for better readability
- **Dual Interface Modes**: Enhanced terminal (TUI) and graphical (GUI) interfaces
- **Real-Time Stats**: Live response times, token counts, and session statistics
- **Concurrent Input**: Type next commands while AI is processing responses
- **Enhanced GUI**: Dark theme, larger input fields, emoji integration, and progress indicators
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

# Build TUI version (default)
go build -o tala

# Build GUI version
go build -tags gui -o tala-gui
```

### First Run

```bash
# Install and start Ollama (if using default setup)
ollama serve
ollama pull deepseek-r1

# Run Terminal Interface (TUI)
./tala

# Run Graphical Interface (GUI)
./tala-gui
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
  "max_tokens": 0,
  "system_prompt": "You are a helpful AI assistant."
}
```

### Configuration Parameters Explained

- **provider**: AI service to use (`ollama`, `openai`, or `anthropic`)
- **model**: Specific AI model name for the chosen provider
- **api_key**: Authentication key (required for OpenAI/Anthropic, not needed for Ollama)
- **temperature**: Response creativity level (0.0-2.0)
  - `0.0`: Very focused, deterministic responses
  - `0.7`: Balanced creativity (recommended)
  - `2.0`: Very creative, varied responses
- **max_tokens**: Maximum response length (`0` = unlimited)
- **system_prompt**: Initial instruction for the AI assistant

### Supported Providers

- **Ollama** (default): Local AI models, no API key required
  - Models: `deepseek-r1`, `llama2`, `mistral`, `codellama`, etc.
  - Setup: Install Ollama and pull desired model
- **OpenAI**: GPT models, requires API key
  - Models: `gpt-3.5-turbo`, `gpt-4`, `gpt-4-turbo`, etc.
  - Setup: Get API key from OpenAI platform
- **Anthropic**: Claude models, requires API key
  - Models: `claude-3-sonnet`, `claude-3-haiku`, `claude-3-opus`, etc.
  - Setup: Get API key from Anthropic console

### Switching Providers

To use OpenAI or Anthropic, edit your config file:

```json
{
  "api_key": "your-api-key-here",
  "provider": "openai",
  "model": "gpt-4",
  "temperature": 0.7,
  "max_tokens": 0,
  "system_prompt": "You are a helpful AI assistant."
}
```

## Usage

### Interface Controls

**Terminal (TUI) Mode:**
- **Enter**: Send message
- **Ctrl+C**: Quit application
- **Ctrl+L**: Clear screen and reset session stats
- **Backspace**: Delete characters from input

**GUI Mode:**
- **Enter**: New line in input field
- **Shift+Enter**: Send message
- **Ctrl+N**: New chat (clear history)
- **Ctrl+Q**: Quit application

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

- [Ollama](https://github.com/ollama/ollama) for local AI model serving
- The Go community for excellent tooling and libraries
- Terminal color standards and ANSI escape sequences for cross-platform compatibility
- The Unix philosophy of simple, composable tools

---

**Tala** - Making AI conversations as natural as terminal commands.