package eventmanager

import (
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/eventhandler"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/gamestatemanager"
	"github.com/krishnaZawar/LevelCraft/utils/models"
)

// EventManager routes the event to their respective handlers
type EventManager struct {
	// maps the event name -> []handler
	handlers map[string][]models.EventHandler
}

// Creates a new EventManager with all the events and their corresponding handlers registered
func NewEventManager() *EventManager {
	evtManager := &EventManager{
		handlers: map[string][]models.EventHandler{},
	}

	return evtManager
}

func registerHandlers(
	evtManager *EventManager,
	gsm *gamestatemanager.GameStateManager,
) {
	evtManager.Register(base.Event_MoveEvent, eventhandler.NewMoveEventHandler(gsm))
}

// registers a new handler for a particular event
func (em *EventManager) Register(eventName string, handler models.EventHandler) {
	em.handlers[eventName] = append(em.handlers[eventName], handler)
}

// fetches the corresponding handler for that event
func (em *EventManager) GetHandler(eventName string) ([]models.EventHandler, bool) {
	handlers, ok := em.handlers[eventName]
	return handlers, ok
}

func (em *EventManager) Size() int {
	return len(em.handlers)
}

func newSingleton() *EventManager {
	evtManager := NewEventManager()
	registerHandlers(evtManager, gamestatemanager.Get())
	return evtManager
}

var evtManager = newSingleton()

func Get() *EventManager {
	return evtManager
}
