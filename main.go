package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	connHost = "localhost"
	connPort = "8080"
	connType = "tcp"
)

func main() {
	fmt.Println("Starting", connType, "server on", connHost, ":", connPort)

	l, err := net.Listen(connType, connHost+":"+connPort)

	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	defer l.Close()

	for {
		c, err := l.Accept()

		if err != nil {
			fmt.Println("Error connecting:", err.Error())
			return
		}

		fmt.Println("Client connected:", c.RemoteAddr().String())

		go handleConnection(c)
	}
}

func handleConnection(conn net.Conn) {
	for {
		buffer, err := bufio.NewReader(conn).ReadBytes('\n')

		if err != nil {
			fmt.Println("Client disconnected")
			conn.Close()
			return
		}

		input := string(buffer)
		input = strings.ReplaceAll(input, "\r", "")
		input = strings.ReplaceAll(input, "\n", "")

		fmt.Println("Client message:", input)
		fmt.Println("Length:", len(input))

		if input == "/quit" {
			fmt.Println("Disconnecting client...")
			conn.Close()
			return
		}
	}
}
