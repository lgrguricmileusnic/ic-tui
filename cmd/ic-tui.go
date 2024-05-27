package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/lgrguricmileusnic/ic-tui/internal/styles"
	"github.com/lgrguricmileusnic/ic-tui/pkg/bubbles/blinkers"
)

const (
	padding  = 2
	maxWidth = 80
	maxSpeed = 250.00
)

type updatePostData struct {
	Speed        float64
	Blinkers     bool
	WinCondition bool
}
type responseMsg struct {
	speed    float64
	blinkers bool
}

type WinMsg struct{}

func listenForActivity(sub chan updatePostData) tea.Cmd {
	return func() tea.Msg {
		for {
			mux := http.NewServeMux()
			mux.HandleFunc("POST /update", func(w http.ResponseWriter, r *http.Request) {

				var data updatePostData

				err := json.NewDecoder(r.Body).Decode(&data)

				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				sub <- data
			})
			http.ListenAndServe(":8080", mux)
		}
	}
}

func waitForActivity(sub chan updatePostData) tea.Cmd {
	return func() tea.Msg {
		data := updatePostData(<-sub)

		if data.WinCondition {
			return WinMsg{}
		}

		return responseMsg{data.Speed, data.Blinkers}
	}
}

type Window struct {
	width  int
	heigth int
}
type model struct {
	sub         chan updatePostData
	blinkers    blinkers.Model
	displayFlag bool
	speedbar    progress.Model
	speed       float64
	window      Window
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.blinkers.Init(),
		listenForActivity(m.sub),
		waitForActivity(m.sub),
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
		m.window.heigth = msg.Height
		m.window.width = msg.Width
		return m, nil

	case responseMsg:
		m.speed = msg.speed
		scmd := m.speedbar.SetPercent(msg.speed / maxSpeed)
		bcmd := m.blinkers.SetBlinking(msg.blinkers)
		return m, tea.Batch(waitForActivity(m.sub), scmd, bcmd)

		// win condition msg
	case WinMsg:
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

	default:
		return m, nil
	}
}

func (m model) View() string {

	sm := m.speedbar.View()
	ic := styles.IcStyle.Render(lipgloss.JoinVertical(lipgloss.Center,
		m.blinkers.View(),
		"\n",
		sm,
		fmt.Sprintf("%.f km/h", m.speed)))
	s := ic
	return lipgloss.Place(m.window.width, m.window.heigth, lipgloss.Center, lipgloss.Center, s)
}

func main() {
	// Progress model init
	pm := progress.New(progress.WithSolidFill("#FF2800"))
	pm.ShowPercentage = false

	// Blinkers model init

	bm := blinkers.New()

	p := tea.NewProgram(model{
		sub:      make(chan updatePostData),
		blinkers: bm,
		speedbar: pm},
		tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
		os.Exit(1)
	}
}
