# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

Tala is a terminal-based AI language assistant built with Go and Bubble Tea. It provides an interactive interface for communicating with various AI providers including OpenAI, Anthropic, and Ollama.

## Development Commands

### Testing
```bash
go test ./...
```

### Building
```bash
go build -o tala
```

**Note**: When helping with Tala development, do not build the application or delete the `tala` binary file. The user will build it themselves.

### Dependencies
```bash
go mod tidy
```

### Running with Debug
```bash
export DEBUG=1
./tala
```

## Project Structure

```
tala/
├── main.go              # Application entry point
├── config/              # Configuration management
│   ├── config.go        # Config struct and file operations
│   └── config_test.go   # Configuration tests
├── ai/                  # AI provider implementations
│   ├── provider.go      # Provider interface and implementations
│   ├── provider_test.go # Provider tests
│   ├── tools.go         # File operation tools for AI
│   └── tools_test.go    # AI tools tests
├── fileops/             # File system operations
│   ├── fileops.go       # File and directory CRUD operations
│   ├── commands.go      # Command parsing and execution
│   ├── fileops_test.go  # File operations tests
│   └── commands_test.go # Command tests
├── tui/                 # Terminal UI components
│   └── model.go         # Bubble Tea model implementation
├── go.mod              # Go module file
├── go.sum              # Go module checksums
├── README.md           # Main documentation
├── CHANGELOG.md        # Version history (Keep a Changelog format)
├── ROADMAP.md          # Development roadmap
└── CLAUDE.md           # This file
```

## Configuration

### Default Configuration Location
`~/.config/tala/config.json`

### Default Configuration Values
- Provider: `ollama`
- Model: `deepseek-r1`
- Temperature: `0.7`
- Max Tokens: `0` (no limit)
- System Prompt: `"You are a helpful AI assistant."`

### Configuration Validation
The application validates:
- API key presence (required for OpenAI and Anthropic, not needed for Ollama)
- Provider name (must be supported: ollama, openai, anthropic)
- Model name (must be specified)

## Architecture Overview

