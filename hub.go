// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Upgrade by deep

package wschat

import (
	"fmt"
)

type Broadcast struct {
	Topic   string
	Message []byte
}

type Private struct {
	ID      string
	Message []byte
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	Clients map[string]*Client

	// Inbound messages from the clients to topic.
	Broadcast chan []byte

	// Inbound messages from the client to client.
	Private chan *Private

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client

	// Kill client
	kill chan string
}

// NewHub create new hub
func NewHub() *Hub {
	hub := &Hub{
		// Broadcast:  make(chan []byte),
		Clients:    make(map[string]*Client),
		Broadcast:  make(chan []byte),
		Private:    make(chan *Private),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		kill:       make(chan string),
	}

	go hub.run()

	return hub
}

// Run run hub
func (h *Hub) Kill(id string) {
	h.kill <- id
}

// Run run hub
func (h *Hub) Run() {
	go h.run()
}

// Run run hub
func (h *Hub) run() {
	for {
		select {
		case newclient := <-h.Register:
			fmt.Printf("%s\n", newclient.id)
			if client, ok := h.Clients[newclient.id]; ok {
				client.quit <- true
				delete(h.Clients, client.id)
			}
			h.Clients[newclient.id] = newclient
			fmt.Printf("Total connected: %d\n", len(h.Clients))
		case client := <-h.Unregister:
			fmt.Printf("Unregister: %s\n", client.id)
			delete(h.Clients, client.id)
		case id := <-h.kill:
			if client, ok := h.Clients[id]; ok {
				client.quit <- true
				delete(h.Clients, client.id)
			}

			// case broadcast := <-h.Broadcast:
			// 	for _, client := range h.Clients {
			// 		if arrayhelper.StringInSlice(broadcast.Topic, client.topics) {
			// 			select {
			// 			case client.send <- broadcast.Message:
			// 			default:
			// 				close(client.send)
			// 				delete(h.Clients, client.id)
			// 			}
			// 		}
			// 	}
		case message := <-h.Broadcast:
			for _, client := range h.Clients {
				select {
				case client.send <- message:
				default:
					client.quit <- true
					delete(h.Clients, client.id)
				}
			}

		case private := <-h.Private:
			fmt.Printf("Searching...")
			if client, ok := h.Clients[private.ID]; ok {
				fmt.Printf("Found : %s\n", private.Message)
				select {
				case client.send <- private.Message:
				default:
					client.quit <- true
					delete(h.Clients, client.id)
				}
			}
		}
	}
}
