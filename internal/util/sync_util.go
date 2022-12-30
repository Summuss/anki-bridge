package util

import "sync"

type SafeList[T any] struct {
	data []T
	mu   sync.Mutex
}

func (receiver *SafeList[T]) Add(t T) {
	receiver.mu.Lock()
	receiver.data = append(receiver.data, t)
	receiver.mu.Unlock()
}

func (receiver *SafeList[T]) ToSlice() []T {
	return receiver.data
}
