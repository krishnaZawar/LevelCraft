package eventmanager

import (
	"testing"

	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/event"
	"github.com/krishnaZawar/LevelCraft/utils/models"
	"github.com/stretchr/testify/assert"
)

type MockEventHandler struct {
	mockHandle func(models.Event) *models.EventResponse
}

func (meh *MockEventHandler) Handle(evt models.Event) *models.EventResponse {
	return meh.mockHandle(evt)
}

func Test_EventManager(t *testing.T) {
	evtManager := NewEventManager()

	name := "temp"
	handler := &MockEventHandler{
		mockHandle: func(e models.Event) *models.EventResponse {
			return event.NewEmittableResponse(true, "dummy resp", nil)
		},
	}

	evtManager.Register(name, handler)
	val, ok := evtManager.GetHandler(name)
	assert.Equal(t, val, handler)
	assert.Equal(t, true, ok)
}

func Test_NewSingleton(t *testing.T) {
	evtManager := newSingleton()

	assert.NotNil(t, evtManager)
}
