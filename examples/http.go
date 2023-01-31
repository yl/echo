package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/yl/echo"
	"golang.org/x/exp/slog"
)

func Broadcast(e *echo.Echo) (string, func(http.ResponseWriter, *http.Request)) {
	return "/broadcast", func(w http.ResponseWriter, r *http.Request) {
		msg, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("io read error", err)
			return
		}
		message := &Message{}
		err = json.Unmarshal(msg, message)
		if err != nil {
			slog.Error("json unmarshal error", err, "msg", msg)
		}
		e.Broadcast(message.Channel, msg)
	}
}
