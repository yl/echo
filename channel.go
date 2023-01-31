package echo

import (
	"github.com/olahol/melody"
	"golang.org/x/exp/slog"
)

type channel struct {
	name       string
	closed     bool
	sessions   map[*melody.Session]bool
	enterC     chan *melody.Session
	leaveC     chan *melody.Session
	closeC     chan bool
	broadcastC chan []byte
}

func newChannel(n string) *channel {
	return &channel{
		name:       n,
		closed:     false,
		sessions:   make(map[*melody.Session]bool),
		enterC:     make(chan *melody.Session),
		leaveC:     make(chan *melody.Session),
		closeC:     make(chan bool),
		broadcastC: make(chan []byte),
	}
}

func (c *channel) enter(s *melody.Session) {
	c.enterC <- s
}

func (c *channel) leave(s *melody.Session) {
	c.leaveC <- s
}

func (c *channel) close() {
	c.closeC <- true
}

func (c *channel) broadcast(m []byte) {
	c.broadcastC <- m
}

func (c *channel) handleEnter(s *melody.Session) {
	s.Set(c.name, true)
	c.sessions[s] = true
}

func (c *channel) handleLeave(s *melody.Session) {
	s.Set(c.name, false)
	delete(c.sessions, s)
}

func (c *channel) handleBroadcast(m []byte) {
	for s := range c.sessions {
		err := s.Write(m)
		if err != nil {
			slog.Error("Message broadcast error", err, "channel", c.name, "message", m)
		}
	}
}

func (c *channel) handleClose() {
	if c.closed {
		return
	}

	for s := range c.sessions {
		s.Set(c.name, false)
	}
	close(c.enterC)
	close(c.leaveC)
	close(c.closeC)
	close(c.broadcastC)

	c.closed = true

	slog.Info("Channel closed", "channel", c.name)
}

func (c *channel) run() {
	slog.Info("Channel running", "channel", c.name)

	for {
		select {
		case s := <-c.enterC:
			c.handleEnter(s)
		case s := <-c.leaveC:
			c.handleLeave(s)
		case m := <-c.broadcastC:
			c.handleBroadcast(m)
		case <-c.closeC:
			c.handleClose()
			return
		}
	}
}
