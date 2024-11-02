package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var serverInfo = struct {
	host        string
	port        string
	httpVersion string
}{
	host:        "0.0.0.0",
	port:        "4221",
	httpVersion: "1.1",
}

func main() {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%s", serverInfo.host, serverInfo.port))
	if err != nil {
		fmt.Printf("Failed to bind to port %s", serverInfo.port)
		os.Exit(1)
	}
	fmt.Printf("%v", os.Args)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) error {
	defer conn.Close()

	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading connection request: ", err.Error())
		return err
	}
	request := strings.Split(string(buffer), "\r\n")

	for _, val := range request {
		fmt.Printf("info: %v\n", val)
	}

	requestLine, userAgentString := request[0], request[2]
	requestLineValues := strings.Split(requestLine, " ")
	method, path, httpVersion := requestLineValues[0], requestLineValues[1], requestLineValues[2]
	fmt.Printf("userAgentString: %v\n", userAgentString)
	fmt.Printf("method: %v, path: %v, httpVersion: %v\n", method, path, httpVersion)

	switch {
	case method == "GET" && path == "/":
		conn.Write([]byte("HTTP/" + serverInfo.httpVersion + " 200 OK\r\n\r\n"))

	case method == "GET" && strings.HasPrefix(path, "/echo"):
		responseValue := strings.TrimPrefix(path, "/echo/")
		fmt.Printf("echo response value: %v", responseValue)
		responseSize := strconv.Itoa(len(responseValue))
		conn.Write([]byte("HTTP/" + serverInfo.httpVersion + " 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + responseSize + "\r\n\r\n" + responseValue))

	case method == "GET" && strings.HasPrefix(path, "/user-agent"):
		userAgentValues := strings.Split(userAgentString, " ")
		fmt.Printf("%v", userAgentValues)
		userAgent := userAgentValues[len(userAgentValues)-1]
		responseSize := strconv.Itoa(len(userAgent))
		conn.Write([]byte("HTTP/" + serverInfo.httpVersion + " 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + responseSize + "\r\n\r\n" + userAgent))

	case method == "GET" && strings.HasPrefix(path, "/files"):
		// reading filepath and sending its content back to client
		filename := strings.TrimPrefix(path, "/files")
		directory := os.Args[2]
		file, err := os.ReadFile(directory + filename)
		if err != nil {
			conn.Write([]byte("HTTP/" + serverInfo.httpVersion + " 404 Not Found\r\n\r\n"))
			return nil
		}
		responseSize := strconv.Itoa(len(file))
		conn.Write([]byte("HTTP/" + serverInfo.httpVersion + " 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: " + responseSize + "\r\n\r\n" + string(file)))

	case method == "POST" && strings.HasPrefix(path, "/files"):
		filename := strings.TrimPrefix(path, "/files")
		directory := os.Args[2]
		requestBody := request[len(request)-1] // body is always the last part of a request
		fmt.Printf("%v", requestBody)
		if err := os.WriteFile(directory+filename, []byte(strings.Trim(requestBody, "\x00")), 0644); err != nil {
			conn.Write([]byte("HTTP/" + serverInfo.httpVersion + " 500 Error\r\n\r\n"))
			break
		}
		conn.Write([]byte("HTTP/" + serverInfo.httpVersion + " 201 Created\r\n\r\n"))

	default:
		conn.Write([]byte("HTTP/" + serverInfo.httpVersion + " 404 Not Found\r\n\r\n"))
	}
	return nil
}
