package main

type HttpRequest struct {
	Method      string
	Target      string
	HttpVersion string
	Headers     map[string]string
	Body        string
}
