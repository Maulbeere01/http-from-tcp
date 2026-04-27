package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const crlf = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	reqLine, err := parseRequestLine(data)
	if err != nil {
		return nil, err
	}
	return &Request{RequestLine: *reqLine}, nil
}

func parseRequestLine(data []byte) (*RequestLine, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, fmt.Errorf("no CRLF found in requesst line")
	}

	reqLine, err := requestLineFromString(string(data[:idx]))
	if err != nil {
		return nil, err
	}
	return reqLine, nil
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
