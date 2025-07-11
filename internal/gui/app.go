package gui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"tala/internal/ai"
	"tala/internal/config"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type App struct {
	fyneApp  fyne.App
	window   fyne.Window
	provider ai.Provider
	config   *config.Config
	
	// UI components
	chatHistory *widget.RichText
	input       *widget.Entry
	sendButton  *widget.Button
	statusLabel *widget.Label
	
	// State
	isLoading bool
}

func NewApp(cfg *config.Config) (*App, error) {
	provider, err := ai.CreateProviderFromConfig(cfg)
	if err != nil {
		return nil, err
	}

	fyneApp := app.New()
	fyneApp.SetIcon(nil) // TODO: Add app icon
	
	window := fyneApp.NewWindow("Tala - AI Assistant")
	window.Resize(fyne.NewSize(800, 600))
	
	guiApp := &App{
		fyneApp:  fyneApp,
		window:   window,
		provider: provider,
		config:   cfg,
	}
	
	guiApp.setupUI()
	return guiApp, nil
}

func (a *App) setupUI() {
	// Chat history
	a.chatHistory = widget.NewRichText()
	a.chatHistory.Wrapping = fyne.TextWrapWord
	a.chatHistory.Scroll = container.ScrollVerticalOnly
	
	// Input field
	a.input = widget.NewEntry()
	a.input.SetPlaceHolder("Type your message...")
	a.input.MultiLine = true
	a.input.OnSubmitted = func(text string) {
		a.sendMessage(text)
	}
	
	// Send button
	a.sendButton = widget.NewButton("Send", func() {
		a.sendMessage(a.input.Text)
	})
	
	// Status label
	a.statusLabel = widget.NewLabel(fmt.Sprintf("Ready - Provider: %s", a.provider.GetName()))
	
	// Layout
	inputContainer := container.NewBorder(nil, nil, nil, a.sendButton, a.input)
	
	content := container.NewBorder(
		nil,                                           // top
		container.NewVBox(inputContainer, a.statusLabel), // bottom
		nil,                                           // left
		nil,                                           // right
		container.NewScroll(a.chatHistory),            // center
	)
	
	a.window.SetContent(content)
	
	// Setup menu
	a.setupMenu()
}

func (a *App) setupMenu() {
	// File menu
	newItem := fyne.NewMenuItem("New Chat", func() {
		a.chatHistory.ParseMarkdown("")
	})
	
	quitItem := fyne.NewMenuItem("Quit", func() {
		a.fyneApp.Quit()
	})
	
	fileMenu := fyne.NewMenu("File", newItem, fyne.NewMenuItemSeparator(), quitItem)
	
	// Settings menu
	settingsItem := fyne.NewMenuItem("Preferences", func() {
		a.showSettings()
	})
	
	aboutItem := fyne.NewMenuItem("About", func() {
		dialog.ShowInformation("About Tala", 
			"Tala - Terminal AI Language Assistant\n\nBuilt with Go and Fyne", 
			a.window)
	})
	
	settingsMenu := fyne.NewMenu("Settings", settingsItem, fyne.NewMenuItemSeparator(), aboutItem)
	
	// Main menu
	mainMenu := fyne.NewMainMenu(fileMenu, settingsMenu)
	a.window.SetMainMenu(mainMenu)
}

func (a *App) showSettings() {
	providerEntry := widget.NewEntry()
	providerEntry.SetText(a.config.Provider)
	
	modelEntry := widget.NewEntry()
	modelEntry.SetText(a.config.Model)
	
	apiKeyEntry := widget.NewPasswordEntry()
	apiKeyEntry.SetText(a.config.APIKey)
	
	tempSlider := widget.NewSlider(0, 2)
	tempSlider.SetValue(a.config.Temperature)
	tempSlider.Step = 0.1
	
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Provider", Widget: providerEntry},
			{Text: "Model", Widget: modelEntry},
			{Text: "API Key", Widget: apiKeyEntry},
			{Text: "Temperature", Widget: tempSlider},
		},
		OnSubmit: func() {
			a.config.Provider = providerEntry.Text
			a.config.Model = modelEntry.Text
			a.config.APIKey = apiKeyEntry.Text
			a.config.Temperature = tempSlider.Value
			
			if err := a.config.Save(); err != nil {
				dialog.ShowError(err, a.window)
				return
			}
			
			// Recreate provider with new config
			provider, err := ai.CreateProviderFromConfig(a.config)
			if err != nil {
				dialog.ShowError(err, a.window)
				return
			}
			
			a.provider = provider
			a.statusLabel.SetText(fmt.Sprintf("Ready - Provider: %s", a.provider.GetName()))
			
			dialog.ShowInformation("Settings", "Configuration saved successfully!", a.window)
		},
	}
	
	dialog.ShowForm("Settings", "Save", "Cancel", form.Items, func(submitted bool) {
		if submitted {
			form.OnSubmit()
		}
	}, a.window)
}

func (a *App) sendMessage(text string) {
	if strings.TrimSpace(text) == "" {
		return
	}
	
	if a.isLoading {
		return
	}
	
	// Add user message to chat
	a.addMessage("You", text, false)
	
	// Clear input
	a.input.SetText("")
	
	// Set loading state
	a.isLoading = true
	a.sendButton.Disable()
	a.statusLabel.SetText("AI is thinking...")
	
	// Process message in goroutine
	go func() {
		defer func() {
			a.isLoading = false
			a.sendButton.Enable()
			a.statusLabel.SetText(fmt.Sprintf("Ready - Provider: %s", a.provider.GetName()))
		}()
		
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		startTime := time.Now()
		
		// Check if message should use tools
		if a.provider.SupportsTools() && a.shouldUseTool(text) {
			response, toolResults, err := a.provider.GenerateResponseWithTools(ctx, text)
			if err != nil {
				a.addMessage("Error", fmt.Sprintf("Error: %v", err), true)
				return
			}
			
			// Add tool results if any
			if len(toolResults) > 0 {
				for _, result := range toolResults {
					a.addMessage("System", fmt.Sprintf("Executed: %s\nResult: %s", result.Name, result.Content), true)
				}
			}
			
			duration := time.Since(startTime)
			a.addMessage("AI", fmt.Sprintf("%s\n\n*Response time: %v*", response, duration), false)
		} else {
			response, err := a.provider.GenerateResponse(ctx, text)
			if err != nil {
				a.addMessage("Error", fmt.Sprintf("Error: %v", err), true)
				return
			}
			
			duration := time.Since(startTime)
			a.addMessage("AI", fmt.Sprintf("%s\n\n*Response time: %v*", response, duration), false)
		}
	}()
}

func (a *App) shouldUseTool(text string) bool {
	// Simple heuristic to determine if we should use tools
	keywords := []string{"file", "directory", "create", "read", "write", "list", "command", "run", "execute"}
	lowText := strings.ToLower(text)
	
	for _, keyword := range keywords {
		if strings.Contains(lowText, keyword) {
			return true
		}
	}
	
	return false
}

func (a *App) addMessage(sender, message string, isSystem bool) {
	timestamp := time.Now().Format("15:04:05")
	
	var formattedMessage string
	if isSystem {
		formattedMessage = fmt.Sprintf("**[%s] %s:**\n%s\n\n", timestamp, sender, message)
	} else {
		formattedMessage = fmt.Sprintf("**[%s] %s:**\n%s\n\n", timestamp, sender, message)
	}
	
	// Append to existing content
	currentContent := a.chatHistory.String()
	newContent := currentContent + formattedMessage
	a.chatHistory.ParseMarkdown(newContent)
}

func (a *App) Run() {
	a.window.ShowAndRun()
}