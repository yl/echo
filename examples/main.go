package main

import (
	"github.com/yl/echo"
)

func main() {
	e := echo.New()
	e.HandleMessage(HandleMessage(e))
	e.HandleDisconnect(HandleDisconnect(e))

	go Subscribe(e)
	e.HandleHttp(Broadcast)

	e.Run(":8000")
}
