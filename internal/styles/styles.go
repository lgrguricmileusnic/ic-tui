package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	IcStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FFFFFF")).
		Padding(1, 1).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true).
		Align(lipgloss.Center)
)
