package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

type Status int16

const (
	HOST                              = "localhost"
	PORT                              = "8080"
	TYPE                              = "tcp"
	HTTP_1_1                          = "HTTP/1.1"
	HTTP_GET                          = "GET"
	HTTP_OK                    Status = 200
	HTTP_NOTFOUND              Status = 404
	HTTP_VERSION_NOT_SUPPORTED Status = 505
)

func main() {
	listen, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	// close listener
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	// incoming request
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	texts := strings.Fields(string(buf[:]))
	verb := texts[0]
	path := texts[1]
	protocol := texts[2]
	var responseStartLine string
	var responseStr string
	if protocol != HTTP_1_1 {
		responseStartLine = fmt.Sprintf("%v %v %v\n", HTTP_1_1, 505, "HTTP Version Not Supported")
		conn.Write([]byte(responseStartLine))
	} else if verb != HTTP_GET {
		responseStartLine = fmt.Sprintf("%v %v %v\n", HTTP_1_1, 501, "Not Implemented")
		conn.Write([]byte(responseStartLine))
	} else {
		fileContent, err := readFileFromPath(path)
		if err != nil {
			responseStartLine = fmt.Sprintf("%v %v %v\n", HTTP_1_1, 404, "Not Found")
			conn.Write([]byte(responseStartLine))
		} else {
			responseStartLine = fmt.Sprintf("%v %v %v\n", HTTP_1_1, 200, "OK")
			responseStr = fmt.Sprintf("%v\n%v", responseStartLine, fileContent)
			conn.Write([]byte(responseStr))
		}
	}
	fmt.Printf("%v %v %v -- %v", verb, path, protocol, responseStartLine)

	// close conn
	conn.Close()
}

func readFileFromPath(path string) (string, error) {
	splitPath := strings.Split(path, "#")
	splitPath = strings.Split(splitPath[0], "?")
	splitPath = strings.Split(splitPath[0], "/")
	filename := splitPath[len(splitPath)-1]
	if !strings.Contains(filename, ".") {
		splitPath = append(splitPath, "index.html")
	}
	path = strings.Join(splitPath, "/")
	path = fmt.Sprintf("html/%v", path)
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
    // close fi on exit and check for its returned error
    defer func() {
        if err := file.Close(); err != nil {
            log.Fatal(err)
        }
    }()
	fileInfo, err := file.Stat()
	if err != nil {
		return "", err
	}
	fileSize := fileInfo.Size()
	// make a buffer to keep chunks that are read
    buf := make([]byte, fileSize)
    for {
        // read a chunk
        n, err := file.Read(buf)
        if err != nil && err != io.EOF {
            return "", err
        }
        if n == 0 {
            break
        }
    }
	return string(buf), nil
}
