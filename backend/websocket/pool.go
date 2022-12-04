package websocket

import (
	"log"
)

// A Pool acts as room, which keep tracks of users that belong to specific chat
type Pool struct {
	ID         string
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan map[string]interface{}
}

func NewPool() *Pool {
	return &Pool{
		ID:         "",
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan map[string]interface{}),
	}
}

// Start method handles adding and removing users to pool
// and also broadcasting user messages to all pool clients
func (p *Pool) Start() {
	for {
		select {
		case client := <-p.Register:
			p.Clients[client] = true
			log.Printf("Client %+v joined the pool %v", client, p.ID)

		case client := <-p.Unregister:
			delete(p.Clients, client)

		case message := <-p.Broadcast:
			log.Println("Sending msg to all clients")
			for client := range p.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					log.Panic(err)
					return
				}
			}
		}
	}
}
