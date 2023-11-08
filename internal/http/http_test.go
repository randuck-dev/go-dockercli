package http

import (
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func XyzHandler(w http.ResponseWriter, r *http.Request) {
	resp := "Hello World"
	slog.Info("Received request with protocol type", "proto", r.Proto)
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(HttpStatusCodeOK)

	w.Write(([]byte(resp)))
}

func MiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("got a request", "request", r)
		next.ServeHTTP(w, r)
	})
}

func TestGet(t *testing.T) {

	mux := http.NewServeMux()
	mux.Handle("/", MiddleWare(http.HandlerFunc(XyzHandler)))

	server := httptest.NewServer(mux)
	defer server.Close()

	url := strings.Split(server.URL, "://")[1]
	conn, err := net.Dial("tcp", url)
	if err != nil {
		t.Errorf("failed to dial socket %s", err)
	}
	http_client := HttpClient{conn}
	res, err := http_client.Get("/")

	if err != nil {
		t.Errorf("Unexpected error %s", err)
	}

	if ct, err := res.ContentType(); err == nil && ct != "text/plain" {
		t.Errorf("got %s want %s", ct, "text/plain")
	}

	if !res.Ok() {
		t.Errorf("got %d want %d", res.StatusLine.StatusCode, HttpStatusCodeOK)
	}

	body := string(res.Body)

	if body != "Hello World" {
		t.Errorf("got %s want %s", body, "Hello World")
	}
}

func TestHead(t *testing.T) {
	return
	server, client := net.Pipe()
	http_client := HttpClient{client}

	go func() {
		defer client.Close()
		_, err := http_client.Head("xyz")

		if err != nil {
			t.Errorf("Unexpected error %s", err)
		}
	}()

	res, err := io.ReadAll(server)

	if err != nil {
		t.Errorf("Unexpected error %s", err)
	}

	expected := "HEAD xyz HTTP/1.1\nHost: localhost\r\n\r\n"

	if string(res) != expected {
		t.Errorf("got %s want %s", string(res), expected)
	}
}

func TestDo(t *testing.T) {
	return
	request := Request{
		Method: "INVALID",
	}

	c := HttpClient{}

	_, err := c.Do(request)

	if err != ErrImplementationDoesNotSupportMethod {
		t.Errorf("got %s want %s", err, ErrImplementationDoesNotSupportMethod)
	}
}

func TestParseResponse(t *testing.T) {
	t.Run("successfull response", func(t *testing.T) {
		responseRaw := "HTTP/1.1 200 OK\nHost: localhost\r\n\r\n"

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
		responseRaw := "HTTP/2.0 200 OK\nHost: localhost\r\n\r\n"

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
