package tui

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"tala/internal/ai"
	"tala/internal/config"
	"tala/internal/fileops"
)

// ANSI color codes for better UX
const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Gray    = "\033[37m"
	Bold    = "\033[1m"
	Dim     = "\033[2m"
)

// getTerminalWidth returns the terminal width, or 80 as default
func getTerminalWidth() int {
	// Try to get terminal width using tput
	if cmd := exec.Command("tput", "cols"); cmd != nil {
		if output, err := cmd.Output(); err == nil {
			if width, err := strconv.Atoi(strings.TrimSpace(string(output))); err == nil && width > 0 {
				return width - 2 // Leave a small margin
			}
		}
	}
	
	// Fallback to default width
	return 78
}

// SimpleTUI provides a basic terminal interface without external dependencies
type SimpleTUI struct {
	provider      ai.Provider
	config        *config.Config
	totalTokens   int
	totalRequests int
	totalTime     time.Duration
	currentTokens int // Track tokens for current response
}

// NewSimpleTUI creates a new simple TUI instance
func NewSimpleTUI(cfg *config.Config) (*SimpleTUI, error) {
	provider, err := ai.CreateProvider(cfg.Provider, cfg.APIKey, cfg.Model, cfg.Temperature, cfg.MaxTokens)
	if err != nil {
		return nil, err
	}

	return &SimpleTUI{
		provider: provider,
		config:   cfg,
	}, nil
}

// Run starts the simple TUI
func (s *SimpleTUI) Run() error {
	// Setup signal handling for clean exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	
	// Print colorful header
	fmt.Printf("\n%s🗣️ Tala - Terminal AI Language Assistant%s\n", Bold+Cyan, Reset)
	fmt.Printf("%sProvider:%s %s%s%s %s|%s %sModel:%s %s%s%s\n", 
		Dim, Reset, Green, s.provider.GetName(), Reset,
		Dim, Reset, Dim, Reset, Yellow, s.config.Model, Reset)
	fmt.Printf("%sType '%s/help%s' for file operations or chat normally with AI%s\n", 
		Gray, Cyan, Gray, Reset)
	fmt.Printf("%sCtrl+C to exit%s\n\n", Dim, Reset)

	// Channel for input
	inputChan := make(chan string)
	aiBusy := false
	
	// Start input reader goroutine
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			inputChan <- scanner.Text()
		}
		close(inputChan)
	}()
	
	// Show initial prompt with color
	fmt.Printf("%s> %s", Blue+Bold, Reset)
	
	for {
		select {
		case <-c:
			fmt.Println("\nGoodbye!")
			return nil
			
		case input, ok := <-inputChan:
			if !ok {
				return nil // EOF
			}
			
			input = strings.TrimSpace(input)
			if input == "" {
				if !aiBusy {
					fmt.Printf("%s> %s", Blue+Bold, Reset)
				}
				continue
			}

			// Handle exit commands
			if input == "exit" || input == "quit" || input == "/quit" || input == "/exit" {
				fmt.Println("Goodbye!")
				return nil
			}

			// If AI is busy, queue the input for later
			if aiBusy {
				fmt.Printf("\n%s[Queued]:%s %s\n", Yellow, Reset, input)
				// TODO: Implement proper queuing
				continue
			}

			// Handle slash commands
			if strings.HasPrefix(input, "/") {
				s.handleSlashCommand(input)
				fmt.Printf("%s> %s", Blue+Bold, Reset)
				continue
			}

			// Handle AI conversation
			aiBusy = true
			go func(prompt string) {
				s.handleAIConversation(prompt)
				aiBusy = false
				fmt.Printf("%s> %s", Blue+Bold, Reset)
			}(input)
		}
	}
}

