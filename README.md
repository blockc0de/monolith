# Monolith
Golang implementation of the blockc0de engine server.

## Prerequisites
- Redis server

## How to start
Copy `etc/monolith-api.yaml.example` to `monolith-api.yaml`, Modify your redis server conf in the `etc/monolith-api.yaml` file, and run it:
```bash
go run monolith.go
```
