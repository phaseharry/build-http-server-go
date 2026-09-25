package main

import (
	"fmt"
	"strconv"
)

type HttpResponse struct {
	Status      Status
	ContentType ContentType
	Body        string
}

func NewHttpResponse(
	status Status,
	contentType ContentType,
	body string,
) HttpResponse {
	return HttpResponse{
		Status:      status,
		ContentType: contentType,
		Body:        body,
	}
}

func (h *HttpResponse) ToBytes() []byte {
	var responseBytes []byte
	body := []byte(h.Body)

	// Status line information
	// sending protocol header info
	responseBytes = append(responseBytes, fmt.Sprintf("HTTP/%s", VERSION)...)
	// appending empty byte for a space
	responseBytes = append(responseBytes, ' ')
	// appending status code
	responseBytes = append(responseBytes, strconv.Itoa(h.Status.Code)...)
	responseBytes = append(responseBytes, ' ')
	// appending status code phase (OK, UNAUTHORIZED, etc.)
	responseBytes = append(responseBytes, h.Status.Phrase...)
	responseBytes = append(responseBytes, CRLF...)

	// Header information
	if h.ContentType != "" {
		responseBytes = append(responseBytes, fmt.Sprintf("%s: ", headerContentType)...)
		responseBytes = append(responseBytes, []byte(h.ContentType)...)
		responseBytes = append(responseBytes, CRLF...)
	}
	if len(body) != 0 {
		responseBytes = append(responseBytes, fmt.Sprintf("%s: ", headerContentLength)...)
		responseBytes = append(responseBytes, fmt.Sprintf("%d", len(body))...)
		responseBytes = append(responseBytes, CRLF...)
	}
	responseBytes = append(responseBytes, CRLF...)

	// Response Body
	responseBytes = append(responseBytes, []byte(h.Body)...)
	fmt.Println(string(responseBytes))
	return responseBytes
}
