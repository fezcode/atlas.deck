package ui

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fezcode/atlas.deck/internal/model"
)

type Model struct {
	Deck     *model.Deck
	Status   string
	Running  bool
	Width    int
	Height   int
	LastKey  string

	RunningLabel string
	RunningCmd   string
	Cmd          *exec.Cmd

	// Logs
	Viewport   viewport.Model
	Logs       []string
	Spinner    spinner.Model
	OutputChan chan string
	FinishChan chan error
}

func NewModel(deck *model.Deck) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#bd93f9"))

	status := "Ready"
	if deck == nil {
		status = "No deck.piml found."
	}

	vp := viewport.New(0, 0)
	vp.Style = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(lipgloss.Color("#44475a")).
		Padding(0, 1)

	return Model{
		Deck:       deck,
		Status:     status,
		Spinner:    s,
		Viewport:   vp,
		LastKey:    "None",
		OutputChan: make(chan string),
		FinishChan: make(chan error),
	}
}

func (m Model) Init() tea.Cmd {
	return m.Spinner.Tick
}

type outputMsg string
type finishMsg struct{ err error }
type commandStartedMsg struct{ cmd *exec.Cmd }

func waitForOutput(ch chan string) tea.Cmd {
	return func() tea.Msg {
		return outputMsg(<-ch)
	}
}

func waitForFinish(ch chan error) tea.Cmd {
	return func() tea.Msg {
		return finishMsg{err: <-ch}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			if m.Running && m.Cmd != nil && m.Cmd.Process != nil {
				m.Cmd.Process.Kill()
			}
			return m, tea.Quit
		case "ctrl+x": // Kill running process
			if m.Running && m.Cmd != nil && m.Cmd.Process != nil {
				if runtime.GOOS == "windows" {
					// Kill process tree on Windows
					exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprintf("%d", m.Cmd.Process.Pid)).Run()
				} else {
					m.Cmd.Process.Kill()
				}
				m.Status = "Killed: " + m.RunningLabel
				m.Logs = append(m.Logs, lipgloss.NewStyle().Foreground(RedColor).Bold(true).Render(fmt.Sprintf(">>> [KILL] Terminated: %s", m.RunningCmd)))
				m.Viewport.SetContent(strings.Join(m.Logs, "\n"))
				m.Viewport.GotoBottom()
			}
			return m, nil
		case "ctrl+l": // Clear logs
			m.Logs = []string{}
			if m.Running {
				logLine := lipgloss.NewStyle().Foreground(lipgloss.Color("#bd93f9")).Render(fmt.Sprintf(">>> [%s] Executing: %s", m.LastKey, m.RunningCmd))
				m.Logs = append(m.Logs, logLine)
			} else {
				m.Status = "Logs cleared"
			}
			m.Viewport.SetContent(strings.Join(m.Logs, "\n"))
			return m, nil
		default:
			if !m.Running && m.Deck != nil {
				for _, pad := range m.Deck.Pads {
					if msg.String() == pad.Key {
						m.LastKey = pad.Key
						m.Running = true
						m.RunningLabel = pad.Label
						m.RunningCmd = pad.Command
						m.Status = fmt.Sprintf("Running: %s", pad.Label)

						// Add execution log
						logLine := lipgloss.NewStyle().Foreground(lipgloss.Color("#bd93f9")).Render(fmt.Sprintf(">>> [%s] Executing: %s", pad.Key, pad.Command))
						m.Logs = append(m.Logs, logLine)
						m.Viewport.SetContent(strings.Join(m.Logs, "\n"))
						m.Viewport.GotoBottom()

						return m, tea.Batch(
							m.runCommand(pad.Command),
							waitForOutput(m.OutputChan),
							waitForFinish(m.FinishChan),
						)
					}
				}
			}
		}

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.Viewport.Width = msg.Width - 4
		m.Viewport.Height = 10

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd

	case commandStartedMsg:
		m.Cmd = msg.cmd
		return m, nil

	case outputMsg:
		m.Logs = append(m.Logs, string(msg))
		m.Viewport.SetContent(strings.Join(m.Logs, "\n"))
		m.Viewport.GotoBottom()
		return m, waitForOutput(m.OutputChan)

	case finishMsg:
		m.Running = false
		m.Cmd = nil
		if msg.err != nil {
			// Ignore error if it was a kill signal
			if !strings.Contains(msg.err.Error(), "exit status 1") && !strings.Contains(msg.err.Error(), "killed") {
				m.Status = fmt.Sprintf("Error: %v", msg.err)
				m.Logs = append(m.Logs, lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5555")).Render(fmt.Sprintf("!!! Error: %v", msg.err)))
			} else if m.Status != "Process terminated by user" {
				m.Status = "Completed with exit code"
			}
		} else {
			m.Status = "Completed"
			m.Logs = append(m.Logs, lipgloss.NewStyle().Foreground(lipgloss.Color("#50fa7b")).Render("--- Finished Successfully ---"))
		}

		m.Viewport.SetContent(strings.Join(m.Logs, "\n"))
		m.Viewport.GotoBottom()
		return m, nil
	}

	var vpCmd tea.Cmd
	m.Viewport, vpCmd = m.Viewport.Update(msg)
	cmds = append(cmds, vpCmd)

	return m, tea.Batch(cmds...)
}

