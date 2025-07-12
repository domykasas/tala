# Changelog

All notable changes to Tala will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.13] - 2025-07-11

### Fixed
- **Cross-Platform Compatibility**: Fixed multiple platform-specific issues in GitHub Actions workflows
  - **macOS Build**: Fixed checksum command error by using `shasum -a 256` instead of `sha256sum` on macOS
  - **Windows Build**: Added `shell: bash` directive to all build steps to prevent PowerShell syntax errors
  - **Linux Build**: Removed failed apt-get snapcraft installation (handled separately in snapcraft workflow)
  - **Platform Detection**: Enhanced platform-specific command selection for better cross-platform support

### Enhanced
- **Build Reliability**: Significantly improved workflow success rate across all supported platforms (Linux, Windows, macOS, FreeBSD)

## [1.0.12] - 2025-07-11

### Fixed
- **Windows Build**: Fixed PowerShell syntax error in GitHub Actions Windows build workflow
  - Added explicit `shell: bash` directive to checksum generation step
  - Resolved PowerShell vs Bash syntax incompatibility on Windows runners
  - Fixed "Missing '(' after 'if' in if statement" error in Windows environment
  - Ensured cross-platform compatibility for all shell commands

### Enhanced
- **Cross-Platform Build**: Improved workflow reliability across all supported platforms (Linux, Windows, macOS)

## [1.0.11] - 2025-07-11

### Fixed
- **Release Workflow**: Fixed checksum generation failures in GitHub Actions release workflow
  - Added proper file existence checks before creating checksums
  - Improved error handling for missing artifacts
  - Added debug logging to identify file availability issues
  - Enhanced checksum creation for package files (DEB, RPM)
  - Prevented build failures when optional files don't exist

### Enhanced
- **Build Reliability**: Improved release workflow robustness with better artifact handling

## [1.0.10] - 2025-07-11

### Fixed
- **Snapcraft Installation**: Fixed persistent snapcraft installation failures in GitHub Actions
  - Removed apt-get attempt to install snapcraft (not available in Ubuntu repos)
  - Improved snap daemon initialization with systemctl daemon-reload
  - Added core snap installation (required for classic snaps)
  - Enhanced snapd socket readiness with proper symbolic link creation
  - Increased retry attempts from 3 to 5 with longer wait times
  - Added squashfs-tools dependency for snap building
  - Improved service startup sequence with proper systemctl commands

### Enhanced
- **Snap Build Reliability**: Significantly improved snap package build success rate with robust installation process

## [1.0.9] - 2025-07-11

### Fixed
- **GitHub Actions Permissions**: Fixed 403 release creation errors in GitHub Actions workflows
  - Added `contents: write` permission to release.yml workflow
  - Added `contents: write` permission to snapcraft.yml workflow
  - Resolved GitHub token permission issues preventing release creation
  - Enabled proper GitHub release automation for tagged versions

### Enhanced
- **Release Automation**: Improved workflow reliability for consistent release creation

## [1.0.8] - 2025-07-11

### Fixed
- **Snapcraft Workflow**: Fixed snapcraft installation error in GitHub Actions
  - Added proper systemd service start/enable for snapd
  - Implemented retry logic for snapcraft installation
  - Added service readiness waiting with sleep
  - Enhanced error handling and verification
  - Made multipass installation optional with graceful failure

### Enhanced
- **CI/CD Reliability**: Improved snap build workflow robustness for consistent package creation

## [1.0.7] - 2025-07-11

### Updated
- **Project Documentation**: Updated CLAUDE.md with accurate project structure and implementation details
  - Corrected TUI implementation from "Bubble Tea" to "Simple terminal I/O" 
  - Updated default model from "deepseek-r1" to "llama3.2:1b"
  - Added missing files and proper directory structure
  - Fixed architecture descriptions to match actual codebase

### Enhanced
- **Development Guidance**: Improved CLAUDE.md accuracy for better development workflow

## [1.0.6] - 2025-07-11

### Fixed
- **Release Automation**: Fixed missing GitHub releases by properly creating and pushing version tags
- **Workflow Documentation**: Clarified that releases require tag pushes, not just commit pushes

### Enhanced
- **Release Process**: Improved documentation for proper release creation workflow

## [1.0.5] - 2025-07-11

### Fixed
- **Snapcraft Installation**: Fixed snapcraft package installation error in GitHub Actions
  - Added proper snapd installation before snapcraft
  - Improved snap build fallback strategies with better error handling
  - Enhanced debugging output for snap build process
  - Made multipass installation optional with graceful failure handling

