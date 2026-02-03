package main

import (
	"log"
	"net"
	"time"

	"p2_ssh_stream/common"
)

func connect() net.Conn {
	for {
		conn, err := net.Dial("tcp", "127.0.0.1:8000")
		if err != nil {
			log.Println("Failed to connect, retrying in 2s:", err)
			time.Sleep(2 * time.Second)
			continue
		}
		log.Println("Connected to receiver")
		return conn
	}
}

func main() {
	stream1 := make(chan string, 5)
	stream2 := make(chan string, 5)

	// producer goroutines
	go func() {
		for {
			stream1 <- "stream1 message " + time.Now().Format(time.RFC3339)
			time.Sleep(500 * time.Millisecond)
		}
	}()

	go func() {
		for {
			stream2 <- "stream2 message " + time.Now().Format(time.RFC3339)
			time.Sleep(700 * time.Millisecond)
		}
	}()

	conn := connect()
	defer conn.Close()

	// heartbeat goroutine
	go func() {
		for {
			time.Sleep(3 * time.Second)
			if err := common.WriteFrameCompressed(conn, 0, []byte("PING")); err != nil {
				log.Println("Heartbeat failed, reconnecting:", err)
				conn.Close()
				conn = connect()
			}
		}
	}()

	// sender loop
	for {
		select {
		case msg := <-stream1:
			if err := common.WriteFrameCompressed(conn, 1, []byte(msg)); err != nil {
				log.Println("Write failed, reconnecting:", err)
				conn.Close()
				conn = connect()
				common.WriteFrame(conn, 1, []byte(msg))
			}
		case msg := <-stream2:
			if err := common.WriteFrameCompressed(conn, 2, []byte(msg)); err != nil {
				log.Println("Write failed, reconnecting:", err)
				conn.Close()
				conn = connect()
				common.WriteFrame(conn, 2, []byte(msg))
			}
		}
	}
}
