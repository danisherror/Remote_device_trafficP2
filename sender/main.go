package main

import (
	"log"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for {
		conn.Write([]byte("hello over ssh tunnel\n"))
		time.Sleep(1 * time.Second)
	}
}