// handleAIConversation processes AI chat with streaming paragraph updates
func (s *SimpleTUI) handleAIConversation(input string) {
	fmt.Printf("%sYou:%s %s\n", Green+Bold, Reset, input)
	
	start := time.Now()
	
	// Show thinking indicator with live stats
	done := make(chan bool, 1)
	go s.showThinkingProgress(start, done)
	
	ctx := context.Background()
	var response string
	var err error
	var toolResults []ai.ToolResult

	// Get the response (still non-streaming to avoid API complexity)
	if s.provider.SupportsTools() {
		response, toolResults, err = s.provider.GenerateResponseWithTools(ctx, input)
	} else {
		response, err = s.provider.GenerateResponse(ctx, input)
	}

	// Stop thinking indicator
	done <- true
	fmt.Print("\r\033[K") // Clear the thinking line

	// Handle errors
	if err != nil {
		fmt.Printf("%sError:%s %s\n\n", Red+Bold, Reset, err.Error())
		return
	}

	// Display tool results if any
	if len(toolResults) > 0 {
		fmt.Printf("%sSystem:%s File operations executed:\n", Cyan+Bold, Reset)
		for _, result := range toolResults {
			fmt.Printf("  %s✓%s %s: %s\n", Green, Reset, result.Name, result.Content)
		}
		fmt.Println()
	}

	// Display AI response with paragraph-based streaming simulation
	fmt.Printf("%sAI:%s ", Magenta+Bold, Reset)
	s.displayResponseByParagraphs(response)

	// Update and display colorful stats
	duration := time.Since(start)
	s.totalRequests++
	tokens := len(strings.Fields(response))
	s.totalTokens += tokens
	s.totalTime += duration

	// Display colorful stats
	fmt.Printf("%s[%sTokens:%s %s%d%s %s|%s %sTime:%s %s%s%s%s]%s\n\n", 
		Dim, Reset+Cyan, Dim, Yellow, tokens, Dim, Reset+Dim, Dim, Reset+Cyan, Dim, 
		Green, duration.Round(time.Millisecond), Dim, Reset+Dim, Reset)
}

// handleSlashCommand processes slash commands
func (s *SimpleTUI) handleSlashCommand(cmd string) {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return
	}

	command := parts[0]

	switch command {
	case "/help":
		s.showHelp()
	case "/clear":
		s.clearScreen()
	case "/stats":
		s.showStats()
	case "/config":
		s.showConfig()
	case "/exit", "/quit":
		fmt.Printf("%sGoodbye!%s\n", Green+Bold, Reset)
		os.Exit(0)
	default:
		// Try file operation
		result := fileops.ExecuteCommand(cmd)
		if strings.Contains(result.Message, "✓") || strings.Contains(result.Message, "success") {
			fmt.Printf("%sSystem:%s %s\n\n", Green+Bold, Reset, result.Message)
		} else {
			fmt.Printf("%sSystem:%s %s\n\n", Red+Bold, Reset, result.Message)
		}
	}
}

// showHelp displays help information
func (s *SimpleTUI) showHelp() {
	fmt.Printf("%sAvailable Commands:%s\n\n", Cyan+Bold, Reset)
	
	fmt.Printf("%sSystem Commands:%s\n", Yellow+Bold, Reset)
	fmt.Printf("  %s/clear%s           Clear screen and reset session\n", Green, Reset)
	fmt.Printf("  %s/stats%s           Show session statistics\n", Green, Reset)
	fmt.Printf("  %s/config%s          Show current configuration\n", Green, Reset)
	fmt.Printf("  %s/help%s            Show this help message\n", Green, Reset)
	fmt.Printf("  %s/exit, /quit%s     Exit application\n\n", Green, Reset)
	
	fmt.Printf("%sFile Operations:%s\n", Yellow+Bold, Reset)
	fmt.Printf("  %s/ls [path]%s       List files and directories\n", Green, Reset)
	fmt.Printf("  %s/cat <file>%s      Display file content\n", Green, Reset)
	fmt.Printf("  %s/pwd%s             Show current directory\n", Green, Reset)
	fmt.Printf("  %s/cd <path>%s       Change directory\n", Green, Reset)
	fmt.Printf("  %s/create <file>%s   Create new file\n", Green, Reset)
	fmt.Printf("  %s/mkdir <dir>%s     Create directory\n\n", Green, Reset)
	
	fmt.Printf("%sKeyboard Shortcuts:%s\n", Yellow+Bold, Reset)
	fmt.Printf("  %sCtrl+C%s           Exit application\n", Green, Reset)
	fmt.Printf("  %sEnter%s            Send message\n\n", Green, Reset)
}

// clearScreen clears the terminal
func (s *SimpleTUI) clearScreen() {
	fmt.Print("\033[2J\033[H")
	s.totalTokens = 0
	s.totalRequests = 0
	s.totalTime = 0
	
	fmt.Printf("%s🗣️ Tala - Terminal AI Language Assistant%s\n", Bold+Cyan, Reset)
	fmt.Printf("%sProvider:%s %s%s%s %s|%s %sModel:%s %s%s%s\n", 
		Dim, Reset, Green, s.provider.GetName(), Reset,
		Dim, Reset, Dim, Reset, Yellow, s.config.Model, Reset)
	fmt.Printf("%sType '%s/help%s' for commands or chat normally with AI%s\n\n", 
		Gray, Cyan, Gray, Reset)
}

