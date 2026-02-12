package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func reverse_string(str string) string {
	if str == "" {
		return ""
	}

	runes := []rune(str)
	if len(runes) < 2 {
		return str
	}

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	fmt.Println("-------", string(runes))
	return string(runes)
}

func find_duplicates(arr []int) []int {
	if len(arr) < 2 {
		return []int{}
	}

	seen := make(map[int]bool)
	added := make(map[int]bool)
	result := []int{}

	for _, v := range arr {
		if seen[v] && !added[v] {
			result = append(result, v)
			added[v] = true
		}
		seen[v] = true
	}

	return result
}

func fizzbuzz(n int) []any {
	if n == 0 {
		return []any{}
	}
	if n < 2 {
		return []any{1}
	}
	result := []any{}
	for i := 1; i <= n; i++ {
		if i%3 == 0 && i%5 == 0 {
			result = append(result, "FIzzBuzz")
		}
		if i%3 == 0 && i%5 != 0 {
			result = append(result, "Fizz")
		}
		if i%5 == 0 && i%3 != 0 {
			result = append(result, "Buzz")
		}
		if i%3 != 0 && i%5 != 0 {
			result = append(result, i)
		}
	}
	return result
}

func is_palindrome(str string) bool {
	if str == "" {
		return true
	}
	s := strings.ReplaceAll(str, " ", "")
	s = strings.ToLower(s)
	runes := []rune(s)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	fmt.Println("-------", string(runes))
	return string(runes) == s
}

func sumArray1(arr []interface{}) int {
	sum := 0

	for _, item := range arr {
		switch v := item.(type) {
		case int:
			// Якщо елемент - число, додаємо до суми
			sum += v
		case []interface{}:
			// Якщо елемент - масив, рекурсивно обчислюємо його суму
			sum += sumArray1(v)
		default:
			// Ігноруємо інші типи (strings, nil, etc.)
			continue
		}
	}

	return sum
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	_ = name
	if name == "" {
		name = "Guest"
	}
	fmt.Fprintf(w, "Hello, %s!\n", name)
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		log.Printf("error - %v", err)
	}
}

//func main() {
//	//a := reverse_string("a!b@c#")
//	//fmt.Println(a)
//
//	//a := find_duplicates([]int{
//	//	1, 2, 3, 2, 4, 5, 1,
//	//})
//	//a := find_duplicates([]int{
//	//	5, 5, 3, 3, 1, 1,
//	//})
//	//a := fizzbuzz(0)
//	//a := is_palindrome("racecar")
//	//fmt.Println(a)
//	//b := is_palindrome("A man a plan a canal Panama")
//	//b := is_palindrome("race a car")
//	//b := sumArray1([]interface{}{1, []interface{}{2, 3}, 4, []interface{}{5, 6}})
//	//fmt.Println(b)
//	//arr2 := sumArray1([]interface{}{1, []interface{}{2, 3}, 4, []interface{}{5, 6}})
//	////fmt.Printf("sumArray(%v) = %d\n", arr2, sumArray(arr2))
//	//fmt.Println(arr2)
//
//}
