package http

import (
	"io"
	"net"
	"testing"
)

func TestGet(t *testing.T) {

	server, client := net.Pipe()
	http_client := HttpClient{client}

	go func() {
		defer client.Close()
		_, err := http_client.Get("xyz")

		if err != nil {
			t.Errorf("Unexpected error %s", err)
		}
	}()

	res, err := io.ReadAll(server)

	if err != nil {
		t.Errorf("Unexpected error %s", err)
	}

	expected := "GET xyz HTTP/1.1\nHost: localhost \r\n\r\n"

	if string(res) != expected {
		t.Errorf("got %s want %s", string(res), expected)
	}
}
