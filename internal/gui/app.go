package gui

import (
	"context"
	"fmt"
	"image/color"
	"strings"
	"time"

	"tala/internal/ai"
	"tala/internal/config"
	"tala/internal/fileops"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Custom colors matching terminal theme
var (
	UserColor    = color.NRGBA{R: 0, G: 255, B: 0, A: 255}     // Green
	AIColor      = color.NRGBA{R: 255, G: 0, B: 255, A: 255}   // Magenta
	SystemColor  = color.NRGBA{R: 0, G: 255, B: 255, A: 255}   // Cyan
	ErrorColor   = color.NRGBA{R: 255, G: 0, B: 0, A: 255}     // Red
	StatsColor   = color.NRGBA{R: 255, G: 255, B: 0, A: 255}   // Yellow
	PromptColor  = color.NRGBA{R: 0, G: 100, B: 255, A: 255}   // Blue
)

// CustomTheme extends the dark theme with better text colors
type CustomTheme struct {
	fyne.Theme
}

func (t *CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	// Override text colors to make them more readable
	switch name {
	case theme.ColorNameForeground:
		return color.NRGBA{R: 240, G: 240, B: 240, A: 255} // Light text
	case theme.ColorNamePlaceHolder:
		return color.NRGBA{R: 180, G: 180, B: 180, A: 255} // Medium gray for placeholders
	}
	// Use the default dark theme for everything else
	return theme.DarkTheme().Color(name, variant)
}

func (t *CustomTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DarkTheme().Font(style)
}

func (t *CustomTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DarkTheme().Icon(name)
}

func (t *CustomTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DarkTheme().Size(name)
}

type App struct {
	fyneApp  fyne.App
	window   fyne.Window
	provider ai.Provider
	config   *config.Config
	
	// UI components
	chatHistory   *widget.Entry // Using Entry for copy-paste functionality
	input         *widget.Entry
	sendButton    *widget.Button
	statusLabel   *widget.Label
	statsLabel    *widget.Label
	progressBar   *widget.ProgressBarInfinite
	
	// Enhanced UI components
	providerLabel *widget.Label
	modelLabel    *widget.Label
	clearButton   *widget.Button
	
	// State
	isLoading      bool
	totalRequests  int
	totalTokens    int
	totalTime      time.Duration
	currentStart   time.Time
	
	// Concurrent input handling
	inputQueue     chan string
	processingLock bool
	
	// Chat history management
	chatContent    string
}

func NewApp(cfg *config.Config) (*App, error) {
	provider, err := ai.CreateProvider(cfg.Provider, cfg.APIKey, cfg.Model, cfg.Temperature, cfg.MaxTokens)
	if err != nil {
		return nil, err
	}

	fyneApp := app.New()
	fyneApp.Settings().SetTheme(&CustomTheme{}) // Use custom theme with better text colors
	fyneApp.SetIcon(nil) // TODO: Add app icon
	
	window := fyneApp.NewWindow("Tala - Terminal AI Language Assistant")
	window.Resize(fyne.NewSize(1000, 700)) // Larger window
	
	guiApp := &App{
		fyneApp:    fyneApp,
		window:     window,
		provider:   provider,
		config:     cfg,
		inputQueue: make(chan string, 100), // Buffered channel for input queue
	}
	
	guiApp.setupUI()
	guiApp.startInputProcessor()
	return guiApp, nil
}

func (a *App) setupUI() {
	// Chat history with copy-paste functionality - use Entry but keep it enabled for better text color
	a.chatHistory = widget.NewEntry()
	a.chatHistory.MultiLine = true
	a.chatHistory.Wrapping = fyne.TextWrapWord
	
	// Make it read-only by preventing changes
	a.chatHistory.OnChanged = func(content string) {
		// If content was changed by user input (not by our SetText calls), revert it
		if a.chatContent != "" && content != a.chatContent {
			a.chatHistory.SetText(a.chatContent)
		}
	}
	
	// Add welcome message
	a.addWelcomeMessage()
	
	// Larger input field
	a.input = widget.NewEntry()
	a.input.SetPlaceHolder("Type your message here... (Enter for new line, Shift+Enter to send)")
	a.input.MultiLine = true
	a.input.Resize(fyne.NewSize(600, 100)) // Much larger input field
	
	// Enhanced input handling - Shift+Enter sends, Enter adds new line
	a.input.OnSubmitted = func(text string) {
		// OnSubmitted is called on Shift+Enter in multiline mode
		if !a.processingLock {
			a.queueMessage(text)
		}
	}
	
	// Enhanced send button
	a.sendButton = widget.NewButton("Send", func() {
		if !a.processingLock {
			a.queueMessage(a.input.Text)
		}
	})
	a.sendButton.Importance = widget.HighImportance
	
	// Clear chat button
	a.clearButton = widget.NewButton("Clear Chat", func() {
		a.chatContent = ""
		a.chatHistory.SetText("")
		a.addWelcomeMessage()
		a.totalRequests = 0
		a.totalTokens = 0
		a.totalTime = 0
		a.updateStats()
	})
	
	// Provider and model labels - clean text without emojis for better compatibility
	a.providerLabel = widget.NewLabel(fmt.Sprintf("Provider: %s", a.provider.GetName()))
	a.modelLabel = widget.NewLabel(fmt.Sprintf("Model: %s", a.config.Model))
	
	// Make labels more prominent
	a.providerLabel.Importance = widget.MediumImportance
	a.modelLabel.Importance = widget.MediumImportance
	
	// Status label with color
	a.statusLabel = widget.NewLabel("Ready - Type your message below")
	a.statusLabel.Importance = widget.MediumImportance
	
	// Statistics label
	a.statsLabel = widget.NewLabel("Session: 0 requests, 0 tokens, 0.0s avg")
	a.statsLabel.Importance = widget.LowImportance
	
	// Progress bar (hidden initially)
	a.progressBar = widget.NewProgressBarInfinite()
	a.progressBar.Hide()
	
	// Enhanced layout
	headerContainer := container.NewHBox(
		a.providerLabel,
		widget.NewSeparator(),
		a.modelLabel,
	)
	
	inputContainer := container.NewBorder(
		nil, nil, nil, 
		container.NewVBox(a.sendButton, a.clearButton),
		a.input,
	)
	
	statusContainer := container.NewVBox(
		a.statusLabel,
		a.statsLabel,
		a.progressBar,
	)
	
	// Main layout with better spacing
	content := container.NewBorder(
		headerContainer,        // top
		container.NewVBox(      // bottom
			widget.NewSeparator(),
			inputContainer,
			statusContainer,
		),
		nil,                    // left
		nil,                    // right
		container.NewScroll(a.chatHistory), // center
	)
	
	a.window.SetContent(content)
	
	// Setup menu
	a.setupMenu()
}

func (a *App) addWelcomeMessage() {
	welcome := fmt.Sprintf(`Welcome to Tala!

Provider: %s  
Model: %s  
Tools: %v

Type your message below and press Enter to chat with AI. You can:
- Ask questions naturally
- Request file operations: "create a file called test.txt"
- Execute commands: "list files in current directory"
- Get help: "what can you do?"

=================================================================

`, a.provider.GetName(), a.config.Model, a.provider.SupportsTools())
	
	a.chatContent = welcome
	a.chatHistory.SetText(welcome)
}

func (a *App) setupMenu() {
	// File menu
	newItem := fyne.NewMenuItem("New Chat", func() {
		a.chatContent = ""
		a.chatHistory.SetText("")
		a.addWelcomeMessage()
		a.totalRequests = 0
		a.totalTokens = 0
		a.totalTime = 0
		a.updateStats()
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
			"Tala - Terminal AI Language Assistant\n\n"+
			"Built with Go and Fyne\n"+
			"Enhanced GUI with professional interface\n"+
			"Multi-provider AI support\n"+
			"Intelligent file operations\n\n"+
			"Features:\n"+
			"• Professional, responsive interface\n"+
			"• Concurrent input handling\n"+
			"• Real-time statistics\n"+
			"• File operations support\n"+
			"• Cross-platform compatibility", 
			a.window)
	})
	
	settingsMenu := fyne.NewMenu("Settings", settingsItem, fyne.NewMenuItemSeparator(), aboutItem)
	
	// Help menu
	helpItem := fyne.NewMenuItem("Help", func() {
		helpText := `# Tala Help

## Basic Usage
- Type messages in the input field
- Press Enter for new lines
- Use Shift+Enter to send messages

## File Operations
- "create a file called example.txt"
- "list files in current directory"
- "read the content of README.md"
- "write hello world to test.txt"

## Commands
- Natural language commands work automatically
- AI will understand your intent and execute appropriate actions

## Interface Features
- **Colorful Messages**: Different colors for user, AI, and system messages
- **Real-time Stats**: Session statistics displayed below
- **Concurrent Input**: Type next message while AI processes current one
- **Progress Indicators**: Visual feedback during AI processing

## Keyboard Shortcuts
- **Enter**: New line in input
- **Shift+Enter**: Send message
- **Ctrl+N**: New chat (clear history)
- **Ctrl+Q**: Quit application
`
		dialog.ShowInformation("Help", helpText, a.window)
	})
	
	helpMenu := fyne.NewMenu("Help", helpItem, aboutItem)
	
	// Main menu
	mainMenu := fyne.NewMainMenu(fileMenu, settingsMenu, helpMenu)
	a.window.SetMainMenu(mainMenu)
}

