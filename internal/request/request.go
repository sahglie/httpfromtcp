package request

import (
	"bytes"
	"errors"
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strconv"
	"strings"
)

type state int

const (
	initialized state = iota
	done
	parsingHeaders
	parsingBody
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
	state       state
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const crlf = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	req := Request{
		state:   initialized,
		Headers: headers.NewHeaders(),
		Body:    make([]byte, 0),
	}

	buf := make([]byte, 8)

	numBytesRead := 0

	for req.state != done {
		if numBytesRead >= len(buf) {
			tmp := make([]byte, len(buf)*2)
			copy(tmp, buf)
			buf = tmp
		}

		nBytes, err := reader.Read(buf[numBytesRead:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if req.state != done {
					return nil, fmt.Errorf("incomplete request, in state: %d, read n bytes on EOF: %d", req.state, numBytesRead)
				}
				break
			}

			return nil, err
		}

		numBytesRead += nBytes
		numBytesParsed, err := req.parse(buf[:numBytesRead])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[numBytesParsed:])
		numBytesRead -= numBytesParsed
	}

	return &req, nil
}

func parseRequestLine(rawBytes []byte) (*RequestLine, int, error) {
	idx := bytes.Index(rawBytes, []byte(crlf))
	if idx == -1 {
		return nil, 0, nil
	}

	startLine := string(rawBytes[:idx])
	startLine = strings.TrimSpace(startLine)

	numBytesRead := idx + 2

	parts := strings.Split(startLine, " ")
	if len(parts) != 3 {
		return nil, numBytesRead, fmt.Errorf("poorly formatted request line %s", startLine)
	}

	method := parts[0]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, numBytesRead, fmt.Errorf("invalid method: %s", method)

		}
	}

	target := parts[1]

	versionParts := strings.Split(parts[2], "/")
	if (len(versionParts) != 2) || (versionParts[0] != "HTTP") || (versionParts[1] != "1.1") {
		return nil, numBytesRead, fmt.Errorf("invalid HTTP-version")
	}

	requestLine := &RequestLine{
		HttpVersion:   versionParts[1],
		RequestTarget: target,
		Method:        method,
	}

	return requestLine, numBytesRead, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.state != done {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		totalBytesParsed += n
		if n == 0 {
			break
		}
	}
	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.state {
	case initialized:
		requestLine, n, err := parseRequestLine(data)
		if err != nil {
			// something actually went wrong
			return 0, err
		}
		if n == 0 {
			// just need more data
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.state = parsingHeaders
		return n, nil
	case parsingHeaders:
		n, isDone, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if isDone {
			r.state = parsingBody
		}
		return n, nil
	case parsingBody:
		size := r.Headers.Get("Content-Length")
		if size == "" {
			r.state = done
			return len(data), nil
		}

		r.Body = append(r.Body, data...)
		s, err := strconv.Atoi(size)
		if err != nil {
			return 0, err
		}

		if len(r.Body) > s {
			return 0, fmt.Errorf("body is longer than content-length\n")
		}
		if len(r.Body) == s {
			r.state = done
		}

		return len(data), nil
	case done:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("unknown state")
	}
}
