package blinkers

import "github.com/charmbracelet/lipgloss"

const (
	left  string = "🡄"
	right string = "🡆"
)

type styles struct {
	on  lipgloss.Style
	off lipgloss.Style
}

func makeStyles(r *lipgloss.Renderer) styles {
	return styles{
		on:  r.NewStyle().Foreground(lipgloss.Color("#33FF57")).Align(lipgloss.Center),
		off: r.NewStyle().Foreground(lipgloss.Color("#808080")).Align(lipgloss.Center),
	}

}
