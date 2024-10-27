package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type serverInfo struct {
	host        string
	port        string
	httpVersion string
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	serverInfo := serverInfo{
		host:        "0.0.0.0",
		port:        "4221",
		httpVersion: "1.1",
	}

	// Uncomment this block to pass the first stage
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%s", serverInfo.host, serverInfo.port))
	if err != nil {
		fmt.Printf("Failed to bind to port %s", serverInfo.port)
		os.Exit(1)
	}

	for {
		connection, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		request := make([]byte, 1024)
		_, err = connection.Read(request)
		if err != nil {
			fmt.Println("Error reading connection request: ", err.Error())
			os.Exit(1)
		}
		requestInfo := strings.Split(string(request), "\r\n")

		requestLine := requestInfo[0]
		requestLineValues := strings.Split(requestLine, " ")
		method, path, httpVersion := requestLineValues[0], requestLineValues[1], requestLineValues[2]

		fmt.Printf("method: %v, path: %v, httpVersion: %v", method, path, httpVersion)
		switch {
		case method == "GET" && path == "/":
			connection.Write([]byte("HTTP/" + serverInfo.httpVersion + " 200 OK\r\n\r\n"))
		case method == "GET" && strings.HasPrefix(path, "/echo/"):
			pathValues := strings.Split(path, "/")
			responseValue := pathValues[len(pathValues)-1]
			responseSize := strconv.Itoa(len(responseValue))
			connection.Write([]byte("HTTP/" + serverInfo.httpVersion + " 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + responseSize + "\r\n\r\n" + responseValue))
		default:
			connection.Write([]byte("HTTP/" + serverInfo.httpVersion + " 404 Not Found\r\n\r\n"))
		}
	}
}
