package main

import (
	"fmt"
	"log"
	"os"

	"tala/config"
	"tala/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
		fmt.Fprintf(os.Stderr, "Please set your API key and other configuration options.\n")
		fmt.Fprintf(os.Stderr, "Configuration file location: ~/.config/tala/config.json\n")
		os.Exit(1)
	}

	model, err := tui.NewModel(cfg)
	if err != nil {
		log.Fatal(err)
	}

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}