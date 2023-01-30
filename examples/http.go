package main

import (
	"encoding/json"
	"golang.org/x/exp/slog"
	"io"
	"net/http"

	"github.com/yl/echo"
)

func Broadcast(e *echo.Echo) (string, func(http.ResponseWriter, *http.Request)) {
	return "/broadcast", func(w http.ResponseWriter, r *http.Request) {
		msg, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("io read error", err)
			return
		}
		message := &echo.Message{}
		err = json.Unmarshal(msg, message)
		if err != nil {
			slog.Error("json unmarshal error", err, "msg", msg)
		}
		e.Broadcast(message.Channel, msg)
	}
}
