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

// This function listens for incoming connections. When connections arrive, they run on their
// own go routines that is a dedicated loop for that specific client.
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
		go client(c)
	}
}

// This is the go routine that services one client connection. This receives data from the
// client, applies game logic, then writes output.
func client(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		buffer, err := reader.ReadBytes('\n')

		if err != nil {
			fmt.Println("Client disconnected")
			return
		}

		input := string(buffer)
		input = strings.ReplaceAll(input, "\r", "")
		input = strings.ReplaceAll(input, "\n", "")

		writer.WriteString("Server received: " + input + "\n")
		writer.Flush()

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

	<-make(chan bool)
}
