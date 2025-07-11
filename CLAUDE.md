# CLAUDE.md - Project Memory

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview
**Tala** - Terminal-based AI language assistant with multiple interface modes
- **Language**: Go 1.24.4+
- **UI Frameworks**: Simple terminal I/O (TUI), Fyne (GUI)  
- **Architecture**: Modular internal packages (ai, config, fileops, tui, gui)
- **Build Modes**: TUI (default), GUI
- **Current Version**: 1.0.11

## Development Philosophy

### Communication Principles
- **Direct and Concise**: Provide specific technical context without unnecessary explanations
- **Problem-Focused**: Understand root causes systematically before implementing solutions
- **Minimal Changes**: Make targeted, minimal modifications to achieve goals
- **Documentation-First**: Always update documentation after modifications

### Dependency Philosophy
- **Minimal Dependencies**: Use as few external modules as possible
- **Standard Library First**: Prefer Go standard library over external packages
- **Simple Solutions**: Avoid complex frameworks when basic I/O suffices
- **Copy-Paste Friendly**: Terminal interface must support full history access and copy-paste

### Version Management
- **Semantic Versioning**: MAJOR.MINOR.PATCH (breaking.feature.bugfix)
- **Consistent References**: Update version badges in README.md and documentation
- **Changelog Maintenance**: Document all changes following Keep a Changelog format

## Version Increment Rules

**Every commit must include a version bump, no exceptions**

### Version Increment Triggers:
1. **Documentation changes**: PATCH version (1.0.3 → 1.0.4)
2. **Bug fixes**: PATCH version (1.0.3 → 1.0.4)
3. **New features**: MINOR version (1.0.3 → 1.1.0)
4. **Breaking changes**: MAJOR version (1.0.3 → 2.0.0)

### Version Number Format:
- **Pre-release**: "v1.0.0-rc.N" (testing versions)
- **Stable release**: "v1.0.0" (production releases)
- **Always includes 'v' prefix**

### Documentation Requirements:
Files to update for every version bump:
1. **CLAUDE.md** - Update current version references
2. **README.md** - Update version badge and status section
3. **CHANGELOG.md** - Add new version section with changes
4. **Git tag**: `git tag v1.0.X`

### Key Principles:
- **Proactive versioning**: Increment versions with every meaningful change
- **Manual control**: Human controls git push operations
- **Documentation consistency**: Maintain version alignment across all files
- **Clear categorization**: Properly classify changes as breaking, feature, or fix

## Prerequisites

- Go 1.24.4 or later (as specified in go.mod)
- For default setup: Ollama with llama3.2:1b model
- For OpenAI/Anthropic: Valid API keys

## Development Commands

### Testing
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test ./... -v

# Run specific package tests
go test ./internal/config -v
go test ./internal/ai -v
go test ./internal/fileops -v

# Run tests with coverage
go test ./... -cover

# Run specific test function
go test ./internal/ai -run TestOllamaProvider

# Run tests with race detection
go test -race ./...

# Run tests with timeout
go test -timeout=30s ./...
```

### Building
```bash
# Build TUI version (default)
go build -o tala

