package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/BigJk/end_of_eden/settings"
	"github.com/BigJk/end_of_eden/ui/menus/mainmenu"
	"github.com/BigJk/end_of_eden/ui/root"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	zone "github.com/lrstanley/bubblezone"
	"github.com/muesli/termenv"
	"github.com/olahol/melody"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

const CookieId = "id"

type session struct {
	prog   *tea.Program
	input  *ConcurrentRW
	output *ConcurrentRW
	wg     *sync.WaitGroup
	cancel func()
}

var sessionMtx = &sync.RWMutex{}
var sessions = map[string]*session{}

func main() {
	// Parse flags
	debug := flag.Bool("debug", false, "enable debug mode")
	flag.Parse()

	// Logger
	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.TimeOnly,
	})

	// Force true color
	termenv.SetDefaultOutput(termenv.NewOutput(os.Stdout, termenv.WithProfile(termenv.TrueColor)))
	lipgloss.SetColorProfile(termenv.TrueColor)

	e := echo.New()
	m := melody.New()

	e.HideBanner = true

	e.Static("/", "./dist")

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if ck, err := c.Request().Cookie(CookieId); err == nil {
				if ck.Valid() == nil {
					c.Set(CookieId, ck.Value)
					return next(c)
				}
			}

			id := fmt.Sprintf("%d-%d-%d", time.Now().Unix(), rand.Intn(99999999), rand.Intn(99999999))
			c.SetCookie(&http.Cookie{
				Name:    CookieId,
				Value:   id,
				Expires: time.Now().AddDate(1, 0, 0),
			})
			c.Set(CookieId, id)

			logger.Info("new id assigned", "id", id)

			return next(c)
		}
	})

	api := e.Group("/api")
	api.POST("/resize", func(c echo.Context) error {
		id := c.Get(CookieId).(string)

		sessionMtx.Lock()
		defer sessionMtx.Unlock()

		if _, ok := sessions[id]; ok {
			return c.String(http.StatusOK, "session already exists")
		}

		type resize struct {
			Cols int `json:"cols"`
			Rows int `json:"rows"`
		}

		var r resize
		if err := c.Bind(&r); err != nil {
			return c.String(http.StatusBadRequest, "invalid resize")
		}

		sessions[id].prog.Send(tea.WindowSizeMsg{Width: r.Cols, Height: r.Rows})

		return c.String(http.StatusOK, "ok")
	})
	api.GET("/create", func(c echo.Context) error {
		id := c.Get(CookieId).(string)

		sessionMtx.Lock()
		defer sessionMtx.Unlock()

		if _, ok := sessions[id]; ok {
			return c.String(http.StatusOK, "session already exists")
		}

		var baseModel tea.Model
		zones := zone.New()
		baseModel = root.New(zones, mainmenu.NewModel(zones, settings.GetGlobal(), nil, nil))

		input := NewConcurrentRW()
		output := NewConcurrentRW()

		go input.Run()
		go output.Run()

		prog := tea.NewProgram(baseModel, tea.WithInput(input), tea.WithOutput(output), tea.WithAltScreen(), tea.WithMouseAllMotion(), tea.WithANSICompressor())
		wg := &sync.WaitGroup{}
		ctx, cancel := context.WithCancel(context.Background())

		wg.Add(1)
		go func() {
			defer wg.Done()

			if _, err := prog.Run(); err != nil {
				logger.Warn("alas, there's been an error", "err", err, "id", id)

				input.Close()
				output.Close()
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer output.Close()

			nextSend := make([]byte, 0, 2048+256)
			buff := make([]byte, 1024)
			ticker := time.NewTicker(time.Millisecond)
			for {
				select {
				case <-ctx.Done():
					logger.Info("session closed", "id", id)
					prog.Send(tea.Quit())
				case <-ticker.C:
					for {
						n, err := output.Read(buff)
						if err == io.EOF || len(nextSend) >= 2048 {
							_ = m.BroadcastFilter(nextSend, func(q *melody.Session) bool {
								return q.Keys["watching"].(string) == id
							})
							nextSend = nextSend[:0]
						} else if err != nil {
							logger.Warn("error reading output", "err", err, "id", id)
							return
						} else {
							nextSend = append(nextSend, buff[:n]...)
						}
					}
				}
			}
		}()

		session := &session{
			prog:   prog,
			input:  input,
			output: output,
			cancel: cancel,
			wg:     wg,
		}

		sessions[id] = session

		return c.String(http.StatusOK, "session created")
	})

	e.GET("/ws/:id", func(c echo.Context) error {
		return m.HandleRequestWithKeys(c.Response().Writer, c.Request(), map[string]interface{}{
			"watching": c.Param("id"),
			"id":       c.Get(CookieId).(string),
			"isOwner":  c.Param("id") == c.Get(CookieId).(string),
		})
	})

	m.HandleConnect(func(s *melody.Session) {
		logger.Info("connected", "id", s.Keys["id"], "watching", s.Keys["watching"], "isOwner", s.Keys["isOwner"])
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {

	})

	if *debug { // If debug is enabled pass frontend requests to vite dev server.
		viteUrl, err := url.Parse("http://127.0.0.1:3000")
		if err != nil {
			panic(err)
		}

		e.Use(middleware.ProxyWithConfig(middleware.ProxyConfig{Skipper: func(c echo.Context) bool {
			return strings.HasPrefix(c.Request().URL.Path, "/api") ||
				strings.HasPrefix(c.Request().URL.Path, "/ws")
		}, Balancer: middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{{URL: viteUrl}})}))
	}

	if err := e.Start(":8080"); err != nil {
		panic(err)
	}
}
