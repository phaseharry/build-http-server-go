package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

func main() {
	directoryArg := flag.String("directory", directory, "directory where files are stored")
	flag.Parse()

	// setting the global directory path config so files can be read / written there
	directory = *directoryArg
	fmt.Println(*directoryArg)
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		go handleConnection(conn)
	}
}
