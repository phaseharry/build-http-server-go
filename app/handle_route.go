package main

import (
	"errors"
	"os"
	"path"
	"strings"
)

func handleRoute(request HttpRequest) HttpResponse {
	var response HttpResponse

	if request.RequestLine.Target == "/" {
		response = HttpResponse{
			Status: StatusOk,
		}
	} else if strings.HasPrefix(request.RequestLine.Target, "/echo/") {
		toEcho := strings.Split(request.RequestLine.Target, "/echo/")[1]
		response = HttpResponse{
			Status:      StatusOk,
			Body:        toEcho,
			ContentType: contentTypePlainText,
		}
	} else if strings.HasPrefix(request.RequestLine.Target, "/files/") {
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
			response = HttpResponse{
				Status:      StatusOk,
				Body:        string(data),
				ContentType: contentTypeApplicationOctectStream,
			}
		}
	} else if request.RequestLine.Target == "/user-agent" {
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
