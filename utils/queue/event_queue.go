package queue

import (
	"errors"
	"sync"

	"github.com/krishnaZawar/LevelCraft/utils/helper"
	"github.com/krishnaZawar/LevelCraft/utils/models"
)

var (
	ErrNoEventsFound = errors.New("error: No Events found")
)

type EventQueue struct {
	eqMu       sync.RWMutex
	eventQueue *helper.Queue[models.Event]
}

func NewEventQueue() *EventQueue {
	return &EventQueue{
		eventQueue: helper.NewQueue[models.Event](),
	}
}

// Ingests an event into the queue for processing
func (eq *EventQueue) Ingest(event models.Event) {
	eq.eqMu.Lock()
	defer eq.eqMu.Unlock()
	eq.eventQueue.Push(event)
}

// ConsumeEvent emits the first Event in the queue for further processing
func (eq *EventQueue) ConsumeEvent() (models.Event, error) {
	eq.eqMu.Lock()
	defer eq.eqMu.Unlock()
	event, ok := eq.eventQueue.Pop()
	if !ok {
		return nil, ErrNoEventsFound
	}
	return event, nil
}

// Length function returns the length of the eventQueue
func (eq *EventQueue) Length() int {
	eq.eqMu.RLock()
	defer eq.eqMu.RUnlock()
	return eq.eventQueue.Length()
}
