// package main
//
// import (
//
//	"context"
//	"fmt"
//	"sync"
//	"time"
//
// )
//
// // 13. ErrGroup Pattern - Групова обробка з помилками
//
// // Паттерн: Запуск горутин з обробкою першої помилки
// //
// // Task1 ──┐
// // Task2 ──┼→ ErrGroup → First Error → Cancel All
// // Task3 ──┘
//
//	type ErrGroup struct {
//		ctx    context.Context
//		cancel context.CancelFunc
//		wg     *sync.WaitGroup
//		errCh  chan error
//	}
//
//	func NewErrGroup(ctx context.Context) *ErrGroup {
//		ctx, cancel := context.WithCancel(ctx)
//		return &ErrGroup{
//			ctx:    ctx,
//			cancel: cancel,
//			wg:     &sync.WaitGroup{},
//			errCh:  make(chan error, 1),
//		}
//	}
//
//	func (eg *ErrGroup) Go(f func(ctx context.Context) error) {
//		eg.wg.Add(1)
//		go func() {
//			defer eg.wg.Done()
//
//			if err := f(eg.ctx); err != nil {
//				select {
//				case eg.errCh <- err:
//					eg.cancel() // Скасовуємо інші
//				default:
//				}
//			}
//		}()
//	}
//
//	func (eg *ErrGroup) Wait() error {
//		// Чекаємо завершення всіх горутин
//		done := make(chan struct{})
//		go func() {
//			eg.wg.Wait()
//			close(done)
//		}()
//
//		select {
//		case <-done:
//			// Всі завершились без помилок
//			return nil
//		case err := <-eg.errCh:
//			// Отримали помилку
//			return err
//		}
//	}
//
// // Робочі функції
//
//	func task1(ctx context.Context) error {
//		for i := 0; i < 3; i++ {
//			select {
//			case <-ctx.Done():
//				fmt.Println("Task1: cancelled")
//				return ctx.Err()
//			default:
//				fmt.Println("Task1: working", i)
//				time.Sleep(300 * time.Millisecond)
//			}
//		}
//		fmt.Println("Task1: completed")
//		return nil
//	}
//
//	func task2(ctx context.Context) error {
//		time.Sleep(500 * time.Millisecond)
//		fmt.Println("Task2: ERROR occurred!")
//		return fmt.Errorf("task2 failed")
//	}
//
//	func task3(ctx context.Context) error {
//		for i := 0; i < 5; i++ {
//			select {
//			case <-ctx.Done():
//				fmt.Println("Task3: cancelled")
//				return ctx.Err()
//			default:
//				fmt.Println("Task3: working", i)
//				time.Sleep(200 * time.Millisecond)
//			}
//		}
//		fmt.Println("Task3: completed")
//		return nil
//	}
//
//	func main() {
//		fmt.Println("=== ErrGroup Pattern ===")
//		fmt.Println()
//
//		eg := NewErrGroup(context.Background())
//
//		// Запускаємо 3 tasks
//		eg.Go(task1)
//		eg.Go(task2) // Ця задача зафейлиться
//		eg.Go(task3)
//
//		// Чекаємо завершення або помилки
//		if err := eg.Wait(); err != nil {
//			fmt.Println()
//			fmt.Printf("❌ Error occurred: %v\n", err)
//			fmt.Println("   Other tasks were cancelled")
//		} else {
//			fmt.Println()
//			fmt.Println("✅ All tasks completed successfully")
//		}
//	}
//
// // Use cases:
// // - Parallel API calls (stop on first error)
// // - Database migrations
// // - File processing
// // - Health checks
// // - Validation tasks
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
