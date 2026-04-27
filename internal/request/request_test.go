package request

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestFromReader(t *testing.T) {
	// Test: Good GET Request line
	r, err := RequestFromReader(strings.NewReader("GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	assert.Equal(t, "HTTP/1.1", r.RequestLine.HttpVersion)

	// Test: Good GET Request line with path
	r, err = RequestFromReader(strings.NewReader("GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "HTTP/1.1", r.RequestLine.HttpVersion)

	// Test: Invalid number of parts in request line
	_, err = RequestFromReader(strings.NewReader("/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.Error(t, err)
}

func TestRequestLineParse(t *testing.T) {
	// Test: Good GET Request line
	l, err := parseRequestLine("GET / HTTP/1.1")
	require.NoError(t, err)
	assert.Equal(t, "GET", l.Method)
	assert.Equal(t, "/", l.RequestTarget)
	assert.Equal(t, "HTTP/1.1", l.HttpVersion)

	// Test: Good GET Request line with path
	l, err = parseRequestLine("GET /coffee HTTP/1.1")
	require.NoError(t, err)
	assert.Equal(t, "GET", l.Method)
	assert.Equal(t, "/coffee", l.RequestTarget)
	assert.Equal(t, "HTTP/1.1", l.HttpVersion)

	// Test: Good POST Request with path
	l, err = parseRequestLine("POST /milk HTTP/1.1")
	require.NoError(t, err)
	assert.Equal(t, "POST", l.Method)
	assert.Equal(t, "/milk", l.RequestTarget)
	assert.Equal(t, "HTTP/1.1", l.HttpVersion)

	// Test: Path with query parameters
	l, err = parseRequestLine("GET /search?q=coffee&size=large HTTP/1.1")
	require.NoError(t, err)
	assert.Equal(t, "/search?q=coffee&size=large", l.RequestTarget)

	// Test: Custom valid method (token)
	l, err = parseRequestLine("UPDATE /profile HTTP/1.1")
	require.NoError(t, err)
	assert.Equal(t, "UPDATE", l.Method)

	// Test: Lowercase method (invalid)
	l, err = parseRequestLine("get /coffee HTTP/1.1")
	require.Error(t, err)

	// Test: Invalid version in Request line
	l, err = parseRequestLine("GET /coffee HTTP/2.1")
	require.Error(t, err)

	// Test: Invalid number of parts (too many)
	l, err = parseRequestLine("GET /coffee HTTP/1.1 END")
	require.Error(t, err)

	// Test: Invalid number of parts (missing version)
	_, err = parseRequestLine("GET /only-two-parts")
	require.Error(t, err)

	// Test: Multiple spaces between parts (RFC forbids this)
	l, err = parseRequestLine("GET  /coffee  HTTP/1.1")
	require.Error(t, err)

	// Test: Empty input
	_, err = parseRequestLine("")
	require.Error(t, err)

	// Test: Only whitespace
	_, err = parseRequestLine("   ")
	require.Error(t, err)

	// Test: Leading whitespace (invalid start-line)
	_, err = parseRequestLine(" GET / HTTP/1.1")
	require.Error(t, err)
}
