package eventmanager

import (
	"testing"

	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/event"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/eventhandler"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/gamestatemanager"
	"github.com/krishnaZawar/LevelCraft/utils/models"
	"github.com/stretchr/testify/assert"
)

type MockEventHandler struct {
	mockHandle func(models.Event, float64) *models.EventResponse
}

func (meh *MockEventHandler) Handle(evt models.Event, dt float64) *models.EventResponse {
	return meh.mockHandle(evt, dt)
}

func Test_registerHandlers(t *testing.T) {
	gsm := &gamestatemanager.GameStateManager{}
	evtManager := NewEventManager()
	registerHandlers(evtManager, gsm)

	assert.Equal(t, 1, evtManager.Size())

	tests := []struct {
		name             string
		eventName        string
		expectedHandlers []models.EventHandler
	}{
		{
			name:      "MoveEvent Handlers",
			eventName: base.Event_MoveEvent,
			expectedHandlers: []models.EventHandler{
				eventhandler.NewMoveEventHandler(gsm),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers, ok := evtManager.GetHandler(tt.eventName)
			assert.Equal(t, true, ok)
			for _, expected := range tt.expectedHandlers {
				assert.Contains(t, handlers, expected)
			}
		})
	}
}

func Test_EventManager(t *testing.T) {
	evtManager := NewEventManager()

	name := "temp"
	handler := &MockEventHandler{
		mockHandle: func(e models.Event, dt float64) *models.EventResponse {
			return event.NewEmittableResponse(true, "dummy resp", nil)
		},
	}

	evtManager.Register(name, handler)
	val, ok := evtManager.GetHandler(name)
	assert.Equal(t, val, []models.EventHandler{handler})
	assert.Equal(t, true, ok)
}

func Test_NewSingleton(t *testing.T) {
	evtManager := newSingleton()

	assert.NotNil(t, evtManager)
}
