package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net"
	"os"
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

	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)

	if err != nil {
		fmt.Println("Error reading connection request: ", err.Error())
		return err
	}

	httpRequest, err := NewHttpRequest(buffer)
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
		httpResponse.SetStatus(OK)
		conn.Write(httpResponse.ToString())

	case method == "GET" && strings.HasPrefix(path, "/echo"):
		responseValue := strings.TrimPrefix(path, "/echo/")

		handleEncoding(headers["Accept-Encoding"], responseValue, &httpResponse)
		httpResponse.SetStatus(OK)
		conn.Write(httpResponse.ToString())

	case method == "GET" && strings.HasPrefix(path, "/user-agent"):
		userAgent := headers["User-Agent"]
		responseSize := strconv.Itoa(len(userAgent))

		httpResponse.SetStatus(OK)
		httpResponse.AppendHeader("Content-Type", "text/plain")
		httpResponse.AppendHeader("Content-Length", responseSize)
		httpResponse.SetBody(userAgent)
		conn.Write(httpResponse.ToString())

	case method == "GET" && strings.HasPrefix(path, "/files"):
		// reading filepath and sending its content back to client
		filename := strings.TrimPrefix(path, "/files")
		directory := os.Args[2]
		file, err := os.ReadFile(directory + filename)
		if err != nil {
			httpResponse.SetStatus(NOT_FOUND)
			conn.Write(httpResponse.ToString())
			return nil
		}
		responseSize := strconv.Itoa(len(file))

		httpResponse.SetStatus(OK)
		httpResponse.AppendHeader("Content-Type", "application/octet-stream")
		httpResponse.AppendHeader("Content-Length", responseSize)
		httpResponse.SetBody(string(file))

		conn.Write(httpResponse.ToString())

	case method == "POST" && strings.HasPrefix(path, "/files"):
		filename := strings.TrimPrefix(path, "/files")
		directory := os.Args[2]
		requestBody := httpRequest.Body
		if err := os.WriteFile(directory+filename, []byte(strings.Trim(requestBody, "\x00")), 0644); err != nil {
			httpResponse.SetStatus(INTERNAL)
			conn.Write(httpResponse.ToString())
			break
		}
		httpResponse.SetStatus(CREATED)
		conn.Write(httpResponse.ToString())

	default:
		httpResponse.SetStatus(NOT_FOUND)
		conn.Write(httpResponse.ToString())
	}
	return nil
}

func handleEncoding(encodingsHeaderStr string, content string, httpResponse *HttpResponse) {
	encodings := strings.Split(encodingsHeaderStr, ", ")
	encodedResponse := false
	for _, encoding := range encodings {
		fmt.Printf("encoding Val: %v\n", encoding)
		// only support GZIP for now
		if encoding == GZIP {
			encodedResponse = true
			var buffer bytes.Buffer
			gzipWriter := gzip.NewWriter(&buffer)
			gzipWriter.Write([]byte(content))
			gzipWriter.Close()

			compressedContent := buffer.String()
			httpResponse.AppendHeader("Content-Encoding", encoding)
			httpResponse.AppendHeader("Content-Type", "text/plain")
			httpResponse.AppendHeader("Content-Length", strconv.Itoa(len(compressedContent)))
			httpResponse.SetBody(compressedContent)
			break
		}
	}

	// if no encoding was done then set response as regular content
	if !encodedResponse {
		httpResponse.AppendHeader("Content-Type", "text/plain")
		httpResponse.AppendHeader("Content-Length", strconv.Itoa(len(content)))
		httpResponse.SetBody(content)
	}
}
