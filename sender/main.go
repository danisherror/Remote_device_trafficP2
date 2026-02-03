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

	metrics := common.NewMetrics()
	conn := connect()
	defer conn.Close()

	// Producer goroutines
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

	// Heartbeat goroutine
	go func() {
		for {
			time.Sleep(3 * time.Second)
			if err := common.WriteFrameCompressed(conn, 0, []byte("PING")); err != nil {
				log.Println("Heartbeat failed, reconnecting:", err)
				conn.Close()
				metrics.Reconnect()
				conn = connect()
			} else {
				metrics.Heartbeat()
			}
		}
	}()

	// Metrics logging goroutine
	go func() {
		for {
			time.Sleep(5 * time.Second)
			metrics.Lock()
			log.Printf("Sender Metrics: BytesSent=%v, FramesSent=%v, Heartbeats=%d, Reconnects=%d\n",
				metrics.BytesSent, metrics.FramesSent, metrics.Heartbeats, metrics.Reconnects)
			metrics.Unlock()
		}
	}()

	// Sender loop
	for {
		select {
		case msg := <-stream1:
			if err := common.WriteFrameCompressed(conn, 1, []byte(msg)); err != nil {
				log.Println("Write failed, reconnecting:", err)
				conn.Close()
				metrics.Reconnect()
				conn = connect()
				common.WriteFrameCompressed(conn, 1, []byte(msg))
			} else {
				metrics.Sent(1, len(msg))
			}

		case msg := <-stream2:
			if err := common.WriteFrameCompressed(conn, 2, []byte(msg)); err != nil {
				log.Println("Write failed, reconnecting:", err)
				conn.Close()
				metrics.Reconnect()
				conn = connect()
				common.WriteFrameCompressed(conn, 2, []byte(msg))
			} else {
				metrics.Sent(2, len(msg))
			}
		}
	}
}
