package main

import (
	"fmt"

	// Uncomment this block to pass the first stage
	"flag"
	"net"
	"os"
	"strings"
)

var directory = flag.String("directory", "CurrentDirectory", "Serves the files from this directory")

func main() {

	flag.Parse()

	fmt.Println("the given directory is" + *directory)
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go HandleFunc(c)
	}
}

func HandleFunc(c net.Conn) {

	data := make([]byte, 1000000)

	defer c.Close()

	_, err := c.Read(data)
	if err != nil {
		fmt.Println("Failed to read data: ", err.Error())
		os.Exit(1)
	}

	lines := strings.Split(string(data), "\r\n")

	splitStrings := strings.Split(lines[0], " ")
	pathval := splitStrings[1]
	method := splitStrings[0]

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

	} else if pathval == "/user-agent" {
		usrAgentString := lines[2][12:]
		result := len(usrAgentString)
		res := fmt.Sprintf("%d", result)

		c.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + res + "\r\n\r\n" + usrAgentString))

	} else if strings.HasPrefix(pathval, "/files") {
		fileName := strings.TrimPrefix(pathval, "/files/")

		filePath := *directory + "/" + fileName
		// ["get", "/file", "/filename/"]
		switch method {
		case "POST":

			_, err := os.Stat(filePath)
			if err != nil {
				fmt.Println("no file found in given directory ", filePath)
				c.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
				return
			}

			/* Get the index of the first string line which is empty
			   and the post request body is next line onwards
			*/

			indexOfReqBody := 0
			for index, line := range lines {
				if line == "" {
					indexOfReqBody = index + 1
				}
			}

			fileData := ""
			for i := indexOfReqBody; i <= (len(lines) - 1); i++ {
				fileData += lines[i]
				fileData += "\r\n"

				if i != (len(lines) - 1) {
					fileData += "\r\n"
				}

			}

			os.WriteFile(fileName, []byte(fileData), 0666)
			c.Write([]byte("HTTP/1.1 201 Created\r\n\r\n"))
			fmt.Println("Changes to file saved")

		case "GET":
			_, err := os.Stat(filePath)
			if err != nil {
				fmt.Println("no file found in given directory ", filePath)
				c.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
				return

			} else {
				readFile, err := os.ReadFile(filePath)
				if err != nil {
					fmt.Println("file found in given directory ", err)
					c.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
					return
				}

				readFileLength := len(readFile)
				readFileBody := string(readFile)
				readFileLengthStr := fmt.Sprintf("%d", readFileLength)

				c.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: " + readFileLengthStr + "\r\n\r\n" + readFileBody))

			}

		}

	} else {
		c.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}

}
