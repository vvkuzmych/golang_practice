package main

import (
	"context"
	"errors"
	"testing"
	"time"
)

// TestURLService_ProcessAsync_Success tests successful sorting
func TestURLService_ProcessAsync_Success(t *testing.T) {
	// Arrange: Create mock with successful behavior
	mock := &MockURLSorter{
		SortFunc: func(ctx context.Context, urls []string) ([]string, error) {
			return []string{"sorted.com", "urls.com"}, nil
		},
	}
	service := NewURLService(mock)

	// Act: Process URLs asynchronously
	resultCh := service.ProcessAsync(context.Background(), []string{"input.com"})
	result := <-resultCh

	// Assert: Check results
	if result.Err != nil {
		t.Errorf("Expected no error, got: %v", result.Err)
	}
	if len(result.URLs) != 2 {
		t.Errorf("Expected 2 URLs, got: %d", len(result.URLs))
	}
	if result.URLs[0] != "sorted.com" {
		t.Errorf("Expected first URL to be 'sorted.com', got: %s", result.URLs[0])
	}
}

// TestURLService_ProcessAsync_Error tests error handling
func TestURLService_ProcessAsync_Error(t *testing.T) {
	// Arrange: Create mock that returns error
	expectedErr := errors.New("sorting failed")
	mock := &MockURLSorter{
		SortFunc: func(ctx context.Context, urls []string) ([]string, error) {
			return nil, expectedErr
		},
	}
	service := NewURLService(mock)

	// Act: Process URLs
	resultCh := service.ProcessAsync(context.Background(), []string{"input.com"})
	result := <-resultCh

	// Assert: Error is propagated
	if result.Err == nil {
		t.Error("Expected error, got nil")
	}
	if result.Err.Error() != expectedErr.Error() {
		t.Errorf("Expected error '%v', got: '%v'", expectedErr, result.Err)
	}
}

// TestURLService_ProcessAsync_ContextCancellation tests context cancellation
func TestURLService_ProcessAsync_ContextCancellation(t *testing.T) {
	// Arrange: Create mock with delay
	mock := &MockURLSorter{
		SortFunc: func(ctx context.Context, urls []string) ([]string, error) {
			// Simulate slow operation
			select {
			case <-time.After(1 * time.Second):
				return []string{"slow.com"}, nil
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		},
	}
	service := NewURLService(mock)

	// Act: Use context with short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	resultCh := service.ProcessAsync(ctx, []string{"input.com"})
	result := <-resultCh

	// Assert: Context cancellation error
	if result.Err == nil {
		t.Error("Expected context cancellation error, got nil")
	}
	if !errors.Is(result.Err, context.DeadlineExceeded) {
		t.Errorf("Expected DeadlineExceeded error, got: %v", result.Err)
	}
}

// TestURLService_ProcessWithCallback tests callback pattern
func TestURLService_ProcessWithCallback(t *testing.T) {
	// Arrange: Create mock
	mock := &MockURLSorter{
		SortFunc: func(ctx context.Context, urls []string) ([]string, error) {
			return []string{"callback.com"}, nil
		},
	}
	service := NewURLService(mock)

	// Act: Process with callback
	done := make(chan bool)
	var callbackResult []string
	var callbackErr error

	service.ProcessWithCallback(
		context.Background(),
		[]string{"input.com"},
		func(urls []string, err error) {
			callbackResult = urls
			callbackErr = err
			done <- true
		},
	)

	// Wait for callback
	<-done

	// Assert: Callback received correct data
	if callbackErr != nil {
		t.Errorf("Expected no error, got: %v", callbackErr)
	}
	if len(callbackResult) != 1 || callbackResult[0] != "callback.com" {
		t.Errorf("Expected ['callback.com'], got: %v", callbackResult)
	}
}

// TestMockURLSorter_NilFunction tests default behavior when SortFunc is nil
func TestMockURLSorter_NilFunction(t *testing.T) {
	// Arrange: Mock without SortFunc
	mock := &MockURLSorter{}

	// Act: Call Sort
	result, err := mock.Sort(context.Background(), []string{"a.com", "b.com"})

	// Assert: Returns input unchanged
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("Expected 2 URLs, got: %d", len(result))
	}
}

// BenchmarkURLService_ProcessAsync benchmarks async processing
func BenchmarkURLService_ProcessAsync(b *testing.B) {
	mock := &MockURLSorter{
		SortFunc: func(ctx context.Context, urls []string) ([]string, error) {
			return urls, nil
		},
	}
	service := NewURLService(mock)
	ctx := context.Background()
	urls := []string{"test.com"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resultCh := service.ProcessAsync(ctx, urls)
		<-resultCh
	}
}
