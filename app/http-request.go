package main

import (
	"fmt"
	"strings"
)

type HttpRequest struct {
	Method      string
	Target      string
	HttpVersion string
	Headers     map[string]string
	Body        string
}

func NewHttpRequest(buffer []byte) (*HttpRequest, error) {
	httpRequest := HttpRequest{}
	request := strings.Split(string(buffer), CRLF)
	fmt.Printf("request len: %v", len(request))
	for _, val := range request {
		fmt.Printf("info: %v\n", val)
	}

	requestLine := request[0]
	parseRequestLine(&httpRequest, requestLine)

	i := 1
	prevLineEmptyString := false

	headers := make(map[string]string)
	for i < len(request) {
		headerLine := request[i]
		fmt.Printf("headerLine: %v\n", headerLine)
		i++
		// if prevLine was an empty string then this line is the request body
		if prevLineEmptyString {
			httpRequest.Body = headerLine
			continue
		}
		if headerLine == "" {
			prevLineEmptyString = true
			continue
		}
		headerLineValues := strings.Split(headerLine, ": ")
		key, value := headerLineValues[0], headerLineValues[1]
		headers[key] = value
	}
	httpRequest.Headers = headers
	fmt.Printf("request: %v\n", httpRequest)
	return &httpRequest, nil
}

func parseRequestLine(httpRequest *HttpRequest, requestLine string) error {
	requestLineValues := strings.Split(requestLine, " ")
	fmt.Printf("requestLineValues: %v, len: %v\n", requestLineValues, len(requestLineValues))
	if len(requestLineValues) != 3 {
		return fmt.Errorf("invalid request line")
	}
	httpRequest.Method = requestLineValues[0]
	httpRequest.Target = requestLineValues[1]
	httpRequest.HttpVersion = requestLineValues[2]
	fmt.Printf("method: %v, path: %v, httpVersion: %v\n", httpRequest.Method, httpRequest.Target, httpRequest.HttpVersion)
	return nil
}
