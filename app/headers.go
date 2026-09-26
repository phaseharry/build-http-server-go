
package main

// common headers shared by both Request and Response
const (
	headerContentType   = "Content-Type"
	headerContentLength = "Content-Length"
)

// request headers
const (
	requestHeaderAccept         = "Accept"
	requestHeaderAcceptEncoding = "Accept-Encoding"
	requestHeaderConnection     = "Connection"
)

// response headers
const (
	responseContentEncoding = "Content-Encoding"
)

// Connection
const (
	connectionClose     = "close"
	connectionKeepAlive = "keep-alive"
)

// Content-Type
const (
	contentTypePlainText               = "text/plain"
	contentTypeApplicationOctectStream = "application/octet-stream"
)
