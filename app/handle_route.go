package main

import "strings"

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
