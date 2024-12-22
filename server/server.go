package server

import (
	"encoding/json"
	"fan-out-subscription/client"
	"fan-out-subscription/event"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Client = client.Client
type Event = event.Event

// Server manages client connections and event broadcasting
type Server struct {
	clients    map[string]*Client
	mu         sync.RWMutex
	eventsChan chan Event
}

func New() *Server {
	return &Server{
		clients:    make(map[string]*client.Client),
		eventsChan: make(chan Event, 100),
	}
}

// AcceptClient adds a new client to the server and starts handling its connection
func (s *Server) AcceptClient(c *Client) {
	s.mu.Lock()
	s.clients[c.ID()] = c
	s.mu.Unlock()
	log.Printf("New client connected: %s", c.ID())

	// Simple ping-pong to keep connection alive
	buffer := make([]byte, 1024)
	for {
		log.Printf("Client %s ping", c.ID())
		_, err := c.Read(buffer)
		if err != nil {
			log.Printf("Client %s disconnected: %v", c.ID(), err)
			return
		}
	}
}

func (s *Server) HandleEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse event from query parameters
	event := Event{
		ID:        uuid.New().String(),
		Type:      "update",
		Data:      map[string]interface{}{"value": time.Now().Unix()},
		Timestamp: time.Now(),
	}

	// Validate required fields
	if event.ID == "" || event.Type == "" {
		http.Error(w, "Missing required event fields", http.StatusBadRequest)
		return
	}

	// Send event to channel for processing
	s.eventsChan <- event

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "Event received",
	})
}

func (s *Server) Start(tcpPort string, httpPort string) error {
	// Register HTTP handler
	http.HandleFunc("/event", s.HandleEvent)

	// Start TCP server for client connections
	tcpListener, err := net.Listen("tcp", ":"+tcpPort)
	if err != nil {
		return fmt.Errorf("failed to start TCP server: %v", err)
	}

	// Start HTTP server on a different port for the API endpoint
	httpListener, err := net.Listen("tcp", ":"+httpPort)
	if err != nil {
		return fmt.Errorf("failed to start HTTP server: %v", err)
	}

	// Start HTTP server
	go func() {
		if err := http.Serve(httpListener, nil); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	go s.handleEvents()

	log.Printf("Server listening on ports %s (TCP) and %s (HTTP)", tcpPort, httpPort)

	// Accept connections
	for {
		conn, err := tcpListener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		log.Printf("Accepted connection %v", conn)
	}
}

func (s *Server) handleEvents() {
	for event := range s.eventsChan {
		s.broadcastEvent(event)
	}
}

func (s *Server) broadcastEvent(event Event) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, client := range s.clients {
		go func(c *Client, e Event) {
			log.Printf("Broadcasting event to client %s %v", c.ID(), e)
			if err := c.HandleEvent(e); err != nil {
				log.Printf("Error sending event to client %s: %v", c.ID(), err)
			}
		}(client, event)
	}
}
