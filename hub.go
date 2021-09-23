// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wschat

import "github.com/mysayasan/arrayhelper"

type Broadcast struct {
	Topic   string
	Message []byte
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	Clients map[*Client]bool

	// Inbound messages from the clients.
	Broadcast chan *Broadcast

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client
}

// NewHub create new hub
func NewHub() *Hub {
	hub := &Hub{
		// Broadcast:  make(chan []byte),
		Broadcast:  make(chan *Broadcast),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}

	go hub.run()

	return hub
}

// Run run hub
func (h *Hub) Run() {
	go h.run()
}

// Run run hub
func (h *Hub) run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		case broadcast := <-h.Broadcast:
			for client := range h.Clients {
				if arrayhelper.StringInSlice(broadcast.Topic, client.Topics) {
					select {
					case client.Send <- broadcast.Message:
					default:
						close(client.Send)
						delete(h.Clients, client)
					}
				}
			}
		}
	}
}
