package main

import (
	"log"
	"net"
	"time"

	"p2_ssh_stream/common"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// buffered channel for backpressure
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

        go func() {
            for {
                time.Sleep(3 * time.Second) // heartbeat interval
                    if err := common.WriteFrame(conn, 0, []byte("PING")); err != nil {
                        log.Println("Heartbeat failed:", err)
                            return
                    }
            }
        }()

	// sender loop
	for {
		select {
		case msg := <-stream1:
			if err := common.WriteFrame(conn, 1, []byte(msg)); err != nil {
				log.Fatal(err)
			}
		case msg := <-stream2:
			if err := common.WriteFrame(conn, 2, []byte(msg)); err != nil {
				log.Fatal(err)
			}
		}
	}
}