func (a *App) showSettings() {
	// Create much larger input fields - increased width significantly for better usability
	providerEntry := widget.NewEntry()
	providerEntry.SetText(a.config.Provider)
	providerEntry.Resize(fyne.NewSize(800, 60))
	
	modelEntry := widget.NewEntry()
	modelEntry.SetText(a.config.Model)
	modelEntry.Resize(fyne.NewSize(800, 60))
	
	apiKeyEntry := widget.NewPasswordEntry()
	apiKeyEntry.SetText(a.config.APIKey)
	apiKeyEntry.Resize(fyne.NewSize(800, 60))
	
	// Temperature as input field instead of slider - make it much larger
	tempEntry := widget.NewEntry()
	tempEntry.SetText(fmt.Sprintf("%.1f", a.config.Temperature))
	tempEntry.Resize(fyne.NewSize(800, 60))
	
	maxTokensEntry := widget.NewEntry()
	maxTokensEntry.SetText(fmt.Sprintf("%d", a.config.MaxTokens))
	maxTokensEntry.Resize(fyne.NewSize(800, 60))
	
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Provider", Widget: providerEntry},
			{Text: "Model", Widget: modelEntry},
			{Text: "API Key", Widget: apiKeyEntry},
			{Text: "Temperature (0.0-2.0)", Widget: tempEntry},
			{Text: "Max Tokens (0=unlimited)", Widget: maxTokensEntry},
		},
		OnSubmit: func() {
			// Validate and save configuration
			a.config.Provider = providerEntry.Text
			a.config.Model = modelEntry.Text
			a.config.APIKey = apiKeyEntry.Text
			
			// Parse temperature
			if tempText := tempEntry.Text; tempText != "" {
				if temp, err := fmt.Sscanf(tempText, "%f", &a.config.Temperature); err != nil || temp != 1 {
					a.config.Temperature = 0.7 // Default value
				}
				// Clamp temperature between 0.0 and 2.0
				if a.config.Temperature < 0.0 {
					a.config.Temperature = 0.0
				} else if a.config.Temperature > 2.0 {
					a.config.Temperature = 2.0
				}
			}
			
			// Parse max tokens
			if maxTokensText := maxTokensEntry.Text; maxTokensText != "" {
				if maxTokens, err := fmt.Sscanf(maxTokensText, "%d", &a.config.MaxTokens); err != nil || maxTokens != 1 {
					a.config.MaxTokens = 0 // Default to unlimited
				}
			}
			
			if err := a.config.Save(); err != nil {
				dialog.ShowError(err, a.window)
				return
			}
			
			// Recreate provider with new config
			provider, err := ai.CreateProvider(a.config.Provider, a.config.APIKey, a.config.Model, a.config.Temperature, a.config.MaxTokens)
			if err != nil {
				dialog.ShowError(err, a.window)
				return
			}
			
			a.provider = provider
			a.providerLabel.SetText(fmt.Sprintf("Provider: %s", a.provider.GetName()))
			a.modelLabel.SetText(fmt.Sprintf("Model: %s", a.config.Model))
			a.statusLabel.SetText("Ready - Configuration updated")
			
			dialog.ShowInformation("Settings", "Configuration saved successfully!\n\nNew provider and model are now active.", a.window)
		},
	}
	
	// Create a custom dialog with larger size
	var customDialog *dialog.CustomDialog
	
	saveButton := widget.NewButton("Save", func() {
		form.OnSubmit()
		customDialog.Hide()
	})
	saveButton.Importance = widget.HighImportance
	
	cancelButton := widget.NewButton("Cancel", func() {
		customDialog.Hide()
	})
	
	buttons := container.NewHBox(saveButton, cancelButton)
	
	content := container.NewVBox(
		widget.NewLabel("Settings - AI Provider Configuration"),
		widget.NewSeparator(),
		container.NewGridWithColumns(2,
			widget.NewLabel("Provider:\n(AI service: ollama, openai, anthropic)"), providerEntry,
			widget.NewLabel("Model:\n(AI model name for the provider)"), modelEntry,
			widget.NewLabel("API Key:\n(Required for OpenAI/Anthropic, not needed for Ollama)"), apiKeyEntry,
			widget.NewLabel("Temperature (0.0-2.0):\n(Response creativity: 0.0=focused, 2.0=creative)"), tempEntry,
			widget.NewLabel("Max Tokens (0=unlimited):\n(Maximum response length, 0 for no limit)"), maxTokensEntry,
		),
		widget.NewSeparator(),
		buttons,
	)
	
	customDialog = dialog.NewCustom("Settings", "Close", content, a.window)
	customDialog.Resize(fyne.NewSize(900, 400)) // Make dialog larger
	customDialog.Show()
}

