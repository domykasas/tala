//go:build gui
// +build gui

package main

import (
	"fmt"
	"log"
	"os"

	"tala/internal/config"
	"tala/internal/gui"
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

	app, err := gui.NewApp(cfg)
	if err != nil {
		log.Fatal(err)
	}

	app.Run()
}