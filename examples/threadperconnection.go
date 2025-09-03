package main

import (
	"io"
	"log"
	"net"
	"time"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		// Set a deadline to prevent a client from holding the connection open indefinitely
		conn.SetReadDeadline(time.Now().Add(time.Second * 10))

		// Read data from the client. Using a smaller buffer for demonstration.
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Println("Client close the connection.")
			} else {
				log.Println("Read error: ", err)
			}
			return // Exit the loop and close the connection
		}

		requestData := buf[:n]
		log.Printf("Received data: %s", string(requestData))

		// Process the request
		time.Sleep(time.Second * 1)

		// Reply to the client
		_, err = conn.Write([]byte("Hello world"))
		if err != nil {
			log.Println("Write error:", err)
			return // Exit on write error
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Println("Server listening on :3000")

	for {
		conn, err := listener.Accept() // conn == socket == dedicated communication channel
		if err != nil {
			log.Fatal(err)
		}

		go handleConnection(conn) // Use a goroutine to handle each connection concurrently
	}
}
