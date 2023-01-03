package common

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

func DoParallel[T any](list *[]T, processor func(*T) error) error {
	size := len(*list)
	var wg sync.WaitGroup
	wg.Add(size)

	errList := SafeList[error]{}
	for i := 0; i < size; i++ {
		i := i
		go func() {
			defer wg.Done()
			err := processor(&(*list)[i])
			if err != nil {
				errList.Add(err)
				return
			}
		}()
	}
	wg.Wait()
	return MergeErrors(errList.ToSlice())
}

func ComputeParallel[T any, R any](list *[]T, processor func(*T) (*R, error)) (*[]*R, *[]error) {
	size := len(*list)
	var wg sync.WaitGroup
	wg.Add(size)

	errList := SafeList[error]{}
	resList := SafeList[*R]{}
	for i := 0; i < size; i++ {
		i := i
		go func() {
			defer wg.Done()
			res, err := processor(&(*list)[i])
			resList.Add(res)
			errList.Add(err)
		}()
	}
	wg.Wait()
	rs := resList.ToSlice()
	errs := errList.ToSlice()
	return &rs, &errs
}
