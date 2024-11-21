package main

import (
	"fmt"
	"net"
	"os"

	"github.com/Ultimate-conscious/redis-FS/resp"
)

var _ = net.Listen
var _ = os.Exit

func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		resp_ := resp.NewResp(conn)
		val, err := resp_.Read()
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			break
		}
		fmt.Println("Received: ", val)
		conn.Write([]byte("+OK\r\n"))

	}
}

func eventLoop(conn <-chan net.Conn) {
	for conn := range conn {
		go handleConnection(conn)
	}
}

func main() {

	fmt.Println("Hello, playground")

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
