package main

import (
	"context"
	"fmt"
	"time"
)

type DefaultSortingStruct struct{}

func (s *DefaultSortingStruct) Swaper(strings []string) []string {
	return strings
}

type SortInterface interface {
	Swaper(urls []string) []string
	Swap(urls []string) []string
	//Less([]string) []string
}

//func (s *DefaultSortingStruct) Less([]string) []string {
//	return []string{}
//}

func (s *DefaultSortingStruct) Swap(urls []string) []string {
	return urls
}

type MockSortingStruct struct {
	SwapFunc func(urls []string) []string
	//CallCount int
	//LastURLs  []string
}

func (m *MockSortingStruct) Swap(urls []string) []string {
	return urls
}

// Swap реалізує SortInterface (БЕЗ context.Context!)
func (m *MockSortingStruct) Swaper(urls []string) []string {
	//m.CallCount++
	//m.LastURLs = append([]string{}, urls...)

	if m.SwapFunc == nil {
		return urls
	}
	return m.SwapFunc(urls)
}

type SorterService struct {
	sorter SortInterface
}

func NewSorterService(s SortInterface) *SorterService {
	return &SorterService{sorter: s}
}

func (newSorter *SorterService) SortMethod(ctx context.Context, urls []string) <-chan []string {
	results := make(chan []string)

	go func() {
		defer close(results)
		select {
		case <-ctx.Done():
			return
		default:
		}

		res := newSorter.sorter.Swaper(urls)

		select {
		case results <- res:
		case <-ctx.Done():
		}
	}()
	return results
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	urls := []string{
		"https://www.baidu1.com.cn/",
		"https://www.baidu2.com.cn/",
		"https://www.baidu3.com.cn/",
	}

	//// Example 1: DefaultSortingStruct
	fmt.Println("=== Example 1: DefaultSortingStruct ===")
	s := &DefaultSortingStruct{}
	newSorter := NewSorterService(s)
	results := newSorter.SortMethod(ctx, urls)

	select {
	case url := <-results:
		fmt.Println("Result:", url)
	case <-ctx.Done():
		fmt.Println("Timeout або cancellation")
	}

	fmt.Println()

	// Example 2: MockSortingStruct з reverse
	fmt.Println("=== Example 2: MockSortingStruct (reverse) ===")
	m := &MockSortingStruct{
		SwapFunc: func(urls []string) []string {
			reversed := make([]string, len(urls))
			for i, url := range urls {
				reversed[len(urls)-1-i] = url
			}
			return reversed
		},
	}
	mockSorter := NewSorterService(m)
	res := mockSorter.SortMethod(ctx, urls)

	select {
	case url := <-res:
		fmt.Println("Mock result (reversed):", url)
		//fmt.Printf("CallCount: %d\n", m.CallCount)
	case <-ctx.Done():
		fmt.Println("Timeout або cancellation")
	}

}
