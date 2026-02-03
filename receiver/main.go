package main

import (
	"log"
	"net"
        "time"
	"p2_ssh_stream/common"
)

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:9000")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Receiver listening on 9000")

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handle(conn)
	}
}

func handle(c net.Conn) {
	defer c.Close()
	lastHeartbeat := time.Now()

	for {
		streamID, msg, err := common.ReadFrame(c)
		if err != nil {
			log.Println("connection closed:", err)
			return
		}

		if streamID == 0 {
			log.Println("Received heartbeat:", string(msg))
			lastHeartbeat = time.Now()
			continue
		}

		log.Printf("[Stream %d] %s\n", streamID, msg)

		// check if heartbeat timeout
		if time.Since(lastHeartbeat) > 10*time.Second {
			log.Println("No heartbeat received in 10s, closing connection")
			return
		}
	}
}
