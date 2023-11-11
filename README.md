# Custom HTTP Parser

A custom HTTP Parser that attempts to be compliant with some of HTTP Semantics defined as of [RFC 9110](https://datatracker.ietf.org/doc/html/rfc9110). Furthermore the focus will be on HTTP/1.1 as per [RFC 9112](https://datatracker.ietf.org/doc/html/rfc9112.html)

This is built purely for educational purposes as to get a deeper understanding for how the HTTP protocol works underneath the hood.

## TODO
- Body parsing
  - Partially implemented. Support for text/plain
- Streaming
- Error handling
- Etc?


## Example of using the library


```go
c, err := http.NewHttpClient(http.TcpDialContext("example.com:80"))
if err != nil {
  slog.Error("Error while creating http client", "err", err)
  return
}

resp, err := c.Get("/")
if err != nil {
  slog.Error("Error calling GET on uri", "err", err)
  return
}

slog.Info("Received response", "body", string(resp.Body))
```