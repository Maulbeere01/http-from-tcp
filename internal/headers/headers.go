package headers

import (
	"bytes"
	"fmt"
	"strings"
)

const crlf = "\r\n"

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	lineEnd := bytes.Index(data, []byte(crlf))
	if lineEnd == -1 {
		return 0, false, nil
	}
	if lineEnd == 0 {
		return 2, true, nil
	}

	idx := bytes.Index(data[:lineEnd], []byte(`:`))

	if idx == -1 {
		return 0, false, fmt.Errorf("invalid header format: %s", string(data))
	}

	fieldName := data[:idx]
	for _, c := range fieldName {
		if c == ' ' {
			return 0, false, fmt.Errorf("invalid header format: %s", fieldName)
		}
	}

	fieldValue := bytes.TrimSpace(data[idx+1 : lineEnd])
	h[strings.ToLower(string(fieldName))] = string(fieldValue)

	return lineEnd + 2, false, nil
}
