package main

import (
	"fan-out-subscription/client"
	"fan-out-subscription/server"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	s := server.New()
	go func() {
		// Try to establish connection to server
		for {
			conn, err := net.Dial("tcp", ":8080")
			if err != nil {
				log.Println("Waiting for server to be available...")
				time.Sleep(100 * time.Millisecond)
				continue
			}
			conn.Close()
			break
		}
		log.Println("Server is available")

		// Create and connect clients
		for i := 0; i < 10; i++ {
			conn, err := net.Dial("tcp", ":8080")
			if err != nil {
				log.Printf("Failed to connect client %d: %v", i, err)
				continue
			}
			log.Println("Creating client", i)
			go s.AcceptClient(client.NewClient(fmt.Sprintf("client-%d", i), conn))
		}
	}()

	if err := s.Start("8080", "8081"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
