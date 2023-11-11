package http

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func XyzHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		resp := "Hello World"
		slog.Info("Received request with protocol type", "proto", r.Proto)
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(HttpStatusCodeOK)

		w.Write(([]byte(resp)))
	}

	if r.Method == "HEAD" {
		w.WriteHeader(HttpStatusCodeOK)
	}

	w.WriteHeader(HttpStatusCodeNotFound)
}

func PermanentlyMovedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Add("Location", "/redirecttarget")
		w.WriteHeader(HttpStatusCodeMovedPermanently)
	}

	w.WriteHeader(HttpStatusCodeNotFound)
}

func TemporarilyMovedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Add("Location", "/redirecttarget")
		w.WriteHeader(HttpStatusCodeTemporaryRedirect)
	}

	w.WriteHeader(HttpStatusCodeNotFound)
}

func TargetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		resp := "RedirectTargetHandler"
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(HttpStatusCodeOK)

		w.Write(([]byte(resp)))
	}
}

func MiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("got a request", "request", r)
		next.ServeHTTP(w, r)
	})
}

func BuildServer(t *testing.T) (*httptest.Server, string) {
	t.Helper()
	mux := http.NewServeMux()
	mux.Handle("/", MiddleWare(http.HandlerFunc(XyzHandler)))
	mux.Handle("/movedpermanently", MiddleWare(http.HandlerFunc(PermanentlyMovedHandler)))
	mux.Handle("/movedtemporarily", MiddleWare(http.HandlerFunc(TemporarilyMovedHandler)))
	mux.Handle("/redirecttarget", MiddleWare(http.HandlerFunc(TargetHandler)))

	server := httptest.NewServer(mux)
	url := strings.Split(server.URL, "://")[1]
	return server, url
}

func TestGet(t *testing.T) {
	server, url := BuildServer(t)
	defer server.Close()

	http_client, err := NewHttpClient(TcpDialContext(url))

	if err != nil {
		t.Errorf("failed to create new http client %s", err)
	}
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
	server, url := BuildServer(t)
	http_client, err := NewHttpClient(TcpDialContext(url))

	if err != nil {
		t.Errorf("unexpected error when creating http client %s", err)
	}
	defer server.Close()
	resp, err := http_client.Head("/")

	if err != nil {
		t.Errorf("Unexpected error %s", err)
	}

	if !resp.Ok() {
		t.Errorf("got %d want %d", resp.StatusLine.StatusCode, 200)
	}
}

func TestDo(t *testing.T) {
	request := Request{
		Method: "INVALID",
	}

	c := HttpClient{}

	_, err := c.Do(request)

	if err != ErrImplementationDoesNotSupportMethod {
		t.Errorf("got %s want %s", err, ErrImplementationDoesNotSupportMethod)
	}
}

func TestRedirect(t *testing.T) {
	t.Run("moved permanently", func(t *testing.T) {
		server, url := BuildServer(t)
		http_client, err := NewHttpClient(TcpDialContext(url))

		if err != nil {
			t.Errorf("unexpected error when creating http client %s", err)
		}
		defer server.Close()
		resp, err := http_client.Get("/movedpermanently")

		if err != nil {
			t.Errorf("Unexpected error %s", err)
		}

		if !resp.Ok() {
			t.Errorf("got %d want %d", resp.StatusLine.StatusCode, HttpStatusCodeOK)
		}

		body := string(resp.Body)

		if body != "RedirectTargetHandler" {
			t.Errorf("got %s want %s", body, "RedirectTargetHandler")
		}
	})

	t.Run("moved temporarily", func(t *testing.T) {
		server, url := BuildServer(t)
		http_client, err := NewHttpClient(TcpDialContext(url))

		if err != nil {
			t.Errorf("unexpected error when creating http client %s", err)
		}
		defer server.Close()
		resp, err := http_client.Get("/movedtemporarily")

		if err != nil {
			t.Errorf("Unexpected error %s", err)
		}

		if !resp.Ok() {
			t.Errorf("got %d want %d", resp.StatusLine.StatusCode, HttpStatusCodeOK)
		}

		body := string(resp.Body)

		if body != "RedirectTargetHandler" {
			t.Errorf("got %s want %s", body, "RedirectTargetHandler")
		}
	})
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

		if resp.StatusLine.StatusCode != HttpStatusCodeOK {
			t.Errorf("got %d want %d", resp.StatusLine.StatusCode, HttpStatusCodeOK)
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
