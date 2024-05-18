package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"testing"
)

const (
	TEST_HOST = "localhost"
	TEST_PORT = "8080"
	TEST_TYPE = "tcp"
)

func startServer() func() {
	// Create a listener on the given host and port
	listener, err := net.Listen(TEST_TYPE, TEST_HOST+":"+TEST_PORT)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	var wg sync.WaitGroup

	// Goroutine to accept connections
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			conn, err := listener.Accept()
			if err != nil {
				return // Stop accepting connections when listener is closed
			}
			go handleRequest(conn)
		}
	}()

	// Function to stop the server
	stopServer := func() {
		listener.Close()
		wg.Wait()
	}

	// Return the stop function
	return stopServer
}

func connectToServer(t *testing.T) net.Conn {
	conn, err := net.Dial(TEST_TYPE, TEST_HOST+":"+TEST_PORT)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	return conn
}

func sendRequest(t *testing.T, conn net.Conn, request string) string {
	fmt.Fprintf(conn, request)
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}
	return response
}

func TestGetExistingFile(t *testing.T) {
	stopServer := startServer()
	defer stopServer()

	conn := connectToServer(t)
	defer conn.Close()

	request := "GET /test.html HTTP/1.1\r\nHost: localhost\r\n\r\n"
	expected := "HTTP/1.1 200 OK"

	response := sendRequest(t, conn, request)
	if !strings.HasPrefix(response, expected) {
		t.Errorf("Expected response to start with %q, got %q", expected, response)
	}

	// Read until the end of the response
	for {
		line, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			break
		}
		if strings.TrimSpace(line) == "" {
			break
		}
	}
}

func TestGetNonExistentFile(t *testing.T) {
	stopServer := startServer()
	defer stopServer()

	conn := connectToServer(t)
	defer conn.Close()

	request := "GET /nonexistent.html HTTP/1.1\r\nHost: localhost\r\n\r\n"
	expected := "HTTP/1.1 404 Not Found"

	response := sendRequest(t, conn, request)
	if !strings.HasPrefix(response, expected) {
		t.Errorf("Expected response to start with %q, got %q", expected, response)
	}

	// Read until the end of the response
	for {
		line, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			break
		}
		if strings.TrimSpace(line) == "" {
			break
		}
	}
}

func TestUnsupportedHTTPVersion(t *testing.T) {
	stopServer := startServer()
	defer stopServer()

	conn := connectToServer(t)
	defer conn.Close()

	request := "GET / HTTP/1.0\r\nHost: localhost\r\n\r\n"
	expected := "HTTP/1.1 505 HTTP Version Not Supported"

	response := sendRequest(t, conn, request)
	if !strings.HasPrefix(response, expected) {
		t.Errorf("Expected response to start with %q, got %q", expected, response)
	}

	// Read until the end of the response
	for {
		line, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			break
		}
		if strings.TrimSpace(line) == "" {
			break
		}
	}
}
