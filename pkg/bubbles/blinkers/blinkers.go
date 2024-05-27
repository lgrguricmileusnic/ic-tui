package blinkers

import (
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
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

type Model struct {
	ID       int
	blinking bool
	on       bool
	Spacing  int
	Width    int
}

type TickMsg struct {
	ID int
}

type OnOffMsg struct {
	ID       int
	blinking bool
}

func (m Model) tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg{m.ID}
	})
}

func (m Model) SetBlinking(blinking bool) tea.Cmd {
	if m.blinking != blinking {
		return func() tea.Msg {
			return OnOffMsg{ID: m.ID, blinking: blinking}
		}
	}
	return nil
}
func (m *Model) On() tea.Cmd {
	return m.SetBlinking(true)
}

func (m *Model) Off() tea.Cmd {
	return m.SetBlinking(false)

}

func (m Model) Init() tea.Cmd {
	return m.tick()
}

func (m Model) View() string {
	s := left + strings.Repeat(" ", m.Spacing) + right
	if m.blinking && m.on {
		return BlinkersActiveStyle.Render(s)
	}
	return BlinkersInactiveStyle.Render(s)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case OnOffMsg:
		if msg.ID != m.ID {
			return m, nil
		}
		m.blinking = msg.blinking
		return m, nil
	case TickMsg:
		if msg.ID != m.ID {
			return m, nil
		}
		m.on = !m.on
		return m, m.tick()
	}

	return m, nil
}

func New() Model {
	return Model{
		ID:       nextID(),
		blinking: false,
		Spacing:  10,
		on:       false,
		Width:    12,
	}
}
