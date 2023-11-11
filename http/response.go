package http

import (
	"errors"
	"slices"
	"strconv"
)

type Response struct {
	StatusLine StatusLine
	Headers    map[string]string

	Body []byte
}

type StatusLine struct {
	HttpVersion  string
	StatusCode   uint16
	ReasonPhrase string
}

var ErrNoContentTypefound = errors.New("no content type")
var ErrHeaderNotFound = errors.New("header not found")
var ErrInvalidContentLengthFormat = errors.New("invalid content lenght format")

var redirectCodes = []int{
	HttpStatusCodeTemporaryRedirect,
	HttpStatusCodeMovedPermanently,
}

func (r Response) ContentType() (string, error) {

	res, ok := r.Headers["Content-Type"]

	if !ok {
		return "", ErrNoContentTypefound
	}

	return res, nil
}

func (r Response) TransferEncoding() (string, error) {
	res, ok := r.Headers["Transfer-Encoding"]

	if !ok {
		return "", ErrHeaderNotFound
	}

	return res, nil
}

func (r Response) Ok() bool {
	return r.StatusLine.StatusCode == HttpStatusCodeOK
}

func (r Response) ContentLength() (int64, error) {
	res, ok := r.Headers["Content-Length"]

	if !ok {
		return -1, ErrHeaderNotFound
	}

	val, err := strconv.ParseInt(res, 10, 64)

	if err != nil {
		return -1, ErrInvalidContentLengthFormat
	}

	return val, nil
}

func (r Response) location() (string, error) {
	res, ok := r.Headers["Location"]

	if !ok {
		return "", ErrHeaderNotFound
	}

	return res, nil
}

func (r Response) redirected() bool {
	return slices.Contains(redirectCodes, int(r.StatusLine.StatusCode))
}
