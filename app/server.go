package main

import (
	"fmt"
	// Uncomment this block to pass the first stage
	"net"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	data := make([]byte, 1024)

	for {
		HandleFunc(l, data)
	}
}

func HandleFunc(l net.Listener, data []byte) {

	c, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	defer c.Close()

	_, err = c.Read(data)
	if err != nil {
		fmt.Println("Failed to read data: ", err.Error())
		os.Exit(1)
	}

	lines := strings.Split(string(data), "\r\n")

	pathval := strings.Split(lines[0], " ")[1]

	if pathval == "/" {
		c.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else {
		c.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
	c.Close()
}
