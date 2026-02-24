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
