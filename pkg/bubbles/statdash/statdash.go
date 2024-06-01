package statdash

import (
	"strings"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	lastID int
	idMtx  sync.Mutex
)

func nextID() int {
	idMtx.Lock()
	defer idMtx.Unlock()
	lastID++
	return lastID
}

var leds = [5]string{"BELT", "ENG", "BAT", "DOOR", "OIL"}

type Model struct {
	ID     int
	Width  int
	status [5]bool
	styles styles
}

type LedStatusMsg struct {
	ID       int
	seatbelt bool
	engine   bool
	battery  bool
	doors    bool
	oil      bool
}

func (m Model) SetLedStatus(seatbelt bool, engine bool, battery bool, doors bool, oil bool) tea.Cmd {
	return func() tea.Msg {
		return LedStatusMsg{m.ID, seatbelt, engine, battery, doors, oil}
	}
}

func (m Model) View() string {
	rLeds := make([]string, 5)
	for i, led := range leds {
		if m.status[i] {
			rLeds[i] = m.styles.on.Render(led)
		} else {
			rLeds[i] = m.styles.off.Render(led)
		}
	}
	return m.styles.dash.Render(strings.Join(rLeds, strings.Repeat(" ", 6)))
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case LedStatusMsg:
		if msg.ID != m.ID {
			return m, nil
		}
		m.status[0] = msg.seatbelt
		m.status[1] = msg.engine
		m.status[2] = msg.battery
		m.status[3] = msg.doors
		m.status[4] = msg.oil

		return m, nil
	default:
		return m, nil
	}
}

type Option func(*Model)

func WithRenderer(r *lipgloss.Renderer) Option {
	return func(m *Model) {
		m.styles = makeStyles(r)
	}
}
func New(opts ...Option) Model {
	m := Model{
		ID:     nextID(),
		Width:  21,
		styles: makeStyles(lipgloss.DefaultRenderer()),
	}
	for _, opt := range opts {
		opt(&m)
	}
	return m
}
