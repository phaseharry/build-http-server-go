
package main

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
)

const (
	GET  = "GET"
	POST = "POST"
)

type HttpRequest struct {
	RequestLine requestLine
	Headers     headers
	Body        []byte
}

type headers map[string]string

type requestLine struct {
	Method  string
	Target  string
	Version string
}

func NewRequest(reader *bufio.Reader) (HttpRequest, error) {
	// request line + headers bytes only.
	// reads until the blank line that only contains the CRLF and then looks at the incoming Content-Length header to determine if there was a body within the request so we can read those bytes in as well.
	// that how we'll if the following bytes after the request line + header still belongs to this request or is it part of another request on the same TCP connection.
	var head []byte
	request := HttpRequest{}
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return request, err
		}
		head = append(head, line...)
		// if the line only contains the CRLF then that's the end of the request line + header bytes
		if string(line) == CRLF {
			break
		}
	}

	reqLine, headers, err := parseHeadBytes(head)
	if err != nil {
		return request, err
	}
	request.RequestLine = reqLine
	request.Headers = headers

	// reading potential body to the HttpRequest using the Content-Length header value to know how many bytes to read
	if contentLength, ok := request.Headers[headerContentLength]; ok {
		size, err := strconv.Atoi(contentLength)
		if err != nil || size < 0 {
			return request, ErrInvalidRequest
		}
		// allocating a slice of just the right size and calling io.ReadFull to continuously
		// read bytes into request.Body until it has filled up the slice.
		request.Body = make([]byte, size)
		if _, err := io.ReadFull(reader, request.Body); err != nil {
			return request, err
		}
	}

	return request, nil
}

func parseHeadBytes(head []byte) (requestLine, headers, error) {
	reqLine, headers := requestLine{}, make(map[string]string)
	requestParts := bytes.Split(head, []byte(CRLF))

	// the first CRLF is always the request line
	requestLineParts := bytes.Split(requestParts[0], []byte(" "))
	if len(requestLineParts) != 3 {
		return reqLine, headers, ErrInvalidRequest
	}

	reqLine.Method = string(requestLineParts[0])
	reqLine.Target = string(requestLineParts[1])
	reqLine.Version = string(requestLineParts[2])

	// every header entry get its own CRLF
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

	return reqLine, headers, nil
}
