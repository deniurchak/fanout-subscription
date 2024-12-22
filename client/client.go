package client

import (
	"encoding/json"
	"fan-out-subscription/event"
	"log"
	"net"
	"sync"
)

type Event = event.Event

// Client represents a connected client and its state
type Client struct {
	id              string
	conn            net.Conn
	processedEvents map[string]bool
	state           map[string]interface{}
	mu              sync.RWMutex
}

func (c *Client) ID() string {
	return c.id
}

func NewClient(id string, conn net.Conn) *Client {
	return &Client{
		id:              id,
		conn:            conn,
		processedEvents: make(map[string]bool),
		state:           make(map[string]interface{}),
		mu:              sync.RWMutex{},
	}
}

func (c *Client) GetMutex() *sync.RWMutex {
	return &c.mu
}

func (c *Client) CloseConn() {
	c.conn.Close()
}

func (c *Client) Read(b []byte) (int, error) {
	return c.conn.Read(b)
}

func (c *Client) Write(b []byte) (int, error) {
	return c.conn.Write(b)
}

func (c *Client) HandleEvent(event Event) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.processedEvents[event.ID] {
		return nil // Skip duplicate events
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	log.Printf("Sending event to client %s: %s", c.ID(), string(data))

	if _, err := c.Write(append(data, '\n')); err != nil {
		return err
	}

	c.processedEvents[event.ID] = true
	c.state[event.Type] = event.Data
	return nil
}
