package main

import (
	"log/slog"
	"net/http"
)

func main() {
	res, err := http.Get("http://example.com")
	if err != nil {
		slog.Error("error when fetching", "err", err)
		return
	}

	slog.Info("got response", "res", res)
}
