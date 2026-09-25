package main

import (
	"bytes"
	"compress/gzip"
	"errors"
)

const (
	GZIP = "gzip"
)

func encode(encoding string, data []byte) ([]byte, error) {
	var buf bytes.Buffer
	var encodedData []byte
	if encoding == GZIP {
		gzipWriter := gzip.NewWriter(&buf)
		if _, err := gzipWriter.Write(data); err != nil {
			return encodedData, errors.New("error gziping data")
		}
		encodedData = buf.Bytes()
	} else {
		return encodedData, ErrInvalidEncodingFormat
	}
	return encodedData, nil
}
