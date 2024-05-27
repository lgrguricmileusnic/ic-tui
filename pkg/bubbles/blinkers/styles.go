package blinkers

import "github.com/charmbracelet/lipgloss"

const (
	left  string = "🡄"
	right string = "🡆"
)

var (
	BlinkersActiveStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#33FF57")).Align(lipgloss.Center)
	BlinkersInactiveStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080")).Align(lipgloss.Center)
)
