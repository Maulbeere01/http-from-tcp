package request

import (
	"fmt"
	"io"
	"regexp"
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

	reqLine, err := parseRequestLine(string(data))
	if err != nil {
		return nil, err
	}
	return &Request{RequestLine: *reqLine}, nil
}

func parseRequestLine(input string) (*RequestLine, error) {
	msg := strings.Split(input, "\r\n")
	parts := strings.Split(msg[0], " ")

	if len(parts) != 3 {
		return &RequestLine{}, fmt.Errorf("invalid number of fields in request line: %d", len(parts))
	}

	var reqLine RequestLine
	reqLine.Method = parts[0]
	reqLine.RequestTarget = parts[1]
	reqLine.HttpVersion = parts[2]

	var validMethod = regexp.MustCompile(`^[A-Z]+$`)
	if !validMethod.MatchString(reqLine.Method) || reqLine.HttpVersion != "HTTP/1.1" {
		return &RequestLine{}, fmt.Errorf("invalid request line: %s", input)
	}
	return &reqLine, nil
}
