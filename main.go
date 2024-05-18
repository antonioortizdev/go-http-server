package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type Status int16

const (
	HOST                       = "localhost"
	PORT                       = "8080"
	TYPE                       = "tcp"
	HTTP_1_1                   = "HTTP/1.1"
	HTTP_GET                   = "GET"
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
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}
	texts := strings.Fields(string(buffer[:]))
	verb := texts[0]
	path := texts[1]
	protocol := texts[2]
	var responseStr string
	if protocol != HTTP_1_1 {
		responseStr = fmt.Sprintf("%v %v %v\n", HTTP_1_1, 505, "HTTP Version Not Supported")
		conn.Write([]byte(responseStr))
	} else if verb != HTTP_GET {
		responseStr = fmt.Sprintf("%v %v %v\n", HTTP_1_1, 501, "Not Implemented")
		conn.Write([]byte(responseStr))
	} else {
		responseStr = fmt.Sprintf("%v %v %v\n", HTTP_1_1, 404, "Not Found")
		conn.Write([]byte(responseStr))
	}
	fmt.Printf("%v %v %v -- %v", verb, path, protocol, responseStr)

	// close conn
	conn.Close()
}
