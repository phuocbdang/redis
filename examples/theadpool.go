package main

import (
	"io"
	"log"
	"net"
	"time"
)

// Element in the queue
type Job struct {
	conn net.Conn
}

// Represent the thread in the pool
type Worker struct {
	id      int
	jobChan chan Job
}

func NewWorker(id int, jobChan chan Job) *Worker {
	return &Worker{
		id:      id,
		jobChan: jobChan,
	}
}

func (w *Worker) Start() {
	go func() {
		for job := range w.jobChan {
			handleConnection(job.conn)
		}
	}()
}

// Represent the thread pool
type Pool struct {
	jobQueue chan Job
	workers  []*Worker
}

func NewPool(numOfWorker int) *Pool {
	return &Pool{
		jobQueue: make(chan Job),
		workers:  make([]*Worker, numOfWorker),
	}
}

func (p *Pool) Start() {
	for i := 0; i < len(p.workers); i++ {
		worker := NewWorker(i, p.jobQueue)
		p.workers[i] = worker
		worker.Start()
	}
}

func (p *Pool) AddJob(conn net.Conn) {
	p.jobQueue <- Job{conn: conn}
}

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

	pool := NewPool(2)
	pool.Start()

	log.Println("Server listening on :3000")

	for {
		conn, err := listener.Accept() // conn == socket == dedicated communication channel
		if err != nil {
			log.Fatal(err)
		}

		pool.AddJob(conn)
	}
}
