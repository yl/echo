package echo

import (
	"encoding/json"
	"golang.org/x/exp/slog"

	"github.com/olahol/melody"
)

type Message struct {
	Channel string `json:"channel"`
	Payload any    `json:"payload"`
}

type request struct {
	ID     any      `json:"id"`
	Method string   `json:"method"`
	Params []string `json:"params"`
}

type response struct {
	ID     any    `json:"id"`
	Error  string `json:"error,omitempty"`
	Result any    `json:"result,omitempty"`
}

func HandleMessage(e *Echo) func(*melody.Session, []byte) {
	return func(s *melody.Session, m []byte) {
		request := &request{}
		if err := json.Unmarshal(m, request); err != nil {
			slog.Error("Message unmarshal error", err, "msg", m)
		}
		response := &response{ID: request.ID}

		switch request.Method {
		case "subscribe":
			for _, n := range request.Params {
				e.Enter(n, s)
			}
			response.Result = true
		case "unsubscribe":
			for _, n := range request.Params {
				e.Leave(n, s)
			}
			response.Result = true
		default:
			response.Error = "Method Not Found"
		}

		msg, err := json.Marshal(response)
		if err != nil {
			slog.Error("Message marshal error", err, "msg", response)
		}
		if err := s.Write(msg); err != nil {
			slog.Error("Message send error", err)
		}
	}
}

func HandleDisconnect(e *Echo) func(session *melody.Session) {
	return func(s *melody.Session) {
		for n, v := range s.Keys {
			if exist, ok := v.(bool); ok && exist {
				e.Leave(n, s)
			}
		}
	}
}
