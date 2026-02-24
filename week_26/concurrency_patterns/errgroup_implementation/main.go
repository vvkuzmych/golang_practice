package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type ErrGroup struct {
	err  error
	wg   sync.WaitGroup
	once sync.Once

	doneCh chan struct{}
}

func NewErrGroup() (*ErrGroup, chan struct{}) {
	doneCh := make(chan struct{})
	return &ErrGroup{
		doneCh: doneCh,
	}, doneCh
}

func (eg *ErrGroup) Go(task func() error) {
	eg.wg.Add(1)
	go func() {
		defer eg.wg.Done()

		select {
		case <-eg.doneCh:
			return
		default:
			if err := task(); err != nil {
				eg.once.Do(func() {
					eg.err = err
					close(eg.doneCh)
				})
			}
		}
	}()
}

func (eg *ErrGroup) Wait() error {
	eg.wg.Wait()
	return eg.err
}

func main() {
	group, groupDone := NewErrGroup()
	for i := 0; i < 5; i++ {
		group.Go(func() error {
			timeout := time.Second * time.Duration(rand.Intn(10))
			timer := time.NewTimer(timeout)
			defer timer.Stop()

			select {
			case <-groupDone:
				fmt.Println("canceled")
				return errors.New("error")
			case <-timer.C:
				fmt.Println("timeout")
				return errors.New("error")
			}
		})
	}

	if err := group.Wait(); err != nil {
		fmt.Println(err.Error())
	}
}
