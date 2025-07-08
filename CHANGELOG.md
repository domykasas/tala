# Changelog

All notable changes to Tala will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.0.8] - 2025-07-08

### Security
- **Enhanced file permissions**: Fixed gosec security issues
  - Changed file permissions from 0644 to 0600 (user read/write only)
  - Changed directory permissions from 0755 to 0750 (more restrictive)
  - Added input sanitization for process filtering to prevent command injection
  - Fixed unhandled error in process kill operation
- **Improved security posture**: All gosec high/medium severity issues resolved

### Fixed
- **Command injection prevention**: Sanitized filter input in process listing
- **File security**: More restrictive permissions for config and data files
- **Error handling**: Proper error handling for system operations

## [0.0.7] - 2025-07-08

### Fixed
- **Tests**: Fixed failing test in config/config_test.go
  - Updated test expectation for MaxTokens from 1000 to 0 (unlimited)
  - Aligned test with configuration changes from version 0.0.3

### Technical
- **Test coverage**: All tests now pass successfully
- **CI/CD reliability**: Improved automated testing pipeline stability

## [0.0.6] - 2025-07-08

### Fixed
- **Code quality**: Fixed staticcheck linting errors across codebase
  - Corrected error message capitalization in ai/provider.go
  - Removed unnecessary fmt.Sprintf in ai/tools.go
  - Replaced conditional with unconditional strings.TrimPrefix in fileops/commands.go
  - Removed unused min/max functions in tui/model.go
- **Documentation**: Updated ROADMAP.md with completed version history
- **Development workflow**: Enhanced code maintainability and consistency

### Enhanced
- **Code standards**: Improved adherence to Go best practices
- **Maintainability**: Cleaner, more readable code structure
- **Development process**: Better integration with static analysis tools

## [0.0.5] - 2025-07-08

### Fixed
- **CI/CD workflows**: Fixed import paths and tool configurations
  - Corrected Gosec import path from `securecodewarrior/gosec` to `securego/gosec`
  - Replaced deprecated golint with staticcheck for better linting
  - Simplified security workflow to remove duplicate vulnerability checks
  - Enhanced workflow reliability and tool compatibility

### Enhanced
- **Code quality tools**: Updated to use maintained and supported linters
- **Security scanning**: Improved security workflow with correct tool paths

## [0.0.4] - 2025-07-08

### Added
- **Comprehensive CI/CD workflows**: Enhanced GitHub Actions for better quality assurance
  - **Security scanning**: Added Gosec security scanner, dependency checks, and vulnerability scanning
  - **Multi-platform testing**: Cross-platform builds on Ubuntu, Windows, and macOS
  - **Go version matrix**: Testing across Go 1.21, 1.22, and 1.23
  - **Code coverage**: Race condition detection and coverage reporting with Codecov integration
  - **Automated releases**: Tag-based release workflow with cross-platform binary builds
- **Code quality tools**: Integrated linting, vetting, and static analysis
- **Security automation**: Weekly scheduled security scans and SARIF reporting

### Enhanced
- **CI pipeline reliability**: Improved test stability and comprehensive coverage
- **Development workflow**: Better developer experience with automated quality checks
- **Release process**: Streamlined binary distribution for multiple platforms

## [0.0.3] - 2025-07-08

### Changed
- **Configuration defaults**: Removed default token limit (set to 0 for unlimited)
- **Enhanced documentation**: Significantly improved CLAUDE.md with architectural insights
  - Added dual command architecture explanation
  - Enhanced security model documentation
  - Improved provider implementation guidance
  - Better testing command examples
  - Expanded troubleshooting guide

## [0.0.2] - 2025-07-08

### Added
- **AI-integrated file operations**: Natural language file and directory management
  - Create, read, update, delete files with simple commands
  - Directory operations (create, delete, list contents)
  - Copy and move file operations
- **Shell command execution**: Secure bash/shell command execution
  - Built-in security filtering to block dangerous commands
  - Timeout protection (30-second maximum)
  - Output size limits (10KB) to prevent memory issues
  - Cross-platform support (Windows, Linux, macOS)
- **AI-powered intent detection**: Intelligent command interpretation
  - Replaces hardcoded pattern matching with AI-based understanding
  - Fallback pattern matching when AI detection unavailable
  - Support for natural language requests like "create test.txt with hello world"
- **System information tools**: 
  - Get working directory and change directories
  - List running processes with optional filtering
  - Display system information (OS, architecture, CPU count)
- **Direct command interface**: Slash commands for immediate operations
  - `/help` - Show available commands
  - `/create <filename> <content>` - Create files directly
  - `/ls` - List directory contents
  - `/pwd` - Show current directory

### Enhanced
- **Improved error handling**: Better fallback responses when AI providers unavailable
- **Enhanced provider interface**: All providers now support tool execution
- **Comprehensive test coverage**: Tests for all new file operations and AI integration
- **Security-first design**: Extensive command filtering and safety measures

### Security
- **Command execution security**: Whitelist-based approach blocking dangerous operations
- **Input validation**: Protection against command injection and directory traversal
- **Resource protection**: Timeout and output size limits prevent abuse
- **Safe command patterns**: Only allow vetted, safe shell commands

### Technical
- **New packages**: 
  - `fileops/` - File system operations with structured error handling
  - `ai/intent.go` - AI-powered intent detection system
  - `ai/tools.go` - Tool execution framework
- **Enhanced architecture**: Tool-based system for extensible AI capabilities
- **Improved testing**: Added comprehensive test suites for new functionality

## [0.0.1] - 2025-07-08

### Added
- Initial release of Tala (Terminal AI Language Assistant)
- Support for multiple AI providers:
  - **Ollama** (default): Local AI models with deepseek-r1 as default
  - **OpenAI**: GPT models (gpt-3.5-turbo, gpt-4, etc.)
  - **Anthropic**: Claude models (claude-3-sonnet, claude-3-haiku, etc.)
- Copy-paste friendly terminal interface (no alt-screen mode)
- Real-time statistics tracking:
  - Per-response token counts and timing
  - Session-wide statistics (total requests, tokens, average response time)
  - Live elapsed time display during AI processing
- JSON-based configuration system with validation
- Interactive terminal UI built with Bubble Tea framework
- Comprehensive test suite for core functionality
- Configuration file at `~/.config/tala/config.json`
- Keyboard shortcuts:
  - `Enter`: Send message
  - `Ctrl+C`: Quit application
  - `Ctrl+L`: Clear screen and reset session stats
  - `Backspace`: Delete input characters
- HTTP client integration for Ollama API
- Error handling and graceful degradation
- Project documentation (README, ROADMAP, development guide)

### Security
- Configuration files created with secure permissions (0644)
- API keys stored in user configuration directory
- No sensitive data logged or exposed

---

## Release Notes

### Version 0.0.1 - Initial Release

This is the first public release of Tala, a terminal-based AI language assistant designed for developers and power users who want a simple, efficient way to interact with AI models directly from their terminal.

**Key Highlights:**
- **Local-first approach**: Defaults to Ollama with deepseek-r1, no cloud API required
- **Terminal-native**: Full copy-paste support, works with all terminal features
- **Multi-provider**: Easy switching between Ollama, OpenAI, and Anthropic
- **Real-time feedback**: Live statistics and response timing
- **Zero configuration**: Works out of the box with sensible defaults

**Getting Started:**
1. Install Ollama and pull deepseek-r1 model
2. Build Tala: `go build -o tala`
3. Run: `./tala`

For detailed installation and usage instructions, see [README.md](README.md).