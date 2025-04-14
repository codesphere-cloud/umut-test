package main

import (
	"fmt"
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", sseHandler)

	e.Logger.Fatal(e.Start(":3000"))
}

func sseHandler(c echo.Context) error {
	log.Printf("SSE client connected, ip: %v", c.RealIP())

	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-c.Request().Context().Done():
			log.Printf("SSE client disconnected, ip: %v", c.RealIP())
			return nil
		case <-ticker.C:
			fmt.Fprint(w, ": heartbeat\n\n")
			w.Flush()
		}
	}
}
