package eventmanager

import (
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/event"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/gamestatemanager"
	"github.com/krishnaZawar/LevelCraft/utils/helper"
)

// EventManager routes the event to their respective handlers
type EventManager struct {
	// maps the event name -> handler
	handlers *helper.Registry[string, event.EventHandler]
}

// Creates a new EventManager with all the events and their corresponding handlers registered
func NewEventManager() *EventManager {
	evtManager := &EventManager{
		handlers: helper.NewRegistry[string, event.EventHandler](),
	}

	return evtManager
}

func registerHandlers(evtManager *EventManager) {
	gsm := gamestatemanager.Get()
	evtManager.Register(base.Event_AddGameobject, event.NewAddGameobjectEventHandler(gsm))
	evtManager.Register(base.Event_DeleteGameobject, event.NewDeleteGameobjectEventHandler(gsm))
	evtManager.Register(base.Event_AttachComponent, event.NewAttachComponentEventHandler(gsm))
	evtManager.Register(base.Event_DetachComponent, event.NewDetachComponentEventHandler(gsm))
	evtManager.Register(base.Event_UpdateComponent, event.NewUpdateComponentEventHandler(gsm))
}

// registers a new handler for a particular event
func (em *EventManager) Register(eventName string, handler event.EventHandler) {
	em.handlers.Register(eventName, handler)
}

// fetches the corresponding handler for that event
func (em *EventManager) GetHandler(eventName string) (event.EventHandler, bool) {
	return em.handlers.GetValue(eventName)
}

func newSingleton() *EventManager {
	evtManager := NewEventManager()
	registerHandlers(evtManager)
	return evtManager
}

var evtManager = newSingleton()

func Get() *EventManager {
	return evtManager
}
