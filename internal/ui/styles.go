package ui

import "github.com/charmbracelet/lipgloss"

var (
	BaseColor = lipgloss.Color("#8be9fd") // Cyan
	GoldColor = lipgloss.Color("#f1fa8c") // Yellow/Gold
	RedColor  = lipgloss.Color("#ff5555") // Red
	GreenColor = lipgloss.Color("#50fa7b") // Green
	
	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#bd93f9")). // Purple
			Bold(true).
			MarginBottom(1)

	PadStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#6272a4")). // Muted blue
			Padding(1, 2).
			Margin(0, 1).
			Width(20).
			Height(5).
			Align(lipgloss.Center, lipgloss.Center)

	ActivePadStyle = PadStyle.Copy().
			BorderForeground(BaseColor).
			Bold(true)

	StatusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272a4")).
			Italic(true).
			MarginTop(1)

	KeyStyle = lipgloss.NewStyle().
			Foreground(BaseColor).
			Bold(true)
)
