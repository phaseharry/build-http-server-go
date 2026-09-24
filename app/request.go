package main

import "bytes"

type HttpRequest struct {
	RequestLine requestLine
	Headers     map[string]string
	Body        string
}

type requestLine struct {
	Method  string
	Target  string
	Version string
}

func NewRequest(req []byte) (HttpRequest, error) {
	httpRequest := HttpRequest{}
	requestParts := bytes.Split(req, []byte(CRLF))

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

	return httpRequest, nil
}
