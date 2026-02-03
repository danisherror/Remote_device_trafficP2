
#  Multiplexed SSH Tunnel in Go

## Overview

**P2** is a lightweight SSH-based tunnel written in **Go** that allows multiple logical streams of data to flow over a single TCP/SSH connection.
It is designed for **reliable, efficient, and observable data streaming**, with features like:

* Multiplexed streams
* Backpressure / flow control
* Heartbeats and tunnel health
* Automatic reconnect and session resume
* Compression of payloads
* Metrics and observability

This project is ideal for testing **custom data streams, logs, or traffic simulation** over SSH tunnels.

---

## Features

| Feature                 | Description                                                              |
| ----------------------- | ------------------------------------------------------------------------ |
| Multiplexing            | Multiple logical streams over one TCP connection (stream IDs)            |
| Backpressure            | Buffered channels per stream prevent memory spikes when receiver is slow |
| Heartbeats              | Detect dead connections using periodic PING frames                       |
| Auto-reconnect          | Automatically reconnect and resume sending messages if the tunnel fails  |
| Compression             | Compress payloads using gzip to reduce bandwidth                         |
| Metrics & Observability | Track bytes sent/received, frames per stream, and heartbeat health       |

---

## Folder Structure

```
p2-ssh-stream/
├── common/            # Shared utilities (frame handling, compression, metrics)
│   └── frame.go
├── receiver/          # Receiver application
│   └── main.go
├── sender/            # Sender application
│   └── main.go
├── go.mod
└── README.md
```

---

## Getting Started

### Prerequisites

* Go 1.25 or later
* SSH installed (macOS / Linux)

---

### Running Receiver

```bash
cd ~/Desktop/projects/p2-ssh-stream
go run ./receiver
```

* Listens on `127.0.0.1:9000` by default
* Receives multiplexed, compressed frames and prints them per stream

---

### Running Sender

```bash
cd ~/Desktop/projects/p2-ssh-stream
go run ./sender
```

* Connects to `127.0.0.1:8000` via SSH tunnel
* Sends multiple streams with compression, backpressure, and heartbeats

---

### Running with SSH Tunnel

```bash
ssh -N -L 8000:127.0.0.1:9000 localhost
```

* `-L` forwards local port `8000` to receiver’s port `9000`
* All streams are sent over a single secure SSH connection

---

## Example Output

```
[Stream 1] stream1 message 2026-02-03T20:15:00Z
[Stream 2] stream2 message 2026-02-03T20:15:00Z
[Stream 0] PING
[Stream 1] stream1 message 2026-02-03T20:15:01Z
[Stream 2] stream2 message 2026-02-03T20:15:01Z
```

* Stream `0` → Heartbeat
* Other streams → Logical data streams

---

## How It Works

1. **Sender** produces messages per stream (channels with backpressure)
2. Messages are **compressed** (gzip) and framed `[streamID][length][payload]`
3. Sender writes frames to a TCP connection tunneled over SSH
4. **Receiver** reads frames, decompresses, and routes messages by stream ID
5. Heartbeats detect dead connections, triggering automatic reconnects
6. Metrics track per-stream traffic and connection health

---

## Extending the Project

* Add **more logical streams** dynamically
* Add **custom compression algorithms** (zstd, snappy)
* Add **Prometheus/Grafana metrics export**
* Implement **stream prioritization or QoS**

---
## Optional enhancements for P2:

1. Export metrics over HTTP for Prometheus/Grafana
2. Add logging levels (info/debug/warning)
3. Add configuration for stream IDs, heartbeat intervals, and max reconnect retries





