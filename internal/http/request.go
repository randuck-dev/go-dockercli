package http

import "fmt"

type Request struct {
	Method  string
	Uri     string
	Version string

	requestHeader string
}

func (r Request) ToRaw() string {
	host := "Host: localhost"
	return fmt.Sprintf("%s %s %s\n%s%s", r.Method, r.Uri, r.Version, host, EndOfMessage)
}
