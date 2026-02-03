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

	i := 0
	for {
		msg := []byte("message #" + time.Now().Format(time.RFC3339))
		err := common.WriteFrame(conn, msg)
		if err != nil {
			log.Fatal(err)
		}

		i++
		time.Sleep(1 * time.Second)
	}
}
