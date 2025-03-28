package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type state int

const (
	initialized state = 0
	done        state = 1
)

type Request struct {
	RequestLine RequestLine
	state       state
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const crlf = "\r\n"

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case initialized:
		reqLine, err, numBytesRead := parseRequestLine(data)
		if err != nil {
			return 0, err
		}

		if numBytesRead == 0 {
			return 0, nil
		}

		r.RequestLine = *reqLine
		r.state = done

		return numBytesRead, nil
	case done:
		return 0, fmt.Errorf("tryint to read data in a done state")
	default:
		return 0, fmt.Errorf("unknown state")
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, 8)
	readToIndex := 0
	req := &Request{state: initialized}

	for req.state != done {
		if readToIndex >= len(buf) {
			tmp := make([]byte, len(buf)*2)
			copy(tmp, buf)
			buf = tmp
		}

		numBytesRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(io.EOF, err) {
				req.state = done
				break
			}
			return nil, err
		}

		readToIndex += numBytesRead
		numBytesParsed, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[numBytesParsed:])
		readToIndex -= numBytesParsed
	}

	return req, nil
}

func parseRequestLine(data []byte) (*RequestLine, error, int) {
	idx := bytes.Index(data, []byte(crlf))

	if idx == -1 {
		return nil, nil, 0
	}

	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, err, 0
	}
	return requestLine, nil, idx + 2
}

func requestLineFromString(str string) (*RequestLine, error) {
	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("poorly formatted request-line: %s", str)
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
