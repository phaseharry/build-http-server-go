package main

import (
	"fmt"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
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

	response := handleRoute(request)

	if _, err := conn.Write(response.ToBytes()); err != nil {
		fmt.Println(err.Error())
	}
}
