package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
)

func handleRoute(request HttpRequest) HttpResponse {
	var response HttpResponse

	if request.RequestLine.Target == "/" && request.RequestLine.Method == GET {
		response = HttpResponse{
			Status: StatusOk,
		}
	} else if strings.HasPrefix(request.RequestLine.Target, "/echo/") && request.RequestLine.Method == GET {
		toEcho := strings.Split(request.RequestLine.Target, "/echo/")[1]
		response = HttpResponse{
			Status:      StatusOk,
			Body:        toEcho,
			ContentType: contentTypePlainText,
		}
	} else if strings.HasPrefix(request.RequestLine.Target, "/files/") && request.RequestLine.Method == GET {
		filename := strings.Split(request.RequestLine.Target, "/files/")[1]
		filepath := path.Join(directory, filename)
		data, err := os.ReadFile(filepath)
		fmt.Println("hello")
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
			response = HttpResponse{
				Status:      StatusOk,
				Body:        string(data),
				ContentType: contentTypeApplicationOctectStream,
			}
		}
	} else if strings.HasPrefix(request.RequestLine.Target, "/files/") && request.RequestLine.Method == POST {
		filename := strings.Split(request.RequestLine.Target, "/files/")[1]
		filepath := path.Join(directory, filename)

		data := request.Body
		// if content-length was not sent as part of the request, write the entire body including
		// the empty bytes availalble due to HttpRequest's 1024 byte buffer.
		contentLength, ok := request.Headers[string(headerContentLength)]
		if ok {
			size, err := strconv.Atoi(contentLength)
			if err == nil {
				data = data[:size]
			}
		}
		err := os.WriteFile(filepath, request.Body, 0644)
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
		response = HttpResponse{
			Status:      StatusOk,
			Body:        request.Headers["User-Agent"],
			ContentType: contentTypePlainText,
		}
	} else {
		response = HttpResponse{
			Status: StatusNotFound,
		}
	}

	return response
}
