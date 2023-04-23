package game

import (
	"context"
	"fmt"
	"github.com/BigJk/end_of_eden/luhelp"
	"github.com/labstack/echo/v4"
	"github.com/olahol/melody"
	"github.com/samber/lo"
	lua "github.com/yuin/gopher-lua"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

func ExposeDebug(bind string, session *Session, l *lua.LState, log *log.Logger) func() error {
	e := echo.New()
	mtx := sync.Mutex{}
	m := melody.New()
	mapper := luhelp.NewMapper(l)

	e.GET("/ws", func(c echo.Context) error {
		_ = m.HandleRequest(c.Response().Writer, c.Request())
		return nil
	})

	l.SetGlobal("debug_r", l.NewFunction(func(state *lua.LState) int {
		_ = m.Broadcast([]byte(strings.Join(lo.Map(make([]any, state.GetTop()), func(_ any, index int) string {
			val := state.Get(1 + index)
			return luhelp.ToString(val, mapper)
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
		log.Println("Debug connected:", session.RemoteAddr())
	})

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "end_of_eden debug server")
	})

	e.GET("/state", func(c echo.Context) error {
		return c.JSON(http.StatusOK, session.GetGameState())
	})

	e.GET("/fight", func(c echo.Context) error {
		return c.JSON(http.StatusOK, session.GetFight())
	})

	e.GET("/merchant", func(c echo.Context) error {
		return c.JSON(http.StatusOK, session.GetMerchant())
	})

	e.GET("/actor/:guid", func(c echo.Context) error {
		return c.JSON(http.StatusOK, session.GetActor(c.Param("guid")))
	})

	e.GET("/instance/:guid", func(c echo.Context) error {
		return c.JSON(http.StatusOK, session.GetInstance(c.Param("guid")))
	})

	e.GET("/actors", func(c echo.Context) error {
		return c.JSON(http.StatusOK, lo.Map(session.GetActors(), func(guid string, _ int) Actor {
			return session.GetActor(guid)
		}))
	})

	e.GET("/instances", func(c echo.Context) error {
		return c.JSON(http.StatusOK, lo.Map(session.GetInstances(), func(guid string, _ int) any {
			return session.GetInstance(guid)
		}))
	})

	e.GET("svg", func(c echo.Context) error {
		svg, _, err := session.ToSVG()
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		return c.Blob(http.StatusOK, "image/svg+xml", svg)
	})

	e.GET("d2", func(c echo.Context) error {
		_, diag, err := session.ToSVG()
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		return c.String(http.StatusOK, diag)
	})

	go func() {
		e.StdLogger = log
		e.HideBanner = true
		if err := e.Start(bind); err != nil && err != http.ErrServerClosed {
			log.Fatal("shutting down the server")
		}
	}()

	return func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return e.Shutdown(ctx)
	}
}