# Build GUI version
go build -tags gui -o tala-gui
```

### GUI Interface Features
The GUI mode provides an enhanced graphical interface with all terminal features:

- **Dark Theme**: Professional dark theme with semantic color coding
- **Larger Input Fields**: 600x100px input area with multiline support
- **Emoji Integration**: Consistent emoji usage (🤖 Provider, 👤 User, 🔧 System, ❌ Error)
- **Professional Layout**: Header with provider/model info, enhanced menus
- **Progress Indicators**: Visual progress bars and loading states
- **Concurrent Input**: Queue-based input handling for responsive interaction
- **Session Statistics**: Real-time display of requests, tokens, and timing
- **Paragraph Streaming**: AI responses appear paragraph by paragraph
- **Enhanced Settings**: Larger configuration dialog with emoji labels
- **Comprehensive Help**: Built-in help system with keyboard shortcuts
- **File Operations**: Full slash command support and natural language processing

**Note**: When helping with Tala development, do not build the application or delete the `tala` binary file. The user will build it themselves.

**🚨 CRITICAL**: Always use `go test ./...` for compilation verification, never `go build` unless explicitly requested by the user.

**Important**: Always remember to update ROADMAP.md when making significant changes or releases. The roadmap should reflect current development status and future plans.

**Version Management**: When releasing new versions, remember to update the version badge in README.md from `![Version](https://img.shields.io/badge/version-X.X.X-blue.svg)` to match the new version number.

**Change Documentation**: After making any significant file updates, bug fixes, or feature additions, always update CHANGELOG.md with the changes following the Keep a Changelog format. This ensures proper version history tracking.

**Commit Messages**: Do not include Claude Code attribution or co-authorship lines in commit messages. Keep commits clean and professional without AI assistant signatures.

**🚨 GUI Chat History Widget**: CRITICAL - The chat history widget in `internal/gui/app.go` MUST remain as `widget.Entry` to maintain copy-paste functionality. Do NOT change it back to `widget.Label` or `widget.RichText` as this breaks text selection and copying. The current implementation uses Entry with OnChanged handler for read-only protection while preserving full text selection capabilities. This was specifically requested by users and is essential for GUI usability.

## Bug Tracking and Resolution

### Issue Analysis Process
1. **Systematic Diagnosis**: Carefully analyze error messages and symptoms
2. **Root Cause Investigation**: Understand underlying causes before implementing fixes
3. **Targeted Solutions**: Implement minimal, focused changes that address specific issues
4. **Comprehensive Testing**: Verify fixes across different platforms and scenarios
5. **Documentation Updates**: Record solutions and architectural decisions

### Testing Strategy
- **Priority**: Always use `go test ./...` instead of `go build` for validation
- **Cross-Platform**: Test compilation and functionality across supported environments
- **Systematic Coverage**: Use structured test approaches with race detection and timeouts
- **Regression Prevention**: Maintain comprehensive test coverage for fixed issues

## 🚨 Critical Build Testing Rule

**NEVER USE `go build` - Use `go test ./...` instead**

### Why This Matters:
- **Repository Size**: Executables are large (5-30MB) and bloat git repository
- **Compilation Verification**: `go test ./...` verifies code compiles correctly without generating files
- **User Preference**: Respects user preference for handling builds independently
- **Professional Practice**: Standard workflow for code verification

### Exceptions:
- **User Request**: Only build when explicitly requested by the user
- **Release Process**: Automated CI/CD builds for distribution
- **Debugging**: Specific troubleshooting scenarios where binary is needed

### Standard Workflow:
1. **Always**: `go test ./...` for compilation and functionality verification
2. **Only if requested**: `go build` for specific binary creation
3. **Never**: Proactive building without user request

**Release Workflows**: The project uses a comprehensive workflow architecture inspired by Shario:

**go.yml** (Continuous Integration):
- **Trigger**: Push/PR to main and develop branches
- **Purpose**: Continuous testing and validation on all platforms
- **Platforms**: Ubuntu, Windows, macOS
- **Features**: 
  - Cross-platform testing with proper dependencies
  - TUI and GUI build validation
  - Static analysis, race detection, and coverage
  - Artifact upload for coverage reports

**release.yml** (Release Automation):
- **Trigger**: Tag pushes (`v*`)
- **Purpose**: Multi-platform binary builds and package creation
- **Build Matrix**: 
  - Linux: amd64/arm64 (DEB, RPM, AppImage, Snap)
  - Windows: amd64/arm64 (ZIP archives)
  - macOS: amd64/arm64 (DMG files)
  - FreeBSD: amd64 (raw binaries)
- **Package Types**: Comprehensive packaging for all major distributions
- **Checksums**: SHA256 verification for all artifacts

**snapcraft.yml** (Snap Package Workflow):
- **Trigger**: Tag pushes + manual dispatch option
- **Purpose**: Dedicated Snap package creation with multiple fallback strategies
- **Strategies**: Destructive mode → Multipass → Docker fallbacks
- **Features**:
  - **Snap Testing**: Local installation and version verification
  - **Store Publishing**: Ready for Snap Store integration (requires credentials)
  - **Artifact Management**: Automatic upload to GitHub releases
  - **Quality Assurance**: Comprehensive testing before publication

**Development Workflow**:
```
Development → PR → go.yml (CI testing)
Ready for release → Create tag → Push tag → release.yml + snapcraft.yml → GitHub Release
```

**Manual Release Process**:
1. Create tag: `git tag v1.0.3`
2. Push tag: `git push origin v1.0.3`
3. GitHub Actions automatically:
   - Builds all platform binaries and packages
   - Creates comprehensive release with download tables
   - Generates SHA256 checksums for verification
   - Builds and tests Snap packages

**⚠️ CRITICAL REMINDER FOR CLAUDE**:
- **`git push`** = Pushes commits only (triggers CI workflow)
- **`git push origin v1.0.X`** = Pushes tags (triggers RELEASE workflows)
- **GitHub releases ONLY appear when you push TAGS, not commits!**
- **Always do BOTH**: First `git push` for commits, then `git push origin vX.X.X` for releases
- **Missing step**: If no GitHub release appears, you forgot to push the tag!

**Architecture Benefits**:
- ✅ Comprehensive platform coverage (Linux, Windows, macOS, FreeBSD)
- ✅ Professional packaging (DEB, RPM, AppImage, DMG, ZIP, Snap)
- ✅ Reliable Snap building with fallback strategies
- ✅ Dedicated workflows for specialized packaging
- ✅ Proven approach based on Shario's architecture

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
├── main.go              # TUI application entry point
├── main_gui.go          # GUI application entry point  
├── internal/            # Internal packages
│   ├── config/          # Configuration management
│   │   ├── config.go    # Config struct and file operations
│   │   └── config_test.go # Configuration tests
│   ├── ai/              # AI provider implementations
│   │   ├── provider.go  # Provider interface and implementations
│   │   ├── provider_test.go # Provider tests
│   │   ├── intent.go    # AI-powered intent detection
│   │   ├── intent_test.go # Intent detection tests
│   │   ├── tools.go     # File operation tools for AI
│   │   └── tools_test.go # AI tools tests
│   ├── fileops/         # File system operations
│   │   ├── fileops.go   # File and directory CRUD operations
│   │   ├── commands.go  # Command parsing and execution
│   │   ├── fileops_test.go # File operations tests
│   │   └── commands_test.go # Command tests
│   ├── tui/             # Terminal UI components
│   │   └── simple.go    # Simple terminal implementation
│   └── gui/             # GUI components
│       └── app.go       # Fyne GUI application
├── .github/workflows/   # CI/CD workflows
│   ├── go.yml          # Continuous integration
│   ├── release.yml     # Release automation
│   └── snapcraft.yml   # Snap package workflow
├── go.mod              # Go module file
├── go.sum              # Go module checksums
├── snapcraft.yaml      # Snap package configuration
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
- Model: `llama3.2:1b`
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
- **main.go**: TUI entry point that initializes config and starts simple terminal interface
- **main_gui.go**: GUI entry point for Fyne-based graphical interface
- **internal/config/**: Configuration management with JSON file at `~/.config/tala/config.json`
- **internal/ai/**: Provider interface pattern supporting OpenAI, Anthropic, and Ollama
- **internal/tui/**: Simple terminal interface with standard library I/O and ANSI color support
- **internal/gui/**: Fyne-based graphical interface with chat window and settings dialog
- **internal/fileops/**: File system operations with command parsing and AI tool integration

### Architecture Patterns

#### Build Constraint System
- **TUI Mode**: `//go:build !gui` - Default terminal interface (simple I/O)
- **GUI Mode**: `//go:build gui` - Graphical interface using Fyne framework
- **Conditional Compilation**: Allows platform-specific builds and deployment flexibility
- **Provider Compatibility**: Both modes support the same AI provider interface

#### Event-Driven Architecture
- **Provider Interface**: Pluggable AI providers with consistent API
- **Tool Integration**: AI-powered file operations and system commands
- **Multi-Modal Support**: Seamless switching between TUI, GUI, and headless modes

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

### TUI Implementation
The `tui.SimpleTUI` struct manages terminal interface:
- Uses standard library I/O with bufio.Scanner
- ANSI color support for enhanced readability
- Real-time statistics display
- Copy-paste friendly (no alt-screen mode)
- Signal handling for graceful shutdown

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
   ollama pull llama3.2:1b
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

5. **GUI build failures**
   ```bash
   # Install GUI dependencies (platform-specific)
   # Linux: sudo apt-get install libgl1-mesa-dev xorg-dev
   # macOS: No additional dependencies needed
   # Windows: No additional dependencies needed
   go build -tags gui -o tala-gui
   ```

### Debugging Workflow
1. **Error Analysis**: Examine error messages for specific technical details
2. **Minimal Reproduction**: Create minimal test cases to isolate issues
3. **Systematic Testing**: Use `go test ./...` with race detection and timeouts
4. **Cross-Platform Verification**: Test fixes across supported platforms
5. **Documentation Updates**: Record solutions and architectural decisions

### Debug Mode
Set `DEBUG=1` environment variable for verbose logging:
```bash
export DEBUG=1
./tala
```

### Additional Testing Options
```bash
# Test with coverage and generate HTML report
go test ./... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out

# Benchmark tests
go test -bench=. ./...

# Test with memory profiling
go test -memprofile=mem.prof ./...

# Test with CPU profiling
go test -cpuprofile=cpu.prof ./...
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

**Terminal Controls**
```
Enter               - Send message/command
Ctrl+C              - Quit application
Ctrl+L              - Clear screen and reset session stats
Backspace           - Delete characters from input
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
go test ./...           # Run all tests (preferred over go build)
go test ./... -cover    # Run tests with coverage
go test -race ./...     # Run tests with race detection
go test -timeout=30s ./... # Run tests with timeout

# Building
go build -o tala              # Build TUI version
go build -tags gui -o tala-gui     # Build GUI version
go clean                     # Clean build artifacts

# Configuration
~/.config/tala/config.json  # Config file location

# Running
./tala                  # Start Tala (TUI)
./tala-gui             # Start GUI version
export DEBUG=1          # Enable debug mode
DEBUG=1 ./tala          # Debug mode inline
```

## Release Management

### Asset Naming Convention
- **TUI**: `tala-vX.X.X-platform-arch`
- **GUI**: `tala-vX.X.X-platform-arch-gui`
- **Packages**: Platform-specific formats (DEB, RPM, DMG, etc.)

### Build Matrix
- **Platforms**: Linux, Windows, macOS, FreeBSD
- **Architectures**: amd64, arm64
- **Modes**: TUI (default), GUI
- **Package Formats**: Binary, DEB, RPM, Snap, AppImage, DMG

### Key Files
- `main.go` - TUI application entry point
- `main_gui.go` - GUI application entry point
- `internal/config/config.go` - Configuration management
- `internal/ai/provider.go` - AI provider implementations
- `internal/tui/simple.go` - Terminal UI logic
- `internal/gui/app.go` - Fyne GUI implementation

### Important Environment Variables
- `DEBUG=1` - Enable debug logging
- `HOME` - Used for config directory path

---

**Happy coding with Tala!** 🗣️