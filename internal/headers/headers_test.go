package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestHeaderParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid single header with extra whitespace
	headers = NewHeaders()
	data = []byte("Host: local host: 42069                                 \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "local host: 42069", headers["host"])
	assert.Equal(t, 58, n)
	assert.False(t, done)

	// Test: Valid 2 headers with existing headers
	headers = NewHeaders()
	headers["existing_header"] = "i was here first"
	data = []byte("Host: localhost:42069\r\nContent-Type: text/plain\r\n\r\n")

	totalConsumed := 0
	for {
		n, done, err = headers.Parse(data[totalConsumed:])
		require.NoError(t, err)
		totalConsumed += n
		if done {
			break
		}
	}
	require.NotNil(t, headers)
	assert.Equal(t, "i was here first", headers["existing_header"])
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, "text/plain", headers["content-type"])

	assert.Equal(t, len(data), totalConsumed)
	assert.True(t, done)

	// Test: valid done
	headers = NewHeaders()
	data = []byte("\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 2, n)
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid header
	headers = NewHeaders()
	data = []byte("Host localhost42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Case-insensitive duplicate header
	headers = NewHeaders()
	headers["content-type"] = "text/plain"
	data = []byte("Content-Type: application/json\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "application/json", headers["content-type"])
	assert.Equal(t, 1, len(headers))
	assert.False(t, done)

	// Test: Invalid char in field-name
	headers = NewHeaders()
	data = []byte("H@st: localhost42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}
