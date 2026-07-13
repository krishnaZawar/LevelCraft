package queue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	subscriber1 = "sub1"
)

func Test_NewEventQueue(t *testing.T) {
	eq := NewEventQueue()

	assert.NotNil(t, eq)
}

func Test_IngestEvent(t *testing.T) {
	eq := NewEventQueue()

	eq.Ingest(e1)

	assert.Equal(t, 1, eq.Length())
}

func Test_ConsumeEvent(t *testing.T) {
	eq := NewEventQueue()

	t.Run("consume from empty event queue", func(t *testing.T) {
		_, err := eq.ConsumeEvent()
		assert.Equal(t, ErrNoEventsFound, err)
	})

	t.Run("consume when events exist in queue", func(t *testing.T) {
		eq.Ingest(e1)
		event, err := eq.ConsumeEvent()

		assert.Nil(t, err)
		assert.Equal(t, event, e1)
		assert.Equal(t, 0, eq.Length())
	})
}
