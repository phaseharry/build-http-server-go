package main

import (
	"errors"
	"os"
	"path"
	"strconv"
	"strings"
)

func handleRoute(request HttpRequest) HttpResponse {
	var response HttpResponse
	responseHeaders := make(map[string]string)

	if request.RequestLine.Target == "/" && request.RequestLine.Method == GET {
		response = HttpResponse{
			Status: StatusOk,
		}
	} else if strings.HasPrefix(request.RequestLine.Target, "/echo/") && request.RequestLine.Method == GET {
		toEcho := strings.Split(request.RequestLine.Target, "/echo/")[1]
		responseHeaders[headerContentType] = contentTypePlainText
		response = HttpResponse{
			Status:  StatusOk,
			Body:    toEcho,
			Headers: responseHeaders,
		}
	} else if strings.HasPrefix(request.RequestLine.Target, "/files/") && request.RequestLine.Method == GET {
		filename := strings.Split(request.RequestLine.Target, "/files/")[1]
		filepath := path.Join(directory, filename)
		data, err := os.ReadFile(filepath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				response = HttpResponse{
					Status: StatusNotFound,
				}
			} else {
				response = HttpResponse{
					Status: StatusInternalServerError,
				}
			}
		} else { // success case
			responseHeaders[headerContentType] = contentTypeApplicationOctectStream
			response = HttpResponse{
				Status:  StatusOk,
				Headers: responseHeaders,
				Body:    string(data),
			}
		}
	} else if strings.HasPrefix(request.RequestLine.Target, "/files/") && request.RequestLine.Method == POST {
		filename := strings.Split(request.RequestLine.Target, "/files/")[1]
		filepath := path.Join(directory, filename)

		data := request.Body
		// if content-length was not sent as part of the request, write the entire body including
		// the empty bytes availalble due to HttpRequest's 1024 byte buffer.
		contentLength, ok := request.Headers[headerContentLength]
		if ok {
			size, err := strconv.Atoi(contentLength)
			if err == nil {
				data = data[:size]
			}
		}
		err := os.WriteFile(filepath, data, 0644)
		if err != nil {
			response = HttpResponse{
				Status: StatusInternalServerError,
			}
		} else {
			response = HttpResponse{
				Status: StatusCreated,
			}
		}
	} else if request.RequestLine.Target == "/user-agent" && request.RequestLine.Method == GET {
		responseHeaders[headerContentType] = contentTypePlainText
		response = HttpResponse{
			Status:  StatusOk,
			Body:    request.Headers["User-Agent"],
			Headers: responseHeaders,
		}
	} else {
		response = HttpResponse{
			Status: StatusNotFound,
		}
	}

	if encoding, ok := request.Headers[requestHeaderAcceptEncoding]; ok {
		encoded, err := encode(encoding, []byte(response.Body))
		// if there's an error with encoding with a supported encoding format, send an internal server response.
		// if there's an error because client sent an unsupported encoding format, just send back the raw response with no encoding.
		// if encoding is supported and it encoded successfully, send back the encoding response
		if err != nil && !errors.Is(err, ErrInvalidEncodingFormat) {
			response = HttpResponse{
				Status: StatusInternalServerError,
			}
		} else if err == nil { // successfully encoded so use the encoded data
			response.Body = string(encoded)
			response.Headers[responseContentEncoding] = encoding
		}
	}

	return response
}
