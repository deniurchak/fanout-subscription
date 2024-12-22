# fanout-subscription

A simple fan-out subscription service implemented in Go that demonstrates:

- TCP server accepting multiple client connections, by default on port 8080
- HTTP endpoint to trigger mock event sending, by default on port 8081
- Broadcasting events to all connected clients
- Client state management and deduplication
- Thread-safe concurrent operations

## Architecture

The service consists of:

- Server: Manages client connections and broadcasts events
- Client: Maintains connection and handles incoming events
- Event: Data structure for messages passed through the system

## Usage

Install dependencies:

```bash
go mod tidy
```

Start the server:

```bash
go run cmd/main.go
```
