package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	requestBytes := make([]byte, 1024)
	conn.Read(requestBytes)

	request, err := NewRequest(requestBytes)
	if err != nil {
		badRequestResponse := HttpResponse{
			Status: StatusNotFound,
		}
		conn.Write(badRequestResponse.ToBytes())
		return
	}
	fmt.Println(request)

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

	if _, err := conn.Write(response.ToBytes()); err != nil {
		fmt.Println(err.Error())
	}
}
