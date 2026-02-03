# Remote_device_trafficP2
---

#Step  Run apps

Terminal 1:

```bash
go run receiver/main.go
```

Terminal 2:

```bash
ssh -N -L 8000:127.0.0.1:9000 localhost
```

Terminal 3:

```bash
go run sender/main.go
```

