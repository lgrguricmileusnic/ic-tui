package program

import (
	"github.com/charmbracelet/lipgloss"
)

type styles struct {
	icStyle lipgloss.Style
}

func makeStyles(r *lipgloss.Renderer) styles {
	return styles{
		icStyle: r.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FFFFFF")).
			Padding(1, 1).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true).
			Align(lipgloss.Center),
	}
}
