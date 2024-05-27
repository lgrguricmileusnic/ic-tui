package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/lgrguricmileusnic/ic-tui/internal/api"
	"github.com/lgrguricmileusnic/ic-tui/internal/args"
	"github.com/lgrguricmileusnic/ic-tui/internal/program"
	"github.com/lgrguricmileusnic/ic-tui/pkg/bubbles/blinkers"
	"github.com/lgrguricmileusnic/ic-tui/pkg/bubbles/statdash"
)

const (
	sshHost = "localhost"
	sshPort = "1234"
)

func main() {
	cfg := args.ParseArgs()

	if cfg.Ssh {
		startWish(cfg)
	} else {
		startLocal(cfg)
	}
}

func InitTeaProgramModel(cfg args.Args) tea.Model {
	// Progress model init
	pm := progress.New(progress.WithSolidFill("#FFC300"))
	pm.ShowPercentage = false

	// Blinkers model init
	bm := blinkers.New()

	// Status Dashboard model init
	sm := statdash.New()

	m := program.Model{
		ApiAddr:  cfg.Addr,
		Sub:      make(chan api.UpdatePostData),
		Blinkers: bm,
		Speedbar: pm,
		Statdash: sm}
	return m
}

func teaHandlerWrapper(cfg args.Args) func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	return func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
		pm := InitTeaProgramModel(cfg)

		return pm, []tea.ProgramOption{tea.WithAltScreen()}
	}

}

func startWish(cfg args.Args) {
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(sshHost, sshPort)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandlerWrapper(cfg)),
			activeterm.Middleware(), // Bubble Tea apps usually require a PTY.
			logging.Middleware(),
		),
	)

	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", sshHost, "port", sshPort)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}

}

func startLocal(cfg args.Args) {
	pm := InitTeaProgramModel(cfg)
	p := tea.NewProgram(pm, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
		os.Exit(1)
	}
}
