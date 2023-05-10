package game

import (
	"context"
	"fmt"
	"github.com/BigJk/end_of_eden/lua/luhelp"
	"github.com/labstack/echo/v4"
	"github.com/olahol/melody"
	"github.com/samber/lo"
	lua "github.com/yuin/gopher-lua"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// ExposeDebug exposes a debug interface on the given port. This interface can be used to execute lua code on the server.
// This is a very dangerous function, which should only be used for debugging purposes. It should never be exposed to the public.
func ExposeDebug(port int, session *Session, l *lua.LState, log *log.Logger) func() error {
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
		_ = session.Write([]byte("::: Welcome to the End of Eden REPL        :::"))
		_ = session.Write([]byte("::: Use debug_r(args...) to send data back :::"))

		log.Println("Debug connected:", session.RemoteAddr())
	})

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "End Of Eden Debug Server")
	})

	e.GET("/state", func(c echo.Context) error {
		return c.JSONPretty(http.StatusOK, session.GetGameState(), "\t")
	})

	e.GET("/fight", func(c echo.Context) error {
		return c.JSONPretty(http.StatusOK, session.GetFight(), "\t")
	})

	e.GET("/merchant", func(c echo.Context) error {
		return c.JSONPretty(http.StatusOK, session.GetMerchant(), "\t")
	})

	e.GET("/actor/:guid", func(c echo.Context) error {
		return c.JSONPretty(http.StatusOK, session.GetActor(c.Param("guid")), "\t")
	})

	e.GET("/instance/:guid", func(c echo.Context) error {
		return c.JSONPretty(http.StatusOK, session.GetInstance(c.Param("guid")), "\t")
	})

	e.GET("/actors", func(c echo.Context) error {
		return c.JSONPretty(http.StatusOK, lo.Map(session.GetActors(), func(guid string, _ int) Actor {
			return session.GetActor(guid)
		}), "\t")
	})

	e.POST("/exec", func(c echo.Context) error {
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

	e.GET("/instances", func(c echo.Context) error {
		return c.JSONPretty(http.StatusOK, lo.Map(session.GetInstances(), func(guid string, _ int) any {
			return session.GetInstance(guid)
		}), "\t")
	})

	e.GET("svg", func(c echo.Context) error {
		svg, _, err := session.ToSVG()
		if err != nil {
			return c.JSONPretty(http.StatusBadRequest, err.Error(), "\t")
		}
		return c.Blob(http.StatusOK, "image/svg+xml", svg)
	})

	e.GET("d2", func(c echo.Context) error {
		_, diag, err := session.ToSVG()
		if err != nil {
			return c.JSONPretty(http.StatusBadRequest, err.Error(), "\t")
		}
		return c.String(http.StatusOK, diag)
	})

	go func() {
		e.StdLogger = log
		e.HideBanner = true
		if err := e.Start(fmt.Sprintf("127.0.0.1:%d", port)); err != nil && err != http.ErrServerClosed {
			log.Fatal("shutting down the server")
		}
	}()

	// Open the REPL in browser if ttyd and wscat are available.
	var ttydInstance *exec.Cmd
	go func() {
		ttyd, err := exec.LookPath("ttyd")
		if err != nil {
			return
		}

		wscat, err := exec.LookPath("wscat")
		if err != nil {
			return
		}

		time.Sleep(time.Second * 1)

		log, _ := os.OpenFile(fmt.Sprintf("./logs/ttyd-%d.log", port+1), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
		defer log.Close()

		ttydInstance = exec.Command(ttyd, "--port", "0", "--browser", wscat, fmt.Sprintf("--connect=ws://127.0.0.1:%d/ws", port))
		ttydInstance.Stdout = log
		ttydInstance.Stderr = log

		_ = ttydInstance.Run()
	}()

	return func() error {
		// Kill ttyd
		if ttydInstance != nil && ttydInstance.Process != nil {
			_ = ttydInstance.Process.Kill()
		}

		// Shutdown server
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return e.Shutdown(ctx)
	}
}
