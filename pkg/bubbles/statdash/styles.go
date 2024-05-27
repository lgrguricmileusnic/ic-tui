package statdash

import "github.com/charmbracelet/lipgloss"

var (
	LedOnStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF3131")).Align(lipgloss.Center)
	LedOffStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080")).Align(lipgloss.Center)
	StatDashStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("#FFFFFF")).
			Padding(0, 0).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true).
			Align(lipgloss.Center)
)
