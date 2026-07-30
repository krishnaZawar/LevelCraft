package models

// The Event is a well defined breakdown of the Command that the system picks up for updates
type Event interface {
	// GetEventName is used to specify the name of the Event
	//
	// This is the identifier for the Event
	GetEventName() string
}

// EventResponse is used to return back necessary info to the frontend for updates
//
// Not every response is sent to the frontend, only the necessary ones are sent
type EventResponse struct {
	Success    bool        `json:"success"` // whether the event was successful or not
	Msg        string      `json:"msg"`     // the message describing the success or failure of the event
	Data       interface{} `json:"data"`    // the updated data received from the event processing
	ShouldEmit bool        `json:"-"`       // tells whether the data should be sent to the frontend or not
}

// EventHandler is used to handle the corresponding event that occured
type EventHandler interface {
	// Holds the actual logic on what happens when the event occurs
	Handle(Event) *EventResponse
}
