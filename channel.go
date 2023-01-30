package echo

import (
	"github.com/olahol/melody"
	"golang.org/x/exp/slog"
)

type channel struct {
	name      string
	sessions  map[*melody.Session]bool
	enter     chan *melody.Session
	leave     chan *melody.Session
	broadcast chan []byte
}

func newChannel(n string) *channel {
	return &channel{
		name:      n,
		sessions:  make(map[*melody.Session]bool),
		enter:     make(chan *melody.Session),
		leave:     make(chan *melody.Session),
		broadcast: make(chan []byte),
	}
}

func (c *channel) handleEnter(s *melody.Session) {
	c.enter <- s
}

func (c *channel) handleLeave(s *melody.Session) {
	c.leave <- s
}

func (c *channel) handleBroadcast(m []byte) {
	c.broadcast <- m
}

func (c *channel) run() {
	for {
		select {
		case s := <-c.enter:
			c.sessions[s] = true
		case s := <-c.leave:
			delete(c.sessions, s)
		case m := <-c.broadcast:
			for s := range c.sessions {
				err := s.Write(m)
				if err != nil {
					slog.Error("handleBroadcast failed", err)
				}
			}
		}
	}
}
