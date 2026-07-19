package models

// The Event is a well defined breakdown of the Command that the system picks up for updates
type Event interface {
	// GetEventName is used to specify the name of the Event
	//
	// This is the identifier for the Event
	GetEventName() string
}
