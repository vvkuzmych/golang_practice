package main

import (
	"context"
	"fmt"
	"time"
)

// 19. Error Group Pattern - Обробка помилок у горутинах

type ErrorGroup struct {
	ctx    context.Context
	cancel context.CancelFunc
	errCh  chan error
	count  int
}

func NewErrorGroup(ctx context.Context) *ErrorGroup {
	ctx, cancel := context.WithCancel(ctx)
	return &ErrorGroup{
		ctx:    ctx,
		cancel: cancel,
		errCh:  make(chan error, 1),
	}
}

func (eg *ErrorGroup) Go(f func() error) {
	eg.count++
	go func() {
		if err := f(); err != nil {
			select {
			case eg.errCh <- err:
				eg.cancel() // Cancel інші горутини
			default:
			}
		}
	}()
}

func (eg *ErrorGroup) Wait() error {
	// Чекаємо завершення або помилки
	time.Sleep(100 * time.Millisecond)
	select {
	case err := <-eg.errCh:
		return err
	default:
		return nil
	}
}

func task1() error {
	time.Sleep(50 * time.Millisecond)
	fmt.Println("Task 1 completed")
	return nil
}

func task2() error {
	time.Sleep(100 * time.Millisecond)
	return fmt.Errorf("task 2 failed")
}

func task3() error {
	time.Sleep(150 * time.Millisecond)
	fmt.Println("Task 3 completed")
	return nil
}

func main() {
	eg := NewErrorGroup(context.Background())

	eg.Go(task1)
	eg.Go(task2)
	eg.Go(task3)

	if err := eg.Wait(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
