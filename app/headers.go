package main

type Header string

const (
	headerContentType   Header = "Content-Type"
	headerContentLength Header = "Content-Length"
)

type ContentType string

const (
	contentTypePlainText ContentType = "text/plain"
)
