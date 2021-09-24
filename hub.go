// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Upgrade by deep

package wschat

import "github.com/mysayasan/arrayhelper"

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
	Broadcast chan *Broadcast

	// Inbound messages from the client to client.
	Private chan *Private

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
		Private:    make(chan *Private),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[string]*Client),
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
		case newclient := <-h.Register:
			if client, ok := h.Clients[newclient.ID]; ok {
				close(client.Send)
				delete(h.Clients, client.ID)
			}
			h.Clients[newclient.ID] = newclient
		case client := <-h.Unregister:
			if _, ok := h.Clients[client.ID]; ok {
				close(client.Send)
				delete(h.Clients, client.ID)
			}
		case broadcast := <-h.Broadcast:
			for _, client := range h.Clients {
				if arrayhelper.StringInSlice(broadcast.Topic, client.Topics) {
					select {
					case client.Send <- broadcast.Message:
					default:
						close(client.Send)
						delete(h.Clients, client.ID)
					}
				}
			}
		case private := <-h.Private:
			if client, ok := h.Clients[private.ID]; ok {
				select {
				case client.Send <- private.Message:
				default:
					close(client.Send)
					delete(h.Clients, client.ID)
				}
			}
		}
	}
}
