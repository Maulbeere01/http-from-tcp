package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

const crlf = "\r\n"
const bufferSize = 8

type parserState int

const (
	stateInitialized parserState = iota
	stateDone
)

type Request struct {
	RequestLine RequestLine
	state       parserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0

	r := &Request{
		state: stateInitialized,
	}

	for r.state != stateDone {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		n, err := reader.Read(buf[readToIndex:])
		if n > 0 {
			readToIndex += n
		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				r.state = stateDone
				break
			}
			return nil, err
		}

		consumed, err := r.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}
		if consumed > 0 {
			copy(buf, buf[consumed:readToIndex])
			readToIndex -= consumed
		}
	}
	return r, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case stateInitialized:
		reqLine, bytesRead, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if bytesRead == 0 {
			return 0, nil // need more data
		}
		r.RequestLine = *reqLine
		r.state = stateDone
		return bytesRead, nil
	case stateDone:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("error: unknown state")
	}
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, 0, nil
	}

	reqLine, err := requestLineFromString(string(data[:idx]))
	if err != nil {
		return nil, 0, err
	}
	return reqLine, idx + 2, nil // +2 to include CRLF
}

func requestLineFromString(str string) (*RequestLine, error) {
	parts := strings.Split(str, ` `)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid request line format: %s", str)
	}

	method := parts[0]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}

	requestTarget := parts[1]

	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("malformed start-line: %s", str)
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpPart)
	}
	version := versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", version)
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   versionParts[1],
	}, nil
}
