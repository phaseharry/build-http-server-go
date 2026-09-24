package main

type StatusCode struct {
	Code   int
	Phrase string
}

var StatusCodeOk = StatusCode{
	Code:   200,
	Phrase: "OK",
}
