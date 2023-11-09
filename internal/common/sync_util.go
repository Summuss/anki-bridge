package common

import (
	"sync"
)

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

func DoParallel[T any](list *[]T, processor func(*T) error) error {
	return DoParallelWithLimitThread(list, processor, len(*list))
}
func DoParallelWithLimitThread[T any](list *[]T, processor func(*T) error, threadSize int) error {
	size := len(*list)
	errList := make([]error, 0, size)
	ch := make(chan error, size)
	ctrlCh := make(chan struct{}, threadSize)
	for i := 0; i < size; i++ {
		i := i
		ctrlCh <- struct{}{}
		go func() {
			err := processor(&(*list)[i])
			ch <- err
			<-ctrlCh
		}()
	}
	for i := 0; i < size; i++ {
		err := <-ch
		if err != nil {
			errList = append(errList, err)
		}
	}
	return MergeErrors(errList)
}

func ComputeParallel[T any, R any](list *[]T, processor func(*T) (*R, error)) (*[]*R, *[]error) {
	size := len(*list)
	resCh := make(chan [2]interface{}, size)
	for i := 0; i < size; i++ {
		i := i
		go func() {
			res, err := processor(&(*list)[i])
			resCh <- [2]interface{}{res, err}
		}()
	}
	errList := make([]error, 0, size)
	resList := make([]*R, 0, size)
	for i := 0; i < size; i++ {
		res := <-resCh
		resList = append(resList, res[0].(*R))
		errList = append(errList, res[1].(error))
	}
	return &resList, &errList
}
