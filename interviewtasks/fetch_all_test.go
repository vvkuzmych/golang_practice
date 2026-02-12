package main

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
)

type MockFetcher struct {
	Responses map[string]string
	Errors    map[string]error
	mu        sync.Mutex
}

func (m *MockFetcher) Fetch(url string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err, ok := m.Errors[url]; ok {
		return "", err
	}

	if res, ok := m.Responses[url]; ok {
		return res, nil
	}

	return "", fmt.Errorf("unexpected url: %s", url)
}

func TestFetchAllInOrder(t *testing.T) {
	mock := &MockFetcher{
		Responses: map[string]string{
			"url1": "data1",
			"url2": "data2",
			"url3": "data3",
		},
	}

	urls := []string{"url1", "url2", "url3"}

	res, err := FetchAllInOrder(mock, urls)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"data1", "data2", "data3"}
	if !reflect.DeepEqual(res, expected) {
		t.Fatalf("expected %v, got %v", expected, res)
	}
}
