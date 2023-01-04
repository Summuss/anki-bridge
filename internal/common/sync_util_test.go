package common

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestDoParallel(t *testing.T) {
	err := DoParallel(
		&[]int{1, 2, 3, 4, 5}, func(i *int) error {
			s := "hello"
			defer func() {
				println("s:" + s)
			}()
			s = "world"
			fmt.Printf("%d started\n", *i)
			ii := *i
			time.Sleep(time.Duration(ii * int(time.Second)))
			if ii%2 == 0 {
				return errors.New(strconv.Itoa(ii))
			} else {
				return nil
			}
		},
	)
	println(err.Error())
}

func TestChannel(t *testing.T) {
	ch := make(chan int, 5)
	go func() {
		ch <- 1
		<-time.After(time.Second)
		ch <- 2
		<-time.After(time.Second)
		ch <- 3
		<-time.After(2 * time.Second)
		close(ch)
	}()
	for {
		i, ok := <-ch
		if ok {
			fmt.Printf("get %d\n", i)
		} else {
			fmt.Printf("finished\n")
			return
		}
	}
}
