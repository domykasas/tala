package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"tala/ai"
	"tala/config"
	"tala/fileops"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	input         string
	ready         bool
	loading       bool
	provider      ai.Provider
	config        *config.Config
	totalTokens   int
	totalRequests int
	totalTime     time.Duration
	requestStart  time.Time
}

type msgReady struct{}
type msgResponse struct {
	response    string
	err         error
	tokens      int
	duration    time.Duration
	toolResults []ai.ToolResult
}
type msgTick struct{}

func NewModel(cfg *config.Config) (*Model, error) {
	provider, err := ai.CreateProvider(cfg.Provider, cfg.APIKey, cfg.Model, cfg.Temperature, cfg.MaxTokens)
	if err != nil {
		return nil, err
	}

	return &Model{
		input:    "",
		provider: provider,
		config:   cfg,
	}, nil
}

func (m *Model) Init() tea.Cmd {
	return func() tea.Msg {
		return msgReady{}
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case msgReady:
		m.ready = true
		fmt.Print("\n")
		titleStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("39"))
		fmt.Println(titleStyle.Render("🗣️ Tala - Terminal AI Language Assistant"))
		
		statsStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
		fmt.Println(statsStyle.Render(fmt.Sprintf("Provider: %s | Model: %s", m.provider.GetName(), m.config.Model)))
		fmt.Println(statsStyle.Render("Type '/help' for file operations or chat normally with AI"))
		fmt.Print("\n")
		return m, nil
	case msgTick:
		if m.loading {
			return m, m.tickCmd()
		}
		return m, nil
	case msgResponse:
		m.loading = false
		m.totalRequests++
		m.totalTokens += msg.tokens
		m.totalTime += msg.duration
		
		if msg.err != nil {
			fmt.Printf("Error: %s\n", msg.err.Error())
		} else {
			// Display file operations that were executed
			if len(msg.toolResults) > 0 {
				toolStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("226")) // Yellow
				fmt.Printf("%s File operations executed:\n", toolStyle.Render("System:"))
				for _, result := range msg.toolResults {
					var resultStyle lipgloss.Style
					if result.Success {
						resultStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("46")) // Green
					} else {
						resultStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196")) // Red
					}
					fmt.Printf("  %s %s: %s\n", resultStyle.Render("✓"), result.Name, result.Content)
				}
				fmt.Print("\n")
			}
			
			// Display AI response
			aiStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("46"))
			fmt.Printf("%s %s\n", aiStyle.Render("AI:"), msg.response)
			
			if msg.tokens > 0 {
				statsStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
				fmt.Printf("%s\n", statsStyle.Render(fmt.Sprintf("[Tokens: %d, Time: %s]", msg.tokens, msg.duration.Round(time.Millisecond))))
			}
		}
		fmt.Print("\n")
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			fmt.Println("\nGoodbye!")
			return m, tea.Quit
		case "enter":
			if m.input != "" && !m.loading {
				userStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
				fmt.Printf("%s %s\n", userStyle.Render("You:"), m.input)
				
				prompt := m.input
				m.input = ""
				
				// Check if this is a file operation command
				if strings.HasPrefix(prompt, "/") {
					result := fileops.ExecuteCommand(prompt)
					
					// Style the response based on success/failure
					var responseStyle lipgloss.Style
					if result.Success {
						responseStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("46")) // Green
					} else {
						responseStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196")) // Red
					}
					
					fmt.Printf("%s %s\n", responseStyle.Render("System:"), result.Message)
					fmt.Print("\n")
					return m, nil
				}
				
				// Regular AI conversation
				m.loading = true
				m.requestStart = time.Now()
				return m, tea.Batch(m.generateResponse(prompt), m.tickCmd())
			}
			return m, nil
		case "backspace":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
			return m, nil
		case "ctrl+l":
			fmt.Print("\033[2J\033[H")
			m.totalTokens = 0
			m.totalRequests = 0
			m.totalTime = 0
			
			titleStyle := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("39"))
			fmt.Println(titleStyle.Render("🗣️ Tala - Terminal AI Language Assistant"))
			
			statsStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
			fmt.Println(statsStyle.Render(fmt.Sprintf("Provider: %s | Model: %s", m.provider.GetName(), m.config.Model)))
			fmt.Println(statsStyle.Render("Type '/help' for file operations or chat normally with AI"))
			fmt.Print("\n")
			return m, nil
		default:
			if !m.loading {
				m.input += msg.String()
			}
			return m, nil
		}
	}
	return m, nil
}

func (m *Model) tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		return msgTick{}
	})
}

func (m *Model) generateResponse(prompt string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		start := time.Now()
		
		var response string
		var err error
		var toolResults []ai.ToolResult
		
		// Use tool-enabled generation if provider supports it
		if m.provider.SupportsTools() {
			response, toolResults, err = m.provider.GenerateResponseWithTools(ctx, prompt)
		} else {
			response, err = m.provider.GenerateResponse(ctx, prompt)
			toolResults = []ai.ToolResult{}
		}
		
		duration := time.Since(start)
		
		tokens := 0
		if err == nil {
			tokens = len(strings.Fields(response))
		}
		
		return msgResponse{
			response:    response,
			err:         err,
			tokens:      tokens,
			duration:    duration,
			toolResults: toolResults,
		}
	}
}

func (m *Model) View() string {
	if !m.ready {
		return ""
	}

	var s strings.Builder

	if m.loading {
		elapsed := time.Since(m.requestStart)
		loadingStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Italic(true)
		
		loadingText := fmt.Sprintf("\r🤔 AI is thinking... (elapsed: %s)", elapsed.Round(time.Millisecond*100))
		if m.totalRequests > 0 {
			avgTime := m.totalTime / time.Duration(m.totalRequests)
			loadingText += fmt.Sprintf(" | Session: %d requests, %d tokens, avg: %s", 
				m.totalRequests, m.totalTokens, avgTime.Round(time.Millisecond))
		}
		
		s.WriteString(loadingStyle.Render(loadingText))
	} else {
		inputStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("205"))
		s.WriteString(inputStyle.Render("> "))
		s.WriteString(m.input)
	}

	return s.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}