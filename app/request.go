package main

import (
	"bytes"
	"fmt"
)

const (
	GET  = "GET"
	POST = "POST"
)

type HttpRequest struct {
	RequestLine requestLine
	Headers     map[string]string
	Body        []byte
}

type requestLine struct {
	Method  string
	Target  string
	Version string
}

func NewRequest(req []byte) (HttpRequest, error) {
	httpRequest := HttpRequest{}
	requestParts := bytes.Split(req, []byte(CRLF))

	// the first CRLF is always the request line
	requestLineParts := bytes.Split(requestParts[0], []byte(" "))
	if len(requestLineParts) != 3 {
		return httpRequest, ErrInvalidRequest
	}
	method, target, version := requestLineParts[0], requestLineParts[1], requestLineParts[2]

	reqLine := requestLine{
		Method:  string(method),
		Target:  string(target),
		Version: string(version),
	}

	httpRequest.RequestLine = reqLine

	// every header entry get its own CRLF
	headers := make(map[string]string)
	for i := 1; i < len(requestParts)-1; i++ {
		headerEntry := requestParts[i]
		// bytes.Cut returns the first occurrance of the seperator bytes and returns 3 values.
		// left of the seperator, right of the seperator,
		// bool of whether a seperator even exists or not
		key, value, found := bytes.Cut(headerEntry, []byte(": "))
		if !found {
			continue
		}
		headers[string(key)] = string(value)
	}
	fmt.Println(headers)
	httpRequest.Headers = headers
	// the entry after the last CRLF is the request body if it exists
	httpRequest.Body = requestParts[len(requestParts)-1]

	return httpRequest, nil
}
