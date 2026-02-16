package main

import "fmt"

// 12. Filter - Filter map entries by condition

func filter[K comparable, V any](m map[K]V, predicate func(K, V) bool) map[K]V {
	result := make(map[K]V)
	for k, v := range m {
		if predicate(k, v) {
			result[k] = v
		}
	}
	return result
}

func main() {
	scores := map[string]int{
		"Alice": 85,
		"Bob":   45,
		"Carol": 92,
		"Dave":  58,
	}

	// Filter passing scores (>= 60)
	passed := filter(scores, func(k string, v int) bool { return v >= 60 })
	fmt.Println("Passed:", passed)

	// Filter by key prefix
	users := map[string]int{"user_1": 1, "user_2": 2, "admin_1": 3}
	usersOnly := filter(users, func(k string, _ int) bool {
		return len(k) >= 5 && k[:5] == "user_"
	})
	fmt.Println("Users only:", usersOnly)
}
