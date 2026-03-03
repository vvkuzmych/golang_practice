package main

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"time"
)

// URLSorter interface - dependency injection pattern
type URLSorter interface {
	Sort(ctx context.Context, urls []string) ([]string, error)
}

// DefaultURLSorter - concrete implementation
type DefaultURLSorter struct{}

func (s *DefaultURLSorter) Sort(ctx context.Context, urls []string) ([]string, error) {
	sorted := make([]string, len(urls))
	copy(sorted, urls)

	// Check for cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Sort by hostname
	sort.Slice(sorted, func(i, j int) bool {
		ui, _ := url.Parse("http://" + sorted[i])
		uj, _ := url.Parse("http://" + sorted[j])
		return ui.Host < uj.Host
	})

	return sorted, nil
}

// URLService - service that uses URLSorter
type URLService struct {
	sorter URLSorter
}

func NewURLService(s URLSorter) *URLService {
	return &URLService{sorter: s}
}

// Result holds async operation result
type Result struct {
	URLs []string
	Err  error
}

// ProcessAsync - asynchronous processing with channels
func (s *URLService) ProcessAsync(ctx context.Context, urls []string) <-chan Result {
	resultCh := make(chan Result, 1)

	go func() {
		defer close(resultCh)

		res, err := s.sorter.Sort(ctx, urls)
		resultCh <- Result{
			URLs: res,
			Err:  err,
		}
	}()

	return resultCh
}

// ProcessWithCallback - callback-based async processing
func (s *URLService) ProcessWithCallback(
	ctx context.Context,
	urls []string,
	callback func([]string, error),
) {
	go func() {
		res, err := s.sorter.Sort(ctx, urls)
		callback(res, err)
	}()
}

func main() {
	fmt.Println("=== URL Sorting Service Demo ===")
	fmt.Println()

	// Example 1: Using default sorter
	fmt.Println("1. Default Sorter:")
	defaultSorter := &DefaultURLSorter{}
	svc := NewURLService(defaultSorter)

	ctx := context.Background()
	resultCh := svc.ProcessAsync(ctx, []string{
		"github.com",
		"stackoverflow.com",
		"google.com",
		"amazon.com",
	})

	result := <-resultCh
	if result.Err != nil {
		fmt.Printf("Error: %v\n", result.Err)
	} else {
		fmt.Printf("Sorted URLs: %v\n\n", result.URLs)
	}

	// Example 2: Using mock sorter
	fmt.Println("2. Mock Sorter:")
	mockSorter := &MockURLSorter{
		SortFunc: func(ctx context.Context, urls []string) ([]string, error) {
			// Reverse order as mock behavior
			reversed := make([]string, len(urls))
			for i, url := range urls {
				reversed[len(urls)-1-i] = url
			}
			return reversed, nil
		},
	}

	svcMock := NewURLService(mockSorter)
	resultCh2 := svcMock.ProcessAsync(ctx, []string{"a.com", "b.com", "c.com"})
	result2 := <-resultCh2
	fmt.Printf("Mock result: %v\n\n", result2.URLs)

	// Example 3: Context cancellation
	fmt.Println("3. Context Cancellation:")
	ctxCancel, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	time.Sleep(10 * time.Millisecond) // Ensure timeout
	resultCh3 := svc.ProcessAsync(ctxCancel, []string{"x.com", "y.com"})
	result3 := <-resultCh3
	if result3.Err != nil {
		fmt.Printf("Expected error: %v\n\n", result3.Err)
	}

	// Example 4: Callback pattern
	fmt.Println("4. Callback Pattern:")
	done := make(chan bool)
	svc.ProcessWithCallback(ctx, []string{"example.com", "test.com"}, func(urls []string, err error) {
		if err != nil {
			fmt.Printf("Callback error: %v\n", err)
		} else {
			fmt.Printf("Callback result: %v\n", urls)
		}
		done <- true
	})
	<-done

	fmt.Println()
	fmt.Println("=== Demo Complete ===")
}
