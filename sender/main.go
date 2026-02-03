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
		// send on stream 1
		err := common.WriteFrame(conn, 1, []byte("stream1 message "+time.Now().Format(time.RFC3339)))
		if err != nil {
			log.Fatal(err)
		}

		// send on stream 2
		err = common.WriteFrame(conn, 2, []byte("stream2 message "+time.Now().Format(time.RFC3339)))
		if err != nil {
			log.Fatal(err)
		}

		i++
		time.Sleep(1 * time.Second)
	}
}
