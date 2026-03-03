package main

import "context"

// MockURLSorter - test double for URLSorter interface
// Used in tests to inject custom behavior without real sorting logic
type MockURLSorter struct {
	SortFunc func(ctx context.Context, urls []string) ([]string, error)
}

// Sort implements URLSorter interface by delegating to injected function
func (m *MockURLSorter) Sort(ctx context.Context, urls []string) ([]string, error) {
	if m.SortFunc == nil {
		// Default behavior if no function provided
		return urls, nil
	}
	return m.SortFunc(ctx, urls)
}