func (a *App) startInputProcessor() {
	go func() {
		for message := range a.inputQueue {
			a.processMessage(message)
		}
	}()
}

func (a *App) queueMessage(text string) {
	text = strings.TrimSpace(text)
	if text == "" {
		return
	}
	
	select {
	case a.inputQueue <- text:
		// Clear input immediately for better UX
		a.input.SetText("")
	default:
		// Queue is full, show warning
		a.statusLabel.SetText("Message queue full, please wait...")
	}
}

func (a *App) processMessage(text string) {
	if a.processingLock {
		return
	}
	
	a.processingLock = true
	defer func() {
		a.processingLock = false
	}()
	
	// Add user message to chat
	a.addMessage("You", text, UserColor)
	
	// Set loading state
	a.statusLabel.SetText("AI is thinking...")
	a.progressBar.Show()
	a.progressBar.Start()
	a.sendButton.Disable()
	
	// Start timing
	a.currentStart = time.Now()
	
	// Process message
	go func() {
		defer func() {
			a.progressBar.Stop()
			a.progressBar.Hide()
			a.sendButton.Enable()
			a.statusLabel.SetText("Ready - Type your message below")
			a.updateStats()
		}()
		
		ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
		defer cancel()
		
		// Handle slash commands
		if strings.HasPrefix(text, "/") {
			a.handleSlashCommand(text)
			return
		}
		
		// Check if we should use tools
		if a.provider.SupportsTools() {
			response, toolResults, err := a.provider.GenerateResponseWithTools(ctx, text)
			if err != nil {
				a.addMessage("Error", fmt.Sprintf("Error: %v", err), ErrorColor)
				return
			}
			
			// Add tool results if any
			if len(toolResults) > 0 {
				for _, result := range toolResults {
					a.addMessage("System", fmt.Sprintf("🛠️ **%s**: %s", result.Name, result.Content), SystemColor)
				}
			}
			
			// Add AI response with paragraph-based display
			a.addAIResponseWithDelay(response)
		} else {
			response, err := a.provider.GenerateResponse(ctx, text)
			if err != nil {
				a.addMessage("Error", fmt.Sprintf("Error: %v", err), ErrorColor)
				return
			}
			
			// Add AI response with paragraph-based display
			a.addAIResponseWithDelay(response)
		}
		
		// Update statistics
		a.totalRequests++
		tokens := len(strings.Fields(text)) // Simple token approximation
		a.totalTokens += tokens
		a.totalTime += time.Since(a.currentStart)
	}()
}

