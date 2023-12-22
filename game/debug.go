package game

import (
	"context"
	"fmt"
	luhelp2 "github.com/BigJk/end_of_eden/internal/lua/luhelp"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/olahol/melody"
	"github.com/samber/lo"
	lua "github.com/yuin/gopher-lua"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func setupVite(e *echo.Echo) {
	viteUrl, err := url.Parse("http://127.0.0.1:3000")
	if err != nil {
		panic(err)
	}

	e.Use(middleware.ProxyWithConfig(middleware.ProxyConfig{Skipper: func(c echo.Context) bool {
		return strings.HasPrefix(c.Request().URL.Path, "/api") ||
			strings.HasPrefix(c.Request().URL.Path, "/ws")
	}, Balancer: middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{{URL: viteUrl}})}))
}

func setupStaticDir(e *echo.Echo) {
	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}
	e.Static("/", filepath.Join(filepath.Dir(exe), "debug/dist"))
}

// ExposeDebug exposes a debug interface on the given port. This interface can be used to execute lua code on the server.
// This is a very dangerous function, which should only be used for debugging purposes. It should never be exposed to the public.
func ExposeDebug(port int, session *Session, l *lua.LState, log *log.Logger) func() error {
	e := echo.New()
	mtx := sync.Mutex{}
	m := melody.New()
	mapper := luhelp2.NewMapper(l)

	e.GET("/ws", func(c echo.Context) error {
		_ = m.HandleRequest(c.Response().Writer, c.Request())
		return nil
	})

	l.SetGlobal("debug_r", l.NewFunction(func(state *lua.LState) int {
		_ = m.Broadcast([]byte(strings.Join(lo.Map(make([]any, state.GetTop()), func(_ any, index int) string {
			val := state.Get(1 + index)
			return luhelp2.ToString(val, mapper)
		}), " ")))
		return 0
	}))

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		mtx.Lock()
		defer mtx.Unlock()

		if _, err := l.LoadString(string(msg)); err != nil {
			_ = s.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
		} else {
			if err := l.DoString(string(msg)); err != nil {
				_ = s.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
			}
		}
	})

	m.HandleConnect(func(session *melody.Session) {
		_ = session.Write([]byte("::: Welcome to the End of Eden REPL        :::"))
		_ = session.Write([]byte("::: Use debug_r(args...) to send data back :::"))

		log.Println("Debug connected:", session.RemoteAddr())
	})

	api := e.Group("/api")

	api.GET("/state", func(c echo.Context) error {
		return c.JSONPretty(http.StatusOK, session.GetGameState(), "\t")
	})

	api.GET("/fight", func(c echo.Context) error {
		return c.JSONPretty(http.StatusOK, session.GetFight(), "\t")
	})

	api.GET("/merchant", func(c echo.Context) error {
		return c.JSONPretty(http.StatusOK, session.GetMerchant(), "\t")
	})

	api.GET("/actor/:guid", func(c echo.Context) error {
		return c.JSONPretty(http.StatusOK, session.GetActor(c.Param("guid")), "\t")
	})

	api.GET("/instance/:guid", func(c echo.Context) error {
		return c.JSONPretty(http.StatusOK, session.GetInstance(c.Param("guid")), "\t")
	})

	api.GET("/actors", func(c echo.Context) error {
		return c.JSONPretty(http.StatusOK, lo.Map(session.GetActors(), func(guid string, _ int) Actor {
			return session.GetActor(guid)
		}), "\t")
	})

	api.POST("/exec", func(c echo.Context) error {
		lua, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return c.JSONPretty(http.StatusBadRequest, err.Error(), "\t")
		}

		mtx.Lock()
		defer mtx.Unlock()

		if _, err := l.LoadString(string(lua)); err != nil {
			return c.JSONPretty(http.StatusBadRequest, err.Error(), "\t")
		} else {
			if err := l.DoString(string(lua)); err != nil {
				return c.JSONPretty(http.StatusBadRequest, err.Error(), "\t")
			}
		}

		return c.NoContent(http.StatusOK)
	})

	api.GET("/registered/enemies", func(c echo.Context) error {
		return c.JSONPretty(http.StatusOK, session.resources.Enemies, "\t")
	})

	api.GET("/registered/cards", func(c echo.Context) error {
		return c.JSONPretty(http.StatusOK, session.resources.Cards, "\t")
	})

	api.GET("/registered/events", func(c echo.Context) error {
		return c.JSONPretty(http.StatusOK, session.resources.Events, "\t")
	})

	api.GET("/registered/status_effects", func(c echo.Context) error {
		return c.JSONPretty(http.StatusOK, session.resources.StatusEffects, "\t")
	})

	api.GET("/registered/artifacts", func(c echo.Context) error {
		return c.JSONPretty(http.StatusOK, session.resources.Artifacts, "\t")
	})

	api.GET("/registered/story_tellers", func(c echo.Context) error {
		return c.JSONPretty(http.StatusOK, session.resources.StoryTeller, "\t")
	})

	api.GET("/instances", func(c echo.Context) error {
		return c.JSONPretty(http.StatusOK, lo.Map(session.GetInstances(), func(guid string, _ int) any {
			return session.GetInstance(guid)
		}), "\t")
	})

	api.GET("/svg", func(c echo.Context) error {
		svg, _, err := session.ToSVG()
		if err != nil {
			return c.JSONPretty(http.StatusBadRequest, err.Error(), "\t")
		}
		return c.Blob(http.StatusOK, "image/svg+xml", svg)
	})

	api.GET("/d2", func(c echo.Context) error {
		_, diag, err := session.ToSVG()
		if err != nil {
			return c.JSONPretty(http.StatusBadRequest, err.Error(), "\t")
		}
		return c.String(http.StatusOK, diag)
	})

	setupStaticDir(e)

	if os.Getenv("EOE_VITE") == "1" {
		setupVite(e)
	}

	go func() {
		e.StdLogger = log
		e.HideBanner = true
		if err := e.Start(fmt.Sprintf("127.0.0.1:%d", port)); err != nil && err != http.ErrServerClosed {
			log.Fatal("shutting down the server")
		}
	}()

	return func() error {
		// Shutdown server
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return e.Shutdown(ctx)
	}
}
