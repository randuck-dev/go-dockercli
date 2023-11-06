package main

import (
	"http-parser/internal/http"
	"log/slog"
	"net"
)

func main() {

	conn, err := net.Dial("tcp", "example.com:80")
	if err != nil {
		slog.Error("Unable to dial", "err", err)
		return
	}

	c := http.HttpClient{
		Conn: conn,
	}

	resp, err := c.Get("example.com")
	if err != nil {
		slog.Error("Error getting data from uri", "err", err)
		return
	}

	slog.Info("received data", "data", resp)
}
