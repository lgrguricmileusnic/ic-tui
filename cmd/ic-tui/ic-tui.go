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

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/lgrguricmileusnic/ic-tui/internal/args"
	"github.com/lgrguricmileusnic/ic-tui/internal/program"
)

const (
	sshHost = "0.0.0.0"
	sshPort = "22"
)

func main() {
	cfg := args.ParseArgs()

	if cfg.Ssh {
		startWish(cfg)
	} else {
		startLocal(cfg)
	}
}

func teaHandlerWrapper(cfg args.Args) func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	return func(s ssh.Session) (tea.Model, []tea.ProgramOption) {

		r := bubbletea.MakeRenderer(s)
		pm := program.New(cfg.Addr, "ctf{wroooom}", r)

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
	pm := program.New(cfg.Addr, "ctf{wroooom}", lipgloss.DefaultRenderer())
	p := tea.NewProgram(pm, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
		os.Exit(1)
	}
}
