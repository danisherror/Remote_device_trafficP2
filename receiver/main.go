package main

import (
	"log"
	"net"

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
	for {
		streamID, msg, err := common.ReadFrame(c)
		if err != nil {
			log.Println("connection closed:", err)
			return
		}
		log.Printf("[Stream %d] %s\n", streamID, msg)
	}
}
