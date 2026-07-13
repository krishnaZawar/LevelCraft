package queue

import (
	"sync"
	"testing"

	"github.com/krishnaZawar/LevelCraft/utils/helper"
	"github.com/stretchr/testify/assert"
)

type MockCommandFactory struct {
	mockNewCommand func([]byte) Command
}

func (mcf *MockCommandFactory) NewCommand(details []byte) Command {
	return mcf.mockNewCommand(details)
}

type MockCommand struct {
	mockGetCommandName func() string
	mockHandle         func() []Event
}

func (mc *MockCommand) GetCommandName() string {
	return mc.mockGetCommandName()
}
func (mc *MockCommand) Handle() []Event {
	return mc.mockHandle()
}

type MockEvent struct {
	mockGetEventName   func() string
	mockGetChannelName func() string
}

func (me *MockEvent) GetEventName() string {
	return me.mockGetEventName()
}
func (me *MockEvent) GetChannelName() string {
	return me.mockGetChannelName()
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
	}

	c1 = &MockCommand{
		mockGetCommandName: func() string {
			return command1
		},
		mockHandle: func() []Event {
			return []Event{e1}
		},
	}

	cf1 = &MockCommandFactory{
		mockNewCommand: func(b []byte) Command {
			return c1
		},
	}
)

func NewTestDecoder() *helper.Registry[string, CommandFactory] {
	decoder := helper.NewRegistry[string, CommandFactory]()
	decoder.Register(command1, cf1)
	return decoder
}

func Test_NewCommandQueue(t *testing.T) {
	cq := NewCommandQueue(NewTestDecoder())

	assert.NotNil(t, cq)
}

func Test_Ingest(t *testing.T) {
	cq := NewCommandQueue(NewTestDecoder())

	assert.Equal(t, 0, cq.Length())

	cr := CommandRequest{
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

		assert.Equal(t, []Event{}, events)
		assert.Equal(t, err, ErrNoCommandRequestsFound)
	})

	t.Run("consume invalid command", func(t *testing.T) {
		invalidReq := CommandRequest{
			RequestType:    invalidCommand,
			RequestDetails: []byte{},
		}
		cq.Ingest(invalidReq)
		events, err := cq.ConsumeCommand()

		assert.Equal(t, []Event{}, events)
		assert.Equal(t, err, ErrFactoryNotFound)
	})

	t.Run("consume valid command", func(t *testing.T) {
		req := CommandRequest{
			RequestType:    command1,
			RequestDetails: []byte{},
		}

		cq.Ingest(req)

		events, err := cq.ConsumeCommand()

		assert.Nil(t, err)

		expectedEvents := []Event{e1}
		assert.Equal(t, expectedEvents, events)
	})
}

func Test_ConcurrentPushPop(t *testing.T) {
	workers := 10
	requests := 100

	invalidReq := CommandRequest{
		RequestType:    invalidCommand,
		RequestDetails: []byte{},
	}
	validReq := CommandRequest{
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
		go func(req CommandRequest) {
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
