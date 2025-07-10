package main

import (
	"kafka-producer/api"
	"log"
	"net/http"
)

func main() {
	// Create websocket manager
	manager := api.NewManager()

	// Create Event handlers
	manager.RegisterHandler(api.WS_EVENT_REGISTER_STUDENTS, api.HandleRegisterStudents)

	// List websocket handler
	http.HandleFunc("/ws", manager.Upgrade)

	// Start WebSocket server
	log.Println("Starting server at port :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
