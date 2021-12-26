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

func lobby(heartbeat chan int) {
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

		go handleConnection(c, heartbeat)
	}
}

func handleConnection(conn net.Conn, heartbeat chan int) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	go func() {
		for {
			<-heartbeat
			fmt.Println("CONN heartbeat")
			writer.Flush()
		}
	}()

	for {
		buffer, err := reader.ReadBytes('\n')

		if err != nil {
			fmt.Println("Client disconnected")
			return
		}

		input := string(buffer)
		input = strings.ReplaceAll(input, "\r", "")
		input = strings.ReplaceAll(input, "\n", "")

		fmt.Println("Client message:", input)
		fmt.Println("Length:", len(input))

		bytes, err := writer.WriteString("Server received: " + input)

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

func heartbeat(heartbeatChan chan int) {
	const rate = 2000

	for {
		time.Sleep(rate * time.Millisecond)
		fmt.Println("===== Heartbeat")

		heartbeatChan <- 1
	}
}

func main() {
	heartbeatChan := make(chan int)

	go lobby(heartbeatChan)
	go heartbeat(heartbeatChan)

	for {
	}
}
