package main

import (
	"fmt"
	"log"
	"net"
)

func handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading for connection:", conn.RemoteAddr())
	}
	fmt.Println("Received data from connection:", conn.RemoteAddr(), string(buf))

	conn.Write([]byte("Hello from server"))

	conn.Close()
}

func main() {
	lsnr, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
	}
	defer lsnr.Close()

	fmt.Println("Listening on 8080")

	for {
		// Listen for an incoming connection.
		conn, err := lsnr.Accept()
		if err != nil {
			log.Fatal("Error accepting: ", err.Error())
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)

	}
}