### Enhanced
- **Snap Workflow Reliability**: Improved snap building process with multiple fallback strategies
  - Strategy 1: Destructive mode (fastest for CI)
  - Strategy 2: Default snapcraft build (automatic backend selection)
  - Strategy 3: LXD backend (most isolated)
  - Better error reporting and directory listing for debugging

## [1.0.4] - 2025-07-11

### Added
- **Dedicated Snap Workflow**: Added `.github/workflows/snapcraft.yml` following Shario's proven approach
  - Multi-strategy fallback system: destructive mode → multipass → docker
  - Comprehensive testing with local snap installation and version verification
  - Automatic release attachment for snap packages
  - Support for both tag pushes and manual dispatch triggers

### Fixed
- **Release Workflow Triggers**: Updated tag patterns to `v*.*.*` and `v*.*.*-*` for better compatibility
- **Release Automation**: Fixed GitHub releases not appearing automatically for version tags

### Enhanced
- **Workflow Architecture**: Now includes comprehensive 3-workflow system (CI + Release + Snap)
- **Documentation**: Updated CLAUDE.md with complete workflow architecture details

## [1.0.3] - 2025-07-11

### Fixed
- **Staticcheck Error**: Removed unused `currentTokens` field from SimpleTUI struct
- **Workflow Simplification**: Replaced complex 4-workflow system with Shario's proven 2-workflow approach
  - Simplified CI/CD maintenance with go.yml (CI) and release.yml (releases)
  - Removed unnecessary pre-release.yml and snapcraft.yml workflows
  - Enhanced build matrix for comprehensive multi-platform support
  - Proper CGO handling for GUI builds with platform-specific dependencies

### Enhanced
- **GUI Build System**: Fixed build constraints and CGO dependencies for proper GUI compilation
- **Workflow Reliability**: Adopted Shario's proven approach for better build success rates
- **Documentation**: Updated CLAUDE.md with simplified workflow architecture

### Added
- **Comprehensive Release System**: Complete multi-platform packaging inspired by Shario
  - **Cross-Platform Builds**: Linux (amd64/arm64), Windows (amd64/arm64), macOS (Intel/Apple Silicon), FreeBSD (amd64)
  - **Package Formats**: DEB, RPM, Snap, AppImage, DMG, tar.xz archives, and raw binaries
  - **Professional Packaging**: Proper control files, app bundles, and desktop integration
  - **Automated Checksums**: SHA256 verification for all packages and binaries
  - **Comprehensive Testing**: Cross-compilation validation, security scanning, and coverage reporting
- **Advanced CI/CD Workflows**: 
  - **Release Workflow**: Matrix-based builds with artifact collection and GitHub release creation
  - **Snapcraft Workflow**: Dedicated snap building with multiple fallback strategies
  - **Enhanced Testing**: Multi-OS testing, security scanning, linting, and coverage reporting
  - **Build Validation**: Cross-compilation tests for all target platforms
- **Colorful Terminal Interface**: Comprehensive ANSI color support throughout SimpleTUI
  - Color-coded user prompts, AI responses, and system messages
  - Enhanced statistics display with colorful formatting
  - Improved help and configuration displays with semantic color coding
  - Better visual hierarchy and readability in terminal output
- **Enhanced GUI Interface**: Complete redesign to match terminal improvements
  - **Colorful Theme**: Dark theme with semantic color coding matching terminal
  - **Larger Input Fields**: Expanded input area (600x100px) with better placeholder text
  - **Emoji Integration**: Consistent emoji usage throughout interface (🤖, 👤, 🔧, ❌)
  - **Professional Layout**: Header with provider/model info, enhanced menus, progress indicators
  - **Concurrent Input**: Queue-based input handling allowing typing while AI processes
  - **Session Statistics**: Real-time stats display with request/token/timing information
  - **Paragraph Streaming**: AI responses appear paragraph by paragraph with natural timing
  - **Enhanced Settings**: Larger configuration dialog with emoji labels and better validation

### Enhanced
- **Paragraph-Based Streaming Display**: Natural response flow with paragraph-by-paragraph updates
  - AI responses appear paragraph by paragraph with natural timing (200ms delays)
  - Maintains stable cursor positioning without jumping
  - Better reading experience with visual paragraph breaks
  - Proper word wrapping for terminal width compatibility
- **Concurrent Input Handling**: Improved non-blocking input system
  - Users can type next commands while AI is processing responses
  - Stable cursor positioning without jumping during AI thinking
  - Clean progress indicators with live session statistics
