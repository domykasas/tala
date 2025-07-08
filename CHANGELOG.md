# Changelog

All notable changes to Tala will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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