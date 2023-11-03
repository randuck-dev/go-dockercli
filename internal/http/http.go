package http

import (
	"bufio"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/textproto"
	"strconv"
	"strings"
	"sync"
)

const (
	HTTP11 = "HTTP/1.1"
)

var UnsupportedHttpVersion = errors.New("Unsupported HTTP Version")
var IncompleteStatusLine = errors.New("Incomplete StatusLine. Needs 3 parts")
var StatusCodeOutsideOfRange = errors.New("Statuscode is outside of allowed range 100-599")

var InvalidHeaderFormat = errors.New("Invalid Header Format detected. Expected Format: \"key: value\"")

func Raw_http_parsing_docker_socket(docker_socket string, wg *sync.WaitGroup) {

	socket, err := dial(docker_socket)

	if err != nil {
		slog.Error("Unable to connect to socket", "err", err)
		return
	}

	defer socket.Close()
	go listen(socket, wg)
	written, err := socket.Write([]byte("GET /containers/json HTTP/1.1\nHost: localhost \r\n\r\n"))

	if err != nil {
		slog.Error("Error while writing to connection", "err", err)
		return
	}

	slog.Debug("Written bytes", "written", written)

	wg.Wait()
}

func parseStatusLine(payload string) (StatusLine, error) {

	split_line := strings.Split(payload, " ")
	if len(split_line) != 3 {
		return StatusLine{}, IncompleteStatusLine
	}

	httpVersion := split_line[0]

	if httpVersion != HTTP11 {
		return StatusLine{}, UnsupportedHttpVersion
	}

	status_code, err := strconv.ParseUint(split_line[1], 10, 16)
	if err != nil {
		slog.Error("Error occured while parsing status code status line", "line", payload)
		return StatusLine{}, err
	}

	if status_code < 100 || status_code > 599 {
		return StatusLine{}, StatusCodeOutsideOfRange
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
		return "", "", InvalidHeaderFormat
	}

	key := strings.TrimSpace(headers_split[0])
	value := strings.TrimSpace(headers_split[1])

	if len(key) == 0 {
		return "", "", InvalidHeaderFormat
	}

	if len(value) == 0 {
		return "", "", InvalidHeaderFormat
	}
	return key, value, nil
}

func listen(conn net.Conn, wg *sync.WaitGroup) {
	reader := bufio.NewReader(conn)
	tp := textproto.NewReader(reader)

	current_line := 0
	parsing_headers := true
	headers := make(map[string]string)
	for {
		line, err := tp.ReadLine()
		if err == io.EOF {
			slog.Info("End of file reached")
			break
		}
		if err != nil {
			slog.Error("Error occurred while reading line", "err", err)
			return
		}
		if current_line == 0 {
			sl, err := parseStatusLine(line)
			if err != nil {
				slog.Error("Error when parsing status line", "err", err, "line", line)
				return
			}
			slog.Info("Parsed status line", "statusline", sl)
			current_line += 1
			continue
		}

		// Now we are at a point where there might be headers
		if line == "" && parsing_headers {
			parsing_headers = false
			slog.Info("Finished parsing headers", "headers", headers)
		} else if line != "" && parsing_headers {
			slog.Debug("Header", "header", line)

			key, value, err := parseHeader(line)
			if err != nil {
				slog.Error("Error when parsing header", "err", err, "line", line)
			}
			headers[key] = value
		}

		slog.Info("Read a line from connection", "line", line)
		current_line += 1
	}

	wg.Done()
}

type Response struct {
	StatusLine StatusLine
	Headers    map[string]string
}

type StatusLine struct {
	HttpVersion  string
	StatusCode   uint16
	ReasonPhrase string
}

func dial(addr string) (net.Conn, error) {

	conn, err := net.Dial("unix", addr)

	if err != nil {
		return nil, err
	}

	return conn, nil
}
