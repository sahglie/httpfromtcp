package headers

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

type Headers map[string]string

const crlf = "\r\n"

func NewHeaders() Headers {
	return make(map[string]string)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))

	if idx == 0 {
		return 2, true, nil
	}

	if idx == -1 {
		return 0, false, nil
	}

	parts := strings.SplitN(string(data[:idx]), ":", 2)
	key := parts[0]

	if key != strings.TrimRight(key, " ") {
		return 0, false, fmt.Errorf("invalid header name: %s", key)
	}

	value := strings.TrimSpace(parts[1])
	key = strings.TrimSpace(key)

	if !validHeaderChars(key) {
		return 0, false, fmt.Errorf("header contains invalid character: %s", key)
	}

	h.Set(key, value)

	return idx + 2, done, err
}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	val, exists := h[key]
	if exists {
		value = fmt.Sprintf("%s, %s", val, value)
	}

	h[key] = value
}

func (h Headers) Replace(key, value string) {
	key = strings.ToLower(key)
	h[key] = value
}

func (h Headers) Delete(key string) {
	delete(h, key)
}

func (h Headers) Get(key string) string {
	return h[strings.ToLower(key)]
}

const ValidHeaderSpecialChars = "!#$%&'*+-.^_`|~]+"

func validHeaderChars(chars string) bool {
	for _, c := range chars {
		if !unicode.IsLetter(c) &&
			!unicode.IsDigit(c) &&
			!strings.Contains(ValidHeaderSpecialChars, string(c)) {
			return false
		}
	}

	return true
}
