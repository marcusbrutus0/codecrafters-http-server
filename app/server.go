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
	} else if strings.HasPrefix(pathval, "/echo") {
		//trimming pathval by slash
		// /echo/abc/hugr -> abc/hugr
		randStr := strings.TrimPrefix(pathval, "/echo/")

		contentLength := len(randStr)
		contentLengthStr := fmt.Sprintf("%d", contentLength)
		/*
			HTTP/1.1 200 OK
			Content-Type: text/plain
			Content-Length: 3

			abc
		*/
		c.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + contentLengthStr + "\r\n\r\n" + randStr))

	} else {
		c.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}

}
