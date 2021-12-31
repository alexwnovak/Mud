package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

const (
	connHost = "localhost"
	connPort = "8080"
	connType = "tcp"
)

// This function listens for incoming connections. When connections arrive, they run on their
// own go routines that is a dedicated loop for that specific client.
func lobby(shutdown chan int) {
	fmt.Println("Lobby started, waiting for connections")

	l, err := net.Listen(connType, connHost+":"+connPort)

	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	listener := l.(*net.TCPListener)
	defer listener.Close()

	var exit bool

	go func() {
		<-shutdown
		exit = true
	}()

	for {
		fmt.Println("Calling Accept")

		listener.SetDeadline(time.Now().Add(time.Second * 1))
		c, err := listener.Accept()

		if err != nil {
			if exit {
				return
			}

			continue
		}

		fmt.Println("Client connected:", c.RemoteAddr().String())
		go client(c, shutdown)
	}
}

// This is the go routine that services one client connection. This receives data from the
// client, applies game logic, then writes output.
func client(conn net.Conn, shutdown chan int) {
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

		if input == "/shutdown" {
			shutdown <- 1
		}
	}
}
