package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/lgrguricmileusnic/ic-tui/internal/api"
	"github.com/lgrguricmileusnic/ic-tui/internal/styles"
	"github.com/lgrguricmileusnic/ic-tui/pkg/bubbles/blinkers"
	"github.com/lgrguricmileusnic/ic-tui/pkg/bubbles/statdash"
)

const (
	padding  = 2
	maxWidth = 80
	maxSpeed = 250.00
)

type Window struct {
	width  int
	heigth int
}
type model struct {
	sub         chan api.UpdatePostData
	statdash    statdash.Model
	blinkers    blinkers.Model
	displayFlag bool
	speedbar    progress.Model
	speed       float64
	window      Window
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.blinkers.Init(),
		api.ListenForActivity(m.sub, ":8080"),
		api.WaitForActivity(m.sub),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	case tea.WindowSizeMsg:
		m.speedbar.Width = msg.Width - padding*2 - 4
		if m.speedbar.Width > maxWidth {
			m.speedbar.Width = maxWidth
		}
		m.blinkers.Width = msg.Width
		m.statdash.Width = msg.Width
		m.window.heigth = msg.Height
		m.window.width = msg.Width
		return m, nil

	case api.UpdateMsg:
		m.speed = msg.Speed
		scmd := m.speedbar.SetPercent(msg.Speed / maxSpeed)
		bcmd := m.blinkers.SetBlinking(msg.Blinkers)
		sdcmd := m.statdash.SetLedStatus(msg.Seatbelt, msg.Engine, msg.Battery, msg.Doors, msg.Oil)
		return m, tea.Batch(api.WaitForActivity(m.sub), scmd, bcmd, sdcmd)

	// win condition msg
	case api.WinMsg:
		m.displayFlag = true
		return m, tea.Quit

	// progress messages
	case progress.FrameMsg:
		progressModel, cmd := m.speedbar.Update(msg)
		m.speedbar = progressModel.(progress.Model)
		return m, cmd

	// blinkers messages
	case blinkers.TickMsg:
		var cmd tea.Cmd
		m.blinkers, cmd = m.blinkers.Update(msg)
		return m, cmd

	case blinkers.OnOffMsg:
		var cmd tea.Cmd
		m.blinkers, cmd = m.blinkers.Update(msg)
		return m, cmd

	// status dashboard messages
	case statdash.LedStatusMsg:
		var cmd tea.Cmd
		m.statdash, cmd = m.statdash.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
}

func (m model) View() string {
	sm := m.speedbar.View()
	sd := m.statdash.View()
	ic := styles.IcStyle.Render(lipgloss.JoinVertical(lipgloss.Center,
		m.blinkers.View(),
		"\n",
		sm,
		fmt.Sprintf("%.f km/h\n", m.speed)),
		"\n",
		sd,
	)
	s := ic
	return lipgloss.Place(m.window.width, m.window.heigth, lipgloss.Center, lipgloss.Center, s)
}

func main() {
	// Progress model init
	pm := progress.New(progress.WithSolidFill("#FFC300"))
	pm.ShowPercentage = false

	// Blinkers model init
	bm := blinkers.New()

	// Status Dashboard model init

	sm := statdash.New()

	p := tea.NewProgram(model{
		sub:      make(chan api.UpdatePostData),
		blinkers: bm,
		speedbar: pm,
		statdash: sm},
		tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
		os.Exit(1)
	}
}
