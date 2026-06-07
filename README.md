# HTTP from TCP

**⚠️ Work in progress.**

A HTTP/1.1 implementation in Go, built directly on raw TCP sockets. 

Built by following Boot.dev's "HTTP from TCP" course, implementing the message format from [RFC 9112 §2.1](https://datatracker.ietf.org/doc/html/rfc9112#name-message-format).

## Status

- [x] Request line parsing
- [x] Header parsing
- [x] Request body parsing
- [x] TCP listener
- [ ] HTTP responses
- [ ] Chunked encoding
- [ ] Binary data

## What it does

Parses raw bytes off a TCP connection into HTTP requests, following the HTTP/1.1 message grammar: a start line, header field lines, a blank line marking the end of the headers, then an optional body. Response writing is the next step.

## Build and run

```bash
go build ./...
go run ./cmd/tcplistener
```

The listener accepts a connection on port 42069 and prints the parsed request. Send one with `curl` or `nc`.
