package main

//
//import (
//	"fmt"
//)
//
//type DefaultSortingStruct struct{}
//
//type SortInterface interface {
//	Swap([]string) []string
//}
//
//func (s *DefaultSortingStruct) Swap(urls []string) []string {
//	return urls
//}
//
//type SorterService struct {
//	sorter SortInterface
//}
//
//func NewSorterService(s SortInterface) *SorterService {
//	return &SorterService{sorter: s}
//}
//
//func (newSorter *SorterService) SortMethod(urls []string) <-chan []string {
//	results := make(chan []string)
//
//	go func() {
//		defer close(results)
//
//		res := newSorter.sorter.Swap(urls)
//		results <- res
//	}()
//	return results
//}
//
//func main() {
//	urls := []string{
//		"https://www.baidu1.com.cn/",
//		"https://www.baidu2.com.cn/",
//		"https://www.baidu3.com.cn/",
//	}
//
//	s := &DefaultSortingStruct{}
//	newSorter := NewSorterService(s)
//
//	results := newSorter.SortMethod(urls)
//
//	// ✅ ВИПРАВЛЕНО: Видалено default - тепер це блокуючий receive
//	// Програма чекатиме поки goroutine надішле дані
//	url := <-results
//	fmt.Println("Отримані URLs:", url)
//}