- **Release Architecture**: Professional-grade release management
  - **Dual Workflow System**: Separate release and snapcraft workflows for reliability
  - **Artifact Management**: Centralized artifact collection and distribution
  - **Version Management**: Automatic version injection and changelog extraction
  - **Quality Assurance**: Comprehensive testing and validation before release

### Fixed
- **Cursor Stability**: Eliminated cursor jumping during AI response generation
- **Display Interference**: Removed live streaming updates that caused terminal disruption
- **Input Responsiveness**: Improved concurrent input handling for better user experience
- **Progress Display Stability**: Fixed thinking progress to use consistent character width formatting
  - Elapsed time formatted as fixed-width (X.Xs format)
  - Session statistics always visible with consistent digit padding (3-digit requests, 5-digit tokens)
  - Fixed line clearing to prevent multiple progress lines appearing on same line
  - Proper use of `\r\033[K` to clear entire line before updating progress
  - Fixed session stats visibility from the first request onwards
- **GUI Interface Issues**: Resolved layout and formatting problems
  - Removed duplicate "Clear Chat" button that was overlaying text
  - Fixed message formatting to ensure each message starts on a new line
  - Significantly enlarged settings input fields from 300x40px to 800x60px for excellent usability
  - Replaced temperature slider with input field for precise value entry
  - Added helpful labels with value ranges (Temperature: 0.0-2.0, Max Tokens: 0=unlimited)
  - Improved message spacing and line breaks for better readability
  - Custom settings dialog with larger dimensions (900x400px) to accommodate bigger input fields
- **GUI Text Selection and Copy-Paste**: Enhanced text interaction capabilities
  - Replaced RichText with Entry widget for better text selection and copy-paste functionality
  - Added clear separator lines (=================================================================) between messages
  - Improved message formatting with proper line breaks and spacing
  - Enabled full text selection of conversations for copying
  - Fixed AI response streaming to display properly with line breaks
  - Made chat history read-only but fully selectable
  - Implemented proper read-only protection while maintaining copy-paste functionality
  - Enhanced text selection with Entry widget instead of Label for better accessibility
- **GUI Readability and Compatibility**: Improved text rendering and display
  - Removed emoji characters that don't render properly on all systems
  - Enhanced text contrast with custom theme and widget styling
  - Clean message prefixes (USER, AI, SYSTEM, ERROR) instead of emojis
  - Fixed font rendering issues and improved overall readability
  - Removed problematic Unicode characters from all UI elements
  - Switched to Label widget for better text color visibility
  - Custom theme overrides for better foreground text contrast
- **GUI Keyboard Behavior and UX**: Enhanced user interaction patterns
  - Fixed keyboard behavior: Enter for new lines, Shift+Enter to send messages
  - Updated all placeholder text and help documentation to reflect correct shortcuts
  - Removed problematic emojis from buttons and menu items for better compatibility
  - Enhanced settings dialog with detailed explanations for each configuration option
  - Added comprehensive configuration parameter explanations in both GUI and documentation

### Technical
- **Build System**: Adopted Shario's proven multi-platform build approach
  - **CGO Handling**: Disabled for cross-compilation, enabled for native testing
  - **Package Creation**: Automated creation of all major package formats
  - **Fallback Strategies**: Multiple build strategies for challenging platforms like Snap
  - **Quality Control**: Comprehensive testing and validation at every stage

## [1.0.2] - 2025-07-11

### Changed
- **Simplified TUI Interface**: Replaced Bubble Tea framework with basic terminal I/O
  - Removed external UI framework dependency for cleaner, simpler implementation
  - Full copy-paste functionality now available - users can scroll and copy entire conversation history
  - No more rendering conflicts or display bugs from complex TUI frameworks
  - Direct terminal output allows natural terminal behavior (scrolling, selecting, copying)
- **Dependency Minimization**: Reduced external dependencies following new development philosophy
  - Removed Bubble Tea and related charmbracelet dependencies
  - Only essential dependencies remain: Fyne for GUI mode
  - Simpler build process and reduced binary size
- **Enhanced User Experience**: 
  - Clean, readable terminal output without formatting artifacts
  - Word wrapping that respects terminal width
  - Natural terminal interaction with standard copy-paste support
  - Preserved all functionality while improving usability

### Added
- **Development Philosophy**: Added dependency minimization principles to CLAUDE.md
  - Prefer Go standard library over external packages
  - Avoid complex frameworks when basic I/O suffices
  - Ensure copy-paste friendly terminal interfaces

