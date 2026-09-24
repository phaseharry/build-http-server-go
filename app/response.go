package main

import (
	"fmt"
	"strconv"
)

type HttpResponse struct {
	StatusCode StatusCode
	Body       string
}

func (h *HttpResponse) ToBytes() []byte {
	var b []byte

	// Status line information
	// sending protocol header info
	b = append(b, fmt.Sprintf("HTTP/%s", VERSION)...)
	// appending empty byte for a space
	b = append(b, ' ')
	// appending status code
	b = append(b, strconv.Itoa(h.StatusCode.Code)...)
	b = append(b, ' ')
	// appending status code phase (OK, UNAUTHORIZED, etc.)
	b = append(b, h.StatusCode.Phrase...)

	// Header information
	b = append(b, CRLF...)
	// Response Body
	b = append(b, CRLF...)
	fmt.Println(string(b))
	return b
}
