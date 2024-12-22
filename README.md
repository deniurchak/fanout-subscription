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

## Before production

Before deploying this to production, we might want to:

### Features
- Add endpoint to post a proper message with configurable body
- Make max number of clients configurable, use connection pool

### Robustness
- Put in proper error handling, gracefully handle errors and panics
- Make clients automatically reconnect when they drop (ie with exponential backoff timing)
- Handle clients joining and leaving with separate endpoints

### Scalability
- Add a message queue (RabbitMQ or Kafka) so we don't lose messages when clients are offline (if message volume is high)
- Add rate limiting so that no single client can use up all the resources

### Security
- Add proper auth 
- Use TLS/SSL

### Monitoring
- Add proper logging 
- Use metrics (Prometheus + Grafana)
- Add health checks for the server
- Set up alerts 

### Testing
- Write some unit tests
- Add integration tests for Db and message queue if we use them

