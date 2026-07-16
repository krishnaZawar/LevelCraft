package queue

import (
	"encoding/json"
	"sync"
	"testing"

	"github.com/krishnaZawar/LevelCraft/utils/helper"
	"github.com/krishnaZawar/LevelCraft/utils/models"
	"github.com/stretchr/testify/assert"
)

type MockCommandFactory struct {
	mockNewCommand func(json.RawMessage) models.Command
}

func (mcf *MockCommandFactory) NewCommand(details json.RawMessage) models.Command {
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
	mockGetEventName   func() string
	mockGetChannelName func() string
	mockGetSubsribers  func() []string
}

func (me *MockEvent) GetEventName() string {
	return me.mockGetEventName()
}
func (me *MockEvent) GetChannelName() string {
	return me.mockGetChannelName()
}
func (me *MockEvent) GetSubscribers() []string {
	return me.mockGetSubsribers()
}

const (
	event1 = "event1"

	channel1 = "channel1"

	command1       = "command1"
	invalidCommand = "invalid"

	commandFactory1 = "commandFactory1"
)

var (
	e1 = &MockEvent{
		mockGetEventName: func() string {
			return event1
		},
		mockGetChannelName: func() string {
			return channel1
		},
		mockGetSubsribers: func() []string {
			return []string{}
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
		mockNewCommand: func(d json.RawMessage) models.Command {
			return c1
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
