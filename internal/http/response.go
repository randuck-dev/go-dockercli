package http

import "errors"

type Response struct {
	StatusLine StatusLine
	Headers    map[string]string
}

type StatusLine struct {
	HttpVersion  string
	StatusCode   uint16
	ReasonPhrase string
}

var ErrNoContentTypefound = errors.New("no content type")

func (r Response) ContentType() (string, error) {

	res, ok := r.Headers["Content-Type"]

	if !ok {
		return "", ErrNoContentTypefound
	}

	return res, nil
}
