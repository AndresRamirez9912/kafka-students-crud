package api

import (
	"encoding/json"
)

// event represents the msg received by the user, is a wrapper of the handler
type event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// eventHandler represents a data type for the handlers per event
type eventHandler func(event event, c *client) error

// Events registered
const (
	WS_EVENT_REGISTER_STUDENTS = "register_students"
	WS_EVENT_GET_STUDENT       = "get_student"
	WS_EVENT_DELETE_STUDENT    = "delete_students"
	WS_EVENT_UPDATE_STUDENT    = "update_students"
)

// Handlers
func HandleRegisterStudents(event event, c *client) error {
	c.buffer <- []byte("Hello")
	return nil
}
