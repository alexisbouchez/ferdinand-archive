package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	s, err := wish.NewServer(
		wish.WithAddress(os.Getenv("SSH_ADDR")),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.SetReportTimestamp(false)

	go func() {
		log.Info("Starting SSH server", "addr", os.Getenv("SSH_ADDR"))
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Username").Key("username"),
			huh.NewInput().Title("Password").Key("password").EchoMode(huh.EchoModePassword),
		),
	)
	r := bubbletea.MakeRenderer(s)
	style := r.NewStyle().
		Border(lipgloss.NormalBorder()).
		Padding(1, 2).
		BorderForeground(lipgloss.Color("#444444")).
		Foreground(lipgloss.Color("#7571F9"))

	m := model{form: form, style: style}
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

type model struct {
	form     *huh.Form
	style    lipgloss.Style
	loggedIn bool
}

func (m model) Init() tea.Cmd {
	if m.form == nil {
		return nil
	}
	return m.form.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if m.form != nil {
		f, cmd := m.form.Update(msg)
		m.form = f.(*huh.Form)
		cmds = append(cmds, cmd)
	}

	m.loggedIn = m.form.State == huh.StateCompleted
	if m.form.State == huh.StateAborted {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.form == nil {
		return "Starting..."
	}
	if m.loggedIn {
		return m.style.Render(
			"Welcome, " + m.form.GetString("username") + "!" +
				" Your password is " + m.form.GetString("password") + ".",
		)
	}
	return m.form.View()
}
