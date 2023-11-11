package http

import (
	"fmt"
)

type Request struct {
	Method  string
	Uri     string
	Version string

	Host          string
	requestHeader string
}

func (r Request) ToRaw() string {
	host := fmt.Sprintf("Host: %s", r.Host)
	return fmt.Sprintf("%s %s %s\n%s%s", r.Method, r.Uri, r.Version, host, EndOfMessage)
}

func (r *Request) setLocation(location string) {
	r.Uri = location
}
