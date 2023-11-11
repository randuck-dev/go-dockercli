package main

import (
	"log/slog"

	"github.com/randuck-dev/http-parser/http"
)

func main() {

	c, err := http.NewHttpClient("example.com:80")
	if err != nil {
		slog.Error("Error while creating client", "err", err)
		return
	}

	resp, err := c.Get("/")
	if err != nil {
		slog.Error("Error getting data from uri", "err", err)
		return
	}

	slog.Info("received data", "body", string(resp.Body))
}
