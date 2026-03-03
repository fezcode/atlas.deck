package ui

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fezcode/atlas.deck/internal/model"
)

type Model struct {
	Deck     *model.Deck
	LastCmd  string
	Status   string
	Running  bool
	Width    int
	Height   int
}

func NewModel(deck *model.Deck) Model {
	status := "Ready"
	if deck == nil {
		status = "No deck.piml found in current or global directory."
	}
	return Model{
		Deck:   deck,
		Status: status,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		default:
			if m.Deck != nil {
				for _, pad := range m.Deck.Pads {
					if msg.String() == pad.Key {
						m.Status = fmt.Sprintf("Running: %s...", pad.Label)
						m.Running = true
						return m, m.runCommand(pad.Command)
					}
				}
			}
		}

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case commandFinishedMsg:
		m.Running = false
		if msg.err != nil {
			m.Status = fmt.Sprintf("Error: %v", msg.err)
		} else {
			m.Status = "Success"
		}
		return m, nil
	}

	return m, nil
}

type commandFinishedMsg struct{ err error }

func (m Model) runCommand(command string) tea.Cmd {
	return func() tea.Msg {
		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("powershell", "-Command", command)
		} else {
			cmd = exec.Command("sh", "-c", command)
		}
		
		err := cmd.Run()
		return commandFinishedMsg{err}
	}
}

func (m Model) View() string {
	if m.Deck == nil {
		return TitleStyle.Render("Atlas Deck") + "\n\n" + StatusStyle.Render(m.Status)
	}

	var s strings.Builder

	s.WriteString(TitleStyle.Render("🚀 " + m.Deck.Name))
	s.WriteString("\n\n")

	// Render Pads in a grid
	var pads []string
	for _, pad := range m.Deck.Pads {
		style := PadStyle
		switch pad.Color {
		case "gold":
			style = style.BorderForeground(GoldColor)
		case "cyan":
			style = style.BorderForeground(BaseColor)
		case "red":
			style = style.BorderForeground(RedColor)
		case "green":
			style = style.BorderForeground(GreenColor)
		}

		content := fmt.Sprintf("%s\n\n%s", KeyStyle.Render("["+pad.Key+"]"), pad.Label)
		pads = append(pads, style.Render(content))
	}

	for i := 0; i < len(pads); i += 3 {
		end := i + 3
		if end > len(pads) {
			end = len(pads)
		}
		s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, pads[i:end]...))
		s.WriteString("\n")
	}

	s.WriteString("\n")
	s.WriteString(StatusStyle.Render("Status: " + m.Status))
	s.WriteString("\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#44475a")).Render("Press 'q' to quit"))

	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, s.String())
}
