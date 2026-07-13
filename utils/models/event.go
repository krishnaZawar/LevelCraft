package models

// The Event is a well defined breakdown of the Command that the system picks up for updates
type Event interface {
	// GetEventName is used to specify the name of the Event
	//
	// This is the identifier for the Event
	GetEventName() string

	// Returns the name of the eventStream where the event should be published
	GetChannelName() string

	// Returns the name of all the components which should be triggered for this event
	GetSubscribers() []string
}
