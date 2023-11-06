package http

import (
	"bufio"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/textproto"
	"slices"
	"sync"
)

const (
	HTTP11 = "HTTP/1.1"
)

var ErrUnsupportedHTTPVersion = errors.New("unsupported HTTP version")
var ErrIncompleteStatusLine = errors.New("incomplete StatusLine. Needs 3 parts")
var ErrStatusCodeOutsideOfRange = errors.New("statuscode is outside of allowed range 100-599")
var ErrConnectionIsNil = errors.New("expected connection to be set, was nil")

var ErrInvalidHeaderFormat = errors.New("invalid header format detected. Expected format: \"key: value\"")

var ErrImplementationDoesNotSupportMethod = errors.New("the implementation does not support the method")

type Client interface {
	Get(string) (Response, error)
}

type HttpClient struct {
	net.Conn
}

var supported_methods = []string{"GET", "HEAD"}

func (hc *HttpClient) Do(request Request) (Response, error) {
	if !slices.Contains(supported_methods, request.Method) {
		return Response{}, ErrImplementationDoesNotSupportMethod
	}

	written, err := hc.Write([]byte(request.ToRaw()))
	if err != nil {
		slog.Error("Error while writing to connection", "err", err)
		return Response{}, err
	}

	slog.Debug("Written bytes", "written", written)
	return Response{}, nil
}

func (hc *HttpClient) Get(uri string) (Response, error) {
	request := Request{
		Version: HTTP11,
		Method:  "GET",
		Uri:     uri,
	}

	return hc.Do(request)
}

func (hc *HttpClient) Head(uri string) (Response, error) {
	// TODO: Is not allowed to have a body
	request := Request{
		Version: HTTP11,
		Method:  "HEAD",
		Uri:     uri,
	}

	return hc.Do(request)
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

func parseResponse(conn io.Reader) (Response, error) {
	if conn == nil {
		return Response{}, ErrConnectionIsNil
	}

	reader := bufio.NewReader(conn)
	tp := textproto.NewReader(reader)

	var status_line StatusLine
	headers := make(map[string]string)
	line, err := tp.ReadLine()
	if err != nil {
		slog.Error("Error occurred while reading line", "err", err)
		return Response{}, err
	}
	status_line, err = parseStatusLine(line)
	if err != nil {
		slog.Error("Error when parsing status line", "err", err, "line", line)
		return Response{}, err
	}

	parsing_headers := true
	for {
		line, err := tp.ReadLine()
		if err != nil {
			return Response{}, err
		}

		if line == "" && parsing_headers {
			parsing_headers = false
			slog.Info("Finished parsing headers", "headers", headers)
			break
		} else if line != "" && parsing_headers {
			slog.Debug("Header", "header", line)

			key, value, err := parseHeader(line)
			if err != nil {
				slog.Error("Error when parsing header", "err", err, "line", line)
			}
			headers[key] = value
		}
	}

	return Response{
		StatusLine: status_line,
		Headers:    headers,
	}, nil
}

func dial(addr string) (net.Conn, error) {

	conn, err := net.Dial("unix", addr)

	if err != nil {
		return nil, err
	}

	return conn, nil
}
