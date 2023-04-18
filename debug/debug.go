package debug

import (
	"context"
	"fmt"
	"github.com/BigJk/project_gonzo/luhelp"
	"github.com/olahol/melody"
	"github.com/samber/lo"
	lua "github.com/yuin/gopher-lua"
	"log"
	"net/http"
	"strings"
	"sync"
)

func Expose(bind string, l *lua.LState, log *log.Logger) func() error {
	srv := &http.Server{Addr: bind}
	mtx := sync.Mutex{}
	m := melody.New()
	mapper := luhelp.NewMapper(l)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		_ = m.HandleRequest(w, r)
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

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	return func() error {
		return srv.Shutdown(context.Background())
	}
}
