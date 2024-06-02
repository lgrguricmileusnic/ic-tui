package program

import (
	"fmt"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lgrguricmileusnic/ic-tui/internal/api"

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
type Model struct {
	ApiAddr     string
	Flag        string
	displayFlag bool

	Sub      chan api.UpdatePostData
	Statdash statdash.Model
	Blinkers blinkers.Model
	Speedbar progress.Model
	speed    float64

	window Window
	styles styles
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.Blinkers.Init(),
		api.ListenForActivity(m.Sub, m.ApiAddr),
		api.WaitForActivity(m.Sub),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Speedbar.Width = msg.Width - padding*2 - 4
		if m.Speedbar.Width > maxWidth {
			m.Speedbar.Width = maxWidth
		}
		m.Blinkers.Width = msg.Width
		m.Statdash.Width = msg.Width
		m.window.heigth = msg.Height
		m.window.width = msg.Width
		return m, nil

	case api.UpdateMsg:
		m.speed = msg.Speed
		scmd := m.Speedbar.SetPercent(msg.Speed / maxSpeed)
		bcmd := m.Blinkers.SetBlinking(msg.Blinkers)
		sdcmd := m.Statdash.SetLedStatus(msg.Seatbelt, msg.Engine, msg.Battery, msg.Doors, msg.Oil)
		return m, tea.Batch(api.WaitForActivity(m.Sub), scmd, bcmd, sdcmd)

	// win condition msg
	case api.WinMsg:
		m.displayFlag = true
		return m, nil

	// progress messages
	case progress.FrameMsg:
		progressModel, cmd := m.Speedbar.Update(msg)
		m.Speedbar = progressModel.(progress.Model)
		return m, cmd

	// blinkers messages
	case blinkers.TickMsg:
		var cmd tea.Cmd
		m.Blinkers, cmd = m.Blinkers.Update(msg)
		return m, cmd

	case blinkers.OnOffMsg:
		var cmd tea.Cmd
		m.Blinkers, cmd = m.Blinkers.Update(msg)
		return m, cmd

	// status dashboard messages
	case statdash.LedStatusMsg:
		var cmd tea.Cmd
		m.Statdash, cmd = m.Statdash.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
}

func (m Model) View() string {
	if m.displayFlag {
		return lipgloss.Place(m.window.width, m.window.heigth, lipgloss.Center, lipgloss.Center, fmt.Sprintf("The flag is \"%s\"", m.Flag))
	}
	sm := m.Speedbar.View()
	sd := m.Statdash.View()
	ic := m.styles.icStyle.Render(lipgloss.JoinVertical(lipgloss.Center,
		m.Blinkers.View(),
		"\n",
		sm,
		fmt.Sprintf("%.f km/h\n", m.speed)),
		"\n",
		sd,
	)
	s := ic
	return lipgloss.Place(m.window.width, m.window.heigth, lipgloss.Center, lipgloss.Center, s)
}

func New(addr string, flag string, r *lipgloss.Renderer) Model {

	pm := progress.New(progress.WithColorProfile(r.ColorProfile()), progress.WithSolidFill("#FFC300"))
	pm.ShowPercentage = false

	// Blinkers model init
	bm := blinkers.New(blinkers.WithRenderer(r))

	// Status Dashboard model init
	sm := statdash.New(statdash.WithRenderer(r))

	m := Model{
		Flag:     flag,
		ApiAddr:  addr,
		Sub:      make(chan api.UpdatePostData),
		Blinkers: bm,
		Speedbar: pm,
		Statdash: sm,
		styles:   makeStyles(r)}
	return m
}
