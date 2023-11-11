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

https://github.com/randuck-dev/http-parser/blob/4d9d4dbcfc9473eee11ea78edcbce3e33068fc18/cmd/externals/externals.go#L1