// showStats displays session statistics
func (s *SimpleTUI) showStats() {
	if s.totalRequests > 0 {
		avgTime := s.totalTime / time.Duration(s.totalRequests)
		fmt.Printf("%sSession Stats:%s %s%d%s requests, %s%d%s tokens, avg %s%s%s\n\n", 
			Cyan+Bold, Reset, Green, s.totalRequests, Reset, 
			Green, s.totalTokens, Reset, Yellow, avgTime.Round(time.Millisecond), Reset)
	} else {
		fmt.Printf("%sNo requests made yet%s\n\n", Dim, Reset)
	}
}

// showConfig displays current configuration
func (s *SimpleTUI) showConfig() {
	fmt.Printf("%sCurrent Configuration:%s\n", Cyan+Bold, Reset)
	fmt.Printf("  %sProvider:%s %s%s%s\n", Yellow, Reset, Green, s.config.Provider, Reset)
	fmt.Printf("  %sModel:%s %s%s%s\n", Yellow, Reset, Green, s.config.Model, Reset)
	fmt.Printf("  %sTemperature:%s %s%.1f%s\n", Yellow, Reset, Green, s.config.Temperature, Reset)
	fmt.Printf("  %sMax Tokens:%s %s%d%s\n", Yellow, Reset, Green, s.config.MaxTokens, Reset)
	fmt.Printf("  %sTools:%s %s%v%s\n\n", Yellow, Reset, Green, s.provider.SupportsTools(), Reset)
}

// displayResponseByParagraphs displays AI response paragraph by paragraph with natural timing
func (s *SimpleTUI) displayResponseByParagraphs(response string) {
	// Split response into paragraphs (double newlines or single newlines)
	paragraphs := strings.Split(response, "\n")
	
	// Process each paragraph
	for i, paragraph := range paragraphs {
		paragraph = strings.TrimSpace(paragraph)
		
		if paragraph == "" {
			// Empty paragraph - just add a newline
			fmt.Println()
			continue
		}
		
		// Wrap the paragraph text
		wrappedParagraph := s.wrapText(paragraph, getTerminalWidth())
		
		// Display the paragraph
		fmt.Print(wrappedParagraph)
		
		// Add newline after paragraph (except for the last one)
		if i < len(paragraphs)-1 {
			fmt.Println()
		}
		
		// Add a slight delay between paragraphs for natural reading flow
		// (but not too long to avoid feeling slow)
		if i < len(paragraphs)-1 && paragraph != "" {
			time.Sleep(200 * time.Millisecond)
		}
	}
	
	// Ensure we end with a newline
	fmt.Println()
}

// wrapText performs basic word wrapping
func (s *SimpleTUI) wrapText(text string, width int) string {
	if width <= 0 {
		width = 80
	}
	
	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}
	
	var result strings.Builder
	currentLine := ""
	
	for _, word := range words {
		testLine := currentLine
		if testLine != "" {
			testLine += " " + word
		} else {
			testLine = word
		}
		
		if len(testLine) <= width {
			currentLine = testLine
		} else {
			if currentLine != "" {
				result.WriteString(currentLine + "\n")
			}
			currentLine = word
		}
	}
	
	if currentLine != "" {
		result.WriteString(currentLine)
	}
	
	return result.String()
}


// showThinkingProgress displays clean thinking progress with stats
func (s *SimpleTUI) showThinkingProgress(start time.Time, done chan bool) {
	ticker := time.NewTicker(400 * time.Millisecond)
	defer ticker.Stop()
	
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			elapsed := time.Since(start)
			
			// Format elapsed time with consistent width (always shows as X.Xs format)
			elapsedSeconds := elapsed.Seconds()
			timeStr := fmt.Sprintf("%4.1fs", elapsedSeconds)
			
			// Always show session stats with consistent formatting
			var avgSeconds float64
			if s.totalRequests > 0 {
				avgTime := s.totalTime / time.Duration(s.totalRequests)
				avgSeconds = avgTime.Seconds()
			}
			
			// Create complete progress line with consistent formatting
			progressText := fmt.Sprintf("%s🤔 AI is thinking...%s %s(%s)%s %s|%s %sSession:%s %s%3d%s req, %s%5d%s tokens, avg %s%4.1fs%s", 
				Yellow, Reset, Dim, timeStr, Reset,
				Dim, Reset, Cyan, Reset, Green, s.totalRequests, Reset,
				Green, s.totalTokens, Reset, Yellow, avgSeconds, Reset)
			
			// Clear the line completely and write the new progress
			fmt.Print("\r\033[K")  // Clear entire line
			fmt.Print(progressText)
		}
	}
}