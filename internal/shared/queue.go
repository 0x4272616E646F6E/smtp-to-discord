package shared

import "sync"

type MessageQueue[T any] struct {
	mu    sync.Mutex
	queue []T
}

func NewMessageQueue[T any]() *MessageQueue[T] {
	return &MessageQueue[T]{queue: make([]T, 0)}
}

func (q *MessageQueue[T]) Enqueue(msg T) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.queue = append(q.queue, msg)
}

func (q *MessageQueue[T]) Dequeue() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.queue) == 0 {
		var zero T
		return zero, false
	}
	msg := q.queue[0]
	q.queue = q.queue[1:]
	return msg, true
}

func (q *MessageQueue[T]) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.queue)
}
