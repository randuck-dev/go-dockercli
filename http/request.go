package http

import (
	"fmt"
	"strings"
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
	if strings.HasPrefix(r.Host, "/") {
		host = "Host: localhost"
	}
	return fmt.Sprintf("%s %s %s\n%s%s", r.Method, r.Uri, r.Version, host, EndOfMessage)
}

func (r *Request) setLocation(location string) {
	r.Uri = location
}
