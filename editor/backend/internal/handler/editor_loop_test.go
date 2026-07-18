package handler

import (
	"encoding/json"
	"testing"

	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/event"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/eventmanager"
	"github.com/krishnaZawar/LevelCraft/utils/helper"
	"github.com/krishnaZawar/LevelCraft/utils/models"
	"github.com/krishnaZawar/LevelCraft/utils/queue"
	"github.com/stretchr/testify/assert"
)

type MockCommandFactory struct {
	mockNewCommand func(json.RawMessage) (models.Command, error)
}

func (mcf *MockCommandFactory) NewCommand(details json.RawMessage) (models.Command, error) {
	return mcf.mockNewCommand(details)
}

type MockCommand struct {
	mockGetCommandName func() string
	mockHandle         func() []models.Event
}

func (mc *MockCommand) GetCommandName() string {
	return mc.mockGetCommandName()
}
func (mc *MockCommand) Handle() []models.Event {
	return mc.mockHandle()
}

type MockEvent struct {
	mockGetEventName func() string
}

func (me *MockEvent) GetEventName() string {
	return me.mockGetEventName()
}

type MockEventHandler struct {
	mockHandle func(models.Event) *event.EventResponse
}

func (meh *MockEventHandler) Handle(evt models.Event) *event.EventResponse {
	return meh.mockHandle(evt)
}

const (
	commandName             = "command"
	invalidCommandName      = "invalidCommand"
	invalidEventCommandName = "invalidEventCommand"

	eventName        = "event"
	invalidEventName = "invalidEvent"
)

var (
	cmd = &MockCommand{
		mockGetCommandName: func() string {
			return commandName
		},
		mockHandle: func() []models.Event {
			return []models.Event{evt}
		},
	}

	cmdWithInvalidEvent = &MockCommand{
		mockGetCommandName: func() string {
			return invalidEventCommandName
		},
		mockHandle: func() []models.Event {
			return []models.Event{invalidEvt}
		},
	}

	evt = &MockEvent{
		mockGetEventName: func() string {
			return eventName
		},
	}

	invalidEvt = &MockEvent{
		mockGetEventName: func() string {
			return invalidEventName
		},
	}

	cmdFactory = &MockCommandFactory{
		mockNewCommand: func(rm json.RawMessage) (models.Command, error) {
			return cmd, nil
		},
	}

	invalidEventCmdFactory = &MockCommandFactory{
		mockNewCommand: func(rm json.RawMessage) (models.Command, error) {
			return cmdWithInvalidEvent, nil
		},
	}

	evtHandler = &MockEventHandler{
		mockHandle: func(e models.Event) *event.EventResponse {
			return event.NewEmittableResponse(true, "success", nil)
		},
	}
)

func MockLoopData() (*queue.CommandQueue, *queue.EventQueue, *helper.Queue[event.EventResponse], *eventmanager.EventManager) {
	decoder := helper.NewRegistry[string, models.CommandFactory]()
	decoder.Register(commandName, cmdFactory)
	decoder.Register(invalidEventCommandName, invalidEventCmdFactory)

	cmdQueue := queue.NewCommandQueue(decoder)
	evtQueue := queue.NewEventQueue()
	respQueue := helper.NewQueue[event.EventResponse]()

	evtManager := eventmanager.NewEventManager()
	evtManager.Register(eventName, evtHandler)

	return cmdQueue, evtQueue, respQueue, evtManager
}

func Test_UpdateEditor(t *testing.T) {
	cmdQueue, evtQueue, respQueue, evtManager := MockLoopData()

	t.Run("all valid commands and events", func(t *testing.T) {
		cmdRequest := models.CommandRequest{
			RequestType:    commandName,
			RequestDetails: []byte{},
		}
		cmdQueue.Ingest(cmdRequest)
		updateEditor(cmdQueue, evtQueue, respQueue, evtManager)

		assert.Equal(t, 1, respQueue.Length())

		resp, ok := respQueue.Pop()

		assert.Equal(t, true, ok)
		assert.Equal(t, true, resp.Success)
		assert.Equal(t, "success", resp.Msg)
		assert.Nil(t, resp.Data)
	})

	t.Run("invalid command", func(t *testing.T) {
		cmdRequest := models.CommandRequest{
			RequestType:    invalidCommandName,
			RequestDetails: []byte{},
		}
		cmdQueue.Ingest(cmdRequest)
		updateEditor(cmdQueue, evtQueue, respQueue, evtManager)

		assert.Equal(t, true, respQueue.IsEmpty())
	})

	t.Run("invalid event", func(t *testing.T) {
		cmdRequest := models.CommandRequest{
			RequestType:    invalidEventCommandName,
			RequestDetails: []byte{},
		}
		cmdQueue.Ingest(cmdRequest)
		updateEditor(cmdQueue, evtQueue, respQueue, evtManager)

		assert.Equal(t, true, respQueue.IsEmpty())
	})

	t.Run("no events to process", func(t *testing.T) {
		updateEditor(cmdQueue, evtQueue, respQueue, evtManager)

		assert.Equal(t, true, respQueue.IsEmpty())
	})
}