### Fixed
- **Terminal Display Issues**: Eliminated formatting conflicts and misaligned text
- **Copy-Paste Functionality**: Full conversation history now accessible via standard terminal selection
- **Word Wrapping**: Proper text wrapping without UI framework interference

## [1.0.1] - 2025-07-11

### Added
- **Multi-Mode Architecture**: Implemented Shario-inspired best practices
  - TUI Mode: Default terminal interface using Bubble Tea
  - GUI Mode: Graphical interface using Fyne framework
  - Build constraints for conditional compilation
- **Enhanced Project Structure**: Reorganized codebase following Go best practices
  - Moved packages to `/internal/` directory
  - Clean separation between UI layers and core logic
  - Provider interface abstraction for AI backends
- **Comprehensive Build System**: Updated GitHub Actions workflow
  - Cross-platform builds for TUI and GUI modes
  - Streamlined asset naming convention
  - Removed redundant headless mode
- **Documentation Improvements**: Enhanced CLAUDE.md with
  - Development philosophy and communication principles
  - Version management rules and increment triggers
  - Critical build testing rule (go test vs go build)
  - Bug tracking and debugging workflows
  - Release management best practices

### Changed
- **Build Constraints**: Updated from `!headless && !gui` to `!gui` for TUI mode
- **Architecture**: Unified provider system across all interface modes
- **Dependencies**: Added Fyne v2.4.5 for GUI functionality

### Removed
- **Headless Mode**: Eliminated redundant headless executable
  - No unique functionality beyond TUI mode
  - Reduced build complexity and asset duplication
  - Simplified deployment options

### Fixed
- **Import Paths**: Corrected all imports to use internal package structure
- **Provider Interface**: Added backward compatibility for config-based provider creation
- **Build System**: Resolved build constraint conflicts between modes

## [1.0.0] - 2025-07-08

### Added
- **Comprehensive Release Packaging**: Multi-format release system supporting:
  - Linux: DEB, RPM, AppImage, Snap, tar.xz archives, and raw binaries
  - Windows: Executables with blockmaps and Squirrel auto-updater packages
  - macOS: DMG disk images with blockmaps and raw binaries
  - FreeBSD: Native binaries
- **Professional Release Workflows**: Dual-workflow architecture for releases and pre-releases
- **Snap Package Support**: Dedicated snap building and publishing workflow
- **Cross-Platform Distribution**: Comprehensive package formats for all major platforms
- **Release Documentation**: Professional release notes with organized download tables
- **Version Management**: Automated version injection into all package formats

### Enhanced
- **GitHub Workflows**: Robust CI/CD with graceful error handling for package creation
- **Release Notes**: Detailed platform-specific download instructions and verification
- **Package Metadata**: Proper control files, app bundles, and desktop integration
- **Security Features**: SHA256 blockmaps for download verification
- **Documentation**: Comprehensive development guides and release process documentation

### Fixed
- **Workflow Reliability**: Improved error handling for packaging tools in CI environments
- **Package Permissions**: Correct file permissions and executable flags across platforms
- **Cross-Platform Compatibility**: Verified builds for multiple architectures

### Technical
- **Release Architecture**: Separated snap packaging into dedicated workflow for reliability
- **Package Formats**: Professional-grade packaging with proper metadata and signatures
- **Automated Publishing**: GitHub Actions integration for streamlined release process

## [0.1.1] - 2025-07-08

### Enhanced
- **Release Workflows**: Improved dual-workflow architecture for release management
  - Standardized pre-release format to use `-rc.X` (with dot) for consistency
  - Enhanced pre-release workflow to handle both `-rc` and `-rc.` formats automatically
  - Added comprehensive documentation for release process in CLAUDE.md
- **Documentation**: Enhanced CLAUDE.md with detailed workflow explanations
  - Added manual release process instructions
  - Documented preferred RC tag format and tagging reminders
  - Improved architecture benefits and development workflow documentation

### Fixed
- **Version Management**: Updated README.md version badge to reflect current version
- **Workflow Logic**: Pre-release workflow now properly converts between RC formats

## [0.0.9] - 2025-07-08

### Fixed
- **Cross-platform testing**: Fixed path resolution test failure on macOS
  - Added symlink resolution for directory comparison in tests
  - Improved cross-platform compatibility for path operations
  - Ensures tests pass consistently across macOS, Linux, and Windows

### Technical
- **Test reliability**: Enhanced test robustness for different filesystem configurations
- **Path handling**: Better handling of symbolic links in test environments

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