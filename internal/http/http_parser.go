package http

import (
	"log/slog"
	"strconv"
	"strings"
)

func parseStatusLine(payload string) (StatusLine, error) {

	split_line := strings.Split(payload, " ")
	if len(split_line) != 3 {
		return StatusLine{}, ErrIncompleteStatusLine
	}

	httpVersion := split_line[0]

	if httpVersion != HTTP11 {
		return StatusLine{}, ErrUnsupportedHTTPVersion
	}

	status_code, err := strconv.ParseUint(split_line[1], 10, 16)
	if err != nil {
		slog.Error("Error occured while parsing status code status line", "line", payload)
		return StatusLine{}, err
	}

	if status_code < 100 || status_code > 599 {
		return StatusLine{}, ErrStatusCodeOutsideOfRange
	}

	sl := StatusLine{
		split_line[0],
		uint16(status_code),
		split_line[2],
	}

	slog.Info("parsed status line", "status_line", sl)
	return sl, nil
}

func parseHeader(rawHeader string) (string, string, error) {
	headers_split := strings.SplitN(rawHeader, ":", 2)

	if len(headers_split) < 2 {
		return "", "", ErrInvalidHeaderFormat
	}

	key := strings.TrimSpace(headers_split[0])
	value := strings.TrimSpace(headers_split[1])

	if len(key) == 0 {
		return "", "", ErrInvalidHeaderFormat
	}

	if len(value) == 0 {
		return "", "", ErrInvalidHeaderFormat
	}
	return key, value, nil
}
