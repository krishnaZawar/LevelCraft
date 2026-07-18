package models

import "encoding/json"

// CommandRequest holds the request data sent by the UI to the backend for updates
type CommandRequest struct {
	RequestType    string          `json:"requestType"`    // The type of request made to determine how it should be handled
	RequestDetails json.RawMessage `json:"requestDetails"` // holds the request specific details
}

// Command is what holds the actual request data and performs processing on it
//
// This is broken down to events accordingly
type Command interface {
	// GetCommandName is used to specify the name of the command
	//
	// This is the identifier for the Command
	GetCommandName() string

	// Handle how the Command should broken down into Events
	Handle() []Event
}

// CommandFactory holds the implementation to convert the CommandRequest details to the Command of our choice
type CommandFactory interface {
	NewCommand(json.RawMessage) (Command, error)
}
