package api

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type ClientList map[*client]bool

// manager handle the websocket connection and clients
type manager struct {
	wsUpgrader websocket.Upgrader
	Clients    ClientList
	Handlers   map[string]eventHandler
	sync.RWMutex
}

// NewManager creates an instance of the Manager struct
func NewManager() *manager {
	wsUpgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	return &manager{
		wsUpgrader: wsUpgrader,
		Clients:    make(ClientList),
		RWMutex:    sync.RWMutex{},
		Handlers:   make(map[string]eventHandler),
	}
}

// Upgrade upgrades the http connection from http rest to Websockets
func (m *manager) Upgrade(w http.ResponseWriter, r *http.Request) {
	// Upgrade regular connection into Web socket
	conn, err := m.wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading the connection:", err)
		return
	}

	// Create the Client
	client := NewClient(conn, m)

	// Add the Client to the Client list
	m.AddClient(client)

	// Start Client process
	go client.ReadMessages()
	go client.WriteMessages()
}

// AddClient adds the Client to the Clients list
func (m *manager) AddClient(client *client) {
	m.Lock()
	defer m.Unlock()

	// Add client to the Client List
	m.Clients[client] = true
}

// RemoveClient removes the Client from the Clients list
func (m *manager) RemoveClient(client *client) {
	m.Lock()
	defer m.Unlock()

	// Check if the client exists
	_, ok := m.Clients[client]
	if !ok {
		log.Println("User does not have a connection with the server")
		return
	}

	// Delete client from the list and close the connection
	delete(m.Clients, client)
	client.conn.Close()
	close(client.buffer) // close buffer channel
}

// RegisterHandler registers the EventHandler on the manager list
func (m *manager) RegisterHandler(eventType string, handler eventHandler) {
	// Check if the event has been already registered
	_, ok := m.Handlers[eventType]
	if ok {
		log.Println("Event type already registered")
		return
	}

	// Register the event handler
	m.Handlers[eventType] = handler
}
