package main

import (
	"fmt"
	"log"
)

type HttpStatus struct {
	code    int
	message string
}

var (
	OK = HttpStatus{
		code:    200,
		message: "OK",
	}
	CREATED = HttpStatus{
		code:    201,
		message: "Created",
	}
	NOT_FOUND = HttpStatus{
		code:    404,
		message: "Not Found",
	}
	INTERNAL = HttpStatus{
		code:    500,
		message: "Error",
	}
)

type HttpResponse struct {
	Status  *HttpStatus
	Headers map[string]string
	Body    string
}

func (httpResponse *HttpResponse) ToString() []byte {
	if httpResponse.Status == nil {
		log.Println("http response missing status")
		// if missing status, set as internal server error and send that as response
		httpResponse.Status = &INTERNAL
		return []byte(
			fmt.Sprintf("HTTP/%v %v %v%v", serverInfo.httpVersion, httpResponse.Status.code, httpResponse.Status.message, CRLF),
		)
	}

	response := []byte(
		fmt.Sprintf("HTTP/%v %v %v%v", serverInfo.httpVersion, httpResponse.Status.code, httpResponse.Status.message, CRLF),
	)

	for key, val := range httpResponse.Headers {
		headerLine := fmt.Sprintf("%v: %v %v", key, val, CRLF)
		response = append(response, []byte(headerLine)...)
	}

	response = append(response, CRLF...)

	if httpResponse.Body != "" {
		response = append(response, httpResponse.Body...)
	}

	return response
}

func (httpResponse *HttpResponse) SetStatus(httpStatus *HttpStatus) {
	httpResponse.Status = httpStatus
}

func (httpResponse *HttpResponse) AppendHeader(key string, value string) {
	httpResponse.Headers[key] = value
}

func (httpResponse *HttpResponse) SetBody(value string) {
	httpResponse.Body = value
}
