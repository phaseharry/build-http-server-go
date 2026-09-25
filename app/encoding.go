package main

import (
	"bytes"
	"compress/gzip"
	"errors"
)

const (
	GZIP = "gzip"
)

var SUPPORTED_ENCODINGS = map[string]struct{}{
	GZIP: {},
}

func encode(encoding string, data []byte) ([]byte, error) {
	var buf bytes.Buffer
	var encodedData []byte
	if _, ok := SUPPORTED_ENCODINGS[encoding]; !ok {
		return encodedData, ErrInvalidEncodingFormat
	}

	if encoding == GZIP {
		gzipWriter := gzip.NewWriter(&buf)
		if _, err := gzipWriter.Write(data); err != nil {
			return encodedData, errors.New("error gziping data")
		}
		gzipWriter.Close()
		encodedData = buf.Bytes()
	}

	return encodedData, nil
}
