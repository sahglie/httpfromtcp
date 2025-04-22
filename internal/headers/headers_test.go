package headers

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestHeaders_Parse_ValidSingleHeader(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers.Get("Host"), fmt.Sprintf("%#v", headers))
	assert.Equal(t, 23, n)
	assert.False(t, done)
}

func TestHeaders_Parse_InvalidSpacingHeader(t *testing.T) {
	headers := NewHeaders()
	data := []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

func TestHeaders_Parse_ValidSingleHeaderExtraWhitespace(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("            Host: localhost:42069              \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers.Get("Host"))
	assert.Equal(t, 49, n)
	assert.False(t, done)
}

func TestHeaders_Parse_Valid2Headers(t *testing.T) {
	// Test: Valid 2 headers with existing headers
	headers := Headers(map[string]string{"host": "localhost:42069"})

	data := []byte("User-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers.Get("host"))
	assert.Equal(t, "curl/7.81.0", headers.Get("User-Agent"))
	assert.Equal(t, 25, n)
	assert.False(t, done)
}

func TestHeaders_Parse_CaseInsensitive(t *testing.T) {
	headers := NewHeaders()
	data := []byte("User-agenT: curl/7.81.0\r\nAcCepT: */*\r\n\r\n")

	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "curl/7.81.0", headers.Get("user-agent"))
	assert.Equal(t, 25, n)
	assert.False(t, done)

	n, done, err = headers.Parse(data[n:])
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "*/*", headers.Get("accept"))
	assert.Equal(t, 13, n)
	assert.False(t, done)
}

func TestHeaders_Parse_HandlesDuplicateHeaders(t *testing.T) {
	headers := Headers(map[string]string{"set-person": "lane-loves-go"})
	data := []byte("Set-Person: prime-loves-zig\r\n\r\n")

	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "lane-loves-go, prime-loves-zig", headers.Get("Set-Person"), fmt.Sprintf("%#v", headers))
	assert.Equal(t, 29, n)
	assert.False(t, done)
}

func TestHeaders_Parse_RequiresValidTokenCharacterSet(t *testing.T) {
	headers := NewHeaders()
	data := []byte("(User)-Agent: curl/7.81.0\r\n\r\n")

	n, done, err := headers.Parse(data)
	require.Error(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "", headers.Get("user-agent"))
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

func TestShit(t *testing.T) {
	target := "/httpbin/stream/100"
	parts := strings.Split(target, "/")
	fmt.Printf("%#v\n", parts)

	path := strings.TrimPrefix(target, "/httpbin/")
	fmt.Printf("%#v\n", path)
}