func (a *App) handleSlashCommand(cmd string) {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return
	}
	
	command := parts[0]
	
	switch command {
	case "/help":
		helpText := `## Available Commands

### System Commands
- **/clear** - Clear chat history
- **/stats** - Show session statistics
- **/help** - Show this help message
- **/quit** - Exit application

### File Operations
- **/ls [path]** - List files and directories
- **/cat <file>** - Display file content
- **/pwd** - Show current directory
- **/cd <path>** - Change directory
- **/create <file>** - Create new file
- **/mkdir <dir>** - Create directory

### Tips
- You can also use natural language: "create a file called test.txt"
- AI will understand and execute appropriate file operations
- Use the input field below for normal conversations
`
		a.addMessage("System", helpText, SystemColor)
		
	case "/clear":
		a.chatContent = ""
		a.chatHistory.SetText("")
		a.addWelcomeMessage()
		a.totalRequests = 0
		a.totalTokens = 0
		a.totalTime = 0
		a.updateStats()
		
	case "/stats":
		if a.totalRequests > 0 {
			avgTime := a.totalTime / time.Duration(a.totalRequests)
			statsText := fmt.Sprintf("📊 **Session Statistics:**\n\n- **Requests**: %d\n- **Tokens**: %d\n- **Average Time**: %v\n- **Total Time**: %v", 
				a.totalRequests, a.totalTokens, avgTime.Round(time.Millisecond), a.totalTime.Round(time.Millisecond))
			a.addMessage("System", statsText, StatsColor)
		} else {
			a.addMessage("System", "📊 No requests made yet", SystemColor)
		}
		
	case "/quit":
		a.fyneApp.Quit()
		
	default:
		// Try file operation
		result := fileops.ExecuteCommand(cmd)
		if strings.Contains(result.Message, "✓") || strings.Contains(result.Message, "success") {
			a.addMessage("System", fmt.Sprintf("✅ %s", result.Message), SystemColor)
		} else {
			a.addMessage("System", fmt.Sprintf("❌ %s", result.Message), ErrorColor)
		}
	}
}

