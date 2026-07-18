package queue

import (
	"encoding/json"
	"errors"
	"sync"
	"testing"

	"github.com/krishnaZawar/LevelCraft/utils/helper"
	"github.com/krishnaZawar/LevelCraft/utils/models"
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

const (
	event1 = "event1"

	command1       = "command1"
	invalidCommand = "invalid"

	factoryErr = "error"
)

var (
	e1 = &MockEvent{
		mockGetEventName: func() string {
			return event1
		},
	}

	c1 = &MockCommand{
		mockGetCommandName: func() string {
			return command1
		},
		mockHandle: func() []models.Event {
			return []models.Event{e1}
		},
	}

	cf1 = &MockCommandFactory{
		mockNewCommand: func(d json.RawMessage) (models.Command, error) {
			return c1, nil
		},
	}
	cf2 = &MockCommandFactory{
		mockNewCommand: func(d json.RawMessage) (models.Command, error) {
			return nil, errors.New(factoryErr)
		},
	}
)

func NewTestDecoder() *helper.Registry[string, models.CommandFactory] {
	decoder := helper.NewRegistry[string, models.CommandFactory]()
	decoder.Register(command1, cf1)
	return decoder
}

func Test_NewCommandQueue(t *testing.T) {
	cq := NewCommandQueue(NewTestDecoder())

	assert.NotNil(t, cq)
}

func Test_IngestCommand(t *testing.T) {
	cq := NewCommandQueue(NewTestDecoder())

	assert.Equal(t, 0, cq.Length())

	cr := models.CommandRequest{
		RequestType:    command1,
		RequestDetails: []byte{},
	}

	cq.Ingest(cr)

	assert.Equal(t, 1, cq.Length())
}

func Test_ConsumeCommand(t *testing.T) {
	cq := NewCommandQueue(NewTestDecoder())

	t.Run("consume from empty queue", func(t *testing.T) {
		events, err := cq.ConsumeCommand()

		assert.Equal(t, []models.Event{}, events)
		assert.Equal(t, err, ErrNoCommandRequestsFound)
	})

	t.Run("consume invalid command", func(t *testing.T) {
		invalidReq := models.CommandRequest{
			RequestType:    invalidCommand,
			RequestDetails: []byte{},
		}
		cq.Ingest(invalidReq)
		events, err := cq.ConsumeCommand()

		assert.Equal(t, []models.Event{}, events)
		assert.Equal(t, err, ErrFactoryNotFound)
	})

	t.Run("consume valid command", func(t *testing.T) {
		req := models.CommandRequest{
			RequestType:    command1,
			RequestDetails: []byte{},
		}

		cq.Ingest(req)

		events, err := cq.ConsumeCommand()

		assert.Nil(t, err)

		expectedEvents := []models.Event{e1}
		assert.Equal(t, expectedEvents, events)
	})

	t.Run("test error from factory", func(t *testing.T) {
		decoder := helper.NewRegistry[string, models.CommandFactory]()
		decoder.Register(command1, cf2)
		cq := NewCommandQueue(decoder)

		req := models.CommandRequest{
			RequestType:    command1,
			RequestDetails: []byte{},
		}
		cq.Ingest(req)

		_, err := cq.ConsumeCommand()

		assert.NotNil(t, err)
	})
}

func Test_ConcurrentPushPop(t *testing.T) {
	workers := 10
	requests := 100

	invalidReq := models.CommandRequest{
		RequestType:    invalidCommand,
		RequestDetails: []byte{},
	}
	validReq := models.CommandRequest{
		RequestType:    command1,
		RequestDetails: []byte{},
	}

	cq := NewCommandQueue(NewTestDecoder())

	var wg sync.WaitGroup

	for i := 1; i <= workers; i++ {
		wg.Add(1)
		req := invalidReq
		if i%2 == 0 {
			req = validReq
		}
		go func(req models.CommandRequest) {
			defer wg.Done()
			for i := 1; i <= requests; i++ {
				cq.Ingest(req)
			}
		}(req)
	}

	wg.Wait()

	expectedLen := workers * requests

	assert.Equal(t, expectedLen, cq.Length())

	var resMu sync.RWMutex
	var res struct {
		Correct   int
		Incorrect int
	}

	for i := 1; i <= workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 1; i <= requests; i++ {
				_, err := cq.ConsumeCommand()
				resMu.Lock()
				if err != nil {
					res.Incorrect++
				} else {
					res.Correct++
				}
				resMu.Unlock()
			}
		}()
	}

	wg.Wait()

	expectedCorrect := workers / 2 * requests
	expectedIncorrect := expectedLen - expectedCorrect

	assert.Equal(t, 0, cq.Length())
	assert.Equal(t, expectedCorrect, res.Correct)
	assert.Equal(t, expectedIncorrect, res.Incorrect)
}
