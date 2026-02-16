// package main
//
// import (
//
//	"fmt"
//	"sync"
//	"time"
//
// )
//
// // 15. Rate Limiter Pattern - Обмеження швидкості виконання
//
// // Паттерн: Контроль частоти виконання операцій
// //
// // Requests → Rate Limiter (N per second) → Allowed Requests
//
// // Token Bucket Rate Limiter
//
//	type RateLimiter struct {
//		tokens chan struct{}
//		rate   time.Duration
//	}
//
//	func NewRateLimiter(maxRequests int, interval time.Duration) *RateLimiter {
//		rl := &RateLimiter{
//			tokens: make(chan struct{}, maxRequests),
//			rate:   interval / time.Duration(maxRequests),
//		}
//
//		// Заповнюємо bucket токенами
//		for i := 0; i < maxRequests; i++ {
//			rl.tokens <- struct{}{}
//		}
//
//		// Поповнюємо токени періодично
//		go func() {
//			ticker := time.NewTicker(rl.rate)
//			defer ticker.Stop()
//
//			for range ticker.C {
//				select {
//				case rl.tokens <- struct{}{}:
//				default:
//					// Bucket повний
//				}
//			}
//		}()
//
//		return rl
//	}
//
//	func (rl *RateLimiter) Allow() bool {
//		select {
//		case <-rl.tokens:
//			return true
//		default:
//			return false
//		}
//	}
//
//	func (rl *RateLimiter) Wait() {
//		<-rl.tokens
//	}
//
// // Sliding Window Rate Limiter
//
//	type SlidingWindowLimiter struct {
//		mu       sync.Mutex
//		requests []time.Time
//		limit    int
//		window   time.Duration
//	}
//
//	func NewSlidingWindowLimiter(limit int, window time.Duration) *SlidingWindowLimiter {
//		return &SlidingWindowLimiter{
//			requests: make([]time.Time, 0),
//			limit:    limit,
//			window:   window,
//		}
//	}
//
//	func (swl *SlidingWindowLimiter) Allow() bool {
//		swl.mu.Lock()
//		defer swl.mu.Unlock()
//
//		now := time.Now()
//
//		// Видаляємо старі requests
//		cutoff := now.Add(-swl.window)
//		newRequests := make([]time.Time, 0)
//		for _, t := range swl.requests {
//			if t.After(cutoff) {
//				newRequests = append(newRequests, t)
//			}
//		}
//		swl.requests = newRequests
//
//		// Перевіряємо ліміт
//		if len(swl.requests) < swl.limit {
//			swl.requests = append(swl.requests, now)
//			return true
//		}
//
//		return false
//	}
//
//	func main() {
//		fmt.Println("=== Rate Limiter Pattern ===")
//		fmt.Println()
//
//		// 1. Token Bucket Limiter
//		fmt.Println("Token Bucket: 5 requests per second")
//		fmt.Println()
//
//		limiter := NewRateLimiter(5, time.Second)
//
//		// Відправляємо 10 requests
//		for i := 1; i <= 10; i++ {
//			if limiter.Allow() {
//				fmt.Printf("Request %d: ✅ Allowed\n", i)
//			} else {
//				fmt.Printf("Request %d: ❌ Rate limited\n", i)
//			}
//		}
//
//		// Чекаємо поповнення
//		fmt.Println()
//		fmt.Println("Waiting 1 second for token refill...")
//		time.Sleep(1 * time.Second)
//
//		for i := 11; i <= 15; i++ {
//			if limiter.Allow() {
//				fmt.Printf("Request %d: ✅ Allowed\n", i)
//			} else {
//				fmt.Printf("Request %d: ❌ Rate limited\n", i)
//			}
//		}
//
//		// 2. Sliding Window Limiter
//		fmt.Println()
//		fmt.Println("=== Sliding Window: 3 requests per 2 seconds ===")
//		fmt.Println()
//
//		swLimiter := NewSlidingWindowLimiter(3, 2*time.Second)
//
//		for i := 1; i <= 8; i++ {
//			if swLimiter.Allow() {
//				fmt.Printf("Request %d: ✅ Allowed at %v\n", i, time.Now().Format("15:04:05"))
//			} else {
//				fmt.Printf("Request %d: ❌ Rate limited at %v\n", i, time.Now().Format("15:04:05"))
//			}
//			time.Sleep(500 * time.Millisecond)
//		}
//
//		// 3. Rate Limiter з Wait
//		fmt.Println()
//		fmt.Println("=== With Wait (blocking) ===")
//		fmt.Println()
//
//		limiter2 := NewRateLimiter(2, time.Second)
//		var wg sync.WaitGroup
//
//		for i := 1; i <= 5; i++ {
//			wg.Add(1)
//			go func(n int) {
//				defer wg.Done()
//				limiter2.Wait() // Чекаємо токен
//				fmt.Printf("Request %d: executed at %v\n", n, time.Now().Format("15:04:05"))
//			}(i)
//		}
//
//		wg.Wait()
//
//		fmt.Println()
//		fmt.Println("✅ Rate limiting complete")
//	}
//
// // Use cases:
// // - API rate limiting
// // - Database query throttling
// // - Request throttling
// // - Resource protection
// // - DoS prevention
package main

import (
	"time"
)

type RateLimiter struct {
	leakyBucketCh chan struct{}

	closeCh     chan struct{}
	closeDoneCh chan struct{}
}

func NewLeakyBucketLimiter(limit int, period time.Duration) RateLimiter {
	limiter := RateLimiter{
		leakyBucketCh: make(chan struct{}, limit),
		closeCh:       make(chan struct{}),
		closeDoneCh:   make(chan struct{}),
	}

	leakInterval := period.Nanoseconds() / int64(limit)
	go limiter.startPeriodicLeak(time.Duration(leakInterval))
	return limiter
}

func (l *RateLimiter) startPeriodicLeak(interval time.Duration) {
	timer := time.NewTicker(interval)
	defer func() {
		timer.Stop()
		close(l.closeDoneCh)
	}()

	for {
		select {
		case <-l.closeCh:
			return
		default:
		}

		select {
		case <-l.closeCh:
			return
		case <-timer.C:
			select {
			case <-l.leakyBucketCh:
			default:
			}
		}
	}
}

func (l *RateLimiter) Allow() bool {
	select {
	case l.leakyBucketCh <- struct{}{}:
		return true
	default:
		return false
	}
}

func (l *RateLimiter) Shutdown() {
	close(l.closeCh)
	<-l.closeDoneCh
}
