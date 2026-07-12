package helper

const (
	// defines the initial size with which the queue is initialized
	initialQueueSize = 8
)

// This is a generic thread-safe implementation of the Queue
//
// This implementation follows the ring buffer for more memory efficient operations
type Queue[T any] struct {
	items []T
	head  int
	tail  int
	size  int
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		items: make([]T, initialQueueSize),
		head:  0,
		tail:  -1,
		size:  0,
	}
}

// pushes an element to the end of the queue
func (q *Queue[T]) Push(val T) {
	if q.size == len(q.items) {
		// grow the queue size when the no space to accomodate new elements
		q.grow()
	}
	q.tail = (q.tail + 1) % len(q.items)
	q.items[q.tail] = val
	q.size++
}

// pops the first element in the queue
func (q *Queue[T]) Pop() (T, bool) {
	var zero T
	if q.size == 0 {
		return zero, false
	}
	val := q.items[q.head]
	q.items[q.head] = zero
	q.head = (q.head + 1) % len(q.items)
	q.size--
	return val, true
}

// returns the value present at the head of the queue
//
// Return types:
//   - val T: returns the value present at the head, default if not present
//   - bool: returns True if the item exists else returns False
func (q *Queue[T]) Peek() (T, bool) {
	var zero T
	if q.size == 0 {
		return zero, false
	}
	return q.items[q.head], true
}

// returns true if the queue is empty
func (q *Queue[T]) IsEmpty() bool {
	return q.size == 0
}

// returns the length of the queue
func (q *Queue[T]) Length() int {
	return q.size
}

// returns the length of items array in the queue
func (q *Queue[T]) arrayLen() int {
	return len(q.items)
}

// doubles the queue size
func (q *Queue[T]) grow() {
	new_items := make([]T, 2*len(q.items))
	for i := 0; i < q.size; i++ {
		new_items[i] = q.items[(q.head+i)%len(q.items)]
	}
	q.items = new_items
	q.head = 0
	q.tail = q.size - 1
}