### Core Components
- **main.go**: Entry point that initializes config and starts the TUI
- **config/**: Configuration management with JSON file at `~/.config/tala/config.json`
- **ai/**: Provider interface pattern supporting OpenAI, Anthropic, and Ollama
- **tui/**: Bubble Tea-based terminal interface with no alt-screen mode for copy-paste functionality

### Provider System
The `ai.Provider` interface allows pluggable AI providers:
- `OpenAIProvider`: OpenAI API integration (GPT models)
- `AnthropicProvider`: Anthropic API integration (Claude models)  
- `OllamaProvider`: Local Ollama integration with HTTP API calls

### Configuration Flow
1. Load config from `~/.config/tala/config.json` or create default
2. Validate provider requirements (API keys for OpenAI/Anthropic, not needed for Ollama)
3. Create provider instance based on config
4. Initialize TUI with provider and config

### Enhanced Operations System
Tala includes comprehensive operations capabilities:

**1. Direct Commands (User-initiated)**
- **Command Detection**: Input starting with `/` triggers file operations instead of AI chat
- **Operation Types**: Create, read, update, delete files and directories, plus navigation (ls, cd, pwd)
- **Safety**: Operations are restricted to current working directory and subdirectories
- **Feedback**: Color-coded responses (green for success, red for errors)

**2. AI-Integrated Tools (AI-initiated)**
- **Intent Detection**: AI-powered natural language understanding of user requests
- **File Operations**: Create, read, update, delete files and directories
- **Shell Commands**: Execute safe shell/bash commands with security restrictions
- **System Information**: Get system details, process lists, working directory
- **Security**: Comprehensive command filtering and timeout protection

## Key Implementation Details

### TUI Model Structure
The `tui.Model` struct manages application state:
- `input`: Current user input string
- `loading`: Boolean for request state
- `provider`: AI provider instance
- `config`: Configuration reference
- Statistics tracking (tokens, requests, timing)

### Message Flow
1. User types input and presses Enter
2. Input displayed with "You:" prefix
3. **Branch A - File Operations** (if input starts with `/`):
   - Parse command and arguments
   - Execute file operation via `fileops.ExecuteCommand()`
   - Display result with "System:" prefix (green/red based on success)
4. **Branch B - AI Chat** (normal input):
   - Loading state shows with elapsed time
   - Provider generates response via `GenerateResponse()`
   - Response displayed with "AI:" prefix and stats
   - Statistics updated for session totals

### Provider Interface Requirements
All providers must implement:
```go
type Provider interface {
    GenerateResponse(ctx context.Context, prompt string) (string, error)
    GenerateResponseWithTools(ctx context.Context, prompt string) (string, []ToolResult, error)
    GetName() string
    SupportsTools() bool
}
```

### AI Tool System
- **Intent Detection**: `IntentDetector` class uses AI to understand user requests
- **Tool Interface**: Standardized tool execution via `ExecuteTool()` function
- **Available Tools**: File operations, shell commands, system information, process management
- **Security**: Command filtering, timeout protection, output size limits
- **Provider Integration**: All providers support tool calling through enhanced interface
- **Fallback System**: Pattern matching fallback when AI intent detection fails

#### Security Model for AI Tools
- **Command Blacklist**: Dangerous commands (`rm -rf`, `sudo`, `chmod 777`) are blocked
- **Timeout Protection**: All shell commands timeout after 30 seconds maximum
- **Output Limits**: Command output truncated to prevent memory exhaustion
- **Path Validation**: File operations restricted to working directory and subdirectories
- **Process Cleanup**: Proper cleanup of timed-out or failed processes

### Configuration Management
- Auto-creates config file with defaults if missing
- Validates provider-specific requirements
- Uses JSON marshaling for persistence
- File permissions set to 0644 for security

### Dual Command Architecture
Tala implements a unique dual-interface system:

**1. Direct Commands (Immediate Execution)**
- **Trigger**: Input prefixed with `/` (e.g., `/ls`, `/cat file.txt`)
- **Processing**: Bypasses AI, directly executes via `fileops.ExecuteCommand()`
- **Response**: Immediate system response with color coding
- **Use Case**: Quick file operations without AI interpretation overhead

**2. Natural Language Commands (AI-Mediated)**
- **Trigger**: Regular input without `/` prefix
- **Processing**: AI interprets intent and decides whether to use tools
- **Response**: AI-generated response enhanced with tool execution results
- **Use Case**: Complex requests requiring context understanding

**Decision Flow**:
```
User Input → Starts with '/'? → YES → Direct Command Execution
                              → NO  → AI Intent Detection → Tool Selection → AI Response
```

**Example Patterns**:
- `/ls` vs "show me all files" - both list files, different execution paths
- `/cat config.json` vs "what's in the config file?" - both read files
- Complex: "create a backup of main.go and show me the differences" - only works via AI

## Troubleshooting Guide

### Common Development Issues

1. **Module import errors**
   ```bash
   go mod tidy
   ```

2. **Ollama connection failed**
   ```bash
   ollama serve
   ollama pull deepseek-r1
   ```

3. **Config file permissions**
   ```bash
   chmod 644 ~/.config/tala/config.json
   ```

4. **Build errors**
   ```bash
   go clean
   go mod tidy
   go build -o tala
   ```

### Debug Mode
Set `DEBUG=1` environment variable for verbose logging:
```bash
export DEBUG=1
./tala
```

### Testing Individual Components
```bash
# Test specific package
go test ./config -v

# Test with coverage
go test ./... -cover

# Run specific test
go test ./ai -run TestOllamaProvider

# Run tests with race detection
go test -race ./...

# Run tests with timeout
go test -timeout=30s ./...
```

## Adding New Providers

To add a new AI provider:

1. **Create Provider Struct**: Define struct with required fields (API key, model, etc.)
2. **Implement Interface**: Add `GenerateResponse()` and `GetName()` methods
3. **Update Factory**: Add case to `CreateProvider()` function in `ai/provider.go`
4. **Add Tests**: Create test file with provider-specific test cases
5. **Update Config**: Add provider to validation in `config/config.go`

### Provider Implementation Template
```go
type NewProvider struct {
    APIKey      string
    Model       string
    Temperature float64
    MaxTokens   int
}

func (p *NewProvider) GenerateResponse(ctx context.Context, prompt string) (string, error) {
    // Implement API call logic
}

func (p *NewProvider) GenerateResponseWithTools(ctx context.Context, prompt string) (string, []ToolResult, error) {
    // Implement tool-aware API call logic
    // Return both response and tool execution results
}

func (p *NewProvider) GetName() string {
    return "NewProvider"
}

func (p *NewProvider) SupportsTools() bool {
    return true // or false based on provider capabilities
}
```

## Testing Strategy

### Unit Tests
- Mock HTTP clients for provider tests
- Test configuration loading and validation
- Test TUI state transitions
- Use table-driven tests for multiple scenarios
- Temporary directory isolation for file operation tests

### Integration Tests
- Test with actual Ollama instance when available
- Validate end-to-end message flow
- Test error handling and recovery
- AI tool execution with real shell commands (in safe test environment)

### Test Coverage
Current coverage focuses on:
- Configuration management (`config/config_test.go`)
- Provider implementations (`ai/provider_test.go`)
- AI tool system (`ai/tools_test.go`)
- File operations (`fileops/fileops_test.go`)
- Command parsing (`fileops/commands_test.go`)
- Run `go test ./... -cover` to check coverage

### Available File Commands

**Direct Commands (prefix with `/`)**
```
/help               - Show all available commands
/ls [path]          - List files and directories
/cat <file>         - Display file content
/create <file> [content] - Create new file with optional content
/write <file> <content>  - Write content to file (create or overwrite)
/update <file> <content> - Update existing file content
/rm <file>          - Remove file
/mkdir <dir>        - Create directory
/rmdir <dir>        - Remove directory (and contents)
/cp <src> <dst>     - Copy file
/mv <src> <dst>     - Move/rename file
/pwd                - Show current directory
/cd <path>          - Change directory
```

**AI Natural Language Examples**
```
"create a file called hello.txt with Hello World content"
"list all files in the current directory"
"show me what's in the config.json file"
"make a new folder called projects"
"what is my current working directory?"
"delete the old_file.txt"
"copy main.go to backup.go"
```

---

## Quick Reference

### Essential Commands
```bash
# Development
go mod tidy              # Update dependencies
go test ./...           # Run all tests
go test ./... -cover    # Run tests with coverage
go build -o tala        # Build binary
go clean                # Clean build artifacts

# Configuration
~/.config/tala/config.json  # Config file location

# Running
./tala                  # Start Tala
export DEBUG=1          # Enable debug mode
DEBUG=1 ./tala          # Debug mode inline
```

### Key Files
- `main.go` - Application entry point
- `config/config.go` - Configuration management
- `ai/provider.go` - AI provider implementations
- `tui/model.go` - Terminal UI logic

### Important Environment Variables
- `DEBUG=1` - Enable debug logging
- `HOME` - Used for config directory path

---

**Happy coding with Tala!** 🗣️