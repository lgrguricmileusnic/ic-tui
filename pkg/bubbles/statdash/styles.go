package statdash

import "github.com/charmbracelet/lipgloss"

type styles struct {
	on   lipgloss.Style
	off  lipgloss.Style
	dash lipgloss.Style
}

func makeStyles(r *lipgloss.Renderer) styles {
	return styles{
		on:  r.NewStyle().Foreground(lipgloss.Color("#FF3131")).Align(lipgloss.Center),
		off: r.NewStyle().Foreground(lipgloss.Color("#808080")).Align(lipgloss.Center),
		dash: r.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("#FFFFFF")).
			Padding(0, 0).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true).
			Align(lipgloss.Center),
	}
}
