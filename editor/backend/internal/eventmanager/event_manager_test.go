package eventmanager

import (
	"testing"

	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/event"
	"github.com/krishnaZawar/LevelCraft/utils/models"
	"github.com/stretchr/testify/assert"
)

type MockEventHandler struct {
	mockHandle func(models.Event) *event.EventResponse
}

func (meh *MockEventHandler) Handle(evt models.Event) *event.EventResponse {
	return meh.mockHandle(evt)
}

func Test_EventManager(t *testing.T) {
	evtManager := NewEventManager()

	name := "temp"
	handler := &MockEventHandler{
		mockHandle: func(e models.Event) *event.EventResponse {
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

	var handler event.EventHandler
	var ok bool

	handler, ok = evtManager.GetHandler(base.Event_AddGameobject)
	assert.Equal(t, true, ok)
	assert.IsType(t, &event.AddGameobjectEventHandler{}, handler)

	handler, ok = evtManager.GetHandler(base.Event_DeleteGameobject)
	assert.Equal(t, true, ok)
	assert.IsType(t, &event.DeleteGameobjectEventHandler{}, handler)

	handler, ok = evtManager.GetHandler(base.Event_AttachComponent)
	assert.Equal(t, true, ok)
	assert.IsType(t, &event.AttachComponentEventHandler{}, handler)

	handler, ok = evtManager.GetHandler(base.Event_DetachComponent)
	assert.Equal(t, true, ok)
	assert.IsType(t, &event.DetachComponentEventHandler{}, handler)
}
