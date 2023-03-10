package echo

import (
	"net/http"
	"sync"

	"github.com/olahol/melody"
	"golang.org/x/exp/slog"
)

type HttpHandleFunc func(*Echo) (string, func(http.ResponseWriter, *http.Request))

type Echo struct {
	rwmutex  *sync.RWMutex
	mux      *http.ServeMux
	melody   *melody.Melody
	channels map[string]*channel
}

func New() *Echo {
	return &Echo{
		rwmutex:  &sync.RWMutex{},
		mux:      http.NewServeMux(),
		melody:   melody.New(),
		channels: make(map[string]*channel),
	}
}

func (e *Echo) HandleConnect(fn func(*melody.Session)) {
	e.melody.HandleConnect(fn)
}

func (e *Echo) HandleDisconnect(fn func(*melody.Session)) {
	e.melody.HandleDisconnect(fn)
}

func (e *Echo) HandlePong(fn func(*melody.Session)) {
	e.melody.HandlePong(fn)
}

func (e *Echo) HandleMessage(fn func(*melody.Session, []byte)) {
	e.melody.HandleMessage(fn)
}

func (e *Echo) HandleMessageBinary(fn func(*melody.Session, []byte)) {
	e.melody.HandleMessageBinary(fn)
}

func (e *Echo) HandleSentMessage(fn func(*melody.Session, []byte)) {
	e.melody.HandleSentMessage(fn)
}

func (e *Echo) HandleSentMessageBinary(fn func(*melody.Session, []byte)) {
	e.melody.HandleSentMessageBinary(fn)
}

func (e *Echo) HandleError(fn func(*melody.Session, error)) {
	e.melody.HandleError(fn)
}

func (e *Echo) HandleClose(fn func(*melody.Session, int, string) error) {
	e.melody.HandleClose(fn)
}

func (e *Echo) channel(n string) *channel {
	e.rwmutex.Lock()
	defer e.rwmutex.Unlock()

	if c, ok := e.channels[n]; ok {
		return c
	}

	c := newChannel(n)
	go c.run()
	e.channels[n] = c
	return c
}

func (e *Echo) Enter(n string, s *melody.Session) {
	e.channel(n).enter(s)
}

func (e *Echo) Leave(n string, s *melody.Session) {
	e.channel(n).leave(s)
}

func (e *Echo) Broadcast(n string, m []byte) {
	e.channel(n).broadcast(m)
}

func (e *Echo) Close(n string) {
	e.channel(n).close()
	delete(e.channels, n)
}

func (e *Echo) HandleHttp(fn HttpHandleFunc) {
	e.mux.HandleFunc(fn(e))
}

func (e *Echo) Run(addr string) {
	e.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := e.melody.HandleRequest(w, r)
		if err != nil {
			slog.Error("HandleRequest error", err)
		}
	})

	err := http.ListenAndServe(addr, e.mux)
	if err != nil {
		slog.Error("ListenAndServe error", err)
	}
}
