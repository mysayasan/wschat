// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wschat

import (
	"fmt"

	"github.com/gorilla/websocket"
)

// const (
// 	// Time allowed to write a message to the peer.
// 	writeWait = 10 * time.Second

// 	// Time allowed to read the next pong message from the peer.
// 	pongWait = 10 * time.Second

// 	// Send pings to peer with this period. Must be less than pongWait.
// 	pingPeriod = (pongWait * 9) / 10

// 	// Maximum message size allowed from peer.
// 	maxMessageSize = 512
// )

// var (
// 	newline = []byte{'\n'}
// 	space   = []byte{' '}
// )

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	id string

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	quit chan bool
}

// NewBrokerUcase will create new an brokerUcase object representation of broker.Usecase interface
func NewClient(hub *Hub, id string, conn *websocket.Conn, send chan []byte) *Client {
	return &Client{
		hub:  hub,
		id:   id,
		conn: conn,
		send: send,
		quit: make(chan bool),
	}
}

func (c *Client) Run() {
	// c.hub.Register <- c
	go c.runReader()
	c.runWriter()

}

func (c *Client) Quit() {
	c.quit <- true
	c.conn.Close()
	// c.hub.Unregister <- c
}

func (c *Client) GetID() string {
	return c.id
}

func (c *Client) Send(data []byte) {
	c.send <- data
}

func (c *Client) runWriter() {
	for {
		select {
		case data := <-c.send:
			fmt.Printf("Received: %s\n", data)
			err := c.conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Printf("Writer: %s\n", err)
				c.Quit()
				// c.Logger().Error(err)
			}
		case <-c.quit:
			fmt.Printf("writer stop\n")
			return
		default:
			continue
		}
	}
}

func (c *Client) runReader() {
	for {
		select {
		case <-c.quit:
			fmt.Printf("reader stop\n")
			return
		default:
			// Read
			_, msg, err := c.conn.ReadMessage()
			if err != nil {
				fmt.Printf("Reader: %s\n", err)
				// if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// 	fmt.Printf("error: %v", err)
				// }
				c.Quit()
				return
			}

			c.send <- msg
			fmt.Printf("%s >> %s\n", c.id, msg)
		}
	}
}
