package event

import "github.com/krishnaZawar/LevelCraft/utils/models"

// To create an EventResponse that is sent to the frontend
func NewEmittableResponse(success bool, msg string, data interface{}) *models.EventResponse {
	return &models.EventResponse{
		Success:    success,
		Msg:        msg,
		Data:       data,
		ShouldEmit: true,
	}
}

// To create an EventResponse that is not sent to the frontend
func NewNonEmittableResponse(success bool, msg string, data interface{}) *models.EventResponse {
	return &models.EventResponse{
		Success:    success,
		Msg:        msg,
		Data:       data,
		ShouldEmit: false,
	}
}
