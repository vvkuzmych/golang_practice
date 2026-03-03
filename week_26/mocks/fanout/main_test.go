package main

import (
	"sync/atomic"
	"testing"
	"time"
)

// TestFanOut_BasicFunctionality tests fan-out with multiple workers
func TestFanOut_BasicFunctionality(t *testing.T) {
	var callCount int32
	process := func(n int) int {
		atomic.AddInt32(&callCount, 1)
		time.Sleep(10 * time.Millisecond)
		return n * 2
	}

	numbers := []int{1, 2, 3, 4, 5, 6}
	results := ProcessFanOut(numbers, process, 3)

	if len(results) != len(numbers) {
		t.Errorf("Expected %d results, got %d", len(numbers), len(results))
	}

	if atomic.LoadInt32(&callCount) != int32(len(numbers)) {
		t.Errorf("Expected %d calls, got %d", len(numbers), callCount)
	}

	for i, result := range results {
		expected := numbers[i] * 2
		if result != expected {
			t.Errorf("Expected %d, got %d", expected, result)
		}
	}
}

// TestFanOut_PerformanceImprovement tests that fan-out is faster
func TestFanOut_PerformanceImprovement(t *testing.T) {
	process := func(n int) int {
		time.Sleep(100 * time.Millisecond)
		return n * 2
	}

	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8}

	// Sequential
	startSeq := time.Now()
	_ = ProcessSequential(numbers, process)
	durationSeq := time.Since(startSeq)

	// Fan-out
	startFan := time.Now()
	_ = ProcessFanOut(numbers, process, 4)
	durationFan := time.Since(startFan)

	if durationFan >= durationSeq {
		t.Errorf("Fan-out (%v) should be faster than sequential (%v)", durationFan, durationSeq)
	}

	speedup := float64(durationSeq) / float64(durationFan)
	t.Logf("Speedup: %.2fx (Sequential: %v, Fan-out: %v)", speedup, durationSeq, durationFan)

	if speedup < 2.0 {
		t.Errorf("Expected at least 2x speedup, got %.2fx", speedup)
	}
}

// TestFanOut_OrderPreserved tests that output order matches input order
func TestFanOut_OrderPreserved(t *testing.T) {
	process := func(n int) int {
		// Add variable delay to test ordering
		if n%2 == 0 {
			time.Sleep(50 * time.Millisecond)
		} else {
			time.Sleep(10 * time.Millisecond)
		}
		return n * 2
	}

	numbers := []int{1, 2, 3, 4, 5, 6}
	results := ProcessFanOut(numbers, process, 3)

	for i, result := range results {
		expected := numbers[i] * 2
		if result != expected {
			t.Errorf("Order not preserved at index %d: expected %d, got %d", i, expected, result)
		}
	}
}

// TestFanOut_DifferentFunctions tests using different processing functions
func TestFanOut_DifferentFunctions(t *testing.T) {
	tests := []struct {
		name     string
		process  ProcessFunc
		input    []int
		expected []int
	}{
		{
			name:     "Double",
			process:  func(n int) int { return n * 2 },
			input:    []int{1, 2, 3},
			expected: []int{2, 4, 6},
		},
		{
			name:     "Square",
			process:  func(n int) int { return n * n },
			input:    []int{2, 3, 4},
			expected: []int{4, 9, 16},
		},
		{
			name:     "AddTen",
			process:  func(n int) int { return n + 10 },
			input:    []int{5, 10, 15},
			expected: []int{15, 20, 25},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := ProcessFanOut(tt.input, tt.process, 2)

			if len(results) != len(tt.expected) {
				t.Fatalf("Expected %d results, got %d", len(tt.expected), len(results))
			}

			for i, result := range results {
				if result != tt.expected[i] {
					t.Errorf("Index %d: expected %d, got %d", i, tt.expected[i], result)
				}
			}
		})
	}
}

// BenchmarkFanOut compares sequential vs fan-out
func BenchmarkFanOut(b *testing.B) {
	process := func(n int) int {
		time.Sleep(10 * time.Millisecond)
		return n * 2
	}

	numbers := []int{1, 2, 3, 4}

	b.Run("Sequential", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = ProcessSequential(numbers, process)
		}
	})

	b.Run("FanOut_2Workers", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = ProcessFanOut(numbers, process, 2)
		}
	})

	b.Run("FanOut_4Workers", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = ProcessFanOut(numbers, process, 4)
		}
	})
}
