package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
)

type StatusCode int

const (
	OK                  = 200
	BadRequest          = 400
	InternalServerError = 500
)

type Writer struct {
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	reason := ""

	switch statusCode {
	case OK:
		reason = "OK"
	case BadRequest:
		reason = "Bad Request"
	case InternalServerError:
		reason = "Internal Server Error"
	}

	statusLine := fmt.Sprintf("HTTP/1.1 %d %s \r\n", statusCode, reason)
	_, err := w.Write([]byte(statusLine))

	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")

	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		h := fmt.Sprintf("%s: %s\r\n", k, v)
		_, err := w.Write([]byte(h))
		if err != nil {
			return err
		}
	}

	_, err := w.Write([]byte("\r\n"))
	return err
}
