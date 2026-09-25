package main

type Status struct {
	Code   int
	Phrase string
}

var StatusOk = Status{
	Code:   200,
	Phrase: "OK",
}

var StatusBadRequest = Status{
	Code:   400,
	Phrase: "Bad Request",
}

var StatusNotFound = Status{
	Code:   404,
	Phrase: "Not Found",
}

var StatusInternalServerError = Status{
	Code:   500,
	Phrase: "Internal Server Error",
}
