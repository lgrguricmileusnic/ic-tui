package api

import (
	"encoding/json"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
)

type UpdateMsg struct {
	Speed    float64
	Blinkers bool
	Seatbelt bool
	Engine   bool
	Battery  bool
	Doors    bool
	Oil      bool
}

type WinMsg struct{}

type UpdatePostData struct {
	WinCondition bool

	Speed    float64
	Blinkers bool
	Seatbelt bool
	Engine   bool
	Battery  bool
	Doors    bool
	Oil      bool
}

func ListenForActivity(sub chan UpdatePostData, addr string) tea.Cmd {
	return func() tea.Msg {
		for {
			mux := http.NewServeMux()
			mux.HandleFunc("POST /update", func(w http.ResponseWriter, r *http.Request) {

				var data UpdatePostData

				err := json.NewDecoder(r.Body).Decode(&data)

				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				sub <- data
			})
			http.ListenAndServe(addr, mux)
		}
	}
}

func WaitForActivity(sub chan UpdatePostData) tea.Cmd {
	return func() tea.Msg {
		d := UpdatePostData(<-sub)

		if d.WinCondition {
			return WinMsg{}
		}

		return UpdateMsg{d.Speed, d.Blinkers, d.Seatbelt, d.Engine, d.Battery, d.Doors, d.Oil}
	}
}
