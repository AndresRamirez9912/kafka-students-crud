package api

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type client struct {
	conn    *websocket.Conn
	manager *manager
	buffer  chan []byte // This is a unBuffer channel where the received msg will be stored there
}

func NewClient(conn *websocket.Conn, manamanager *manager) *client {
	return &client{
		conn:    conn,
		manager: manamanager,
		buffer:  make(chan []byte),
	}
}

func (c *client) ReadMessages() {
	// if the connection ends remove the connection
	defer func() {
		c.manager.RemoveClient(c)
	}()

	for {
		// Check if the connection has an error
		_, payload, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message: %v", err)
			}
			break
		}

		// Parse events
		evt := event{}
		err = json.Unmarshal(payload, &evt)
		if err != nil {
			log.Println("Invalid event format:", err)
			continue
		}

		// Search the registered event handler
		handler, ok := c.manager.Handlers[evt.Type]
		if !ok {
			log.Printf("No handler for event type: %s\n", evt.Type)
			continue
		}

		// Execute the handler
		err = handler(evt, c)
		if err != nil {
			log.Printf("Handler error for type %s: %v\n", evt.Type, err)
		}

	}
}

func (c *client) WriteMessages() {
	defer func() {
		c.manager.RemoveClient(c)
	}()

	for msg := range c.buffer {
		err := c.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("Error sending message: ", err)
		}

		log.Println("Message sent")
	}

	// If the look ends means the channel has been closed
	err := c.conn.WriteMessage(websocket.CloseMessage, nil)
	if err != nil {
		log.Println("Error sending the close conn message:", err)
	}
}
