package main

import (
	"fmt"
	"net"
	"os"
)

var _ = net.Listen
var _ = os.Exit

func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, 1024)

		_, err := conn.Read(buf)

		if err != nil {
			fmt.Println("Error reading:", err.Error())
			break
		}
		conn.Write([]byte("+PONG\r\n"))

	}
}

func eventLoop(conn <-chan net.Conn) {
	for conn := range conn {
		go handleConnection(conn)
	}
}

func main() {

	fmt.Println("Logs from your program will appear here!")

	lsnr, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	eventQueue := make(chan net.Conn)

	go eventLoop(eventQueue)

	for {
		conn, err := lsnr.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		eventQueue <- conn
	}

}
