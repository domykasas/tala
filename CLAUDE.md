# CLAUDE.md - Tala Development Guide

This file contains development information and configuration specific to Tala (Terminal AI Language Assistant).

## Development Commands

### Testing
```bash
go test ./...
```

### Building
```bash
go build -o tala
```

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
│   └── provider_test.go # Provider tests
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
- Max Tokens: `1000`
- System Prompt: `"You are a helpful AI assistant."`

### Configuration Validation
The application validates:
- API key presence (required for OpenAI and Anthropic, not needed for Ollama)
- Provider name (must be supported: ollama, openai, anthropic)
- Model name (must be specified)

## Architecture Decisions

### Why Tala?
**Tala** means "to speak" or "language" in several languages, reflecting our focus on AI language assistance in the terminal.

### Why No Alt-Screen Mode?
- Enables native terminal copy-paste functionality
- Users can scroll back through conversation history
- Works with all terminal features (search, selection, etc.)
- More natural terminal experience

### Why Ollama as Default?
- Local-first approach for privacy
- No API key required for immediate use
- Fast responses for development tasks
- Supports many popular models

### Why Bubble Tea?
- Excellent TUI framework for Go
- Good documentation and community
- Handles terminal complexity well
- Supports concurrent operations

### Provider Interface Pattern
- Allows easy addition of new AI providers
- Enables testing with mock providers
- Provides consistent API across providers
- Supports provider-specific configuration

## Development Guidelines

### Code Style
- Follow Go naming conventions
- Use meaningful variable and function names
- Add comments for exported functions
- Keep functions focused and small
- Use `gofmt` for formatting

### Testing Strategy
- Write tests for all new functionality
- Use table-driven tests where appropriate
- Mock external dependencies (HTTP calls)
- Maintain test coverage above 80%
- Test error conditions and edge cases

### Error Handling
- Use Go's idiomatic error handling
- Provide meaningful error messages
- Handle network timeouts gracefully
- Fail fast for configuration errors
- Log errors appropriately (without exposing secrets)

### Dependencies
- Minimize external dependencies
- Use well-maintained libraries
- Pin dependency versions in go.mod
- Regularly update dependencies for security
- Document any breaking changes

## Security Considerations

### API Key Handling
- Store API keys in user config directory only
- Set appropriate file permissions (0644)
- Never log API keys or include in error messages
- Consider encryption for future versions
- Support environment variable overrides

### Input Validation
- Validate all user inputs
- Sanitize data before sending to providers
- Handle malformed responses gracefully
- Prevent injection attacks
- Limit input sizes appropriately

### Network Security
- Use HTTPS for all external API calls
- Implement request timeouts
- Handle TLS certificate validation
- Consider rate limiting for abuse prevention

## Performance Considerations

### Memory Usage
- Stream large responses when possible
- Limit conversation history in memory
- Clean up unused resources promptly
- Consider garbage collection impact
- Monitor memory usage in long sessions

### Network Operations
- Implement proper timeouts (30s default)
- Handle network errors gracefully
- Consider request retry logic
- Cache responses when appropriate
- Minimize request overhead

### Startup Performance
- Lazy load configurations
- Minimize startup dependencies
- Fast config validation
- Quick error reporting

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
```

## Release Process

### Version Numbering
Following [Semantic Versioning](https://semver.org/):
- `MAJOR.MINOR.PATCH`
- Current: `0.0.1` (initial release)

### Release Checklist
1. Update version in relevant files
2. Update CHANGELOG.md following [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)
3. Run full test suite: `go test ./...`
4. Test with all supported providers
5. Update documentation if needed
6. Create git tag: `git tag v0.0.1`
7. Build release binaries
8. Create GitHub release with notes

### Build for Multiple Platforms
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o tala-linux-amd64

# macOS
GOOS=darwin GOARCH=amd64 go build -o tala-darwin-amd64

# Windows
GOOS=windows GOARCH=amd64 go build -o tala-windows-amd64.exe
```

## Contributing Guidelines

### Before Contributing
1. Read this development guide
2. Check existing issues and roadmap
3. Run tests locally: `go test ./...`
4. Follow coding standards
5. Update documentation as needed

### Pull Request Process
1. Create feature branch from main
2. Write tests for new functionality
3. Ensure all tests pass
4. Update relevant documentation
5. Submit PR with clear description
6. Address review feedback promptly

### Code Review Criteria
- **Functionality**: Does it work as intended?
- **Tests**: Are there adequate tests?
- **Security**: Any security implications?
- **Performance**: Impact on performance?
- **Documentation**: Is documentation updated?
- **Style**: Follows Go conventions?

## Future Development

### Provider Integration Roadmap
- Google Gemini API integration
- Cohere model support
- Hugging Face Inference API
- Local llama.cpp integration
- Streaming response support

### UI/UX Improvements
- Syntax highlighting for code blocks
- Better message formatting
- Custom themes and colors
- Configurable prompt templates

### Advanced Features
- Session persistence
- Conversation export
- File upload and analysis
- Git integration
- Plugin system

### Performance Optimizations
- Response caching
- Concurrent request handling
- Memory usage optimization
- Startup time improvements

---

## Quick Reference

### Essential Commands
```bash
# Development
go mod tidy              # Update dependencies
go test ./...           # Run all tests
go build -o tala        # Build binary

# Configuration
~/.config/tala/config.json  # Config file location

# Running
./tala                  # Start Tala
export DEBUG=1          # Enable debug mode
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