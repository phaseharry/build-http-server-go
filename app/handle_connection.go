
package main

import (
	"bufio"
	"errors"
	"io"
	"net"
	"strings"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	// creates a buffer to continously read bytes from open connections.
	// only close if the client sends a Connection: close header entry or
	// they hung up the kept-alive connection after they're done with it.
	reader := bufio.NewReader(conn)
	for {
		request, err := NewRequest(reader)
		// io.EOF = client hung up so we can return to close the connection
		if err != nil && errors.Is(err, io.EOF) {
			return
		} else if err != nil { // if there was an error reading the request, then it is an invalid request so send a 400 response
			badRequestResponse := HttpResponse{
				Status: StatusBadRequest,
			}
			conn.Write(badRequestResponse.ToBytes())
			return
		}

		response := handleRoute(request)
		closeConn := strings.EqualFold(request.Headers[requestHeaderConnection], connectionClose)
		if closeConn {
			response.Headers[requestHeaderConnection] = connectionClose
		}
		conn.Write(response.ToBytes())

		if closeConn {
			return
		}
	}
}
