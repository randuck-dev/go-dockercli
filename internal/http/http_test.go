package http

import (
	"testing"
)

func TestParseStatusLine(t *testing.T) {
	t.Run("Correct status line 200 OK", func(t *testing.T) {
		rawStatusLine := "HTTP/1.1 200 OK"
		resp, err := parseStatusLine(rawStatusLine)

		if err != nil {
			t.Error(err)
		}

		if resp.HttpVersion != "HTTP/1.1" {
			t.Errorf("got %s want %s", resp.HttpVersion, "HTTP/1.1")
		}

		if resp.StatusCode != 200 {
			t.Errorf("got %d want %d", resp.StatusCode, 200)
		}

		if resp.ReasonPhrase != "OK" {
			t.Errorf("got %s want %s", resp.ReasonPhrase, "OK")
		}
	})

	t.Run("Status line fails when statuscode is not an integer", func(t *testing.T) {
		rawStatusLine := "HTTP/1.1 FAIL OK"
		_, err := parseStatusLine(rawStatusLine)
		if err == nil {
			t.Errorf("Expected test to fail")
		}
	})

	t.Run("Status line fails when the status line is incomplete", func(t *testing.T) {
		rawStatusLine := "200 OK"
		_, err := parseStatusLine(rawStatusLine)
		if err != IncompleteStatusLine {
			t.Errorf("got %s want %s", err, IncompleteStatusLine)
		}
	})

	t.Run("Status line fails when the HTTP version is not recognized", func(t *testing.T) {
		rawStatusLine := "HTTP/-1 200 OK"
		_, err := parseStatusLine(rawStatusLine)
		if err != UnsupportedHttpVersion {
			t.Errorf("got %s want %s", err, UnsupportedHttpVersion)
		}
	})

	t.Run("Statuscode must fail if it is outside of allowed range", func(t *testing.T) {
		rawStatusLine := "HTTP/1.1 09 OK"
		_, err := parseStatusLine(rawStatusLine)
		if err != StatusCodeOutsideOfRange {
			t.Errorf("got %s want %s", err, StatusCodeOutsideOfRange)
		}
	})
}
