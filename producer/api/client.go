package api

import (
	"encoding/json"
	"log"
	"time"

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

	// Set how long I will wait for a pong response
	err := c.conn.SetReadDeadline(time.Now().Add(PONG_WAIT))
	if err != nil {
		log.Println("Error setting pong wait time: ", err)
		return
	}

	// Set pong handler
	c.conn.SetPongHandler(c.PongHandler)

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

	ticker := time.NewTicker(PINT_INTERVAL)

	for {
		select {
		// Handle when I receive data from the conn buffer
		case msg, ok := <-c.buffer:
			// If the ok is false means the channel has been closed
			if !ok {
				err := c.conn.WriteMessage(websocket.CloseMessage, nil)
				if err != nil {
					log.Println("Error sending the close conn message:", err)
				}
				return
			}

			// Write the content of the channel
			err := c.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("Error sending message: ", err)
			}

			log.Println("Message sent")

		// Handle when I the ping interval ends
		case <-ticker.C:
			// Send Ping
			err := c.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				log.Println("Error pinging connection: ", err)
				return
			}
		}
	}
}

// PongHandler handles when a Pong message is received, reset the deadline timer
func (c client) PongHandler(pongMsg string) error {
	return c.conn.SetReadDeadline(time.Now().Add(PONG_WAIT))
}
