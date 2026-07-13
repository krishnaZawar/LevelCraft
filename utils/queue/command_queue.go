package queue

import (
	"errors"
	"sync"

	"github.com/krishnaZawar/LevelCraft/utils/helper"
)

var (
	ErrNoCommandRequestsFound = errors.New("error: No CommandRequests found")
	ErrFactoryNotFound        = errors.New("error: CommandFactory not found")
)

// All the incoming CommandRequests pass through this
//
// It is the communication channel for Frontend -> Backend
//
// Note: Every valid CommandRequest should have a Factory registered with the decoder. If not then it is treated as an invalid CommandRequest
type CommandQueue struct {
	rqMu         sync.RWMutex
	requestQueue *helper.Queue[CommandRequest] // Holds all the incoming CommandRequests sequentially

	// decoder is used to convert the CommandRequest -> Command
	decoder *helper.Registry[string, CommandFactory]
}

func NewCommandQueue(decoder *helper.Registry[string, CommandFactory]) *CommandQueue {
	return &CommandQueue{
		requestQueue: helper.NewQueue[CommandRequest](),
		decoder:      decoder,
	}
}

// Ingests a CommandRequest into the queue for processing
func (cq *CommandQueue) Ingest(req CommandRequest) {
	cq.rqMu.Lock()
	defer cq.rqMu.Unlock()
	cq.requestQueue.Push(req)
}

// ConsumeCommand processes the first CommandRequest in the queue and converts it into their corresponding events
//
// Return types:
//   - []Event: all the events that the command should emit
//   - error : returns error when unable to fetch CommandRequest or the corresponding Factory
func (cq *CommandQueue) ConsumeCommand() ([]Event, error) {
	cq.rqMu.Lock()
	req, ok := cq.requestQueue.Pop()
	if !ok {
		cq.rqMu.Unlock()
		return []Event{}, ErrNoCommandRequestsFound
	}
	cq.rqMu.Unlock()

	// fetch the factory for the corresponding CommandRequest
	factory, ok := cq.decoder.GetValue(req.RequestType)
	if !ok {
		return []Event{}, ErrFactoryNotFound
	}

	// fetch Command with all the details
	command := factory.NewCommand(req.RequestDetails)

	// fetch all the corresponding events to be published from the Command
	return command.Handle(), nil
}

// Length function returns the length of the commandQueue
func (cq *CommandQueue) Length() int {
	return cq.requestQueue.Length()
}
