package main

import (
	"github.com/yl/echo"
)

func main() {
	e := echo.New()

	go Subscribe(e)
	e.HandleHttp(Broadcast)

	e.Run(":8000")
}
