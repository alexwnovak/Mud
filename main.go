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

func lobby() {
	fmt.Println("Lobby started, waiting for connections")

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
	defer conn.Close()

	for {
		reader, err := bufio.NewReader(conn).ReadBytes('\n')
		writer := bufio.NewWriter(conn)

		if err != nil {
			fmt.Println("Client disconnected")
			return
		}

		input := string(reader)
		input = strings.ReplaceAll(input, "\r", "")
		input = strings.ReplaceAll(input, "\n", "")

		fmt.Println("Client message:", input)
		fmt.Println("Length:", len(input))

		bytes, err := writer.WriteString("Server received: " + input)
		writer.Flush()

		fmt.Println("Bytes written to client:", bytes)

		if err != nil {
			fmt.Println("Write error:", err.Error())
		}

		if input == "/quit" {
			fmt.Println("Disconnecting client...")
			return
		}
	}
}

func main() {
	go lobby()

	for {
	}
}
