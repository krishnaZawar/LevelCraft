package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
The tests below carry coverage for functions other than:
- Push()
- Pop()
- Peek()
Hence, explicit tests are not written for those
*/

func Test_NewQueue(t *testing.T) {
	queue := NewQueue[string]()
	assert.NotNil(t, queue)
	assert.Equal(t, 0, queue.Length())
	assert.Equal(t, initialQueueSize, queue.arrayLen())
}

func Test_Push(t *testing.T) {
	t.Run("Push elements from start", func(t *testing.T) {
		queue := NewQueue[int]()
		for i := 0; i < initialQueueSize; i++ {
			queue.Push(i)
		}

		assert.Equal(t, initialQueueSize, queue.Length())
		assert.Equal(t, initialQueueSize, queue.arrayLen())

		expected := 0
		for !queue.IsEmpty() {
			val, ok := queue.Pop()
			assert.Equal(t, true, ok)
			assert.Equal(t, expected, val)
			expected++
		}

		assert.Equal(t, 0, queue.Length())
	})

	t.Run("Push elements from the middle", func(t *testing.T) {
		queue := NewQueue[int]()

		queue.Push(1)
		queue.Push(2)
		queue.Pop()
		queue.Pop()

		for i := 0; i < initialQueueSize; i++ {
			queue.Push(i)
		}

		assert.Equal(t, initialQueueSize, queue.Length())
		assert.Equal(t, initialQueueSize, queue.arrayLen())

		expected := 0
		for !queue.IsEmpty() {
			val, ok := queue.Pop()
			assert.Equal(t, true, ok)
			assert.Equal(t, expected, val)
			expected++
		}

		assert.Equal(t, initialQueueSize, queue.arrayLen())
		assert.Equal(t, 0, queue.Length())
	})

	t.Run("grow queue", func(t *testing.T) {
		queue := NewQueue[int]()
		for i := 0; i < 2*initialQueueSize; i++ {
			queue.Push(i)
		}

		assert.Equal(t, 2*initialQueueSize, queue.Length())
		assert.Equal(t, 2*initialQueueSize, queue.arrayLen())
	})
}

func Test_Pop(t *testing.T) {
	t.Run("pop elements from start", func(t *testing.T) {
		queue := NewQueue[int]()
		for i := 0; i < initialQueueSize; i++ {
			queue.Push(i)
		}

		expected := 0
		for !queue.IsEmpty() {
			val, ok := queue.Pop()
			assert.Equal(t, true, ok)
			assert.Equal(t, expected, val)
			expected++
		}
		assert.Equal(t, 0, queue.Length())
	})

	t.Run("pop elements from the middle", func(t *testing.T) {
		queue := NewQueue[int]()

		queue.Push(1)
		queue.Push(2)
		queue.Pop()
		queue.Pop()

		for i := 0; i < initialQueueSize; i++ {
			queue.Push(i)
		}

		expected := 0
		for !queue.IsEmpty() {
			val, ok := queue.Pop()
			assert.Equal(t, true, ok)
			assert.Equal(t, true, ok)
			assert.Equal(t, expected, val)
			expected++
		}
		assert.Equal(t, 0, queue.Length())
	})

	t.Run("pop from empty queue", func(t *testing.T) {
		queue := NewQueue[string]()
		_, ok := queue.Pop()
		assert.Equal(t, false, ok)
	})
}

func Test_Peek(t *testing.T) {
	queue := NewQueue[int]()
	value := 1
	t.Run("when value does not exists", func(t *testing.T) {
		_, ok := queue.Peek()
		assert.Equal(t, false, ok)
	})

	queue.Push(value)

	t.Run("when value exists", func(t *testing.T) {
		val, ok := queue.Peek()
		assert.Equal(t, true, ok)
		assert.Equal(t, value, val)
	})
}
