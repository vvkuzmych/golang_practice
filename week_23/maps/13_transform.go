package main

import "fmt"

// 13. Transform - Transform map values

func transform[K comparable, V any, R any](m map[K]V, fn func(K, V) R) map[K]R {
	result := make(map[K]R, len(m))
	for k, v := range m {
		result[k] = fn(k, v)
	}
	return result
}

func main() {
	prices := map[string]float64{
		"apple":  1.0,
		"banana": 0.5,
		"cherry": 2.0,
	}

	// Apply 10% discount
	discounted := transform(prices, func(k string, v float64) float64 {
		return v * 0.9
	})
	fmt.Println("Discounted:", discounted)

	// Transform to strings
	asStrings := transform(prices, func(k string, v float64) string {
		return fmt.Sprintf("$%.2f", v)
	})
	fmt.Println("As strings:", asStrings)
}
