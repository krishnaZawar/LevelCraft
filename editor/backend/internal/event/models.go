package event

import "github.com/krishnaZawar/LevelCraft/utils/models"

// EventResponse is used to return back necessary info to the frontend for updates
//
// Not every response is sent to the frontend, only the necessary ones are sent
type EventResponse struct {
	Success    bool        `json:"success"` // whether the event was successful or not
	Msg        string      `json:"msg"`     // the message describing the success or failure of the event
	Data       interface{} `json:"data"`    // the updated data received from the event processing
	ShouldEmit bool        `json:"-"`       // tells whether the data should be sent to the frontend or not
}

// To create an EventResponse that is sent to the frontend
func NewEmittableResponse(success bool, msg string, data interface{}) *EventResponse {
	return &EventResponse{
		Success:    success,
		Msg:        msg,
		Data:       data,
		ShouldEmit: true,
	}
}

// To create an EventResponse that is not sent to the frontend
func NewNonEmittableResponse(success bool, msg string, data interface{}) *EventResponse {
	return &EventResponse{
		Success:    success,
		Msg:        msg,
		Data:       data,
		ShouldEmit: false,
	}
}

// EventHandler is used to handle the corresponding event that occured
type EventHandler interface {
	// Holds the actual logic on what happens when the event occurs
	Handle(models.Event) *EventResponse
}
