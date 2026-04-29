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
	if !isValidFieldName(string(fieldName)) {
		return 0, false, fmt.Errorf("invalid header name: %s", fieldName)
	}

	fieldValue := bytes.TrimSpace(data[idx+1 : lineEnd])

	h.Set(string(fieldName), string(fieldValue))

	return lineEnd + 2, false, nil
}

func (h Headers) Set(key string, value string) {
	key = strings.ToLower(key)
	v, ok := h[key]
	if ok {
		value = v + ", " + value
	}
	h[key] = value
}

func isValidFieldName(s string) bool {
	for i := 0; i < len(s); i++ {
		if !isValidChar(s[i]) {
			return false
		}
	}
	return true
}

func isValidChar(c byte) bool {
	return (c >= 'A' && c <= 'Z') ||
		(c >= 'a' && c <= 'z') ||
		(c >= '0' && c <= '9') ||
		(c == '!' || c == '#' || c == '$' || c == '%' || c == '&' ||
			c == '\'' || c == '*' || c == '+' || c == '-' || c == '.' ||
			c == '^' || c == '_' || c == '`' || c == '|' || c == '~')
}