func (m Model) runCommand(command string) tea.Cmd {
	return func() tea.Msg {
		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("powershell", "-NoProfile", "-Command", command)
		} else {
			cmd = exec.Command("sh", "-c", command)
		}

		stdout, _ := cmd.StdoutPipe()
		stderr, _ := cmd.StderrPipe()
		multi := io.MultiReader(stdout, stderr)

		if err := cmd.Start(); err != nil {
			m.FinishChan <- err
			return nil
		}

		go func() {
			scanner := bufio.NewScanner(multi)
			for scanner.Scan() {
				m.OutputChan <- scanner.Text()
			}
			m.FinishChan <- cmd.Wait()
		}()

		return commandStartedMsg{cmd: cmd}
	}
}

func (m Model) View() string {
	if m.Deck == nil {
		return TitleStyle.Render("Atlas Deck") + "\n\n" + StatusStyle.Render(m.Status)
	}

	var s strings.Builder

	// Header
	header := "🚀 " + m.Deck.Name
	if m.Running {
		header += " " + m.Spinner.View()
	}
	s.WriteString(TitleStyle.Render(header))
	s.WriteString("\n\n")

	// Render Pads in a responsive grid
	var pads []string
	for _, pad := range m.Deck.Pads {
		style := PadStyle
		if m.Running && pad.Key == m.LastKey {
			style = ActivePadStyle
		} else {
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
		}

		content := fmt.Sprintf("%s\n\n%s", KeyStyle.Render("["+pad.Key+"]"), pad.Label)
		pads = append(pads, style.Render(content))
	}

	// Dynamic grid
	padWidth := 24 
	cols := m.Width / padWidth
	if cols < 1 {
		cols = 1
	}

	for i := 0; i < len(pads); i += cols {
		end := i + cols
		if end > len(pads) {
			end = len(pads)
		}
		s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, pads[i:end]...))
		s.WriteString("\n")
	}

	// Footer / Logs
	s.WriteString("\n")
	statusLine := fmt.Sprintf("Last Key: [%s] • Status: %s", m.LastKey, m.Status)
	if m.Running {
		statusLine = fmt.Sprintf("Last Key: [%s] • %s: %s", m.LastKey, lipgloss.NewStyle().Foreground(BaseColor).Bold(true).Render("Running"), m.RunningCmd)
	}
	s.WriteString(StatusStyle.Render(statusLine))
	s.WriteString("\n")
	s.WriteString(m.Viewport.View())
	
	footerHints := "Press ctrl+c to quit • ctrl+l to clear logs"
	if m.Running {
		footerHints += " • " + lipgloss.NewStyle().Foreground(RedColor).Bold(true).Render("ctrl+x to kill process")
	}
	s.WriteString("\n" + StatusStyle.Render(footerHints))

	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, s.String())
}
