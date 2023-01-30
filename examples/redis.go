package main

import (
	"context"
	"encoding/json"
	"golang.org/x/exp/slog"

	"github.com/redis/go-redis/v9"

	"github.com/yl/echo"
)

func Subscribe(e *echo.Echo) {
	c := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx := context.Background()
	sub := c.Subscribe(ctx, "ws")
	defer func() {
		_ = sub.Close()
	}()

	for {
		select {
		case m := <-sub.Channel():
			message := &echo.Message{}
			if err := json.Unmarshal([]byte(m.Payload), message); err != nil {
				slog.Error("Message unmarshal error", err)
				continue
			}
			msg, err := json.Marshal(message)
			if err != nil {
				slog.Error("Message marshal error", err)
				continue
			}
			e.Broadcast(message.Channel, msg)
		case <-ctx.Done():
			return
		}
	}
}