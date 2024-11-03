package main

import (
	"fmt"
	"net"
	"os"
	"slices"
	"strconv"
	"strings"
)

const (
	CRLF = "\r\n"
	GZIP = "gzip"
)

var SUPPORTED_COMPRESSIONS = []string{GZIP}

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

	httpRequest, err := parseRequest(conn)
	if err != nil {
		return err
	}

	method := httpRequest.Method
	path := httpRequest.Target
	headers := httpRequest.Headers

	httpResponse := HttpResponse{
		Headers: map[string]string{},
	}
	switch {
	case method == "GET" && path == "/":
		httpResponse.SetStatus(&OK)

		fmt.Printf("response: %v\n", string(httpResponse.ToString()))

		conn.Write(httpResponse.ToString())

	case method == "GET" && strings.HasPrefix(path, "/echo"):
		responseValue := strings.TrimPrefix(path, "/echo/")
		responseSize := strconv.Itoa(len(responseValue))

		encoding := headers["Accept-Encoding"]

		httpResponse.SetStatus(&OK)
		if encoding != "" && slices.Contains(SUPPORTED_COMPRESSIONS, encoding) {
			httpResponse.AppendHeader("Content-Encoding", headers["Accept-Encoding"])
		}
		httpResponse.AppendHeader("Content-Type", "text/plain")
		httpResponse.AppendHeader("Content-Length", responseSize)
		httpResponse.SetBody(responseValue)
		fmt.Printf("response: %v\n", string(httpResponse.ToString()))
		conn.Write(httpResponse.ToString())

	case method == "GET" && strings.HasPrefix(path, "/user-agent"):
		userAgent := headers["User-Agent"]
		responseSize := strconv.Itoa(len(userAgent))

		httpResponse.SetStatus(&OK)
		httpResponse.AppendHeader("Content-Type", "text/plain")
		httpResponse.AppendHeader("Content-Length", responseSize)
		httpResponse.SetBody(userAgent)
		fmt.Printf("response: %v\n", string(httpResponse.ToString()))
		conn.Write(httpResponse.ToString())

	case method == "GET" && strings.HasPrefix(path, "/files"):
		// reading filepath and sending its content back to client
		filename := strings.TrimPrefix(path, "/files")
		directory := os.Args[2]
		file, err := os.ReadFile(directory + filename)
		if err != nil {
			httpResponse.SetStatus(&NOT_FOUND)
			conn.Write(httpResponse.ToString())
			return nil
		}
		responseSize := strconv.Itoa(len(file))

		httpResponse.SetStatus(&OK)
		httpResponse.AppendHeader("Content-Type", "application/octet-stream")
		httpResponse.AppendHeader("Content-Length", responseSize)
		httpResponse.SetBody(string(file))

		fmt.Printf("response: %v\n", string(httpResponse.ToString()))
		conn.Write(httpResponse.ToString())

	case method == "POST" && strings.HasPrefix(path, "/files"):
		filename := strings.TrimPrefix(path, "/files")
		directory := os.Args[2]
		requestBody := httpRequest.Body
		if err := os.WriteFile(directory+filename, []byte(strings.Trim(requestBody, "\x00")), 0644); err != nil {
			httpResponse.SetStatus(&INTERNAL)
			conn.Write(httpResponse.ToString())
			break
		}
		httpResponse.SetStatus(&CREATED)
		fmt.Printf("response: %v\n", string(httpResponse.ToString()))
		conn.Write(httpResponse.ToString())

	default:
		httpResponse.SetStatus(&NOT_FOUND)
		conn.Write(httpResponse.ToString())
	}
	return nil
}

func parseRequest(connection net.Conn) (HttpRequest, error) {
	buffer := make([]byte, 1024)

	httpRequest := HttpRequest{}

	_, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading connection request: ", err.Error())
		return httpRequest, err
	}
	request := strings.Split(string(buffer), CRLF)
	fmt.Printf("request len: %v", len(request))
	for _, val := range request {
		fmt.Printf("info: %v\n", val)
	}

	requestLine := request[0]
	requestLineValues := strings.Split(requestLine, " ")
	fmt.Printf("requestLineValues: %v, len: %v\n", requestLineValues, len(requestLineValues))
	if len(requestLineValues) != 3 {
		return httpRequest, fmt.Errorf("invalid request line")
	}
	httpRequest.Method = requestLineValues[0]
	httpRequest.Target = requestLineValues[1]
	httpRequest.HttpVersion = requestLineValues[2]
	fmt.Printf("method: %v, path: %v, httpVersion: %v\n", httpRequest.Method, httpRequest.Target, httpRequest.HttpVersion)

	i := 1
	prevLineEmptyString := false

	headers := make(map[string]string)
	for i < len(request) {
		headerLine := request[i]
		fmt.Printf("headerLine: %v\n", headerLine)
		i++
		// if prevLine was an empty string then this line is the request body
		if prevLineEmptyString {
			httpRequest.Body = headerLine
			continue
		}
		if headerLine == "" {
			prevLineEmptyString = true
			continue
		}
		headerLineValues := strings.Split(headerLine, ": ")
		key, value := headerLineValues[0], headerLineValues[1]
		headers[key] = value
	}
	httpRequest.Headers = headers
	fmt.Printf("request: %v\n", httpRequest)
	return httpRequest, nil
}
