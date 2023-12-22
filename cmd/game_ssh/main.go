package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/BigJk/end_of_eden/system/settings"
	"github.com/BigJk/end_of_eden/ui/menus/mainmenu"
	"github.com/BigJk/end_of_eden/ui/root"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"
	"github.com/muesli/termenv"
)

var mtx sync.Mutex
var instanceLimit int
var instances int

func main() {
	fmt.Println("End Of Eden :: SSH Server")
	fmt.Println("Each SSH session creates it's own game session. Modding and audio not supported.")
	fmt.Println()

	bind := flag.String("bind", ":8273", "ip and port to bind to")
	timeout := flag.Int("timeout", 0, "ssh idle timeout")
	maxInstance := flag.Int("max_inst", 10, "maximum of game instances")
	help := flag.Bool("help", false, "show help")
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		return
	}

	if err := os.Setenv("NO_CLIPBOARD", "true"); err != nil {
		panic(err)
	}

	termenv.SetDefaultOutput(termenv.NewOutput(os.Stdout, termenv.WithProfile(termenv.TrueColor)))

	options := []ssh.Option{
		wish.WithAddress(*bind),
		wish.WithHostKeyPath(".ssh/term_info_ed25519"),
		wish.WithMiddleware(
			func(handler ssh.Handler) ssh.Handler {
				return func(session ssh.Session) {
					mtx.Lock()
					instances -= 1
					mtx.Unlock()

					handler(session)
				}
			},
			gameMiddleware(),
			func(handler ssh.Handler) ssh.Handler {
				return func(session ssh.Session) {
					mtx.Lock()
					if instanceLimit > 0 && instances >= instanceLimit {
						mtx.Unlock()

						log.Warn("Denying instance because of limit!")
						_, _ = session.Write([]byte("Too many instances... Please try again later."))
						time.Sleep(time.Second * 2)
						_ = session.Close()

						return
					}
					instances += 1
					mtx.Unlock()

					handler(session)
				}
			},
			lm.Middleware(),
		),
	}

	if *timeout > 0 {
		options = append(options, wish.WithIdleTimeout(time.Duration(*timeout)*time.Minute))
	}

	instanceLimit = *maxInstance

	s, err := wish.NewServer(options...)
	if err != nil {
		log.Error("could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server:", "bind", *bind, "max_inst", *maxInstance, "timeout", *timeout)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("could not start server", "error", err)
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("could not stop server", "error", err)
	}
}

func gameMiddleware() wish.Middleware {
	newProg := func(m tea.Model, opts ...tea.ProgramOption) *tea.Program {
		p := tea.NewProgram(m, opts...)
		lipgloss.SetColorProfile(termenv.TrueColor)
		return p
	}
	teaHandler := func(s ssh.Session) *tea.Program {
		_, _, active := s.Pty()
		if !active {
			wish.Fatalln(s, "no active terminal, skipping")
			return nil
		}
		zones := zone.New()
		return newProg(
			root.New(zones, mainmenu.NewModel(zones, settings.GetGlobal(), nil, nil)),
			tea.WithInput(s),
			tea.WithOutput(s),
			tea.WithANSICompressor(),
			tea.WithMouseAllMotion(),
			tea.WithAltScreen(),
		)
	}
	return bm.MiddlewareWithProgramHandler(teaHandler, termenv.ANSI256)
}
