package main

import (
	"fmt"
	"net"
	"os"
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

	// reading only one connection and returning an http 200
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
	if requestLine != "GET / HTTP/1.1" {
		connection.Write([]byte("HTTP/" + serverInfo.httpVersion + " 404 Not Found\r\n\r\n"))
		return
	}

	// fmt.Println(requestString)
	connection.Write([]byte("HTTP/" + serverInfo.httpVersion + " 200 OK\r\n\r\n"))
}