func (a *App) addAIResponseWithDelay(response string) {
	// Simply add the AI response as a regular message
	a.addMessage("AI", response, AIColor)
}

func (a *App) addMessage(sender, message string, textColor color.Color) {
	timestamp := time.Now().Format("15:04:05")
	
	// Map sender to clean prefix
	var prefix string
	switch sender {
	case "You":
		prefix = "USER"
	case "AI":
		prefix = "AI"
	case "System":
		prefix = "SYSTEM"
	case "Error":
		prefix = "ERROR"
	default:
		prefix = "MSG"
	}
	
	// Get current content
	currentContent := a.chatContent
	
	// Create formatted message with proper spacing and clear separators
	var formattedMessage string
	if strings.TrimSpace(currentContent) != "" {
		// Add clear separator line between messages
		formattedMessage = fmt.Sprintf("\n=================================================================\n\n[%s] %s:\n\n%s\n\n", timestamp, prefix, message)
	} else {
		formattedMessage = fmt.Sprintf("[%s] %s:\n\n%s\n\n", timestamp, prefix, message)
	}
	
	// Append to existing content
	newContent := currentContent + formattedMessage
	a.chatContent = newContent
	a.chatHistory.SetText(newContent)
}

func (a *App) updateStats() {
	if a.totalRequests > 0 {
		avgTime := a.totalTime / time.Duration(a.totalRequests)
		a.statsLabel.SetText(fmt.Sprintf("Session: %d requests, %d tokens, %v avg", 
			a.totalRequests, a.totalTokens, avgTime.Round(time.Millisecond)))
	} else {
		a.statsLabel.SetText("Session: 0 requests, 0 tokens, 0.0s avg")
	}
}

func (a *App) Run() {
	a.window.ShowAndRun()
}