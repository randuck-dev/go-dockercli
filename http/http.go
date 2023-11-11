package http

import (
	"bufio"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/textproto"
	"slices"
	"strconv"
)

const (
	HTTP11 = "HTTP/1.1"
)

var ErrUnsupportedHTTPVersion = errors.New("unsupported HTTP version")
var ErrIncompleteStatusLine = errors.New("incomplete StatusLine. Needs 3 parts")
var ErrStatusCodeOutsideOfRange = errors.New("statuscode is outside of allowed range 100-599")
var ErrConnectionIsNil = errors.New("expected connection to be set, was nil")
var ErrLocationNotFoundOnRedirect = errors.New("expected Location header to be sent on redirect")

var ErrInvalidHeaderFormat = errors.New("invalid header format detected. Expected format: \"key: value\"")

var ErrImplementationDoesNotSupportMethod = errors.New("the implementation does not support the method")

type Client interface {
	Get(string) (Response, error)
}

type HttpClient struct {
	net.Conn
	address string
}

func NewHttpClient(address string) (HttpClient, error) {
	dial, err := net.Dial("tcp", address)
	if err != nil {
		return HttpClient{}, err
	}

	return HttpClient{dial, address}, nil
}

var supported_methods = []string{"GET", "HEAD"}

func (hc *HttpClient) Do(request Request) (Response, error) {
	if !slices.Contains(supported_methods, request.Method) {
		return Response{}, ErrImplementationDoesNotSupportMethod
	}

	slog.Debug("Do: Performing request to URI", "uri", request.Uri, "address", hc.RemoteAddr())
	written, err := hc.Write([]byte(request.ToRaw()))
	if err != nil {
		slog.Error("Error while writing to connection", "err", err)
		return Response{}, err
	}

	slog.Debug("Written bytes", "written", written)

	resp, err := parseResponse(hc)
	if err != nil {
		return Response{}, err
	}

	if resp.redirected() {
		loc, err := resp.location()
		if err != nil {
			return Response{}, ErrLocationNotFoundOnRedirect
		}

		req := request
		req.setLocation(loc)
		return hc.Do(req)
	}

	return resp, nil
}

func (hc *HttpClient) Get(uri string) (Response, error) {
	request := Request{
		Version: HTTP11,
		Method:  "GET",
		Uri:     uri,
		Host:    hc.address,
	}

	return hc.Do(request)
}

func (hc *HttpClient) Head(uri string) (Response, error) {
	// TODO: Is not allowed to have a body
	request := Request{
		Version: HTTP11,
		Method:  "HEAD",
		Uri:     uri,
		Host:    hc.address,
	}

	return hc.Do(request)
}

func Raw_http_parsing_docker_socket(docker_socket string) (Response, error) {
	client, err := NewHttpClient(docker_socket)
	if err != nil {
		return Response{}, err
	}

	defer client.Close()
	return client.Get("http://localhost/containers/json")
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
			slog.Debug("Finished parsing headers", "headers", headers)
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

	resp := Response{
		StatusLine: status_line,
		Headers:    headers,
	}

	content_length, err := resp.ContentLength()

	if err == ErrHeaderNotFound {
		return resp, nil
	}

	if err != nil {
		return Response{}, err
	}

	limit_reader := io.LimitReader(reader, content_length)

	buf := make([]byte, content_length)

	_, err = limit_reader.Read(buf)

	if err == io.EOF {
		return resp, nil
	}

	if err != nil {
		return Response{}, nil
	}

	resp.Body = buf
	return resp, nil
	line, err = tp.ReadLine()

	if err == io.EOF {
		return resp, nil
	}

	if err != nil && err != io.EOF {
		return Response{}, nil
	}

	if _, err := resp.TransferEncoding(); err == nil {
		val, err := strconv.ParseInt(line, 16, 32)
		if err != nil {
			return Response{}, err
		}

		limitReader := io.LimitReader(reader, val)
		buf := make([]byte, val)
		_, err = limitReader.Read(buf)

		if err != nil {
			return Response{}, err
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
