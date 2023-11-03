package http

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/textproto"
	"sync"
)

const (
	HTTP11 = "HTTP/1.1"
)

var UnsupportedHttpVersion = errors.New("Unsupported HTTP Version")
var IncompleteStatusLine = errors.New("Incomplete StatusLine. Needs 3 parts")
var StatusCodeOutsideOfRange = errors.New("Statuscode is outside of allowed range 100-599")

var InvalidHeaderFormat = errors.New("Invalid Header Format detected. Expected Format: \"key: value\"")

type Client interface {
	Get(string) (Response, error)
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

type HttpClient struct {
	net.Conn
}

func (hc *HttpClient) Get(uri string) (Response, error) {
	written, err := hc.Write([]byte(fmt.Sprintf("GET %s HTTP/1.1\nHost: localhost \r\n\r\n", uri)))
	if err != nil {
		slog.Error("Error while writing to connection", "err", err)
		return Response{}, err
	}

	slog.Debug("Written bytes", "written", written)
	return Response{}, nil
}

func Raw_http_parsing_docker_socket(docker_socket string, wg *sync.WaitGroup) {

	socket, err := dial(docker_socket)

	if err != nil {
		slog.Error("Unable to connect to socket", "err", err)
		return
	}

	client := HttpClient{socket}

	defer client.Close()

	wg.Wait()
}

func listen(conn io.Reader, wg *sync.WaitGroup) {
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

func dial(addr string) (net.Conn, error) {

	conn, err := net.Dial("unix", addr)

	if err != nil {
		return nil, err
	}

	return conn, nil
}
