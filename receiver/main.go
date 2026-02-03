package main

import (
	"log"
	"net"
	"time"

	"p2_ssh_stream/common"
)

func main() {
	metrics := common.NewMetrics()

	// Log receiver metrics every 5 seconds
	go func() {
		for {
			time.Sleep(5 * time.Second)
			metrics.Lock()
			heartbeatDelay := time.Since(metrics.LastHeartbeat)
			if heartbeatDelay > 10*time.Second {
				log.Println("WARNING: Heartbeat delay > 10s")
			}
			log.Printf("Receiver Metrics: BytesReceived=%v, FramesReceived=%v, Heartbeats=%d\n",
				metrics.BytesReceived, metrics.FramesReceived, metrics.Heartbeats)
			metrics.Unlock()
		}
	}()

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
		go handle(conn, metrics)
	}
}

func handle(c net.Conn, metrics *common.Metrics) {
	defer c.Close()
	for {
		streamID, msg, err := common.ReadFrameCompressed(c)
		if err != nil {
			log.Println("Connection closed:", err)
			return
		}

		if streamID == 0 {
			log.Println("Received heartbeat:", string(msg))
			metrics.Heartbeat()
			continue
		}

		log.Printf("[Stream %d] %s\n", streamID, msg)
		metrics.Received(streamID, len(msg))
	}
}
