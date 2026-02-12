package main

import (
	"sync"
)

//import (
//	"net/http"
//	"sync"
//)
//
//type Fetcher interface {
//	Fetch(url string) (int, error)
//}
//
//// HTTPFetcher
//
//type HTTPFetcher struct{}
//
//func (f HTTPFetcher) Fetch(url string) (int, error) {
//	s, err := http.Get(url)
//
//	if err != nil {
//		return 0, err
//	}
//	return s.StatusCode, nil
//}
//
//// MockFetcher { Status int }
//
//type MockFetcher struct {
//	Status int
//}
//
//func (f MockFetcher) Fetch(url string) (int, error) {
//	return f.Status, nil
//}
//
//func main() {
//
//	httpFetcher := HTTPFetcher{}
//	FetchAll(httpFetcher, []string{"url1", "url2"})
//
//	mockFetcher := MockFetcher{Status: 200}
//	FetchAll(mockFetcher, []string{"url1", "url2"})
//}
//
//// - Fetch each URL with `http.Get`
//// - Return a slice of HTTP status codes in the same order
//// - Stop and return error immediately if any fetch fails
//func FetchAll(f Fetcher, urls []string) ([]int, error) {
//	statuses := make([]int, len(urls))
//	var wg sync.WaitGroup
//
//	wg.Add(len(urls))
//	for i, url := range urls {
//
//		// var s int
//		go func() {
//			s, err := checkStatus(f, url)
//			if err != nil {
//				//return s, err
//			}
//
//			statuses[i] = s
//
//			wg.Done()
//		}()
//	}
//
//	wg.Wait()
//	return statuses, nil
//}
//
//func checkStatus(f Fetcher, url string) (int, error) {
//	s, err := f.Fetch(url)
//
//	if err != nil {
//		return 0, err
//	}
//	return s, nil
//}

type URLFetcher interface {
	Fetch(url string) (string, error)
}

//type HTTPFetcher struct{}
//
//func (h *HTTPFetcher) Fetch(url string) (string, error) {
//	resp, err := http.Get(url)
//	if err != nil {
//		return "", err
//	}
//	defer resp.Body.Close()
//
//	body, err := io.ReadAll(resp.Body)
//	if err != nil {
//		return "", err
//	}
//
//	return string(body), nil
//}

func FetchAllInOrder(fetcher URLFetcher, urls []string) ([]string, error) {
	results := make([]string, len(urls))

	var (
		wg   sync.WaitGroup
		mu   sync.Mutex
		errs error
	)

	wg.Add(len(urls))

	for i, url := range urls {
		go func(idx int, u string) {
			defer wg.Done()

			data, err := fetcher.Fetch(u)
			if err != nil {
				mu.Lock()
				if errs == nil {
					errs = err
				}
				mu.Unlock()
				return
			}

			mu.Lock()
			results[idx] = data
			mu.Unlock()
		}(i, url)
	}

	wg.Wait()

	if errs != nil {
		return nil, errs
	}

	return results, nil
}
