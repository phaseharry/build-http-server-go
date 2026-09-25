package main

import "errors"

var ErrInvalidRequest = errors.New("invalid HTTP request")

var ErrInvalidEncodingFormat = errors.New("invalid encoding format")
