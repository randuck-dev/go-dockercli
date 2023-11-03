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

func TestHeaderParsing(t *testing.T) {
	t.Run("Header parsing is successful", func(t *testing.T) {
		rawHeader := "Content-Type: application/json"
		key, value, err := parseHeader(rawHeader)

		if err != nil {
			t.Errorf("Did not expect to fail got %s", err)
		}

		if key != "Content-Type" {
			t.Errorf("got %s want %s", key, "Content-Type")
		}

		if value != "application/json" {
			t.Errorf("got %s want %s", value, "application/json")
		}
	})

	t.Run("Header fails because of missing key", func(t *testing.T) {
		rawHeader := ": application/json"
		_, _, err := parseHeader(rawHeader)

		if err != InvalidHeaderFormat {
			t.Errorf("got %s want %s", err, InvalidHeaderFormat)
		}
	})

	t.Run("Header fails because of missing value", func(t *testing.T) {
		rawHeader := "Content-Type: "
		_, _, err := parseHeader(rawHeader)

		if err != InvalidHeaderFormat {
			t.Errorf("got %s want %s", err, InvalidHeaderFormat)
		}
	})

	t.Run("Header fails because it is not possible to split into key and value", func(t *testing.T) {
		rawHeader := "Content-Type application/json"
		_, _, err := parseHeader(rawHeader)

		if err != InvalidHeaderFormat {
			t.Errorf("got %s want %s", err, InvalidHeaderFormat)
		}
	})
}
