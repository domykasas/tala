//go:build !gui
// +build !gui

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"tala/internal/ai"
	"tala/internal/config"
	"tala/internal/tui"
)

var version = "1.0.3" // Version can be overridden at build time

func main() {
	// Parse command line flags
	var (
		prompt = flag.String("p", "", "Direct prompt mode - execute prompt and exit")
		model = flag.String("model", "", "Override model for this session")
		provider = flag.String("provider", "", "Override provider for this session")
		help = flag.Bool("help", false, "Show help message")
		versionFlag = flag.Bool("version", false, "Show version information")
	)
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	if *versionFlag {
		showVersion()
		return
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Apply command-line overrides
	if *model != "" {
		cfg.Model = *model
	}
	if *provider != "" {
		cfg.Provider = *provider
	}

	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
		fmt.Fprintf(os.Stderr, "Please set your API key and other configuration options.\n")
		fmt.Fprintf(os.Stderr, "Configuration file location: ~/.config/tala/config.json\n")
		os.Exit(1)
	}

	// Handle direct prompt mode (headless)
	if *prompt != "" {
		runDirectPrompt(*prompt, cfg)
		return
	}

	// Handle remaining args as prompt
	args := flag.Args()
	if len(args) > 0 {
		promptText := strings.Join(args, " ")
		runDirectPrompt(promptText, cfg)
		return
	}

	// Default TUI mode
	simpleTUI, err := tui.NewSimpleTUI(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := simpleTUI.Run(); err != nil {
		log.Fatal(err)
	}
}

// runDirectPrompt executes a single prompt and exits (headless mode)
func runDirectPrompt(prompt string, cfg *config.Config) {
	provider, err := ai.CreateProvider(cfg.Provider, cfg.APIKey, cfg.Model, cfg.Temperature, cfg.MaxTokens)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating provider: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	var response string

	// Use tools if available
	if provider.SupportsTools() {
		response, _, err = provider.GenerateResponseWithTools(ctx, prompt)
	} else {
		response, err = provider.GenerateResponse(ctx, prompt)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Output response directly to stdout (Unix-philosophy)
	fmt.Print(response)
	if !strings.HasSuffix(response, "\n") {
		fmt.Print("\n")
	}
}

// showHelp displays usage information
func showHelp() {
	fmt.Printf(`Tala - Terminal AI Language Assistant

Usage:
  tala [flags] [prompt...]

Flags:
  -p, --prompt string     Direct prompt mode - execute prompt and exit
  --model string          Override model for this session
  --provider string       Override provider for this session
  --help                  Show this help message
  --version               Show version information

Examples:
  tala                           # Interactive mode
  tala "Hello world"             # Direct prompt
  tala -p "Explain Go channels"  # Direct prompt with flag
  tala --model gpt-4 "Help me"   # Override model
  tala --provider openai -p "Hi" # Override provider

Interactive Commands:
  /help                   Show available commands
  /clear                  Clear screen and reset session
  /ls, /cat, /pwd, etc.   File operations
  Ctrl+C                  Exit
  Ctrl+L                  Clear screen

Configuration:
  ~/.config/tala/config.json

For more information, visit: https://github.com/domykasas/tala
`)
}

// showVersion displays version information
func showVersion() {
	fmt.Printf("Tala v%s\n", version)
	fmt.Printf("Terminal AI Language Assistant\n")
	fmt.Printf("Built with Go 1.24.4\n")
}