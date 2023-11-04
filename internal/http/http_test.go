package http

import (
	"io"
	"net"
	"strings"
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

func TestParseResponse(t *testing.T) {
	t.Run("successfull response", func(t *testing.T) {
		responseRaw := "HTTP/1.1 200 OK\nHost: localhost \r\n\r\n"

		resp, err := parseResponse(strings.NewReader(responseRaw))
		if err != nil {
			t.Errorf("Unexecpted error %s", err)
		}

		if resp.StatusLine.HttpVersion != HTTP11 {
			t.Errorf("got %s want %s", resp.StatusLine.HttpVersion, HTTP11)
		}

		if resp.StatusLine.StatusCode != 200 {
			t.Errorf("got %d want %d", resp.StatusLine.StatusCode, 200)
		}

		if resp.StatusLine.ReasonPhrase != "OK" {
			t.Errorf("got %s want %s", resp.StatusLine.ReasonPhrase, "OK")
		}

		val, ok := resp.Headers["Host"]
		if !ok {
			t.Errorf("unable to find Host header in headers")
		}

		if val != "localhost" {
			t.Errorf("got %s want %s", val, "localhost")
		}
	})

	t.Run("Failed response: UnsupportedHttpVersion", func(t *testing.T) {
		responseRaw := "HTTP/2.0 200 OK\nHost: localhost \r\n\r\n"

		_, err := parseResponse(strings.NewReader(responseRaw))
		if err != ErrUnsupportedHTTPVersion {
			t.Errorf("Unexecpted error %s", err)
		}
	})

	t.Run("Failed Response: ConnectionIsNil", func(t *testing.T) {
		_, err := parseResponse(nil)

		if err != ErrConnectionIsNil {
			t.Errorf("got %s want %s", err, ErrConnectionIsNil)
		}
	})
}
