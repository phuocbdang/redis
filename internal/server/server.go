package server

import (
	"io"
	"log"
	"net"
	"phuocbdang/internal/config"
	"phuocbdang/internal/core/io_multiplexing"
	"syscall"
)

// func readCommand(fd int) (*core.Command, error) {
// 	buf := make([]byte, 512)
// 	n, err := syscall.Read(fd, buf)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if n == 0 {
// 		return nil, io.EOF
// 	}
// 	return core

// }

func RunIOMultiplexingServer() {
	log.Println("Starting an I/O Multiplexing TCP server on", config.Port)
	listener, err := net.Listen(config.Protocol, config.Port)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	// Get the file descriptor from the client
	tcpListener, ok := listener.(*net.TCPListener)
	if !ok {
		log.Fatal("Listener is not a TCPListener")
	}
	listenerFile, err := tcpListener.File()
	if err != nil {
		log.Fatal(err)
	}
	defer listenerFile.Close()

	serverFD := int(listenerFile.Fd())

	// Create an I/O Multiplexing instance (epoll in Linux, kqueue in MacOS)
	ioMultiplexing, err := io_multiplexing.CreateIOMultiplexer()
	if err != nil {
		log.Fatal(err)
	}
	defer ioMultiplexing.Close()

	// Monitor "read" events on the Server FD
	if err = ioMultiplexing.Monitor(io_multiplexing.Event{
		Fd: serverFD,
		Op: io_multiplexing.OpRead,
	}); err != nil {
		log.Fatal(err)
	}

	events := make([]io_multiplexing.Event, config.MaxConnection)
	for {
		// Wait for file descriptors in the monitoring list to be ready for I/O
		events, err := ioMultiplexing.Wait() // Blocking call
		if err != nil {
			continue
		}

		for i := 0; i < len(events); i++ {
			if events[i].Fd == serverFD {
				log.Print("New client is trying to connect")
				connFD, _, err := syscall.Accept(serverFD)
				if err != nil {
					log.Print("- err: ", err)
					continue
				}
				log.Print("- setting up a new connection")
				if err = ioMultiplexing.Monitor(io_multiplexing.Event{
					Fd: connFD,
					Op: io_multiplexing.OpRead,
				}); err != nil {
					log.Fatal(err)
				}
			} else {
				cmd, err := readCommand(events[i].Fd)
				if err != nil {
					if err == io.EOF || err == syscall.ECONNRESET {
						log.Println("Client disconnected")
						_ = syscall.Close(events[i].Fd)
						continue
					}
					log.Println("Read error: ", err)
					continue
				}
			}
		}
	}
}
